package balancer

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/naming"

	"jingoal.com/skylb-api/internal/health"
)

// mainService implements interface grpc.Balancer.
// According to the weight, a main service is selected.
// When a new service is launched, if the service has more weight,
// it will become the new main service.
// When the main service goes offline, a main service will be re-elected.
//
// Note that when the weight is less than 0,
// it will never come to serve as the main service.
//
// This balancer is mostly a copy of the grpc RoundRobin balancer.
type mainService struct {
	r          naming.Resolver
	w          naming.Watcher
	addrs      []*addrInfo // all the addresses the client should potentially connect
	mainAddr   *addrInfo
	mu         sync.Mutex
	nextHealth int // index of the next health check address to return for Get()
	// The channel to notify gRPC internals the list of addresses.
	// The client should connect to.
	addrCh chan []grpc.Address
	// The channel to block when there is no connected address available.
	waitCh chan struct{}
	// The Balancer is closed.
	done bool

	debugPrintStopCh chan struct{}
}

// MainService returns a Balancer that selects a main service address.
// It uses r to watch the name resolution
// updates and updates the addresses available correspondingly.
func MainService(r naming.Resolver) grpc.Balancer {
	return &mainService{
		r:                r,
		debugPrintStopCh: make(chan struct{}, 1),
	}
}

func (ms *mainService) watchAddrUpdates() error {
	updates, err := ms.w.Next()
	if err != nil {
		grpclog.Warningf("grpc: the naming watcher stops working due to %v.", err)
		return err
	}
	ms.mu.Lock()
	defer ms.mu.Unlock()
	for _, update := range updates {
		addr := grpc.Address{
			Addr:     update.Addr,
			Metadata: update.Metadata,
		}
		switch update.Op {
		case naming.Add:
			var exist bool
			for _, v := range ms.addrs {
				if addr.Addr == v.addr.Addr {
					v.addr.Metadata = addr.Metadata
					if addr.Metadata != nil {
						v.weight, _ = addr.Metadata.(int32)
					} else {
						v.weight = 0
					}
					exist = true
					grpclog.Infoln("grpc: The name resolver wanted to add an existing address: ", addr)
					break
				}
			}
			if exist {
				continue
			}
			newAddrInfo := &addrInfo{addr: addr}
			if update.Metadata != nil {
				newAddrInfo.weight, _ = update.Metadata.(int32)
			}
			ms.addrs = append(ms.addrs, newAddrInfo)
		case naming.Delete:
			for i, v := range ms.addrs {
				if addr.Addr == v.addr.Addr {
					copy(ms.addrs[i:], ms.addrs[i+1:])
					ms.addrs = ms.addrs[:len(ms.addrs)-1]
					break
				}
			}
		default:
			grpclog.Errorln("Unknown update.Op ", update.Op)
		}
	}
	// Make a copy of rr.addrs and write it onto rr.addrCh
	// so that gRPC internals gets notified.
	open := make([]grpc.Address, len(ms.addrs))
	for i, v := range ms.addrs {
		open[i] = v.addr
		if ms.done {
			return grpc.ErrClientConnClosing
		}
	}
	select {
	case <-ms.addrCh:
	default:
	}
	ms.addrCh <- open
	return nil
}

func (ms *mainService) Start(target string, config grpc.BalancerConfig) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if ms.done {
		return grpc.ErrClientConnClosing
	}
	if ms.r == nil {
		// If there is no name resolver installed, it is not needed to
		// do name resolution. In this case, target is added into rr.addrs
		// as the only address available and rr.addrCh stays nil.
		ms.addrs = append(ms.addrs, &addrInfo{addr: grpc.Address{Addr: target}})
		return nil
	}
	w, err := ms.r.Resolve(target)
	if err != nil {
		return err
	}
	ms.w = w
	ms.addrCh = make(chan []grpc.Address, 1)
	go func() {
		for {
			if err := ms.watchAddrUpdates(); err != nil {
				return
			}
		}
	}()
	return nil
}

