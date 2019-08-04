package(
    default_visibility = ["//visibility:public"],
)

java_library(
    name = "grpc",
    srcs = glob([
        "*.java",
    ]),
    deps = [
        "//java/src/main/java/com/jingoal/skylb:skylb-client",
        "//java/src/main/java/com/jingoal/skylb/balancer/consistenthash",
        "//java/src/main/java/com/jingoal/skylb/balancer/roundrobin",
        "//java/src/main/java/com/jingoal/skylb/grpchealth",
        "//java/src/main/java/com/jingoal/skylb/healthcheck",
        "//java/src/main/java/com/jingoal/skylb/util",
        "//third-party-java/dependencies/lib:commons-lang3-3.4",
        "//third-party-java/dependencies/lib:grpc-1.18.0",
        "//third-party-java/dependencies/lib:simpleclient-0.0.15",
        "//third-party-java/dependencies/lib:slf4j-api-1.7.21",
    ],
)
