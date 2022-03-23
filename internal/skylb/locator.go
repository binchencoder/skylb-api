package skylb

import (
	vexpb "github.com/binchencoder/gateway-proto/data"
	"github.com/binchencoder/skylb-apiv2/client/option"
	pb "github.com/binchencoder/skylb-apiv2/proto"
	"google.golang.org/grpc/resolver"
)

var (
	withPrometheusHistogram = true
)

// serviceLocator implements interface skylb-apiv2/client/ServiceLocator.
type serviceLocator struct {
	clientServiceId   vexpb.ServiceId
	clientServiceName string
	keeper            option.SkyLbKeeper
	specs             []*pb.ServiceSpec
	lbHandlers        map[string]option.LoadBalanceHandler
	failFast          bool

	debugSvcEndpoints map[string]string
}

// Resolve resolves a service spec.
// It needs to be called for every service used by the client.
func (sl *serviceLocator) Resolve(lbHandler option.LoadBalanceHandler) {
	spec := lbHandler.ServiceSpec()
	sl.lbHandlers[spec.String()] = lbHandler

	sl.resolveService(spec)
}

func (sl *serviceLocator) resolveService(spec *pb.ServiceSpec) {
	// Note: non-gRPC handlers also use the same resolving mechanism from gRPC.
	var r resolver.Resolver
	if ep, ok := sl.debugSvcEndpoints[spec.ServiceName]; ok {
		// Debug mode, do not need to register service.

		ie := ParseEndpoint(ep)
		r = NewDebugResolver(&ie, spec)
	} else {
		// TODO(chenbin)
		r = nil
	}

	sl.specs = append(sl.specs, spec)

	handler := sl.lbHandlers[spec.String()]
	handler.BeforeResolve(spec, r)
}

// NewServiceLocator returns a new service locator with the given debug service
// endpoints map.
func NewServiceLocator(clientServiceId vexpb.ServiceId, dseps map[string]string) *serviceLocator {
	return &serviceLocator{
		clientServiceId: clientServiceId,
		keeper:          NewSkyLbKeeper(),
		specs:           []*pb.ServiceSpec{},
		lbHandlers:      map[string]option.LoadBalanceHandler{},

		debugSvcEndpoints: dseps,
	}
}
