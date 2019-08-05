package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "github.com/binchencoder/skylb-api/client/option",
    deps = [
        "//skylb-api/proto:go_default_library",
        "//third-party-go/vendor/github.com/opentracing/opentracing-go:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/naming:go_default_library",
        "//ease-gateway/proto/data:go_default_library",
    ],
)
