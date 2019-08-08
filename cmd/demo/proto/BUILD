package(
    default_visibility = ["//visibility:public"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "api_proto",
    srcs = ["api.proto"],
    deps = [
        "//ease-gateway/gateway/options:options_proto",
        "//ease-gateway/proto/data:data_proto",
    ],
)

go_proto_library(
    name = "api_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
        "//ease-gateway/gateway/protoc-gen-grpc-gateway:go_gen_grpc_gateway",  # keep
    ],
    importpath = "github.com/binchencoder/skylb-api/cmd/demo/proto",
    proto = ":api_proto",
    deps = [
        "//ease-gateway/proto/data:data_go_proto",
        "//ease-gateway/gateway/options:options_go_proto",
    ],
)

go_library(
    name = "go_default_library",
    embed = [
        ":api_go_proto",
    ],
    importpath = "github.com/binchencoder/skylb-api/cmd/demo/proto",
    deps = [
        "//ease-gateway/gateway/options:go_default_library",
         "//ease-gateway/gateway/runtime:go_default_library",
        "//ease-gateway/proto/data:go_default_library",
        "//ease-gateway/proto/frontend:go_default_library",
        "//letsgo/grpc:go_default_library",
        "//skylb-api/balancer:go_default_library",
        "//skylb-api/client:go_default_library",
        "//skylb-api/client/option:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/protobuf/proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/protobuf/protoc-gen-go/descriptor:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/codes:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/grpclog:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/naming:go_default_library",
        "@grpc_ecosystem_grpc_gateway//runtime:go_default_library",
        "@grpc_ecosystem_grpc_gateway//utilities:go_default_library",
    ],
)

#genproto_java(
#    name = "greeting_src",
#    srcs = ["api.proto"],
#    has_service = 1,
#    deps = [
#        "//ease-gateway/gateway/options:options_proto",
#        "//third-party-go/vendor/github.com/google/protobuf/src/google/protobuf:protos_java",
#        "//ease-gateway/proto/data:data_srcjar",
#    ],
#)

#java_library(
#    name = "greeting",
#    srcs = [":greeting_src"],
#    deps = [
#        "//ease-gateway/gateway/options",
#        "//third-party-java/dependencies/lib:grpc-1.18.0",
#    ],
#)

java_proto_library(
    name = "api_java_proto",
    deps = [
        ":api_proto",
    ],
)
