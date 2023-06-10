#!/usr/bin/env bash

set -e

# Pack a bundle of dependency files

cd -- "$( dirname -- "${BASH_SOURCE[0]}" )"/..

rm -f deps-bundle.tar.zst

tar -Izstd -cf deps-bundle.tar.zst \
  ./opt/{include,lib,lib64}

echo "[+] Created deps-bundle.tar.zst"

# Now you can commit this file to the deps branch of the repository
