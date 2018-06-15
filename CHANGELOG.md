# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add environment variable `DATA_DIR` in CLI's
- `USE_CSRF` environment variable for CLI, if the remote node has CSRF enabled (CSRF is enabled by default, use `-disable-csrf` to disable)
- `cli showConfig` command to echo the cli's configuration back to the user
- Option to generate 12/24 word seed when creating new wallet
- libskycoin 0.0.1 released with bindings for cipher/address, cipher/hash, cipher/crypto, cli/create_rawtx
- Add `-version` flag to show node version
- Add transaction verification step to "Send" page
- Add more details about transaction in transaction history
- Add CLI `encryptWallet` command
- Add CLI `decryptWallet` command
- Add CLI `showSeed` command
- Add `password` argument to the CLI commands of `addPrivateKey`, `createRawTransaction`, `generateAddresses`, `generateWallet`, `send`
- Support for decoding map values in cipher binary encoder
- Expose known block height of peer in brand new `height` field added in responses of `GET /api/v1/network/connections` API endpoints
- `-verify-db` option (default true), will verify the database integrity during startup and exit if a problem is found
- `-reset-corrupt-db` option (default false) will verify the database integrity during startup and reset the db if a problem is found
- `GET /explorer/address`: add `fee` to transaction objects and `calculated_hours` to transaction inputs
- Test data generator and test suite for verification of alternative `cipher` implementations
- Begin `/api/v2` API endpoints. These endpoints are in beta and subject to change.
- Add `POST /api/v2/transaction/verify` API endpoint
- Add `POST /api/v2/address/verify` API endpoint
- Add advanced spend UI
- Add `ignore_unconfirmed` option to `POST /api/v1/wallet/transaction` to allow transactions to be created or spent even if there are unspent outputs in the unconfirmed pool.
- Add `uxouts` to `POST /api/v1/wallet/transaction`, to allow specific unspent outputs to be used in a transaction.
- Add Dockerfile in docker/images/dev-cli to build a docker image suitable for development.
- Coin creator tool, `cmd/newcoin`, to quickly bootstrap a new fiber coin

### Fixed

- Reduce connection disconnects, improves syncing
- Fix #1171, update CLI to support wallet encryption
- Use `bolt.Tx` correctly for read operations
- Docker images for `arm32v5` and `ar32v7` architectures by using busybox as base in docker/images/mainnet/Dockerfile and docker/images/mainnet/hooks/

### Changed

- JSON 2.0 RPC interface (used by the CLI tool) is now served on the same host interface as the REST API, port `6420`. The additional listener has been removed.
- CLI's `RPC_ADDR` environment variable must now start with a scheme e.g. `http://127.0.0.1:6420`, previously it did not use a scheme.
- API response will be gzip compressed if client sends request with 'Accept-Encoding' contains 'gzip' in the header.
- `GET /api/v1/wallet/balance` and `GET /api/v1/balance` now return an address balance list as well.
- API endpoints are prefixed with `/api/v1/`. API endpoints without the `/api/v1/` prefix are deprecated but can be enabled with `-enable-unversioned-api`. Please migrate to use `/api/v1/` prefix in URLs.
- Enable message protocol upgrade
- `change_address` is no longer required in `POST /api/v1/wallet/transaction`. If not provided, `change_address` will default to one of the addresses being spent from.
- In `POST /api/v1/wallet/transaction`, for `auto` type `share` mode requests, if extra coinhours remain after applying the `share_factor` but change cannot be made due to insufficient coins, the `share_factor` will switch to `1.0`.
- Support automatic port allocation of the API interface by specifying port 0
- The web interface / API port is randomly allocated for the precompiled standalone client and electron client released on the website.
  If you are using the CLI tool or another API client to communicate with the standalone client, use `-web-interface-port=6420` to continue using port 6420.
  If the program is run from source (e.g. `go run`, `run.sh`, `make run`) there is no change, the API will still be on port 6420.
- Change number of outgoing connections to 8 from 16

### Removed

- Remove `-rpc-interface-addr`, `-rpc-interface-port` options.  The RPC interface is now on default port `6420` with the REST API.
- Remove `-rpc-thread-num` option
- Remove `-connect-to` option
- Remove `-print-web-interface-address` option
- Remove support for go1.9

