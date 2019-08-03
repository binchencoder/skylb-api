package balancer

import "google.golang.org/grpc"

type addrInfo struct {
	addr      grpc.Address
	connected bool
	weight    int32
}
