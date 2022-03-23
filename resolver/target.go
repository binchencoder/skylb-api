package resolver

import (
	"fmt"
	"strings"

	"github.com/binchencoder/skylb-apiv2/resolver/internal"
)

// BuildSkyLBTarget returns a string that represents the given endpoints with skylb schema.
func BuildSkyLBTarget(endpoints []string, key string) string {
	return fmt.Sprintf("%s://%s/%s", internal.SkyLBScheme,
		strings.Join(endpoints, internal.EndpointSep), key)
}

func BuildDebugTarget() string {
	return ""
}
