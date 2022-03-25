package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	pb "github.com/binchencoder/skylb-apiv2/cmd/skytest/proto"
)

var (
	grpcEndpoint = flag.String("grpc-endpoint", "192.168.221.104:18000", "The gRPC server endpoint")
	timeout      = flag.Duration("timeout", time.Second, "The timeout to call gRPC service")
)

func main() {
	flag.Parse()

	for {
		testDirectGrpc()
		time.Sleep(100 * time.Millisecond)
	}
}

func testDirectGrpc() {
	conn, err := grpc.Dial(*grpcEndpoint, grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	cli := pb.NewSkytestClient(conn)

	req := pb.GreetingRequest{
		Name: fmt.Sprintf("John Doe %d", time.Now().Second()),
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	// fuyc: "cancel" may not be used on all execution paths. Safely ignore
	// warnings from "go tool vet".
	resp, err := cli.Greeting(ctx, &req, grpc.FailFast(false))
	if err != nil {
		cancel()
		glog.Errorf("Failed to greet service, %v", err)
		fmt.Printf("Failed to greet service, %v", err)
		time.Sleep(100 * time.Millisecond)
		return
	}
	fmt.Printf("Demo Reply: %s\n", resp.Greeting)
}
