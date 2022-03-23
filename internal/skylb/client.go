package skylb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	vexpb "github.com/binchencoder/gateway-proto/data"
	jg "github.com/binchencoder/letsgo/grpc"
	jn "github.com/binchencoder/letsgo/service/naming"
	"github.com/binchencoder/skylb-apiv2/client/option"
	"github.com/binchencoder/skylb-apiv2/internal/flags"
	"github.com/binchencoder/skylb-apiv2/metrics"
	pb "github.com/binchencoder/skylb-apiv2/proto"
)

// serviceClient implements interface skylb-api/client/ServiceClient.
type serviceClient struct {
	clientServiceId vexpb.ServiceId
	keeper          option.SkyLbKeeper
	specs           []*pb.ServiceSpec
	conns           []*grpc.ClientConn
	// lbs                map[string]grpc.Balancer
	healthCheckClosers []chan<- struct{}
	opts               map[string]option.ResolveOptions
	unaryInterceptors  []grpc.UnaryClientInterceptor
	failFast           bool
	resolveCount       int

	debugSvcEndpoints map[string]string

	resolveFullEps bool
}

// Resolve resolves a service spec and returns a load balancer handle.
// It needs to be called for every service used by the client.
func (sc *serviceClient) Resolve(spec *pb.ServiceSpec, opts ...option.ResolveOption) {
	ropts := option.ResolveOptions{}
	for _, opt := range opts {
		opt(&ropts)
	}
	sc.opts[spec.String()] = ropts
	sc.resolve(spec)
}

func (sc *serviceClient) resolve(spec *pb.ServiceSpec) {
	ropts := sc.opts[spec.String()]

	var r naming.Resolver
	if ep, ok := sc.debugSvcEndpoints[spec.ServiceName]; ok {
		// Debug mode, do not need to register service.

		parts := strings.SplitN(ep, ":", 2)
		if len(parts) != 2 {
			panic(fmt.Sprintf("Service instance endpoint should in format host:port, got %s", ep))
		}
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}

		ie := pb.InstanceEndpoint{
			Op:   pb.Operation_Add,
			Host: parts[0],
			Port: int32(port),
		}
		r = NewDebugResolver(&ie, spec)
	} else {
		r, _ = NewResolver(sc.keeper, spec)
		sc.resolveCount++
	}

	sc.specs = append(sc.specs, spec)
}

// EnableHistogram enables historgram in client api metrics.
// (This function doesn't have to be member method, but done so anyway
// in order to expose to interface ServiceCli)
func (sc *serviceClient) EnableHistogram() {
	withPrometheusHistogram = true
}

// EnableResolveFullEps enables to resolve full endpoints.
func (sc *serviceClient) EnableResolveFullEps() {
	sc.resolveFullEps = true
}

// DisableResolveFullEps disables resolving full endpoints.
func (sc *serviceClient) DisableResolveFullEps() {
	sc.resolveFullEps = false
}

// AddUnaryInterceptor adds a unary client interceptor to the client.
func (sc *serviceClient) AddUnaryInterceptor(incept grpc.UnaryClientInterceptor) {
	sc.unaryInterceptors = append(sc.unaryInterceptors, incept)
}

func getClientServiceName(clientServiceId vexpb.ServiceId) string {
	name, err := jn.ServiceIdToName(clientServiceId)
	if nil != err {
		glog.Errorf("Invalid client service id %d\n", clientServiceId)
		return "unknownService"
	}

	return name
}

