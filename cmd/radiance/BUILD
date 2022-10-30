load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "radiance_lib",
    srcs = ["main.go"],
    importpath = "go.firedancer.io/radiance/cmd/radiance",
    tags = ["manual"],
    visibility = ["//visibility:private"],
    deps = [
        "//cmd/radiance/blockstore",
        "//cmd/radiance/car",
        "//cmd/radiance/gossip",
        "//cmd/radiance/replay",
        "@com_github_gagliardetto_solana_go//programs/vote",
        "@com_github_spf13_cobra//:cobra",
        "@io_k8s_klog_v2//:klog",
    ],
)

go_binary(
    name = "radiance",
    embed = [":radiance_lib"],
    gotags = ["rocksdb"],
    visibility = ["//visibility:public"],
)
