load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_docker//docker:docker.bzl", "docker_build")

pkg_tar(
    name = "demo-java-tar",
    srcs = [
        "//demo/src/main/java/com/jingoal/skylb/demo:chc_deploy.jar",
        "//demo/src/main/java/com/jingoal/skylb/demo:chcms_deploy.jar",
        "//demo/src/main/java/com/jingoal/skylb/demo:gc_deploy.jar",
        "//demo/src/main/java/com/jingoal/skylb/demo:gs_deploy.jar",
        "//demo/src/main/java/com/jingoal/skylb/demo:sgs_deploy.jar",
    ],
)

docker_build(
    name = "latest",
    base = "//bld_tools/docker/jdk:latest",
    directory = "/skylb",
    repository = "harbor.jingoal.ltd",
    tag = "skylb/demo-java",
    tars = [
        ":demo-java-tar",
    ],
)
