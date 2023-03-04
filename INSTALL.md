## Building with Go

Radiance commands can be built with standard Go tooling (Go >= 1.19 is required),
for example:

    go run go.firedancer.io/radiance/cmd/radiance

By default, commands that require extra dependencies (e.g. RocksDB) are disabled.

### With RocksDB

Radiance tools that require direct access to the blockstore (such as `blockstore` and `car`)
require a working C toolchain and extra compiler arguments to link against rocksdb.

You'll need a working C compiler toolchain as well as prerequisites listed
by [grocksdb](https://github.com/linxGnu/grocksdb#prerequisite) and
[RocksDB itself](https://github.com/linxGnu/grocksdb#prerequisite).

For RHEL/Centos/Fedora:

    dnf -y install "@Development Tools" gflags-devel snappy snappy-devel zlib zlib-devel bzip2 bzip2-devel lz4-devel libzstd-devel

First, check out the rocksdb submodule:

    git submodule update --init

Then, build the rocksdb C library:

    cd third_party/rocksdb
    make -j $(nproc) static_lib
    cd -

Finally, build the `radiance` command with the `rocksdb` tag:

    export CGO_CFLAGS="-I$(pwd)/third_party/rocksdb/include"
    export CGO_LDFLAGS="-L$(pwd)/third_party/rocksdb"

    go run -tags=rocksdb go.firedancer.io/radiance/cmd/radiance
