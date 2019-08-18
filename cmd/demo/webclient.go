package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"binchencoder.com/letsgo"
	skylb "binchencoder.com/skylb-api/client"
	pb "binchencoder.com/skylb-api/cmd/demo/proto"
	"binchencoder.com/skylb-api/cmd/demo/rpc"
	skypb "binchencoder.com/skylb-api/proto"
	vexpb "binchencoder.com/gateway-proto/data"
)

const (
	defaultNameSpace = "default"
)

var (
	hostAddr = flag.String("lsn-address", ":9000", "listen address")
)

var (
	skycli  skylb.ServiceCli
	demoCli pb.DemoClient
)

func usage() {
	fmt.Println(`SkyLB demo web client.

Usage:
	webclient [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func initSkyCli() {
	// Initialize gRPC service client.
	skycli = skylb.NewServiceCli(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)

	// Resolve service "vexillary-demo".
	demoSpec := skylb.NewServiceSpec(defaultNameSpace, vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, rpc.PortName)
	skycli.Resolve(demoSpec)

	skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		switch spec.String() {
		case demoSpec.String():
			demoCli = pb.NewDemoClient(conn)
		}
	})
}

func stop() {
	skycli.Shutdown()
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	name := values.Get("name")
	if len(name) == 0 {
		name = "SkyLB"
	}

	glog.V(3).Infof("Received a client request, name %s, from %s", name, r.RemoteAddr)

	req := &pb.GreetingRequest{
		Name: name,
	}

	resp, err := demoCli.Greeting(context.Background(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := io.WriteString(w, resp.Greeting); err != nil {
		glog.Error("Failed to write response,", err)
	}
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	initSkyCli()

	http.HandleFunc("/demo", handleRequest)

	glog.Infof("Starting web server at %s\n", *hostAddr)
	if err := http.ListenAndServe(*hostAddr, nil); err != nil {
		glog.Fatal("ListenAndServe: ", err)
	}

	stop()
}
