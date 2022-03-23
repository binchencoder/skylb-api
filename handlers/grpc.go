package handlers

import (
	"context"
	"fmt"
	"time"

	vexpb "github.com/binchencoder/gateway-proto/data"
	jg "github.com/binchencoder/letsgo/grpc"
	"github.com/binchencoder/skylb-apiv2/client/option"
	"github.com/binchencoder/skylb-apiv2/internal/flags"
	"github.com/binchencoder/skylb-apiv2/metrics"
	pb "github.com/binchencoder/skylb-apiv2/proto"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// GrpcLoadBalanceHandler implements interface LoadBalanceHandler defined in
// skylb-apiv2/client/option/options.go.

type GrpcLoadBalanceHandler struct {
	spec               *pb.ServiceSpec
	callback           func(conn *grpc.ClientConn)
	conns              []*grpc.ClientConn
	healthCheckClosers []chan<- struct{}
	options            []grpc.DialOption
	unaryInterceptors  []grpc.UnaryClientInterceptor
}

func (glbh *GrpcLoadBalanceHandler) ServiceSpec() *pb.ServiceSpec {
	return glbh.spec
}

// AddUnaryInterceptor adds a unary client interceptor to the client.
func (glbh *GrpcLoadBalanceHandler) AddUnaryInterceptor(incept grpc.UnaryClientInterceptor) {
	glbh.unaryInterceptors = append(glbh.unaryInterceptors, incept)
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

	// TODO(chenbin) use WithTransportCredentials
	glbh.options = append(glbh.options, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(jg.ChainUnaryClient(incepts...)),
		grpc.WithStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			ctx2 := jg.WithServiceId(ctx, int(csId))
			return metrics.StreamClientInterceptor(csName, spec, ctx2, desc, cc, method, streamer, opts...)
		}),
	)
	// TODO(chenbin) FailFast
	if failFast {
		glbh.options = append(glbh.options,
			grpc.WithBlock(),
			grpc.WithTimeout(*flags.SkylbRetryInterval),
		)
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

			conn, err = grpc.Dial(spec.ServiceName, glbh.options...)
		}()

		if err == nil {
			break
		}

		// if lb2, ok := balancer.(jb.DebugBalancer); ok {
		// 	glog.Infof("StopDebugPrint from err for %v", reflect.TypeOf(lb2))
		// 	lb2.StopDebugPrint()
		// }

		glog.Warningf("Failed to dial service %q, %v.", spec.ServiceName, err)
		if failFast {
			break
		}
		time.Sleep(*flags.SkylbRetryInterval)
	}

	glbh.conns = append(glbh.conns, conn)
	glbh.callback(conn)
	// TODO(chenbin)
	// if *cflags.EnableHealthCheck {
	// 	closer := health.StartHealthCheck(conn, balancer, spec.ServiceName)
	// 	if closer != nil {
	// 		glbh.healthCheckClosers = append(glbh.healthCheckClosers, closer)
	// 	}
	// }
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

	// for _, balancer := range glbh.lbs {
	// 	if lb, ok := balancer.(jb.DebugBalancer); ok {
	// 		glog.Infof("StopDebugPrint for %v", reflect.TypeOf(lb))
	// 		lb.StopDebugPrint()
	// 	}
	// }
	// glbh.lbs = nil
}

// NewGrpcLoadBalanceHandler returns a new LoadBalanceHandler for gRPC service.
func NewGrpcLoadBalanceHandler(spec *pb.ServiceSpec, callback func(conn *grpc.ClientConn), options ...grpc.DialOption) *GrpcLoadBalanceHandler {
	return &GrpcLoadBalanceHandler{
		spec:               spec,
		callback:           callback,
		options:            options,
		conns:              []*grpc.ClientConn{},
		healthCheckClosers: []chan<- struct{}{},
		unaryInterceptors:  []grpc.UnaryClientInterceptor{},
	}
}
