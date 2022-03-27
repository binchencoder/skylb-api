package server

import (
	"errors"
	"flag"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/soheilhy/cmux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	hpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/peer"

	vexpb "github.com/binchencoder/gateway-proto/data"
	jg "github.com/binchencoder/letsgo/grpc"
	lmetrics "github.com/binchencoder/letsgo/metrics"
	"github.com/binchencoder/letsgo/runtime/pprof"
	"github.com/binchencoder/skylb-api/internal/flags"
	"github.com/binchencoder/skylb-api/internal/rpccli"
	"github.com/binchencoder/skylb-api/metrics"
	"github.com/binchencoder/skylb-api/naming"
	pb "github.com/binchencoder/skylb-api/proto"
)

const (
	defaultNameSpace = "default"
	defaultWeight    = 100
)

var (
	// TODO(chenbin): change default value to true when most apps are k8s-ready.
	withinK8s             = flag.Bool("within-k8s", false, "Whether this grpc service runs in kubernetes cluster")
	reportInterval        = flag.Duration("skylb-report-interval", 3*time.Second, "The SkyLB load-report interval.")
	EnablePromOnSamePort  = flag.Bool("enable-prom-on-same-port", false, "Enable serving prometheus metrics on the same port as gRPC.")
	enablePprofOnSamePort = flag.Bool("enable-pprof-on-same-port", false, "Enable serving pprof on the same port as gRPC.")
	withWeight            = flag.Int("with-weight", 0, "The serving server with weight for load balancing decisions. For compatibility with services that do not need to provide weights, the default value is 0.")

	skylbPanicRecoveryCounts = prom.NewCounter(
		prom.CounterOpts{
			Namespace: "skylb",
			Subsystem: "server",
			Name:      "panic_recovery_counts",
			Help:      "Skylb server panic recovery counts.",
		},
	)

	once        sync.Once
	spec        *pb.ServiceSpec
	servicePort int

	withPrometheusHistogram = true

	// For multi-service use. Independent from single service version.
	serviceIds   []vexpb.ServiceId
	servicePorts []int
	specs        []*pb.ServiceSpec
	hostAddrs    []string

	// Optional customized grpc server options.
	// Set with SetGrpcServerOptions() before call to Start() or StartMulti().
	grpcServerOptions []grpc.ServerOption
)

func init() {
	prom.MustRegister(skylbPanicRecoveryCounts)
}

// A simple grpc health service implementation, similar to "echo".
type healthServer struct {
}

func (hs *healthServer) Check(ctx context.Context, in *hpb.HealthCheckRequest) (*hpb.HealthCheckResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return &hpb.HealthCheckResponse{
			Status: hpb.HealthCheckResponse_UNKNOWN,
		}, errors.New("Failed to get peer client info from context.")
	}
	glog.V(2).Infof("Health check from client %s", p.Addr.String())
	return &hpb.HealthCheckResponse{
		Status: hpb.HealthCheckResponse_SERVING,
	}, nil
}

func (hs *healthServer) Watch(in *hpb.HealthCheckRequest, stream hpb.Health_WatchServer) error {
	return errors.New("not implemented")
}

// Register starts the SkyLB load-report worker with the given service spec
// and port.
func Register(serviceId vexpb.ServiceId, portName string, port int) {
	servicePort = port

	serviceName, err := naming.ServiceIdToName(serviceId)
	if err != nil {
		panic(err)
	}

	spec = &pb.ServiceSpec{
		Namespace:   defaultNameSpace,
		ServiceName: serviceName,
		PortName:    portName,
	}
}

// RegisterMulti is similar to Register but supports multiple services.
func RegisterMulti(serviceId vexpb.ServiceId, portName string, port int, hostAddr string) {
	Register(serviceId, portName, port)
	serviceIds = append(serviceIds, serviceId)
	servicePorts = append(servicePorts, port)
	specs = append(specs, spec)
	hostAddrs = append(hostAddrs, hostAddr)
}

func EnableHistogram() {
	withPrometheusHistogram = true
}

// Set customized GRPC server options, such as max send/recv message size.
// This function must be called ahead of Start() or StartMulti().
// grpc.NewServer(...) panics if opts contains UnaryInterceptor or StreamInterceptor.
func SetGrpcServerOptions(opts ...grpc.ServerOption) {
	grpcServerOptions = opts
}

func Start(hostAddr string, callback func(s *grpc.Server) error, interceptors ...grpc.UnaryServerInterceptor) {
	lis, s := start0(hostAddr, servicePort, spec, interceptors...)

	if err := callback(s); err != nil {
		panic(err)
	}

	serve(lis, s)
}

func StartMulti(callback func(serviceId vexpb.ServiceId, s *grpc.Server) error, interceptors ...grpc.UnaryServerInterceptor) {
	for i, sid := range serviceIds {
		hostAddr := hostAddrs[i]
		port := servicePorts[i]
		spec := specs[i]
		lis, s := start0(hostAddr, port, spec, interceptors...)

		if err := callback(sid, s); err != nil {
			panic(err)
		}

		if i < len(serviceIds)-1 {
			// Not the last, use goroutine to start services.
			go serve(lis, s)
			continue
		}

		// Now this is the last one
		serve(lis, s)
	}
}

func serve(lis net.Listener, grpcs *grpc.Server) {
	if withPrometheusHistogram {
		metrics.EnableHandlingTimeHistogram()
	}

	glog.Infoln("Registering health check service.")
	hpb.RegisterHealthServer(grpcs, &healthServer{})

	if !*EnablePromOnSamePort && !*enablePprofOnSamePort {
		glog.Infoln("Start serving grpc.")
		startGrcpServer(lis, grpcs)
		return
	}

	m := cmux.New(lis)

	// Match connections in order: first grpc, then HTTP.
	grpcl := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpl := m.Match(cmux.HTTP1Fast())

	go startHTTPServer(httpl, http.DefaultServeMux)
	go startGrcpServer(grpcl, grpcs)

	m.Serve()
}

