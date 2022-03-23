package resolver

import (
	"fmt"
	"strings"

	"github.com/binchencoder/skylb-apiv2/resolver/internal"
)

// BuildSkyLBTarget returns a string that represents the given endpoints with skylb schema.
func BuildSkyLBTarget(addrs string) string {
	if len(addrs) == 0 {
		return addrs
	}

	addrs = strings.TrimSpace(addrs)
	return fmt.Sprintf("%s://%s", internal.SkyLBScheme, addrs)
}

// BuildDirectTarget returns a string that represents the given endpoints with direct schema.
func BuildDirectTarget(addrs string) string {
	if len(addrs) == 0 {
		return addrs
	}

	addrs = strings.TrimSpace(addrs)
	return fmt.Sprintf("%s://%s", internal.DirectScheme, addrs)
}
