// Package flags defines flags for client SkyLB APIs.
package client

import (
	"flag"
	"time"
)

var (
	DebugSvcInterval = flag.Duration("debug-svc-interval", time.Minute,
		"Interval of printing services connected to client.")

	EnableHealthCheck = flag.Bool("enable-health-check", true, "True to enable health check")
)
