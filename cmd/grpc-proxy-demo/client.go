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

	"github.com/binchencoder/letsgo"
	"github.com/binchencoder/letsgo/hashring"
	"github.com/binchencoder/skylb-api/balancer"
	skylb "github.com/binchencoder/skylb-api/client"
	"github.com/binchencoder/skylb-api/client/option"
	pb "github.com/binchencoder/skylb-api/cmd/grpc-proxy-demo/proto"
	"github.com/binchencoder/skylb-api/handlers"
	vexpb "github.com/binchencoder/gateway-proto/data"
)

var (
	flagConsistentHashing = flag.Bool("enable-consistent-hashing", false, "True to enable consistent hashing load balancer")

	sl      skylb.ServiceLocator
	demoCli pb.DemoClient
)

func initSkylb() {
	demoSpec := skylb.NewServiceSpec(skylb.DefaultNameSpace, vexpb.ServiceId_GRPC_PROXY, skylb.DefaultPortName)
	sl = skylb.NewServiceLocator(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)

	grpcHandler := handlers.NewGrpcLoadBalanceHandler(demoSpec,
		func(conn *grpc.ClientConn) {
			demoCli = pb.NewDemoClient(conn)
		})
	if *flagConsistentHashing {
		sl.Resolve(grpcHandler,
			option.WithBalancerCreator(func(r naming.Resolver) grpc.Balancer {
				return balancer.ConsistentHashing(r)
			}))
	} else {
		sl.Resolve(grpcHandler)
	}

	sl.Start()
}

func usage() {
	fmt.Println(`Grpc-Proxy demo client.

Usage:
	grpc-proxy-demo [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	initSkylb()

	go func() {
		for range time.Tick(5 * time.Second) {
			req := pb.GreetingRequest{
				Name: fmt.Sprintf("John Doe %d", time.Now().Second()),
			}
			ctx := context.Background()
			if *flagConsistentHashing {
				ctx = hashring.NewHashKey(ctx)
			}
			resp, err := demoCli.Greeting(ctx, &req)
			if err != nil {
				glog.Errorf("Failed to greet service, %v", err)
				continue
			}
			glog.Infof("Demo Reply: %s", resp.Greeting)
		}
	}()

	select {}
}
