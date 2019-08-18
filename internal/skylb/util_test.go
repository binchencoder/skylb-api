package skylb

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/naming"

	pb "binchencoder.com/skylb-api/proto"
)

func TestCalcServiceKey(t *testing.T) {
	Convey("Calculate service key from service spec", t, func() {
		spec := pb.ServiceSpec{
			Namespace:   "ns",
			ServiceName: "svc",
			PortName:    "ptn",
		}
		key := calcServiceKey(&spec)
		So(key, ShouldEqual, "ns.svc:ptn")
	})
}

func TestToNamingOp(t *testing.T) {
	Convey("Convert proto Operation enums to grpc naming.Operation", t, func() {
		Convey("Convert Operation_Add", func() {
			op := toNamingOp(pb.Operation_Add)
			So(op, ShouldEqual, naming.Add)
		})
		Convey("Convert Operation_Delete", func() {
			op := toNamingOp(pb.Operation_Delete)
			So(op, ShouldEqual, naming.Delete)
		})
		Convey("Convert invalid operation", func() {
			op := toNamingOp(pb.Operation(100))
			So(op, ShouldEqual, naming.Add)
		})
	})
}

func TestOpToString(t *testing.T) {
	Convey("Convert grpc naming.Operation to string", t, func() {
		Convey("Convert naming.Add", func() {
			op := opToString(naming.Add)
			So(op, ShouldEqual, "ADD")
		})
		Convey("Convert naming.Delete", func() {
			op := opToString(naming.Delete)
			So(op, ShouldEqual, "DELETE")
		})
		Convey("Convert invalid operation", func() {
			op := opToString(100)
			So(op, ShouldEqual, "")
		})
	})
}
