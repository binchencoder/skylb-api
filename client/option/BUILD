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
        "//proto:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/opentracing/opentracing-go:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//naming:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)