func startGrcpServer(grpcl net.Listener, grpcs *grpc.Server) {
	if err := grpcs.Serve(grpcl); err != nil {
		panic(err)
	}
}

func startHTTPServer(lis net.Listener, mux *http.ServeMux) {
	if *EnablePromOnSamePort {
		lmetrics.EnablePrometheus(mux)
	}
	if *enablePprofOnSamePort {
		pprof.EnablePprof(mux)
	}
	once.Do(func() {
		if err := http.Serve(lis, mux); err != nil {
			panic(err)
		}
	})
}

func start0(hostAddr string, servicePort int, spec *pb.ServiceSpec, interceptors ...grpc.UnaryServerInterceptor) (net.Listener, *grpc.Server) {
	if !*withinK8s {
		glog.Infoln("Outside k8s, starting load reporter.")
		go StartSkylbReportLoad(spec, servicePort)
	} else {
		glog.Warningln("WARNING: Inside k8s, load-reporting disabled.")
	}

	lis, err := net.Listen("tcp", hostAddr)
	if err != nil {
		panic(err)
	}

	tracer := jg.InitOpenTracing(spec.ServiceName)
	openTracingInterceptor := otgrpc.OpenTracingServerInterceptor(tracer)

	metricsInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return metrics.UnaryServerInterceptor(spec, ctx, req, info, handler)
	}

	incepts := make([]grpc.UnaryServerInterceptor, 0, 3+len(interceptors))
	// ServerFromMetadataInterceptor needs to be the first.
	incepts = append(incepts, jg.ServerFromMetadataInterceptor)
	incepts = append(incepts, metricsInterceptor, openTracingInterceptor)
	incepts = append(incepts, interceptors...)

	glog.Infof("Creating grpc server[%s] at [%s]", spec.ServiceName, hostAddr)

	opts := make([]grpc.ServerOption, 0, 10)
	opts = append(opts, grpcServerOptions...)

	midwareOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			skylbPanicRecoveryCounts.Inc()
			glog.Errorf("[grpc-middleware-recovery] recovered from panic: %+v. Trace: %s", p, getStackTrace())
			return grpc.Errorf(codes.Internal, "panic triggered: %v", p)
		}),
	}

	// Recovery handlers should typically be last in the chain so that other
	// middleware (e.g. logging) can operate on the recovered state instead of
	// being directly affected by any panic.
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		jg.ChainUnaryServer(incepts...),
		grpc_recovery.UnaryServerInterceptor(midwareOpts...),
	))
	opts = append(opts, grpc_middleware.WithStreamServerChain(
		func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return metrics.StreamServerInterceptor(spec, srv, ss, info, handler)
		},
		grpc_recovery.StreamServerInterceptor(midwareOpts...),
	))

	return lis, grpc.NewServer(opts...)
}

// StartSkylbReportLoad starts load reporting to SkyLB server
// at the fixed interval.
func StartSkylbReportLoad(spec *pb.ServiceSpec, port int) {
	for {
		ctx, cancel := context.WithCancel(context.Background())
		cli, err := rpccli.NewGrpcClient(ctx)
		if err != nil {
			glog.Errorf("Failed to create gRPC client to skyLB, %v", err)
			cancel()

			time.Sleep(*flags.SkylbRetryInterval)
			continue
		}

		if err = reportLoad(ctx, cli, spec, "", port); err != nil {
			glog.Errorf("Failed to report load to skyLB, %v", err)
			cancel()

			time.Sleep(*flags.SkylbRetryInterval)
			continue
		}
	}
}

// StartSkylbReportLoadWithFixedHost starts load reporting to SkyLB server
// at the fixed interval with the fixed host.
//
// This is used by use cases like an agent program registering service for
// a 3rd party program, like a SQL database, for which we are not able
// to let them register themselves.
//
// If your service relies on SkyLB's auto service registry/discovery, please
// call StartSkylbReportLoad() instead.
func StartSkylbReportLoadWithFixedHost(spec *pb.ServiceSpec, fixedHost string, port int) {
	for {
		ctx, cancel := context.WithCancel(context.Background())
		cli, err := rpccli.NewGrpcClient(ctx)
		if err != nil {
			glog.Errorf("Failed to create gRPC client to skyLB, %v", err)
			cancel()

			time.Sleep(*flags.SkylbRetryInterval)
			continue
		}

		if err = reportLoad(ctx, cli, spec, fixedHost, port); err != nil {
			glog.Errorf("Failed to report load to skyLB, %v", err)
			cancel()

			time.Sleep(*flags.SkylbRetryInterval)
			continue
		}
	}
}

func reportLoad(ctx context.Context, cli pb.SkylbClient, spec *pb.ServiceSpec, fixedHost string, port int) error {
	stream, err := cli.ReportLoad(ctx)
	if err != nil {
		glog.Errorf("Failed to call ReportLoad, %v", err)
		return err
	}
	defer stream.CloseSend()

	req := pb.ReportLoadRequest{
		Spec:      spec,
		FixedHost: fixedHost,
		Port:      int32(port),
		Weight:    int32(*withWeight),
	}
	for {
		if err := stream.Send(&req); err != nil {
			glog.Errorf("Failed to send load report, %v", err)
			return err
		}
		time.Sleep(*reportInterval)
	}
}

func getStackTrace() string {
	stack := strings.Split(string(debug.Stack()), "\n")
	stacks := make([]string, 0, len(stack))
	stacks = append(stack[0:1], stack[7:]...)
	return strings.Join(stacks, "\n")
}
