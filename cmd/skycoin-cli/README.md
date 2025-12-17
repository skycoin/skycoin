# Skycoin CLI Documentation

Skycoin command line interface

The CLI is accessed via the `skycoin cli` subcommand. All CLI command APIs can be used directly from a Go application, see [Skycoin CLI Godoc](https://godoc.org/github.com/skycoin/skycoin/src/cli).

**Note:** Skycoin now uses **GNU/POSIX standard flags** (`-f` or `--flag`) instead of Go-style flags (`-flag`).

<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3" -->

- [Installation](#installation)
- [Environment Settings](#environment-settings)
	- [RPC_ADDR](#rpc_addr)
	- [RPC_USER](#rpc_user)
	- [RPC_PASS](#rpc_pass)
	- [COIN](#coin)
	- [DATA_DIR](#data_dir)
- [Usage](#usage)
- [Commands](#commands)

<!-- /MarkdownTOC -->

## Installation

The CLI is included in the main `skycoin` binary. See the main [README.md](../../README.md#installation) for installation instructions.

**Quick installation:**

```bash
# Run directly
$ go run github.com/skycoin/skycoin@develop cli --help

# Or install to PATH
$ go install github.com/skycoin/skycoin@develop
$ skycoin cli --help
```

## Environment Settings

The CLI uses environment variables for configuration.

### RPC_ADDR

CLI connects to the skycoin node REST API at `http://127.0.0.1:6420` by default.

Change the address by setting the `RPC_ADDR` environment variable:

```bash
$ export RPC_ADDR=http://127.0.0.1:6420
```

Note: `RPC_ADDR` must be in `scheme://host` format.

### RPC_USER

Username for authenticating requests to the skycoin node.

```bash
$ export RPC_USER=...
```

### RPC_PASS

Password for authenticating requests to the skycoin node.

```bash
$ export RPC_PASS=...
```

### COIN

Name of the coin (used for data directory).

```bash
$ export COIN=skycoin
```

Default: `skycoin`

### DATA_DIR

Directory where wallets and data are stored.

```bash
$ export DATA_DIR=$HOME/.skycoin
```

Default: `$HOME/.$COIN/`

## Usage

Run `skycoin cli` to see available commands:

```
$ skycoin cli --help
┌─┐┬┌─┬ ┬┌─┐┌─┐┬┌┐┌   ┌─┐┬  ┬
└─┐├┴┐└┬┘│  │ │││││───│  │  │
└─┘┴ ┴ ┴ └─┘└─┘┴┘└┘   └─┘┴─┘┴
skycoin command line interface

ENVIRONMENT VARIABLES:
  RPC_ADDR: Address of RPC node. Must be in scheme://host format. Default "http://127.0.0.1:6420"
  RPC_USER: Username for RPC API, if enabled in the RPC.
  RPC_PASS: Password for RPC API, if enabled in the RPC.
  COIN: Name of the coin. Default "skycoin"
  DATA_DIR: Directory where everything is stored. Default "$HOME/.$COIN/"

Usage:
  skycoin cli

Available Commands:
  addPrivateKey        Add a private key to wallet
  addressBalance       Check the balance of specific addresses
  addressGen           Generate skycoin or bitcoin addresses
  addressOutputs       Display outputs of specific addresses
  addressTransactions  Show detail for transaction associated with one or more specified addresses
  addresscount         Get the count of addresses with unspent outputs (coins)
  blocks               Lists the content of a single block or a range of blocks
  broadcastTransaction Broadcast a raw transaction to the network
  checkDBDecoding      Verify the database data encoding
  checkdb              Verify the database
  createRawTransaction Create a raw transaction that can be broadcast to the network later
  createRawTransactionV2 Create a raw transaction that can be broadcast to the network later
  decodeRawTransaction Decode raw transaction
  decryptWallet        Decrypt a wallet
  distributeGenesis    Distributes the genesis block coins into the configured distribution addresses
  encodeJsonTransaction Encode JSON transaction
  encryptWallet        Encrypt wallet
  fiberAddressGen      Generate addresses and seeds for a new fiber coin
  lastBlocks           Displays the content of the most recently N generated blocks
  listAddresses        Lists all addresses in a given wallet
  listWallets          Lists all wallets stored in the wallet directory
  pendingTransactions  Get all unconfirmed transactions
  richlist             Get skycoin richlist
  send                 Send skycoin from a wallet or an address to a recipient address
  showConfig           Show cli configuration
  showSeed             Show wallet seed and seed passphrase
  signTransaction      Sign an unsigned transaction with specific wallet
  status               Check the status of current Skycoin node
  transaction          Show detail info of specific transaction
  verifyAddress        Verify a skycoin address
  verifyTransaction    Verify if the specific transaction is spendable
  version              List the current version of Skycoin components
  walletAddAddresses   Generate additional addresses for a deterministic, bip44 or xpub wallet
  walletBalance        Check the balance of a wallet
  walletCreate         Create a new wallet
  walletCreateTemp     Create a new temporary wallet
  walletHistory        Display the transaction history of specific wallet. Requires skycoin node rpc.
  walletKeyExport      Export a specific key from an HD wallet
  walletOutputs        Display outputs of specific wallet
  walletScanAddresses  Scan addresses ahead for deterministic, bip44 or xpub wallet
```

## Commands

### General Pattern

All commands follow this pattern:

```bash
$ skycoin cli <command> [arguments] [flags]
```

### Get Help for Any Command

```bash
$ skycoin cli <command> --help
```

### Examples

**Check address balance:**
```bash
$ skycoin cli addressBalance 2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc
```
```json
{
    "confirmed": {
        "coins": "0.000000",
        "hours": "0"
    },
    "spendable": {
        "coins": "0.000000",
        "hours": "0"
    },
    "expected": {
        "coins": "0.000000",
        "hours": "0"
    },
    "addresses": [
        {
            "confirmed": {
                "coins": "0.000000",
                "hours": "0"
            },
            "spendable": {
                "coins": "0.000000",
                "hours": "0"
            },
            "expected": {
                "coins": "0.000000",
                "hours": "0"
            },
            "address": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc"
        }
    ]
}
```

**Generate addresses:**
```bash
$ skycoin cli addressGen --num 2
```
```json
{
    "meta": {
        "coin": "skycoin",
        "cryptoType": "scrypt-chacha20poly1305",
        "encrypted": "false",
        "filename": "2025_12_11_4eb0.wlt",
        "label": "default",
        "lastSeed": "ea460b8a0355075801d4c9f5621db30145080b60185f120d54cdfc1899345469",
        "seed": "pilot venture quick approve resemble false wear purchase child health annual able",
        "tm": "1765496128",
        "type": "deterministic",
        "version": "0.4"
    },
    "entries": [
        {
            "address": "w3Q4gXFdynmJ7vkzZ72qGjv4ZvXKgv9Q37",
            "public_key": "02c95e80c1c2a6869126841f04549e324edc9b3306c8e8f0ebbcd29ec8cb33fbb4",
            "secret_key": "46f4513e12620a951455e93c22e868ba908d1aa52efba7467bdb9607c9cace08"
        },
        {
            "address": "2TNhAaESivijY1tJRvdtj1r6j1eLMNd735a",
            "public_key": "03184673de208a804bd66f62357c85715d90c25a0532f2b2609e8483a6afcea86b",
            "secret_key": "142f5b7d13b5466a3df5ccdf5b03c26e560c5ff981dcd5f0acb78eeb3a52118e"
        }
    ]
}
```

**Create a wallet:**
```bash
$ skycoin cli walletCreate my-wallet
```
```json
{
    "meta": {
        "coin": "skycoin",
        "cryptoType": "scrypt-chacha20poly1305",
        "encrypted": "false",
        "filename": "my-wallet.wlt",
        "label": "my-wallet",
        "lastSeed": "",
        "seed": "your twelve word seed phrase will appear here for backup",
        "tm": "1765496128",
        "type": "deterministic",
        "version": "0.4"
    },
    "entries": [
        {
            "address": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc",
            "public_key": "02c95e80c1c2a6869126841f04549e324edc9b3306c8e8f0ebbcd29ec8cb33fbb4",
            "secret_key": "46f4513e12620a951455e93c22e868ba908d1aa52efba7467bdb9607c9cace08"
        }
    ]
}
```

**List wallets:**
```bash
$ skycoin cli listWallets
```
```json
{
    "directory": "/home/user/.skycoin/wallets",
    "wallets": [
        {
            "name": "my-wallet.wlt",
            "label": "my-wallet",
            "address_num": 1
        },
        {
            "name": "trading.wlt",
            "label": "trading",
            "address_num": 5
        }
    ]
}
```

**Check wallet balance:**
```bash
$ skycoin cli walletBalance my-wallet.wlt
```
```json
{
    "confirmed": {
        "coins": "125.500000",
        "hours": "5420"
    },
    "spendable": {
        "coins": "125.500000",
        "hours": "5420"
    },
    "expected": {
        "coins": "125.500000",
        "hours": "5420"
    },
    "addresses": [
        {
            "confirmed": {
                "coins": "125.500000",
                "hours": "5420"
            },
            "spendable": {
                "coins": "125.500000",
                "hours": "5420"
            },
            "expected": {
                "coins": "125.500000",
                "hours": "5420"
            },
            "address": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc"
        }
    ]
}
```

**Check node status:**
```bash
$ skycoin cli status
```
```json
{
    "status": {
        "blockchain": {
            "head": {
                "seq": 193949,
                "block_hash": "d06a1d4cc7b2b9ae0799b8f9f19746138486a6bebc3314ea76a9744866dd5811",
                "timestamp": 1765470105
            },
            "unspents": 183961,
            "unconfirmed": 10,
            "time_since_last_block": "7h13m47s"
        },
        "version": {
            "version": "0.27.1",
            "commit": "",
            "branch": ""
        },
        "coin": "skycoin",
        "open_connections": 22,
        "outgoing_connections": 3,
        "incoming_connections": 19,
        "uptime": "5h15m49s"
    }
}
```

**Send coins:**
```bash
$ skycoin cli send my-wallet.wlt 2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc 10.5
```
```json
{
    "txid": "89578005d8730fe1789288ee7dea036160a9bd43234fb673baa6abd91289a48b",
    "transaction": {
        "length": 220,
        "type": 0,
        "txid": "89578005d8730fe1789288ee7dea036160a9bd43234fb673baa6abd91289a48b",
        "inner_hash": "7b28bd34890f989c7637d19d555f419d4764ba1b3b985bd4e359bc0512ff1213",
        "fee": "48",
        "sigs": [...],
        "inputs": [...],
        "outputs": [...]
    }
}
```

**Create raw transaction:**
```bash
$ skycoin cli createRawTransactionV2 my-wallet.wlt 2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc 10.5
```
```json
{
    "rawtx": "dc00000000b8e5e8cfa97d9654044f8cd2da212dfe...",
    "transaction": {
        "length": 220,
        "type": 0,
        "txid": "89578005d8730fe1789288ee7dea036160a9bd43234fb673baa6abd91289a48b",
        "inner_hash": "7b28bd34890f989c7637d19d555f419d4764ba1b3b985bd4e359bc0512ff1213",
        "fee": "48",
        "sigs": [...],
        "inputs": [...],
        "outputs": [...]
    }
}
```

**Broadcast transaction:**
```bash
$ skycoin cli broadcastTransaction dc00000000b8e5e8cfa97d9654044f8cd2da212dfe...
```
```json
{
    "txid": "89578005d8730fe1789288ee7dea036160a9bd43234fb673baa6abd91289a48b"
}
```

## Detailed Command Documentation

For detailed documentation on each command including all flags and examples, use the `--help` flag:

```bash
$ skycoin cli addressGen --help
$ skycoin cli createRawTransactionV2 --help
$ skycoin cli walletCreate --help
```

### Common Flags

Most commands support these common patterns:

**Password for encrypted wallets:**
```bash
-p, --password string      Wallet password
```

**JSON output:**
```bash
-j, --json                 Returns the results in JSON format
```

**Number of items:**
```bash
-n, --num int              Number of addresses/items to generate
```

## Security Notes

**Wallet Passwords:**
Use caution when using the `-p` command. If you have command history enabled, your wallet encryption password can be recovered from the history log. If you do not include the `-p` option you will be prompted to enter your password after you enter your command.

**Sensitive Data:**
- Keep wallet files and seeds secure
- Never commit wallet files to version control
- Store seeds offline for production use
- Use encrypted wallets for added security

## Integration

The CLI can be used from:
1. **Command line**: Direct execution via `skycoin cli`
2. **Scripts**: Automate operations via shell scripts
3. **Go programs**: Import `github.com/skycoin/skycoin/src/cli` package

See [INTEGRATION.md](../../INTEGRATION.md) for integration guidelines.
