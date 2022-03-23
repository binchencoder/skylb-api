package tests

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"google.golang.org/grpc/resolver"
)

func TestParseTarget(t *testing.T) {
	target := "direct://vxserver:4100,vxserver1:4100"
	u, err := url.Parse(target)
	if err != nil {
		fmt.Println("parset err: ", err)
		return
	}
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
	fmt.Printf("resolver.Target: %+v \n", resolver.Target{
		Scheme:    u.Scheme,
		Authority: u.Host,
		Endpoint:  endpoint,
		URL:       *u,
	})
}
