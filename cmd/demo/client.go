package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"

	"binchencoder.com/letsgo"
	"binchencoder.com/letsgo/metrics"
	pb "binchencoder.com/skylb-api/cmd/demo/proto"
	"binchencoder.com/skylb-api/cmd/demo/rpc"
)

var (
	loopInterval = flag.Duration("loop-interval", time.Second, "The interval of grpc method call looping.")
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

	go metrics.StartPrometheusServer(30001)

	// Initialize gRPC service client.
	rpc.Init()
	defer rpc.Shutdown()

	go func() {
		for range time.Tick(*loopInterval) {
			req := pb.GreetingRequest{
				Name: fmt.Sprintf("John Doe %d", time.Now().Second()),
			}
			resp, err := rpc.DemoCli.Greeting(context.Background(), &req)
			if err != nil {
				glog.Errorf("Failed to greet service, %v", err)
				continue
			}
			glog.Infof("Demo Reply: %s", resp.Greeting)
		}
	}()

	if !*rpc.FlagSingleService {
		go func() {
			for range time.Tick(*loopInterval) {
				req := pb.GreetingRequest{
					Name: fmt.Sprintf("Joe Blow %d", time.Now().Second()),
				}
				resp, err := rpc.TestCli.Greeting(context.Background(), &req)
				if err != nil {
					glog.Errorf("Failed to greet service 2, %v", err)
					continue
				}
				glog.Infof("Test Reply: %s", resp.Greeting)
			}
		}()
	}

	select {}
}
