module github.com/binchencoder/skylb-apiv2

go 1.17

require (
	github.com/binchencoder/gateway-proto v0.0.5
	github.com/binchencoder/letsgo v0.0.3
)

require (
	github.com/golang/glog v1.0.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/prometheus/client_golang v1.7.1
	github.com/soheilhy/cmux v0.1.4
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.7.1
	go.uber.org/atomic v1.6.0 // indirect
	google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106 // indirect
	google.golang.org/grpc/examples v0.0.0-20200630190442-3de8449f8555 // indirect
)

replace (
	github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
	github.com/coreos/bbolt v1.3.5 => go.etcd.io/bbolt v1.3.5
)
