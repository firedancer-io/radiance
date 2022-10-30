load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "create",
    srcs = ["create.go"],
    importpath = "go.firedancer.io/radiance/cmd/radiance/car/create",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/blockstore",
        "//pkg/ipld/cargen",
        "@com_github_pkg_profile//:profile",
        "@com_github_spf13_cobra//:cobra",
        "@io_k8s_klog_v2//:klog",
    ],
)
