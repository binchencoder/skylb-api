package main

// To run the Test*Client* tests, service "vexillary-demo" and "vexillary-test"
// must be running. (If they are not running, Start() will block)
//
// To run the Test*Service tests, client "vexillary-client" had better be
// running, so as to actually see these services being called.
//
//
// Other pre-requisites:
//
// host name "skylbserver" is reachable (e.g. set it in /etc/hosts)
// skylb server is running in host "skylbserver".

import (
	"flag"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	skylbclient "github.com/binchencoder/skylb-api/client"
	pb "github.com/binchencoder/skylb-api/cmd/demo/proto"
	"github.com/binchencoder/skylb-api/cmd/demo/rpc"
	skypb "github.com/binchencoder/skylb-api/proto"
	skylbserver "github.com/binchencoder/skylb-api/server"
	vexpb "github.com/binchencoder/gateway-proto/data"
)

var (
	namespace   = "default"
	myServiceId vexpb.ServiceId

	letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

var (
	skycli skylbclient.ServiceCli

	skycli2 skylbclient.ServiceCli

	// The gRPC client for vexillary-demo API.
	DemoCli pb.DemoClient
	// The gRPC client for vexillary-test API.
	TestCli pb.DemoClient
)

var (
	port  = flag.Int("port", 38080, "The gRPC port of the server")
	port2 = flag.Int("port2", 38082, "The gRPC port of serivce 2")
	port3 = flag.Int("port3", 38084, "The gRPC port of serivce 3")
)

func init() {
	// Set some default flags for testing.
	fmt.Println("Init default flags for tests")
	flag.Set("logtostderr", "true")
	flag.Set("v", "4")

	outsidek8s := true
	if outsidek8s {
		fmt.Println("Init default flags for outside k8s")
		flag.Set("within-k8s", "false")
		flag.Set("skylb-endpoints", "skylbserver:1900")
	}
}

// Test single client.
func TestOneClientToOne(t *testing.T) {
	skycli = skylbclient.NewServiceCli(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)
	demoSpec := skylbclient.NewServiceSpec(namespace, vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, rpc.PortName)
	skycli.Resolve(demoSpec)
	skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		switch spec.String() {
		case demoSpec.String():
			DemoCli = pb.NewDemoClient(conn)
		}
	})
}

func TestOneClientToMulti(t *testing.T) {
	skycli = skylbclient.NewServiceCli(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)

	demoSpec := skylbclient.NewServiceSpec(namespace, vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, rpc.PortName)
	skycli.Resolve(demoSpec)

	testSpec := skylbclient.NewServiceSpec(namespace, vexpb.ServiceId_VEXILLARY_TEST_SERVICE, rpc.PortName)
	skycli.Resolve(testSpec)

	// Enables histogram.
	skycli.EnableHistogram()

	skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		switch spec.String() {
		case demoSpec.String():
			DemoCli = pb.NewDemoClient(conn)
		case testSpec.String():
			TestCli = pb.NewDemoClient(conn)
		}
	})
}

func TestMultiClientToMulti(t *testing.T) {
	skycli = skylbclient.NewServiceCli(vexpb.ServiceId_AD_EXCHANGE)

	demoSpec := skylbclient.NewServiceSpec(namespace, vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, rpc.PortName)
	skycli.Resolve(demoSpec)

	skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		switch spec.String() {
		case demoSpec.String():
			DemoCli = pb.NewDemoClient(conn)
			glog.Infof("%#v", DemoCli)
		}
	})

	skycli2 = skylbclient.NewServiceCli(vexpb.ServiceId_AD_EXCHANGE)
	testSpec := skylbclient.NewServiceSpec(namespace, vexpb.ServiceId_VEXILLARY_TEST_SERVICE, rpc.PortName)
	skycli2.Resolve(testSpec)
	skycli2.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		switch spec.String() {
		case testSpec.String():
			TestCli = pb.NewDemoClient(conn)
			glog.Infof("%#v", TestCli)
		}
	})
}

func TestOneService(t *testing.T) {
	myServiceId = vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE
	skylbserver.Register(myServiceId, rpc.PortName, *port)
	skylbserver.EnableHistogram()
	go func() {
		skylbserver.Start(fmt.Sprintf(":%d", *port), func(s *grpc.Server) error {
			pb.RegisterDemoServer(s, &greetingServer{})
			return nil
		})
	}()
	time.Sleep(time.Second)
	glog.Infof("Let %s run for a while", myServiceId)
	// Let this service handle some requests and then we exit this test.
	time.Sleep(4 * time.Second)
}

func TestMultiService(t *testing.T) {
	skylbserver.RegisterMulti(vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, rpc.PortName, *port3, fmt.Sprintf(":%d", *port3))
	skylbserver.RegisterMulti(vexpb.ServiceId_VEXILLARY_TEST_SERVICE, rpc.PortName, *port2, fmt.Sprintf(":%d", *port2))

	skylbserver.EnableHistogram() // optional

	go func() {
		skylbserver.StartMulti(func(serviceId vexpb.ServiceId, s *grpc.Server) error {
			switch serviceId {
			case vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE:
				pb.RegisterDemoServer(s, &greetingServer{
					token:       randString(20),
					myServiceId: serviceId,
				})
				return nil
			case vexpb.ServiceId_VEXILLARY_TEST_SERVICE:
				pb.RegisterDemoServer(s, &greetingServer{
					token:       randString(20),
					myServiceId: serviceId,
				})
				return nil
			}
			return nil
		})
	}()
	time.Sleep(time.Second)
	glog.Infoln("Let services run for a while")
	time.Sleep(4 * time.Second)
}

func TestClientAndService(t *testing.T) {
	TestOneClientToOne(t)

	// To void service port conflict with previous TestOneService()
	*port += 99
	TestOneService(t)
}

type greetingServer struct {
	token       string
	myServiceId vexpb.ServiceId
}

func (gs greetingServer) Greeting(ctx context.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	glog.Infof("getting request from client, name %s", req.Name)

	g := pb.GreetingResponse{
		Greeting: fmt.Sprintf("Hello %s <- %s|%d", req.Name, gs.token, gs.myServiceId),
	}
	return &g, nil
}

func (gs greetingServer) GreetingForEver(req *pb.GreetingRequest, stream pb.Demo_GreetingForEverServer) error {
	glog.Infof("getting for ever request from client, name %s", req.Name)
	resp := pb.GreetingResponse{}
	for range time.Tick(time.Second) {
		resp.Greeting = fmt.Sprintf("Hello %s, from stream %s", req.Name, gs.token)
		err := stream.Send(&resp)
		if err != nil {
			glog.Errorf("Failed to greet service, %v", err)
			return err
		}
		glog.V(3).Infof("Reply: %s", resp.Greeting)
	}

	return nil
}

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
