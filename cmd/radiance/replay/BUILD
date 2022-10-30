load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "replay",
    srcs = ["replay.go"],
    importpath = "go.firedancer.io/radiance/cmd/radiance/replay",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/blockstore",
        "//pkg/genesis",
        "//pkg/merkletree",
        "//pkg/poh",
        "//pkg/runtime",
        "@com_github_gagliardetto_solana_go//:solana-go",
        "@com_github_spf13_cobra//:cobra",
        "@io_k8s_klog_v2//:klog",
    ],
)
