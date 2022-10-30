load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "dump",
    srcs = ["dump.go"],
    importpath = "go.firedancer.io/radiance/cmd/radiance/car/dump",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/ipld/ipldgen",
        "@com_github_gagliardetto_binary//:binary",
        "@com_github_gagliardetto_solana_go//:solana-go",
        "@com_github_ipld_go_car//:go-car",
        "@com_github_spf13_cobra//:cobra",
        "@io_k8s_klog_v2//:klog",
    ],
)
