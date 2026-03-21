# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## 0.28.5

### Added

- **BIP44 Chain Selection**
  - `/api/v1/wallet/newAddress` now accepts `chain` parameter (`external`/`change`) and `account` parameter for BIP44 wallets
  - CLI `walletAddAddresses` has `--chain` flag for generating on external or change chain
  - GUI wallet shows "New Address" button for both external and change address sections

- **Node Shutdown API**
  - `POST /api/v1/shutdown` endpoint for graceful node shutdown (NET_CTRL API set)
  - CLI `halt` command to shut down a running node via RPC

- **FIBER_TOML Help Menus**
  - All help menus (top-level, daemon, cli, web, explorer) dynamically update ASCII art and descriptions based on `FIBER_TOML` env var
  - Default RPC port in CLI help reflects the fibercoin's configured `web_interface_port`

- **Newcoin Templates**
  - `newcoin templates` subcommand to export embedded templates for customization
  - `--template-dir` flag in `createcoin` now functional (was dead code)

- **Fibercoin Integration Test**
  - CI test exercising full fibercoin lifecycle: genesis wallet, fiber.toml, distribution, block publishing
  - Tests BIP44 chain selection via API
  - Tests `distributeGenesis` end-to-end

- **Custom GUI Override**
  - `skycoin web` now has `--gui-dir` flag to serve custom GUI instead of embedded

### Changed

- **Release Workflow**
  - Replaced goreleaser with direct `go install` builds for proper version embedding
  - All platforms build in parallel with matrix strategy
  - No repository checkout needed for builds (only specific files fetched when needed)
  - macOS `.pkg` installers built inline via `pkgbuild`
  - Updated to Go 1.26

- **API Client**
  - Removed `DisallowUnknownFields()` from JSON decoders for forward compatibility (older clients can talk to newer nodes)

- **Peer Connection Logging**
  - "Already connected to this peer" no longer logged as a warning (silently skipped)

### Fixed

- **fiber.toml Empty String Overrides**
  - Setting `peer_list_url = ""`, `explorer_url = ""`, `default_connections = []`, etc. in fiber.toml now correctly overrides hardcoded defaults (previously ignored empty values)

- **Windows Release Build**
  - Added msys2 mingw64 to PATH for `pkg-config` to find `libusb-1.0.pc`
  - Fixed MSI installer expecting separate `skyhw.exe` (now included in combined binary)

- **cmd/release CGO Split**
  - Hardware wallet import behind `cgo` build tag so `cmd/release` builds with `CGO_ENABLED=0`

- **Flaky Windows Test**
  - Increased `wait()` timeout in `TestStopListen` for slow CI runners

### Security

- Bumped `fast-xml-parser` to 5.5.7 in `/explorer`
- Bumped `flatted` to 3.4.2 in `/src/gui/static` and `/src/skycoin-web`
- Removed deprecated `protractor` from `/explorer` (eliminated 10 vulnerabilities)
- Ran `npm audit fix` across all npm directories
- Updated golangci-lint to v2.11.3 for Go 1.26 compatibility
- Updated GitHub Actions to Node.js 24 compatible versions

## 0.28.1

### Added

- **Explorer & Web Wallet**
  - `skycoin explorer` subcommand - blockchain explorer with embedded compiled UI
  - `skycoin web` subcommand - thin client web wallet using WebAssembly (skycoin-lite)
  - Embedded GUI sources for both wallet and explorer (Go embed directives)
  - Explorer runs standalone on port 8001, connects to skycoin node on port 6420
  - Web wallet provides lightweight browser-based access without full node

- **Hardware Wallet Integration**
  - Hardware wallet daemon and CLI integration in `cmd/hardware-wallet/` (binary: `skyhw`)
  - Hardware wallet binary (`skyhw`) included in Linux amd64 releases
  - Comprehensive hardware wallet setup documentation for Linux, macOS, and Windows
  - CI testing for hardware wallet utilities with libusb dependency
  - Restrict hardware-wallet to linux/amd64 only (build tag)

