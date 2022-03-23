package balancer

import (
	"time"
)

// DebugBalancer defines the grpc load balancer which prints debug info.
type DebugBalancer interface {
	// StartDebugPrint starts printing debug info.
	StartDebugPrint(interval time.Duration)
	// StopDebugPrint stops printing debug info.
	StopDebugPrint()
}
