package(default_visibility = ["//visibility:public"])

load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

filegroup(
    name = "logback",
    srcs = ["logback.xml"],
)

pkg_tar(
    name = "conf-tar",
    srcs = [
        "logback.xml",
    ],
    package_dir = "/conf",
)
