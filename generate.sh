#!/usr/bin/env bash
set -e

(
  cd third_party/tools
  ./build.sh
)

third_party/tools/bin/buf generate

# cargo install serde-generate
if command -v serdegen &> /dev/null
then
  serdegen ./pkg/gossip/schema.yaml \
    --language=Go \
    --with-runtimes=Bincode \
    --serde-package-name=gossip \
    > ./pkg/gossip/schema.go
fi
