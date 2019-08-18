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

	"binchencoder.com/letsgo"
	cli "binchencoder.com/skylb-api/client"
	pb "binchencoder.com/skylb-api/cmd/grpc-proxy-demo/proto"
	skylb "binchencoder.com/skylb-api/server"
	vexpb "binchencoder.com/gateway-proto/data"
)

var (
	port = flag.Int("port", 8901, "The gRPC port of the server")
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
	fmt.Println(`Grpc-Proxy demo server.

Usage:
	server [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	myServiceId = vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE
	skylb.Register(myServiceId, cli.DefaultPortName, *port)
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

	g := pb.GreetingResponse{
		Greeting: fmt.Sprintf("Hello %s <- %s@:%d", req.Name, token, *port),
	}
	return &g, nil
}
