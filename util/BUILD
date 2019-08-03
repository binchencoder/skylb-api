package(default_visibility = [
    "//skylb:__subpackages__",
    "//skylb-api:__subpackages__",
])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "jingoal.com/skylb-api/util",
)