- **Release Infrastructure**
  - GitHub Actions workflow for automated releases on tag push
  - GoReleaser configurations for multi-platform builds
    - Linux: amd64, 386, arm64, arm, armhf, riscv64
    - macOS: amd64, arm64
    - Windows: amd64, 386, arm64
  - Changelog generator script (`ci_scripts/changelog.sh`) using GitHub CLI
  - Token permissions to CI workflow for enhanced security

- **CLI/Command Improvements**
  - Version display and build info flags (`--bv`, `--info`) to main command
  - Cobra-based CLI with standard help format
  - Embed template files in newcoin binary

- **Newcoin/Fibercoin Enhancements**
  - `FIBER_TOML` environment variable for runtime fibercoin configuration
  - Skycoin daemon can now run a fibercoin node by specifying `FIBER_TOML=/path/to/fibercoin/fiber.toml`
  - Optimized newcoin and fibercoin creation process
  - Dynamic ASCII art generation using calvin.AsciiFont
  - Dynamic help menu ASCII art and coin name based on FIBER_TOML
  - Complete automated fibercoin blockchain initialization
  - Allow external module lookups in newcoin make target

- **Documentation**
  - Add DEVELOPMENT.md and DOCKER.md documentation
  - Add dependency graph visualization to docs
  - Update CLI documentation with actual command outputs and coinhour conversions
  - Update documentation to reflect modern skycoin CLI structure and Go 1.23+ requirement

- **API Endpoints**
  - `/api/v1/wallet/xpub` API endpoint to get xpub key of a bip44 wallet
  - `/api/v1/transactions/num` to get total transactions number
  - `POST /api/v1/wallet/createTemp` API to create a temporary wallet
  - `POST /api/v1/wallet/scan` API to scan wallet addresses ahead (default 20)
  - `GET /api/v2/transactions` API to get transactions with pagination
  - `qr_uri_prefix` field to `/api/v1/health` endpoint
  - `src_txid` to verbose transaction inputs (affects multiple APIs)
  - `calculated_hours` to `/api/v1/uxout` and `/api/v1/address_uxouts`

- **CLI Commands**
  - `walletCreateTemp` command to create a temporary wallet
  - `walletScanAddresses` command to scan wallet addresses ahead
  - `private-keys` param to `walletCreate`, `walletCreateTemp`, and `walletNewAddresses` commands

