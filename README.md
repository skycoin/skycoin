skycoin
=======

Skycoin is a next-generation cryptocurrency.

Skycoin improves on Bitcoin in too many ways to be addressed here.

Skycoin is small part of OP Redecentralize and OP Darknet Plan.

Installation
------------

*For detailed installation instructions, see [Installing Skycoin](../../wiki/Installation)*

For linux:
sudo apt-get install curl git mercurial make binutils gcc bzr bison libgmp3-dev -y

OSX:
brew install mercurial bzr

```
./setup.sh
./run.sh -h
```

*Running Wallet

```
./run.sh
Goto http://127.0.0.1:6402

OR

go run ./cmd/skycoin/skycoin.go
```

Public API
----

This is a public server. You can use these urls on local host too, with the skycoin client running.

http://skycoin-chompyz.c9.io/outputs
http://skycoin-chompyz.c9.io/blockchain/blocks?start=0&end=500
http://skycoin-chompyz.c9.io/blockchain
http://skycoin-chompyz.c9.io/connections

Modules
-----

```
/src/cipher - cryptography library
/src/coin - the blockchain
/src/daemon - networking and wire protocol
/src/visor - the top level, client
/src/gui - the web wallet and json client interface
/src/wallet - the private key storage library
```
