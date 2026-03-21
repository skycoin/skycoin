# Fiber Coin Creation Documentation

Create new fiber coins from a [fiber.toml](../../config/fiber.toml) config file.

**Two workflows:**
1. **Quick Start (Recommended)**: Run with `skycoin` binary using environment variables (no compilation)
2. **Standalone Binary**: Generate coin-specific source code and compile a dedicated binary

---

## Running Skycoin Commands

Throughout this documentation, `skycoin` refers to any of the following interchangeable methods:

```bash
# From within cloned source directory
go run .

# From outside source (after PR #2740 merge)
go run github.com/skycoin/skycoin@develop

# With release binary (>v0.28.0)
skycoin

# From skywire release binary (>v1.3.31)
skywire skycoin

# Using skywire from outside source
go run github.com/skycoin/skywire@develop skycoin
```

All commands produce the same help menu:

```
$ skycoin

    ┌─┐┬┌─┬ ┬┌─┐┌─┐┬┌┐┌
    └─┐├┴┐└┬┘│  │ │││││
    └─┘┴ ┴ ┴ └─┘└─┘┴┘└┘

 Available Commands:    
  cli       skycoin command line interface    
  daemon    skycoin wallet  
  explorer  skycoin blockchain explorer  
  newcoin   newcoin is a helper tool for creating new fiber coins  
  web       skycoin thin client web wallet

Flags:
  -b, --bv        print runtime/debug.BuildInfo.Main.Version
  -d, --info      print runtime/debug.BuildInfo
  -v, --version   version for skycoin
```

**For brevity, this documentation uses `skycoin` - substitute any method above.**

---

## Prerequisites

Depends on your chosen method:

| Method | Requires |
|--------|----------|
| Release binary | Nothing (pre-compiled) |
| `go run github.com/skycoin/skycoin@develop` | `go` 1.23+ only |
| `go run .` (from source) | `git` and `go` 1.23+ |

**Optional - Clone the source:**
```bash
mkdir -p $HOME/go/src/github.com/skycoin
cd $HOME/go/src/github.com/skycoin
git clone https://github.com/skycoin/skycoin
cd skycoin
```

---

## Quick Start - Runtime Configuration (Recommended)

**Time: ~30 seconds** | No compilation | Auto-configuration

### 1. Generate Genesis Wallet
```bash
skycoin cli addressGen > genesis.json
```

Creates deterministic wallet with genesis address, blockchain public/secret keys.

### 2. Create Configuration
```bash
skycoin newcoin config > mycoin.toml
```

**Edit these fields:**
- `display_name`: "MyCoin" *
- `ticker`: "MYC" *
- `genesis_coin_volume`: 100000000000000
- `port`: 6000 *
- `web_interface_port`: 6420 *
- `default_connections`: ["your.server.ip:6000"] *
- `initial_unlocked_count`: Number of distribution addresses to unlock initially

`*` Required changes

**Leave blank** (auto-populated):
- `genesis_address_str`, `blockchain_pubkey_str`, `genesis_signature_str`, `distribution_addresses`

### 3. Generate Distribution Addresses
```bash
FIBER_TOML=mycoin.toml skycoin cli fiberAddressGen -n 10
```

Generates addresses and **automatically updates mycoin.toml**. Creates `addresses.txt` and `seeds.csv` (keep secure!).

### 4. Initialize Blockchain
```bash
GENESIS=genesis.json FIBER_TOML=mycoin.toml skycoin daemon --block-publisher --download-peerlist=false --disable-default-peers
```

Automatically:
- Loads credentials from genesis.json
- Loads config from mycoin.toml
- Creates genesis block and signature
- **Writes genesis credentials back to mycoin.toml**

### 5. Distribute Genesis Coins

In a separate terminal:
```bash
RPC_ADDR="http://127.0.0.1:7420" COIN="mycoin" skycoin cli distributeGenesis $(jq -r '.entries[0].secret_key' genesis.json)
```

Creates block #1, distributing coins equally to all distribution addresses.

---

## Subsequent Runs

**After first run**, mycoin.toml is fully populated:

```bash
# Without GENESIS (using saved credentials)
FIBER_TOML=mycoin.toml skycoin daemon --block-publisher --blockchain-secret-key=$(jq -r '.entries[0].secret_key' genesis.json)

# Peer node (different port/directory)
FIBER_TOML=mycoin.toml skycoin daemon --port=7998 --data-dir=$HOME/.mycoin-peer --web-interface-port=7418
```

**Import distribution wallets** using seeds from `seeds.csv`:
```bash
skycoin cli walletCreate -s "seed phrase from seeds.csv" -l "Distribution 1"
```

---

## Traditional Workflow - Standalone Binary

For a standalone binary with compiled-in defaults:

