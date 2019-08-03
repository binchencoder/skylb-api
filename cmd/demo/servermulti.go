package main

/*
Demonstrates registering multiple sub services.

This program can work standalone with client.go, unlike server.go which needs
to start two instances for demo and test service each.
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
	port  = flag.Int("port", 8082, "The gRPC port of service 1")
	port2 = flag.Int("port2", 8092, "The gRPC port of serivce 2")
)

var (
	letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	rand.Seed(time.Now().UnixNano())
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

	skylb.RegisterMulti(vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, rpc.PortName, *port, fmt.Sprintf(":%d", *port))
	skylb.RegisterMulti(vexpb.ServiceId_VEXILLARY_TEST_SERVICE, rpc.PortName, *port2, fmt.Sprintf(":%d", *port2))

	skylb.EnableHistogram() // optional

	skylb.StartMulti(func(serviceId vexpb.ServiceId, s *grpc.Server) error {
		switch serviceId {
		case vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE:
			pb.RegisterDemoServer(s, &greetingServer{
				token:       randString(20),
				myServiceId: serviceId,
			})
			return nil
		case vexpb.ServiceId_VEXILLARY_TEST_SERVICE:
			pb.RegisterDemoServer(s, &greetingServer{
				token:       randString(20),
				myServiceId: serviceId,
			})
			return nil
		}
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
	token       string
	myServiceId vexpb.ServiceId
}

func (gs greetingServer) Greeting(ctx context.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	glog.Infof("getting request from client, name %s", req.Name)

	g := pb.GreetingResponse{
		Greeting: fmt.Sprintf("Hello %s <- %s|%d", req.Name, gs.token, gs.myServiceId),
	}
	return &g, nil
}

func (gs greetingServer) GreetingForEver(req *pb.GreetingRequest, stream pb.Demo_GreetingForEverServer) error {
	glog.Infof("getting for ever request from client, name %s", req.Name)
	resp := pb.GreetingResponse{}
	for range time.Tick(time.Second) {
		resp.Greeting = fmt.Sprintf("Hello %s, from stream %s", req.Name, gs.token)
		err := stream.Send(&resp)
		if err != nil {
			glog.Errorf("Failed to greet service, %v", err)
			return err
		}
		glog.V(3).Infof("Reply: %s", resp.Greeting)
	}

	return nil
}
