# Source file descriptors: `c<layer>-<domain>-<system>`

Every hand written Go file in this repository opens its package comment with
the package name, the file's path and a short code:

```go
// Package visor src/visor/visor.go c2-coin-visor
// Package cipher src/cipher/address.go c0-com-crypto
// Package main cmd/explorer/explorer.go c4-ui-explorer
```

The code is a taxonomy label from the SSS3 Go standard, the same one
[skywire](https://github.com/skycoin/skywire/blob/develop/docs/code-descriptors.md)
uses. The standard itself lives in
[skywire/docs/sss3](https://github.com/skycoin/skywire/tree/develop/docs/sss3).
Its governing rule is that idiomatic Go beats a transliterated house style.

## What the three parts mean

```
c<layer>-<domain>-<system>
   |        |        +-- the subsystem inside the domain
   |        +----------- one of four: com | coin | app | ui
   +-------------------- c0..c4, the dependency tier
```

**Layer.** `c0` is the most generic tier and a lower tier never imports a
higher one. A package's layer is the highest layer it imports, or that plus
one. Go rejects import cycles, so this part cannot silently rot.

**Domain.** Deliberately few. Detail goes into systems, and a new domain is
only made for something that cannot be a system.

| domain | what it holds |
|---|---|
| `com` | primitives every other domain uses: cryptography, utilities, test helpers |
| `coin` | the Skycoin node: chain data, storage, networking, wallets, API, CLI |
| `app` | products built on the chain: the DEX, the hardware wallet, the lite client |
| `ui` | user interfaces: the wallet UIs, the explorer, the DEX client UI |

**System.** The subsystem inside a domain. Most new packages join an
existing system.

## The systems

### `com`

| system | packages | layer |
|---|---|---|
| `com.crypto` | `src/cipher` and everything under it | c0 |
| `com.util` | `src/util/*` (`fee` is c1 because it reads `coin`) | c0, c1 |
| `com.test` | `src/testutil`, `src/testutil/assert`, `src/testutil/require` | c0, c1 |

### `coin`

| system | packages | layer |
|---|---|---|
| `coin.chain` | `coin`, `transaction`, `params`, `fiber`, `consensus`, `config` | c1 |
| `coin.wallet` | `wallet` and its kinds, `btc`, `eth`, `electrum` | c1 |
| `coin.visor` | `visor/blockdb`, `visor/historydb`, `visor/dbutil`, `kvstorage` (c1), `visor` (c2) | c1, c2 |
| `coin.daemon` | `daemon/gnet`, `daemon/pex`, `daemon/strand` (c1), `daemon` (c2) | c1, c2 |
| `coin.api` | `readable` (c2), `api` (c3) | c2, c3 |
| `coin.cli` | `src/cli` (c3), `cmd/skycoin-cli` (c4) | c3, c4 |
| `coin.node` | `src/skycoin` (c3), `cmd/skycoin` (c4) | c3, c4 |
| `coin.tool` | `templates` (c0), the root `main`, `cmd/skycoin-wallet`, `cmd/newcoin`, `cmd/release`, `cmd/address_gen`, `cmd/cipher-testdata`, `cmd/monitor-peers` (c4) | c0, c4 |

### `app`

| system | packages | layer |
|---|---|---|
| `app.dex` | `skydex/*` (c1), `skydex/server` (c2), `skydex/chain` (c3), `cmd/skydex-client`, `cmd/skydex-market` (c4) | c1 to c4 |
| `app.hw` | `hardware-wallet/*` (c1), `hardware-wallet/cli`, `hardware-wallet-daemon/*` (c2), `cmd/hardware-wallet` (c4) | c1, c2, c4 |
| `app.lite` | `skycoin-lite/wasmgz`, `wasm-go`, `wasm-tinygo`, `cmd/gen` (c1), the rest of `skycoin-lite` (c2) | c1, c2 |

### `ui`

| system | packages | layer |
|---|---|---|
| `ui.wallet` | `skycoin-web/src/gui` (c1), `ui2` (c3), `src/skycoin-web`, `cmd/skycoin-web` (c4) | c1, c3, c4 |
| `ui.explorer` | `explorer` (c0), `explorerui` (c3), `cmd/explorer` (c4) | c0, c3, c4 |
| `ui.dex` | `cmd/skydex-client/skydexui` | c3 |

## Where config goes

Config is never its own system. A config package takes the domain and system
of the thing it configures. `params` and `config` (the embedded `fiber.toml`)
configure the chain, so both are `coin.chain`.

## What this is not

It is not a rename. The layer, domain and system are labels to read and grep.
They never appear in an identifier, a package name or an import path.

Domain and system are not a lint rule. An `app` package importing `coin` is
normal. Only the layer matches something the compiler checks, and it checks
it already.

So `make check-descriptors` verifies only that

1. every non-test, non-generated Go file has a `// Package <name> <path> <descriptor>` line,
2. the path in that line is the file's own path, and
3. all files of one package carry the same descriptor.

There is no list of exceptions. Generated files (those with a
`// Code generated ... DO NOT EDIT.` line) are skipped.

## Adding a file

Copy the header from the file next to it and change the path. For a new
package, find its system in the tables above and pick the layer from what it
imports. The header goes directly above the `package` clause, or first in an
existing package comment followed by a `//` line:

```go
// Package wallet src/wallet/meta.go c1-coin-wallet
//
// Package wallet implements wallets and the wallet service.
package wallet
```

`cmd/skycoin/skycoin.go`, `cmd/skycoin/commands/root.go` and
`src/params/params.go` are written by `make newcoin` from `templates/`, so
their headers live in the templates.
