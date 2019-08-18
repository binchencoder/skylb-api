package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"binchencoder.com/letsgo"
	grpcerr "binchencoder.com/letsgo/grpc"
	pb "binchencoder.com/skylb-api/cmd/demo/proto"
	"binchencoder.com/skylb-api/cmd/demo/rpc"
	skylb "binchencoder.com/skylb-api/server"
	vexpb "binchencoder.com/gateway-proto/data"
	fepb "binchencoder.com/gateway-proto/frontend"
)

var (
	port = flag.Int("port", 8080, "The gRPC port of the server")
	test = flag.Bool("test", false, "If true register as service vexillary-test")
)

var (
	letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	token       string
	myServiceId vexpb.ServiceId
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

	// This is to reuse the same binary both as vexillary-demo and vexillary-test
	// services.
	if *test {
		myServiceId = vexpb.ServiceId_VEXILLARY_TEST_SERVICE
	} else {
		myServiceId = vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE
	}
	skylb.Register(myServiceId, rpc.PortName, *port)
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

	randNum := rand.Intn(100)
	if randNum < 5 {
		// Return error at 5% possibility.
		e := fepb.Error{
			Code:   fepb.ErrorCode_NORIGHT_ERROR,
			Params: []string{"simulateErr", "a", "b", "c"},
		}
		return nil, grpcerr.ToGrpcError(codes.InvalidArgument, &e)
	} else if randNum < 10 {
		panic("simulate panic")
	}

	g := pb.GreetingResponse{
		Greeting: fmt.Sprintf("Hello %s <- %s@:%d", req.Name, token, *port),
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
