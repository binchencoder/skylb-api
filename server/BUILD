package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = ["*_test.go"],
    ),
    importpath = "github.com/binchencoder/skylb-api/server",
    deps = [
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@binchencoder_letsgo//grpc:go_default_library",
        "@binchencoder_letsgo//metrics:go_default_library",
        "@binchencoder_letsgo//runtime/pprof:go_default_library",
        "@binchencoder_letsgo//service/naming:go_default_library",
        "//internal/flags:go_default_library",
        "//internal/rpccli:go_default_library",
        "//metrics:go_default_library",
        "//proto:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/grpc-ecosystem/go-grpc-middleware:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/grpc-ecosystem/go-grpc-middleware/recovery:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/prometheus/client_golang/prometheus:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/soheilhy/cmux:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/codes:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/health/grpc_health_v1:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc/peer:go_default_library",
    ],
)
