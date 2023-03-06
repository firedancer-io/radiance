FROM alpine:3.17 AS builder

# Add build dependencies
RUN apk add --no-cache \
    git go gcc g++ cmake ninja linux-headers \
    zlib-dev zlib-static bzip2-dev bzip2-static lz4-dev lz4-static zstd-dev snappy-dev snappy-static

# Drop privileges
RUN addgroup -g 1000 user \
 && adduser -G user -u 1000 user -D -h /home/user
USER user

# Fetch RocksDB
RUN cd /home/user \
 && git clone https://github.com/facebook/rocksdb --branch v7.10.2 --depth 1
WORKDIR /home/user/rocksdb
RUN mkdir -p build \
 && cd build \
 && cmake -G Ninja .. \
      -DFAIL_ON_WARNINGS=OFF \
      -DCMAKE_C_FLAGS="-march=x86-64-v2" \
      -DCMAKE_CXX_FLAGS="-march=x86-64-v2" \
      -DCMAKE_BUILD_TYPE=Release \
      -DROCKSDB_BUILD_SHARED=OFF \
      -DWITH_GFLAGS=OFF \
      -DWITH_BZ2=ON \
      -DWITH_SNAPPY=ON \
      -DWITH_ZLIB=ON \
      -DWITH_ZSTD=ON \
      -DWITH_ALL_TESTS=OFF \
      -DWITH_BENCHMARK_TOOLS=OFF \
      -DWITH_CORE_TOOLS=OFF \
      -DWITH_RUNTIME_DEBUG=OFF \
      -DWITH_TESTS=OFF \
      -DWITH_TOOLS=OFF \
      -DWITH_TRACE_TOOLS=OFF \
      -DWITH_URING=OFF \
 && ninja

# Fetch Go dependencies
RUN mkdir /home/user/radiance
WORKDIR /home/user/radiance
COPY go.mod go.sum ./
RUN go mod download

# Copy tree
COPY . ./

# Build
RUN mkdir build \
 && sh && CGO_CFLAGS="-I/home/user/rocksdb/include" \
    CGO_LDFLAGS="-L/home/user/rocksdb/build" \
    go build -o build/radiance -buildvcs=false -ldflags "-extldflags '-static -lbz2'" ./cmd/radiance

FROM alpine
# Copy CA certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# Copy build artifacts
COPY --from=builder /home/user/radiance/build/radiance /usr/bin/radiance
