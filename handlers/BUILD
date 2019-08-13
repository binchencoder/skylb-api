package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "grpc.go",
        "sql.go",
    ],
    importpath = "github.com/binchencoder/skylb-api/handlers",
    deps = [
        "@binchencoder_letsgo//grpc:go_default_library",
        "@binchencoder_letsgo//sql:go_default_library",
        "//balancer:go_default_library",
        "//client/option:go_default_library",
        "//internal/flags:go_default_library",
        "//internal/flags/client:go_default_library",
        "//internal/health:go_default_library",
        "//metrics:go_default_library",
        "//proto:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/opentracing/opentracing-go:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//naming:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

go_test(
    name = "tests",
    size = "small",
    srcs = [
        "grpc_test.go",
    ],
    embed = [
        ":go_default_library",
    ],
    deps = [
        "@binchencoder_letsgo//testing/mocks/grpc:go_default_library",
        "@binchencoder_letsgo//testing/mocks/tracing:go_default_library",
        "//client/option:go_default_library",
        "//proto:go_default_library",
        "//testing:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)
