# Release Process

Releases are automated via GitHub Actions. Pushing a tag matching `v*` triggers the release workflow (`.github/workflows/release.yml`), which creates a draft GitHub release and builds binaries for all supported platforms in parallel.

## Binaries

A single combined binary named `skycoin` (or `skycoin.exe` on Windows) is produced for each platform. This binary is compiled from `cmd/release/` and includes all of the following as subcommands:

- `daemon` — skycoin wallet node
- `cli` — command line interface
- `web` — thin client web wallet
- `explorer` — blockchain explorer
- `newcoin` — fibercoin creation tool
- `skyhw` — hardware wallet utilities (only on platforms with CGO + libusb support)

The `skyhw` subcommand requires CGO and libusb. On platforms where CGO is disabled (Windows 386, Windows ARM64), the `skyhw` subcommand is excluded at compile time via build tags. All other subcommands are always available.

On Windows amd64, a separate `skyhw.exe` standalone binary is also built and included in the release archive alongside the required DLLs (`libhidapi-0.dll`, `libusb-1.0.dll`).

## Build Method

All binaries are built using `go install` directly from the tagged module version:

```
go install github.com/skycoin/skycoin/cmd/release@<tag>
```

This approach embeds version information automatically via `runtime/debug.BuildInfo` without requiring ldflags. No repository checkout is needed for the build step — only specific support files (libusb build script, Windows installer template) are fetched individually when required.

## Platforms

### Linux

| Architecture | Archive | CGO | Linking | Hardware Wallet |
|---|---|---|---|---|
| amd64 | `skycoin-<tag>-linux-amd64.tar.gz` | Yes | Static (musl) | Yes |
| arm64 | `skycoin-<tag>-linux-arm64.tar.gz` | Yes | Static (musl) | Yes |
| arm (ARMv6) | `skycoin-<tag>-linux-arm.tar.gz` | Yes | Static (musl) | Yes |
| armhf (ARMv7) | `skycoin-<tag>-linux-armhf.tar.gz` | Yes | Static (musl) | Yes |
| 386 | `skycoin-<tag>-linux-386.tar.gz` | Yes | Static (musl) | Stub (gousb unsupported) |
| riscv64 | `skycoin-<tag>-linux-riscv64.tar.gz` | Yes | Static (musl) | Yes |

All Linux binaries are **statically linked** using musl cross-compilers. This eliminates runtime dependencies — the binary runs on any Linux distribution without requiring glibc or any shared libraries.

The static linking process:
1. Download musl cross-compiler toolchain for the target architecture
2. Cross-compile libusb-1.0 as a static library using the musl toolchain (`ci-scripts/build-libusb-musl.sh`)
3. Build the Go binary with `-linkmode external -extldflags '-static'` using the musl GCC as the C compiler

Musl toolchains are cached between builds to avoid repeated downloads.

### macOS

| Architecture | Archive | Installer | CGO | Hardware Wallet |
|---|---|---|---|---|
| amd64 (Intel) | `skycoin-<tag>-darwin-amd64.tar.gz` | `skycoin-<tag>-darwin-amd64.pkg` | Yes | Yes |
| arm64 (Apple Silicon) | `skycoin-<tag>-darwin-arm64.tar.gz` | `skycoin-<tag>-darwin-arm64.pkg` | Yes | Yes |

macOS binaries are dynamically linked (standard for macOS). CGO is enabled with libusb and hidapi installed via Homebrew. Each architecture is built on its native runner (`macos-15-intel` for amd64, `macos-latest` for arm64).

The `.pkg` installers are created with `pkgbuild` and install the binary to `/usr/local/bin/skycoin`.

### Windows

| Architecture | Archive | Installer | CGO | Hardware Wallet |
|---|---|---|---|---|
| amd64 | `skycoin-<tag>-windows-amd64.zip` | `skycoin-installer-<tag>-windows-amd64.msi` | Yes | Yes (+ separate `skyhw.exe`) |
| 386 | `skycoin-<tag>-windows-386.zip` | `skycoin-installer-<tag>-windows-386.msi` | No | No |
| arm64 | `skycoin-<tag>-windows-arm64.zip` | _(none)_ | No | No |

Windows amd64 is the only Windows architecture with CGO enabled. libusb and hidapi are installed via msys2/mingw. The 386 and arm64 builds use `CGO_ENABLED=0`, which excludes hardware wallet support but keeps all other functionality.

The amd64 archive includes a standalone `skyhw.exe` binary and the required DLLs (`libhidapi-0.dll`, `libusb-1.0.dll`) for hardware wallet operations.

MSI installers are built using WiX 3.11 for amd64 and 386 architectures. The installer adds `skycoin.exe` to the system PATH. No MSI installer is produced for arm64.

## Release Artifacts

Each release includes the following for every platform:

- Binary archive (`.tar.gz` for Linux/macOS, `.zip` for Windows)
- Per-archive SHA256 checksum (`.sha256`)
- Combined `checksums.txt` with all checksums

macOS releases additionally include `.pkg` installers with checksums.

Windows amd64 and 386 releases additionally include `.msi` installers with checksums.

## Triggering a Release

```bash
git tag v0.28.5
git push origin v0.28.5
```

The release is created as a **draft**. After verifying the artifacts, publish it manually from the GitHub releases page.

## Local Builds

To build the release binary locally:

```bash
# With hardware wallet support (requires libusb-1.0-dev)
CGO_ENABLED=1 go build -tags cgo -o skycoin ./cmd/release/

# Without hardware wallet support
CGO_ENABLED=0 go build -o skycoin ./cmd/release/
```

Or install from the module directly:

```bash
go install github.com/skycoin/skycoin/cmd/release@latest
```
