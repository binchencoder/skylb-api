package handlers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	lgrpc "binchencoder.com/letsgo/testing/mocks/grpc"
	ltracing "binchencoder.com/letsgo/testing/mocks/tracing"
	"binchencoder.com/skylb-api/client/option"
	pb "binchencoder.com/skylb-api/proto"
	jt "binchencoder.com/skylb-api/testing"
)

func TestGrpcLoadBalanceHandler(t *testing.T) {
	Convey("gRPC load balance handler", t, func() {
		var callbackCalled bool
		callback := func(conn *grpc.ClientConn) {
			callbackCalled = true
		}
		spec := pb.ServiceSpec{
			Namespace:   "ns",
			ServiceName: "svc",
			PortName:    "ptn",
		}
		Convey("when creating a gRPC load balance handler", func() {
			handler := NewGrpcLoadBalanceHandler(&spec, callback)
			So(handler.spec, ShouldResemble, &spec)
			Convey("when getting back service spec", func() {
				s := handler.ServiceSpec()
				So(s, ShouldResemble, &spec)
			})

			unary := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				return nil
			}
			Convey("when adding a unary interceptor", func() {
				handler.AddUnaryInterceptor(unary)
				So(len(handler.unaryInterceptors), ShouldEqual, 1)
				So(handler.unaryInterceptors[0], ShouldEqual, unary)
			})
			Convey("before resolving service", func() {
				r := lgrpc.ResolverMock{}
				opts := option.ResolveOptions{}
				handler.BeforeResolve(&spec, &r, &opts)
				So(handler.lbs, ShouldNotBeNil)
				So(len(handler.lbs), ShouldEqual, 1)
			})
			Convey("after resolving service", func() {
				keeper := jt.SkyLbKeeperMock{}
				keeper.On("LoopUntilAvailable").Return(true)
				keeper.On("WaitUntilReady")

				tracer := ltracing.TracerMock{}
				handler.AfterResolve(&spec, 1, "test-client", &keeper, &tracer, true)
				So(callbackCalled, ShouldBeTrue)
				So(len(handler.conns), ShouldEqual, 1)
				So(len(handler.options), ShouldEqual, 6)

				Convey("when shutdown", func() {
					handler.OnShutdown()
					So(handler.conns, ShouldBeNil)
					So(handler.lbs, ShouldBeNil)
				})
			})
		})
	})
}
