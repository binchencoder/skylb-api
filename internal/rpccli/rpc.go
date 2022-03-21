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

	"github.com/binchencoder/letsgo/strings"
	"github.com/binchencoder/skylb-apiv2/internal/flags"
	pb "github.com/binchencoder/skylb-apiv2/proto"
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
	if len(eps) == 0 {
		glog.Fatalln("Valid flag --skylb-endpoints is required.")
	}

	// Randomly pick one SkyLB gRPC instance to connect.
	var idx int
	if len(eps) > 1 {
		idx = rand.Intn(len(eps))
	}
	endpoint := eps[idx]
	glog.Infof("Picked SKyLB gRPC instance %s", endpoint)

	parts := stdstr.SplitN(endpoint, ":", 2)
	ep, port := parts[0], parts[1]
	addrs, err := net.LookupHost(ep)
	if err != nil {
		glog.Errorf("Failed to lookup SkyLB endpoint %s, %v.\n", endpoint, err)
		return nil, err
	}
	glog.Infof("Resolved SkyLB instances %s", addrs)

	switch len(addrs) {
	case 0:
		return nil, fmt.Errorf("No SkyLB instances found")
	case 1:
		ep = addrs[0]
	default:
		idx = rand.Intn(len(addrs))
		ep = addrs[idx]
	}
	glog.Infof("Connecting SkyLB instance %s on port %s", ep, port)

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", ep, port), grpc.WithInsecure(), grpc.WithTimeout(time.Second), grpc.WithBlock())
	if err != nil {
		glog.Errorf("Failed to dial to SkyLB instance %s, %+v.", ep, err)
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
