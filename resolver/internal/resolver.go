package internal

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/resolver"
)

var (
	ErrUnsupportSchema = errors.New("unsupport schema skylb")
	ErrMissServiceName = errors.New("target miss service name")
	ErrNoInstances     = errors.New("no valid instance")
)

const (
	// DirectScheme stands for direct scheme.
	DirectScheme = "direct"
	// SkyLBScheme stands for skylb scheme.
	SkyLBScheme = "skylb"
	// EndpointSepChar is the separator cha in endpoints.
	EndpointSepChar = ','

	// EqualSpeChar is the separator cha in endpoints.
	EqualSpeChar = "="
	// SlashSpeChar is the separator cha in endpoints.
	SlashSpeChar = "/"

	subsetSize = 32
)

var (
	// EndpointSep is the separator string in endpoints.
	EndpointSep = fmt.Sprintf("%c", EndpointSepChar)

	skylbResolverBuilder  skylbBuilder
	directResolverBuilder directBuilder
)

// RegisterResolver registers the skylb schemes to the resolver.
func RegisterResolver() {
	resolver.Register(&directResolverBuilder)
	resolver.Register(&skylbResolverBuilder)
}

type skylbResolver struct {
	cc resolver.ClientConn
}

func (r *skylbResolver) Close() {
}

func (r *skylbResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

type nopResolver struct {
	cc resolver.ClientConn
}

func (r *nopResolver) Close() {
}

func (r *nopResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
