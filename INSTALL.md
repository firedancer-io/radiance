# Building Radiance

Radiance depends on many external dependencies including optional C libraries.
This makes the build procedure slightly more complicated than a pure Go project.

For installing C deps, the following options are available, each with their own tradeoffs:
- [Setup Script](#building-with-depssh) **(recommended)**
  - Stable, one-time ~3 minute installation
  - Supported on Debian, Ubuntu 20.04+, Fedora, RHEL 8+, Alpine 3, and macOS
- [Nixpkgs](#building-with-nix)
  - Somewhat stable, one-time 60 second installation
  - Supported on any platform with a pre-existing installation of Nix (Linux, macOS)
- [Without Cgo](#without-cgo)
  - Plain old pure Go build
  - Missing a lot of functionality (no Solana Labs compatibility layer)
- Manual Installation: If you know what you're doing, feel free to install C deps manually.
  - Good luck. Plan in 30 minutes of debugging time

Radiance requires further requires Go 1.20. Other Go versions are not supported.

Here's a trick to download another Go version in case you have the wrong one.
(See [Managing Go versions](https://golang.org/doc/manage-install))

    go install golang.org/dl/go1.20.5@latest
    "$(go env GOPATH)/bin/go1.20.5" download
    alias go="$(go env GOPATH)/bin/go1.20.5"

Once your Go toolchain and build dependencies are installed, you can build Radiance as usual:

    go mod download
    go run ./cmd/radiance

## Building with deps.sh

`deps.sh` fetches deps using Git, compiles them, and installs them into the `opt` dir of your Radiance checkout.

To complete the one-time installation of deps, run

    ./deps.sh

You may get prompted to install basic system packages like `pkg-config` or `cmake` if they are missing.

To add the installed deps to your shell env (and Go build tools), each time you open a new shell, run:

    source activate-opt

To clean up, simply run `./deps.sh nuke` or `rm -rf opt`

### How deps.sh works

Managing C deps is only as complicated as you want it to be.

`deps.sh` is a ~500 line long handwritten shell script. It just calls `git clone`, `make`, and `make install` or their
respective equivalents for each dependency. And yet it is more reliable than any tool that tries to be smart.

It will not pollute your system -- The C deps are not installed globally.

## Building with Nix

Radiance provides a [Nix](https://nixos.org/) package.

    nix-build
    ./result/bin/radiance --help

Note that the resulting binary is not freestanding.
It is linked against libraries provided by Nix and will not run on non-Nix hosts.

### Without Cgo

The full set of functionality requires C dependencies via Cgo.
To create a pure-Go build use the `lite` tag.

    go build -tags=lite ./cmd/radiance
