package handlers

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	jg "jingoal.com/letsgo/grpc"
	jb "jingoal.com/skylb-api/balancer"
	"jingoal.com/skylb-api/client/option"
	"jingoal.com/skylb-api/internal/flags"
	cflags "jingoal.com/skylb-api/internal/flags/client"
	"jingoal.com/skylb-api/internal/health"
	"jingoal.com/skylb-api/metrics"
	pb "jingoal.com/skylb-api/proto"
	vexpb "jingoal.com/vexillary-client/proto/data"
)

// GrpcLoadBalanceHandler implements interface LoadBalanceHandler defined in
// skylb-api/client/option/option.go.
type GrpcLoadBalanceHandler struct {
	spec               *pb.ServiceSpec
	callback           func(conn *grpc.ClientConn)
	conns              []*grpc.ClientConn
	healthCheckClosers []chan<- struct{}
	options            []grpc.DialOption
	lbs                map[string]grpc.Balancer
	unaryInterceptors  []grpc.UnaryClientInterceptor
}

func (glbh *GrpcLoadBalanceHandler) ServiceSpec() *pb.ServiceSpec {
	return glbh.spec
}

// AddUnaryInterceptor adds a unary client interceptor to the client.
func (glbh *GrpcLoadBalanceHandler) AddUnaryInterceptor(incept grpc.UnaryClientInterceptor) {
	glbh.unaryInterceptors = append(glbh.unaryInterceptors, incept)
}

func (glbh *GrpcLoadBalanceHandler) BeforeResolve(spec *pb.ServiceSpec, r naming.Resolver, ropts *option.ResolveOptions) {
	var balancer grpc.Balancer
	bc := ropts.BalancerCreator()
	if bc == nil {
		if *cflags.EnableHealthCheck {
			balancer = jb.RoundRobin(r)
		} else {
			balancer = grpc.RoundRobin(r)
		}
		glog.Infof("Skylb client created round-robin load balancer.")
	} else {
		balancer = bc(r)

		glog.Infof("Skylb client created %v load balancer.", reflect.TypeOf(balancer))
	}
	glbh.lbs[spec.String()] = balancer

	if lb, ok := balancer.(jb.DebugBalancer); ok {
		glog.Infof("StartDebugPrint for %v", reflect.TypeOf(lb))
		lb.StartDebugPrint(*cflags.DebugSvcInterval)
	}
}

func (glbh *GrpcLoadBalanceHandler) AfterResolve(spec *pb.ServiceSpec, csId vexpb.ServiceId, csName string, keeper option.SkyLbKeeper, tracer opentracing.Tracer, failFast bool) {
	openTracingInterceptor := otgrpc.OpenTracingClientInterceptor(tracer)

	metricsInterceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx2 := jg.WithServiceId(ctx, int(csId))
		return metrics.UnaryClientInterceptor(csName, spec, ctx2, method, req, reply, cc, invoker, opts...)
	}

	metadataInterceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return jg.ClientToMetadataInterceptor(csName, ctx, method, req, reply, cc, invoker, opts...)
	}

	incepts := make([]grpc.UnaryClientInterceptor, 0, len(glbh.unaryInterceptors)+3)
	// openTracingInterceptor needs to be after
	// ExpBackoffUnaryClientInterceptor so that retry request will
	// have separate tracing.
	incepts = append(incepts, metricsInterceptor, jg.ExpBackoffUnaryClientInterceptor, openTracingInterceptor)
	incepts = append(incepts, glbh.unaryInterceptors...)
	// ClientToMetadataInterceptor needs to be the last.
	incepts = append(incepts, metadataInterceptor)

	lb := glbh.lbs[spec.String()]
	glbh.options = append(glbh.options, grpc.WithInsecure(),
		grpc.WithBalancer(lb),
		grpc.WithUnaryInterceptor(jg.ChainUnaryClient(incepts...)),
		grpc.WithStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			ctx2 := jg.WithServiceId(ctx, int(csId))
			return metrics.StreamClientInterceptor(csName, spec, ctx2, desc, cc, method, streamer, opts...)
		}),
	)
	if failFast {
		glbh.options = append(glbh.options,
			grpc.WithBlock(),
			grpc.WithTimeout(*flags.SkylbRetryInterval),
		)
	}

	var conn *grpc.ClientConn
	var err error
	balancer := glbh.lbs[spec.String()]
	for {
		func() {
			defer func() {
				if p := recover(); p != nil {
					err = fmt.Errorf("%v", p)
				}
			}()

			conn, err = grpc.Dial(spec.ServiceName, glbh.options...)
		}()

		if err == nil {
			break
		}

		if lb2, ok := balancer.(jb.DebugBalancer); ok {
			glog.Infof("StopDebugPrint from err for %v", reflect.TypeOf(lb2))
			lb2.StopDebugPrint()
		}

		glog.Warningf("Failed to dial service %q, %v.", spec.ServiceName, err)
		if failFast {
			break
		}
		time.Sleep(*flags.SkylbRetryInterval)
	}

	glbh.conns = append(glbh.conns, conn)
	glbh.callback(conn)
	if *cflags.EnableHealthCheck {
		closer := health.StartHealthCheck(conn, balancer, spec.ServiceName)
		if closer != nil {
			glbh.healthCheckClosers = append(glbh.healthCheckClosers, closer)
		}
	}
}

func (glbh *GrpcLoadBalanceHandler) OnShutdown() {
	for _, closer := range glbh.healthCheckClosers {
		close(closer)
	}
	for _, conn := range glbh.conns {
		if nil != conn {
			conn.Close()
		}
	}
	glbh.conns = nil

	for _, balancer := range glbh.lbs {
		if lb, ok := balancer.(jb.DebugBalancer); ok {
			glog.Infof("StopDebugPrint for %v", reflect.TypeOf(lb))
			lb.StopDebugPrint()
		}
	}
	glbh.lbs = nil
}

// NewGrpcLoadBalanceHandler returns a new LoadBalanceHandler for gRPC service.
func NewGrpcLoadBalanceHandler(spec *pb.ServiceSpec, callback func(conn *grpc.ClientConn), options ...grpc.DialOption) *GrpcLoadBalanceHandler {
	return &GrpcLoadBalanceHandler{
		spec:               spec,
		callback:           callback,
		options:            options,
		conns:              []*grpc.ClientConn{},
		healthCheckClosers: []chan<- struct{}{},
		lbs:                map[string]grpc.Balancer{},
		unaryInterceptors:  []grpc.UnaryClientInterceptor{},
	}
}
