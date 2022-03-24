package resolver

import (
	"fmt"
	"testing"

	pb "github.com/binchencoder/skylb-apiv2/proto"
	"github.com/stretchr/testify/assert"
)

func TestDirectTarget(t *testing.T) {
	target := DirectTarget("localhost:123,localhost:456")
	fmt.Println(target)
	assert.Equal(t, "direct://localhost:123,localhost:456", target)
}

func TestSkyLBTarget(t *testing.T) {
	target := SkyLBTarget(&pb.ServiceSpec{
		Namespace:   "namespace",
		ServiceName: "serviceName",
		PortName:    "portName",
	})
	fmt.Println(target)
	assert.Equal(t, "skylb://serviceName?ns=namespace&pn=portName", target)
}
