package(
    default_visibility = ["//visibility:public"],
)

load("//bld_tools/bazel/rules_jingoal/protobuf:def.bzl", "genproto_go", "genproto_java")
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
    importpath = "github.com/binchencoder/skylb-api/cmd/skytest/proto",
    deps = [
        "@binchencoder_third_party_go//vendor/github.com/golang/protobuf/proto:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
    ],
)

genproto_java(
    name = "skytest_src",
    srcs = ["api.proto"],
    has_service = 1,
    deps = [
        "//ease-gateway/gateway/options:options_proto",
        "@binchencoder_third_party_go//vendor/github.com/google/protobuf/src/google/protobuf:protos_java",
        "//ease-gateway/proto/data:data_srcjar",
    ],
)

java_library(
    name = "skytest",
    srcs = [":skytest_src"],
    deps = [
        "//ease-gateway/gateway/options",
        "//third-party-java/dependencies/lib:grpc-1.18.0",
    ],
)
