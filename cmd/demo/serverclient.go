package main

/**
Demonstrates a grpc service, which also acts as a grpc client calling
other services.

In this example, it serves as ServiceId_SHARED_TEST_SERVER_SERVICE, and calls
ServiceId_VEXILLARY_TEST_SERVICE, so server.go should be started as well,
if you want to run client.go for testing.
*/

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"jingoal.com/letsgo"
	pb "jingoal.com/skylb-api/cmd/demo/proto"
	"jingoal.com/skylb-api/cmd/demo/rpc"
	skylb "jingoal.com/skylb-api/server"
	vexpb "jingoal.com/vexillary-client/proto/data"
)

var (
	port = flag.Int("port", 8090, "The gRPC port of the server")
)

var (
	letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	token       string
)

func init() {
	rand.Seed(time.Now().UnixNano())
	token = randString(20)
	fmt.Printf("I am %s\n", token)
}

func usage() {
	fmt.Println(`SkyLB demo server.

Usage:
	server [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	// Initialize gRPC service client.
	rpc.Init1()
	defer rpc.Shutdown()

	// Initialize gRPC service server.
	skylb.Register(vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, rpc.PortName, *port)

	// Enables histogram for server api.
	skylb.EnableHistogram()

	skylb.Start(fmt.Sprintf(":%d", *port), func(s *grpc.Server) error {
		pb.RegisterDemoServer(s, &greetingServer{})
		return nil
	})
}

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type greetingServer struct {
}

func (greetingServer) Greeting(ctx context.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	glog.Infof("getting request from client, name %s", req.Name)

	myReq := pb.GreetingRequest{
		Name: fmt.Sprintf("%s <- %s|%d", req.Name, token, vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE),
	}
	glog.Infof("Calling %#v", rpc.DemoCli)
	myResp, err := rpc.DemoCli.Greeting(context.Background(), &myReq)
	reply := "N/A"
	if err != nil {
		glog.Errorf("Failed to call service, %v", err)
	} else {
		reply = myResp.Greeting
	}

	g := pb.GreetingResponse{
		Greeting: reply,
	}
	return &g, nil
}

func (greetingServer) GreetingForEver(req *pb.GreetingRequest, stream pb.Demo_GreetingForEverServer) error {
	glog.Infof("getting for ever request from client, name %s", req.Name)
	resp := pb.GreetingResponse{}
	for range time.Tick(time.Second) {
		resp.Greeting = fmt.Sprintf("Hello %s, from stream %s", req.Name, token)
		err := stream.Send(&resp)
		if err != nil {
			glog.Errorf("Failed to greet service, %v", err)
			return err
		}
		glog.V(3).Infof("Reply: %s", resp.Greeting)
	}

	return nil
}
