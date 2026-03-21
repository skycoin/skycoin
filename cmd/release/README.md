# skycoin release

Release binary compilation including all skycoin utilities with hardware wallet support.

This binary includes:
- `daemon` - skycoin wallet daemon
- `cli` - skycoin command line interface
- `web` - skycoin thin client web wallet
- `explorer` - skycoin blockchain explorer
- `newcoin` - fibercoin creation tool
- `skyhw` - skycoin hardware wallet utilities (requires CGO and libusb)

## Build Requirements

Hardware wallet support requires CGO and libusb-1.0:

```bash
# Debian/Ubuntu
sudo apt-get install libusb-1.0-0-dev

# macOS
brew install libusb hidapi

# Fedora
sudo dnf install libusb1-devel
```

## Building

With hardware wallet support (CGO enabled):

```bash
CGO_ENABLED=1 go build -tags cgo -o skycoin ./cmd/release/
```

Without hardware wallet support (CGO disabled):

```bash
CGO_ENABLED=0 go build -o skycoin ./cmd/release/
```

When built without CGO, the `skyhw` subcommand is not available. All other subcommands work normally.

## Static Compilation (Release Builds)

Release binaries are statically compiled with musl cross-compilers to eliminate runtime dependencies. The release workflow uses `go install` to build directly from the tagged module version without cloning the repository.

## Architecture Support

Hardware wallet support (`skyhw`) is automatically excluded on:
- **386** — gousb/libusb does not support 32-bit
- **Windows ARM64** — gousb/libusb not available

This is handled via Go build tags (`//go:build cgo && !386 && !(windows && arm64)`). A stub command is provided on unsupported architectures.

## Fibercoin Support

Set `FIBER_TOML` to customize the binary for any fibercoin:

```bash
FIBER_TOML=path/to/fiber.toml ./skycoin
```

The ASCII art, descriptions, and default ports in all help menus will reflect the configured coin name.
