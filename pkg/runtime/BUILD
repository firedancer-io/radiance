load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "runtime",
    srcs = [
        "runtime.go",
        "serde.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/runtime",
    visibility = ["//visibility:public"],
    deps = ["@com_github_gagliardetto_binary//:binary"],
)
