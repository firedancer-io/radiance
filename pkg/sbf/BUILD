load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "sbf",
    srcs = [
        "asm.go",
        "interpreter.go",
        "opcode.go",
        "program.go",
        "sbf.go",
        "stack.go",
        "syscalls.go",
        "verifier.go",
        "vm.go",
    ],
    importpath = "go.firedancer.io/radiance/pkg/sbf",
    visibility = ["//visibility:public"],
    deps = ["@com_github_spaolacci_murmur3//:murmur3"],
)

go_test(
    name = "sbf_test",
    srcs = ["opcode_test.go"],
    embed = [":sbf"],
    deps = ["@com_github_stretchr_testify//require"],
)
