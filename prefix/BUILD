package(default_visibility = [
    "//skylb:__subpackages__",
    "//skylb-api:__subpackages__",
])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "prefix.go",
    ],
    importpath = "jingoal.com/skylb-api/prefix",
    deps = [
        "//third-party-go/vendor/github.com/coreos/etcd/client:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
    ],
)
