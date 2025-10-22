# Skycoin Web Wallet - Go Edition

A modernized, self-contained web wallet for Skycoin with embedded GUI and API proxy.

## Features

- Single binary with embedded web interface  
- Built with Go 1.25+ and Angular 12
- No external dependencies - everything embedded
- Secure proxy architecture - API requests proxied through server
- Cross-platform support (Linux, macOS, Windows)
- Simple CLI built with Cobra

## Architecture

The server acts as a proxy between the web UI and the Skycoin node:

```
Browser -> localhost:8001 -> Server (Go) -> node.skycoin.com
                              ↓
                          Embedded Web UI
```

All `/api/*` requests are proxied to the configured node, avoiding CORS issues and keeping the node URL server-side only.

## Quick Start

```bash
# Run with default settings (connects to https://node.skycoin.com)
./skycoin-web

# Specify a custom node URL
./skycoin-web --node-url https://your-node.example.com

# Custom host and port
./skycoin-web --host 0.0.0.0 --port 8080

# All together
./skycoin-web --node-url https://node.skycoin.com --host 0.0.0.0 --port 8080

# View help
./skycoin-web --help
```

## Building from Source

### Prerequisites
- Node.js 16+ and npm
- Go 1.25+

### Build Steps

```bash
# 1. Install dependencies
npm install --legacy-peer-deps

# 2. Build the web interface (outputs to src/gui/dist)
npm run build

# 3. Build the Go binary
go build -o skycoin-web .

# 4. Run it!
./skycoin-web
```

You can also run directly with:
```bash
# From repository root
go run .

# Or from anywhere once merged upstream
go run github.com/skycoin/skycoin-web@develop
```

## CLI Commands

- `skycoin-web` - Start the web server with defaults
- `skycoin-web serve` - Explicitly start the server
- `skycoin-web version` - Show version information
- `skycoin-web --help` - Show all available options

### Flags

- `--host` / `-H` - Host to bind to (default: 127.0.0.1)
- `--port` / `-p` - Port to serve on (default: 8001)
- `--node-url` / `-n` - Skycoin node URL to connect to (default: https://node.skycoin.com)

## Development

### Frontend Development
```bash
npm start  # Runs dev server on http://localhost:4200
```

### Backend Development  
```bash
go run . --port 8001
```

## Recent Modernization

This wallet has been upgraded from Angular 5 → Angular 12:

- Updated all dependencies and fixed 200+ vulnerabilities
- Migrated to modern Angular Material imports
- Added webpack 5 polyfills for crypto libraries
- Replaced node-sass with dart-sass
- Fixed all TypeScript compilation errors
- Created single-binary distribution with Go

## License

MIT

## Contributing

Contributions welcome! This is a fork maintained at github.com/0pcom/skycoin-web

Original project: github.com/skycoin/skycoin-web
