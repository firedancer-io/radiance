load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "confirmed_block",
    srcs = ["confirmed_block.pb.go"],
    importpath = "go.firedancer.io/radiance/third_party/solana_proto/confirmed_block",
    visibility = ["//visibility:public"],
    deps = [
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
)
