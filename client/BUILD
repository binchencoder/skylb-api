package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "github.com/binchencoder/skylb-api/client",
    deps = [
         "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@binchencoder_letsgo//flags:go_default_library",
        "@binchencoder_letsgo//service/naming:go_default_library",
        "//proto:go_default_library",
        "//client/option:go_default_library",
        "//internal/skylb:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
    ],
)
