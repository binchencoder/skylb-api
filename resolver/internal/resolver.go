package internal

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

const (
	// DiscovScheme stands for skylb scheme.
	SkyLBScheme = "skylb"
	// EndpointSepChar is the separator cha in endpoints.
	EndpointSepChar = ','

	subsetSize = 32
)

var (
	// EndpointSep is the separator string in endpoints.
	EndpointSep = fmt.Sprintf("%c", EndpointSepChar)

	skylbResolverBuilder skylbBuilder
)

// RegisterResolver registers the skylb schemes to the resolver.
func RegisterResolver() {
	resolver.Register(&skylbResolverBuilder)
}

type skylbResolver struct {
	cc resolver.ClientConn
}

func (r *skylbResolver) Close() {
}

func (r *skylbResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
