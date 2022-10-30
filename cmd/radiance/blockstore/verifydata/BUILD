load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "verifydata",
    srcs = [
        "verifydata.go",
        "worker.go",
    ],
    importpath = "go.firedancer.io/radiance/cmd/radiance/blockstore/verifydata",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/blockstore",
        "//pkg/shred",
        "@com_github_linxgnu_grocksdb//:grocksdb",
        "@com_github_mattn_go_isatty//:go-isatty",
        "@com_github_spf13_cobra//:cobra",
        "@com_github_vbauerster_mpb_v8//:mpb",
        "@com_github_vbauerster_mpb_v8//decor",
        "@com_github_vividcortex_ewma//:ewma",
        "@io_k8s_klog_v2//:klog",
        "@org_golang_x_sync//errgroup",
    ],
)
