# Running Custom Fibercoins with Configuration Files

## Overview

You can now run custom fibercoin nodes without compiling coin-specific binaries by using the `--fiber-config` flag with a custom `fiber.toml` configuration file.

## Usage

### 1. Create a Custom Fiber Configuration

```bash
cp fiber.toml mycoin.toml
# Edit mycoin.toml to customize your fibercoin parameters
```

### 2. Run the Node with Custom Config

```bash
skycoin daemon --fiber-config mycoin.toml
```

The daemon will:
- Load all parameters from `mycoin.toml`
- Use the custom genesis settings, ticker, display name, etc.
- Set the data directory from the config file

### 3. Override Individual Parameters

CLI flags take precedence over config file values:

```bash
skycoin daemon \
  --fiber-config mycoin.toml \
  --port 7000 \
  --ticker MYC \
  --display-name "My Custom Coin"
```

## Configuration Precedence

1. **CLI flags** (highest priority)
2. **fiber.toml values** (from --fiber-config)
3. **Hardcoded defaults** (lowest priority)

## Available Flags

All fiber coin parameters can now be configured via flags:

### Core Parameters
- `--coin-name` - Name of the coin
- `--ticker` - Coin ticker symbol (e.g., SKY)
- `--display-name` - Display name of the coin
- `--port` - P2P network port
- `--web-interface-port` - Web interface port
- `--data-dir` - Data directory path

### Branding Parameters
- `--coin-hours-name` - Display name for coin hours
- `--coin-hours-name-singular` - Singular form
- `--coin-hours-ticker` - Ticker for coin hours
- `--qr-uri-prefix` - QR code URI prefix
- `--explorer-url` - Block explorer URL
- `--version-url` - Version check URL
- `--bip44-coin` - BIP44 coin type number

### Blockchain Parameters
- `--genesis-address` - Genesis address
- `--genesis-signature` - Genesis block signature
- `--genesis-timestamp` - Genesis timestamp
- `--blockchain-public-key` - Blockchain public key
- `--blockchain-secret-key` - Blockchain secret key (for block publishers)
- `--blockchain-pq-public-key` - Post-quantum (ML-DSA-44) blockchain public key
- `--blockchain-pq-secret-key` - Post-quantum blockchain secret key (for block publishers)
- `--pq-fork-seq` - First block that needs a post-quantum signature
- `--pq-address-fork-seq` - First block that may hold post-quantum addresses and their spends

### Post-quantum block signatures

With `blockchain_pq_pubkey_str` set, every block from `pq_fork_seq` must carry an
ML-DSA-44 (FIPS 204) signature by that key as well as the secp256k1 one. The
genesis block never needs one. `skycoin cli pqKeyGen` makes the key pair, and the
block publisher runs with `--blockchain-pq-secret-key`.

The signature is stored and sent apart from the block, so the block encoding does
not change. Nodes that understand it announce protocol version 3 and pass the
signatures to each other ahead of the blocks. Older nodes still verify only the
secp256k1 signature and keep following the chain, without the post-quantum
protection. A new fibercoin can set `pq_fork_seq = 1` to sign every block from the start.

### Post-quantum addresses

From `pq_address_fork_seq`, coins can be sent to post-quantum addresses. A
post-quantum address is version 1 and holds the first 20 bytes of SHA-256 over an
ML-DSA-44 public key, so it has the same size as a secp256k1 address. Spending
from one takes a transaction of type 1, which carries each input's ML-DSA-44 key
and signature as a witness beside the transaction. The witness is not part of the
transaction ID, and nodes on protocol version 4 pass it to each other ahead of
the transactions and blocks that need it.

This is a hard fork: older nodes reject type 1 transactions and stop at the first
block that holds one. `skycoin cli pqAddressGen` derives post-quantum addresses
from a seed, and `skycoin cli pqSend` spends from one. Reusing a post-quantum
address is safe, because its public key does not give away its secret key.

## Example: Creating a Test Coin

```toml
# testcoin.toml
[node]
genesis_signature_str = "..."
genesis_address_str = "..."
blockchain_pubkey_str = "..."
genesis_timestamp = 1426562704
genesis_coin_volume = 100000000000000
port = 7000
web_interface_port = 7420
data_directory = "$HOME/.testcoin"
display_name = "TestCoin"
ticker = "TST"
coin_hours_display_name = "Test Hours"
qr_uri_prefix = "testcoin"
explorer_url = "https://explorer.testcoin.com"
bip44_coin = 9000

[params]
max_coin_supply = 100000000
distribution_addresses = [
    "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
]
```

Run it:
```bash
skycoin daemon --fiber-config testcoin.toml
```

## Benefits

- **No compilation required** - Use the same `skycoin` binary for multiple coins
- **Easy configuration** - All parameters in one TOML file
- **Flexible deployment** - Override specific values via flags when needed
- **Consistent behavior** - Same code base, different configurations

## Implementation Details

The fiber config is loaded during the command's `PreRun` phase, after flag parsing but before the daemon starts. This ensures:

1. Config file values override hardcoded defaults
2. CLI flags override config file values
3. All parameters are properly validated before starting

See the full implementation in:
- `src/skycoin/config.go` - Config loading logic
- `cmd/skycoin/commands/root.go` - PreRun hook
- `template/coin.template` - Template for generated coins
