package rpc

import (
	"flag"

	"google.golang.org/grpc"

	vexpb "github.com/binchencoder/gateway-proto/data"
	skylb "github.com/binchencoder/skylb-api/client"
	pb "github.com/binchencoder/skylb-api/cmd/demo/proto"
	skypb "github.com/binchencoder/skylb-api/proto"
)

const (
	PortName = "grpc"
)

var (
	FlagSingleService = flag.Bool("single-service", true, "Whether depend on one single service only")

	namespace = "default"
)

var (
	skycli skylb.ServiceCli

	// The gRPC client.
	DemoCli pb.DemoClient
	// The gRPC client pointing to another service spec.
	TestCli pb.DemoClient
)

// Init initializes the gRPC client.
func Init() {
	skycli = skylb.NewServiceCli(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)

	// Resolve service
	demoSpec := skylb.NewServiceSpec(namespace, vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, PortName)
	skycli.Resolve(demoSpec)

	// Resolve another service
	testSpec := skylb.NewServiceSpec(namespace, vexpb.ServiceId_CUSTOM_EASE_GATEWAY_TEST, PortName)

	if !*FlagSingleService {
		skycli.Resolve(testSpec)
	}

	// Enable histogram.
	skycli.EnableHistogram()

	// Enable the new SkyLB protocol.
	skycli.EnableResolveFullEps()

	skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		switch spec.String() {
		case demoSpec.String():
			DemoCli = pb.NewDemoClient(conn)
		case testSpec.String():
			TestCli = pb.NewDemoClient(conn)
		}
	})
}

// Init1 initializes the gRPC client.
// Used by serverclient.go only.
func Init1() {
	skycli = skylb.NewServiceCli(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)

	demoSpec := skylb.NewServiceSpec(namespace, vexpb.ServiceId_CUSTOM_EASE_GATEWAY_TEST, PortName)
	skycli.Resolve(demoSpec)

	// Enables histogram for client api.
	skycli.EnableHistogram()

	skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		switch spec.String() {
		case demoSpec.String():
			DemoCli = pb.NewDemoClient(conn)
		}
	})
}

// Shutdown turns off the client.
func Shutdown() {
	skycli.Shutdown()
}
