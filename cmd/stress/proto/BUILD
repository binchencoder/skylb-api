package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "apipb_proto",
    srcs = glob(["*.proto"]),
    deps = [
        "@com_github_binchencoder_ease_gateway//gateway/options:options_proto",
        "@com_github_binchencoder_gateway_proto//data:data_proto",
        "@com_google_protobuf//:duration_proto",
        "@com_google_protobuf//:empty_proto",
        "@com_google_protobuf//:field_mask_proto",
        "@com_google_protobuf//:timestamp_proto",
        "@com_google_protobuf//:wrappers_proto",
        "@go_googleapis//google/api:annotations_proto",
    ],
)

go_proto_library(
    name = "apipb_go_proto",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
        "@com_github_binchencoder_ease_gateway//gateway/protoc-gen-grpc-gateway:go_gen_grpc_gateway",  # keep
    ],
    importpath = "github.com/binchencoder/skylb-api/cmd/stress/proto",
    proto = ":apipb_proto",
    deps = [
        "@com_github_binchencoder_ease_gateway//gateway/options:options_go_proto",
        "@com_github_binchencoder_gateway_proto//data:data_go_proto",
        "@com_github_binchencoder_gateway_proto//frontend:error_go_proto",
        "@go_googleapis//google/api:annotations_go_proto",
    ],
)

go_library(
    name = "go_default_library",
    embed = [":apipb_go_proto"],
    importpath = "github.com/binchencoder/skylb-api/cmd/stress/proto",
    deps = [
        "@com_github_binchencoder_ease_gateway//gateway/options:go_default_library",
        "@com_github_binchencoder_gateway_proto//data:go_default_library",
        "@com_github_binchencoder_gateway_proto//frontend:go_default_library",
        "@grpc_ecosystem_grpc_gateway//runtime:go_default_library",
        "@grpc_ecosystem_grpc_gateway//utilities:go_default_library",
        "@com_github_binchencoder_letsgo//grpc:go_default_library",
        "//balancer:go_default_library",
        "//client:go_default_library",
        "//client/option:go_default_library",
        "//proto:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//grpclog:go_default_library",
        "@org_golang_google_grpc//naming:go_default_library",
    ],
)
