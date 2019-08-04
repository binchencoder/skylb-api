package main

/**
Demonstrates usage of EnableFailFast(), where client will not block to wait for
service to become available in Start().
*/

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/binchencoder/letsgo"
	skylbclient "github.com/binchencoder/skylb-api/client"
	pb "github.com/binchencoder/skylb-api/cmd/demo/proto"
	skypb "github.com/binchencoder/skylb-api/proto"
	vexpb "github.com/binchencoder/ease-gateway/proto/data"
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
	skycli := skylbclient.NewServiceCli(vexpb.ServiceId_NOTIFICATION_PUSH_SERVICE)

	// Make it fail fast.
	skycli.EnableFailFast()

	// Resolve service
	demoSpec := skylbclient.NewServiceSpec("default", vexpb.ServiceId_DORY_SERVICE, "grpc")
	skycli.Resolve(demoSpec)

	skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		switch spec.String() {
		case demoSpec.String():
			// Init service client only when connection is available.
			if nil != conn {
				svcCln = pb.NewDemoClient(conn)
			}
		}
	})
	defer skycli.Shutdown()

	glog.Infoln("After skylbcli.Start")

	// Quit when connection is not available.
	if nil == svcCln {
		glog.Infoln("Not connected to service")
		return
	}

	go func() {
		for range time.Tick(2 * time.Second) {
			req := pb.GreetingRequest{
				Name: fmt.Sprintf("John Doe %d", time.Now().Second()),
			}
			resp, err := svcCln.Greeting(context.Background(), &req)
			if err != nil {
				glog.Errorf("Failed to greet service, %v", err)
				continue
			}
			glog.Infof("Demo Reply: %s", resp.Greeting)
		}
	}()

	select {}
}
