package testing

import (
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/resolver"

	vexpb "github.com/binchencoder/gateway-proto/data"
	pb "github.com/binchencoder/skylb-apiv2/proto"
)

// SkyLbKeeperMock mocks SkyLB API option.SkyLbKeeper.
type SkyLbKeeperMock struct {
	mock.Mock
}

func (skm *SkyLbKeeperMock) RegisterService(spec *pb.ServiceSpec, cliConn resolver.ClientConn) {
	skm.Called(spec, cliConn)
}

func (skm *SkyLbKeeperMock) Start(csId vexpb.ServiceId, csName string, resolveFullEps bool) {
	skm.Called(csId, csName, resolveFullEps)
}

func (skm *SkyLbKeeperMock) WaitUntilReady() {
	skm.Called()
}

func (skm *SkyLbKeeperMock) Shutdown() {
	skm.Called()
}
