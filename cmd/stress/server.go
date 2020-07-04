package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	vexpb "github.com/binchencoder/gateway-proto/data"
	"github.com/binchencoder/letsgo"
	jg "github.com/binchencoder/letsgo/grpc"
	"github.com/binchencoder/letsgo/trace"
	cli "github.com/binchencoder/skylb-api/client"
	pb "github.com/binchencoder/skylb-api/cmd/stress/proto"
	skylb "github.com/binchencoder/skylb-api/server"
)

var (
	port = flag.Int("port", 5888, "The gRPC port of the server")
)

func usage() {
	fmt.Println(`SkyLB stress test server.

Usage:
	server [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	skylb.Register(vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, cli.DefaultPortName, *port)
	skylb.Start(fmt.Sprintf(":%d", *port), func(s *grpc.Server) error {
		pb.RegisterStressServiceServer(s, &stressServer{})
		return nil
	})
}

type stressServer struct {
}

func (s *stressServer) SayHello(ctx context.Context, req *pb.SayHelloReq) (*pb.SayHelloResp, error) {
	glog.V(5).Infof("getting request from client, name %s, context %+v", req.Name, ctx)

	svcId, err := jg.GetServiceId(ctx)
	if err != nil {
		return nil, err
	}

	tid, ok := trace.GetTraceId(ctx)
	if !ok {
		return nil, fmt.Errorf("Failed to get trace ID from context, %+v", ctx)
	}

	g := pb.SayHelloResp{
		Greeting:  fmt.Sprintf("Hello %s", req.Name),
		Peer:      trace.HeaderFromContext(ctx).Uid,
		ServiceId: int32(svcId),
		Tid:       tid,
	}
	return &g, nil
}

func (s *stressServer) SayHelloDisabled(ctx context.Context, req *pb.SayHelloReq) (*pb.SayHelloResp, error) {
	return s.SayHello(ctx, req)
}
