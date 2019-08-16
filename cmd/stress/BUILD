load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "stress-client",
    srcs = ["client.go"],
    deps = [
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//token:go_default_library",
        "//cmd/stress/proto:go_default_library",
        "@com_github_golang_protobuf//jsonpb:go_default_library",
    ],
)

go_binary(
    name = "stress-server",
    srcs = ["server.go"],
    deps = [
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//grpc:go_default_library",
        "@com_github_binchencoder_letsgo//trace:go_default_library",
        "//client:go_default_library",
        "//cmd/stress/proto:go_default_library",
        "//server:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@com_github_binchencoder_gateway_proto//data:go_default_library",
    ],
)

go_binary(
    name = "stress-gateway",
    srcs = [
        ":gateway.go",
        "@com_github_binchencoder_ease_gateway//cmd/custom-janus:main.go",
    ],
    deps = [
        "@grpc_ecosystem_grpc_gateway//runtime:go_default_library",
        "@com_github_binchencoder_ease_gateway//integrate:go_default_library",
        "@com_github_binchencoder_gateway_proto//data:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
        "//cmd/stress/proto:go_default_library",
        "@com_github_golang_glog//:go_default_library",
    ],
)
