package(
    default_visibility = ["//visibility:public"],
)

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
        "@com_github_golang_protobuf//proto:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

# genproto_java(
#     name = "skytest_src",
#     srcs = ["api.proto"],
#     has_service = 1,
#     deps = [
#         "@binchencoder_ease_gateway//gateway/options:options_proto",
#         "@binchencoder_third_party_go//vendor/github.com/google/protobuf/src/google/protobuf:protos_java",
#         "@com_github_binchencoder_gateway_proto//data:data_srcjar",
#     ],
# )

# java_library(
#     name = "skytest",
#     srcs = [":skytest_src"],
#     deps = [
#         "@binchencoder_ease_gateway//gateway/options",
#         "//third-party-java/dependencies/lib:grpc-1.18.0",
#     ],
# )
