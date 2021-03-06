package skylb

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/golang/glog"
	prom "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/naming"
	"google.golang.org/grpc/status"

	vexpb "github.com/binchencoder/gateway-proto/data"
	"github.com/binchencoder/letsgo/sync"
	"github.com/binchencoder/skylb-api/internal/flags"
	"github.com/binchencoder/skylb-api/internal/rpccli"
	"github.com/binchencoder/skylb-api/metrics"
	pb "github.com/binchencoder/skylb-api/proto"
)

var (
	skylbAliveTimeout = flag.Duration("skylb-alive-timeout", 3*time.Minute, "The timeout duration to keep alive of the SkyLB endpoints updates. Recommended value is 3 times of the auto rectification interval")

	svcKeeperGauge = prom.NewGauge(
		prom.GaugeOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "service_keeper_gauge",
			Help:      "Number of service keepers.",
		},
	)
	svcKeeperRecvStreamGauge = prom.NewGauge(
		prom.GaugeOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "service_keeper_recv_stream_gauge",
			Help:      "Number of service keepers receiving stream.",
		},
	)
	svcWatcherUpdatesGauge = prom.NewGauge(
		prom.GaugeOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "service_watcher_updates_gauge",
			Help:      "Number of unconsumed service watcher updates.",
		},
	)
)

type serviceEntry struct {
	spec     *pb.ServiceSpec
	updateCh chan []*naming.Update
}

// skyLbKeeper keeps connectivity to SkyLb instance.
type skyLbKeeper struct {
	sync.RWLock

	services map[string]*serviceEntry
	readyCh  chan struct{}
	ready    bool
	cancel   context.CancelFunc
	stopped  bool
}

func init() {
	prom.MustRegister(svcKeeperGauge)
	prom.MustRegister(svcKeeperRecvStreamGauge)
	prom.MustRegister(svcWatcherUpdatesGauge)
}

func (sk *skyLbKeeper) RegisterService(spec *pb.ServiceSpec) <-chan []*naming.Update {
	ch := make(chan []*naming.Update, 100)

	key := calcServiceKey(spec)
	se := serviceEntry{
		spec:     spec,
		updateCh: ch,
	}
	sk.services[key] = &se

	glog.Infof("Registered to resolve service spec %s.%s on port name %q.", spec.Namespace, spec.ServiceName, spec.PortName)

	return ch
}

func (sk *skyLbKeeper) Start(csId vexpb.ServiceId, csName string, resolveFullEps bool) {
	svcKeeperGauge.Inc()

	glog.V(4).Infof("Starting SkyLB keeper for caller service ID %#v", csId)
	if len(sk.services) == 0 {
		svcKeeperGauge.Dec()
		return
	}

	req := pb.ResolveRequest{
		Services:             []*pb.ServiceSpec{},
		CallerServiceId:      csId,
		CallerServiceName:    csName,
		ResolveFullEndpoints: resolveFullEps,
	}
	for _, s := range sk.services {
		req.Services = append(req.Services, s.spec)
	}

	ctx, cancel := context.WithCancel(context.Background())
	sk.cancel = cancel

	for {
		if err := sk.start(ctx, resolveFullEps, &req); err != nil {
			if err == io.EOF {
				glog.Info("SkyLB server closed the resolve stream.")
			}
			if st, ok := status.FromError(err); ok && st.Code() == codes.Canceled {
				if sk.stopped {
					break
				}
			}
		}
		time.Sleep(*flags.SkylbRetryInterval)
	}

	for _, s := range sk.services {
		close(s.updateCh)
	}
	svcKeeperGauge.Dec()
}

