package(
    default_visibility = ["//visibility:public"],
)

java_library(
    name = "util",
    srcs = glob([
        "*.java",
    ]),
    deps = [
        "//vendor/github.com/binchencoder/third-party-java/dependencies/lib:commons-lang3-3.4",
        "//vendor/github.com/binchencoder/third-party-java/dependencies/lib:grpc-1.18.0",
        "//vendor/github.com/binchencoder/third-party-java/dependencies/lib:protobuf-java-util-3.4.0",
        "//vendor/github.com/binchencoder/third-party-java/dependencies/lib:slf4j-api-1.7.21",
        "@com_github_binchencoder_gateway_proto//data",
        "@com_github_binchencoder_gateway_proto//frontend:error",
    ],
)
