#!/usr/bin/env bash
set -euo pipefail

# Derived from Firedancer:
# https://github.com/firedancer-io/firedancer/blob/main/deps.sh

# Change into Radiance root directory
cd "$(dirname "${BASH_SOURCE[0]}")"
REPO_ROOT="$(pwd)"

# Fix pkg-config path and environment
# shellcheck source=./activate-opt
source activate-opt

# Load OS information
OS="$(uname -s)"
MAKE=( make -j )
case "$OS" in
  Darwin)
    ID=macos
    ;;
  Linux)
    # Load distro information
    if [[ -f /etc/os-release ]]; then
      source /etc/os-release
    fi
    ;;
  *)
    echo "[!] Unsupported OS $OS"
    ;;
esac

# Figure out how to escalate privileges
SUDO=""
if [[ ! "$(id -u)" -eq "0" ]]; then
  SUDO="sudo "
fi

# Install prefix
PREFIX="$(pwd)/opt"

help () {
cat <<EOF

  Usage: $0 [cmd] [args...]

  deps.sh is a lightweight dependency manager for Radiance.

  It locally installs
  - build tools using the system package manager
  - third party dependencies from Git repo sources

  If cmd is ommitted, default is 'install'.

  Commands are:

    help
    - Prints this message

    check
    - Runs system requirement checks for dep build/install
    - Exits with code 0 on success

    nuke
    - Get rid of dependency checkouts
    - Get rid of all third party dependency files
    - Same as 'rm -rf $(pwd)/opt'

    fetch
    - Fetches dependencies from Git repos into $(pwd)/opt/git

    install
    - Runs 'fetch'
    - Runs 'check'
    - Builds dependencies
    - Re-installs all project dependencies into prefix $(pwd)/opt

EOF
  exit 0
}

nuke () {
  rm -rf ./opt
  echo "[-] Nuked $(pwd)/opt"
  exit 0
}

checkout_repo () {
  # Skip if dir already exists
  if [[ -d ./opt/git/"$1" ]]; then
    echo "[~] Skipping $1 fetch as \"$(pwd)/opt/git/$1\" already exists"
  else
    echo "[+] Cloning $1 from $2"
    git -c advice.detachedHead=false clone "$2" "./opt/git/$1" --branch "$3" --depth=1
  fi
  echo

  echo "[~] Checking out $1 $3"
  (
    cd ./opt/git/"$1"
    git fetch origin "$3" --depth=1
    git -c advice.detachedHead=false checkout "$3"
  )
  echo
}

fetch () {
  mkdir -pv ./opt/git

  checkout_repo zlib    https://github.com/madler/zlib      "v1.2.13"
  checkout_repo zstd    https://github.com/facebook/zstd    "v1.5.4"
  checkout_repo rocksdb https://github.com/facebook/rocksdb "v7.10.2"
}

check_fedora_pkgs () {
  local REQUIRED_RPMS=( cmake pkgconf )

  echo "[~] Checking for required RPM packages"

  local MISSING_RPMS=( )
  for rpm in "${REQUIRED_RPMS[@]}"; do
    if ! rpm -q "$rpm" >/dev/null; then
      MISSING_RPMS+=( "$rpm" )
    fi
  done

  if [[ "${#MISSING_RPMS[@]}" -eq 0 ]]; then
    echo "[~] OK: RPM packages required for build are installed"
    return 0
  fi

  echo "[!] Found missing packages"
  echo "[?] This is fixed by the following command:"
  echo "        ${SUDO}dnf install -y ${MISSING_RPMS[*]}"
  read -r -p "[?] Install missing packages with superuser privileges? (y/N) " choice
  case "$choice" in
    y|Y)
      echo "[+] Installing missing RPMs"
      ${SUDO}dnf install -y "${MISSING_RPMS[@]}"
      echo "[+] Installed missing RPMs"
      ;;
    *)
      echo "[-] Skipping package install"
      ;;
  esac
}

