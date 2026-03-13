# skycoin-web

Thin client web wallet for Skycoin and fibercoins. Unlike the full desktop wallet (`skycoin`), `skycoin-web` does not run a node — it proxies API requests to one or more remote nodes while managing wallets locally.

## Build

```bash
go build ./cmd/skycoin-web/
```

## Usage

```bash
skycoin-web [flags]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--node-url` | `-n` | `https://node.skycoin.com` | Node URL (can be specified multiple times) |
| `--wallet-dir` | `-w` | (none) | Local wallet directory (can be specified multiple times) |
| `--port` | `-p` | `8001` | Port to serve on |
| `--host` | `-H` | `127.0.0.1` | Host to bind to |
| `--enable-seed-api` | | `false` | Enable the wallet seed API (requires `--wallet-dir`) |

### Examples

Connect to a single Skycoin node:

```bash
skycoin-web -n http://localhost:6420 -w ~/.skycoin/wallets -p 8001
```

Connect to multiple fibercoin nodes:

```bash
skycoin-web \
  -n http://localhost:6420 \
  -n http://localhost:8320 \
  -w ~/.skycoin/wallets \
  -w ~/.aix/wallets \
  -p 8006
```

When the number of `-w` flags matches the number of `-n` flags, each wallet directory is mapped 1:1 to the corresponding node. Otherwise, all wallet directories are shared across all nodes.

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
