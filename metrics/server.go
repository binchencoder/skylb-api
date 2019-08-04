// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// gRPC Prometheus monitoring interceptors for server-side gRPC.

package metrics

import (
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	jgrpc "github.com/binchencoder/letsgo/grpc"
	pb "github.com/binchencoder/skylb-api/proto"
)

// PreregisterServices takes a gRPC server and pre-initializes all counters to 0.
// This allows for easier monitoring in Prometheus (no missing metrics), and should be called *after* all services have
// been registered with the server.
func Register(server *grpc.Server) {
	serviceInfo := server.GetServiceInfo()
	for serviceName, info := range serviceInfo {
		for _, mInfo := range info.Methods {
			preRegisterMethod(serviceName, &mInfo)
		}
	}
}

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
func UnaryServerInterceptor(spec *pb.ServiceSpec, ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	_, cln := jgrpc.FromMetadataIncoming(ctx) // cln is guaranteed not nil.
	glog.V(3).Infof("client: %s\n", cln)
	monitor := newServerReporter(spec, Unary, info.FullMethod, cln)
	monitor.ReceivedMessage()
	resp, err := handler(ctx, req)
	monitor.Handled(grpc.Code(err))
	if err == nil {
		monitor.SentMessage()
	}
	return resp, err
}

// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
func StreamServerInterceptor(spec *pb.ServiceSpec, srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// TODO(fuyc): handle stream.
	monitor := newServerReporter(spec, streamRpcType(info), info.FullMethod, "TODO")
	err := handler(srv, &monitoredServerStream{ss, monitor})
	monitor.Handled(grpc.Code(err))
	return err
}

func streamRpcType(info *grpc.StreamServerInfo) grpcType {
	if info.IsClientStream && !info.IsServerStream {
		return ClientStream
	} else if !info.IsClientStream && info.IsServerStream {
		return ServerStream
	}
	return BidiStream
}

// monitoredStream wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type monitoredServerStream struct {
	grpc.ServerStream
	monitor *serverReporter
}

func (s *monitoredServerStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.monitor.SentMessage()
	}
	return err
}

func (s *monitoredServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.monitor.ReceivedMessage()
	}
	return err
}
