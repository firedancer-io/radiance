# Ledger History

Radiance proposes a neutral file format to capture Solana ledger history.

It also ships various reference tool to work with this file format.

## Ledger content

*Ledger data* is broadly defined as transaction and consensus data that can be trustlessly validated.

On Solana, ledger data is made up by only two classes of information.
* Proof-of-History[^1] parameters (contains block hashes)
* Transactions (contains user txs and consensus txs)

[^1]: PoH is a cryptographic delay function based on recursive SHA256 hashing

### Entries

The Solana protocol propagates ledger data in the form of _entries_.

On the wire, entries are additionally _shredded_ into network packets with erasure coding,
but this bears no relevance to the Solana ledger itself.

Entries have the following schema.

```python
class Entry:
  num_hashes: uint64
  prev_hash: Hash
  transactions: list[Tx]
```

Note that a _block_ on Solana is made of up multiple entries.
But on the ledger itself, the concept of blocks is only implied.

### Existing Formats

We find that the following representations of ledger data are widely used.

1. **Shreds** as UDP packets
   - Used in the peer-to-peer network
   - Hard to capture and archive for long-term storage
2. Archives of **blockstore** databases
   - Archives of the RocksDB database used in the Solana Labs validator implementation
   - Technically implementation-defined, forces use librocksdb (C++)
3. Google Cloud **Bigtable** integration for Solana RPC nodes
   - Closed source
   - Locked to one specific vendor
   - Lacks PoH data

We introduce a new format better suited for long-term archival and public distribution than the existing alternatives.

## CARv1 File Format

The **Content-addressable ARchive** is a streaming container format for blobs (files without a name).

[IPLD CARv1 Specification](https://ipld.io/specs/transport/car/carv1/)

### Content Addressing

_Why not .tar.zst, .7z, .rar, etc?_

Unlike with traditional archive formats, all blobs in CARs are content-addressed with a hash function.
Blobs are referred by CIDs (content identifiers) which unambiguously refer to the exact byte contents.

Leveraging the [IPLD Merkle-DAG](https://docs.ipfs.tech/concepts/merkle-dag/) construction,
blobs can recursively refer to other CIDs to build arbitrarily complex acyclic graphs of data.

Thus, if users know and trust a root CID (~35 bytes), they can safely retrieve blobs from any untrusted source.
Notably, users have the ability to verify if untrusted blobs match exactly what was requested.

### Determinism

Ledger CAR files are reproducible and deterministic.
Independent node operators would generate byte-by-byte identical CAR files for the same extent of ledger history,
regardless of where that data is sourced from.

### Header

The header of the ledger CAR file is set to the following.

```json
{
  "roots": ["bafkqaaa"],
  "version": 1
}
```

Rationale: The CAR file does not have a single root so we place the "empty" multihash instead,
as recommended by the [CARv1 spec](https://ipld.io/specs/transport/car/carv1/#number-of-roots).

This implies that any CARv1 file starts with the following byte content (hex).

```
19 a2 65 72 6f 6f 74 73
81 d8 2a 45 00 01 55 00
00 67 76 65 72 73 69 6f
6e 01
```

### IPLD data types

#### Transactions

Each Solana transaction is mapped to an IPLD block in native (bincode) serialization.

```
type Transaction bytes
```

#### Entries

```
type Entry struct {
  numHashes  Int
  hash       Hash
  txs        TransactionList
} representation tuple

type TransactionList [ &Transaction ]
```

#### Blocks

```
type Block struct {
  slot      Int
  entries   [ Link ]
  shredding [ Shredding ]
} representation tuple
```
