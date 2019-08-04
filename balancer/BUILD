package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "balancer.go",
        "consistent.go",
        "debug.go",
        "mainservice.go",
        "roundrobin.go",
    ],
    importpath = "github.com/binchencoder/skylb-api/balancer",
    deps = [
        "//internal/health:go_default_library",
        "@binchencoder_letsgo//grpc:go_default_library",
        "@binchencoder_letsgo//hashring:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/codes:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/grpclog:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/naming:go_default_library",
    ],
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = [
        "roundrobin_test.go",
    ],
    embed = [
        ":go_default_library",
    ],
    deps = [
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/balancer/roundrobin:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/codes:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/grpclog/glogger:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/internal/leakcheck:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/peer:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/resolver:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/resolver/manual:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/status:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/test/grpc_testing:go_default_library",
    ],
)
