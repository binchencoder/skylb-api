package skylb

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	pb "binchencoder.com/skylb-api/proto"
)

func TestSkylbResolver(t *testing.T) {
	Convey("Skylb resolver", t, func() {
		spec := pb.ServiceSpec{
			Namespace:   "ns",
			ServiceName: "svc",
			PortName:    "ptn",
		}
		keeper := skyLbKeeper{
			services: map[string]*serviceEntry{},
		}
		Convey("Resolve service 'test-service' and get back a watcher", func() {
			r, err := NewResolver(&keeper, &spec)
			So(err, ShouldBeNil)
			So(r, ShouldNotBeNil)

			w, err := r.Resolve("test-service")
			So(err, ShouldBeNil)
			So(w, ShouldNotBeNil)
		})
	})
}
