package(default_visibility = ["//skylb-api:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "handler.go",
        "keeper.go",
    ],
    importpath = "github.com/binchencoder/skylb-api/testing",
    deps = [
        "//client/option:go_default_library",
        "//proto:go_default_library",
        "//third-party-go/vendor/github.com/opentracing/opentracing-go:go_default_library",
        "//third-party-go/vendor/github.com/stretchr/testify/mock:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/naming:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)
