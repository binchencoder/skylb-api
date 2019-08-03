load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_docker//docker:docker.bzl", "docker_build")
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "skytest-client",
    srcs = ["client.go"],
    deps = [
        "//letsgo:go_default_library",
        "//skylb-api/client:go_default_library",
        "//skylb-api/cmd/skytest/proto:go_default_library",
        "//skylb-api/handlers:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/prometheus/client_golang/prometheus:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/health/grpc_health_v1:go_default_library",
        "//vexillary-client/proto/data:go_default_library",
    ],
)

go_binary(
    name = "skytest-server",
    srcs = ["server.go"],
    deps = [
        "//letsgo:go_default_library",
        "//skylb-api/client:go_default_library",
        "//skylb-api/cmd/skytest/proto:go_default_library",
        "//skylb-api/server:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/prometheus/client_golang/prometheus:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//vexillary-client/proto/data:go_default_library",
    ],
)

go_binary(
    name = "direct-grpc",
    srcs = ["grpc.go"],
    deps = [
        "//skylb-api/cmd/skytest/proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
    ],
)

pkg_tar(
    name = "binary",
    srcs = [
        ":skytest-client",
        ":skytest-server",
    ],
    package_dir = "/skylb",
)

docker_build(
    name = "latest",
    base = "//bld_tools/docker/ubuntu:base",
    cmd = "/skylb/start.sh",
    repository = "harbor.jingoal.ltd",
    tag = "skylb/skytest",
    tars = [
        ":binary",
    ],
)
