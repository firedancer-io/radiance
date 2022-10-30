load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "blockstore",
    srcs = [
        "bincode.go",
        "bincode_rocks.go",
        "blockstore.go",
        "blockstore_rocks.go",
        "blockwalk.go",
        "blockwalk_rocks.go",
        "meta.go",
        "meta_rocks.go",
        "shreds.go",
        "shreds_rocks.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/blockstore",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/shred",
        "@com_github_gagliardetto_binary//:binary",
        "@com_github_linxgnu_grocksdb//:grocksdb",
        "@io_k8s_klog_v2//:klog",
    ],
)

go_test(
    name = "blockstore_test",
    srcs = [
        "blockwalk_test.go",
        "shreds_bench_test.go",
        "shreds_test.go",
    ],
    data = ["//fixtures:blobs"],
    embed = [":blockstore"],
    gotags = ["rocksdb"],
    rundir = ".",
    deps = [
        "//fixtures",
        "//pkg/shred",
        "@com_github_gagliardetto_solana_go//:solana-go",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
