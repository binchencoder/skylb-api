load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "fake-service",
    srcs = ["main.go"],
    deps = [
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//service/naming:go_default_library",
        "//client:go_default_library",
        "//internal/skylb:go_default_library",
        "//server:go_default_library",
        "@com_github_golang_glog//:go_default_library",
    ],
)
