package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"

	"jingoal.com/letsgo"
	"jingoal.com/letsgo/service/naming"
	skycli "jingoal.com/skylb-api/client"
	"jingoal.com/skylb-api/internal/skylb"
	skysvr "jingoal.com/skylb-api/server"
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
