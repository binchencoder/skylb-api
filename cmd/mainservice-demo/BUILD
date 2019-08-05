load("@io_bazel_rules_docker//docker:docker.bzl", "docker_build")
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "mainservice-client",
    srcs = ["client.go"],
    deps = [
        "//letsgo:go_default_library",
        "//balancer:go_default_library",
        "//client:go_default_library",
        "//client/option:go_default_library",
        "//cmd/mainservice-demo/proto:go_default_library",
        "//handlers:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/naming:go_default_library",
        "//ease-gateway/proto/data:go_default_library",
    ],
)

go_binary(
    name = "mainservice-server",
    srcs = ["server.go"],
    deps = [
        "//letsgo:go_default_library",
        "//client:go_default_library",
        "//cmd/mainservice-demo/proto:go_default_library",
        "//server:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//ease-gateway/proto/data:go_default_library",
    ],
)

docker_build(
    name = "latest",
    base = "//bld_tools/docker/ubuntu:base",
    directory = "/skylb",
    files = [
        ":mainservice-client",
        ":mainservice-server",
    ],
    ports = [
        "8901",
    ],
    repository = "harbor.jingoal.ltd",
    tag = "skylb/mainservice-demo",
    workdir = "/skylb",
)
