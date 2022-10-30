load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "dumpshreds",
    srcs = ["dumpshreds.go"],
    importpath = "go.firedancer.io/radiance/cmd/radiance/blockstore/dumpshreds",
    visibility = ["//visibility:public"],
    deps = [
        "//cmd/radiance/blockstore/util",
        "//pkg/blockstore",
        "@com_github_linxgnu_grocksdb//:grocksdb",
        "@com_github_spf13_cobra//:cobra",
        "@io_k8s_klog_v2//:klog",
    ],
)
