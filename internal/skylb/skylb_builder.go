package skylb

import (
	"fmt"
	"net/url"

	"github.com/binchencoder/skylb-apiv2/client/option"
	pb "github.com/binchencoder/skylb-apiv2/proto"
	skyRs "github.com/binchencoder/skylb-apiv2/resolver"
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
	if target.Scheme != skyRs.SkyLBScheme {
		return nil, ErrUnsupportSchema
	}

	u, err := url.Parse(target.Endpoint)
	if err != nil {
		panic(ErrInvalidTarget)
	}

	values := u.Query()
	servSpec := &pb.ServiceSpec{
		ServiceName: u.Host,
		Namespace:   values.Get("ns"),
		PortName:    values.Get("pn"),
	}

	fmt.Printf("servSpec: %+v\n", servSpec)
	b.keeper.RegisterServiceCliConn(servSpec, cc)

	resolver := &skylbResolver{cliConn: cc}
	return resolver, nil
}

func (b *skylbBuilder) Scheme() string {
	return skyRs.SkyLBScheme
}

func (r *skylbResolver) Close() {
}

func (r *skylbResolver) ResolveNow(options resolver.ResolveNowOptions) {
}