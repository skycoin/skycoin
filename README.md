# Skycoin

[![Build Status](https://travis-ci.org/skycoin/skycoin.svg)](https://travis-ci.org/skycoin/skycoin)
[![GoDoc](https://godoc.org/github.com/skycoin/skycoin?status.svg)](https://godoc.org/github.com/skycoin/skycoin)
[![Go Report Card](https://goreportcard.com/badge/github.com/skycoin/skycoin)](https://goreportcard.com/report/github.com/skycoin/skycoin)

Skycoin is a next-generation cryptocurrency.

Skycoin improves on Bitcoin in too many ways to be addressed here.

Skycoin is small part of OP Redecentralize and OP Darknet Plan.

## Links

* [skycoin.net](https://www.skycoin.net)
* [Skycoin Blog](https://blog.skycoin.net)
* [Skycoin Blockchain Explorer](https://explorer.skycoin.net)
* [Skycoin Distribution Event](https://event.skycoin.net)

## Table of Contents

<!-- MarkdownTOC depth="2" autolink="true" bracket="round" -->

- [Installation](#installation)
    - [Go Installation and Setup](#go-19-installation-and-setup)
    - [Go get skycoin](#go-get-skycoin)
    - [Run Skycoin from the command line](#run-skycoin-from-the-command-line)
    - [Show Skycoin node options](#show-skycoin-node-options)
    - [Run Skycoin with options](#run-skycoin-with-options)
- [API Documentation](#api-documentation)
    - [Wallet REST API](#wallet-rest-api)
    - [JSON-RPC 2.0 API](#json-rpc-20-api)
    - [Skycoin command line interface](#skycoin-command-line-interface)
- [Development](#development)
    - [Modules](#modules)
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

[Golang 1.9+ Installation/Setup](./Installation.md)

### Go get skycoin

```sh
go get https://github.com/skycoin/skycoin/...
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

## API Documentation

### Wallet REST API

[Wallet REST API](src/gui/README.md).

### JSON-RPC 2.0 API

[JSON-RPC 2.0 README](src/api/webrpc/README.md).

### Skycoin command line interface

[CLI command API](cmd/cli/README.md).

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

### Running Tests

```sh
make test
```

### Formatting

All `.go` source files should be formatted with `gofmt` or `goimports`.

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

After adding a new dependency (with `dep ensure`), run `dep prune` to remove any unnecessary subpackages from the dependency.

When updating or initializing, `dep` will find the latest version of a dependency that will compile.

Examples:

Initialize all dependencies:

```sh
dep init
dep prune
```

Update all dependencies:

```sh
dep ensure -update -v
dep prune
```

Add a single dependency (latest version):

```sh
dep ensure github.com/foo/bar
dep prune
```

Add a single dependency (more specific version), or downgrade an existing dependency:

```sh
dep ensure github.com/foo/bar@tag
dep prune
```

### Wallet GUI Development

The compiled wallet source should be checked in to the repo, so that others do not need to install node to run the software.

Compile and add the wallet source to git:

```sh
cd src/gui/static
npm install
npm run build
git add .
```

### Releases

*TODO: Full instructions on doing a release. Need instructions on updating version number in source code, and running scripts to build releases*

When ready to do a release, a pull request merging `develop` into `master` must be made.

After merging to `master`, tag the branch with the version number.

Once `master` branch is updated, `git checkout master` and create the release builds.

#### Creating release builds

[Create Release builds](electron/README.md).

## Changelog

[CHANGELOG.md](CHANGELOG.md)
