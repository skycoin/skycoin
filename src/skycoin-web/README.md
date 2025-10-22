![skycoin web](https://user-images.githubusercontent.com/8619106/56056491-f4bd8d80-5d79-11e9-8e55-8ebf2668edac.png)

# Skycoin web client

[![Build Status](https://travis-ci.org/skycoin/skycoin-web.svg?branch=master)](https://travis-ci.org/skycoin/skycoin-web)
[![Known Vulnerabilities](https://snyk.io/test/github/skycoin/skycoin-web/badge.svg)](https://snyk.io/test/github/skycoin/skycoin-web)

The Skycoin web client provides a lite multi-coin wallet which can be ran from the browser, using a remote full node exposing selected back-end functions.

## Prerequisites

The Skycoin web interface requires Node 6.9.0 or higher, together with NPM 3 or higher.

## Installation

This project is generated using Angular CLI, therefore it is adviced to first run `npm install -g @angular/cli`.

To install all Angular, Angular CLI and all other libraries, you will then have to run `npm install`.

You will only have to run this again, if any dependencies have been changed in the `package.json` file.

## Compiling new target files

To compile new target files, you will have to run: `npm run build`. To create a dev build, use `npm run build-dev`
instead. To compile a build capable of being run from the local file system, use `npm run build-for-local-fs`.

## Docker

To build a production ready Docker image, you have to run:

```sh
docker build -t image_name:tag -f docker/images/Dockerfile .
```

or you can pull the image from the official repository:

```sh
docker pull skycoin/skycoin-web:<tag>
```

To create a container you must run:

```sh
docker run -d --name <name> -p 8001:80 skycoin/skycoin-web:latest
```

Then, you can access to the web wallet by visiting `http://your_IP:8001`

## Development server

Run `npm start` for a dev server. Navigate to `http://localhost:4200/`.

## Back-end

As this is a lite client, it requires a back-end to retrieve the blockchain state and inject new transactions. While using
a dev build, the wallet tries to connect to a Skycoin node running locally, which must be started with the `-disable-csrf`
option. The production builds connects to a public node set up at `https://node.skycoin.net`. The user can connect the wallet
to another backend by entering its URL.

### Minimum back-end version

The wallet is currently compatible with nodes running v0.24.0 and above. As it is not possible to recover the number of active
decimals in version 0.24.0, the wallet assumes that it is 6. The same happens with burn rate, so it is assumed to be 50%.
Because of this, it is advisable to connect to a node running v0.25.0 or higher.

## Multi-coin compatibility

This wallet have options for switching between different coins and is compatible with most coins based on the Skycoin
technology (as long as they meet some version requirements and do not have modifications in areas necessary for
compatibility). You can find information about how to add a new coins in [MULTICOIN-INFO.md](MULTICOIN-INFO.md).

## Cryptography

For security reasons, this wallet does not store the the seeds, neither in the backend nor in the client, so the user must
maintain a backup of them. The seeds are also not transmitted anywhere outside the application code.

To maintain maximum compatibility with Skycoin's original code, the cryptographic functions used by this wallet are in the
[skycoin-lite](https://github.com/skycoin/skycoin-lite) repository, coded in Go (Golang), and are compiled to wasm.
The compiled wasm file is [src/assets/scripts/skycoin-lite.wasm](/src/assets/scripts/skycoin-lite.wasm).

The mnemonic phrases are created using the [bip39 library](https://www.npmjs.com/package/bip39). To create mnemonic phrases
with enough entropy, that library uses the [randombytes library](https://www.npmjs.com/package/randombytes), which uses
[Crypto.getRandomValues()](https://developer.mozilla.org/en-US/docs/Web/API/Crypto/getRandomValues) for obtaining
cryptographically strong random values. The use of the [bip39 library](https://www.npmjs.com/package/bip39) is limited to
the creation of mnemonic phrases, the code that creates the binary seeds is in
[skycoin-lite](https://github.com/skycoin/skycoin-lite).

## Pre-release testing

Perform these actions before releasing:

- Test if all options are working correctly when opening the wallet on the test server (using `npm start`). Load a wallet with balance and try to send coins, to verify if the functions for loading wallets and creating transactions work correctly.

- Test if all options are working correctly when opening the wallet from the local file system (the compilation is done by executing `make build-for-local-fs`). In the local version it is important to check all the pages and make sure that all the images load correctly. For this it is useful to go through all pages while the development console is open.
