package(default_visibility = ["//visibility:public"])

load("//bld_tools/bazel/rules_jingoal/protobuf:def.bzl", "genproto_go", "genproto_grpc_gateway")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

genproto_go(
    name = "proto_go",
    srcs = glob(["*.proto"]),
    has_service = 1,
    deps = [
        "//janus/gateway/options:options_proto_go",
    ],
)

genproto_grpc_gateway(
    name = "gateway_go",
    srcs = glob(["*.proto"]),
    deps = [
        "//janus/gateway/options:options_proto_go",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        ":gateway_go",
        ":proto_go",
    ],
    importpath = "github.com/binchencoder/skylb-api/cmd/stress/proto",
    deps = [
        "//janus/gateway/options:go_default_library",
        "//janus/gateway/runtime:go_default_library",
        "//janus/gateway/utilities:go_default_library",
        "@binchencoder_letsgo//grpc:go_default_library",
        "//balancer:go_default_library",
        "//client:go_default_library",
        "//client/option:go_default_library",
        "//proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/protobuf/proto:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/codes:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/grpclog:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/naming:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@binchencoder_ease_gateway//proto/frontend:go_default_library",
    ],
)
