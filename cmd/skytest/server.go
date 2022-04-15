package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	vexpb "github.com/binchencoder/gateway-proto/data"
	"github.com/binchencoder/letsgo"
	cli "github.com/binchencoder/skylb-api/client"
	pb "github.com/binchencoder/skylb-api/cmd/skytest/proto"
	skylb "github.com/binchencoder/skylb-api/server"
)

var (
	host       = flag.String("host", "localhost", "The host")
	port       = flag.Int("port", 18000, "The gRPC port of the server")
	scrapeAddr = flag.String("scrape-addr", "0.0.0.0:18001", "The prometheus scrape port")
)

var (
	letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	token       string
	myServiceId vexpb.ServiceId
)

func init() {
	rand.Seed(time.Now().UnixNano())
	token = randString(20)
	fmt.Printf("I am %s \n", token)
}

func usage() {
	fmt.Println(`Skytest gRPC server.

Usage:
	skytest-server [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	// if *host == "" {
	// 	glog.Fatalf("Flag --host is required.")
	// }

	// go registerFakeServices()
	go registerPrometheus()

	myServiceId = vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE
	skylb.Register(myServiceId, cli.DefaultPortName, *port)
	skylb.EnableHistogram()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	glog.Infof("Starting gRPC service at %s\n", addr)
	skylb.Start(addr, func(s *grpc.Server) error {
		pb.RegisterSkytestServer(s, &greetingServer{})
		return nil
	})
}

// func registerFakeServices() {
// 	for k := range vexpb.ServiceId_name {
// 		if k > 1000 && k < 220000 {
// 			glog.Infof("Registered service %d", k)
// 			spec := skycli.NewServiceSpec(skycli.DefaultNameSpace, vexpb.ServiceId(k), "grpc")
// 			go skylb.StartSkylbReportLoadWithFixedHost(spec, *host, *port)
// 		}
// 	}
// }

func registerPrometheus() {
	http.Handle("/_/metrics", prometheus.UninstrumentedHandler())
	if err := http.ListenAndServe(*scrapeAddr, nil); err != nil {
		log.Fatal("ListenServerError:", err)
	}
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

	time.Sleep(time.Duration(10+rand.Intn(40)) * time.Millisecond)

	g := pb.GreetingResponse{
		Greeting: fmt.Sprintf("Hello %s <- %s@:%d", req.Name, token, *port),
	}
	return &g, nil
}
