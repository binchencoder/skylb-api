load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    go_repository(
        name = "binchencoder_third_party_go",
        commit = "d315c8c6aeab36114ee245515150906d434e92b3",
        importpath = "gitee.com/binchencoder/third-party-go",
    )
    go_repository(
        name = "binchencoder_ease_gateway",
        commit = "544d50be5ccd1d8956eef3da33ed90ec7d6281e6",
        importpath = "gitee.com/binchencoder/ease-gateway",
    )
    go_repository(
        name = "binchencoder_letsgo",
        commit = "16c8caf20f0a9601808ec77da4ae5d26ed60f5ac",
        importpath = "gitee.com/binchencoder/letsgo",
    )

    go_repository(
        name = "grpc_ecosystem_grpc_gateway",
        commit = "ad529a448ba494a88058f9e5be0988713174ac86",
        importpath = "github.com/grpc-ecosystem/grpc-gateway",
    )