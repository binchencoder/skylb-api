package internal

import "google.golang.org/grpc/resolver"

type directBuilder struct{}

// Build
func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {

	return nil, nil
}
