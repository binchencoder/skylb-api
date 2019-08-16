package(default_visibility = ["//gigaform:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_test(
    name = "integration_test",
    size = "medium",
    srcs = glob(["*_test.go"]),
    deps = [
        "@com_github_binchencoder_letsgo//trace:go_default_library",
        "//cmd/stress/proto:go_default_library",
        "@com_github_golang_protobuf//jsonpb:go_default_library",
        "@com_github_smartystreets_goconvey//convey:go_default_library",
    ],
)
