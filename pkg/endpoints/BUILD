load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "endpoints",
    srcs = ["endpoints.go"],
    importpath = "go.firedancer.io/radiance/pkg/endpoints",
    visibility = ["//visibility:public"],
    deps = ["@com_github_gagliardetto_solana_go//rpc"],
)
