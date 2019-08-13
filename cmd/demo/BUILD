#load("@io_bazel_rules_docker//docker:docker.bzl", "docker_build")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_test")

go_binary(
    name = "client",
    srcs = ["client.go"],
    deps = [
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//metrics:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
    ],
)

go_binary(
    name = "client_failfast",
    srcs = ["client_failfast.go"],
    deps = [
        "//client:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//proto:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_binary(
    name = "consistenthash",
    srcs = ["consistenthash.go"],
    deps = [
        "//balancer:go_default_library",
        "//client:go_default_library",
        "//client/option:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//proto:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//hashring:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//naming:go_default_library",
    ],
)

go_binary(
    name = "webclient",
    srcs = ["webclient.go"],
    deps = [
        "//client:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//proto:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_binary(
    name = "streamclient",
    srcs = ["streamclient.go"],
    deps = [
        "//client:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//proto:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_binary(
    name = "server",
    srcs = ["server.go"],
    deps = [
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//server:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@binchencoder_ease_gateway//proto/frontend:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@com_github_binchencoder_letsgo//grpc:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
    ],
)

go_binary(
    name = "servermulti",
    srcs = ["servermulti.go"],
    deps = [
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//server:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_binary(
    name = "serverclient",
    srcs = ["serverclient.go"],
    deps = [
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//server:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_binary(
    name = "gateway",
    srcs = ["gateway.go"],
    deps = [
        "//cmd/demo/proto:go_default_library",
        "@binchencoder_ease_gateway//gateway/runtime:go_default_library",
        "@com_github_binchencoder_letsgo//:go_default_library",
    ],
)

go_test(
    name = "int_test",
    size = "medium",
    srcs = glob(["*int_test.go"]),
    deps = [
        "//client:go_default_library",
        "//cmd/demo/proto:go_default_library",
        "//cmd/demo/rpc:go_default_library",
        "//proto:go_default_library",
        "//server:go_default_library",
        "@binchencoder_ease_gateway//proto/data:go_default_library",
        "@binchencoder_third_party_go//vendor/github.com/golang/glog:go_default_library",
        "@binchencoder_third_party_go//vendor/golang.org/x/net/context:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

# Build and test this image locally. No need to push to repository.
#docker_build(
#    name = "latest",
#    base = "//bld_tools/docker/ubuntu:base",
#    directory = "/skylb",
#    files = [
#        ":client",
#        ":consistenthash",
#        ":server",
#    ],
#    repository = "harbor.jingoal.ltd",
#    tag = "skylb/demo-go",
#    workdir = "/skylb",
#)
