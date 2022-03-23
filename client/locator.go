package client

import (
	"flag"

	vexpb "github.com/binchencoder/gateway-proto/data"
	"github.com/binchencoder/letsgo/flags"
	"github.com/binchencoder/letsgo/service/naming"
	"github.com/binchencoder/skylb-apiv2/client/option"

	"github.com/binchencoder/skylb-apiv2/internal/skylb"
	pb "github.com/binchencoder/skylb-apiv2/proto"
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

// ServiceLocator defines the interface through which the client app locates
// the position of services and obtains load balancing support from SkyLB.
// It can be used for both gRPC or non-gRPC services.
type ServiceLocator interface {
	// Resolve resolves a service for load balance with lbHandler providing
	// connection handling logic.
	//
	// It should to be called for every service used by the client.
	Resolve(lbHandler option.LoadBalanceHandler, opts ...option.ResolveOption)

	// EnableHistogram enables historgram in client api metrics.
	//
	// Even if there are multiple ServiceCli instances, EnableHistogram
	// only needs to be called once, on any of those instances.
	EnableHistogram()

	// EnableFailFast makes service client doesn't wait for service to become
	// available in Start().
	EnableFailFast()

	// Start starts the service locator.
	//
	// Start can only be called once for each ServiceCli instance in the whole
	// lifecycle of an application.
	Start()

	// Shutdown turns the service client down. After shutdown, all grpc.Balancer
	// objects returned from Resolve() call can not be used any more.
	Shutdown()
}

// NewServiceLocator returns a new service client. It should be
// called only once in the whole lifecycle of an application.
func NewServiceLocator(clientServiceId vexpb.ServiceId) ServiceLocator {
	return skylb.NewServiceLocator(clientServiceId, map[string]string(DebugSvcEndpoints))
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
