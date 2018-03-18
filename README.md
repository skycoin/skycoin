![skycoin logo](https://user-images.githubusercontent.com/26845312/32426705-d95cb988-c281-11e7-9463-a3fce8076a72.png)

# Skycoin

[![Build Status](https://travis-ci.org/skycoin/skycoin.svg)](https://travis-ci.org/skycoin/skycoin)
[![GoDoc](https://godoc.org/github.com/skycoin/skycoin?status.svg)](https://godoc.org/github.com/skycoin/skycoin)
[![Go Report Card](https://goreportcard.com/badge/github.com/skycoin/skycoin)](https://goreportcard.com/report/github.com/skycoin/skycoin)

Skycoin is a next-generation cryptocurrency.

Skycoin improves on Bitcoin in too many ways to be addressed here.

Skycoin is a small part of OP Redecentralize and OP Darknet Plan.

## Links

* [skycoin.net](https://www.skycoin.net)
* [Skycoin Blog](https://blog.skycoin.net)
* [Skycoin Blockchain Explorer](https://explorer.skycoin.net)
* [Skycoin Distribution Event](https://event.skycoin.net)

## Table of Contents

<!-- MarkdownTOC depth="2" autolink="true" bracket="round" -->

- [Installation](#installation)
    - [Go 1.9+ Installation and Setup](#go-19-installation-and-setup)
    - [Go get skycoin](#go-get-skycoin)
    - [Run Skycoin from the command line](#run-skycoin-from-the-command-line)
    - [Show Skycoin node options](#show-skycoin-node-options)
    - [Run Skycoin with options](#run-skycoin-with-options)
    - [Docker image](#docker-image)
- [API Documentation](#api-documentation)
    - [Wallet REST API](#wallet-rest-api)
    - [JSON-RPC 2.0 API](#json-rpc-20-api)
    - [Skycoin command line interface](#skycoin-command-line-interface)
- [Integrating Skycoin with your application](#integrating-skycoin-with-your-application)
- [Contributing a node to the network](#contributing-a-node-to-the-network)
- [Development](#development)
    - [Modules](#modules)
    - [Client libraries](#client-libraries)
    - [Running Tests](#running-tests)
    - [Formatting](#formatting)
    - [Code Linting](#code-linting)
    - [Dependency Management](#dependency-management)
    - [Wallet GUI Development](#wallet-gui-development)
    - [Releases](#releases)
- [Changelog](#changelog)

<!-- /MarkdownTOC -->

## Installation

### Go 1.9+ Installation and Setup

[Golang 1.9+ Installation/Setup](./INSTALLATION.md)

### Go get skycoin

```sh
go get github.com/skycoin/skycoin/...
```

This will download `github.com/skycoin/skycoin` to `$GOPATH/src/github.com/skycoin/skycoin`.

You can also clone the repo directly with `git clone https://github.com/skycoin/skycoin`,
but it must be cloned to this path: `$GOPATH/src/github.com/skycoin/skycoin`.

### Run Skycoin from the command line

```sh
cd $GOPATH/src/github.com/skycoin/skycoin
make run
```

### Show Skycoin node options

```sh
cd $GOPATH/src/github.com/skycoin/skycoin
make run-help
```

### Run Skycoin with options

```sh
cd $GOPATH/src/github.com/skycoin/skycoin
make ARGS="--launch-browser=false" run
```

### Docker image

```sh
$ docker volume create skycoin-data
$ docker volume create skycoin-wallet
$ docker run -ti --rm \
    -v skycoin-data:/data \
    -v skycoin-wallet:/wallet \
    -p 6000:6000 \
    -p 6420:6420 \
    -p 6430:6430 \
    skycoin/skycoin
```

Access the dashboard: [http://localhost:6420](http://localhost:6420).

Access the API: [http://localhost:6420/version](http://localhost:6420/version).

## API Documentation

### Wallet REST API

[Wallet REST API](src/gui/README.md).

### JSON-RPC 2.0 API

[JSON-RPC 2.0 README](src/api/webrpc/README.md).

### Skycoin command line interface

[CLI command API](cmd/cli/README.md).

## Integrating Skycoin with your application

[Skycoin Integration Documentation](INTEGRATION.md)

## Contributing a node to the network

Add your node's `ip:port` to the [peers.txt](peers.txt) file.
This file will be periodically uploaded to https://downloads.skycoin.net/blockchain/peers.txt
and used to seed client with peers.

*Note*: Do not add Skywire nodes to `peers.txt`.
Only add Skycoin nodes with high uptime and a static IP address (such as a Skycoin node hosted on a VPS).

## Development

We have two branches: `master` and `develop`.

`develop` is the default branch and will have the latest code.

`master` will always be equal to the current stable release on the website, and should correspond with the latest release tag.

### Modules

* `/src/cipher` - cryptography library
* `/src/coin` - the blockchain
* `/src/daemon` - networking and wire protocol
* `/src/visor` - the top level, client
* `/src/gui` - the web wallet and json client interface
* `/src/wallet` - the private key storage library
* `/src/api/webrpc` - JSON-RPC 2.0 API
* `/src/api/cli` - CLI library

### Client libraries

Skycoin implements client libraries which export core functionality for usage from
other programming languages. Read the corresponding README file for further details.

* `lib/cgo/` - libskycoin C client library ( [read more](lib/cgo/README.md) )

### Running Tests

```sh
make test
```

### Running Integration Tests

Run stable integration tests:

```sh
make integration-test-stable
```

or

```sh
./ci-scripts/integration-test-stable.sh -v -w
```

The `-w` option, run wallet integrations tests.

The `-v` option, show verbose logs.

Run live integration tests:

The live integration tests run against a live runnning skycoin node, so before running the test, we
need to start a skycoin node.

After the skycoin node is up, run the following command to start the live tests:

```sh
./ci-scripts/integration-test.live.sh -v
```

The above command will run all tests except the wallet tests. To run live wallet integration tests, we
need to manually specify a wallet file, the wallet must have at least `2 coins` and `256 coinhours`, and
it also must already have been loaded by the node.

We can specify the wallet by setting two environment variables: `WALLET_DIR` and `WALLET_NAME`. The `WALLET_DIR`
represents the absolute path of the wallet directory, and `WALLET_NAME` represents the wallet file name.

```sh
export WALLET_DIR=$HOME/.skycoin/wallets
export WALLET_NAME=$wallet-file-name-meet-the-requirements
```

Then run the tests with the follwoing command:

```sh
make integration-test-live
```

or

```sh
./ci-scripts/integration-test-live.sh -v -w
```

### Debugging integration tests

Run specific test case:

It's annoying and a waste of time to run all tests to see if the test we real care
is working correctly. There's an option: `-r`, which can be used to run specific test case.
For exampe: if we only want to test `TestStableAddressBalance` and see the result, we can run:

```sh
./ci-scripts/integration-test-stable.sh -v -r TestStableAddressBalance
```

Update golden files in test-fixtures files:

To update all golden files:

```sh
./ci-scripts/integration-test-live.sh -v -u
./ci-scripts/integration-test-stable.sh -v -u
```

We can also update specific test case's golden file with the `-r` option.

### Formatting

All `.go` source files should be formatted `goimports`.  You can do this with:

```sh
make format
```

### Code Linting

Install prerequisites:

```sh
make install-linters
```

Run linters:

```sh
make lint
```

### Dependency Management

Dependencies are managed with [dep](https://github.com/golang/dep).

To install `dep`:

```sh
go get -u github.com/golang/dep
```

`dep` vendors all dependencies into the repo.

If you change the dependencies, you should update them as needed with `dep ensure`.

Use `dep help` for instructions on vendoring a specific version of a dependency, or updating them.

When updating or initializing, `dep` will find the latest version of a dependency that will compile.

Examples:

Initialize all dependencies:

```sh
dep init
```

Update all dependencies:

```sh
dep ensure -update -v
```

Add a single dependency (latest version):

```sh
dep ensure github.com/foo/bar
```

Add a single dependency (more specific version), or downgrade an existing dependency:

```sh
dep ensure github.com/foo/bar@tag
```

### Wallet GUI Development

The compiled wallet source should be checked in to the repo, so that others do not need to install node to run the software.

Instructions for doing this:

[Wallet GUI Development README](src/gui/static/README.md)

### Releases

0. If the `master` branch has commits that are not in `develop` (e.g. due to a hotfix applied to `master`), merge `master` into `develop`
1. Compile the `src/gui/dist/` to make sure that it is up to date (see [Wallet GUI Development README](src/gui/static/README.md))
2. Update all version strings in the repo (grep for them) to the new version
3. Update `CHANGELOG.md`: move the "unreleased" changes to the version and add the date
4. Merge these changes to `develop`
5. On the `develop` branch, make sure that the client runs properly from the command line (`./run.sh`)
6. Build the releases and make sure that the Electron client runs properly on Windows, Linux and macOS. Delete these releases when done.
7. Make a PR merging `develop` into `master`
8. Review the PR and merge it
9. Tag the master branch with the version number. Version tags start with `v`, e.g. `v0.20.0`.
10. Make sure that the client runs properly from the `master` branch
11. Create the release builds from the `master` branch (see [Create Release builds](electron/README.md))

If there are problems discovered after merging to master, start over, and increment the 3rd version number.
For example, `v0.20.0` becomes `v0.20.1`, for minor fixes.

#### Creating release builds

[Create Release builds](electron/README.md).

## Changelog

[CHANGELOG.md](CHANGELOG.md)
