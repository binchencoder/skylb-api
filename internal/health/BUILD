package(default_visibility = [
    "//skylb:__subpackages__",
    "//skylb-api:__subpackages__",
])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "ctx.go",
        "health.go",
    ],
    importpath = "jingoal.com/skylb-api/internal/health",
    deps = [
        "//letsgo/hashring:go_default_library",
        "//third-party-go/vendor/github.com/golang/glog:go_default_library",
        "//third-party-go/vendor/github.com/prometheus/client_golang/prometheus:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/codes:go_default_library",
        "//third-party-go/vendor/google.golang.org/grpc/health/grpc_health_v1:go_default_library",
    ],
)

go_test(
    name = "small_tests",
    size = "small",
    srcs = [
        "ctx_test.go",
    ],
    embed = [
        ":go_default_library",
    ],
    deps = [
        "//third-party-go/vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//third-party-go/vendor/golang.org/x/net/context:go_default_library",
    ],
)
