load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "nftables",
    srcs = ["nftables.go"],
    importpath = "go.firedancer.io/radiance/pkg/nftables",
    target_compatible_with = select({
        "@platforms//os:linux": [],
        "//conditions:default": ["@platforms//:incompatible"],
    }),
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_google_nftables//:nftables",
        "@com_github_google_nftables//binaryutil",
        "@com_github_google_nftables//expr",
        "@org_golang_x_sys//unix",
    ],
)
