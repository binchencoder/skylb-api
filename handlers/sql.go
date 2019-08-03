package handlers

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	lsql "jingoal.com/letsgo/sql"
	jb "jingoal.com/skylb-api/balancer"
	"jingoal.com/skylb-api/client/option"
	"jingoal.com/skylb-api/internal/flags"
	cflags "jingoal.com/skylb-api/internal/flags/client"
	health "jingoal.com/skylb-api/internal/health"
	pb "jingoal.com/skylb-api/proto"
	vexpb "jingoal.com/vexillary-client/proto/data"
)

const (
	sessionAvailable   = 0
	sessionUnavailable = 1
)

var (
	flagDBLivenessMonitorInterval = flag.Duration("db-liveness-monitor-interval", 15*time.Second, "The interval for database liveness monitor")
	flagDBLivenessMonitorTimeout  = flag.Duration("db-liveness-monitor-timeout", time.Second, "The timeout for database liveness monitor")
)

type sessionState int

// DBConnector defines a database connector.
type DBConnector struct {
	// OnConnect is called to create a database session.
	OnConnect func(host string, port int) (lsql.DBSession, error)

	// OnReady is called when a DBCli is ready for use.
	OnReady func(dbCli DBCli)
}

// DBCli defines the exposed interface for a load balanced SQL database client
// through SkyLB.
type DBCli interface {
	// GetNext returns the next database session to use.
	GetNext(ctx context.Context, opts grpc.BalancerGetOptions) (lsql.DBSession, error)
}

// dbClient implements interface DBCli and jb.DebugBalancer.
type dbClient struct {
	spec              *pb.ServiceSpec
	lb                grpc.Balancer
	dbConnector       DBConnector
	sessions          map[grpc.Address]lsql.DBSession
	states            map[grpc.Address]sessionState
	monitorCancelFunc func()
	down              map[grpc.Address]func(error)
	mu                sync.RWMutex
}

func (dbc *dbClient) GetNext(ctx context.Context, opts grpc.BalancerGetOptions) (lsql.DBSession, error) {
	dbc.mu.RLock()
	defer dbc.mu.RUnlock()

	n := len(dbc.sessions)
	for i := 0; i < n; i++ {
		addr, _, err := dbc.lb.Get(ctx, opts)
		if err != nil {
			return nil, err
		}

		if s, ok := dbc.states[addr]; ok && s == sessionUnavailable {
			continue
		}
		if s, ok := dbc.sessions[addr]; ok {
			glog.V(5).Infof("Chosen next addr: %s", addr.Addr)
			return s, nil
		}
	}
	return nil, errors.New("SkyLB SQL API: no database session available")
}

func (dbc *dbClient) close() error {
	if lb, ok := dbc.lb.(jb.DebugBalancer); ok {
		glog.Infof("StopDebugPrint for %v", reflect.TypeOf(lb))
		lb.StopDebugPrint()
	}

	dbc.mu.RLock()
	defer dbc.mu.RUnlock()

	if dbc.monitorCancelFunc != nil {
		dbc.monitorCancelFunc()
	}

	for _, s := range dbc.sessions {
		s.Close()
	}
	return nil
}

func (dbc *dbClient) startLivenessMonitor(ctx context.Context) {
	ticker := time.NewTicker(*flagDBLivenessMonitorInterval)
	done := ctx.Done()

loop:
	for {
		select {
		case <-done:
			ticker.Stop()
			break loop
		case <-ticker.C:
			dbc.mu.RLock()
			for addr, sess := range dbc.sessions {
				go dbc.checkDBHealth(addr, sess)
			}
			dbc.mu.RUnlock()
		}
	}
}

func (dbc *dbClient) checkDBHealth(addr grpc.Address, sess lsql.DBSession) {
	ctx, _ := context.WithTimeout(context.Background(), *flagDBLivenessMonitorTimeout)
	start := time.Now()
	if err := sess.PingContext(ctx); err != nil {
		dbc.mu.Lock()
		if dbc.states[addr] == sessionAvailable {
			dbc.states[addr] = sessionUnavailable
		}
		dbc.mu.Unlock()
		glog.Errorf("SkyLB SQL client failed to ping instance %s, %v", addr, err)
		health.HealthCheckCounts.WithLabelValues(err.Error(), dbc.spec.ServiceName).Inc()
	} else {
		dbc.mu.Lock()
		if dbc.states[addr] == sessionUnavailable {
			dbc.states[addr] = sessionAvailable
			glog.Infof("SkyLB SQL client resumed to ping instance %s", addr)
		}
		dbc.mu.Unlock()
		health.HealthCheckCounts.WithLabelValues("OK", dbc.spec.ServiceName).Inc()
		health.HealthCheckLatency.WithLabelValues(dbc.spec.ServiceName).Observe(time.Since(start).Seconds())
	}
}

func (dbc *dbClient) StartDebugPrint(interval time.Duration) {
	if lb, ok := dbc.lb.(jb.DebugBalancer); ok {
		glog.Infof("StartDebugPrint for %v", reflect.TypeOf(lb))
		lb.StartDebugPrint(*cflags.DebugSvcInterval)
	}
}

func (dbc *dbClient) StopDebugPrint() {
	if lb, ok := dbc.lb.(jb.DebugBalancer); ok {
		glog.Infof("StopDebugPrint from err for %v", reflect.TypeOf(lb))
		lb.StopDebugPrint()
	}
}

