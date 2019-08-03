package health

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
)

func TestHealthCheckContext(t *testing.T) {
	Convey("Create a health check context", t, func() {
		ctx := withHealthCheck(context.Background())
		So(IsHealthCheck(ctx), ShouldBeTrue)
		So(IsHealthCheck(context.Background()), ShouldBeFalse)
	})
}
