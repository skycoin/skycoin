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