### Steps 1-3: Same as Quick Start

Follow Quick Start steps 1-3 to create genesis.json, configure mycoin.toml, and generate distribution addresses.

### 4. Generate Source Code
```bash
skycoin newcoin createcoin --coin mycoin --config-file mycoin.toml
```

Creates:
- `cmd/mycoin/mycoin.go` - Standalone executable
- `cmd/mycoin/commands/root.go` - CLI commands
- `src/params/params.go` - Parameters with **compiled-in distribution addresses**

**Note:** Distribution address changes require re-running `createcoin`.

### 5. Build & Run
```bash
go build ./cmd/mycoin

# Initialize
go run cmd/mycoin/mycoin.go --block-publisher --blockchain-secret-key=$(jq -r '.entries[0].secret_key' genesis.json)

# Distribute genesis
RPC_ADDR="http://127.0.0.1:7420" COIN="mycoin" skycoin cli distributeGenesis $(jq -r '.entries[0].secret_key' genesis.json)
```

### Comparison

| Feature | Quick Start | Traditional |
|---------|-------------|-------------|
| Distribution addresses | Loaded at runtime from fiber.toml | Compiled into binary |
| Changes | Edit fiber.toml (no recompile) | Requires recompilation |
| Multiple coins | One binary with different configs | Separate binary per coin |
| Deployment | Requires fiber.toml | Binary only |

---

## Environment Variables

### `GENESIS`
Points to genesis wallet JSON (from `cli addressGen`). Takes precedence over fiber.toml for credentials. Secret key **never** written to fiber.toml.

```bash
GENESIS=/path/to/genesis.json skycoin daemon --block-publisher
```

### `FIBER_TOML`
Points to coin configuration file. Auto-updated with genesis credentials on first run.

```bash
FIBER_TOML=/path/to/mycoin.toml skycoin daemon
```

### Configuration Precedence
1. CLI flags (`--genesis-address`, `--blockchain-secret-key`)
2. `GENESIS` environment variable
3. `FIBER_TOML` environment variable
4. Template defaults

---

## Commands Reference

### `cli fiberAddressGen` - Generate Distribution Addresses

```bash
skycoin cli fiberAddressGen [flags]
```

**Flags:**
- `-n, --num int` - Number of addresses (default 1)
- `-a, --addr-file string` - Output file (default "addresses.txt")
- `-s, --seed-file string` - Seed file (default "seeds.csv")
- `--overwrite` - Overwrite existing files

With `FIBER_TOML` set, automatically updates fiber.toml with `distribution_addresses`.

**Security:** `seeds.csv` contains wallet seeds - store offline, encrypt for production, never commit to version control.

### `cli distributeGenesis` - Distribute Genesis Coins

```bash
skycoin cli distributeGenesis [genesis_secret_key] [flags]
```

**Prerequisites:**
- Daemon running with `--block-publisher`
- Distribution addresses configured
- Blockchain at block 0 (genesis only)

