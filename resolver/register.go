package resolver

import "github.com/binchencoder/skylb-apiv2/resolver/internal"

// Register registers schemes defined skylb.
// Keep it in a separated package to let third party register manually.
func Register() {
	internal.RegisterResolver()
}
