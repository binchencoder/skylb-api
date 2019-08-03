package client

import (
	"google.golang.org/grpc"

	"jingoal.com/skylb-api/handlers"
	vexpb "jingoal.com/vexillary-client/proto/data"
)

func ExampleNewServiceLocator() {
	// Create a service locator.
	sl := NewServiceLocator(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)

	// Resolve services.
	grpcHandler := handlers.NewGrpcLoadBalanceHandler(
		NewDefaultServiceSpec(vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE),
		func(conn *grpc.ClientConn) {
			// hold the connecton for use later.
		},
	)
	sl.Resolve(grpcHandler)

	// Start the locator.
	sl.Start()

	// Use the connection to create grpc clients.

	// Shutdown before exit.
	sl.Shutdown()
}
