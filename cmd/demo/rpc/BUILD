package(default_visibility = ["//skylb-api/cmd/demo:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "jingoal.com/skylb-api/cmd/demo/rpc",
    deps = [
        "//skylb-api/client:go_default_library",
        "//skylb-api/cmd/demo/proto:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//vexillary-client/proto/data:go_default_library",
    ],
)
