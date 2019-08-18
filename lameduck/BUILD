package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "lameduck.go",
    ],
    importpath = "binchencoder.com/skylb-api/lameduck",
    deps = [
        "//prefix:go_default_library",
        "//util:go_default_library",
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_golang_glog//:go_default_library",
    ],
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = ([
        "lameduck_test.go",
    ]),
    embed = [
        ":go_default_library",
    ],
    deps = [
        "//util:go_default_library",
        "@com_github_coreos_etcd//client:go_default_library",
        "@com_github_smartystreets_goconvey//convey:go_default_library",
    ],
)
