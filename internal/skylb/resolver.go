package skylb

import (
	"github.com/golang/glog"
	"google.golang.org/grpc/naming"

	"github.com/binchencoder/skylb-api/client/option"
	pb "github.com/binchencoder/skylb-api/proto"
)

// skylbResolver implements grpc naming.Resolver.
type skylbResolver struct {
	keeper   option.SkyLbKeeper
	spec     *pb.ServiceSpec
	updateCh <-chan []*naming.Update
}

// Resolve returns a watcher for the service name. Argument target is not used.
func (sr *skylbResolver) Resolve(target string) (naming.Watcher, error) {
	glog.Infof("Request to resolve service %s.%s on port name \"%s\".", sr.spec.Namespace, sr.spec.ServiceName, sr.spec.PortName)
	return NewWatcher(sr.updateCh, sr.spec), nil
}

// NewResolver returns a grpc naming.Resolver for the given service.
func NewResolver(keeper option.SkyLbKeeper, spec *pb.ServiceSpec) (naming.Resolver, error) {
	return &skylbResolver{
		keeper:   keeper,
		spec:     spec,
		updateCh: keeper.RegisterService(spec),
	}, nil
}
