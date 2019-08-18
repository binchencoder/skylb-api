package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"

	"binchencoder.com/letsgo"
	"binchencoder.com/letsgo/service/naming"
	skycli "binchencoder.com/skylb-api/client"
	"binchencoder.com/skylb-api/internal/skylb"
	skysvr "binchencoder.com/skylb-api/server"
)

func usage() {
	fmt.Println(`Fake service gRPC server.

Usage:
	fake-service [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	for svcName, ep := range skycli.DebugSvcEndpoints {
		glog.Infof("service name: %s, endpoint: %s", svcName, ep)
		ie := skylb.ParseEndpoint(ep)
		sid, err := naming.ServiceNameToId(svcName)
		if nil != err {
			glog.Errorf("Unrecognized service name: %s", svcName)
			continue
		}
		spec := skycli.NewServiceSpec("", sid, "")
		go skysvr.StartSkylbReportLoadWithFixedHost(spec, ie.Host, int(ie.Port))
	}

	select {}
}
