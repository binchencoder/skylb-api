package rpccli

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	stdstr "strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/binchencoder/letsgo/strings"
	"github.com/binchencoder/skylb-api/internal/flags"
	pb "github.com/binchencoder/skylb-api/proto"
	"github.com/binchencoder/skylb-api/resolver"
)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

// NewGrpcClient returns a new SkyLB grpc client.
func NewGrpcClient(ctx context.Context) (pb.SkylbClient, error) {
	if *flags.SkylbEndpoints == "" {
		glog.Error("Flag --skylb-endpoints is required.")
		os.Exit(2)
	}

	eps := strings.CsvToSlice(*flags.SkylbEndpoints)
	// if len(eps) == 0 {
	// 	glog.Fatalln("Valid flag --skylb-endpoints is required.")
	// }

	// err, ep, port := randomSkylbEndpoint(eps)
	// if err != nil {
	// 	panic(err)
	// }
	// addrs := fmt.Sprintf("%s:%s", ep, port)

	for _, ep := range eps {
		checkValidHost(ep)
	}
	addrs := *flags.SkylbEndpoints

	conn, err := grpc.Dial(resolver.DirectTarget(addrs),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(time.Second), grpc.WithBlock())
	if err != nil {
		glog.Errorf("Failed to dial to SkyLB instance %s, %+v.", addrs, err)
		return nil, err
	}

	go func(ctx context.Context, conn *grpc.ClientConn) {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			glog.Errorf("Failed to close gRPC connection, %+v", err)
		}
		glog.V(1).Infoln("SkyLB client connection closed.")
	}(ctx, conn)

	glog.V(1).Infoln("SkyLB client connection created.")
	return pb.NewSkylbClient(conn), nil
}

func randomSkylbEndpoint(eps []string) (error, string, string) {
	// Randomly pick one SkyLB gRPC instance to connect.
	var idx int
	if len(eps) > 1 {
		idx = rand.Intn(len(eps))
	}
	endpoint := eps[idx]
	glog.Infof("Picked SKyLB gRPC instance %s", endpoint)

	addrs := checkValidHost(endpoint)
	glog.Infof("Resolved SkyLB instances %s", addrs)

	// TODO(chenbin) 这里不要随机取
	parts := stdstr.SplitN(endpoint, ":", 2)
	ep, port := parts[0], parts[1]
	switch len(addrs) {
	case 0:
		return fmt.Errorf("No SkyLB instances found"), "", ""
	case 1:
		ep = addrs[0]
	default:
		idx = rand.Intn(len(addrs))
		ep = addrs[idx]
	}
	glog.Infof("Connecting SkyLB instance %s on port %s", ep, port)

	return nil, ep, port
}

func checkValidHost(endpoint string) (addres []string) {
	parts := stdstr.SplitN(endpoint, ":", 2)
	if addrs, err := net.LookupHost(parts[0]); err != nil {
		glog.Errorf("Failed to lookup SkyLB endpoint %s, %v.\n", endpoint, err)
		panic(err)
	} else {
		return addrs
	}
}
