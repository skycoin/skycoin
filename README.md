skycoin [![GoDoc](https://godoc.org/github.com/skycoin/skycoin?status.svg)](https://godoc.org/github.com/skycoin/skycoin) [![Go Report Card](https://goreportcard.com/badge/github.com/skycoin/skycoin)](https://goreportcard.com/report/github.com/skycoin/skycoin)
=======

Skycoin is a next-generation cryptocurrency.

Skycoin improves on Bitcoin in too many ways to be addressed here.

Skycoin is small part of OP Redecentralize and OP Darknet Plan.

Installation
------------

* For detailed installation instructions, see [Installing Skycoin](../../wiki/Installation)*

## For linux:

```sh
$ sudo apt-get install curl git mercurial make binutils gcc bzr bison libgmp3-dev screen -y
```

## For OSX:

1) Install [homebrew](brew.sh), if you don't have it yet
```
$ /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

2) Install the latest version of golang
```
$ brew install go
```

3) Setup $GOPATH variable, add it to ~/.bash_profile (or bashrc). After editing, open a new tab
Add to `bashrc` or `bash_profile`
```sh
$ export GOPATH=/Users/<username>/go
$ export PATH=$PATH:$GOPATH/bin

```

4) Install Mercurial and Bazaar
```
$ brew install mercurial bzr
```

5) Fetch the latest code of skycoin from the github repository
```
$ go get github.com/skycoin/skycoin
```

6) Change your current directory to $GOPATH/src/github.com/skycoin/skycoin
```
$ cd $GOPATH/src/github.com/skycoin/skycoin
```

7) Run the node ;)
```
$ ./run.sh -h
```

8) Running Wallet

```
$ ./run.sh
```

Then open `http://127.0.0.1:6402` in a browser.

## Golang ENV setup with gvm

In China, use `--source=https://github.com/golang/go` to bypass firewall when fetching golang source

```
$ sudo apt-get install bison curl git mercurial make binutils bison gcc build-essential
$ bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
$ source $HOME/.gvm/scripts/gvm

$ gvm install go1.4 --source=https://github.com/golang/go
$ gvm use go1.4
$ gvm install go1.6
$ gvm use go1.6 --default
```

If you open up new terminal and the go command is not found then add this to .bashrc . GVM should add this automatically
```
$ [[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"
$ gvm use go1.6 >/dev/null
```


The skycoin repo must be in $GOPATH, under `src/github.com/skycoin`. Otherwise golang programs cannot import the libraries.

```
#pull skycoin repo into the gopath
#note: puts the skycoin folder in $GOPATH/src/github.com/skycoin/skycoin
go get -v github.com/skycoin/skycoin/...

#create symlink of the repo
$ cd $HOME
$ ln -s $GOPATH/src/github.com/skycoin/skycoin skycoin
```

Dependencies
------------

Dependencies are managed with [gvt](https://github.com/FiloSottile/gvt).

To install gvt:
```
$ go get -u github.com/FiloSottile/gvt
```

gvt vendors all dependencies into the repo.

If you change the dependencies, you should update them as needed with `gvt fetch`, `gvt update`, `gvt delete`, etc.

Refer to the [gvt documentation](https://github.com/FiloSottile/gvt) or `gvt help` for further instructions.

Running A Skycoin Node
----------------------

```
$ cd skycoin
$ screen
$ go run ./cmd/skycoin/skycoin.go
#then ctrl+A then D to exit screen
#screen -x to reattach screen
```

##Todo

Use gvm package set, so repo does not need to be symlinked. Does this have a default option?
```
$ gvm pkgset create skycoin
$ gvm pkgset use skycoin
$ git clone https://github.com/skycoin/skycoin
$ cd skycoin
$ go install
```

##Cross Compilation

Install Gox:
```
$ go get github.com/mitchellh/gox
```

Compile:
```
$ gox --help
$ gox [options] cmd/skycoin/
```

Local Server API
----------------

Run the skycoin client then
```
http://127.0.0.1:6420/wallets
http://127.0.0.1:6420/outputs
http://127.0.0.1:6420/blockchain/blocks?start=0&end=500
http://127.0.0.1:6420/blockchain
http://127.0.0.1:6420/connections
```

```
http://127.0.0.1:6420/wallets
- to get your wallet seed. Write this down

http://127.0.0.1:6420/wallet/balance?id=2016_02_17_9671.wlt
- to get wallet balance (use wallet filename as id)
- TODO: allow addresses for balance check

http://127.0.0.1:6420/outputs to see outputs (address balances)

http://127.0.0.1:6420/blockchain/blocks?start=0&end=5000 to see all blocks and transactions.

http://127.0.0.1:6420/network/connections to check network connections

http://127.0.0.1:6420/blockchain to check blockchain head
```

Public API
----------

This is a public server. You can use these urls on local host too, with the skycoin client running.
```
http://skycoin-chompyz.c9.io/outputs
http://skycoin-chompyz.c9.io/blockchain/blocks?start=0&end=500
http://skycoin-chompyz.c9.io/blockchain
http://skycoin-chompyz.c9.io/connections
```

Modules
-------

```
/src/cipher - cryptography library
/src/coin - the blockchain
/src/daemon - networking and wire protocol
/src/visor - the top level, client
/src/gui - the web wallet and json client interface
/src/wallet - the private key storage library
```

Meshnet
-------

```
$ go run ./cmd/mesh/*.go -config=cmd/mesh/sample/config_a.json
$ go run ./cmd/mesh/*.go -config=cmd/mesh/sample/config_b.json
```

Meshnet reminders
-----------------

- one way latency
- two way latency (append), latency between packet and ack
- service handler (ability to append services to meshnet)
- uploading bandwidth, latency measurements over time
- end-to-end route instrumentation

Rebuilding Wallet HTML
----------------------

```sh
$ npm install
$ gulp build
```

Release Builds
----

```sh
$ npm install
$ gulp dist
```
