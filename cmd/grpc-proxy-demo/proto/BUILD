package(
    default_visibility = ["//visibility:public"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "api_proto",
    srcs = ["api.proto"],
    deps = [
        "@com_github_binchencoder_gateway_proto//data:data_proto",
    ],
)

go_proto_library(
    name = "api_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
    ],
    importpath = "binchencoder.com/skylb-api/cmd/grpc-proxy-demo/proto",
    proto = ":api_proto",
)

go_library(
    name = "go_default_library",
    srcs = [
        ":api_go_proto",
    ],
    importpath = "binchencoder.com/skylb-api/cmd/grpc-proxy-demo/proto",
    deps = [
        "@com_github_golang_protobuf//proto:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
