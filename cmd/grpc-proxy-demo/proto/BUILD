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
    importpath = "jingoal.com/skylb-api/cmd/grpc-proxy-demo/proto",
    deps = [
        "//third-party-go/vendor/github.com/golang/protobuf/proto:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
    ],
)
