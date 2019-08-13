package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "prefix.go",
    ],
    importpath = "github.com/binchencoder/skylb-api/prefix",
    deps = [
        "@binchencoder_third_party_go//vendor/github.com/coreos/etcd/client:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
    ],
)