func (dbc *dbClient) Dial(ctx context.Context, svcName string, block bool) error {
	waitC := make(chan error, 1)
	go func() {
		defer close(waitC)

		config := grpc.BalancerConfig{}
		if err := dbc.lb.Start(svcName, config); err != nil {
			waitC <- err
			return
		}
		ch := dbc.lb.Notify()
		if ch != nil {
			if block {
				doneChan := make(chan struct{})
				go dbc.lbWatcher(doneChan)
				<-doneChan
			} else {
				go dbc.lbWatcher(nil)
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-waitC:
		if err != nil {
			return err
		}
	}

	return nil
}

// lbWatcher watches the Notify channel of the balancer in cc and manages
// connections accordingly. If doneChan is not nil, it is closed after the
// first successfull connection is made.
func (dbc *dbClient) lbWatcher(doneChan chan struct{}) {
	for addrs := range dbc.lb.Notify() {
		var (
			add []grpc.Address   // Addresses need to setup connections.
			del []lsql.DBSession // DB sessions need to tear down.
		)
		dbc.mu.Lock()
		for _, a := range addrs {
			if _, ok := dbc.sessions[a]; !ok {
				add = append(add, a)
			}
		}
		for k, c := range dbc.sessions {
			var keep bool
			for _, a := range addrs {
				if k == a {
					keep = true
					break
				}
			}
			if !keep {
				dbc.down[k](nil)
				del = append(del, c)
				delete(dbc.sessions, k)
			}
		}
		dbc.mu.Unlock()

		for _, a := range add {
			var err error
			if doneChan != nil {
				err = dbc.openDatabase(a, true)
				if err == nil {
					close(doneChan)
					doneChan = nil
				}
			} else {
				err = dbc.openDatabase(a, false)
			}
			if err != nil {
				glog.Error(err)
			}
		}
		for _, db := range del {
			db.Close()
		}
	}
}

func (dbc *dbClient) openDatabase(a grpc.Address, block bool) error {
	parts := strings.Split(a.Addr, ":")
	if len(parts) != 2 {
		glog.Errorf("Invalid database address: %s", a.Addr)
		return fmt.Errorf("Invalid database address: %s", a.Addr)
	}

	port, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return fmt.Errorf("Invalid database address: %s", a.Addr)
	}

	session, err := dbc.dbConnector.OnConnect(parts[0], int(port))
	if err != nil {
		return fmt.Errorf("Failed to connect to database address: %s", a.Addr)
	}

	if err := session.Ping(); err != nil {
		session.Close()
		return fmt.Errorf("Failed to ping database address: %s", a.Addr)
	}

	dbc.mu.Lock()
	defer dbc.mu.Unlock()

	dbc.sessions[a] = session
	dbc.down[a] = dbc.lb.Up(a)
	return nil
}

// SQLLoadBalanceHandler implements interface LoadBalanceHandler defined in
// skylb-api/client/option/option.go.
type SQLLoadBalanceHandler struct {
	spec        *pb.ServiceSpec
	dbConnector DBConnector
	dbClient    *dbClient
}

func (slbh *SQLLoadBalanceHandler) ServiceSpec() *pb.ServiceSpec {
	return slbh.spec
}

func (slbh *SQLLoadBalanceHandler) BeforeResolve(spec *pb.ServiceSpec, r naming.Resolver, ropts *option.ResolveOptions) {
	var balancer grpc.Balancer
	bc := ropts.BalancerCreator()
	if bc == nil {
		balancer = grpc.RoundRobin(r)
	} else {
		balancer = bc(r)
	}

	dbcli := dbClient{
		spec:        spec,
		lb:          balancer,
		dbConnector: slbh.dbConnector,
		sessions:    make(map[grpc.Address]lsql.DBSession),
		states:      make(map[grpc.Address]sessionState),
		down:        make(map[grpc.Address]func(error)),
	}
	slbh.dbClient = &dbcli

	dbcli.StartDebugPrint(*cflags.DebugSvcInterval)
}

func (slbh *SQLLoadBalanceHandler) AfterResolve(spec *pb.ServiceSpec, csId vexpb.ServiceId, csName string, keeper option.SkyLbKeeper, tracer opentracing.Tracer, failFast bool) {
	dbcli := slbh.dbClient
	dbconn := slbh.dbConnector

	var err error
	for {
		func() {
			defer func() {
				if p := recover(); p != nil {
					err = fmt.Errorf("%v", p)
				}
			}()

			err = dbcli.Dial(context.Background(), spec.ServiceName, false)
		}()

		if err == nil {
			break
		}

		dbcli.StopDebugPrint()

		glog.Warningf("Failed to dial service %q, %v.", spec.ServiceName, err)
		if failFast {
			break
		}
		time.Sleep(*flags.SkylbRetryInterval)
	}

	if *cflags.EnableHealthCheck {
		ctx, cancel := context.WithCancel(context.Background())
		dbcli.monitorCancelFunc = cancel
		go dbcli.startLivenessMonitor(ctx)
		glog.Infof("Succeeded to start monitor of %s after resolver.", spec.ServiceName)
	}

	dbconn.OnReady(dbcli)
}

func (slbh *SQLLoadBalanceHandler) OnShutdown() {
	slbh.dbClient.close()
}

// NewSQLLoadBalanceHandler returns a new LoadBalanceHandler for SQL
// data sources.
func NewSQLLoadBalanceHandler(spec *pb.ServiceSpec, connector DBConnector) *SQLLoadBalanceHandler {
	return &SQLLoadBalanceHandler{
		spec:        spec,
		dbConnector: connector,
	}
}
