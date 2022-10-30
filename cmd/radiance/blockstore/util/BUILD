load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "util",
    srcs = ["ints.go"],
    importpath = "go.firedancer.io/radiance/cmd/radiance/blockstore/util",
    visibility = ["//visibility:public"],
)

go_test(
    name = "util_test",
    srcs = ["ints_test.go"],
    embed = [":util"],
    deps = ["@com_github_stretchr_testify//assert"],
)
