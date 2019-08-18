package skylb

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/naming"

	pb "binchencoder.com/skylb-api/proto"
)

func TestSkylbWatcher(t *testing.T) {
	Convey("Skylb watcher", t, func() {
		spec := pb.ServiceSpec{
			Namespace:   "ns",
			ServiceName: "svc",
			PortName:    "ptn",
		}
		Convey("Get next update", func() {
			ch := make(chan []*naming.Update, 1)
			w := NewWatcher(ch, &spec)
			ch <- []*naming.Update{
				{Op: naming.Add, Addr: "172.0.1.100:8000"},
				{Op: naming.Delete, Addr: "172.0.1.101:8000"},
			}
			up, err := w.Next()
			So(err, ShouldBeNil)
			So(up, ShouldNotBeNil)
			So(len(up), ShouldEqual, 2)
			So(up[0].Op, ShouldEqual, naming.Add)
			So(up[0].Addr, ShouldEqual, "172.0.1.100:8000")
			So(up[1].Op, ShouldEqual, naming.Delete)
			So(up[1].Addr, ShouldEqual, "172.0.1.101:8000")
		})
	})
}
