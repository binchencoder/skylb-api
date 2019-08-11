package(default_visibility = [
    "//skylb:__subpackages__",
    "//skylb-api:__subpackages__",
])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "lameduck.go",
    ],
    importpath = "github.com/binchencoder/skylb-api/lameduck",
    deps = [
        "//prefix:go_default_library",
        "//util:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/coreos/etcd/client:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
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
        "@binchencoder_third_party_go//vendor/github.com/coreos/etcd/client:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
