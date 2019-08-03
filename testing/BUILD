package(default_visibility = ["//skylb-api:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "handler.go",
        "keeper.go",
    ],
    importpath = "jingoal.com/skylb-api/testing",
    deps = [
        "//skylb-api/client/option:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//third-party-go/vendor/github.com/opentracing/opentracing-go:go_default_library",
        "//third-party-go/vendor/github.com/stretchr/testify/mock:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/naming:go_default_library",
        "//vexillary-client/proto/data:go_default_library",
    ],
)
