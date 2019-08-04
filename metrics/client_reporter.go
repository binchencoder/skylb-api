// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package metrics

import (
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"

	pb "github.com/binchencoder/skylb-api/proto"
)

var (
	clientStartedCounter = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "started_total",
			Help:      "Total number of RPCs started on the client.",
		}, []string{"grpc_type", "self_service", "grpc_method", "grpc_service"})
	// grpc_service now means the callee service. self_service means current service itself.
	// Same for the following.

	clientHandledCounter = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "handled_total",
			Help:      "Total number of RPCs completed by the client, regardless of success or failure.",
		}, []string{"grpc_type", "self_service", "grpc_method", "grpc_code", "grpc_service"})

	clientStreamMsgReceived = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "msg_received_total",
			Help:      "Total number of RPC stream messages received by the client.",
		}, []string{"grpc_type", "self_service", "grpc_method", "grpc_service"})

	clientStreamMsgSent = prom.NewCounterVec(
		prom.CounterOpts{
			Namespace: "skylb",
			Subsystem: "client",
			Name:      "msg_sent_total",
			Help:      "Total number of gRPC stream messages sent by the client.",
		}, []string{"grpc_type", "self_service", "grpc_method", "grpc_service"})

	clientHandledHistogramEnabled = false
	clientHandledHistogramOpts    = prom.HistogramOpts{
		Namespace: "skylb",
		Subsystem: "client",
		Name:      "handling_seconds",
		Help:      "Histogram of response latency (seconds) of the gRPC until it is finished by the application.",
		Buckets:   prom.DefBuckets,
	}
	clientHandledHistogram *prom.HistogramVec
)

func init() {
	prom.MustRegister(clientStartedCounter)
	prom.MustRegister(clientHandledCounter)
	prom.MustRegister(clientStreamMsgReceived)
	prom.MustRegister(clientStreamMsgSent)
}

// EnableClientHandlingTimeHistogram turns on recording of handling time of RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query.
func EnableClientHandlingTimeHistogram(opts ...HistogramOption) {
	for _, o := range opts {
		o(&clientHandledHistogramOpts)
	}
	if !clientHandledHistogramEnabled {
		clientHandledHistogram = prom.NewHistogramVec(
			clientHandledHistogramOpts,
			[]string{"grpc_type", "self_service", "grpc_method", "grpc_service"},
		)
		prom.Register(clientHandledHistogram)
	}
	clientHandledHistogramEnabled = true
}

type clientReporter struct {
	rpcType     grpcType
	serviceName string // Now it means callee service name, not original "grpc service name"
	methodName  string
	startTime   time.Time
	spec        *pb.ServiceSpec
}

func newClientReporter(spec *pb.ServiceSpec, rpcType grpcType, fullMethod string, csName string) *clientReporter {
	r := &clientReporter{
		rpcType:     rpcType,
		spec:        spec,
		serviceName: csName,
	}
	if clientHandledHistogramEnabled {
		r.startTime = time.Now()
	}
	//r.serviceName, r.methodName = splitMethodName(fullMethod)
	_, r.methodName = splitMethodName(fullMethod)
	clientStartedCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.spec.ServiceName).Inc()
	return r
}

func (r *clientReporter) ReceivedMessage() {
	clientStreamMsgReceived.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.spec.ServiceName).Inc()
}

func (r *clientReporter) SentMessage() {
	clientStreamMsgSent.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.spec.ServiceName).Inc()
}

func (r *clientReporter) Handled(code codes.Code) {
	clientHandledCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, code.String(), r.spec.ServiceName).Inc()
	if clientHandledHistogramEnabled {
		clientHandledHistogram.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.spec.ServiceName).Observe(time.Since(r.startTime).Seconds())
	}
}