// Up sets the connected state of addr and sends notification if there are
// pending Get() calls.
func (ms *mainService) Up(addr grpc.Address) func(error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	var cnt int
	for _, a := range ms.addrs {
		if a.addr.Addr == addr.Addr {
			if !a.connected {
				cnt++
				a.connected = true
			}
			if ms.mainAddr == nil {
				break
			}
			// Lose eligibility as a main-service.
			if ms.mainAddr.addr.Addr == a.addr.Addr && a.weight < 0 {
				ms.mainAddr = nil
			}
			// Preemption based on weight.
			if a.weight >= 0 && a.weight > ms.mainAddr.weight {
				ms.mainAddr = a
			}
			break
		}
	}
	// addr is only one which is connected.
	// Notify the Get() callers who are blocking.
	if cnt == 1 && ms.waitCh != nil {
		close(ms.waitCh)
		ms.waitCh = nil
	}
	return func(err error) {
		ms.down(addr, err)
	}
}

// down unsets the connected state of addr.
func (ms *mainService) down(addr grpc.Address, err error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	for _, a := range ms.addrs {
		if addr == a.addr {
			a.connected = false
			break
		}
	}
}

// Get returns the next addr in the rotation.
func (ms *mainService) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	ms.mu.Lock()
	if ms.done {
		ms.mu.Unlock()
		err = grpc.ErrClientConnClosing
		return
	}

	if health.IsHealthCheck(ctx) {
		return ms.getNextHealthAddr(ctx, opts)
	}

	if ms.mainAddr != nil && ms.mainAddr.connected {
		addr = ms.mainAddr.addr
		ms.mu.Unlock()
		return
	}

	return ms.getNextAddr(ctx, opts)
}

func (ms *mainService) getNextAddr(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	if len(ms.addrs) > 0 {
		maxWeightAddr := ms.addrs[0]
		// Weight is preferred.
		for i := 1; i < len(ms.addrs); i++ {
			if maxWeightAddr.weight < ms.addrs[i].weight {
				maxWeightAddr = ms.addrs[i]
			}
		}
		if maxWeightAddr.weight >= 0 && maxWeightAddr.connected {
			ms.mainAddr = maxWeightAddr
			addr = ms.mainAddr.addr
			ms.mu.Unlock()
			return
		}
	}
	if !opts.BlockingWait {
		if len(ms.addrs) == 0 {
			ms.mu.Unlock()
			err = grpc.Errorf(codes.Unavailable, "there is no address available")
			return
		}
		if ms.mainAddr == nil || !ms.mainAddr.connected {
			ms.mu.Unlock()
			err = grpc.Errorf(codes.Unavailable, "there is no suitable weight address available")
			return
		}
		return
	}
	var ch chan struct{}
	// Wait on rr.waitCh for non-failfast RPCs.
	if ms.waitCh == nil {
		ch = make(chan struct{})
		ms.waitCh = ch
	} else {
		ch = ms.waitCh
	}
	ms.mu.Unlock()
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		case <-ch:
			ms.mu.Lock()
			if ms.done {
				ms.mu.Unlock()
				err = grpc.ErrClientConnClosing
				return
			}

			if len(ms.addrs) > 0 {
				maxWeightAddr := ms.addrs[0]
				for i := 1; i < len(ms.addrs); i++ {
					if maxWeightAddr.weight < ms.addrs[i].weight {
						maxWeightAddr = ms.addrs[i]
					}
				}
				if maxWeightAddr.weight >= 0 && maxWeightAddr.connected {
					ms.mainAddr = maxWeightAddr
					addr = ms.mainAddr.addr
					ms.mu.Unlock()
					return
				}
			}
			// The newly added addr got removed by Down() again.
			if ms.waitCh == nil {
				ch = make(chan struct{})
				ms.waitCh = ch
			} else {
				ch = ms.waitCh
			}
			ms.mu.Unlock()
		}
	}
}

