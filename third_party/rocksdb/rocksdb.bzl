load("@rules_foreign_cc//foreign_cc:defs.bzl", "cmake")

filegroup(
    name = "srcs",
    srcs = glob(["**"]),
)

cmake(
    name = "rocksdb",
    lib_source = "//:srcs",
    out_lib_dir = "lib64",
    out_static_libs = ["librocksdb.a"],
    visibility = ["//visibility:public"],
    build_args = ["--parallel 6"],
    generate_args = ["-G Ninja"],
    linkopts = ["-lz", "-lbz2", "-lsnappy", "-lzstd"],
    cache_entries = {
        "ROCKSDB_BUILD_SHARED": "OFF",
        "WITH_BZ2": "ON",
        "WITH_SNAPPY": "ON",
        "WITH_ZLIB": "ON",
        "WITH_ZSTD": "ON",
        "WITH_ALL_TESTS": "OFF",
        "WITH_BENCHMARK_TOOLS": "OFF",
        "WITH_CORE_TOOLS": "OFF",
        "WITH_RUNTIME_DEBUG": "OFF",
        "WITH_TESTS": "OFF",
        "WITH_TOOLS": "OFF",
        "WITH_TRACE_TOOLS": "OFF",
    },
)
