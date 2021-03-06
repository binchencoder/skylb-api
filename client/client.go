package client

import (
	"google.golang.org/grpc"

	"github.com/binchencoder/skylb-api/client/option"
	"github.com/binchencoder/skylb-api/internal/skylb"
	pb "github.com/binchencoder/skylb-api/proto"
	vexpb "github.com/binchencoder/gateway-proto/data"
)

// TODO(zhwang): remove this file once we migrate all references to the
//               new API.

// ServiceCli defines the interface through which the client app obtains
// gRPC load balancing support from SkyLB.
//
// Deprecated: use ServiceLocator instead.
type ServiceCli interface {
	// Resolve resolves a service spec.
	// It needs to be called for every service used by the client.
	Resolve(spec *pb.ServiceSpec, opts ...option.ResolveOption)

	// EnableHistogram enables historgram in client api metrics.
	//
	// Even if there are multiple ServiceCli instances, EnableHistogram
	// only needs to be called once, on any of those instances.
	EnableHistogram()

	// EnableResolveFullEps enables to resolve full endpoints.
	// Deprecated: This method will be removed in future.
	EnableResolveFullEps()

	// EnableFailFast makes service client doesn't wait for service to become
	// available in Start().
	EnableFailFast()

	// DisableResolveFullEps disables resolving full endpoints.
	// Deprecated: This method will be removed in future.
	DisableResolveFullEps()

	// AddUnaryInterceptor adds a unary client interceptor to the client.
	AddUnaryInterceptor(incept grpc.UnaryClientInterceptor)

	// Start starts the service resolver and returns the grpc connection for
	// each service through the callback function.
	//
	// Start can only be called once for each ServiceCli instance in the whole
	// lifecycle of an application.
	Start(callback func(spec *pb.ServiceSpec, conn *grpc.ClientConn), options ...grpc.DialOption)

	// Shutdown turns the service client down. After shutdown, all grpc.Balancer
	// objects returned from Resolve() call can not be used any more.
	Shutdown()
}

// NewServiceCli returns a new service client. NewServiceCli() should be called
// once in the whole lifecycle of an application.
func NewServiceCli(clientServiceId vexpb.ServiceId) ServiceCli {
	return skylb.NewServiceClient(clientServiceId, map[string]string(DebugSvcEndpoints))
}
