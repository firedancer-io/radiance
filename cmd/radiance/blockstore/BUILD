load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "blockstore",
    srcs = ["blockstore.go"],
    importpath = "go.firedancer.io/radiance/cmd/radiance/blockstore",
    visibility = ["//visibility:public"],
    deps = [
        "//cmd/radiance/blockstore/dumpbatches",
        "//cmd/radiance/blockstore/dumpshreds",
        "//cmd/radiance/blockstore/verifydata",
        "//cmd/radiance/blockstore/yaml",
        "@com_github_spf13_cobra//:cobra",
    ],
)

go_library(
    name = "blockstore_missing",
    srcs = ["blockstore_missing.go"],
    importpath = "go.firedancer.io/radiance/cmd/radiance/blockstore",
    visibility = ["//visibility:public"],
    deps = ["@com_github_spf13_cobra//:cobra"],
)
