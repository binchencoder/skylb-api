package skylb

import (
	"errors"
	"fmt"
)

var (
	ErrUnsupportSchema = errors.New("unsupport schema skylb")
	ErrMissServiceName = errors.New("target miss service name")
	ErrInvalidTarget   = errors.New("Invalid target url")
	ErrNoInstances     = errors.New("no valid instance")
)

var (
	// EndpointSep is the separator string in endpoints.
	EndpointSep = fmt.Sprintf("%c", EndpointSepChar)
)

const (
	// EndpointSepChar is the separator cha in endpoints.
	EndpointSepChar = ','

	// EqualSpeChar is the separator cha in endpoints.
	EqualSpeChar = "="
	// SlashSpeChar is the separator cha in endpoints.
	SlashSpeChar = "/"

	subsetSize = 32
)
