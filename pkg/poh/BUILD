load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "poh",
    srcs = ["poh.go"],
    importpath = "go.firedancer.io/radiance/pkg/poh",
    visibility = ["//visibility:public"],
)

go_test(
    name = "poh_test",
    srcs = ["poh_test.go"],
    embed = [":poh"],
    deps = [
        "@com_github_minio_sha256_simd//:sha256-simd",
        "@com_github_stretchr_testify//assert",
    ],
)
