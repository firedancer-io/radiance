load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "yaml",
    srcs = [
        "shreds.go",
        "yaml.go",
    ],
    importpath = "go.firedancer.io/radiance/cmd/radiance/blockstore/yaml",
    visibility = ["//visibility:public"],
    deps = [
        "//cmd/radiance/blockstore/util",
        "//pkg/blockstore",
        "//pkg/shred",
        "@com_github_linxgnu_grocksdb//:grocksdb",
        "@com_github_segmentio_textio//:textio",
        "@com_github_spf13_cobra//:cobra",
        "@in_gopkg_yaml_v3//:yaml_v3",
        "@io_k8s_klog_v2//:klog",
    ],
)
