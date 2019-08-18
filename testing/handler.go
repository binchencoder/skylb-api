package testing

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/naming"

	"binchencoder.com/skylb-api/client/option"
	pb "binchencoder.com/skylb-api/proto"
	vexpb "binchencoder.com/gateway-proto/data"
)

// LoadBalanceHandlerMock mocks SkyLB API option.LoadBalanceHandler.
type LoadBalanceHandlerMock struct {
	mock.Mock
}

func (lbhm *LoadBalanceHandlerMock) ServiceSpec() *pb.ServiceSpec {
	args := lbhm.Called()
	if res, ok := args.Get(0).(*pb.ServiceSpec); ok {
		return res
	}
	return nil
}

func (lbhm LoadBalanceHandlerMock) BeforeResolve(spec *pb.ServiceSpec, r naming.Resolver, ropts *option.ResolveOptions) {
	lbhm.Called(spec, r, ropts)
}

func (lbhm LoadBalanceHandlerMock) AfterResolve(spec *pb.ServiceSpec, csId vexpb.ServiceId, csName string, keeper option.SkyLbKeeper, tracer opentracing.Tracer, failFast bool) {
	lbhm.Called(spec, csId, csName, keeper, tracer, failFast)
}

func (lbhm LoadBalanceHandlerMock) OnShutdown() {
	lbhm.Called()
}
