load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "ipldgen",
    srcs = ["ipldgen.go"],
    importpath = "go.firedancer.io/radiance/pkg/ipld/ipldgen",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/ipld/car",
        "//pkg/ipld/ipldsch",
        "//pkg/shred",
        "@com_github_gagliardetto_binary//:binary",
        "@com_github_gagliardetto_solana_go//:solana-go",
        "@com_github_ipld_go_ipld_prime//datamodel",
        "@com_github_ipld_go_ipld_prime//linking/cid",
    ],
)

go_test(
    name = "ipldgen_test",
    srcs = ["ipldgen_test.go"],
    embed = [":ipldgen"],
    deps = [
        "//pkg/ipld/car",
        "@com_github_ipfs_go_cid//:go-cid",
        "@com_github_multiformats_go_multihash//:go-multihash",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
