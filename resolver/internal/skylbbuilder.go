package internal

import (
	"fmt"

	skyRs "github.com/binchencoder/skylb-apiv2/resolver"
	"google.golang.org/grpc/resolver"
)

type skylbBuilder struct {
}

type skylbResolver struct {
	cc resolver.ClientConn
}

func (b *skylbBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	serverName := target.URL.Host

	fmt.Printf("serverName: %s\n", serverName)

	resolver := &skylbResolver{cc: cc}

	return resolver, nil
}

func (b *skylbBuilder) Scheme() string {
	return skyRs.SkyLBScheme
}

func (r *skylbResolver) Close() {
}

func (r *skylbResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
