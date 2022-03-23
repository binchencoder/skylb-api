package internal

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"
)

type skylbBuilder struct {
}

func (b *skylbBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	hosts := strings.FieldsFunc(target.Authority, func(r rune) bool {
		return r == EndpointSepChar
	})

	fmt.Println(hosts)

	return &skylbResolver{cc: cc}, nil
}

func (b *skylbBuilder) Scheme() string {
	return SkyLBScheme
}
