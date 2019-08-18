package(
    default_visibility = ["//visibility:public"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "api_proto",
    srcs = ["api.proto"]
)

go_library(
    name = "go_default_library",
    srcs = [
        ":api_proto",
    ],
    importpath = "binchencoder.com/skylb-api/cmd/mainservice-demo/proto",
    deps = [
        "@com_github_golang_protobuf//proto:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
