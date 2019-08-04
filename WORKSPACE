workspace(name = "binchencoder_sklb_api")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "6776d68ebb897625dead17ae510eac3d5f6342367327875210df44dbe2aeeb19",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.17.1/rules_go-0.17.1.tar.gz"],
)
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
go_rules_dependencies()
go_register_toolchains()

http_archive(
    name = "bazel_gazelle",
    sha256 = "3c681998538231a2d24d0c07ed5a7658cb72bfb5fd4bf9911157c0e9ac6a2687",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.17.0/bazel-gazelle-0.17.0.tar.gz"],
)
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
gazelle_dependencies()

go_repository(
    name = "binchencoder_ease_gateway",
    commit = "3a37540f37b015fc035702efc0b527f1d7ff699b",
    importpath = "github.com/binchencoder/ease-gateway",
)

go_repository(
    name = "binchencoder_third_party_java",
    commit = "dcac035f578caefefc6cd12a799cbb400a09f004",
    importpath = "github.com/binchencoder/third-party-java",
)

go_repository(
    name = "grpc_ecosystem_grpc_gateway",
    commit = "ad529a448ba494a88058f9e5be0988713174ac86",
    importpath = "github.com/grpc-ecosystem/grpc-gateway",
)

go_repository(
    name = "com_github_bazelbuild_buildtools",
    importpath = "github.com/bazelbuild/buildtools",
    commit = "36bd730dfa67bff4998fe897ee4bbb529cc9fbee",
)

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")
buildifier_dependencies()