// Start starts the service resolver and returns the grpc connection for
// each service through the callback function.
//
// Start can only be called once in the whole lifecycle of an application.
func (sc *serviceClient) Start(callback func(spec *pb.ServiceSpec, conn *grpc.ClientConn), options ...grpc.DialOption) {
	csId := sc.clientServiceId

	glog.Infof("Starting service client with %d service specs to resolve.", sc.resolveCount)

	csName, err := jn.ServiceIdToName(csId)
	if nil != err {
		glog.V(1).Infof("Invalid caller service id %d\n", csId)
		csName = fmt.Sprintf("!%d", csId)
	}

	go sc.keeper.Start(csId, csName, sc.resolveFullEps)

	for _, spec := range sc.specs {
		specCopy := &pb.ServiceSpec{
			Namespace:   spec.Namespace,
			ServiceName: spec.ServiceName,
			PortName:    spec.PortName,
		}

		tracer := jg.InitOpenTracing(getClientServiceName(sc.clientServiceId))
		openTracingInterceptor := otgrpc.OpenTracingClientInterceptor(tracer)

		metricsInterceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			ctx2 := jg.WithServiceId(ctx, int(csId))
			return metrics.UnaryClientInterceptor(csName, specCopy, ctx2, method, req, reply, cc, invoker, opts...)
		}

		metadataInterceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return jg.ClientToMetadataInterceptor(csName, ctx, method, req, reply, cc, invoker, opts...)
		}

		var conn *grpc.ClientConn
		var err error
		for {
			func() {
				defer func() {
					if p := recover(); p != nil {
						err = fmt.Errorf("%v", p)
					}
				}()

				incepts := make([]grpc.UnaryClientInterceptor, 0, len(sc.unaryInterceptors)+3)
				// openTracingInterceptor needs to be after
				// ExpBackoffUnaryClientInterceptor so that retry request will
				// have separate tracing.
				incepts = append(incepts, metricsInterceptor, jg.ExpBackoffUnaryClientInterceptor, openTracingInterceptor)
				incepts = append(incepts, sc.unaryInterceptors...)
				// ClientToMetadataInterceptor needs to be the last.
				incepts = append(incepts, metadataInterceptor)

				options = append(options, grpc.WithInsecure(),
					grpc.WithUnaryInterceptor(jg.ChainUnaryClient(incepts...)),
					grpc.WithStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
						ctx2 := jg.WithServiceId(ctx, int(csId))
						return metrics.StreamClientInterceptor(csName, specCopy, ctx2, desc, cc, method, streamer, opts...)
					}),
				)
				if sc.failFast {
					options = append(options,
						grpc.WithBlock(),
						// TODO(fuyc): here use a relatively longer duration as dial timeout.
						grpc.WithTimeout(*flags.SkylbRetryInterval*6),
					)
				}
				conn, err = grpc.Dial(spec.ServiceName, options...)
			}()

			if err == nil {
				break
			}

			glog.Warningf("Failed to dial service %q, %v.", spec.ServiceName, err)
			time.Sleep(*flags.SkylbRetryInterval)
		}

		sc.conns = append(sc.conns, conn)
		callback(spec, conn)
		// if *cflags.EnableHealthCheck {
		// 	closer := health.StartHealthCheck(conn, balancer, spec.ServiceName)
		// 	if closer != nil {
		// 		sc.healthCheckClosers = append(sc.healthCheckClosers, closer)
		// 	}
		// }
	}

	if withPrometheusHistogram {
		metrics.EnableClientHandlingTimeHistogram()
	}

	if !sc.failFast && sc.resolveCount > 0 {
		sc.keeper.WaitUntilReady()
	}
}

// Shutdown turns the service client down. After shutdown all grpc.Balancer
// objects returned from Resolve() call can not be used any more.
func (sc *serviceClient) Shutdown() {
	glog.V(3).Infof("Shutdown")
	for _, closer := range sc.healthCheckClosers {
		close(closer)
	}
	for _, conn := range sc.conns {
		if nil != conn {
			conn.Close()
		}
	}
	sc.keeper.Shutdown()
}

// EnableFailFast instructs the API framework to not wait all dependent
// services to be available. Here it only delegates the call to the keeper.
func (sc *serviceClient) EnableFailFast() {
	glog.V(3).Infoln("FailFast is enabled")
	sc.failFast = true
}

// NewServiceClient returns a new service client with the given debug service
// endpoints map.
func NewServiceClient(clientServiceId vexpb.ServiceId, dseps map[string]string) *serviceClient {
	return &serviceClient{
		clientServiceId:    clientServiceId,
		specs:              []*pb.ServiceSpec{},
		conns:              []*grpc.ClientConn{},
		healthCheckClosers: []chan<- struct{}{},
		opts:               map[string]option.ResolveOptions{},
		unaryInterceptors:  []grpc.UnaryClientInterceptor{},

		debugSvcEndpoints: dseps,
		resolveFullEps:    true,
	}
}
