workspace(name = "firedancer_radiance")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

################################################################################
# Dependencies                                                                 #
################################################################################

# Dep: gflags (C++)
http_archive(
    name = "com_github_gflags_gflags",
    sha256 = "34af2f15cf7367513b352bdcd2493ab14ce43692d2dcd9dfc499492966c64dcf",
    strip_prefix = "gflags-2.2.2",
    urls = ["https://github.com/gflags/gflags/archive/v2.2.2.tar.gz"],
)

# Dep: RocksDB (C++)
http_archive(
    name = "com_github_facebook_rocksdb",
    build_file = "//:third_party/rocksdb/rocksdb.bzl",
    sha256 = "b8ac9784a342b2e314c821f6d701148912215666ac5e9bdbccd93cf3767cb611",
    strip_prefix = "rocksdb-7.7.3",
    urls = ["https://github.com/facebook/rocksdb/archive/v7.7.3.tar.gz"],
)

# Dep: grocksdb (Go)
http_archive(
    name = "com_github_linxgnu_grocksdb",
    build_file = "//:third_party/go/grocksdb/grocksdb.bzl",
    sha256 = "3e76617aaa74a2658ac59a03b77c632c41971ae01a5ccb6e1b8edeff59f567bf",
    strip_prefix = "grocksdb-1.7.10",
    urls = ["https://github.com/linxGnu/grocksdb/archive/refs/tags/v1.7.10.tar.gz"],
)

################################################################################
# Go toolchain                                                                 #
# Gazelle build file generator                                                 #
################################################################################

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "19ef30b21eae581177e0028f6f4b1f54c66467017be33d211ab6fc81da01ea4d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.38.0/rules_go-v0.38.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.38.0/rules_go-v0.38.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "448e37e0dbf61d6fa8f00aaa12d191745e14f07c31cabfa731f0c8e8a4f41b97",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.28.0/bazel-gazelle-v0.28.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.28.0/bazel-gazelle-v0.28.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.19.5")

################################################################################
# Custom Go dependencies                                                       #
################################################################################

# add `go_repository` rules here to override Gazelle-generated files.

################################################################################
# go mod                                                                       #
################################################################################

gazelle_dependencies()

load("//:third_party/go/repositories.bzl", "go_repositories")

# gazelle:repository_macro third_party/go/repositories.bzl%go_repositories
go_repositories()

################################################################################
# Protobuf                                                                     #
################################################################################

http_archive(
    name = "rules_proto",
    sha256 = "dc3fb206a2cb3441b485eb1e423165b231235a1ea9b031b4433cf7bc1fa460dd",
    strip_prefix = "rules_proto-5.3.0-21.7",
    urls = [
        "https://github.com/bazelbuild/rules_proto/archive/refs/tags/5.3.0-21.7.tar.gz",
    ],
)

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

rules_proto_dependencies()

rules_proto_toolchains()

################################################################################
# Foreign C/C++ build system support                                           #
################################################################################

http_archive(
    name = "rules_foreign_cc",
    sha256 = "2a4d07cd64b0719b39a7c12218a3e507672b82a97b98c6a89d38565894cf7c51",
    strip_prefix = "rules_foreign_cc-0.9.0",
    url = "https://github.com/bazelbuild/rules_foreign_cc/archive/refs/tags/0.9.0.tar.gz",
)

load("@rules_foreign_cc//foreign_cc:repositories.bzl", "rules_foreign_cc_dependencies")

rules_foreign_cc_dependencies()

################################################################################
# Bazel Tools                                                                  #
################################################################################

http_archive(
    name = "com_github_bazelbuild_buildtools",
    sha256 = "e3bb0dc8b0274ea1aca75f1f8c0c835adbe589708ea89bf698069d0790701ea3",
    strip_prefix = "buildtools-5.1.0",
    urls = [
        "https://github.com/bazelbuild/buildtools/archive/5.1.0.tar.gz",
    ],
)

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()
