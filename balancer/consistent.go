package balancer

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/naming"

	lgrpc "github.com/binchencoder/letsgo/grpc"
	"github.com/binchencoder/letsgo/hashring"
	"github.com/binchencoder/skylb-api/internal/health"
)

// consistentHashing implements interface grpc.Balancer. It maintains a
// hashring which implements consistent hashing algorithm. When a gRPC call
// is issued through the balancer, it take out a key from the context and
// pick a proper address in the ring and dispatch the request to that address.
//
// gRPC client should use the context functions provided in letsgo/hashring
// package to create/put/extract hash key from the context.
//
// This balancer is mostly a copy of the grpc RoundRobin balancer, with the
// backing structure replaced with a HashRing.
type consistentHashing struct {
	r          naming.Resolver
	w          naming.Watcher
	hr         *hashring.HashRing
	mu         sync.Mutex
	nextHealth int
	// The channel to notify gRPC internals the list of addresses.
	// The client should connect to.
	addrCh chan []grpc.Address
	// The channel to block when there is no connected address available.
	waitCh chan struct{}
	// The Balancer is closed.
	done bool

	debugPrintStopCh chan struct{}
}

// ConsistentHashing returns a Balancer that selects addresses based on
// consistent hashing algorithm. It uses r to watch the name resolution
// updates and updates the addresses available correspondingly.
func ConsistentHashing(r naming.Resolver) grpc.Balancer {
	return &consistentHashing{
		r:                r,
		hr:               hashring.New(),
		debugPrintStopCh: make(chan struct{}, 1),
	}
}

func (conh *consistentHashing) watchAddrUpdates() error {
	updates, err := conh.w.Next()
	if err != nil {
		glog.Errorf("grpc: the naming watcher stops working due to %+v.\n", err)
		return err
	}
	conh.mu.Lock()
	defer conh.mu.Unlock()
	for _, update := range updates {
		addr := grpc.Address{
			Addr:     update.Addr,
			Metadata: update.Metadata,
		}
		switch update.Op {
		case naming.Add:
			if conh.hr.GetMember(addr.Addr) != nil {
				glog.Errorln("Try to add existing endpoint:", addr.Addr)
				continue
			}
			member := hashring.Member{
				Key: addr.Addr,
				Val: &addrInfo{
					addr: addr,
				},
			}
			conh.hr.Add(&member)
		case naming.Delete:
			if conh.hr.GetMember(addr.Addr) == nil {
				glog.Errorln("Try to remove non-existing endpoint:", addr.Addr)
				continue
			}
			conh.hr.Remove(addr.Addr)
		default:
			glog.Errorln("Unknown update.Op ", update.Op)
		}
	}
	// Make a copy of conh.addrs and write it to conh.addrCh so that
	// gRPC internals gets notified.
	m := conh.hr.Members()
	open := make([]grpc.Address, len(m))
	for i, v := range m {
		ai := v.Val.(*addrInfo)
		open[i] = ai.addr
	}
	if conh.done {
		return grpc.ErrClientConnClosing
	}
	conh.addrCh <- open
	return nil
}

func (conh *consistentHashing) Start(target string, config grpc.BalancerConfig) error {
	conh.mu.Lock()
	defer conh.mu.Unlock()
	if conh.done {
		return grpc.ErrClientConnClosing
	}
	if conh.r == nil {
		// If there is no name resolver installed, it is not needed to
		// do name resolution. In this case, target is added into conh.addrs
		// as the only address available and conh.addrCh stays nil.
		member := hashring.Member{
			Key: target,
			Val: &addrInfo{
				addr: grpc.Address{
					Addr: target,
				},
			},
		}
		conh.hr.Add(&member)
		return nil
	}
	w, err := conh.r.Resolve(target)
	if err != nil {
		return err
	}
	conh.w = w
	conh.addrCh = make(chan []grpc.Address)
	go func() {
		for {
			if err := conh.watchAddrUpdates(); err != nil {
				return
			}
		}
	}()
	return nil
}

