load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "compactindex",
    srcs = [
        "compactindex.go",
        "query.go",
    ] + select({
        "@platforms//os:osx": ["build.go"],
        "@platforms//os:linux": ["build.go"],
        "//conditions:default": [],
    }) + select({
        "@platforms//os:linux": ["fallocate_linux.go"],
        "//conditions:default": ["fallocate_generic.go"],
    }),
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
    ] + select({
        "@platforms//os:osx": ["build_test.go"],
        "@platforms//os:linux": ["build_test.go"],
        "//conditions:default": [],
    }),
    embed = [":compactindex"],
    deps = [
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
