load("@rules_foreign_cc//foreign_cc:defs.bzl", "cmake")

filegroup(
    name = "srcs",
    srcs = glob(["**"]),
)

cmake(
    name = "rocksdb",
    lib_source = "//:srcs",
    out_static_libs = ["librocksdb.a"],
    visibility = ["//visibility:public"],
    build_args = ["--parallel `njobs`"],
)
