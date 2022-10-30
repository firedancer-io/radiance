load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "gossip",
    srcs = ["gossip.go"],
    importpath = "go.firedancer.io/radiance/cmd/radiance/gossip",
    visibility = ["//visibility:public"],
    deps = [
        "//cmd/radiance/gossip/ping",
        "//cmd/radiance/gossip/pull",
        "@com_github_spf13_cobra//:cobra",
    ],
)
