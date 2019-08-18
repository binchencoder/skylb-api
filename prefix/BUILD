package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "prefix.go",
    ],
    importpath = "binchencoder.com/skylb-api/prefix",
    deps = [
        "@com_github_coreos_etcd//client:go_default_library",
        "@org_golang_x_net//context:go_default_library",
    ],
)
