package(
    default_visibility = ["//visibility:public"],
)

java_library(
    name = "skylb-client",
    srcs = glob([
        "*.java",
        "balancer/**/*.java",
        "metrics/*.java",
    ]),
    deps = [
        "//skylb-api/java/src/main/java/com/jingoal/skylb/balancer/consistenthash",
        "//skylb-api/java/src/main/java/com/jingoal/skylb/balancer/roundrobin",
        "//skylb-api/proto:skylb-proto",
        "//third-party-java/dependencies/lib:commons-lang3-3.4",
        "//third-party-java/dependencies/lib:grpc-1.18.0",
        "//third-party-java/dependencies/lib:gson-2.6.2",
        "//third-party-java/dependencies/lib:java-grpc-prometheus-0.3.0",
        "//third-party-java/dependencies/lib:javax.servlet-api-3.1.0",
        "//third-party-java/dependencies/lib:simpleclient-0.0.15",
        "//third-party-java/dependencies/lib:simpleclient_common-0.0.15",
        "//third-party-java/dependencies/lib:simpleclient_servlet-0.0.15",
        "//third-party-java/dependencies/lib:slf4j-api-1.7.21",
        "//third-party-java/dependencies/lib:spark-core-2.5",
    ],
)
