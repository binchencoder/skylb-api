package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/binchencoder/letsgo"
	skylb "github.com/binchencoder/skylb-api/client"
	pb "github.com/binchencoder/skylb-api/cmd/demo/proto"
	"github.com/binchencoder/skylb-api/cmd/demo/rpc"
	skypb "github.com/binchencoder/skylb-api/proto"
	vexpb "github.com/binchencoder/ease-gateway/proto/data"
)

func usage() {
	fmt.Println(`SkyLB demo client.

Usage:
	client [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

var (
	skycli  skylb.ServiceCli
	demoCli pb.DemoClient

	namespace = "default"
)

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	// Initialize gRPC service client.
	skycli = skylb.NewServiceCli(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)
	defer skycli.Shutdown()

	// Resolve service "vexillary-demo".
	demoSpec := skylb.NewServiceSpec(namespace, vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, rpc.PortName)
	skycli.Resolve(demoSpec)

	skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		switch spec.String() {
		case demoSpec.String():
			demoCli = pb.NewDemoClient(conn)
		}
	})

	go func() {
		req := pb.GreetingRequest{
			Name: "balancer",
		}

		for {
			stream, err := demoCli.GreetingForEver(context.Background(), &req)
			if err != nil {
				glog.Warningf("Failed to create stream call %v.", err)
				continue
			}

			for {
				resp, err := stream.Recv()
				if err != nil {
					glog.Warningf("Failed to received from server %v", err)
					break
				}

				glog.Infof("Response %s", resp.Greeting)
			}
		}
	}()

	select {}
}
