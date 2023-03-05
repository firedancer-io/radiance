## Building with Nix

The recommended and most stable way to build Radiance is with [Nix](https://nixos.org/).
This requires an existing Nix installation on your machine.

    nix-build
    ./result/bin/radiance --help

However, note that the resulting binary will only run under the same Nix environment.

## Building with Go

Radiance commands can be built with standard Go 1.19 tooling.

    go run ./cmd/radiance

### Without Cgo

The full set of functionality requires C dependencies via Cgo.
To create a pure-Go build use the `lite` tag.

    go build -tags=lite ./cmd/radiance

### With Cgo dependencies

Radiance tools that require direct access to the blockstore (such as `blockstore` and `car`)
require a working C toolchain and extra compiler arguments to link against rocksdb.

You'll need a working C compiler toolchain as well as prerequisites listed
by [grocksdb](https://github.com/linxGnu/grocksdb#prerequisite) and
[RocksDB itself](https://github.com/linxGnu/grocksdb#prerequisite).

**RHEL/Centos/Fedora**

    dnf -y install "@Development Tools" cmake zlib zlib-devel bzip2 bzip2-devel lz4-devel libzstd-devel

**Debian**

    # With package manager RocksDB
    apt install -y librocksdb-dev

    # With RocksDB from source
    apt install -y build-essential cmake zlib1g-dev libbz2-dev liblz4-dev libzstd-dev

**Building RocksDB**

To build RocksDB from source, run the following commands:

    git clone https://github.com/facebook/rocksdb --branch v7.10.2 --depth 1
    cd rocksdb
    mkdir -p build && cd build
    cmake .. \
      -DCMAKE_BUILD_TYPE=Release \
      -DROCKSDB_BUILD_SHARED=OFF \
      -DWITH_GFLAGS=OFF \
      -DWITH_BZ2=ON \
      -DWITH_SNAPPY=OFF \
      -DWITH_ZLIB=ON \
      -DWITH_ZSTD=ON \
      -DWITH_ALL_TESTS=OFF \
      -DWITH_BENCHMARK_TOOLS=OFF \
      -DWITH_CORE_TOOLS=OFF \
      -DWITH_RUNTIME_DEBUG=OFF \
      -DWITH_TESTS=OFF \
      -DWITH_TOOLS=OFF \
      -DWITH_TRACE_TOOLS=OFF
    make -j
    cd ../..

Finally, rebuild RocksDB with the appropriate Cgo flags.

    export CGO_CFLAGS="-I$(pwd)/rocksdb/include"
    export CGO_LDFLAGS="-L$(pwd)/rocksdb/build"

    go run ./cmd/radiance
