load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "fake-service",
    srcs = ["main.go"],
    deps = [
        "//letsgo:go_default_library",
        "//letsgo/service/naming:go_default_library",
        "//skylb-api/client:go_default_library",
        "//skylb-api/internal/skylb:go_default_library",
        "//skylb-api/server:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
    ],
)
