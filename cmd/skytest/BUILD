load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_docker//docker:docker.bzl", "docker_build")
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "skytest-client",
    srcs = ["client.go"],
    deps = [
        "@binchencoder_letsgo//:go_default_library",
        "//client:go_default_library",
        "//cmd/skytest/proto:go_default_library",
        "//handlers:go_default_library",
        "//proto:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/prometheus/client_golang/prometheus:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/health/grpc_health_v1:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

go_binary(
    name = "skytest-server",
    srcs = ["server.go"],
    deps = [
        "@binchencoder_letsgo//:go_default_library",
        "//client:go_default_library",
        "//cmd/skytest/proto:go_default_library",
        "//server:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/prometheus/client_golang/prometheus:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

go_binary(
    name = "direct-grpc",
    srcs = ["grpc.go"],
    deps = [
        "//cmd/skytest/proto:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
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
