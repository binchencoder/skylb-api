package(
    default_visibility = ["//visibility:public"],
)

java_library(
    name = "skylb-demo-all-in-one",
    srcs = glob(["*.java"]),
    deps = [
        "//cmd/demo/proto:greeting",
        "//java/src/main/java/com/jingoal/skylb:skylb-client",
        "//java/src/main/java/com/jingoal/skylb/balancer/consistenthash",
        "//java/src/main/java/com/jingoal/skylb/balancer/roundrobin",
        "//java/src/main/java/com/jingoal/skylb/grpc",
        "//java/src/main/java/com/jingoal/skylb/grpchealth",
        "//proto:skylb-proto",
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

java_binary(
    name = "gc",
    srcs = ["GreetingClient.java"],
    main_class = "com.binchencoder.skylb.demo.GreetingClient",
    resources = [
        "//demo/src/main/resources:logback",
    ],
    deps = [
        ":skylb-demo-all-in-one",
        "//cmd/demo/proto:greeting",
        "//java/src/main/java/com/jingoal/skylb:skylb-client",
        "//java/src/main/java/com/jingoal/skylb/balancer/consistenthash",
        "//java/src/main/java/com/jingoal/skylb/balancer/roundrobin",
        "//java/src/main/java/com/jingoal/skylb/grpc",
        "//third-party-java/dependencies/lib:grpc-1.18.0",
        "//third-party-java/dependencies/lib:jetty-http-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-io-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-server-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-util-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:logback-classic-1.1.7",
        "//third-party-java/dependencies/lib:logback-core-1.1.7",
        "//third-party-java/dependencies/lib:netty-all-4.1.8.Final",
        "//third-party-java/dependencies/lib:slf4j-api-1.7.21",
        "//third-party-java/dependencies/lib:websocket-server-9.3.6.v20151106",
    ],
)

java_binary(
    name = "chc",
    srcs = ["ConsistentHashClient.java"],
    main_class = "com.binchencoder.skylb.demo.ConsistentHashClient",
    resources = [
        "//demo/src/main/resources:logback",
    ],
    deps = [
        ":skylb-demo-all-in-one",
        "//cmd/demo/proto:greeting",
        "//java/src/main/java/com/jingoal/skylb:skylb-client",
        "//java/src/main/java/com/jingoal/skylb/balancer/consistenthash",
        "//java/src/main/java/com/jingoal/skylb/grpc",
        "//third-party-java/dependencies/lib:grpc-1.18.0",
        "//third-party-java/dependencies/lib:jetty-http-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-io-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-server-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-util-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:logback-classic-1.1.7",
        "//third-party-java/dependencies/lib:logback-core-1.1.7",
        "//third-party-java/dependencies/lib:netty-all-4.1.8.Final",
        "//third-party-java/dependencies/lib:slf4j-api-1.7.21",
        "//third-party-java/dependencies/lib:websocket-server-9.3.6.v20151106",
    ],
)

java_binary(
    name = "chcms",
    srcs = ["ConsistentHashClientMultiSvc.java"],
    main_class = "com.binchencoder.skylb.demo.ConsistentHashClientMultiSvc",
    resources = [
        "//demo/src/main/resources:logback",
    ],
    deps = [
        ":skylb-demo-all-in-one",
        "//cmd/demo/proto:greeting",
        "//java/src/main/java/com/jingoal/skylb:skylb-client",
        "//java/src/main/java/com/jingoal/skylb/balancer/consistenthash",
        "//java/src/main/java/com/jingoal/skylb/grpc",
        "//third-party-java/dependencies/lib:grpc-1.18.0",
        "//third-party-java/dependencies/lib:jetty-http-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-io-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-server-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-util-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:logback-classic-1.1.7",
        "//third-party-java/dependencies/lib:logback-core-1.1.7",
        "//third-party-java/dependencies/lib:netty-all-4.1.8.Final",
        "//third-party-java/dependencies/lib:slf4j-api-1.7.21",
        "//third-party-java/dependencies/lib:websocket-server-9.3.6.v20151106",
    ],
)

# To test GreetingServer (outside k8s), run:
# bazel run gs --jvmopt="-Dwithin-k8s=false" skylb://192.168.38.6:1900
# Or, after bazel build, run:
# bazel-bin/skylb-api/demo/src/main/java/com/jingoal/skylb/demo/gs --jvm_flags="-Dwithin-k8s=false" skylb://192.168.38.6:1900
java_binary(
    name = "gs",
    srcs = ["GreetingServer.java"],
    main_class = "com.binchencoder.skylb.demo.GreetingServer",
    resources = [
        "//demo/src/main/resources:logback",
    ],
    deps = [
        ":skylb-demo-all-in-one",
        "//cmd/demo/proto:greeting",
        "//java/src/main/java/com/jingoal/skylb:skylb-client",
        "//java/src/main/java/com/jingoal/skylb/grpc",
        "//java/src/main/java/com/jingoal/skylb/grpchealth",
        "//third-party-java/dependencies/lib:grpc-1.18.0",
        "//third-party-java/dependencies/lib:jetty-http-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-io-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-server-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-util-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:logback-classic-1.1.7",
        "//third-party-java/dependencies/lib:logback-core-1.1.7",
        "//third-party-java/dependencies/lib:netty-all-4.1.8.Final",
        "//third-party-java/dependencies/lib:slf4j-api-1.7.21",
        "//third-party-java/dependencies/lib:websocket-server-9.3.6.v20151106",
    ],
)

java_binary(
    name = "sgs",
    srcs = ["SingleGreetingServer.java"],
    main_class = "com.binchencoder.skylb.demo.SingleGreetingServer",
    resources = [
        "//demo/src/main/resources:logback",
    ],
    deps = [
        ":skylb-demo-all-in-one",
        "//cmd/demo/proto:greeting",
        "//java/src/main/java/com/jingoal/skylb:skylb-client",
        "//third-party-java/dependencies/lib:grpc-1.18.0",
        "//third-party-java/dependencies/lib:logback-classic-1.1.7",
        "//third-party-java/dependencies/lib:logback-core-1.1.7",
        "//third-party-java/dependencies/lib:netty-all-4.1.8.Final",
        "//third-party-java/dependencies/lib:slf4j-api-1.7.21",
    ],
)
