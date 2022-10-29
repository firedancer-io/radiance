load("@rules_foreign_cc//foreign_cc:make.bzl", "make")

filegroup(
    name = "srcs",
    srcs = glob(["**"]),
)

make(
    name = "librocksdb",
    targets = ["static_lib", "install"],
    visibility = ["//visibility:public"],
    lib_source = "//:srcs",
    args = ["-j `nproc`"],
    env = {
        # Fix `libtool: no output file specified` on Xcode.
        "AR": "",
    },
)

alias(
    name = "rocksdb",
    actual = ":librocksdb",
    visibility = ["//visibility:public"],
)
