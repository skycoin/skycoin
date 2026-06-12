# Skycoin

![Skycoin logo](https://user-images.githubusercontent.com/26845312/32426705-d95cb988-c281-11e7-9463-a3fce8076a72.png){ width="320" }

Skycoin is a next-generation cryptocurrency and the settlement layer of the
Skycoin ecosystem. The node implemented in this repository runs the Skycoin
blockchain, exposes a REST API and CLI, and ships with desktop, web and
hardware-wallet clients.

[:material-rocket-launch: Install Skycoin](guides/installation.md){ .md-button .md-button--primary }
[:material-api: REST API](api/index.md){ .md-button }
[:material-console: CLI](cli/index.md){ .md-button }

## Documentation map

<div class="grid cards" markdown>

-   :material-book-open-variant: __Guides__

    Installation, exchange integration, development workflow, Docker
    images, custom fibercoins and the release process.

    [:octicons-arrow-right-24: Browse guides](guides/index.md)

-   :material-api: __REST API__

    The HTTP API exposed by the node — endpoints for the blockchain,
    wallets, transactions and network.

    [:octicons-arrow-right-24: API reference](api/index.md)

-   :material-console-line: __CLI__

    `skycoin-cli` — command-line wallet and node interaction.

    [:octicons-arrow-right-24: CLI reference](cli/index.md)

-   :material-package-variant: __Packages__

    In-tree library documentation: cipher primitives (bip32/bip39/
    base58/secp256k1) and the peer-exchange / networking layers.

    [:octicons-arrow-right-24: Package docs](packages/index.md)

-   :material-monitor-cellphone: __Clients__

    Desktop GUI, web client, lite client, the block explorer and the
    Electron build system.

    [:octicons-arrow-right-24: Clients](clients/index.md)

-   :material-graph-outline: __Dependency graph__

    The Go package dependency graph for the node.

    [:octicons-arrow-right-24: View graph](dependency-graph.md)

</div>

## Quick start

```bash
# Build and run the node from the repository root
go run . daemon

# Or with the desktop client configuration
go run . gui
```

See the [Installation guide](guides/installation.md) for prerequisites and
platform-specific instructions.

## Resources

- Source: [github.com/skycoin/skycoin](https://github.com/skycoin/skycoin)
- API reference: [REST API](api/index.md)
- CLI reference: [skycoin-cli](cli/index.md)
- Community: [Telegram](https://t.me/skycoin)
