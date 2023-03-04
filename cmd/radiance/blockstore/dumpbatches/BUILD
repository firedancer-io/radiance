load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "dumpbatches",
    srcs = ["dumpbatches.go"],
    importpath = "go.firedancer.io/radiance/cmd/radiance/blockstore/dumpbatches",
    visibility = ["//visibility:public"],
    deps = [
        "//cmd/radiance/blockstore/util",
        "//pkg/blockstore",
        "//pkg/shred",
        "@com_github_linxgnu_grocksdb//:grocksdb",
        "@com_github_spf13_cobra//:cobra",
        "@io_k8s_klog_v2//:klog",
    ],
)
