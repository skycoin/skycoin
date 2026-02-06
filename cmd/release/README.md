# skycoin release

Release binary compilation including all skycoin utilities with hardware wallet support.

This binary includes:
- `daemon` - skycoin wallet daemon
- `cli` - skycoin command line interface
- `web` - skycoin thin client web wallet
- `explorer` - skycoin blockchain explorer
- `newcoin` - fibercoin creation tool
- `skyhw` - skycoin hardware wallet utilities (daemon and cli)

## Build Requirements

Hardware wallet support requires libusb-1.0 for development builds:

```bash
# Debian/Ubuntu
sudo apt-get install libusb-1.0-0-dev

# macOS
brew install libusb

# Fedora
sudo dnf install libusb1-devel
```

## Static Compilation (Release Builds)

Release binaries are statically compiled with musl to eliminate runtime dependencies.
See the goreleaser configuration files for build details.

## Architecture Support

Hardware wallet support is disabled on 386 architecture due to limitations in the
github.com/google/gousb library.
