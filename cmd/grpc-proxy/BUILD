load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
# load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "grpc-proxy",
    srcs = ["main.go"],
    deps = [
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//metrics:go_default_library",
        "@com_github_binchencoder_letsgo//runtime/pprof:go_default_library",
        "//balancer:go_default_library",
        "//client:go_default_library",
        "//client/option:go_default_library",
        "//handlers:go_default_library",
        "//server:go_default_library",
        "@com_github_mwitkow_grpc_proxy//proxy:go_default_library",
        "@com_github_soheilhy_cmux//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//metadata:go_default_library",
        "@org_golang_google_grpc//naming:go_default_library",
        "@com_github_binchencoder_gateway_proto//data:go_default_library",
    ],
)

pkg_tar(
    name = "grpc_proxy_tar",
    srcs = [
        ":grpc-proxy",
    ],
    package_dir = "/skylb/bin",
)

# container_image(
#     name = "latest",
#     base = "@base//image",
#     tars = [
#         ":grpc_proxy_tar",
#         "//production/config/vexconfig",
#     ],
# )

# container_push(
#     name = "release",
#     format = "Docker",
#     image = ":latest",
#     registry = "",
#     repository = "skylb/grpc-proxy",
# )
