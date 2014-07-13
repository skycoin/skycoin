
SkyWire: Skycoin Wire Protocol
==============================

Setup:
```
./setup.sh
```

```
./compile/instal-gvm.sh #install golang
./compile/clean-static-libs.sh #clear out cache
./compile/get-depedencies.sh #get golang package depedencies
./compile/install-to-gopath.sh #setup golang env
```

Skycoin Wire Protocol
=====================

/src/sync is the wire protocol. The wire protocol uses DHT to bootstrap peers. The wire protocol has the ability to syncronize data in a peer-to-peer fashion.

The wire protocol supports "blobs" and supports "hash chain" objects.Transactions and emergency alerts are blob messages. Obelisk consensus nodes and the main blockchain use hash chain replication.


Notes:
- DHT finds bootstrap peers via a hash, which is used to look up peers
- PEX (peer exchange) finds peers by asking peers for new peers

Todo:
- Need to be able to sync multiple blockchains. One sync instance per chain?
- Look up peers for blockchains by roothash of blockchain using DHT
- need multiple connection pools (one per blockchain)
- all instances/connection pools should share common listening port?

Blob Replication
================

A blob is a sequence of bytes. The id of a blob is the SHA256 hash of the blob. Blobs are syncronized using a gossip protocol. There is a callback function that determines if blob is valid data and what to do (ban peer who sent blob, replicate blob). Blobs are syncronized unordered.

When a client receives a blob, it uses the call back function to determine what to do with it. If the blob is valid, the application announces the hash of the blob to its peers. The peers then download the blob if they do not have it.

On connection, a peer transmits a hash list of the blobs it has available for peers to download. Peers download any blobs they do not have.

Hash Chain Replication
======================

A hash chain is a series of "blocks". Each block has a header and a body. The id of a block is a hash of the header of the block. Each block header contains the hash of the parent block. Each hash chain has a "root hash" which is the hash of the first block

Signed Hash Chain
=================

/src/hashchain is an example signed hash chain. A signed hash chain is a blockchain that only a person with a specific private key can mint blocks for. This is example for testing blockchain syncronization


