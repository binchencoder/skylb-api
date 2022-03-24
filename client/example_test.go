package client

import (
	"google.golang.org/grpc"

	vexpb "github.com/binchencoder/gateway-proto/data"
	pb "github.com/binchencoder/skylb-apiv2/proto"
)

func ExampleNewServiceLocator() {
	// Create a service client.
	cli := NewServiceCli(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)

	// Resolve services.
	// grpcHandler := handlers.NewGrpcLoadBalanceHandler(
	// 	NewDefaultServiceSpec(vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE),
	// 	func(conn *grpc.ClientConn) {
	// 		// hold the connecton for use later.
	// 	},
	// )
	// sl.Resolve(grpcHandler)

	cli.Resolve(NewDefaultServiceSpec(vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE))

	// Start the locator.
	// sl.Start()

	cli.Start(func(spec *pb.ServiceSpec, conn *grpc.ClientConn) {
		// hold the connecton for use later.
	})

	// Use the connection to create grpc clients.

	// Shutdown before exit.
	cli.Shutdown()
}
