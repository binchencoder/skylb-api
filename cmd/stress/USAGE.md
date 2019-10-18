# Overrview

sklb-api/cmd/stress 是压测gateway、skylb的一个程序

# Build the stress

build gateway server
```
bazel build cmd/stress:stress-gateway
```

build gRPC server
```
bazel build cmd/stress:stress-server
```

build gRPC client
```
bazel build cmd/stress:stress-client
```

# Run the stress

确定先启动skylb
```
skylb/bazel-bin/cmd/skylb/linux_amd64_stripped/skylb --etcd-endpoints="http://localhost:2377"
```

start gRPC server
```
skylb-api/bazel-bin/cmd/stress/linux_amd64_stripped/stress-server -skylb-endpoints="127.0.0.1:1900" -v=2 -log_dir=.
```

start //cmd/gateway
```
skylb-api/bazel-bin/cmd/stress/linux_amd64_stripped/stress-gateway -skylb-endpoints="127.0.0.1:1900" -v=2 -log_dir=.
```

start gRPC client
```
skylb-api/bazel-bin/cmd/stress/linux_amd64_stripped/stress-client -skylb-endpoints="127.0.0.1:1900" -v=2 -log_dir=.
```