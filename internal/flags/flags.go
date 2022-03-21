// Package flags defines flags common for both client and server SkyLB APIs.
package flags

import (
	"flag"
	"time"
)

var (
	SkylbEndpoints     = flag.String("skylb-endpoints", "", "The SkyLB endpoints")
	SkylbRetryInterval = flag.Duration("skylb-retry-interval", time.Second, "The retry interval to connect to SkyLB")
)
