load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "shred",
    srcs = [
        "entry.go",
        "shred.go",
        "shredder.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/shred",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_gagliardetto_binary//:binary",
        "@com_github_gagliardetto_solana_go//:solana-go",
    ],
)
