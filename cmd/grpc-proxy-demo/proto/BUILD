package(
    default_visibility = ["//visibility:public"],
)

load("//bld_tools/bazel/rules_jingoal/protobuf:def.bzl", "genproto_go")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

genproto_go(
    name = "api_proto",
    srcs = ["api.proto"],
    has_service = 1,
)

go_library(
    name = "go_default_library",
    srcs = [
        ":api_proto",
    ],
    importpath = "github.com/binchencoder/skylb-api/cmd/grpc-proxy-demo/proto",
    deps = [
        "@com_github_golang_protobuf//proto:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
