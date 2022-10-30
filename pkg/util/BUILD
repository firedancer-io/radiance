load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "util",
    srcs = ["util.go"],
    importpath = "go.firedancer.io/radiance/pkg/util",
    visibility = ["//visibility:public"],
)

go_test(
    name = "util_test",
    srcs = ["util_test.go"],
    embed = [":util"],
    deps = ["@com_github_stretchr_testify//assert"],
)
