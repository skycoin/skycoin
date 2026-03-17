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

The `skycoin web` wallet was updated to use a WASM binary compiled with `tinygo` by default, and all necessary resources were embedded in go libraries.

The `skycoin web` wallet was included in the default compilation of skycoin from the [github.com/skycoin/skycoin](https://github.com/skycoin/skycoin) repository root.

In 2026, the `skycoin web` wallet was extended to be able to utilize wallets on the disk by specifying the `--wallet-dir` flag, marking a significant increase in it's usefulness as a thin client wallet for desktop.

Subsequent improvements were made which include:

* multi-fibercoin support
  -- BIP44 / BIP32 GUI support (also added to `skycoin daemon`)
* Bitcoin support (via electrum server or bitcoin full node)
  -- Segwit wallet support

## Operation

The operation without any flags was described above in the [history](#History) section.

* The `--node-url, -n` flag will set a new default node, or multiple nodes. The flag can be specified multiple times to set nodes for fibercoins.

* The `--wallet-dir, -w` flag enables wallets to be used from (or created / saved in) that wallet directory


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

Then rebuild the Go binary.
