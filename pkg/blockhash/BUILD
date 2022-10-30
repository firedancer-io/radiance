load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "blockhash",
    srcs = ["blockhash.go"],
    importpath = "go.firedancer.io/radiance/pkg/blockhash",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/env/v1:env",
        "@com_github_gagliardetto_solana_go//:solana-go",
        "@com_github_gagliardetto_solana_go//rpc",
        "@io_k8s_klog_v2//:klog",
    ],
)
