# Installing Go

Skycoin requires **Go 1.23 or later**.

## Official Installation Guide

Follow the official Go installation instructions for your operating system:

**https://go.dev/doc/install**

The official guide covers:
- Download and installation for Windows, macOS, and Linux
- Setting up your environment
- Verifying your installation
- Getting started with Go

## Quick Installation Links

- **Linux/macOS/Windows:** Download from https://go.dev/dl/
- **Package Managers:**
  - macOS (Homebrew): `brew install go`
  - Ubuntu/Debian: `sudo apt install golang-go` (check version meets minimum requirement)
  - Arch Linux: `sudo pacman -S go`
  - Fedora: `sudo dnf install golang`

## Verify Installation

After installing, verify your Go version meets the minimum requirement:

```sh
$ go version
go version go1.23.0 linux/amd64  # or higher
```

## Module Support

Skycoin uses Go modules, so you **do not need to set up GOPATH** or clone the repository to a specific directory. You can:

- Run directly: `go run github.com/skycoin/skycoin@develop`
- Install to your PATH: `go install github.com/skycoin/skycoin@develop`
- Clone anywhere: `git clone https://github.com/skycoin/skycoin && cd skycoin && go build`

## Troubleshooting

If you encounter issues:

1. Ensure Go version is 1.23 or later: `go version`
2. Check Go is in your PATH: `which go` (Unix) or `where go` (Windows)
3. Clear module cache if needed: `go clean -modcache`

For more help, see:
- [Go Installation Troubleshooting](https://go.dev/doc/install#troubleshooting)
- [Go FAQ](https://go.dev/doc/faq)
