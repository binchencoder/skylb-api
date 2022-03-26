package resolver

import (
	"fmt"
	"strings"

	pb "github.com/binchencoder/skylb-apiv2/proto"
)

// SkyLBTarget returns a string that represents the given endpoints with skylb schema.
func SkyLBTarget(spec *pb.ServiceSpec) string {
	return fmt.Sprintf("%s://%s?ns=%s&pn=%s", SkyLBScheme, spec.ServiceName, spec.Namespace, spec.PortName)
}

// DirectTarget returns a string that represents the given endpoints with direct schema.
func DirectTarget(addrs string) string {
	addrs = strings.TrimSpace(addrs)
	if addrs == "" {
		return addrs
	}

	return fmt.Sprintf("%s://%s", DirectScheme, addrs)
}
