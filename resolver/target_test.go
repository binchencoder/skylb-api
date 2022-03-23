package resolver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDirectTarget(t *testing.T) {
	target := BuildDirectTarget("localhost:123,localhost:456")
	fmt.Println(target)
	assert.Equal(t, "direct://localhost:123,localhost:456", target)
}

func TestBuildSkyLBTarget(t *testing.T) {
	target := BuildSkyLBTarget("localhost:123,localhost:456")
	fmt.Println(target)
	assert.Equal(t, "skylb://localhost:123,localhost:456", target)
}
