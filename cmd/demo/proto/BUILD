package(
    default_visibility = ["//visibility:public"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "api_proto",
    srcs = ["api.proto"],
    deps = [
        "@com_github_binchencoder_ease_gateway//httpoptions:options_proto",
        "@com_github_binchencoder_gateway_proto//data:data_proto",
    ],
)

go_proto_library(
    name = "api_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
        "@com_github_binchencoder_ease_gateway//gateway/protoc-gen-grpc-gateway:go_gen_grpc_gateway",  # keep
    ],
    importpath = "github.com/binchencoder/skylb-api/cmd/demo/proto",
    proto = ":api_proto",
    deps = [
        "@com_github_binchencoder_ease_gateway//gateway/protoc-gen-swagger/options:go_default_library",
        "@com_github_binchencoder_ease_gateway//httpoptions:options_go_proto",
        "@com_github_binchencoder_ease_gateway//gateway/runtime:go_default_library",
        "@com_github_binchencoder_gateway_proto//data:data_go_proto",
        "@com_github_binchencoder_gateway_proto//frontend:error_go_proto",
        "@com_github_binchencoder_letsgo//grpc:go_default_library",
        "//balancer:go_default_library",
        "//client:go_default_library",
        "//client/option:go_default_library",
        "//proto:go_default_library",
        "@org_golang_google_grpc//naming:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    embed = [
        ":api_go_proto",
    ],
    importpath = "github.com/binchencoder/skylb-api/cmd/demo/proto",
    deps = [
        "@com_github_binchencoder_ease_gateway//httpoptions:go_default_library",
        "@com_github_binchencoder_ease_gateway//gateway/runtime:go_default_library",
        "@com_github_binchencoder_gateway_proto//data:go_default_library",
        "@com_github_binchencoder_gateway_proto//frontend:go_default_library",
        "@com_github_binchencoder_letsgo//grpc:go_default_library",
        "//balancer:go_default_library",
        "//client:go_default_library",
        "//client/option:go_default_library",
        "//proto:go_default_library",
        "@org_golang_google_grpc//naming:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//protoc-gen-go/generator:go_default_library",
        "@com_github_grpc_ecosystem_grpc_gateway//utilities:go_default_library",
    ],
)

#genproto_java(
#    name = "greeting_src",
#    srcs = ["api.proto"],
#    has_service = 1,
#    deps = [
#        "@com_github_binchencoder_ease_gateway//httpoptions:options_proto",
#        "@binchencoder_third_party_go//vendor/github.com/google/protobuf/src/google/protobuf:protos_java",
#        "@com_github_binchencoder_gateway_proto//data:data_srcjar",
#    ],
#)

#java_library(
#    name = "greeting",
#    srcs = [":greeting_src"],
#    deps = [
#        "@com_github_binchencoder_ease_gateway//httpoptions",
#        "//third-party-java/dependencies/lib:grpc-1.18.0",
#    ],
#)

java_proto_library(
    name = "api_java_proto",
    deps = [
        ":api_proto",
    ],
)
