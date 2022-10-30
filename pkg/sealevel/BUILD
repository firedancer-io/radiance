load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "sealevel",
    srcs = [
        "cu.go",
        "logging.go",
        "params.go",
        "sealevel.go",
        "syscalls.go",
        "syscalls_log.go",
        "syscalls_misc.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/sealevel",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/sbf",
        "//pkg/sbf/cu",
        "@com_github_gagliardetto_solana_go//:solana-go",
    ],
)

go_test(
    name = "sealevel_test",
    srcs = ["sealevel_test.go"],
    data = ["//fixtures:blobs"],
    embed = [":sealevel"],
    rundir = ".",
    deps = [
        "//fixtures",
        "//pkg/sbf",
        "//pkg/sbf/loader",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
