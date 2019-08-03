package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "grpc.go",
        "sql.go",
    ],
    importpath = "jingoal.com/skylb-api/handlers",
    deps = [
        "//letsgo/grpc:go_default_library",
        "//letsgo/sql:go_default_library",
        "//skylb-api/balancer:go_default_library",
        "//skylb-api/client/option:go_default_library",
        "//skylb-api/internal/flags:go_default_library",
        "//skylb-api/internal/flags/client:go_default_library",
        "//skylb-api/internal/health:go_default_library",
        "//skylb-api/metrics:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc:go_default_library",
        "//third-party-go/vendor/github.com/opentracing/opentracing-go:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/naming:go_default_library",
        "//vexillary-client/proto/data:go_default_library",
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
        "//letsgo/testing/mocks/grpc:go_default_library",
        "//letsgo/testing/mocks/tracing:go_default_library",
        "//skylb-api/client/option:go_default_library",
        "//skylb-api/proto:go_default_library",
        "//skylb-api/testing:go_default_library",
        "//third-party-go/vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
    ],
)
