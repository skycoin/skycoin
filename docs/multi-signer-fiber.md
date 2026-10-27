# Multi-signer blocks for fiber chains

Status: design, not implemented.

A fibercoin today has one block publisher. A node accepts a block only when it
is signed by the single `blockchain_pubkey_str` in its `fiber.toml`
(`executeSignedBlock` in `src/visor/visor.go`). Every node checks every
transaction, so the publisher cannot forge or steal coins, but it alone decides
which transactions go in, in what order, and whether the chain moves at all.

This note describes an option for a fiber chain to have several independent
block signers instead, so that no single key can halt or censor it. It is the
first practical step toward the decentralised consensus that Obelisk aimed at.
It does not change the Skycoin chain, and a fiber chain that does not set the
new fields behaves exactly as today.

## Goals

- A fiber chain can name N signer keys at genesis. A block is final once at
  least k of them have signed it.
- The chain keeps going while up to N − k signers are offline.
- Two different blocks can never both be final at the same height while fewer
  than N/3 signers misbehave.
- Nodes that are not signers need no new trust. They check the signatures
  against the list in their config, as they check the one key today.
- Nothing changes for a chain that does not set the new fields.

Out of scope here: open membership, stake, and Obelisk's statistical voting.
Those can come later on top of this.

## Configuration

New fields in `fiber.toml`, node section:

| Field | Meaning |
|---|---|
| `block_signers` | List of signer public keys. Empty means today's single-publisher rule. |
| `block_signer_threshold` | k, the number of signatures that makes a block final. Defaults to ⌊2N/3⌋ + 1. Must be more than 2N/3. |
| `multi_signer_fork_seq` | First block that needs k signatures. 1 for a new chain. |
| `block_round_timeout` | How long signers wait for a proposal before moving to the next proposer. Defaults to three block intervals. |

A signer node also sets `blockchain_seckey_str` to its own key, which must be
in `block_signers`. Other nodes leave it empty as today. If the post-quantum
block rule is on, each signer also has an ML-DSA-44 key, listed in a matching
`block_pq_signers` field. That is left for a second step.

The genesis block is still signed by `blockchain_pubkey_str`, which may or may
not be one of the signers.

## Rounds and proposers

Block creation moves in rounds. For the block at height h, round r has the
proposer `block_signers[(h + r) mod N]`. Round 0 starts when block h − 1 is
final.

1. The proposer builds a block from its unconfirmed pool with the existing
   rules (`createBlock`), signs its header with its own key, and sends it to the
   other signers as a proposal.
2. A signer that gets a proposal for height h from the right proposer for the
   current round checks it with `VerifyBlock`, against its own chain at h − 1.
   It also checks that the block time is after the previous block and not more
   than a minute ahead of its own clock.
3. If the checks pass, the signer signs the header hash and broadcasts that
   signature to everyone.
4. Any node that holds the block and k valid signatures from distinct signers
   executes it and relays the block and the signatures.
5. If round r produces no final block within `block_round_timeout`, signers
   move to round r + 1, which has the next proposer.

Signers no longer need to build identical blocks, which was the fragile part of
the Obelisk paper. The proposer chooses the transactions, and rotation limits
how long any one signer can censor.

## Safety

The first stage uses one rule: a signer signs at most one block per height,
ever, and records that decision before it sends the signature. With k > 2N/3,
any two sets of k signers share more than N/3 members, so two conflicting
blocks at the same height cannot both reach k unless more than N/3 signers sign
both.

The cost is liveness. If a proposal reaches some signers but not enough of them
before a timeout, those signers cannot sign the next round's block, and the
height can stall. The second stage replaces the rule with Tendermint's
prevote, precommit and locking rules, which keep both safety and liveness with
fewer than N/3 faulty signers. The first stage is enough to prove the format,
the transport and the storage, and to run a test network.

A signer that signs two different blocks at the same height has produced
public evidence: two signatures by its key over different header hashes at one
height. Nodes log it and the API shows it. What to do about it, such as
removing the signer, is a governance decision for the chain.

## Format and transport

The block format does not change. `SignedBlock.Sig` stays the proposer's
signature, which nodes check against the signer list instead of the one key.
The other signatures travel and are stored apart from the block, the same way
the post-quantum block signature does (`PQSig`, `GivePQSigsMessage`,
`blockPQSigs`):

- `coin.BlockCosig` is a signer index and a signature over the header hash.
  `SignedBlock` gains `Cosigs []BlockCosig` tagged `enc:"-"`.
- A new bucket in `src/visor/blockdb` stores the cosignatures by header hash.
- `ProposeBlockMessage` carries a proposal to signers.
  `GiveBlockCosigsMessage` carries signatures. As with `GivePQSigsMessage`
  (`pqSigsProtocolVersion`), a node sends the new messages only to peers whose
  protocol version is at least the one that introduced them, so older nodes
  never receive one.
- A node that receives a block before its signatures holds it, as it does for a
  missing post-quantum signature, and asks for the signatures.
- `coin.MultiSignerRule`, in the style of `coin.PQBlockRule`, decides from the
  height and the config whether a block needs k signatures, and checks them.
  `executeSignedBlock` calls it in place of the single-key check.

The API's block endpoints add each block's signers. The CLI gains a command to
show the signer list and each signer's recent participation.

## Changing the signer set

In the first stage the set changes only by configuration: a new list takes
effect at a height that all nodes set, like a fork height. That is enough for a
test network and for a chain run by a known group.

A later stage adds a transaction type that changes the set from a given height
on. It is valid only when k current signers have signed it. That needs the
transaction type machinery that type 1 post-quantum transactions already use.

## Incentives

The Skycoin chain deliberately pays block makers nothing, and this design does
not change that. A fiber chain may want to reward its signers, and that can be
an option of each chain, off by default. Two sources are possible:

- A share of the coin hours a block burns, split evenly among the signers of
  that block. This needs outputs that hold coin hours and no coins, which are
  not allowed today (`"Zero coin output"` in `src/coin/transactions.go`), so it
  is a later hard fork for chains that want it.
- A fixed payout per block from the chain's own distribution addresses, split
  among the block's signers, which needs no change to the transaction rules.

Either way, a reward should apply only while the signer set has at least two
independent keys, and it should be split evenly, so that taking over the
signer set gains little. This note leaves the incentive for later. The
signatures stored with each block already record who took part.

## Testing

- Unit tests for `MultiSignerRule`, the cosignature encoding and the bucket.
- A local network of four signers and two plain nodes, with k = 3:
  - blocks are made with every proposer in turn;
  - with one signer stopped, blocks keep coming;
  - with two stopped, the chain stops, and it resumes when one comes back;
  - a signer that proposes two blocks at one height never gets both final;
  - a plain node that syncs from scratch accepts the chain and rejects a block
    with k − 1 signatures or a signature from a key outside the list.
- A chain with no `block_signers` passes the existing test suite unchanged.

## Stages

1. Config, `MultiSignerRule`, cosignature storage and messages, round-robin
   proposals with the one-block-per-height rule, the local test network.
2. Tendermint-style rounds with locking, for liveness.
3. Signer-set changes by transaction, post-quantum signer keys, explorer
   support.
4. Optional per-chain signer rewards.
