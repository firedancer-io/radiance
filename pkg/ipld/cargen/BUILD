load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "cargen",
    srcs = ["cargen.go"],
    importpath = "go.firedancer.io/radiance/pkg/ipld/cargen",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/blockstore",
        "//pkg/ipld/car",
        "//pkg/ipld/ipldgen",
        "//pkg/shred",
        "@io_k8s_klog_v2//:klog",
    ],
)

go_test(
    name = "cargen_test",
    srcs = ["cargen_test.go"],
    embed = [":cargen"],
    deps = [
        "//pkg/blockstore",
        "//pkg/shred",
        "@com_github_gagliardetto_solana_go//:solana-go",
        "@com_github_go_logr_logr//testr",
        "@com_github_ipfs_go_cid//:go-cid",
        "@com_github_ipld_go_car//:go-car",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
        "@io_k8s_klog_v2//:klog",
    ],
)
