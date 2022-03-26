package resolver

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

type directBuilder struct{}

type directResolver struct {
	cc resolver.ClientConn
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
		cc.ReportError(err)
		return nil, err
	}

	return &directResolver{cc: cc}, nil
}

func (b *directBuilder) Scheme() string {
	return DirectScheme
}

func (r *directResolver) Close() {
}

func (r *directResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
