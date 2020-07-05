module github.com/binchencoder/skylb-api

go 1.13

require (
	github.com/VividCortex/gohistogram v1.0.0
	github.com/beorn7/perks v1.0.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/cespare/xxhash/v2 v2.1.1
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/color v1.9.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-kit/kit v0.10.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/json-iterator/go v1.1.10
	github.com/jtolds/gls v4.20.0+incompatible
	github.com/klauspost/compress v1.10.10
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v1.0.1
	github.com/mwitkow/grpc-proxy v0.0.0-20181017164139-0f1106ef9c76
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pborman/uuid v1.2.0
	github.com/prometheus/client_golang v1.3.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.10.0
	github.com/prometheus/procfs v0.1.3
	github.com/smartystreets/assertions v1.0.1
	github.com/smartystreets/goconvey v1.6.4
	github.com/soheilhy/cmux v0.1.4
	github.com/stretchr/objx v0.2.0
	github.com/stretchr/testify v1.6.1
	github.com/uber/jaeger-client-go v2.24.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible
	go.uber.org/atomic v1.6.0
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v2 v2.2.4
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	upper.io/db.v3 v3.7.1+incompatible
)

replace (
	github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
	github.com/coreos/bbolt v1.3.5 => go.etcd.io/bbolt v1.3.5
)
