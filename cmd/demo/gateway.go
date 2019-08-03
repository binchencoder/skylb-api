package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"jingoal.com/janus/gateway/runtime"
	"jingoal.com/letsgo"
	gw "jingoal.com/skylb-api/cmd/demo/proto"
)

func usage() {
	fmt.Println(`SkyLB demo gateway.

Usage:
	gateway [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

var (
	endpoint = flag.String("endpoint", "localhost:8080", "endpoint of skylb-api demo server")
)

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	mux := runtime.NewServeMux()
	err := gw.RegisterDemoHandlerFromEndpoint(mux)
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":10000", mux); err != nil {
		panic(err)
	}
}
