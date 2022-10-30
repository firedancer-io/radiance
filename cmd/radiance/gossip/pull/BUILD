load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "pull",
    srcs = ["pull.go"],
    importpath = "go.firedancer.io/radiance/cmd/radiance/gossip/pull",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/gossip",
        "@com_github_spf13_cobra//:cobra",
        "@io_k8s_klog_v2//:klog",
        "@org_golang_x_sync//errgroup",
    ],
)
