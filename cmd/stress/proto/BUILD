package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "apipb_proto",
    srcs = glob(["*.proto"]),
    deps = [
        "//ease-gateway/gateway/options:options_proto",
        "//ease-gateway/proto/data:data_proto",
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
        "//ease-gateway/gateway/protoc-gen-grpc-gateway:go_gen_grpc_gateway",  # keep
    ],
    importpath = "github.com/binchencoder/skylb-api/cmd/stress/proto",
    proto = ":apipb_proto",
    deps = [
        "//ease-gateway/gateway/options:options_go_proto",
        "//ease-gateway/proto/data:data_go_proto",
        "//ease-gateway/proto/frontend:error_go_proto",
        "@go_googleapis//google/api:annotations_go_proto",
    ],
)

go_library(
    name = "go_default_library",
    embed = [":apipb_go_proto"],
    importpath = "github.com/binchencoder/skylb-api/cmd/stress/proto",
    deps = [
        "//ease-gateway/gateway/options:go_default_library",
        "//ease-gateway/proto/data:go_default_library",
        "//ease-gateway/proto/frontend:go_default_library",
        "@grpc_ecosystem_grpc_gateway//runtime:go_default_library",
        "@grpc_ecosystem_grpc_gateway//utilities:go_default_library",
        "//letsgo/grpc:go_default_library",
        "//skylb-api/balancer:go_default_library",
        "//skylb-api/client:go_default_library",
        "//skylb-api/client/option:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/protobuf/proto:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/codes:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/grpclog:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/naming:go_default_library",
    ],
)
