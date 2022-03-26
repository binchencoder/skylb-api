package resolver

import (
	"errors"
	"fmt"

	"github.com/binchencoder/skylb-apiv2/client/option"
	"google.golang.org/grpc/resolver"
)

const (
	// DirectScheme stands for direct scheme.
	DirectScheme = "direct"
	// SkyLBScheme stands for skylb scheme.
	SkyLBScheme = "skylb"
)

const (
	// EndpointSepChar is the separator cha in endpoints.
	EndpointSepChar = ','

	// EqualSpeChar is the separator cha in endpoints.
	EqualSpeChar = "="
	// SlashSpeChar is the separator cha in endpoints.
	SlashSpeChar = "/"

	subsetSize = 32
)

var (
	ErrUnsupportSchema = errors.New("unsupport schema skylb")
	ErrMissServiceName = errors.New("target miss service name")
	ErrInvalidTarget   = errors.New("Invalid target url")
	ErrNoInstances     = errors.New("no valid instance")
)

var (
	// EndpointSep is the separator string in endpoints.
	EndpointSep = fmt.Sprintf("%c", EndpointSepChar)

	directResolverBuilder directBuilder
)

func init() {
	// Registers the direct scheme to the resolver.
	resolver.Register(&directResolverBuilder)
}

// Registers the skylb scheme to the resolver.
func RegisterSkylbResolverBuilder(keeper option.SkyLbKeeper) {
	resolver.Register(&skylbBuilder{
		keeper: keeper,
	})
}