// Up sets the connected state of addr and sends notification if there are
// pending Get() calls.
func (conh *consistentHashing) Up(addr grpc.Address) func(error) {
	conh.mu.Lock()
	defer conh.mu.Unlock()

	var cnt int
	for _, m := range conh.hr.Members() {
		ai := m.Val.(*addrInfo)
		if ai.addr == addr {
			if ai.connected {
				return nil
			}
			ai.connected = true
		}
		if ai.connected {
			cnt++
		}
	}

	// addr is only one which is connected. Notify the Get() callers
	// who are blocking.
	if cnt == 1 && conh.waitCh != nil {
		close(conh.waitCh)
		conh.waitCh = nil
	}
	return func(err error) {
		conh.down(addr, err)
	}
}

// down unsets the connected state of addr.
func (conh *consistentHashing) down(addr grpc.Address, err error) {
	conh.mu.Lock()
	defer conh.mu.Unlock()

	m := conh.hr.GetMember(addr.Addr)
	if m != nil {
		ai := m.Val.(*addrInfo)
		ai.connected = false
	}
}

// Get returns the correct addr in the hash ring.
func (conh *consistentHashing) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	// Try to get hash key from context.
	hashkey, ok := hashring.GetHashKey(ctx)
	if !ok {
		// Try to get hash key from metadata.
		incoming, _ := lgrpc.FromMetadataIncoming(ctx)
		hashkey, _ = hashring.GetHashKey(incoming)
	}

	conh.mu.Lock()
	if conh.done {
		conh.mu.Unlock()
		err = grpc.ErrClientConnClosing
		return
	}

	if health.IsHealthCheck(ctx) {
		members := conh.hr.Members()
		if conh.nextHealth >= len(members) {
			conh.nextHealth = 0
		}
		next := conh.nextHealth
		addr = members[next].Val.(*addrInfo).addr
		// Because nextHealth checks the length above,
		// there is no need to take % here.
		conh.nextHealth += 1
		conh.mu.Unlock()
		return
	}

	next, err := conh.hr.Get(hashkey)
	if err != nil {
		conh.mu.Unlock()
		return
	}

	m := conh.hr.GetMember(next)
	if m != nil {
		ai := m.Val.(*addrInfo)
		if ai.connected {
			addr = ai.addr
			conh.mu.Unlock()
			return
		}
		err = grpc.Errorf(codes.Unavailable, "(consistent-hashing) hash mapped instance is down")
		conh.mu.Unlock()
		return
	}

	conh.mu.Unlock()
	err = grpc.Errorf(codes.Unavailable, "(consistent-hashing) this shouldn't happen, hash key: %s, next: %s.", hashkey, next)
	debug.PrintStack()
	return
}

func (conh *consistentHashing) Notify() <-chan []grpc.Address {
	return conh.addrCh
}

func (conh *consistentHashing) Close() error {
	conh.mu.Lock()
	defer conh.mu.Unlock()
	conh.done = true
	if conh.w != nil {
		conh.w.Close()
	}
	if conh.waitCh != nil {
		close(conh.waitCh)
		conh.waitCh = nil
	}
	if conh.addrCh != nil {
		close(conh.addrCh)
	}
	return nil
}

// StartDebugPrint starts printing debug info.
func (conh *consistentHashing) StartDebugPrint(interval time.Duration) {
	conh.mu.Lock()
	defer conh.mu.Unlock()
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-conh.debugPrintStopCh:
				// Exit to avoid goroutine leak.
				ticker.Stop()
				return

			case <-ticker.C:
				var buf bytes.Buffer
				var ncbuf bytes.Buffer // stores not-connected endpoints.
				for _, m := range conh.hr.Members() {
					ai := m.Val.(*addrInfo)
					if !ai.connected {
						if 0 == ncbuf.Len() {
							ncbuf.WriteString(" [NOT CONNECTED]:")
						}
						ncbuf.WriteString(fmt.Sprintf(" %s", ai.addr.Addr))
						continue
					}
					if 0 == buf.Len() {
						buf.WriteString("[CONNECTED TO]:")
					}
					buf.WriteString(fmt.Sprintf(" %s", ai.addr.Addr))
				}
				if 0 == buf.Len() && 0 == ncbuf.Len() {
					glog.Infoln("[NONE CONNECTED]")
				} else {
					glog.Infoln(buf.String(), ncbuf.String())
				}
			}
		}
	}()
}

func (conh *consistentHashing) StopDebugPrint() {
	conh.debugPrintStopCh <- struct{}{}
}

// Size returns the number of service addresses.
func (conh *consistentHashing) Size() int {
	conh.mu.Lock()
	defer conh.mu.Unlock()

	return len(conh.hr.Members())
}
