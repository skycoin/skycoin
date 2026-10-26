# Upgrading a fibercoin to post-quantum signatures

This guide moves a running fibercoin that was created without post-quantum
protection to ML-DSA-44 (FIPS 204) block signatures and post-quantum addresses.
A new fibercoin does not need it: set the post-quantum fields in its
`fiber.toml` before the genesis block, as [Custom Fibercoins](fibercoin-config.md)
describes.

The upgrade has two parts, each switched on at a block height of your choice:

| Setting | From that height | Nodes that have not upgraded |
|---|---|---|
| `pq_fork_seq` | every block also carries an ML-DSA-44 signature by the publisher | keep following the chain on the secp256k1 signature alone |
| `pq_address_fork_seq` | coins can be sent to post-quantum (version 1) addresses, and type 1 transactions spend them | accept payments to post-quantum addresses, and stop at the first block that spends one |

So the block signatures can be turned on without splitting the network, while
post-quantum addresses are a hard fork. The split does not happen at
`pq_address_fork_seq` itself. It happens at the first block that spends from a
post-quantum address, which can be any block after that height. Treat the
address fork height as the deadline for every node, exchange and explorer to
upgrade.

## 1. Make the publisher's post-quantum key pair

```bash
skycoin cli pqKeyGen > pqkeys.json
```

`secret_key` is a 32-byte seed in hex and `public_key` is the 1312-byte
ML-DSA-44 public key. Keep the secret key with the same care as the existing
`blockchain_seckey_str`. Anyone holding it can sign blocks that upgraded nodes
accept, as long as they also hold the secp256k1 key.

## 2. Pick the fork heights

Choose heights far enough ahead for everyone to install the new release. Keep
`pq_address_fork_seq` above `pq_fork_seq`, so the network runs with
post-quantum block signatures for a while before the hard fork.

Post-quantum inputs are large. Each one carries its public key and a 2420-byte
signature, about 3.7 KB in all, so the default 32 KB transaction and block size
limits fit only about eight post-quantum inputs per transaction. If you raise
`max_block_transactions_size`, `unconfirmed_max_transaction_size`,
`create_block_max_transaction_size` and `user_max_transaction_size`, do it in
the same release.

## 3. Ship the configuration

Add the public key and the heights to the coin's `fiber.toml`:

```toml
[node]
blockchain_pq_pubkey_str = "<public_key from pqkeys.json>"
pq_fork_seq = 50000
pq_address_fork_seq = 60000
```

Every node needs the same three values, from the release's `fiber.toml` or from
`--blockchain-pq-public-key`, `--pq-fork-seq` and `--pq-address-fork-seq`. A
node with a different public key rejects every block from `pq_fork_seq` with
`invalid ML-DSA-44 signature`.

## 4. Restart the publisher before `pq_fork_seq`

```bash
skycoin daemon --block-publisher \
  --blockchain-secret-key=<secp256k1 secret key> \
  --blockchain-pq-secret-key=<secret_key from pqkeys.json>
```

The publisher refuses to start if the post-quantum secret key is missing or does
not match `blockchain_pq_pubkey_str`, so a forgotten key shows up at startup,
not as a halted chain. `/api/v1/health` reports `pq_fork_seq`,
`pq_address_fork_seq` and the public key.

## 5. Between the two heights

From `pq_fork_seq`, upgraded nodes require the signature on every block. Nodes
on protocol version 3 or later pass the signatures to each other, and nodes that
have not upgraded keep following without checking them. Until
`pq_address_fork_seq`, upgraded nodes refuse payments to post-quantum addresses
with `post-quantum addresses cannot receive coins before the post-quantum address
fork`.

Use this time to get the remaining nodes, exchanges and explorers upgraded.

## 6. After `pq_address_fork_seq`

Coins can now be sent to post-quantum addresses:

```bash
skycoin cli walletCreate -t pq -l savings -n 5
skycoin cli pqAddressGen <seed>
```

A payment into a post-quantum address is an ordinary type 0 transaction, so
nodes that have not upgraded still accept it. The first spend from one is a type
1 transaction, and from that block on those nodes stop with `Transaction violates
hard constraint: transaction type invalid` and stay on the old chain.

Holders should move their coins to post-quantum addresses. An ordinary address
whose public key is already on the chain, because it has spent before, stays
exposed to a quantum attacker for as long as it holds coins. Reusing a
post-quantum address is safe.

## Checking the chain

`FIBER_TOML=<coin>.toml skycoin cli checkdb <data.db>` verifies the block
signatures of a stopped node's database, the post-quantum ones included, with
the keys and heights from the coin's config. `--pq-pubkey` and `--pq-fork-seq`
override them.

## Rolling back

Until the first type 1 transaction, nodes that have not upgraded hold the same
chain as upgraded ones, because the post-quantum signatures are kept apart from
the blocks. From that block on they cannot follow it, and the only way back is
to give up the blocks from that point.
