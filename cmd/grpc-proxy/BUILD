load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "grpc-proxy",
    srcs = ["main.go"],
    deps = [
        "@binchencoder_letsgo//:go_default_library",
        "@binchencoder_letsgo//metrics:go_default_library",
        "@binchencoder_letsgo//runtime/pprof:go_default_library",
        "//balancer:go_default_library",
        "//client:go_default_library",
        "//client/option:go_default_library",
        "//handlers:go_default_library",
        "//server:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/mwitkow/grpc-proxy/proxy:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/soheilhy/cmux:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/metadata:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/naming:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

pkg_tar(
    name = "grpc_proxy_tar",
    srcs = [
        ":grpc-proxy",
    ],
    package_dir = "/skylb/bin",
)

container_image(
    name = "latest",
    base = "@base//image",
    tars = [
        ":grpc_proxy_tar",
        "//production/config/vexconfig",
    ],
)

container_push(
    name = "release",
    format = "Docker",
    image = ":latest",
    registry = "",
    repository = "skylb/grpc-proxy",
)
