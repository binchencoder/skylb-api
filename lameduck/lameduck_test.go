package lameduck

import (
	"testing"

	etcd "github.com/coreos/etcd/client"
	. "github.com/smartystreets/goconvey/convey"

	"binchencoder.com/skylb-api/util"
)

func TestHostPort(t *testing.T) {
	Convey("Compose gRPC endpoints", t, func() {
		So(HostPort("192.168.1.101", "4000"), ShouldEqual, "192.168.1.101#4000")
	})
}

func TestIsLameduckMode(t *testing.T) {
	lameducks = []string{"192.168.1.101#4000"}
	Convey("Verfiy endpoint 192.168.1.101#4000 is in lameduck mode", t, func() {
		So(IsLameduckMode("192.168.1.101#4000"), ShouldBeTrue)
	})
	Convey("Verfiy endpoint 192.168.1.102#4000 is not in lameduck mode", t, func() {
		So(IsLameduckMode("192.168.1.102#4000"), ShouldBeFalse)
	})
}

func TestAddLameduckEndpoint(t *testing.T) {
	lameducks = []string{}
	Convey("Verfiy endpoint 192.168.1.101#4000 is not in lameduck mode", t, func() {
		So(IsLameduckMode("192.168.1.101#4000"), ShouldBeFalse)
	})
	success := addLameduckEndpoint("192.168.1.101#4000")
	Convey("Verfiy endpoint 192.168.1.101#4000 is in lameduck mode", t, func() {
		So(IsLameduckMode("192.168.1.101#4000"), ShouldBeTrue)
		So(success, ShouldBeTrue)
	})
}

func TestRemoveLameduckEndpoint(t *testing.T) {
	lameducks = []string{"192.168.1.101#4000"}
	Convey("Verfiy endpoint 192.168.1.101#4000 is in lameduck mode", t, func() {
		So(IsLameduckMode("192.168.1.101#4000"), ShouldBeTrue)
	})
	success := removeLameduckEndpoint("192.168.1.101#4000")
	Convey("Verfiy endpoint 192.168.1.101#4000 is not in lameduck mode", t, func() {
		So(IsLameduckMode("192.168.1.101#4000"), ShouldBeFalse)
		So(success, ShouldBeTrue)
	})
}

func TestExtractLameduck(t *testing.T) {
	Convey("Verfiy endpoints are not in lameduck mode before extracting", t, func() {
		So(IsLameduckMode("192.168.1.101#2000"), ShouldBeFalse)
		So(IsLameduckMode("192.168.1.102#2000"), ShouldBeFalse)
		So(IsLameduckMode("192.168.1.201#3000"), ShouldBeFalse)
		So(IsLameduckMode("192.168.1.202#3000"), ShouldBeFalse)
	})
	root := etcd.Node{
		Key: "/grpc/lameduck/services",
		Nodes: []*etcd.Node{
			{
				Key: "/grpc/lameduck/services/account-service",
				Nodes: []*etcd.Node{
					{
						Key: "/grpc/lameduck/services/account-service/endpoints",
						Nodes: []*etcd.Node{
							{
								Key: "/grpc/lameduck/services/account-service/endpoints/192.168.1.101#2000",
							},
							{
								Key: "/grpc/lameduck/services/account-service/endpoints/192.168.1.102#2000",
							},
						},
					},
				},
			},
			{
				Key: "/grpc/lameduck/services/dory-service",
				Nodes: []*etcd.Node{
					{
						Key: "/grpc/lameduck/services/dory-service/endpoints",
						Nodes: []*etcd.Node{
							{
								Key: "/grpc/lameduck/services/dory-service/endpoints/192.168.1.201#3000",
							},
							{
								Key: "/grpc/lameduck/services/dory-service/endpoints/192.168.1.202#3000",
							},
						},
					},
				},
			},
		},
	}
	Convey("Extract lameduck endpoints from ETCD node", t, func() {
		ExtractLameduck(&root)
		So(IsLameduckMode("192.168.1.101#2000"), ShouldBeTrue)
		So(IsLameduckMode("192.168.1.102#2000"), ShouldBeTrue)
		So(IsLameduckMode("192.168.1.201#3000"), ShouldBeTrue)
		So(IsLameduckMode("192.168.1.202#3000"), ShouldBeTrue)
	})
}

func TestExtractLameduckChange(t *testing.T) {
	lameducks = []string{
		"192.168.1.103#4000",
		"192.168.1.104#4000",
	}
	Convey("Verfiy endpoints are not in lameduck mode", t, func() {
		So(IsLameduckMode("192.168.1.101#4000"), ShouldBeFalse)
		So(IsLameduckMode("192.168.1.102#4000"), ShouldBeFalse)
	})
	Convey("Verfiy endpoints are in lameduck mode", t, func() {
		So(IsLameduckMode("192.168.1.103#4000"), ShouldBeTrue)
		So(IsLameduckMode("192.168.1.104#4000"), ShouldBeTrue)
	})
	Convey("Extract lameduck endpoints from ETCD node watch change", t, func() {
		create := etcd.Response{
			Action: util.ActionCreate,
			Node: &etcd.Node{
				Key: "/grpc/lameduck/services/account-service/endpoints/192.168.1.101#4000",
			},
		}
		ExtractLameduckChange(&create)
		set := etcd.Response{
			Action: util.ActionSet,
			Node: &etcd.Node{
				Key: "/grpc/lameduck/services/account-service/endpoints/192.168.1.102#4000",
			},
		}
		ExtractLameduckChange(&set)
		delete := etcd.Response{
			Action: util.ActionDelete,
			Node: &etcd.Node{
				Key: "/grpc/lameduck/services/account-service/endpoints/192.168.1.103#4000",
			},
		}
		ExtractLameduckChange(&delete)
		expire := etcd.Response{
			Action: util.ActionExpire,
			Node: &etcd.Node{
				Key: "/grpc/lameduck/services/account-service/endpoints/192.168.1.104#4000",
			},
		}
		ExtractLameduckChange(&expire)
		So(IsLameduckMode("192.168.1.101#4000"), ShouldBeTrue)
		So(IsLameduckMode("192.168.1.102#4000"), ShouldBeTrue)
		So(IsLameduckMode("192.168.1.103#4000"), ShouldBeFalse)
		So(IsLameduckMode("192.168.1.104#4000"), ShouldBeFalse)
	})
}
