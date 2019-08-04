package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	prom "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	hpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/binchencoder/letsgo"
	skylb "github.com/binchencoder/skylb-api/client"
	pb "github.com/binchencoder/skylb-api/cmd/skytest/proto"
	"github.com/binchencoder/skylb-api/handlers"
	skypb "github.com/binchencoder/skylb-api/proto"
	vexpb "github.com/binchencoder/ease-gateway/proto/data"
)

var (
	concurrency    = flag.Int("concurrency", 10, "The number of concurrent workers for each service")
	enableStresser = flag.Bool("enable-stresser", false, "True to enable stresser")
	nStresser      = flag.Int("n-stresser", 100, "The number of stressers")
	nBatchRequest  = flag.Int("n-batch-request", 10000, "The number of batched request")
	requestSleep   = flag.Duration("request-sleep", 50*time.Millisecond, "The sleep time after each request")
	requestTimeout = flag.Duration("request-timeout", 100*time.Millisecond, "The timeout of each request")
	scrapeAddr     = flag.String("scrape-addr", "0.0.0.0:18005", "The prometheus scrape port")

	spec = skylb.NewServiceSpec(skylb.DefaultNameSpace, vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE, skylb.DefaultPortName)

	grpcFailCount = prom.NewCounter(
		prom.CounterOpts{
			Namespace: "skytest",
			Subsystem: "client",
			Name:      "grpc_call_failure",
			Help:      "The number of failed gRPC calls.",
		},
	)
)

func startStresser() {
	for i := 0; i < *nStresser; i++ {
		go func() {
			for {
				sl := skylb.NewServiceLocator(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)

				grpcHandler := handlers.NewGrpcLoadBalanceHandler(spec,
					func(conn *grpc.ClientConn) {
					})
				sl.Resolve(grpcHandler)
				sl.EnableHistogram()
				sl.Start()
				time.Sleep(1 * time.Hour)
				sl.Shutdown()
			}
		}()
	}
}

func startOldSkylb(sid vexpb.ServiceId) (skylb.ServiceCli, pb.SkytestClient, hpb.HealthClient) {
	skycli := skylb.NewServiceCli(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)
	skycli.Resolve(skylb.NewServiceSpec(skylb.DefaultNameSpace, sid, skylb.DefaultPortName))
	skycli.EnableHistogram()
	var cli pb.SkytestClient
	var healthCli hpb.HealthClient
	skycli.Start(func(spec *skypb.ServiceSpec, conn *grpc.ClientConn) {
		cli = pb.NewSkytestClient(conn)
		healthCli = hpb.NewHealthClient(conn)
	})
	return skycli, cli, healthCli
}

func startNewSkylb(sid vexpb.ServiceId) (skylb.ServiceLocator, pb.SkytestClient, hpb.HealthClient) {
	var cli pb.SkytestClient
	var healthCli hpb.HealthClient
	sl := skylb.NewServiceLocator(vexpb.ServiceId_SHARED_TEST_CLIENT_SERVICE)
	sl.EnableHistogram()

	grpcHandler := handlers.NewGrpcLoadBalanceHandler(
		skylb.NewServiceSpec(skylb.DefaultNameSpace, sid, skylb.DefaultPortName),
		func(conn *grpc.ClientConn) {
			cli = pb.NewSkytestClient(conn)
			healthCli = hpb.NewHealthClient(conn)
		})
	sl.Resolve(grpcHandler)
	sl.Start()
	return sl, cli, healthCli
}

func usage() {
	fmt.Println(`Skytest gRPC client.

Usage:
	skytest-client [options]

Options:`)

	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	letsgo.Init(letsgo.FlagUsage(usage))

	if *enableStresser {
		startStresser()
	}
	testOldAPI()
	testNewAPI()

	http.Handle("/_/metrics", prometheus.UninstrumentedHandler())
	if err := http.ListenAndServe(*scrapeAddr, nil); err != nil {
		log.Fatal("ListenServerError:", err)
	}
}

var sids = []vexpb.ServiceId{
	vexpb.ServiceId_SHARED_TEST_SERVER_SERVICE,
	vexpb.ServiceId_NETDISK_SERVICE,
	vexpb.ServiceId_CRM_SERVICE,
	vexpb.ServiceId_GIGAFORM_SERVICE,
	vexpb.ServiceId_IGOAL_GRPC_SERVICE,
	vexpb.ServiceId_IDM_SERVICE,
	vexpb.ServiceId_EIM_SERVER,
	vexpb.ServiceId_INDEXER_SERVICE,
	vexpb.ServiceId_PRUINAE_AVATAR_SERVICE,
	vexpb.ServiceId_CRON_SERVER,
	vexpb.ServiceId_WEB_COMMON_SERVICE,
	vexpb.ServiceId_WORKBENCH_SERVICE,
	vexpb.ServiceId_MERGED_SERVER,
}

func testOldAPI() {
	for i := 0; i < len(sids); i++ {
		for j := 0; j < *concurrency; j++ {
			go func(sid vexpb.ServiceId) {
				for {
					sl, cli, healthCli := startOldSkylb(sid)
					for i := 0; i < *nBatchRequest; i++ {
						req := pb.GreetingRequest{
							Name: fmt.Sprintf("John Doe %d", time.Now().Second()),
						}
						ctx, cancel := context.WithTimeout(context.Background(), *requestTimeout)
						resp, err := cli.Greeting(ctx, &req, grpc.FailFast(false))
						if err != nil {
							cancel()
							glog.Errorf("Failed to greet service, %v", err)
							fmt.Printf("Failed to greet service, %v", err)
							grpcFailCount.Inc()
							time.Sleep(*requestTimeout)
							continue
						}
						_ = resp

						healthCli.Check(context.Background(), &hpb.HealthCheckRequest{})
						time.Sleep(*requestSleep)
					}
					sl.Shutdown()
				}
			}(sids[i])
		}
	}
}

func testNewAPI() {
	for i := 0; i < len(sids); i++ {
		for j := 0; j < *concurrency; j++ {
			go func(sid vexpb.ServiceId) {
				for {
					sl, cli, healthCli := startNewSkylb(sid)
					for i := 0; i < *nBatchRequest; i++ {
						req := pb.GreetingRequest{
							Name: fmt.Sprintf("John Doe %d", time.Now().Second()),
						}
						ctx, cancel := context.WithTimeout(context.Background(), *requestTimeout)
						resp, err := cli.Greeting(ctx, &req, grpc.FailFast(false))
						if err != nil {
							cancel()
							glog.Errorf("Failed to greet service, %v", err)
							fmt.Printf("Failed to greet service, %v", err)
							grpcFailCount.Inc()
							time.Sleep(*requestTimeout)
							continue
						}
						// fmt.Printf("Demo Reply: %s\n", resp.Greeting)
						_ = resp
						fmt.Printf(".")

						healthCli.Check(context.Background(), &hpb.HealthCheckRequest{})
						time.Sleep(*requestSleep)
					}
					sl.Shutdown()
				}
			}(sids[i])
		}
	}
}
