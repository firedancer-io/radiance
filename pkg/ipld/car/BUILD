load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "car",
    srcs = [
        "car.go",
        "writer.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/ipld/car",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_filecoin_project_go_leb128//:go-leb128",
        "@com_github_ipfs_go_cid//:go-cid",
        "@com_github_ipld_go_car//:go-car",
        "@com_github_ipld_go_ipld_prime//codec/dagcbor",
        "@com_github_ipld_go_ipld_prime//datamodel",
        "@com_github_multiformats_go_multicodec//:go-multicodec",
    ],
)

go_test(
    name = "car_test",
    srcs = [
        "car_test.go",
        "writer_test.go",
    ],
    embed = [":car"],
    deps = [
        "@com_github_ipfs_go_cid//:go-cid",
        "@com_github_ipfs_go_datastore//:go-datastore",
        "@com_github_ipfs_go_ipfs_blockstore//:go-ipfs-blockstore",
        "@com_github_ipfs_go_ipld_cbor//:go-ipld-cbor",
        "@com_github_ipld_go_car//:go-car",
        "@com_github_multiformats_go_multicodec//:go-multicodec",
        "@com_github_stretchr_testify//assert",
        "@com_github_stretchr_testify//require",
    ],
)
