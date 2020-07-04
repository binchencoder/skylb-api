package lameduck

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strings"
	"sync"

	etcd "github.com/coreos/etcd/client"
	"github.com/golang/glog"

	"github.com/binchencoder/skylb-api/prefix"
	"github.com/binchencoder/skylb-api/util"
)

var (
	// KeyPattern is the regexp pattern of the ETCD key for lameduck endpoints.
	KeyPattern = regexp.MustCompile(`^/grpc/lameduck/services/([^/]*)/endpoints/([^/]*)$`)

	lameducks = make([]string, 0, 100)
	lock      sync.RWMutex

	getOpts = etcd.GetOptions{
		Recursive: true,
	}
	setOpts = etcd.SetOptions{}
	delOpts = etcd.DeleteOptions{}
)

// HostPort combines the given host and port into a lameduck endpoint.
func HostPort(host, port string) string {
	return fmt.Sprintf("%s#%s", host, port)
}

// IsLameduckMode returns true if the given endpoint is in lameduck mode.
func IsLameduckMode(endpoint string) bool {
	lock.RLock()
	defer lock.RUnlock()

	for _, v := range lameducks {
		if v == endpoint {
			return true
		}
	}

	return false
}

// SetLameDuckMode sets the endpoint to lameduck mode.
func SetLameDuckMode(etcdCli etcd.KeysAPI, svcName, ep string) error {
	key := fmt.Sprintf("%s%s/endpoints/%s", prefix.LameduckKey, svcName, ep)
	glog.V(3).Infof("SetLameDuckMode %s", key)
	if _, err := etcdCli.Set(context.Background(), key, "", &setOpts); nil != err {
		return err
	}
	return nil
}

// UnsetLameDuckMode takes the endpoint out of lameduck mode.
func UnsetLameDuckMode(etcdCli etcd.KeysAPI, svcName, ep string) error {
	key := fmt.Sprintf("%s%s/endpoints/%s", prefix.LameduckKey, svcName, ep)
	glog.V(3).Infof("SetLameDuckMode %s", key)
	if _, err := etcdCli.Delete(context.Background(), key, &delOpts); nil != err {
		return err
	}
	return nil
}

// ExtractLameduck recursively extracts the lameduck endpoints from the given
// ETCD nodes.
func ExtractLameduck(node *etcd.Node) {
	if KeyPattern.MatchString(node.Key) {
		ep := path.Base(node.Key)
		t := addLameduckEndpoint(ep)
		glog.V(2).Infof("Added lameduck endpoint %s, changed: %t", ep, t)
	} else {
		for _, n := range node.Nodes {
			ExtractLameduck(n)
		}
	}
}

// ExtractLameduckChange extracts the lameduck endpoints and operation
// from watch response of ETCD.
func ExtractLameduckChange(resp *etcd.Response) {
	if resp == nil || resp.Node == nil {
		return
	}

	key := resp.Node.Key
	ep := path.Base(key)

	switch resp.Action {
	case util.ActionCreate, util.ActionSet:
		t := addLameduckEndpoint(ep)
		glog.V(2).Infof("Added lameduck endpoint %s, changed: %t", ep, t)
	case util.ActionDelete, util.ActionExpire:
		t := removeLameduckEndpoint(ep)
		glog.V(2).Infof("Removed lameduck endpoint %s, changed: %t", ep, t)
	default:
		glog.Errorf("Unexpected action %s in ExtractLameduckChange(), ignore.", resp.Action)
	}
}

// LoadLameducks returns a map from lameduck endpoints to service names.
func LoadLameducks(etcdCli etcd.KeysAPI, serviceName string) map[string]string {
	prefix := prefix.LameduckKey + serviceName + "/"
	resp, err := etcdCli.Get(context.Background(), prefix, &getOpts)
	if err != nil {
		glog.Errorf("Failed to load lameduck instances with key prefix %s, %v", prefix, err)
		return nil
	}
	return extractLameducksAsMap(resp.Node)
}

func extractLameducksAsMap(root *etcd.Node) map[string]string {
	m := map[string]string{}
	for _, node := range root.Nodes {
		keys := make([]string, 0, 10)
		getLeafKeys(node, &keys)
		for _, key := range keys {
			matched := KeyPattern.FindStringSubmatch(key)
			if len(matched) != 3 {
				glog.Errorf("Found invalid lameduck key: %s.", node.Key)
				continue
			}
			eps := strings.Replace(matched[2], "#", ":", 1)
			m[eps] = matched[1]
		}
	}
	return m
}

func getLeafKeys(root *etcd.Node, keys *[]string) {
	if !root.Dir {
		*keys = append(*keys, root.Key)
		return
	}
	for _, node := range root.Nodes {
		getLeafKeys(node, keys)
	}
}

// addLameduckEndpoint adds the ep in lameducks pool.
func addLameduckEndpoint(ep string) bool {
	if IsLameduckMode(ep) {
		return false
	}

	lock.Lock()
	defer lock.Unlock()

	lameducks = append(lameducks, ep)
	return true
}

// removeLameduckEndpoint removes the endpoint out of lameducks pool.
// It returns true if the endpoint was found in the lameducks pool.
func removeLameduckEndpoint(endpoint string) bool {
	lock.Lock()
	defer lock.Unlock()

	found := false
	for i, v := range lameducks {
		if v == endpoint {
			lameducks = append(lameducks[0:i], lameducks[i+1:]...)
			found = true
			break
		}
	}

	return found
}
