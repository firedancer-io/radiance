load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "tpu",
    srcs = ["tpu.go"],
    importpath = "go.firedancer.io/radiance/pkg/tpu",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_gagliardetto_binary//:binary",
        "@com_github_gagliardetto_solana_go//:solana-go",
    ],
)

go_test(
    name = "tpu_test",
    srcs = ["tpu_test.go"],
    embed = [":tpu"],
)
