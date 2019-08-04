package(
    default_visibility = ["//visibility:public"],
)

load("//bld_tools/bazel/rules_jingoal/java:def.bzl", "release_info")
load("//bld_tools/bazel/rules_jingoal/protobuf:def.bzl", "genproto_go", "genproto_grpc_gateway", "genproto_java", "genproto_swagger")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

genproto_go(
    name = "api_proto",
    srcs = ["api.proto"],
    has_service = 1,
    deps = [
        "@binchencoder_ease_gateway//gateway/options:options_proto_go",
    ],
)

genproto_grpc_gateway(
    name = "gateway_go",
    srcs = ["api.proto"],
    deps = [
        "@binchencoder_ease_gateway//gateway/options:options_proto_go",
    ],
)

genproto_swagger(
    name = "swagger",
    srcs = ["api.proto"],
    deps = [
        "@binchencoder_ease_gateway//gateway/options:options_proto_go",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        ":api_proto",
        ":gateway_go",
    ],
    importpath = "github.com/binchencoder/skylb-api/cmd/demo/proto",
    deps = [
        "@binchencoder_ease_gateway//gateway/options:go_default_library",
        "@binchencoder_ease_gateway//gateway/runtime:go_default_library",
        "@binchencoder_ease_gateway//gateway/utilities:go_default_library",
        "//letsgo/grpc:go_default_library",
        "//balancer:go_default_library",
        "//client:go_default_library",
        "//client/option:go_default_library",
        "//proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/protobuf/proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/protobuf/protoc-gen-go/descriptor:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/codes:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/grpclog:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/naming:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@binchencoder_ease_gateway//proto/frontend:go_default_library",
    ],
)

genproto_java(
    name = "greeting_src",
    srcs = ["api.proto"],
    has_service = 1,
    deps = [
        "@binchencoder_ease_gateway//gateway/options:options_proto",
        "//third-party-go/vendor/github.com/google/protobuf/src/google/protobuf:protos_java",
        "@binchencoder_ease_gateway//proto/data:data_srcjar",
    ],
)

java_library(
    name = "greeting",
    srcs = [":greeting_src"],
    deps = [
        "@binchencoder_ease_gateway//gateway/options",
        "//third-party-java/dependencies/lib:grpc-1.18.0",
    ],
)

release_info(
    name = "demo_greeting_src",
    src = "greeting_src",
    group_id = "com.binchencoder.skylb",
)

release_info(
    name = "demo_greeting",
    src = "greeting",
    group_id = "com.binchencoder.skylb",
)
