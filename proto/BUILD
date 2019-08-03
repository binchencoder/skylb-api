package(default_visibility = ["//visibility:public"])

load("//bld_tools/bazel/rules_jingoal/protobuf:def.bzl", "genproto_go", "genproto_java")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

genproto_go(
    name = "api_proto",
    srcs = ["api.proto"],
    has_service = 1,
    deps = [
        "//vexillary-client/proto/data:data_proto",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [":api_proto"],
    importpath = "jingoal.com/skylb-api/proto",
    deps = [
        "//third-party-go/vendor/github.com/golang/protobuf/proto:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//vexillary-client/proto/data:go_default_library",
    ],
)

genproto_java(
    name = "api_proto_java",
    srcs = ["api.proto"],
    has_service = 1,
    deps = [
        "//vexillary-client/proto/data:data_srcjar",
    ],
)

java_library(
    name = "skylb-proto",
    srcs = [":api_proto_java"],
    deps = [
        "//third-party-java/dependencies/lib:grpc-1.18.0",
        "//vexillary-client/proto/data",
    ],
)
