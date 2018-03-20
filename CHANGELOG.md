# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Fixed

### Changed

### Removed

## [0.22.0] - 2018-03-20

### Added

- go1.10 support
- Add Dockerfile
- Add libskycoin C API wrapper
- New wallet UI
- Notify the user when a new version is available
- CLI and GUI integration tests against a stable and live blockchain
- #877, Add `-disable-wallet-api` CLI option
- HTTP API client
- `/richlist` API method, returns top n address balances
- `/addresscount` API method, returns the number of addresses that have any amount of coins
- `/transactions` API method, returns transactions of addresses

### Fixed

- Add Read, Write and Idle timeouts to the HTTP listener, preventing file descriptor leaks
- Support absolute and relative paths for `-data-dir` option
- Prevent creating transactions whose size exceeds the maximum block size
- Check addition and multiplication uint64 overflow
- Keep trusted peers in the peerlist permanently, even if they are unreachable
- #885, Add `Host` header check to localhost HTTP interfaces to prevent DNS rebinding attacks
- #896, Add CSRF check to wallet API
- Fix base58 address parsing, which allowed leading invalid characters and treated unknown characters as a '1'

### Changed

- #1080, `/wallet/transactions` now returns a proper json object with pending transactions under `transactions` key
- #951, cap cli createRawTransaction and send command coinhour distribution, coinhours are capped to a maximum of receiving coins for the address with a minimum of 1 coinhour
- Upgrade to Angular 5
- Add `total_coinhour_supply` and `current_coinhour_supply` to `/coinSupply` endpoint
- #800, Add entropy parameter to `/wallet/newSeed` endpoint. Entropy can be 128 (default) or 256, corresponding to 12- and 24-word seeds respectively
- #866, Include coins and hours in `/explorer/address` inputs

### Removed

- Remove `/lastTxs` API endpoint
- Remove `/logs` and log buffering due to possible crash
- Remove `/wallets/reload` endpoint
- Remove deprecated `/api/getEffectiveOutputs`, use `/coinSupply`.

## [0.21.1] - 2017-12-14

### Fixed

- Fix blank page issue in windows gui wallet, which was caused by misusing the flag of -download-peers-list in electron.

## [0.21.0] - 2017-12-10

### Added

- Require transactions to have an input with non-zero coinhours
- Add `-peerlist-size` and `-max-outgoing-connections` CLI options
- Add `-download-peerlist` and `-peerlist-url` CLI options, to get peers from a URL
- For electron clients, download a list of peers from https://downloads.skycoin.net/blockchain/peers.txt by default

### Fixed

- Fix change hours calculation. Previous gave 1/8 to change and destination addresses; now gives 1/4 to each
- #653, the peerlist size was too small and could be easily filled up; default changed to 65535 from 1000

### Changed

- CLI's `walletBalance` and `addressBalance` commands return aggregate balances for confirmed, spendable and expected balances.  Coins are formatted as droplet strings.  Hours added as strings.
- When splitting an odd number of hours in a spend, give the extra hour to the fee
- Add `block_seq` to `get_outputs` and `/outputs` API response
- Improve UxOut spend selection. Previously, they were spent oldest first. Now they are spent to ensure a non-zero coinhour input and otherwise minimize coinhours.
- `create_rawtx` will try to minimize the number of UxOuts used to create a transaction.
- `/wallet/spend` will try to maximize the number of UxOuts used to create a transaction.
- Update the default peerlist size to 65535 from 1000
- When loading a wallet, 100 addresses will be scanned ahead to find one with a balance

## [0.20.4] - 2017-11-22

### Added

- Add `/logs` api to filter skycoin logs, so that we can add a debug panel to the GUI wallet to show logs

## [0.20.3] - 2017-10-23

### Fixed

- Fix block sync stall (mostly affected Windows users, but any OS was potentially affected)

## [0.20.2] - 2017-10-12

### Fixed

- Fixed Linux .AppImage missing "Category" field
- Clean up electron build script, switch to yarn and remove gulp

## [0.20.1] - 2017-10-12

### Fixed

- Fixed app icon padding

## [0.20.0] - 2017-10-10

### Added

- New wallet frontend in angular4. This is a complete rewrite and fixes many of the old wallet issues.
- New wallet has preliminary support for OTC functionality
- Create `webrpc.Client` for JSON-2.0 RPC calls.
- Add this CHANGELOG.md file.
- Add Installation.md file, with install instructions for go.
- Timelock distribution addresses. The first 25% of the distribution is spendable. After that 25% is spent, a timestamp will be added to the code to enable further distribution.
- Add `/coinSupply` endpoint. Correctly returns total, locked and unlocked coin amounts.
- `testutil` package for common test setup methods.
- `/version` endpoint, which will return the current node version number and the HEAD commit id when build the node
- `-no-ping-log` option to disable ping/pong log output
- Check for invalid block signatures during startup and recreate the database if they are corrupted.
- Add methods for converting fixed-point decimal strings to droplets and vice versa.
- Add `make run`, `make test`, `make lint`, `make check` to `Makefile`

### Changed

- Flag peers as incoming or outgoing.
- Refactor to decouple `wallet` and `visor` package.
- Refactor `cli` package for use as a library.
- `README` improvements.
- Set default wallet's label as "Your Wallet"
- Use BIP32 mnemomic seeds by default in `address_gen`.
- Add `-x` option to `address_gen`, to generate a random base64-encoded 128-bit seed instead of a BIP32 mnemomic seed.
- Add `-v` option to `address_gen` to print all address information (pubkey, seckey, address, seed) to stdout as JSON.
- All API and CLI methods with "coin"-related arguments must be a string and can use decimal notation to specify coin amounts.
- CLI's `walletHistory` command prints amounts as fixed-point decimal strings. Previously, it printed amounts as integers representing whole skycoin amounts, and did not support droplets / fractional skycoins.
- A user is prevented from broadcasting a new transaction with unspent outputs that they have already sent as an unconfirmed transaction.

### Deprecated

- `/api/getEffectiveOutputs` is deprecated in favor of `/coinSupply`.

### Removed

- Old wallet
- `/api/create-address` endpoint (use the `address_gen` tool)

### Fixed

- Wallet folder path loading.
- #371 Fix `/wallet/spend`, will return only when pending transaction is confirmed.
- #443 Fix predicted balance in `/wallet/spend` API call.
- #444 Fix bug in `/blockchain/progress` API call.
- Removed globals in `gui` package that caused race condition with wallet API calls.
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
[0.22.0]: https://github.com/skycoin/skycoin/compare/v0.21.1...v0.22.0
[0.21.1]: https://github.com/skycoin/skycoin/compare/v0.21.0...v0.21.1
[0.21.0]: https://github.com/skycoin/skycoin/compare/v0.20.4...v0.21.0
[0.20.4]: https://github.com/skycoin/skycoin/compare/v0.20.3...v0.20.4
[0.20.3]: https://github.com/skycoin/skycoin/compare/v0.20.2...v0.20.3
[0.20.2]: https://github.com/skycoin/skycoin/compare/v0.20.1...v0.20.2
[0.20.1]: https://github.com/skycoin/skycoin/compare/v0.20.0...v0.20.1
[0.20.0]: https://github.com/skycoin/skycoin/compare/v0.19.1...v0.20.0
[0.19.1]: https://github.com/skycoin/skycoin/compare/v0.19.0...v0.19.1
[0.19.0]: https://github.com/skycoin/skycoin/commit/dd924e1f2de8fab945e05b3245dbeabf267f2910
