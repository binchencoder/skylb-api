package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"jingoal.com/letsgo"
	cli "jingoal.com/skylb-api/client"
	pb "jingoal.com/skylb-api/cmd/mainservice-demo/proto"
	skylb "jingoal.com/skylb-api/server"
	vexpb "jingoal.com/vexillary-client/proto/data"
)

var (
	port   = flag.Int("port", 8901, "The gRPC port of the server")
	server = flag.String("server", "server1", "server no.")
)

var (
	myServiceId vexpb.ServiceId
)

func usage() {
	fmt.Println(`mainservice demo server.

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
		pb.RegisterDemoServer(s, &mainServer{})
		return nil
	})
}

type mainServer struct {
}

func (mainServer) AutoSelectMain(ctx context.Context, req *pb.AutoSelectRequest) (*pb.AutoSelectResponse, error) {
	glog.Infof("getting request from client: %s", req.From)

	e := pb.AutoSelectResponse{
		Server: fmt.Sprintf("server no: %s", *server),
	}
	return &e, nil
}
