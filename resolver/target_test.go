package resolver

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	pb "github.com/binchencoder/skylb-api/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

func TestDirectTarget(t *testing.T) {
	target := DirectTarget("localhost:123,localhost:456")
	fmt.Println(target)
	assert.Equal(t, "direct://localhost:123,localhost:456", target)
}

func TestSkyLBTarget(t *testing.T) {
	target := SkyLBTarget(&pb.ServiceSpec{
		Namespace:   "namespace",
		ServiceName: "serviceName",
		PortName:    "portName",
	})
	fmt.Println(target)
	assert.Equal(t, "skylb://serviceName?ns=namespace&pn=portName", target)
}

func TestParseTarget(t *testing.T) {
	target := SkyLBTarget(&pb.ServiceSpec{
		Namespace:   "namespace",
		ServiceName: "serviceName",
		PortName:    "portName",
	})

	u, _ := url.Parse(target)
	// For targets of the form "[scheme]://[authority]/endpoint, the endpoint
	// value returned from url.Parse() contains a leading "/". Although this is
	// in accordance with RFC 3986, we do not want to break existing resolver
	// implementations which expect the endpoint without the leading "/". So, we
	// end up stripping the leading "/" here. But this will result in an
	// incorrect parsing for something like "unix:///path/to/socket". Since we
	// own the "unix" resolver, we can workaround in the unix resolver by using
	// the `URL` field instead of the `Endpoint` field.
	endpoint := u.Path
	if endpoint == "" {
		endpoint = u.Opaque
	}
	endpoint = strings.TrimPrefix(endpoint, "/")
	resTarget := resolver.Target{
		Scheme:    u.Scheme,
		Authority: u.Host,
		Endpoint:  endpoint,
		URL:       *u,
	}
	fmt.Printf("target: %s, parsedTarget: %+v \n", target, resTarget)
}
