load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "kafka",
    srcs = ["kafka.go"],
    importpath = "go.firedancer.io/radiance/pkg/kafka",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/env/v1:env",
        "@com_github_twmb_franz_go//pkg/kgo",
        "@com_github_twmb_franz_go//pkg/sasl/scram",
    ],
)
