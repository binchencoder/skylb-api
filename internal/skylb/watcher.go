package skylb

import (
	"errors"

	"github.com/golang/glog"
	"google.golang.org/grpc/naming"

	pb "binchencoder.com/skylb-api/proto"
)

// skylbWatcher implements grpc naming.Watcher.
type skylbWatcher struct {
	spec      *pb.ServiceSpec
	updatesCh <-chan []*naming.Update
}

func (sw *skylbWatcher) Close() {
}

// Next returns the load balance updates.
func (sw *skylbWatcher) Next() ([]*naming.Update, error) {
	glog.V(2).Infoln("Waiting for next SkyLB update.")
	if updates, ok := <-sw.updatesCh; ok {
		return updates, nil
	}
	return nil, errors.New("SkyLB watcher update channel has been closed.")
}

// NewWatcher returns a new grpc load balance watcher.
func NewWatcher(updateCh <-chan []*naming.Update, spec *pb.ServiceSpec) naming.Watcher {
	return &skylbWatcher{
		spec:      spec,
		updatesCh: updateCh,
	}
}
