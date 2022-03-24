package skylb

import (
	"strings"

	skyRs "github.com/binchencoder/skylb-apiv2/resolver"
	"google.golang.org/grpc/resolver"
)

var (
	directResolverBuilder directBuilder
)

type directBuilder struct{}

type directResolver struct {
	cc resolver.ClientConn
}

func init() {
	// Registers the direct scheme to the resolver.
	resolver.Register(&directResolverBuilder)
}

// Build
func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	var addrs []resolver.Address
	endpoints := strings.FieldsFunc(target.URL.Host, func(r rune) bool {
		return r == EndpointSepChar
	})

	for _, val := range subset(endpoints, subsetSize) {
		addrs = append(addrs, resolver.Address{
			Addr: val,
		})
	}
	if err := cc.UpdateState(resolver.State{
		Addresses: addrs,
	}); err != nil {
		return nil, err
	}

	return &directResolver{cc: cc}, nil
}

func (b *directBuilder) Scheme() string {
	return skyRs.DirectScheme
}

func (r *directResolver) Close() {
}

func (r *directResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
