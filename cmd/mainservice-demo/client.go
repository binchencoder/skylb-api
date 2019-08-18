package main

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
	"binchencoder.com/skylb-api/balancer"
	skylb "binchencoder.com/skylb-api/client"
	"binchencoder.com/skylb-api/client/option"
	pb "binchencoder.com/skylb-api/cmd/mainservice-demo/proto"
	"binchencoder.com/skylb-api/handlers"
	vexpb "binchencoder.com/gateway-proto/data"
)

var (
	sl      skylb.ServiceLocator
	demoCli pb.DemoClient
)

func initSkylb() {
	demoSpec := skylb.NewServiceSpec(skylb.DefaultNameSpace, vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, skylb.DefaultPortName)
	sl = skylb.NewServiceLocator(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)

	grpcHandler := handlers.NewGrpcLoadBalanceHandler(demoSpec,
		func(conn *grpc.ClientConn) {
			demoCli = pb.NewDemoClient(conn)
		})
	sl.Resolve(grpcHandler,
		option.WithBalancerCreator(func(r naming.Resolver) grpc.Balancer {
			return balancer.MainService(r)
		}))

	sl.Start()
}

func usage() {
	fmt.Println(`mainservice demo client.

Usage:
	mainservice-demo [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	initSkylb()

	glog.Infoln("client finish init skylb.")

	go func() {
		for range time.Tick(3 * time.Minute) {
			req := pb.AutoSelectRequest{
				From: "client1",
			}
			ctx := context.Background()
			resp, err := demoCli.AutoSelectMain(ctx, &req)
			if err != nil {
				glog.Errorf("Failed to auto select service, %v", err)
				continue
			}
			glog.Infof("Demo Reply: %s", resp.Server)
		}
	}()

	select {}
}
