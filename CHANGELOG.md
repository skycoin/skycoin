# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.20.0] - 2017-10-10

### Added

- New wallet frontend in angular4. This is a complete rewrite and fixes many
  of the old wallet issues.
- New wallet has preliminary support for OTC functionality
- Create `webrpc.Client` for JSON-2.0 RPC calls.
- Add this CHANGELOG.md file.
- Add Installation.md file, with install instructions for go.
- Timelock distribution addresses. The first 25% of the distribution is
  spendable. After that 25% is spent, a timestamp will be added to the code to
  enable further distribution.
- Add `/coinSupply` endpoint. Correctly returns total, locked and unlocked coin
  amounts.
- `testutil` package for common test setup methods.
- `/version` endpoint, which will return the current node version number and
  the HEAD commit id when build the node
- `-no-ping-log` option to disable ping/pong log output
- Check for invalid block signatures during startup and recreate the database
  if they are corrupted.
- Add methods for converting fixed-point decimal strings to droplets and
  vice versa.
- Add `make run`, `make test`, `make lint`, `make check` to `Makefile`

### Changed

- Flag peers as incoming or outgoing.
- Refactor to decouple `wallet` and `visor` package.
- Refactor `cli` package for use as a library.
- `README` improvements.
- Set default wallet's label as "Your Wallet"
- Use BIP32 mnemomic seeds by default in `address_gen`.
- Add `-x` option to `address_gen`, to generate a random base64-encoded 128-bit
  seed instead of a BIP32 mnemomic seed.
- Add `-v` option to `address_gen` to print all address information
  (pubkey, seckey, address, seed) to stdout as JSON.
- All API and CLI methods with "coin"-related arguments must be a string and
  can use decimal notation to specify coin amounts.
- CLI's `walletHistory` command prints amounts as fixed-point decimal strings.
  Previously, it printed amounts as integers representing whole skycoin amounts,
  and did not support droplets / fractional skycoins.
- A user is prevented from broadcasting a new transaction with unspent outputs
  that they have already sent as an unconfirmed transaction.

### Deprecated

- `/api/getEffectiveOutputs` is deprecated in favor of `/coinSupply`.

### Removed

- Old wallet
- `/api/create-address` endpoint (use the `address_gen` tool)

### Fixed

- Wallet folder path loading.
- #371 Fix `/wallet/spend`, will return only when pending transaction is
  confirmed.
- #443 Fix predicted balance in `/wallet/spend` API call.
- #444 Fix bug in `/blockchain/progress` API call.
- Removed globals in `gui` package that caused race condition with wallet API
  calls.
- #494 Clean invalid unconfirmed transactions during startup.
- Various race conditions around the bolt.DB blockchain DB
- Missing `strand()` call in `daemon.Visor.AnnounceTxns`.

### Security

## [0.19.1] - 2017-08-26

### Fixed

- #459 dist folder in repo out of date, wallet gui does not load

## [0.19.0] - 2017-07-11

### Added

- Add synchronize indicator when downloading blocks.
- #352 Store unspent pool in db for quick recovery when node restart
- Speed up the time the node start the browser
- Cache unspent pool in memory to speed up query action
- #411 Add button to hide seed
- #380 Move anything with heavy imports into util sub package

### Fixed

- #421 Sort wallet transaction history by time
- #398 Remove seeds from DOM
- #390 Make `go test ./src/...` work
- #383 Error during installation from skycoin source code
- #375 Node can't recovery from zero connections
- #376 Explorer api `/explorer/address` does not return spend transactions
- #373 Master node will be closed if there're no transactions need to execute
- #360 Node will crash when do ctrl+c while downloading blocks
- #350 Wallet name always 'undefined' after loading wallet from seed

[Unreleased]: https://github.com/skycoin/skycoin/compare/master...develop
[0.20.0]: https://github.com/skycoin/skycoin/compare/v0.19.1...v0.20.0
[0.19.1]: https://github.com/skycoin/skycoin/compare/v0.19.0...v0.19.1
[0.19.0]: https://github.com/skycoin/skycoin/commit/dd924e1f2de8fab945e05b3245dbeabf267f2910
