SkyLB API User Guide (in Java)
==============================

准备
----

大部分步骤与golang相同。详情参阅 userguide-golang.md

-   在 vexillary-client/proto/data/data.proto 定义 service id。

-   Java API 代码示例参见 demo/com/binchencoder/skylb/demo

启动
----

注意，只要不是在 Kubernetes 里运行的，就需要设 java property: within-k8s=false。

例如，对于 spring boot 启动的项目，运行命令行为：

~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
mvn spring-boot:run -Dwithin-k8s=false
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
