# Skycoin

[![GoDoc](https://godoc.org/github.com/skycoin/skycoin?status.svg)](https://godoc.org/github.com/skycoin/skycoin) [![Go Report Card](https://goreportcard.com/badge/github.com/skycoin/skycoin)](https://goreportcard.com/report/github.com/skycoin/skycoin)

Skycoin is a next-generation cryptocurrency.

Skycoin improves on Bitcoin in too many ways to be addressed here.

Skycoin is small part of OP Redecentralize and OP Darknet Plan.

## Installation

For detailed installation instructions, see [Installing Skycoin](../../wiki/Installation).

## For OSX

Install [homebrew](brew.sh), if you don't have it yet.

```sh
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

Install the latest version of golang

```sh
brew install go
```

Setup $GOPATH variable, add it to ~/.bash_profile (or bashrc). After editing, open a new tab
Add to `bashrc` or `bash_profile`

```sh
export GOPATH=/Users/<username>/go
export PATH=$PATH:$GOPATH/bin

```

Install Mercurial and Bazaar

```sh
brew install mercurial bzr
```

Fetch the latest code of skycoin from the github repository

```sh
go get github.com/skycoin/skycoin
```

Change your current directory to $GOPATH/src/github.com/skycoin/skycoin

```sh
cd $GOPATH/src/github.com/skycoin/skycoin
```

Run Wallet

```sh
./run.sh

OR
go run ./cmd/skycoin/skycoin.go

For Options
go run ./cmd/skycoin/skycoin.go --help
```

## For linux

```sh
sudo apt-get install curl git mercurial make binutils gcc bzr bison libgmp3-dev screen -y
```

## Setup Golang

use gvm or download binary and follow instructions.

### Golang ENV setup with gvm

In China, use `--source=https://github.com/golang/go` to bypass firewall when fetching golang source.

```sh
sudo apt-get install bison curl git mercurial make binutils bison gcc build-essential
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source $HOME/.gvm/scripts/gvm

gvm install go1.4 --source=https://github.com/golang/go
gvm use go1.4
gvm install go1.8
gvm use go1.8 --default
```

If you open up new terminal and the go command is not found then add this to .bashrc . GVM should add this automatically.

```sh
[[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"
gvm use go1.8 >/dev/null
```

The skycoin repo must be in $GOPATH, under `src/github.com/skycoin`. Otherwise golang programs cannot import the libraries.

Pull skycoin repo into the gopath, note: puts the skycoin folder in $GOPATH/src/github.com/skycoin/skycoin

```sh
go get -v github.com/skycoin/skycoin/...

# create symlink of the repo
cd $HOME
ln -s $GOPATH/src/github.com/skycoin/skycoin skycoin
```

## Dependencies

Dependencies are managed with [gvt](https://github.com/FiloSottile/gvt).

To install gvt:

```sh
go get -u github.com/FiloSottile/gvt
```

gvt vendors all dependencies into the repo.

If you change the dependencies, you should update them as needed with `gvt fetch`, `gvt update`, `gvt delete`, etc.

Refer to the [gvt documentation](https://github.com/FiloSottile/gvt) or `gvt help` for further instructions.

## Run A Skycoin Node

```sh
cd skycoin
screen
go run ./cmd/skycoin/skycoin.go
```

then ctrl+A then D to exit screen
screen -x to reattach screen

### Todo

Use gvm package set, so repo does not need to be symlinked. Does this have a default option?

```sh
gvm pkgset create skycoin
gvm pkgset use skycoin
git clone https://github.com/skycoin/skycoin
cd skycoin
go install
```

### Cross Compilation

Install Gox:

```sh
go get github.com/mitchellh/gox
```

Compile:

```sh
gox --help
gox [options] cmd/skycoin/
```

## Local Server API

See the api details [here](src/gui/READEME.md).

## Skycoin explorer

```link
http://explorer.skycoin.net
```

## Modules

* /src/cipher - cryptography library
* /src/coin - the blockchain
* /src/daemon - networking and wire protocol
* /src/visor - the top level, client
* /src/gui - the web wallet and json client interface
* /src/wallet - the private key storage library

## Meshnet

```sh
go run ./cmd/mesh/*.go -config=cmd/mesh/sample/config_a.json
go run ./cmd/mesh/*.go -config=cmd/mesh/sample/config_b.json
```

## Meshnet reminders

* one way latency
* two way latency (append), latency between packet and ack
* service handler (ability to append services to meshnet)
* uploading bandwidth, latency measurements over time
* end-to-end route instrumentation

## Rebuilding Wallet HTML

```sh
cd src/gui/static
npm install
gulp build
```

## Release Builds

```sh
cd /src/gui/static
npm install
gulp dist
```

## Skycoin command line interface

See the doc of command line interface [here](cmd/cli/README.md).

## WebRPC

See the doc of webrpc [here](src/api/webrpc/README.md).

## Development

We mainly has two branches: master and develop. The develop is the default branch as you can see, all latest code will be updated here.

The master branch will always be in run ready state and will only be updated when we need to release a new version.