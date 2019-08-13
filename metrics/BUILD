package(default_visibility = [
    "//skylb:__subpackages__",
    "//visibility:public",
])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "github.com/binchencoder/skylb-api/metrics",
    deps = [
        "@binchencoder_letsgo//grpc:go_default_library",
        "//proto:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/prometheus/client_golang/prometheus:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/codes:go_default_library",
    ],
)
