package(default_visibility = ["//gigaform:__subpackages__"])

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_test(
    name = "integration_test",
    size = "medium",
    srcs = glob(["*_test.go"]),
    deps = [
        "//letsgo/trace:go_default_library",
        "//skylb-api/cmd/stress/proto:go_default_library",
        "//third-party-go/vendor/github.com/golang/protobuf/jsonpb:go_default_library",
        "//third-party-go/vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
