load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "fixtures",
    srcs = [
        "fixtures.go",
        "shreds.go",
    ],
    importpath = "go.firedancer.io/radiance/fixtures",
    visibility = ["//visibility:public"],
    deps = ["@com_github_stretchr_testify//require"],
)

filegroup(
    name = "blobs",
    srcs = glob(["**"]),
    visibility = ["//visibility:public"],
)
