package(default_visibility = ["//cmd/demo:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "github.com/binchencoder/skylb-api/cmd/demo/rpc",
    deps = [
        "//client:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//proto:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@com_github_binchencoder_gateway_proto//data:go_default_library",
    ],
)
