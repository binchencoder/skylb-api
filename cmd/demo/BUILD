load("@io_bazel_rules_docker//docker:docker.bzl", "docker_build")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_test")

go_binary(
    name = "client",
    srcs = ["client.go"],
    deps = [
        "//letsgo:go_default_library",
        "@binchencoder_letsgo//metrics:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
    ],
)

go_binary(
    name = "client_failfast",
    srcs = ["client_failfast.go"],
    deps = [
        "//letsgo:go_default_library",
        "//client:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

go_binary(
    name = "consistenthash",
    srcs = ["consistenthash.go"],
    deps = [
        "//letsgo:go_default_library",
        "@binchencoder_letsgo//hashring:go_default_library",
        "//balancer:go_default_library",
        "//client:go_default_library",
        "//client/option:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/naming:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

go_binary(
    name = "webclient",
    srcs = ["webclient.go"],
    deps = [
        "//letsgo:go_default_library",
        "//client:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

go_binary(
    name = "streamclient",
    srcs = ["streamclient.go"],
    deps = [
        "//letsgo:go_default_library",
        "//client:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

go_binary(
    name = "server",
    srcs = ["server.go"],
    deps = [
        "//letsgo:go_default_library",
        "@binchencoder_letsgo//grpc:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//server:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/codes:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@binchencoder_ease_gateway//proto/frontend:go_default_library",
    ],
)

go_binary(
    name = "servermulti",
    srcs = ["servermulti.go"],
    deps = [
        "//letsgo:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//server:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

go_binary(
    name = "serverclient",
    srcs = ["serverclient.go"],
    deps = [
        "//letsgo:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//server:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

go_binary(
    name = "gateway",
    srcs = ["gateway.go"],
    deps = [
        "//janus/gateway/runtime:go_default_library",
        "//letsgo:go_default_library",
        "//cmd/demo/proto:go_default_library",
    ],
)

go_test(
    name = "int_test",
    size = "medium",
    srcs = glob(["*int_test.go"]),
    deps = [
        "//client:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//proto:go_default_library",
        "//server:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

# Build and test this image locally. No need to push to repository.
docker_build(
    name = "latest",
    base = "//bld_tools/docker/ubuntu:base",
    directory = "/skylb",
    files = [
        ":client",
        ":consistenthash",
        ":server",
    ],
    repository = "harbor.jingoal.ltd",
    tag = "skylb/demo-go",
    workdir = "/skylb",
)
