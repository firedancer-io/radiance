load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "ipldsch",
    srcs = [
        "ipldsch.go",
        "ipldsch_minima.go",
        "ipldsch_satisfaction.go",
        "ipldsch_types.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/ipld/ipldsch",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_ipld_go_ipld_prime//datamodel",
        "@com_github_ipld_go_ipld_prime//node/mixins",
        "@com_github_ipld_go_ipld_prime//schema",
    ],
)
