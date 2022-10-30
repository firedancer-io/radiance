load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "leaderschedule",
    srcs = ["epoch.go"],
    importpath = "go.firedancer.io/radiance/pkg/leaderschedule",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/env/v1:env",
        "@com_github_gagliardetto_solana_go//:solana-go",
        "@com_github_gagliardetto_solana_go//rpc",
        "@io_k8s_klog_v2//:klog",
    ],
)
