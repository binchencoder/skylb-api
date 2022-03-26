package client

import (
	"flag"

	"google.golang.org/grpc"

	vexpb "github.com/binchencoder/gateway-proto/data"
	"github.com/binchencoder/letsgo/flags"
	"github.com/binchencoder/skylb-api/internal/skylb"
	"github.com/binchencoder/skylb-api/naming"
	pb "github.com/binchencoder/skylb-api/proto"
)

const (
	DefaultNameSpace = "default"
	DefaultPortName  = "grpc"
)

var (
	DebugSvcEndpoints = flags.StringMap{}
)

func init() {
	flag.Var(&DebugSvcEndpoints, "debug-svc-endpoint", "The debug service endpoint. If not empty, disable SkyLB resolving for that service.")
}

// ServiceCli defines the interface through which the client app obtains
// gRPC load balancing support from SkyLB.
type ServiceCli interface {
	// Resolve resolves a service spec.
	// It needs to be called for every service used by the client.
	Resolve(spec *pb.ServiceSpec, opts ...grpc.DialOption)

	// EnableHistogram enables historgram in client api metrics.
	//
	// Even if there are multiple ServiceCli instances, EnableHistogram
	// only needs to be called once, on any of those instances.
	EnableHistogram()

	// EnableFailFast makes service client doesn't wait for service to become
	// available in Start().
	EnableFailFast()

	// AddUnaryInterceptor adds a unary client interceptor to the client.
	AddUnaryInterceptor(incept grpc.UnaryClientInterceptor)

	// Start starts the service resolver and returns the grpc connection for
	// each service through the callback function.
	//
	// Start can only be called once for each ServiceCli instance in the whole
	// lifecycle of an application.
	Start(callback func(spec *pb.ServiceSpec, conn *grpc.ClientConn))

	// Shutdown turns the service client down. After shutdown, all grpc.Balancer
	// objects returned from Resolve() call can not be used any more.
	Shutdown()
}

// NewServiceSpec returns a new ServiceSpec struct with the given parameters.
func NewServiceSpec(namespace string, serviceId vexpb.ServiceId, portName string) *pb.ServiceSpec {
	if namespace == "" {
		namespace = DefaultNameSpace
	}
	if portName == "" {
		portName = DefaultPortName
	}
	serviceName, err := naming.ServiceIdToName(serviceId)
	if err != nil {
		panic("Unknown service ID.")
	}

	return &pb.ServiceSpec{
		Namespace:   namespace,
		ServiceName: serviceName,
		PortName:    portName,
	}
}

// NewDefaultServiceSpec returns a new ServiceSpec struct with the default
// namespace and default port name.
func NewDefaultServiceSpec(serviceId vexpb.ServiceId) *pb.ServiceSpec {
	return NewServiceSpec(DefaultNameSpace, serviceId, DefaultPortName)
}

// NewServiceCli returns a new service client. NewServiceCli() should be called
// once in the whole lifecycle of an application.
func NewServiceCli(clientServiceId vexpb.ServiceId) ServiceCli {
	return skylb.NewServiceClient(clientServiceId, map[string]string(DebugSvcEndpoints))
}
