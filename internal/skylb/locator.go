package skylb

import (
	"fmt"
	"strconv"
	"strings"

	vexpb "github.com/binchencoder/gateway-proto/data"
	jg "github.com/binchencoder/letsgo/grpc"
	"github.com/binchencoder/skylb-apiv2/client/option"
	"github.com/binchencoder/skylb-apiv2/metrics"
	"github.com/binchencoder/skylb-apiv2/naming"
	pb "github.com/binchencoder/skylb-apiv2/proto"
	"github.com/golang/glog"
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
	opts              map[string]option.ResolveOptions
	lbHandlers        map[string]option.LoadBalanceHandler
	failFast          bool
	resolveCount      int

	debugSvcEndpoints map[string]string
}

// Resolve resolves a service spec.
// It needs to be called for every service used by the client.
func (sl *serviceLocator) Resolve(lbHandler option.LoadBalanceHandler, opts ...option.ResolveOption) {
	spec := lbHandler.ServiceSpec()
	sl.lbHandlers[spec.String()] = lbHandler

	// ropts := option.ResolveOptions{}
	// for _, opt := range opts {
	// 	opt(&ropts)
	// }
	// sl.opts[spec.String()] = ropts

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
		sl.resolveCount++
	}

	sl.specs = append(sl.specs, spec)
	ropts := sl.opts[spec.String()]

	handler := sl.lbHandlers[spec.String()]
	handler.BeforeResolve(spec, r, &ropts)
}

func ParseEndpoint(ep string) pb.InstanceEndpoint {
	parts := strings.SplitN(ep, ":", 2)
	if len(parts) != 2 {
		panic(fmt.Sprintf("Service instance endpoint should in format host:port, got %s", ep))
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}

	return pb.InstanceEndpoint{
		Op:   pb.Operation_Add,
		Host: parts[0],
		Port: int32(port),
	}
}

// EnableHistogram enables historgram in client api metrics.
// (This function doesn't have to be member method, but done so anyway
// in order to expose to interface ServiceCli)
func (sl *serviceLocator) EnableHistogram() {
	withPrometheusHistogram = true
}

// Start starts the service resolver.
//
// Start can only be called once in the whole lifecycle of an application.
func (sl *serviceLocator) Start() {
	glog.Infof("Starting service client with %d service specs to resolve.", sl.resolveCount)

	csId := sl.clientServiceId

	csName, err := naming.ServiceIdToName(csId)
	if nil != err {
		glog.V(1).Infof("Invalid caller service id %d\n", csId)
		csName = fmt.Sprintf("!%d", csId)
	}

	go sl.keeper.Start(csId, csName, true /* resolveFullEps */)

	tracer := jg.InitOpenTracing(csName)
	for _, spec := range sl.specs {
		specCopy := &pb.ServiceSpec{
			Namespace:   spec.Namespace,
			ServiceName: spec.ServiceName,
			PortName:    spec.PortName,
		}

		handler := sl.lbHandlers[spec.String()]
		handler.AfterResolve(specCopy, csId, csName, sl.keeper, tracer, sl.failFast)
	}

	if withPrometheusHistogram {
		metrics.EnableClientHandlingTimeHistogram()
	}

	if !sl.failFast && sl.resolveCount > 0 {
		sl.keeper.WaitUntilReady()
	}
}

// Shutdown turns the service client down. After shutdown all grpc.Balancer
// objects returned from Resolve() call can not be used any more.
func (sl *serviceLocator) Shutdown() {
	for _, handler := range sl.lbHandlers {
		handler.OnShutdown()
	}
	sl.keeper.Shutdown()
}

// EnableFailFast makes service client doesn't wait for service to become
// available in Start().
func (sl *serviceLocator) EnableFailFast() {
	glog.V(3).Infoln("FailFast is enabled")
	sl.failFast = true
}

// NewServiceLocator returns a new service locator with the given debug service
// endpoints map.
func NewServiceLocator(clientServiceId vexpb.ServiceId, dseps map[string]string) *serviceLocator {
	return &serviceLocator{
		clientServiceId: clientServiceId,
		keeper:          NewSkyLbKeeper(),
		specs:           []*pb.ServiceSpec{},

		debugSvcEndpoints: dseps,
	}
}
