package internal

import (
	"strings"

	pb "github.com/binchencoder/skylb-apiv2/proto"
	"google.golang.org/grpc/resolver"
)

type directBuilder struct {
	endpoint *pb.InstanceEndpoint
	spec     *pb.ServiceSpec
}

// Build
func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	var addrs []resolver.Address
	endpoints := strings.FieldsFunc(target.Endpoint, func(r rune) bool {
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

	return &nopResolver{cc: cc}, nil
}

func (b *directBuilder) Scheme() string {
	return DirectScheme
}
