load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "loader",
    srcs = [
        "arithmetic.go",
        "copy.go",
        "loader.go",
        "parse.go",
        "relocate.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/sbf/loader",
    visibility = ["//visibility:public"],
    deps = ["//pkg/sbf"],
)

go_test(
    name = "loader_test",
    srcs = [
        "loader_test.go",
        "relocate_test.go",
    ],
    data = ["//fixtures:blobs"],
    embed = [":loader"],
    rundir = ".",
    deps = [
        "//fixtures",
        "//pkg/sbf",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