## [0.23.0] - 2018-04-22

### Added

- Add wallet setup wizard
- Add wallet encryption, using chacha20+poly1305 for encryption and authentication and scrypt for key derivation. Encrypted data is stored in the wallet file in a `"secrets"` metadata field
- Add `GET /health` endpoint
- Add `POST /wallet/transaction` API endpoint, creates a transaction, allowing control of spending address and multiple destinations
- Add `POST /wallet/encrypt` API endpoint, encrypts wallet and returns encrypted wallet without sensitive data
- Add `POST /wallet/decrypt` API endpoint, decrypts wallet and returns decrypted wallet without sensitive data
- Add `POST /wallet/seed` API endpoint, returns the seed of an encrypted wallet. Unencrypted wallets will not expose their seeds over the API. Requires `-enable-seed-api` option
- `-enable-seed-api` option to enable `POST /wallet/seed`
- Add `"size"` to block API response data (affects `GET /block`, `GET /blocks` and `GET /last_blocks`)
- Write [specification for skycoin URIs](https://github.com/skycoin/skycoin#uri-specification) (based upon bip21)

### Fixed

- #1309, Float imprecision error in frontend malformed some spend amounts, preventing the spend
- Fix one aspect of sync stalling caused by a 5-second blocking channel write by switching it to a non-blocking write, decreasing timeouts and increasing buffer sizes

### Changed

- `GET /wallet` API endpoint, remove sensitive data from the response, and fix the data format to be the same as `POST /wallet/create`
- `GET /wallets` API endpoint, remove sensitive data from the response
- `POST /wallet/create` API endpoint, add `encrypt(bool)` and `password` argument
- `POST /wallet/newAddress` API endpoint, add `password` argument
- `POST /wallet/spend` API endpoint, add `password` argument
- Change `-disable-wallet-api` to `-enable-wallet-api`, and disable the wallet API by default
- `-launch-browser` is set to false by default
- A default wallet will not be created on startup if there is no wallet. Instead, the wallet setup wizard will run
- Replace [op/go-logging](https://github.com/op/go-logging) with [logrus](https://github.com/sirupsen/logrus)
- Disable JSON-RPC 2.0 interface when running the application with `run.sh` and electron
- Whitespace will be trimmed from the seed string by the frontend client before creating or loading a wallet
- Notify the user when their wallets have unconfirmed transactions
- Return an error when providing a transaction that spends to the null address in `POST /injectTransaction`
- Change accepted `-log-level` values to `debug`, `info`, `warn`, `error`, `fatal` and `panic` (previously were `debug`, `info`, `notice`, `warning`, `error` and `critical`)
- Default log level is `info`

### Removed

- Remove `"seed"`, `"lastSeed"` and `"secret_key"` in address entries from wallet API responses. A wallet's seed can be accessed through `POST /wallet/seed` only if the wallet is encrypted and the node is run with `-enable-seed-api`
- Remove unused `-logtogui` and `-logbufsize` options

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
- `/wallet/unload` API method, removes the wallet of given id from wallet services

### Fixed

- #1021, remove `SendOr404` and `SendOr500` as they do not work properly due to typed nils
- Add Read, Write and Idle timeouts to the HTTP listener, preventing file descriptor leaks
- Support absolute and relative paths for `-data-dir` option
- Prevent creating transactions whose size exceeds the maximum block size
- Check addition and multiplication uint64 overflow
- Keep trusted peers in the peerlist permanently, even if they are unreachable
- #885, Add `Host` header check to localhost HTTP interfaces to prevent DNS rebinding attacks
- #896, Add CSRF check to wallet API
- Fix base58 address parsing, which allowed leading invalid characters and treated unknown characters as a '1'
- Fix occasional error which causes blockchain progress not to be shown in front-end

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

- CLI's `walletBalance` and `addressBalance` commands return aggregate balances for confirmed, spendable and expected balances. Coins are formatted as droplet strings. Hours added as strings.
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
[0.23.0]: https://github.com/skycoin/skycoin/compare/v0.22.0...v0.23.0
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