- **Frontend/GUI**
  - Format number input fields with commas (e.g., 1,111,111 instead of 1111111)
  - Update Chinese translations to adjust for new GUI changes
  - Update GUI dist files (multiple frontend rebuilds)
  - Embed bip44 account manager instead of using it as member variable
  - Remove all unused unit test files from GUI folder
  - General improvements for the GUI (see #2589)
  - Bump Electron from 15.2.0 to 22.3.25
  - Bump axios from 0.21.1 to 1.6.0
  - Multiple security dependency updates:
    - json5 from 1.0.1 to 1.0.2
    - express from 4.17.1 to 4.18.2
    - qs from 6.5.2 to 6.5.3
    - decode-uri-component from 0.2.0 to 0.2.2
    - loader-utils from 1.4.0 to 1.4.2
    - socket.io-parser from 4.0.4 to 4.0.5
    - moment from 2.29.1 to 2.29.4
    - eventsource from 1.1.0 to 1.1.1
    - nanoid from 3.1.25 to 3.3.4
    - follow-redirects from 1.14.2 to 1.15.0
    - async from 2.6.3 to 2.6.4
    - minimist from 1.2.5 to 1.2.6
    - karma from 6.3.4 to 6.3.16
    - url-parse from 1.5.3 to 1.5.10
    - engine.io from 4.1.1 to 4.1.2
    - log4js from 6.3.0 to 6.4.0
    - nth-check from 2.0.0 to 2.0.1
    - handlebars from 4.0.11 to 4.5.3

- **Flags**
  - `-max-last-blocks-count` flag to limit `/api/v1/last_blocks` results (default 256, fixes OOM panic)
  - `-max-incoming-connection` flag to control maximum allowed incoming connections

- **Testing**
  - Go 1.14+ version support
  - Optimize GitHub Actions CI caching

### Changed

- **CLI/Command Structure**
  - Complete migration to Cobra-based CLI framework
  - Use Cobra standard help format for all commands
  - Remove hardcoded skycoin references from flag descriptions
  - Use lowercase for ASCII art and wallet name in help menu
  - Fix CLI help flag redefinition panic

- **Versioning & Build**
  - Implement buildinfo-based versioning for `cmd/skycoin-wallet/` and `cmd/hardware-wallet/`
  - Update GoReleaser configurations to match skywire release structure and naming conventions
  - Update GoReleaser ldflags to use buildinfo package for version injection
  - Remove documentation files from release archives

- **Binary Names**
  - Hardware wallet binary renamed from `skyhw-daemon` to `skyhw`

- **Configuration**
  - Fix newcoin config-file flag default to `config/fiber.toml`

- **Frontend/GUI**
  - Embed GUI sources into binary using Go embed directives
  - Rebuild skycoin wallet GUI with fixes
  - Fix embedded GUI implementation
  - Revise embedded GUI structure
  - Make GUI work with Node.js 12 and latest hardware wallet daemon
  - Fix web GUI operation in Docker container
  - GUI get the uri prefix from the API
  - GUI update electron and electron builder
  - GUI General improvements for the front-end, see [#2589](https://github.com/skycoin/skycoin/pull/2589)
  - GUI exchange page has been disabled (Swaplab integration not working)
  - Fix Node.js 20 / OpenSSL 3.0 compatibility for UI tests
  - Regenerate package-lock.json after dependency updates
  - Fix npm dependency conflict: upgrade jasmine-core to 3.10.0
  - Update Spanish and Chinese translations
  - Make GUI work with latest version of hardware wallet daemon

- **API Changes**
  - `POST /api/v1/wallet/encrypt` encrypts wallets without 'cryptoType' field using default crypto type
  - Include change addresses for bip44 wallets in `/api/v1/wallet` endpoint
  - `private-keys` param added to `/api/v1/wallet/create`, `/api/v1/wallet/createTemp`, and `/api/v1/wallet/newAddress`

- **CLI Command Changes**
  - `listWallets` removed the `[directory]` argument
  - `encryptWallet`/`decryptWallet` show wallet details after operation
  - `encryptWallet`/`decryptWallet` only return non-sensitive data (no seeds, secrets, private keys)
  - `walletBalance`, `createRawTransactionV2`, `showSeed` use wallet filename instead of full path
  - `createRawTransactionV2` accepts wallet filename instead of full path
  - `walletKeyExport` -p flag replaced with --path, -p now shorthand for --password
  - `walletCreate` changes:
    - Removed `-wallet-file` flag (filename auto-generated)
    - Result consistent with `/api/v1/wallet/create` API response
    - Encrypt wallet by default with password prompt
    - Removed `coin/c` and `crypto-type/x` flags
    - Added `scan` flag to scan ahead addresses with transaction history

- **Code Structure**
  - Move package `src/wallet/crypto` to `src/cipher/crypto`
  - Run `make format` - fix import ordering and formatting
  - Apply formatting from newcoin template generation

- **CI/Build**  
  - Consolidate CI workflows into single 3-job workflow (Linux, macOS, Windows)
  - Remove all Travis CI references
  - Use golangci-lint-action for linting
  - Start Xvfb before UI tests
  - Disable golangci-lint on Windows CI
  - Skip UI e2e tests in CI (Protractor deprecated)
  - Skip test-386 on macOS (Darwin)

### Fixed

- **Explorer & GUI Fixes**
  - Fix wallet balance not showing in GUI - variable shadowing bug
  - Fix explorer opening issue in Electron
  - Fix explorer --node-addr flag being ignored
  - Fix missing unconfirmed transactions from `/explorer/address`
  - Fix CalculatedHours for `/explorer/address`
  - Fix TestDisableGUIAPI 404 message expectation

- **Emergency Wipe**: #2579 Add emergency wipe option for the Skywallet
- **Temporary Wallet Load**: #1109 Temporary wallet load feature
- Fix CLI error message format by removing SilenceErrors
- Fix Windows test failure in TestChooseSpendsMaximizeUxOuts
- Fix TestPoolBroadcastMessage flaky test on 386 architecture
- Fix test failures in droplet and gnet packages
- Fix unit test failures in api and gnet packages
- Fix CORS test failures
- Fix flaky TestPexAddPeers test
- Fix skycoin_test.go for cobra-based flag parsing
- Fix integration test flag parsing
- Fix transaction create test: genesis UTXOs must have null source txn
- Update integration test scripts to use GNU-style flags (--flag)
- Fix 32-bit test failure: ensure slice has minimum length
- Fix integration tests: set ldflags to correct package
- Fix CLI integration test binary compilation
- Fix CLI test: remove flag parsing that conflicts with Cobra
- Filter test flags from CLI test binary before Cobra parses args
- Strip GOCOVERDIR warning from CLI integration test output
- Fix GOCOVERDIR warning regex to preserve error message newlines
- Fix API integration test 404 message expectation
- Fix TestDisableGUIAPI 404 message expectation
- Fix newcoin test file generation
- Apply formatting to root.go after newcoin generation
- Fix coin.template formatting to match goimports output
- Fix Windows CI test failures
- Fix golangci-lint errors in pool_test.go
- Fix formatting and go vet warnings
- Fix all remaining linter issues
- Fix linter issues: add exported var comments and improve file permissions
- Add nolint comments for remaining gosec warnings

### Removed

- Remove endpoint `/api/v2/metrics` (prometheus dependency removed, no longer supported)
- Remove duplicate skycoin_test.go from commands directory
- Remove unused imports in command.template

## [0.27.0] - 2019-11-26

### Added

- Add `createRawTransactionV2` CLI command, which calls the API of `/wallet/transaction` to create the transaction and can create then unsigned transaction. Once the API's performance issue has been fixed, we will replace the `createRawTransaction` with it.
- Add `signTransaction` CLI command to sign transaction.
- Do windows electron builds by travis and abandon the appveyor
- Migrate `skycoin.net` to `skycoin.com`
- Migrate project path to `SkycoinProject/skycoin`
- Use transaction history when scanning wallet addresses, instead of the current address balance
- Document the daemon's CLI options
- Add the ability to save transaction notes
- Add CLI `encodeJsonTransaction` command to retrieve raw transaction given its JSON representation
- Add `package bip44`, implementing the bip44 spec https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
- Codesign daemon and standalone binaries
- Add a guided method for entering the seeds in the GUI
- Add new wallet type `collection` for wallets that are an arbitrary collection of private keys, rather than generated from a seed
- Add new wallet type `bip44` for hierarchical deterministic (HD) wallets obeying the bip44 protocol.
  The Skycoin bip44 `coin` number is `8000`.
  `bip44` wallets avoid address reuse, generating a new change address for each transaction.
  Affects APIs are `POST /api/v1/wallet`, `GET /api/v1/wallets`, `GET /api/v1/wallet`, `POST /api/v1/wallet/seed` and `POST /api/v1/wallet/recover`.
  Refer to the [API documentation](./src/api/README.md) for API changes.
  `bip44` wallets support bip39 "seed passphrases".
  More details are explained in https://github.com/skycoin/skycoin/wiki/Wallet-File-Formats-and-Types
- `cli walletCreate` support for `bip44` wallets added
- Add `bip44_coin` field to `GET /api/v1/health` `fiber` params
- Add the "bulk send" option to the GUI advanced form
- Add `cli walletKeyExport` command to export `xpub`, `xprv`, `pub` and `prv` key from a `bip44` wallet
- Add `xpub` type wallet, which can generate addresses without exposing private keys
- Add `block_publisher` flag to `/api/v1/health`
- Add `no_broadcast` option to `POST /api/v1/injectTransaction` to disable broadcast of the transaction. The transaction will only be added to the local pool.
- Add `cli distributeGenesis` command to split the genesis block into distribution addresses when setting up a new fiber coin