**Environment:**
- `RPC_ADDR` - RPC address (default: http://127.0.0.1:6420)
- `COIN` - Coin name for data directory

Queries daemon for distribution addresses, creates transaction splitting genesis coins equally, signs with genesis secret key, injects into blockchain.

### `newcoin createcoin` - Generate Source Code

```bash
skycoin newcoin createcoin --coin <name> --config-file <path> [flags]
```

Generates coin-specific Go source files from fiber.toml templates. Use `--template-dir` to specify a custom template directory instead of the embedded templates. See `skycoin newcoin createcoin --help` for all flags.

### `newcoin templates` - Export Templates

```bash
# Print all templates to stdout
skycoin newcoin templates

# Export templates to a directory for customization
skycoin newcoin templates ./my-templates
```

Exports the embedded newcoin templates to a directory so they can be customized. After editing, use them with `createcoin`:

```bash
skycoin newcoin createcoin --coin mycoin --template-dir ./my-templates
```

### `newcoin config` - Print Default Config

```bash
skycoin newcoin config > mycoin.toml
```

Prints the embedded default fiber.toml configuration to stdout.

---

## Configuration

### Key fiber.toml Fields

**Auto-populated:**
- `genesis_address_str`, `blockchain_pubkey_str`, `genesis_signature_str` (from GENESIS on first run)
- `distribution_addresses` (from `fiberAddressGen`)

**Customize before first run:**
- `display_name`, `ticker` - Coin branding
- `genesis_coin_volume`, `max_coin_supply` - Supply parameters
- `port`, `web_interface_port` - Network ports
- `default_connections` - Trusted peers
- `initial_unlocked_count`, `unlock_address_rate`, `unlock_time_interval` - Distribution schedule
- `user_burn_factor` - Coinhour burn rate

See [fiber.toml](../../config/fiber.toml) or `skycoin newcoin config` for all options.

---

## Security Best Practices

### Genesis Secret Key

**Extremely sensitive** - allows block publishing. Empty genesis wallet after distribution, but protect the key.

**Development:**
```bash
GENESIS=genesis.json FIBER_TOML=mycoin.toml skycoin daemon --block-publisher
```

**Production:**
```bash
echo "your-secret-key" > /secure/path/publisher-key.txt
chmod 600 /secure/path/publisher-key.txt

FIBER_TOML=mycoin.toml skycoin daemon \
  --block-publisher \
  --blockchain-secret-key=$(cat /secure/path/publisher-key.txt)

history -c
```

### What to Share vs. Keep Private

**Public:**
- `fiber.toml` (after genesis credentials populated)
- Compiled binary
- `peers.txt`, blockchain data
- `addresses.txt`

**Private:**
- `genesis.json` (blockchain secret key)
- `seeds.csv` (wallet seeds)
- Publisher node private keys

---

## Production Deployment

**Mobile Wallet Node:**
```bash
FIBER_TOML=mycoin.toml skycoin daemon \
  --enable-all-api-sets=true \
  --disable-csrf \
  --host-whitelist node.mycoin.com \
  --port=6000 \
  --web-interface-port=6419
```

**Block Publisher (Internal):**
```bash
FIBER_TOML=mycoin.toml skycoin daemon \
  --block-publisher=true \
  --blockchain-secret-key=$(cat /secure/publisher-key.txt) \
  --port=6001 \
  --web-interface-port=6418 \
  --data-dir=$HOME/.mycoin-publisher
```

**Reverse Proxy (Caddy):**
```
node.mycoin.com {
    reverse_proxy 127.0.0.1:6419
}

explorer.mycoin.com {
    reverse_proxy 127.0.0.1:8003
}
```

**Typical Infrastructure:**
- `node.mycoin.com` - Public API for mobile wallets
- `explorer.mycoin.com` - Blockchain explorer
- `downloads.mycoin.com/blockchain/peers.txt` - Peer list
- `version.mycoin.com/mycoin/version.txt` - Version check

---

## Troubleshooting

### Reset Blockchain

```bash
# Stop instances
pkill -f "mycoin\|skycoin daemon"

# Remove blockchain data (keeps wallets)
rm -rf ~/.mycoin/data.db ~/.mycoin/history.db ~/.mycoin/*.log

# Clear genesis fields in fiber.toml
sed -i 's/^genesis_address_str.*/genesis_address_str = ""/' mycoin.toml
sed -i 's/^blockchain_pubkey_str.*/blockchain_pubkey_str = ""/' mycoin.toml  
sed -i 's/^genesis_signature_str.*/genesis_signature_str = ""/' mycoin.toml

# Re-run
GENESIS=genesis.json FIBER_TOML=mycoin.toml skycoin daemon --block-publisher
```

### Common Issues

**"Failed to load GENESIS wallet"**
- Verify path and JSON validity: `jq . genesis.json`
- Check permissions: `chmod 600 genesis.json`

**"Failed to load FIBER_TOML config"**
- Verify path and TOML syntax
- Ensure required fields are present

**Help menu doesn't show custom values**
- Set environment variables before `--help`:
  ```bash
  GENESIS=genesis.json FIBER_TOML=mycoin.toml skycoin --help
  ```

**"Could not connect to any trusted peer"**
- Normal for new coins with no network
- Set up peer network using `default_connections` in fiber.toml

**Genesis signature not in fiber.toml after first run**
- Check logs for "Updated fiber.toml with genesis credentials"
- Verify fiber.toml is writable
- Ensure FIBER_TOML env is set

---

## Key Improvements Over Manual Process

### Before
❌ 8+ manual copy-paste steps  
❌ Parse genesis signature from logs  
❌ Multiple restarts and re-runs  
❌ **10+ minutes, error-prone**

### Now
✅ 4 commands total  
✅ Auto-populates fiber.toml  
✅ Zero manual copy-paste  
✅ **~30 seconds, reliable**

**Benefits:**
- ⚡ Fast: 30 seconds vs. 10+ minutes
- 🎯 Reliable: No manual copy-paste
- 🔒 Secure: Secret key never written to fiber.toml
- 📦 Portable: Single binary, multiple coins
- 🔄 Flexible: Change config without recompiling
- 💰 Working distributeGenesis: Auto-reads daemon config

---

## Getting Help

- [Troubleshooting](#troubleshooting) section above
- [GitHub Issues](https://github.com/skycoin/skycoin/issues)
- Telegram: [@skycoin](https://t.me/skycoin) or [@skywire](https://t.me/skywire)
