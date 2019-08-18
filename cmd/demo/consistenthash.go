package main

/**
Demonstrates usage of consistent hashing load balancer.

You may want to start multiple server instances and restart them randomly to
check out the effect.
*/

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"

	"binchencoder.com/letsgo"
	"binchencoder.com/letsgo/hashring"
	"binchencoder.com/skylb-api/balancer"
	skylbclient "binchencoder.com/skylb-api/client"
	"binchencoder.com/skylb-api/client/option"
	pb "binchencoder.com/skylb-api/cmd/demo/proto"
	skypb "binchencoder.com/skylb-api/proto"
	vexpb "binchencoder.com/gateway-proto/data"
)

var (
	svcCln pb.DemoClient
)

func usage() {
	fmt.Println(`SkyLB demo client.

Usage:
	client [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	// Initialize gRPC service client.
	skycli := skylbclient.NewServiceCli(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)

	// Resolve service
	demoSpec := skylbclient.NewServiceSpec("default",
		vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, "grpc")
	skycli.Resolve(demoSpec,
		// Enable consistent hashing load balancer.
		option.WithBalancerCreator(func(r naming.Resolver) grpc.Balancer {
			return balancer.ConsistentHashing(r)
		}),
	)

	skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		switch spec.String() {
		case demoSpec.String():
			svcCln = pb.NewDemoClient(conn)
		}
	})
	defer skycli.Shutdown()

	glog.Infoln("After skylbcli.Start")

	for range time.Tick(time.Second) {
		req := pb.GreetingRequest{
			Name: fmt.Sprintf("John Doe %d", time.Now().Second()),
		}
		// Provide key value for hashing.
		ctx := hashring.WithHashKey(context.Background(), "valueOfYourKey")
		resp, err := svcCln.Greeting(ctx, &req)
		if err != nil {
			glog.Errorf("Failed to greet service, %v", err)
			continue
		}
		glog.Infof("Demo Reply: %s", resp.Greeting)
	}
}