func (sk *skyLbKeeper) start(ctx context.Context, resolveFullEps bool, req *pb.ResolveRequest) error {
	ctxt, cancel := context.WithCancel(ctx)
	skyCli, err := rpccli.NewGrpcClient(ctxt)
	if err != nil {
		glog.Errorf("Failed to create gRPC client to SkyLB, %+v, retry.", err)
		return err
	}

	stopCh := make(chan struct{}, 1)

	timer := time.NewTimer(*skylbAliveTimeout)
	go func(cancel context.CancelFunc, stopCh <-chan struct{}) {
		select {
		case <-timer.C:
			glog.V(4).Infof("Service keeper timeout to receive updates. Cancel and restart.")
			metrics.SkylbAliveTimeoutCounts.Inc()
		case <-stopCh:
			glog.V(4).Infof("Service keeper is shutdown.")
			if !timer.Stop() {
				<-timer.C
			}
		}
		cancel()
	}(cancel, stopCh)

	glog.V(5).Infof("Resolving request: %+v", req)
	rctx, _ := context.WithCancel(ctxt)
	stream, err := skyCli.Resolve(rctx, req)
	if err != nil {
		glog.Errorf("Failed to call RPC Resolve, %+v, retry.", err)
		close(stopCh)
		return err
	}
	glog.Infoln("Established resolve stream to SkyLB.")

	localEpsMap := make(map[string]map[string]struct{})

	readyMap := map[string]struct{}{}
	for {
		resp, err := stream.Recv()
		if err != nil {
			close(stopCh)
			cancel()
			return err
		}

		svcKeeperRecvStreamGauge.Inc()

		if !timer.Stop() {
			<-timer.C
		}
		timer.Reset(*skylbAliveTimeout)

		var updates []*naming.Update
		if svcEps := resp.GetSvcEndpoints(); svcEps != nil {
			lenEps := len(svcEps.InstEndpoints)
			svcName := svcEps.Spec.ServiceName
			glog.V(2).Infof("Received %d endpoint(s) for service %s", lenEps, svcName)
			metrics.RecordEndpointCount(svcName, lenEps)

			if resolveFullEps {
				localEps, ok := localEpsMap[svcEps.Spec.String()]
				if !ok {
					localEps = make(map[string]struct{})
					localEpsMap[svcEps.Spec.String()] = localEps
				}

				// The response holds full endpoints, we need to calculate
				// the deltas.
				eps := make(map[string]struct{})
				for _, ep := range svcEps.InstEndpoints {
					addr := fmt.Sprintf("%s:%d", ep.Host, ep.Port)
					eps[addr] = struct{}{}
					if _, ok := localEps[addr]; !ok {
						up := naming.Update{
							Op:   naming.Add,
							Addr: addr,
						}
						if ep.Weight != 0 {
							up.Metadata = ep.Weight
						}
						updates = append(updates, &up)
						localEps[addr] = struct{}{}
					}
				}
				for addr := range localEps {
					if _, ok := eps[addr]; !ok {
						up := naming.Update{
							Op:   naming.Delete,
							Addr: addr,
						}
						updates = append(updates, &up)
						delete(localEps, addr)
					}
				}
			} else {
				// TODO(fuyc): Remove the code once all clients changed
				//             to use the new protocol.
				for _, ep := range svcEps.InstEndpoints {
					up := naming.Update{
						Op:   toNamingOp(ep.Op),
						Addr: fmt.Sprintf("%s:%d", ep.Host, ep.Port),
					}
					updates = append(updates, &up)
				}
			}

			if len(updates) == 0 {
				svcKeeperRecvStreamGauge.Dec()
				continue
			}

			key := calcServiceKey(svcEps.Spec)
			if !sk.ready {
				if _, ok := readyMap[key]; !ok {
					readyMap[key] = struct{}{}
					if len(readyMap) == len(sk.services) {
						close(sk.readyCh)
						sk.ready = true
					}
				}
			}

			if glog.V(3) {
				var buf bytes.Buffer
				for i, up := range updates {
					if i > 0 {
						(&buf).WriteString(", ")
					}
					(&buf).WriteString(fmt.Sprintf("[%s]%s", opToString(up.Op), up.Addr))
				}
				glog.Infof("Received endpoints update for %s with value %+v.", key, buf.String())
			}
			if svc, ok := sk.services[key]; ok {
				svc.updateCh <- updates
			} else {
				glog.Warningf("Nil serviceEntry for key %s", key)
			}
		}
		svcKeeperRecvStreamGauge.Dec()
	}
}

func (sk *skyLbKeeper) Shutdown() {
	glog.V(3).Info("Shutting down keeper.")
	sk.stopped = true
	if sk.cancel != nil {
		sk.cancel()
	}
	metrics.ClearEndpointCount()
}

func (sk *skyLbKeeper) WaitUntilReady() {
	<-sk.readyCh
}

// NewSkyLbKeeper returns a new skylb keeper.
func NewSkyLbKeeper() *skyLbKeeper {
	return &skyLbKeeper{
		services: make(map[string]*serviceEntry),
		readyCh:  make(chan struct{}, 1),
	}
}
