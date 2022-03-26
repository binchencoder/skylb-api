package resolver

import (
	"net/url"

	"github.com/binchencoder/skylb-apiv2/client/option"
	pb "github.com/binchencoder/skylb-apiv2/proto"
	"google.golang.org/grpc/resolver"
)

type skylbBuilder struct {
	keeper option.SkyLbKeeper
}

type skylbResolver struct {
	cliConn resolver.ClientConn
}

// Build
func (b *skylbBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	if target.URL.Scheme != SkyLBScheme {
		return nil, ErrUnsupportSchema
	}

	_, err := url.Parse(target.Endpoint)
	if err != nil {
		panic(ErrInvalidTarget)
	}

	url := target.URL
	values := url.Query()
	servSpec := &pb.ServiceSpec{
		ServiceName: url.Host,
		Namespace:   values.Get("ns"),
		PortName:    values.Get("pn"),
	}

	b.keeper.RegisterServiceCliConn(servSpec, cc)

	resolver := &skylbResolver{cliConn: cc}
	return resolver, nil
}

func (b *skylbBuilder) Scheme() string {
	return SkyLBScheme
}

func (r *skylbResolver) Close() {
}

func (r *skylbResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
