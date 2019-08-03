package(
    default_visibility = ["//visibility:public"],
)

java_library(
    name = "grpc",
    srcs = glob([
        "*.java",
    ]),
    deps = [
        "//skylb-api/java/src/main/java/com/jingoal/skylb:skylb-client",
        "//skylb-api/java/src/main/java/com/jingoal/skylb/balancer/consistenthash",
        "//skylb-api/java/src/main/java/com/jingoal/skylb/balancer/roundrobin",
        "//skylb-api/java/src/main/java/com/jingoal/skylb/grpchealth",
        "//skylb-api/java/src/main/java/com/jingoal/skylb/healthcheck",
        "//skylb-api/java/src/main/java/com/jingoal/skylb/util",
        "//third-party-java/dependencies/lib:commons-lang3-3.4",
        "//third-party-java/dependencies/lib:grpc-1.18.0",
        "//third-party-java/dependencies/lib:simpleclient-0.0.15",
        "//third-party-java/dependencies/lib:slf4j-api-1.7.21",
    ],
)
