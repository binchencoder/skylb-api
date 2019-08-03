package prefix

import (
	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

const (
	// EndpointsKey is the prefix for service endpoints key.
	EndpointsKey = "/registry/services/endpoints"

	// GraphKey is the prefix for service graph key.
	GraphKey = "/skylb/graph"

	// LameduckKey is the prefix of the ETCD key for lameduck endpoints.
	LameduckKey = "/grpc/lameduck/services/"
)

var (
	getOpts = etcd.GetOptions{}
)

// Init initializes ETCD keys.
func Init(etcdCli etcd.KeysAPI) {
	mustExist(etcdCli, EndpointsKey)
	mustExist(etcdCli, GraphKey)
	mustExist(etcdCli, LameduckKey)
}

func mustExist(etcdCli etcd.KeysAPI, key string) {
	if _, err := etcdCli.Get(context.Background(), key, &getOpts); err != nil {
		// For whatever reason it failed, let's try to create the key.
		if _, err = etcdCli.Set(context.Background(), key, "", &etcd.SetOptions{
			Dir:       true,
			PrevExist: etcd.PrevNoExist,
		}); err != nil {
			panic(err)
		}
	}
}