check_debian_pkgs () {
  local REQUIRED_DEBS=( build-essential pkg-config )

  echo "[~] Checking for required DEB packages"

  local MISSING_DEBS=( )
  for deb in "${REQUIRED_DEBS[@]}"; do
    if ! dpkg -s "$deb" >/dev/null 2>/dev/null; then
      MISSING_DEBS+=( "$deb" )
    fi
  done

  if [[ ${#MISSING_DEBS[@]} -eq 0 ]]; then
    echo "[~] OK: DEB packages required for build are installed"
    return 0
  fi

  echo "[!] Found missing packages"
  echo "[?] This is fixed by the following command:"
  echo "        ${SUDO}apt-get install -y ${MISSING_DEBS[*]}"
  read -r -p "[?] Install missing packages with superuser privileges? (y/N) " choice
  case "$choice" in
    y|Y)
      echo "[+] Installing missing DEBs"
      ${SUDO}apt-get install -y "${MISSING_DEBS[@]}"
      echo "[+] Installed missing DEBs"
      ;;
    *)
      echo "[-] Skipping package install"
      ;;
  esac
}

check_alpine_pkgs () {
  local REQUIRED_APKS=( build-base pkgconf cmake )

  echo "[~] Checking for required APK packages"

  local MISSING_APKS=( )
  for deb in "${REQUIRED_APKS[@]}"; do
    if ! apk info -e "$deb" >/dev/null; then
      MISSING_APKS+=( "$deb" )
    fi
  done

  if [[ ${#MISSING_APKS[@]} -eq 0 ]]; then
    echo "[~] OK: APK packages required for build are installed"
    return 0
  fi

  echo "[!] Found missing packages"
  echo "[?] This is fixed by the following command:"
  echo "        ${SUDO}apk add ${MISSING_APKS[*]}"
  read -r -p "[?] Install missing packages with superuser privileges? (y/N) " choice
  case "$choice" in
    y|Y)
      echo "[+] Installing missing APKs"
      ${SUDO}apk add "${MISSING_APKS[@]}"
      echo "[+] Installed missing APKs"
      ;;
    *)
      echo "[-] Skipping package install"
      ;;
  esac
}

check_macos_pkgs () {
  local REQUIRED_FORMULAE=( pkg-config cmake )

  echo "[~] Checking for required brew formulae"

  local MISSING_FORMULAE=( )
  for formula in "${REQUIRED_FORMULAE[@]}"; do
    if [[ ! -d "/usr/local/Cellar/$formula" ]]; then
      MISSING_FORMULAE+=( "$formula" )
    fi
  done

  if [[ ${#MISSING_FORMULAE[@]} -eq 0 ]]; then
    echo "[~] OK: brew formulae required for build are installed"
    return 0
  fi

  echo "[!] Found missing formulae"
  echo "[?] This is fixed by the following command:"
  echo "        brew install ${MISSING_FORMULAE[*]}"
  read -r -p "[?] Install missing formulae with brew? (y/N) " choice
  case "$choice" in
    y|Y)
      echo "[+] Installing missing formulae"
      brew install "${MISSING_FORMULAE[@]}"
      echo "[+] Installed missing formulae"
      ;;
    *)
      echo "[-] Skipping formula install"
      ;;
  esac
}

check () {
  DISTRO="${ID_LIKE:-${ID:-}}"
  case "$DISTRO" in
    fedora)
      check_fedora_pkgs
      ;;
    debian)
      check_debian_pkgs
      ;;
    alpine)
      check_alpine_pkgs
      ;;
    macos)
      check_macos_pkgs
      ;;
    *)
      echo "Unsupported distro $DISTRO. Your mileage may vary."
      ;;
  esac
}

install_zlib () {
  cd ./opt/git/zlib

  echo "[+] Configuring zlib"
  ./configure \
    --prefix="$PREFIX"
  echo "[+] Configured zlib"

  echo "[+] Building zlib"
  "${MAKE[@]}" libz.a
  echo "[+] Successfully built zlib"

  echo "[+] Installing zlib to $PREFIX"
  make install -j
  echo "[+] Successfully installed zlib"
}

install_zstd () {
  cd ./opt/git/zstd/lib

  echo "[+] Installing zstd to $PREFIX"
  "${MAKE[@]}" DESTDIR="$PREFIX" PREFIX="" install-pc install-static install-includes
  echo "[+] Successfully installed zstd"
}

install_rocksdb () {
  cd ./opt/git/rocksdb
  mkdir -p build
  cd build
  cmake .. \
    -G"Unix Makefiles" \
    -DCMAKE_INSTALL_PREFIX:PATH="" \
    -DCMAKE_BUILD_TYPE=Release \
    -DROCKSDB_BUILD_SHARED=OFF \
    -DWITH_GFLAGS=OFF \
    -DWITH_LIBURING=OFF \
    -DWITH_BZ2=OFF \
    -DWITH_SNAPPY=OFF \
    -DWITH_ZLIB=ON \
    -DWITH_ZSTD=ON \
    -DWITH_ALL_TESTS=OFF \
    -DWITH_BENCHMARK_TOOLS=OFF \
    -DWITH_CORE_TOOLS=OFF \
    -DWITH_RUNTIME_DEBUG=OFF \
    -DWITH_TESTS=OFF \
    -DWITH_TOOLS=OFF \
    -DWITH_TRACE_TOOLS=OFF \
    -DZLIB_ROOT="$PREFIX" \
    -Dzstd_ROOT_DIR="$PREFIX"

  local NJOBS
  NJOBS=$(( $(nproc) / 2 ))
  NJOBS=$((NJOBS>0 ? NJOBS : 1))
  make -j $NJOBS
  make install DESTDIR="$PREFIX"
}

install () {
  echo "Cleaning install dir"
  for dir in bin include lib lib64 share usr; do
    rm -rf "$PREFIX/$dir"
  done

  # Compression algorithms. Prerequisites of RocksDB, but also used in
  # the Solana protocol.
  ( install_zlib    )
  ( install_zstd    )

  # RocksDB (imported by grocksdb)
  # See https://github.com/linxGnu/grocksdb#prerequisite
  ( install_rocksdb )

  echo "[~] Done! To wire up $(pwd)/opt with Go, run:"
  echo "    source activate-opt"
  echo
}

if [[ $# -eq 0 ]]; then
  echo "[~] This will fetch, build, and install Radiance dependencies into $(pwd)/opt"
  echo "[~] For help, run: $0 help"
  echo
  echo "[~] Running $0 install"

  read -r -p "[?] Continue? (y/N) " choice
  case "$choice" in
    y|Y)
      echo
      fetch
      check
      install
      ;;
    *)
      echo "[!] Stopping." >&2
      exit 1
  esac
fi

while [[ $# -gt 0 ]]; do
  case $1 in
    -h|--help|help)
      help
      ;;
    nuke)
      shift
      nuke
      ;;
    fetch)
      shift
      fetch
      ;;
    check)
      shift
      check
      ;;
    install)
      shift
      fetch
      check
      install
      ;;
    *)
      echo "Unknown command: $1" >&2
      exit 1
      ;;
  esac
done
