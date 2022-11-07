load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix go.firedancer.io/radiance
# gazelle:build_file_name BUILD
gazelle(name = "gazelle")

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=third_party/go/repositories.bzl%go_repositories",
        "-prune",
    ],
    command = "update-repos",
)

# Shortcut for the Go SDK
alias(
    name = "go",
    actual = "@go_sdk//:bin/go",
    visibility = ["//visibility:public"],
)

load("@com_github_bazelbuild_buildtools//buildifier:def.bzl", "buildifier")

buildifier(
    name = "buildifier",
)
