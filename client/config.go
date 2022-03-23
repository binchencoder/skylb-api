package client

import (
	"flag"

	vexpb "github.com/binchencoder/gateway-proto/data"
	"github.com/binchencoder/letsgo/flags"
	"github.com/binchencoder/letsgo/service/naming"
	pb "github.com/binchencoder/skylb-apiv2/proto"
)

const (
	DefaultNameSpace = "default"
	DefaultPortName  = "grpc"
)

var (
	DebugSvcEndpoints = flags.StringMap{}
)

type (
	// A RpcClientConf is a rpc client config.
	RpcClientConf struct {
		// Endpoints []string `json:",optional"`
		// Target    string   `json:",optional"`
		// App       string   `json:",optional"`
		// Token     string   `json:",optional"`
		NonBlock bool  `json:",optional"`
		Timeout  int64 `json:",default=2000"`
	}
)

func init() {
	flag.Var(&DebugSvcEndpoints, "debug-svc-endpoint", "The debug service endpoint. If not empty, disable SkyLB resolving for that service.")
}

// BuildTarget builds the rpc target from the given config.
func (cc RpcClientConf) BuildTarget() (string, error) {
	// TODO(chenbin)
	// if len(cc.Endpoints) > 0 {
	// 	return resolver.BuildDirectTarget(cc.Endpoints), nil
	// } else if len(cc.Target) > 0 {
	// 	return cc.Target, nil
	// }

	// return resolver.BuildDiscovTarget(cc.Etcd.Hosts, cc.Etcd.Key), nil
	return "", nil
}

// HasCredential checks if there is a credential in config.
// func (cc RpcClientConf) HasCredential() bool {
// 	return len(cc.App) > 0 && len(cc.Token) > 0
// }

// NewServiceSpec returns a new ServiceSpec struct with the given parameters.
func NewServiceSpec(namespace string, serviceId vexpb.ServiceId, portName string) *pb.ServiceSpec {
	if namespace == "" {
		namespace = DefaultNameSpace
	}
	if portName == "" {
		portName = DefaultPortName
	}
	serviceName, err := naming.ServiceIdToName(serviceId)
	if err != nil {
		panic("Unknown service ID.")
	}

	return &pb.ServiceSpec{
		Namespace:   namespace,
		ServiceName: serviceName,
		PortName:    portName,
	}
}

// NewDefaultServiceSpec returns a new ServiceSpec struct with the default
// namespace and default port name.
func NewDefaultServiceSpec(serviceId vexpb.ServiceId) *pb.ServiceSpec {
	return NewServiceSpec(DefaultNameSpace, serviceId, DefaultPortName)
}
