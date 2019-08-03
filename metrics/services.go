package metrics

import (
	"sync"

	prom "github.com/prometheus/client_golang/prometheus"
)

var (
	SvcEndpointGauge = prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "service_endpoint_gauge",
			Help:      "Number of service endpoints.",
		},
		[]string{"grpc_service"},
	)
	// TODO(fuyc): SvcEndpointCount is used to check against skylb_web_health_gauge.
	// To consult prometheus community about why there is not public interface to read data from gauge.
	SvcEndpointCount = make(map[string]int)
	epcLock          sync.Mutex

	// Enable flag:
	// -skylb-alive-timeout=3s
	// to test this metric.
	SkylbAliveTimeoutCounts = prom.NewCounter(
		prom.CounterOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "skylb_alive_timeout_counts",
			Help:      "Skylb alive timeout counts.",
		},
	)
)

func init() {
	prom.MustRegister(SvcEndpointGauge)
	prom.MustRegister(SkylbAliveTimeoutCounts)
}

// RecordEndpointCount records each service's instance count known to this client.
func RecordEndpointCount(serviceName string, count int) {
	SvcEndpointGauge.WithLabelValues(serviceName).Set(float64(count))

	epcLock.Lock()
	defer epcLock.Unlock()

	if val, had := SvcEndpointCount[serviceName]; !had || val != count {
		SvcEndpointCount[serviceName] = count
	}
}

// ClearEndpointCount clears previous records of service endpoint count.
func ClearEndpointCount() {
	SvcEndpointGauge.Reset()

	epcLock.Lock()
	defer epcLock.Unlock()

	SvcEndpointCount = make(map[string]int)
}
