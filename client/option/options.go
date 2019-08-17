package option

import (
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	vexpb "github.com/binchencoder/gateway-proto/data"
	pb "github.com/binchencoder/skylb-api/proto"
)

// SkyLbKeeper defines the interface for a SkyLB keeper.
type SkyLbKeeper interface {
	// RegisterService registers the service with the given spec to the keeper
	// and returns a channel for the service instance updates.
	RegisterService(spec *pb.ServiceSpec) <-chan []*naming.Update

	// Start starts the keeper with the given client service ID and name.
	Start(csId vexpb.ServiceId, csName string, resolveFullEps bool)

	// Shutdown shuts down the keeper.
	Shutdown()

	// WaitUntilReady blocks the caller until the keeper receives the initial
	// endpoints for all service specs.
	WaitUntilReady()
}

// LoadBalanceHandler defines the interface to handle the notification logic
// for different clients in SkyLB API load balancing.
type LoadBalanceHandler interface {
	// Returns the service spec for this handler.
	ServiceSpec() *pb.ServiceSpec

	// BeforeResolve is called before SkyLB API resolves the given spec.
	BeforeResolve(spec *pb.ServiceSpec, r naming.Resolver, ropts *ResolveOptions)

	// AfterResolve is called after SkyLB API resolved the given spec.
	AfterResolve(spec *pb.ServiceSpec, csId vexpb.ServiceId, csName string, keeper SkyLbKeeper, tracer opentracing.Tracer, failFast bool)

	// OnShutdown is called when the SkyLB API is shutting down.
	OnShutdown()
}

// BalancerCreator is a function which creates a grpc Balancer.
type BalancerCreator func(naming.Resolver) grpc.Balancer

// ResolveOptions configure a resolve call.
type ResolveOptions struct {
	balancerCreator BalancerCreator
}

// BalancerCreator returns the load balancer creator.
func (ropts *ResolveOptions) BalancerCreator() BalancerCreator {
	return ropts.balancerCreator
}

// ResolveOption configures how we set up the resolve call.
type ResolveOption func(*ResolveOptions)

// WithBalancerCreator returns a ResolveOption which sets a load
// balancer creator.
func WithBalancerCreator(bc BalancerCreator) ResolveOption {
	return func(o *ResolveOptions) {
		o.balancerCreator = bc
	}
}
