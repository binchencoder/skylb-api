load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "stress-client",
    srcs = ["client.go"],
    deps = [
        "//cmd/stress/proto:go_default_library",
        "@binchencoder_letsgo//:go_default_library",
        "@binchencoder_letsgo//token:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/protobuf/jsonpb:go_default_library",
    ],
)

go_binary(
    name = "stress-server",
    srcs = ["server.go"],
    deps = [
        "//client:go_default_library",
        "//cmd/stress/proto:go_default_library",
        "//server:go_default_library",
        "@binchencoder_letsgo//:go_default_library",
        "@binchencoder_letsgo//grpc:go_default_library",
        "@binchencoder_letsgo//trace:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@binchencoder_third_party_go//vendor/google.golang.org/grpc:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)

go_binary(
    name = "stress-gateway",
    srcs = [
        ":gateway.go",
        "//janus/cmd/custom-janus:main.go",
    ],
    deps = [
        "//janus/gateway/runtime:go_default_library",
        "//janus/integrate:go_default_library",
        "//letsgo:go_default_library",
        "//cmd/stress/proto:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
    ],
)
