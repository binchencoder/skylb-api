package testing

import (
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/naming"

	vexpb "github.com/binchencoder/gateway-proto/data"
	pb "github.com/binchencoder/skylb-api/proto"
)

// SkyLbKeeperMock mocks SkyLB API option.SkyLbKeeper.
type SkyLbKeeperMock struct {
	mock.Mock
}

func (skm *SkyLbKeeperMock) RegisterService(spec *pb.ServiceSpec) <-chan []*naming.Update {
	args := skm.Called(spec)
	if res, ok := args.Get(0).(<-chan []*naming.Update); ok {
		return res
	}
	return nil
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
