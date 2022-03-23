package health

import "golang.org/x/net/context"

const (
	reqTypeHealthCheck = 1
)

type requestTypeKey struct{}

// withHealthCheck returns a copy of parent context in which
// the health check is set.
func withHealthCheck(parent context.Context) context.Context {
	return context.WithValue(parent, requestTypeKey{}, reqTypeHealthCheck)
}

// IsHealthCheck returns true if the request is a gRPC health check.
func IsHealthCheck(ctx context.Context) bool {
	if reqType, ok := ctx.Value(requestTypeKey{}).(int); ok {
		return reqType == reqTypeHealthCheck
	}
	return false
}
