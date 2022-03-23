package health

import (
	"flag"
	"time"

	"github.com/golang/glog"
	prom "github.com/prometheus/client_golang/prometheus"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	hpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	"github.com/binchencoder/letsgo/hashring"
	"google.golang.org/grpc/balancer"
)

var (
	healthCheckInterval = flag.Duration("health-check-interval", time.Minute, "The interval of health checking")
	HealthCheckTimeout  = flag.Duration("health-check-timeout", 2000*time.Millisecond, "The timeout for health checking")

	healthCheckCounts = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "health_check_counts",
			Help:      "grpc health check counts.",
		},
		[]string{"code", "grpc_service"},
	)
	// to export
	HealthCheckCounts = healthCheckCounts

	healthCheckSuccessGauge = prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "health_check_success_gauge",
			Help:      "grpc health check success gauge.",
		},
		[]string{"grpc_service"},
	)

	healthCheckSuccessRate = prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "health_check_success_rate",
			Help:      "grpc health check success rate.",
		},
		[]string{"grpc_service"},
	)

	healthCheckLatency = prom.NewHistogramVec(
		prom.HistogramOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "health_check_latency",
			Help:      "grpc health check latency.",
			Buckets:   prom.DefBuckets,
		},
		[]string{"grpc_service"},
	)
	//to export
	HealthCheckLatency = healthCheckLatency
)

func init() {
	prom.MustRegister(healthCheckCounts)
	prom.MustRegister(healthCheckSuccessGauge)
	prom.MustRegister(healthCheckSuccessRate)
	prom.MustRegister(healthCheckLatency)
}

type sizer interface {
	Size() int
}

// StartHealthCheck starts a goroutine to send health check requests
// to gRPC service.
func StartHealthCheck(conn *grpc.ClientConn, balancer balancer.Balancer, svc string) chan<- struct{} {
	if _, ok := balancer.(sizer); !ok {
		glog.Errorf("The load balancer didn't implement Size(), service %s", svc)
		return nil
	}

	stopCh := make(chan struct{}, 1)

	go func() {
		cli := hpb.NewHealthClient(conn)
		ticker := time.NewTicker(*healthCheckInterval)
		req := hpb.HealthCheckRequest{
			Service: svc,
		}
		parent := withHealthCheck(context.Background())
		parent = hashring.WithHashKey(parent, "") // To avoid warning in services like authz.
		for {
			select {
			case <-ticker.C:
				size := balancer.(sizer).Size()
				if size == 0 {
					glog.Errorf("The load balancer has empty instances of service %s", svc)
					healthCheckSuccessGauge.WithLabelValues(svc).Set(0)
					healthCheckSuccessRate.WithLabelValues(svc).Set(0)
					continue
				}
				glog.V(5).Infof("Sending health check for %d instances of service %s", size, svc)

				success := 0
				for i := 0; i < size; i++ {
					start := time.Now()
					ctx, _ := context.WithTimeout(parent, *HealthCheckTimeout)
					if _, err := cli.Check(ctx, &req); err != nil {
						healthCheckCounts.WithLabelValues(status.Code(err).String(), svc).Inc()
						if !IsSafeError(err) {
							glog.Errorf("Failed to send a health check for an instance of service %s: %v", svc, err)
							continue
						}
					}
					healthCheckLatency.WithLabelValues(svc).Observe(time.Since(start).Seconds())
					healthCheckCounts.WithLabelValues("OK", svc).Inc()
					success++
					//time.Sleep(time.Second) // fuyc: uncomment to verify health check working in correct order.
				}
				healthCheckSuccessGauge.WithLabelValues(svc).Set(float64(success))
				healthCheckSuccessRate.WithLabelValues(svc).Set(float64(success) / float64(size))
			case <-stopCh:
				glog.Infof("got stopCh")
				ticker.Stop()
				return
			}
		}
	}()

	return stopCh
}

func IsSafeError(err error) bool {
	code := status.Code(err)
	return codes.Unimplemented == code || codes.InvalidArgument == code || codes.NotFound == code
}
