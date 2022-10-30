load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "merkletree",
    srcs = ["merkletree.go"],
    importpath = "go.firedancer.io/radiance/pkg/merkletree",
    visibility = ["//visibility:public"],
)

go_test(
    name = "merkletree_test",
    srcs = ["merkletree_test.go"],
    embed = [":merkletree"],
    deps = ["@com_github_stretchr_testify//assert"],
)
