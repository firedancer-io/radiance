load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "genesis",
    srcs = [
        "file.go",
        "genesis.go",
        "serde.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/genesis",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/archiveutil",
        "//pkg/runtime",
        "@com_github_gagliardetto_binary//:binary",
    ],
)

go_test(
    name = "genesis_test",
    srcs = ["file_test.go"],
    data = ["//fixtures:blobs"],
    embed = [":genesis"],
    rundir = ".",
    deps = [
        "//fixtures",
        "//pkg/runtime",
        "@com_github_gagliardetto_solana_go//:solana-go",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
