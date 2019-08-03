package skylb

import (
	"fmt"

	"github.com/golang/glog"
	"google.golang.org/grpc/naming"

	pb "jingoal.com/skylb-api/proto"
)

// debugWatcher implements grpc naming.Watcher.
type debugWatcher struct {
	spec *pb.ServiceSpec
	ch   <-chan []*naming.Update
}

// Close closes the watcher.
func (sw *debugWatcher) Close() {
}

// Next returns the load balance updates.
func (sw *debugWatcher) Next() ([]*naming.Update, error) {
	glog.V(2).Infoln("Waiting for next SkyLB update.")
	return <-sw.ch, nil
}

// NewDebugWatcher returns a new grpc load balance watcher for debug purpose.
func NewDebugWatcher(ie *pb.InstanceEndpoint, spec *pb.ServiceSpec) naming.Watcher {
	ch := make(chan []*naming.Update, 1)
	updates := []*naming.Update{
		{
			Op:   naming.Add,
			Addr: fmt.Sprintf("%s:%d", ie.Host, ie.Port),
		},
	}
	ch <- updates

	w := &debugWatcher{
		spec: spec,
		ch:   ch,
	}
	return w
}

type debugResolver struct {
	endpoint *pb.InstanceEndpoint
	spec     *pb.ServiceSpec
}

// Resolve returns a watcher for the service name. Argument target is not used.
func (dr *debugResolver) Resolve(target string) (naming.Watcher, error) {
	glog.Infof("Request to resolve DEBUG service %s.%s on port name \"%s\".", dr.spec.Namespace, dr.spec.ServiceName, dr.spec.PortName)
	return NewDebugWatcher(dr.endpoint, dr.spec), nil
}

// NewDebugResolver returns a grpc naming.Resolver for debug purpose.
func NewDebugResolver(ie *pb.InstanceEndpoint, spec *pb.ServiceSpec) naming.Resolver {
	return &debugResolver{
		endpoint: ie,
		spec:     spec,
	}
}
