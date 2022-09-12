# radiance

Highly experimental Solana playing ground.

## Available tooling

- [solrays](cmd/solrays), an RPC proxy that exports call latency of a Solana RPC node.

## Building

In order to build the available tooling, the following is required:
- Go 1.18+
- Run
  ```
  ./generate.sh
  ```

Building all the tools _should_ be as easy as:
```
:; go build -o _bin/ github.com/certusone/radiance/cmd/...
```

All binaries will be placed in `_bin/` folder.

Or if you're just looking for a single tool, say `solrays`:
```
:; go build -o _bin/ github.com/certusone/radiance/cmd/solrays
```

**NOTE:** Mind yourself, some of the tools here tools may depend on C code (and CGO), and other shenanigans,
so you may have to adapt accordingly.