// getNextHealthAddr is a approximate dup of getNextAddr, with
// all occurrences of "rr.next" replaced with "rr.nextHealth".
func (ms *mainService) getNextHealthAddr(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	if len(ms.addrs) > 0 {
		if ms.nextHealth >= len(ms.addrs) {
			ms.nextHealth = 0
		}
		next := ms.nextHealth
		for {
			a := ms.addrs[next]
			next = (next + 1) % len(ms.addrs)
			if a.connected {
				addr = a.addr
				ms.nextHealth = next
				ms.mu.Unlock()
				return
			}
			if next == ms.nextHealth {
				// Has iterated all the possible address but none is connected.
				break
			}
		}
	}
	if !opts.BlockingWait {
		if len(ms.addrs) == 0 {
			ms.mu.Unlock()
			err = grpc.Errorf(codes.Unavailable, "there is no address available")
			return
		}
		// Returns the next addr on rr.addrs for failfast RPCs.]
		addr = ms.addrs[ms.nextHealth].addr
		ms.nextHealth++
		ms.mu.Unlock()
		return
	}

	var ch chan struct{}
	// Wait on rr.waitCh for non-failfast RPCs.
	if ms.waitCh == nil {
		ch = make(chan struct{})
		ms.waitCh = ch
	} else {
		ch = ms.waitCh
	}
	ms.mu.Unlock()
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		case <-ch:
			ms.mu.Lock()
			if ms.done {
				ms.mu.Unlock()
				err = grpc.ErrClientConnClosing
				return
			}

			if len(ms.addrs) > 0 {
				if ms.nextHealth >= len(ms.addrs) {
					ms.nextHealth = 0
				}
				next := ms.nextHealth
				for {
					a := ms.addrs[next]
					next = (next + 1) % len(ms.addrs)
					if a.connected {
						addr = a.addr
						ms.nextHealth = next
						ms.mu.Unlock()
						return
					}
					if next == ms.nextHealth {
						// Has iterated all the possible address but none is connected.
						break
					}
				}
			}
			// The newly added addr got removed by Down() again.
			if ms.waitCh == nil {
				ch = make(chan struct{})
				ms.waitCh = ch
			} else {
				ch = ms.waitCh
			}
			ms.mu.Unlock()
		}
	}
}

func (ms *mainService) Notify() <-chan []grpc.Address {
	return ms.addrCh
}

func (ms *mainService) Close() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.done = true
	if ms.w != nil {
		ms.w.Close()
	}
	if ms.waitCh != nil {
		close(ms.waitCh)
		ms.waitCh = nil
	}
	if ms.addrCh != nil {
		close(ms.addrCh)
	}
	glog.Warningln("mainservice lb has closed.")
	return nil
}

// StartDebugPrint starts printing debug info.
func (ms *mainService) StartDebugPrint(interval time.Duration) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ms.debugPrintStopCh:
				// Exit to avoid goroutine leak.
				ticker.Stop()
				return

			case <-ticker.C:
				var buf bytes.Buffer
				var ncbuf bytes.Buffer // stores not-connected endpoints.
				for _, addr := range ms.addrs {
					if !addr.connected {
						if 0 == ncbuf.Len() {
							ncbuf.WriteString(" [NOT CONNECTED]:")
						}
						ncbuf.WriteString(fmt.Sprintf(" %s", addr.addr.Addr))
						continue
					}
					if 0 == buf.Len() {
						buf.WriteString("[CONNECTED TO]:")
					}
					buf.WriteString(fmt.Sprintf(" %s", addr.addr.Addr))
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

func (ms *mainService) StopDebugPrint() {
	ms.debugPrintStopCh <- struct{}{}
}

// Size returns the number of service addresses.
func (ms *mainService) Size() int {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return len(ms.addrs)
}
