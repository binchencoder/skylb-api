package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/mwitkow/grpc-proxy/proxy"
	"github.com/soheilhy/cmux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/naming"

	"binchencoder.com/letsgo"
	lmetrics "binchencoder.com/letsgo/metrics"
	"binchencoder.com/letsgo/runtime/pprof"
	"binchencoder.com/skylb-api/balancer"
	skycli "binchencoder.com/skylb-api/client"
	"binchencoder.com/skylb-api/client/option"
	"binchencoder.com/skylb-api/handlers"
	skysrv "binchencoder.com/skylb-api/server"
	vexpb "binchencoder.com/gateway-proto/data"
)

// grpc-proxy is a generic proxy server for SkyLB grpc services.
//
// It's main use case is for clients outside of Kubernetes to call gRPC
// services within Kubernetes. Since gRPC uses long connections while
// Kubernetes uses proxy to load balance (which only happens when the
// connection is established), the direct long connections from outside
// to inside would cause unbalanced traffic to service instances.
//
// grpc-proxy solves the above issue by converting long connections to
// short connections during proxy.
//
// grpc-proxy should run inside Kubernetes and be exposed to outside.
// Traffic to grpc-proxy instances might be unbalanced but the traffic to
// backend gRPC service instances will be balanced.

var (
	host                    = flag.String("host", "", "The grpc proxy host")
	port                    = flag.Int("port", 3101, "The grpc proxy port")
	namespace               = flag.String("namespace", skycli.DefaultNameSpace, "The service namespace")
	portName                = flag.String("port-name", skycli.DefaultPortName, "The service port name")
	targetService           = flag.String("target-service", "", "The target service name")
	sendHeartbeat           = flag.Bool("send-heartbeat", false, "True to send heartbeat to SkyLB")
	enableConsistentHashing = flag.Bool("enable-consistent-hashing", false, "True to enable consistent hashing load balancer")

	sl skycli.ServiceLocator
)

func main() {
	letsgo.Init()

	cc := startSkylbAPI()

	director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		// Copy the inbound metadata explicitly.
		outCtx, _ := context.WithCancel(ctx)
		outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
		return outCtx, cc, nil
	}

	server := grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)))

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Run proxy service at %s:%d\n", *host, *port)

	if *sendHeartbeat {
		go registerToSkylb()
	}

	m := cmux.New(lis)

	// Match connections in order: first grpc, then HTTP.
	grpcl := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpl := m.Match(cmux.HTTP1Fast())

	go startHTTPServer(httpl, http.DefaultServeMux)
	go startGrpcServer(grpcl, server)

	m.Serve()
}

func startSkylbAPI() *grpc.ClientConn {
	sl = skycli.NewServiceLocator(vexpb.ServiceId_GRPC_PROXY)

	spec := skycli.NewServiceSpec(*namespace, vexpb.ServiceId(vexpb.ServiceId_value[*targetService]), *portName)
	var cc *grpc.ClientConn
	handler := handlers.NewGrpcLoadBalanceHandler(spec, func(conn *grpc.ClientConn) {
		cc = conn
	}, grpc.WithCodec(proxy.Codec()))

	options := []option.ResolveOption{}
	if *enableConsistentHashing {
		options = append(options, option.WithBalancerCreator(func(r naming.Resolver) grpc.Balancer {
			return balancer.ConsistentHashing(r)
		}))
	}

	sl.Resolve(handler, options...)
	sl.Start()
	return cc
}

func registerToSkylb() {
	spec := skycli.NewServiceSpec(*namespace, vexpb.ServiceId_GRPC_PROXY, *portName)
	skysrv.StartSkylbReportLoad(spec, *port)
}

func startHTTPServer(lis net.Listener, mux *http.ServeMux) {
	lmetrics.EnablePrometheus(mux)
	pprof.EnablePprof(mux)
	log.Printf("Start grpc server at %s:%d\n", *host, *port)
	if err := http.Serve(lis, mux); err != nil {
		panic(err)
	}
}

func startGrpcServer(grpcl net.Listener, server *grpc.Server) {
	log.Printf("Start http server at %s:%d\n", *host, *port)
	if err := server.Serve(grpcl); err != nil {
		panic(err)
	}
}
