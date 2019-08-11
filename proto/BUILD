package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "api_proto",
    srcs = ["api.proto"],
    deps = [
        "//ease-gateway/proto/data:data_proto",
    ],
)

go_proto_library(
    name = "api_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
    ],
    importpath = "github.com/binchencoder/skylb-api/proto",
    proto = ":api_proto",
    deps = [
        "//ease-gateway/proto/data:data_go_proto",
        "//ease-gateway/gateway/options:options_go_proto",
    ],
)

go_library(
    name = "go_default_library",
    embed = [":api_go_proto"],
    importpath = "github.com/binchencoder/skylb-api/proto",
    deps = [
        "@binchencoder_third_party_go//vendor/github.com/golang/protobuf/proto:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
        "//ease-gateway/proto/data:go_default_library",
    ],
)

java_proto_library(
    name = "api_java_proto",
    deps = [
        ":api_proto",
    ],
)
