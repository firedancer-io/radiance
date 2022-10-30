load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "envfile",
    srcs = [
        "env.go",
        "filter_nodes.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/envfile",
    visibility = ["//visibility:public"],
    deps = [
        "//proto/env/v1:env",
        "@org_golang_google_protobuf//encoding/prototext",
    ],
)

go_test(
    name = "envfile_test",
    srcs = ["env_test.go"],
    data = glob(["testdata/**"]),
    embed = [":envfile"],
    deps = ["//proto/env/v1:env"],
)
