# skycoin-web

`skycoin web` is a thin client for skycoin, fibercoins, and bitcoin.

Unlike the `skycoin daemon` desktop wallet, `skycoin-web` does not run a node — it proxies API requests to one or more full nodes while managing wallets locally.

## History

The original source code for skycoin-web is at: [github.com/skycoin/skycoin-web](https://github.com/skycoin/skycoin-web)

`skycoin web` was originally created as a 'web only' wallet using skycoin libraries in golang webassembly (WASM) for cryptographic operations in the browser. __This is still the default mode of operation when no flags are specified.__

When no flags are specified, `skycoin web` connects to the default node of the skycoin mobile wallet: [node.skycoin.com](https://node.skycoin.com).

**In this mode (without specifying --wallet-dir flag) the seed is generated in the browser and does not leave the browser. All operations of the wallet in this mode require seed re-entry.**

**Note: it is generally considered BAD PRACTICE to enter a wallet seed into the `skycoin web` wallet which you are not personally / locally hosting.**

Hence, the skycoin-web wallet was not used very much, and was unmaintained for many years. The skycoin-web wallet was at some point in the past hosted at `web.skycoin.com`

After considering the following:

* **sending transactions with the `skycoin web` wallet in the default mode of operation requires seed re-entry**
* **there is no fool-proof way to ensure beyond the shadow of a doubt that the `skycoin web` wallet being hosted by any third party (even as hosted on web.skycoin.com) is legitimate**

One may only conclude that the usefulness of the `skycoin web` wallet in the specific application where it was employed (e.g. being hosted on web.skycoin.com) was __dubious at best__ and __a liability at worst.__

## Evolution

In the latter part of 2025, the `skycoin web` wallet source code was migrated to the [github.com/skycoin/skycoin](https://github.com/skycoin/skycoin) repo, around the same time as many other utilities underwent similar migration and updates.

The `skycoin web` wallet was updated to use a WASM binary compiled from Go, and all necessary resources were embedded in Go libraries.

The `skycoin web` wallet was included in the default compilation of skycoin from the [github.com/skycoin/skycoin](https://github.com/skycoin/skycoin) repository root.

In 2026, the `skycoin web` wallet was extended to be able to utilize wallets on the disk by specifying the `--wallet-dir` flag, marking a significant increase in its usefulness as a thin client wallet for desktop.

Subsequent improvements were made which include:

* Multi-fibercoin support
  - BIP44 / BIP32 GUI support (also added to `skycoin daemon`)
* Bitcoin support (via Electrum server or Bitcoin Core full node)
  - Segwit wallet support
* Hardware wallet support (Skycoin hardware wallet via daemon on port 9510)

## Operation

The operation without any flags was described above in the [history](#history) section.

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--node-url` | `-n` | `https://node.skycoin.com` | Node URL. Can be specified multiple times for fibercoins. |
| `--wallet-dir` | `-w` | _(none)_ | Local wallet directory (e.g. `~/.skycoin/wallets`). Enables persistent wallets. |
| `--port` | `-p` | `8001` | Port to serve on. |
| `--host` | `-H` | `127.0.0.1` | Host to bind to. |
| `--enable-seed-api` | | `false` | Enable the wallet seed API (requires `--wallet-dir`). |
| `--btc-node-url` | | _(none)_ | Bitcoin Core RPC URL (e.g. `http://user:pass@127.0.0.1:8332`). |
| `--btc-electrum-url` | | _(none)_ | Electrum server URL (e.g. `ssl://electrum.blockstream.info:50002`). |

Bitcoin support is activated by providing either `--btc-node-url` or `--btc-electrum-url`. If neither is specified, Bitcoin is not available as a coin.

### Hardware Wallet

The Skycoin hardware wallet is supported when the hardware wallet daemon is running on `127.0.0.1:9510`. The frontend communicates directly with the daemon for PIN entry, address generation, and transaction signing. No additional flags are needed.

## Architecture

```
Browser  <-->  skycoin-web proxy  <-->  Skycoin/fibercoin nodes
                    |
              Local wallets
```

- **Coin discovery**: On startup, queries each node's `/api/v1/health` endpoint to discover coin name, ticker, and configuration.
- **Per-coin routing**: The frontend accesses each coin via `/coin/{index}/api/*` routes. The proxy forwards requests to the corresponding node.
- **Wallet management**: Wallet create/list/encrypt/decrypt operations are handled locally using `wallet.Service`. No wallet data is sent to the remote nodes.
- **Transaction signing**: Transactions are signed client-side in the browser via WebAssembly (skycoin-lite).
- **CSRF handling**: Read-only POST endpoints (`balance`, `transactions`, `outputs`) are converted to GET requests. Mutating POST requests (`injectTransaction`) are forwarded with a CSRF token fetched from the target node.
- **Price ticker**: Each coin can specify a `price_ticker_id` and `price_ticker_source` (coinpaprika or coingecko) in its fiber config. The frontend fetches prices directly from these APIs.

## Frontend

The Angular frontend source is at `src/skycoin-web/`. It is built and embedded into the Go binary via `src/skycoin-web/src/gui/`.

To rebuild the frontend:

```bash
cd src/skycoin-web
npx ng build --configuration production
```

To rebuild the WASM cryptographic module:

```bash
GOOS=js GOARCH=wasm go build -o src/skycoin-web/src/assets/scripts/skycoin-lite.wasm ./src/skycoin-lite/wasm/
```

Then rebuild the Go binary from the repository root with `go build .`
