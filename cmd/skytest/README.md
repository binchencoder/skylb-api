# Overrview

sklb-apiv2/cmd/skytest 是一个测试skylb api的程序

## Build the skytest

#### build gRPC server

```shell
bazel build cmd/skytest:skytest-server
```

#### build gRPC client

```skytest
bazel build cmd/skytest:skytest-client
```

## Run the skytest

#### 确定先启动skylb

```
skylb/bazel-bin/cmd/skylb/linux_amd64_stripped/skylb --etcd-endpoints="http://localhost:2377"
```

#### start gRPC server

1. 注册到skylbserver

```
skylb-apiv2/bazel-bin/cmd/skytest/skytest-server_/skytest-server -within-k8s=true -port=18000 -skylb-endpoints="127.0.0.1:1900" -alsologtostderr -v=2 -log_dir=.
```

2. 不注册到skylbserver
```
skylb-apiv2/bazel-bin/cmd/skytest/skytest-server_/skytest-server -port=18000 -alsologtostderr -v=2 -log_dir=.

skylb-apiv2/bazel-bin/cmd/skytest/skytest-server_/skytest-server -port=18001 -alsologtostderr -v=2 -log_dir=.

skylb-apiv2/bazel-bin/cmd/skytest/skytest-server_/skytest-server -port=18002 -alsologtostderr -v=2 -log_dir=.
```

#### start gRPC client

1. 直连

```shell
skylb-apiv2/bazel-bin/cmd/skytest/skytest-client_/skytest-client -debug-svc-endpoint=shared-test-server-service=localhost:18000,localhost:18001,localhost:18002 -alsologtostderr
```

2. 连skylb

```shell
skylb-apiv2/bazel-bin/cmd/skytest/skytest-client_/skytest-client -skylb-endpoints="127.0.0.1:1900" -alsologtostderr -v=2 -log_dir=.
```

   