load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "compactindex",
    srcs = [
        "compactindex.go",
        "query.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/compactindex",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_cespare_xxhash_v2//:xxhash",
    ],
)

go_test(
    name = "compactindex_test",
    srcs = [
        "compactindex_test.go",
        "query_test.go",
    ],
    embed = [":compactindex"],
    deps = [
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ]
)
