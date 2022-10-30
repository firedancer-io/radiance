load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "gossip",
    srcs = [
        "bitvec.go",
        "bloom.go",
        "client.go",
        "crds.go",
        "ping.go",
        "pull.go",
        "schema.go",
        "socketaddr.go",
        "transaction.go",
        "types.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/gossip",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_gagliardetto_solana_go//:solana-go",
        "@com_github_novifinancial_serde_reflection_serde_generate_runtime_golang//bincode",
        "@com_github_novifinancial_serde_reflection_serde_generate_runtime_golang//serde",
        "@io_k8s_klog_v2//:klog",
    ],
)

go_test(
    name = "gossip_test",
    srcs = [
        "bloom_test.go",
        "crds_test.go",
        "message_test.go",
        "ping_test.go",
    ],
    data = ["//fixtures:blobs"],
    embed = [":gossip"],
    rundir = ".",
    deps = [
        "//fixtures",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
        "@org_golang_x_sync//errgroup",
    ],
)
