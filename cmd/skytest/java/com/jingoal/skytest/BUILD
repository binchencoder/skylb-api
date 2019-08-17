java_binary(
    name = "skytest-client-1.5.0",
    srcs = [
        "SkytestClient.java",
    ],
    main_class = "com.binchencoder.skytest.SkytestClient",
    resources = [
        "//cmd/skytest/java:logback",
    ],
    deps = [
        "//cmd/skytest/proto:skytest",
        "//java/src/main/java/com/jingoal/skylb:skylb-client",
        "//java/src/main/java/com/jingoal/skylb/grpc",
        "//third-party-java/dependencies/lib:grpc-1.18.0",
        "//third-party-java/dependencies/lib:jetty-http-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-io-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-server-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:jetty-util-9.3.6.v20151106",
        "//third-party-java/dependencies/lib:netty-all-4.1.8.Final",
        "//third-party-java/dependencies/lib:slf4j-api-1.7.21",
        "//third-party-java/dependencies/lib:slf4j-simple-1.7.21",
        "//third-party-java/dependencies/lib:websocket-server-9.3.6.v20151106",
    ],
)
