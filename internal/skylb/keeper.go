package skylb

import (
	"context"
	"flag"
	"time"

	"github.com/binchencoder/letsgo/sync"
	pb "github.com/binchencoder/skylb-apiv2/proto"
	prom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/resolver"
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
	updateCh chan []*resolver.Address
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

// NewSkyLbKeeper returns a new skylb keeper.
func NewSkyLbKeeper() *skyLbKeeper {
	return &skyLbKeeper{
		services: make(map[string]*serviceEntry),
		readyCh:  make(chan struct{}, 1),
	}
}
