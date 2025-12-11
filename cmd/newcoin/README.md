
# Fiber Coin Creation Documentation

Newcoin is a tool for creating new fiber coins from a [fiber.toml](../../config/fiber.toml) config file.

**Two workflows are supported:**
1. **Quick Start (Recommended)**: Run fibercoins using the `skycoin` binary with environment variables - no compilation needed
2. **Standalone Binary**: Generate coin-specific Go source code and compile a dedicated binary

---

## Quick Start - Automated Initialization (Recommended)

**Time: ~10 seconds** | No compilation required | No manual configuration

### Prerequisites
* `git`
* `go` 1.18+

Clone the skycoin source:
```bash
mkdir -p $HOME/go/src/github.com/skycoin
cd $HOME/go/src/github.com/skycoin
git clone https://github.com/skycoin/skycoin
cd skycoin
```

### Five Simple Steps

**1. Generate Genesis Wallet**
```bash
go run . cli addressGen > genesis.json
```

This creates a deterministic wallet with:
- Genesis address (where all initial coins go)
- Blockchain public key (for block validation)
- Blockchain secret key (for block signing - keep this secure!)

**2. Create Configuration**
```bash
go run . newcoin config > mycoin.toml
```

Edit mycoin.toml to customize:
- `display_name`: "My Coin"
- `ticker`: "MYC"
- `genesis_coin_volume`: 100000000000000
- `port`: 6000 (network port)
- `web_interface_port`: 6420 (web UI port)
- `default_connections`: ["your.server.ip:6000"] (your node's IP)
- `initial_unlocked_count`: Number of distribution addresses to unlock initially

Leave these fields **BLANK** (they auto-populate):
- `genesis_address_str` (from genesis.json)
- `blockchain_pubkey_str` (from genesis.json)
- `genesis_signature_str` (created automatically)
- `distribution_addresses` (generated next)

**3. Generate Distribution Addresses**
```bash
FIBER_TOML=mycoin.toml go run . cli fiberAddressGen -n 10
```

This automatically:
- ✅ Generates 10 distribution addresses
- ✅ **Updates mycoin.toml with distribution_addresses**
- ✅ Creates `addresses.txt` (for reference)
- ✅ Creates `seeds.csv` (for wallet import - keep secure!)

**4. Initialize Blockchain**
```bash
GENESIS=genesis.json FIBER_TOML=mycoin.toml \
  go run . daemon \
    --block-publisher \
    --download-peerlist=false \
    --disable-default-peers \
    --data-dir=$HOME/.mycoin
```

What happens automatically:
- ✅ Loads genesis credentials from genesis.json
- ✅ Loads coin configuration from mycoin.toml
- ✅ Loads distribution addresses from mycoin.toml
- ✅ Creates genesis block and signature
- ✅ **Writes address/pubkey/signature back to mycoin.toml**
- ✅ Starts the blockchain

You'll see in the logs:
```
Loaded genesis credentials from GENESIS: genesis.json
Loaded fiber config from FIBER_TOML: mycoin.toml
Updated fiber.toml with genesis credentials (address, pubkey, signature)
```

**5. Distribute Genesis Coins**

In a separate terminal:
```bash
# Get genesis secret key
export SECRET_KEY=$(cat genesis.json | jq -r '.entries[0].secret_key')

# Set RPC address
export RPC_ADDR="http://127.0.0.1:6420"
export COIN="mycoin"

# Distribute coins to distribution addresses
go run . cli distributeGenesis $SECRET_KEY
```

This creates block #1 distributing coins from the genesis address to all distribution addresses equally.

You'll see:
```
INFO: Block 1 created
INFO: Distributed 100,000,000 coins to 10 addresses (10,000,000 each)
```

**After first run**, mycoin.toml is complete with all credentials and distribution addresses populated.

**Running Without GENESIS** (after initial setup):
```bash
# Extract secret key for block publisher
export SECRET_KEY=$(cat genesis.json | jq -r '.entries[0].secret_key')

# Run using only FIBER_TOML (genesis data now in mycoin.toml)
FIBER_TOML=mycoin.toml go run . daemon \
  --block-publisher \
  --blockchain-secret-key=$SECRET_KEY
```

**Peer Node** (different port/data directory):
```bash
FIBER_TOML=mycoin.toml go run . daemon \
  --port=5998 \
  --data-dir=$HOME/.mycoin-peer \
  --web-interface-port=6418
```

**Import Distribution Wallets:**

Use the seeds from `seeds.csv` to import distribution wallets into the web interface or via CLI:
```bash
# Import a distribution wallet
go run . cli walletCreate -s "seed phrase from seeds.csv" -l "Distribution 1"
```

---

## Traditional Workflow - Standalone Binary

If you need a standalone binary with a custom name compiled into the source code:

### Steps 1-3: Same as Quick Start

Follow Steps 1-3 from the Quick Start above to:
1. Create genesis.json
2. Configure mycoin.toml  
3. Generate distribution addresses with `fiberAddressGen`

You can run the daemon with GENESIS once to auto-populate mycoin.toml, or manually copy the values.

### Step 4: Generate Coin-Specific Source Code

```bash
go run . newcoin createcoin --coin mycoin --config-file mycoin.toml
```

This creates:
- `cmd/mycoin/mycoin.go` - Standalone executable
- `cmd/mycoin/commands/root.go` - CLI command definitions  
- `src/params/params.go` - Importable parameters with **distribution addresses compiled in**

**Important:** The distribution addresses from mycoin.toml are now compiled into `src/params/params.go`. Any future changes to distribution addresses require re-running `createcoin`.

### Step 5: Initialize & Distribute

```bash
# Extract secret key
export SECRET_KEY=$(cat genesis.json | jq -r '.entries[0].secret_key')

# Run block publisher
go run cmd/mycoin/mycoin.go daemon --block-publisher --blockchain-secret-key=$SECRET_KEY

# In another terminal, distribute genesis
export RPC_ADDR="http://127.0.0.1:6420"
export COIN="mycoin"
go run cmd/mycoin/mycoin.go cli distributeGenesis $SECRET_KEY
```

### Step 6: Compile Distributable Binary (Optional)

```bash
go build -o mycoin ./cmd/mycoin
./mycoin daemon --block-publisher --blockchain-secret-key=$SECRET_KEY
```

### Differences from Quick Start

**Quick Start (Runtime Configuration):**
- ✅ Distribution addresses loaded from fiber.toml at runtime
- ✅ Change addresses by editing fiber.toml (no recompilation)
- ✅ One binary can run multiple coins with different configs
- ✅ Faster iteration during development

**Traditional (Compiled Configuration):**
- Distribution addresses compiled into binary
- Changes require recompiling
- Custom binary name for branding
- Slightly smaller deployment (no need for fiber.toml)

---

## Environment Variables Reference

### `GENESIS` - Genesis Wallet Credentials

Points to a wallet JSON file (created by `cli addressGen`) containing:
- Genesis address
- Blockchain public key
- Blockchain secret key

**Usage:**
```bash
GENESIS=/path/to/genesis.json skycoin daemon --block-publisher
```

**Takes precedence over:** fiber.toml values for address/pubkey/seckey

**Security:** The secret key from genesis.json is used but **NEVER** written to fiber.toml

### `FIBER_TOML` - Coin Configuration

Points to a fiber.toml file containing coin parameters.

**Usage:**
```bash
FIBER_TOML=/path/to/mycoin.toml skycoin daemon
```

**What it does:**
- Loads all coin configuration (display name, ports, burn factors, etc.)
- Provides defaults for CLI flags (visible in `--help`)
- Gets auto-updated with genesis credentials on first run

### Configuration Precedence (Highest to Lowest)

1. **CLI flags** - `--genesis-address`, `--blockchain-secret-key`, etc.
2. **GENESIS env** - Credentials from genesis.json
3. **FIBER_TOML env** - Values from fiber.toml
4. **Template defaults** - Hardcoded in generated code

---

## Commands Reference

### `cli fiberAddressGen` - Generate Distribution Addresses

Generates addresses and seeds for distribution, with optional automatic fiber.toml update.

```bash
go run . cli fiberAddressGen [flags]
```

**Flags:**
```
-n, --num int              Number of addresses to generate (default 1)
-a, --addr-file string     Output file for addresses (default "addresses.txt")
-s, --seed-file string     Output file for seeds (default "seeds.csv")
    --overwrite            Overwrite existing files
```

**Environment Variables:**
- `FIBER_TOML` - If set, automatically updates the specified fiber.toml with distribution_addresses

**Example (Standalone):**
```bash
# Generate 100 addresses to files only
go run . cli fiberAddressGen -n 100 -a dist-addresses.txt -s dist-seeds.csv
```

**Example (Auto-update fiber.toml):**
```bash
# Generate 10 addresses AND update mycoin.toml
FIBER_TOML=mycoin.toml go run . cli fiberAddressGen -n 10

# Output:
# Generated 10 addresses
# Saved to: addresses.txt, seeds.csv
# ✓ Updated mycoin.toml with 10 distribution addresses
```

**Files Created:**
- `addresses.txt` - Plain text list of addresses (for reference)
- `seeds.csv` - CSV with address,seed pairs (for wallet import - **keep secure!**)

**Security:**
- The `seeds.csv` file contains wallet seeds - treat it like private keys
- Store offline and encrypt for production use
- Never commit to version control

### `cli distributeGenesis` - Distribute Genesis Coins

Distributes genesis block coins to configured distribution addresses.

```bash
go run . cli distributeGenesis [genesis_secret_key] [flags]
```

**Prerequisites:**
- Daemon must be running with `--block-publisher`
- Distribution addresses must be configured in fiber.toml
- Blockchain must be at block 0 (genesis only)

**Environment Variables:**
- `RPC_ADDR` - RPC address of running daemon (default: http://127.0.0.1:6420)
- `COIN` - Coin name for data directory resolution

**Example:**
```bash
export RPC_ADDR="http://127.0.0.1:6420"
export SECRET_KEY=$(cat genesis.json | jq -r '.entries[0].secret_key')

go run . cli distributeGenesis $SECRET_KEY
```

**What It Does:**
1. Queries `/api/v1/coinSupply` from running daemon to get distribution addresses
2. Creates transaction splitting genesis coins equally
3. Signs transaction with genesis secret key
4. Injects into blockchain (block publisher creates Block #1)

**Validation:**
- Checks that blockchain is at seq 0 (only genesis block exists)
- Verifies MaxCoinSupply divides evenly by number of addresses
- Ensures all addresses are valid

### `newcoin createcoin` - Generate Source Code

Generates coin-specific Go source files from fiber.toml templates.

```bash
go run . newcoin createcoin [flags]
```

**Flags:**
```
-c, --coin string                      name of the coin to create (default "skycoin")
-d, --template-dir string              template directory path (default "./template")
-e, --coin-template-file string        coin template file (importable) (default "coin.template")
-f, --command-template-file string     command template file (executable) (default "command.template")
-g, --coin-test-template-file string   coin test template file (default "coin_test.template")
-i, --params-template-file string      params template file (default "params.template")
-j, --config-dir string                config directory path (default "./")
-k, --config-file string               config file path (default "fiber.toml")
```

**Example:**
```bash
go run . newcoin createcoin --coin privateness --config-file privateness.toml
```

---

## Configuration

### Required fiber.toml Fields (Auto-Populated)

These fields are automatically populated by the automated workflow:

**By GENESIS environment variable (Step 4):**
- `genesis_address_str` - Address that receives genesis coins
- `blockchain_pubkey_str` - Public key for block validation
- `genesis_signature_str` - Signature of the genesis block

**By fiberAddressGen with FIBER_TOML (Step 3):**
- `distribution_addresses` - Array of addresses for coin distribution

### Customizable Fields (Edit Before First Run)

You should customize these in your fiber.toml:

**Basic Info:**
- `display_name` - Coin display name (e.g., "Privateness")
- `ticker` - Price ticker symbol (e.g., "PRIV")
- `coin_hours_display_name` - Name for coin hours
- `coin_hours_ticker` - Ticker for coin hours

**Genesis Block:**
- `genesis_coin_volume` - Total coins in genesis block
- `genesis_timestamp` - Genesis block timestamp (Unix time)

**Network:**
- `port` - Default network port (e.g., 6000)
- `web_interface_port` - Default web UI port (e.g., 6420)
- `default_connections` - Array of trusted peer addresses
- `peer_list_url` - URL for peer discovery

**Supply & Distribution:**
- `max_coin_supply` - Maximum total supply (must divide evenly by number of distribution addresses)
- `distribution_addresses` - Auto-populated by `fiberAddressGen` (or manually add)
- `initial_unlocked_count` - Number of distribution addresses unlocked initially
- `unlock_address_rate` - Addresses to unlock per time interval
- `unlock_time_interval` - Time between unlock events (in seconds)
- `user_burn_factor` - Coinhour burn factor for transactions

**Transaction Limits:**
- `unconfirmed_max_transaction_size` - Max size for unconfirmed txns
- `create_block_max_transaction_size` - Max size when creating blocks
- `max_block_transactions_size` - Max total size of txns in a block

See [fiber.toml](../../config/fiber.toml) for all available options and defaults.

---

## Security Best Practices

### Genesis Secret Key

The secret key in `genesis.json` is **extremely sensitive**:

⚠️ **Never:**
- Commit genesis.json to version control
- Share the secret key publicly
- Store it in plain text on production servers
- Write it to fiber.toml (the system prevents this automatically)

✅ **Do:**
- Keep genesis.json offline in secure storage
- Use environment variables to pass the secret key at runtime
- Use different keys for testing vs production
- Clear your terminal history after using the secret key

### Recommended Setup

**Development:**
```bash
# Use GENESIS for easy testing
GENESIS=genesis.json FIBER_TOML=mycoin.toml skycoin daemon --block-publisher
```

**Production:**
```bash
# Store secret key securely
echo "your-secret-key-here" > /secure/path/publisher-key.txt
chmod 600 /secure/path/publisher-key.txt

# Use fiber.toml + secret key flag
FIBER_TOML=mycoin.toml skycoin daemon \
  --block-publisher \
  --blockchain-secret-key=$(cat /secure/path/publisher-key.txt)

# Clear history
history -c
```

---

## Distribution & Deployment

### What to Distribute

**Public (safe to share):**
- `fiber.toml` (with genesis credentials populated)
- Compiled binary (if using standalone workflow)
- `peers.txt` file
- Blockchain data (after distribution)

**Private (keep secure):**
- `genesis.json` (contains blockchain secret key)
- Publisher node private key files
- Wallet seed phrases

### Production Network Infrastructure

A typical production deployment includes:

**Subdomains:**
- `node.mycoin.com` - Mobile wallet node (public API)
- `explorer.mycoin.com` - Blockchain explorer
- `downloads.mycoin.com/blockchain/peers.txt` - Peer list
- `version.mycoin.com/mycoin/version.txt` - Version check

**Running Production Nodes:**

**Mobile Wallet Node (Public):**
```bash
FIBER_TOML=mycoin.toml skycoin daemon \
  --enable-all-api-sets=true \
  --log-level=info \
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
  --log-level=debug \
  --port=6001 \
  --web-interface-port=6418 \
  --data-dir=$HOME/.mycoin-publisher
```

**Reverse Proxy (Caddy Example):**
```
node.mycoin.com {
    reverse_proxy 127.0.0.1:6419
}

explorer.mycoin.com {
    reverse_proxy 127.0.0.1:8003
}
```

---

## Blockchain Initialization & Distribution

### Initial Blockchain State

After starting the daemon for the first time:
- **Block 0 (Genesis Block)**: Contains all coins in a single output to the genesis address
- **Unspents**: 1 (the genesis output)

### Distributing Genesis Coins

The `distributeGenesis` command creates Block #1, distributing coins from the genesis address to all configured distribution addresses.

**Prerequisites:**
- Daemon running with `--block-publisher` flag
- Distribution addresses configured in fiber.toml (via `fiberAddressGen` or manual entry)

**Command:**
```bash
# Set environment variables
export RPC_ADDR="http://127.0.0.1:6420"
export COIN="mycoin"

# Get genesis secret key from genesis.json
export SECRET_KEY=$(cat genesis.json | jq -r '.entries[0].secret_key')

# Distribute genesis coins
go run . cli distributeGenesis $SECRET_KEY
```

**What Happens:**
1. Reads distribution addresses from the running daemon's configuration
2. Creates a transaction splitting genesis coins equally among all distribution addresses
3. Signs the transaction with the genesis secret key
4. Injects the transaction into the blockchain
5. Block publisher creates Block #1 containing the distribution transaction

**After Distribution:**
- **Block 1 (Distribution Block)**: Contains coins split equally to all distribution addresses
- **Unspents**: N (where N = number of distribution addresses)

**Example Output:**
```
Block 0: 100,000,000 coins → genesis address
Block 1: 100,000,000 coins → 10 addresses (10,000,000 each)
```

**Verification:**
```bash
# Check blockchain has 2 blocks
go run . cli status | jq '.status.blockchain.head.seq'
# Should show: 1

# Check distribution address balance
go run . cli addressBalance <distribution-address>
# Should show: 10,000,000.000000 (for 10 addresses splitting 100M)
```

### Using Distribution Wallets

The distribution wallets can be imported using the seeds from `seeds.csv`:

**Via Web Interface:**
1. Open http://127.0.0.1:6420
2. Go to "Wallets" → "Load Wallet"
3. Paste seed from seeds.csv
4. Label it (e.g., "Distribution 1")

**Via CLI:**
```bash
# Import from seed
go run . cli walletCreate -s "seed phrase from seeds.csv" -l "Distribution 1"

# Check balance
go run . cli walletBalance Distribution-1.wlt
```

### Security Notes

**Keep Secure:**
- `genesis.json` - Contains blockchain secret key
- `seeds.csv` - Contains distribution wallet seeds
- Both files should be stored offline and encrypted

**Can Share:**
- `mycoin.toml` - Public configuration (after genesis credentials populated)
- `addresses.txt` - Just addresses, no keys

---

## Troubleshooting

### Reset/Reinitialize Blockchain

If you need to start over:

```bash
# Stop all running instances
pkill -f "mycoin\|skycoin daemon"

# Remove blockchain data (keep wallets!)
rm -rf ~/.mycoin/data.db ~/.mycoin/history.db ~/.mycoin/*.log
rm -rf ~/.mycoin-peer/data.db ~/.mycoin-peer/history.db ~/.mycoin-peer/*.log

# Keep wallets directory intact
# ls ~/.mycoin/wallets/

# Clear genesis fields in fiber.toml
sed -i 's/^genesis_address_str.*/genesis_address_str = ""/' mycoin.toml
sed -i 's/^blockchain_pubkey_str.*/blockchain_pubkey_str = ""/' mycoin.toml  
sed -i 's/^genesis_signature_str.*/genesis_signature_str = ""/' mycoin.toml

# Re-run with GENESIS to repopulate
GENESIS=genesis.json FIBER_TOML=mycoin.toml skycoin daemon --block-publisher
```

### Common Issues

**"Failed to load GENESIS wallet"**
- Check that genesis.json path is correct
- Verify JSON is valid: `jq . genesis.json`
- Ensure file has proper permissions: `chmod 600 genesis.json`

**"Failed to load FIBER_TOML config"**
- Check that fiber.toml path is correct
- Validate TOML syntax
- Ensure all required fields are present

**Help menu doesn't show custom values**
- Make sure to set environment variables BEFORE running `--help`:
  ```bash
  GENESIS=genesis.json FIBER_TOML=mycoin.toml skycoin --help
  ```

**"Could not connect to any trusted peer"**
- This is normal for new coins with no network
- Ignore ERROR messages about peer connections during initial setup
- Set up your own peer network using `default_connections` in fiber.toml

**Genesis signature not in fiber.toml after first run**
- Check logs for "Updated fiber.toml with genesis credentials"
- Verify fiber.toml file permissions (must be writable)
- Ensure you're using FIBER_TOML env (so the path is tracked)

---

## Key Improvements Over Old Process

### Before (Manual Process)
1. ❌ Run `cli addressGen` → manually copy to terminal
2. ❌ Manually edit fiber.toml with genesis address and pubkey
3. ❌ **Manually generate distribution addresses** → copy into fiber.toml
4. ❌ Run node briefly, watch debug logs
5. ❌ **Parse genesis signature from log output** (painful!)
6. ❌ Manually paste signature into fiber.toml
7. ❌ Re-run newcoin to regenerate code with distribution addresses
8. ❌ **Manually create distribution transaction** or use broken distributeGenesis
9. ❌ Multiple restarts, many potential errors

**Time: 10+ minutes** | **Error-prone** | **Tedious** | **8+ manual steps**

### Now (Automated Process)
1. ✅ Run `cli addressGen > genesis.json`
2. ✅ Edit fiber.toml (one time)
3. ✅ Run `fiberAddressGen` → **auto-updates fiber.toml**
4. ✅ Run daemon with GENESIS + FIBER_TOML → **auto-populates genesis credentials**
5. ✅ Run `distributeGenesis` → **auto-reads addresses from daemon config**
6. ✅ Done!

**Time: ~30 seconds** | **Reliable** | **Simple** | **4 commands total**

### Benefits

- ⚡ **Fast**: 30 seconds instead of 10+ minutes
- 🎯 **Reliable**: Zero manual copy-paste steps
- 🔒 **Secure**: Secret key never written to fiber.toml
- 📦 **Portable**: Single binary can run multiple coins with different configs
- 🔄 **Flexible**: Easy to test different configurations
- ✅ **Auto-updating**: fiber.toml populates itself (genesis + distribution addresses)
- 🌐 **Runtime config**: Change distribution addresses without recompiling
- 💰 **Working distributeGenesis**: Automatically reads config from running daemon
- 👀 **Transparent**: See actual values in `--help` output

---

## Additional Resources

- [fiber.toml reference](../../config/fiber.toml) - Full configuration options
- [FIBERCOIN-CONFIG.md](../../FIBERCOIN-CONFIG.md) - FIBER_TOML environment variable documentation
- [GENESIS-ENV-IMPLEMENTATION.md](../../GENESIS-ENV-IMPLEMENTATION.md) - Technical implementation details
- [Skycoin source](https://github.com/skycoin/skycoin) - Main repository
- [Skycoin Wiki](https://github.com/skycoin/skycoin/wiki) - Development guides

---

## Getting Help

For issues or questions:
- Check [Troubleshooting](#troubleshooting) section above
- Review [GitHub Issues](https://github.com/skycoin/skycoin/issues)
- Ask in [Skycoin Telegram](https://t.me/skycoin)
