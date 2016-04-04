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
./run.sh -h
```

*Running Wallet

```
./run.sh
Goto http://127.0.0.1:6402

OR

go run ./cmd/skycoin/skycoin.go
```

Golang environment setup with gvm
---

```
sudo apt-get install bison curl git mercurial make binutils bison gcc build-essential
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source /root/.gvm/scripts/gvm

<<<<<<< HEAD
gvm install go1.4
gvm use go1.4
gvm install go1.6
gvm use go1.6 --default
```

Install
---

In China, use a VPN during golang installation with gvm. The golang repo is on a server blocked by the chinese firewall.

```
sudo apt-get install curl git mercurial make binutils gcc bzr bison libgmp3-dev -y

bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source /root/.gvm/scripts/gvm

#install golang
gvm install go1.4
gvm use go1.4
gvm install go1.5
gvm use go1.5 --default

#pull skycoin repo into the gopath
cd $GOPATH/src
git clone https://github.com/skycoin/skycoin

#create symlink of the repo
ln -s $GOPATH/src/skycoin/skycoin skycoin
=======
gvm install go1.6
gvm use go1.6

```

The skycoin repo must be in gopath, under "/src/github.com/skycoin". Otherwise golang programs cannot import the libraries.

Dependencies
---

```
go get github.com/robfig/glock
glock sync github.com/skycoin/skycoin

```

Cross Compilation
---

Install Gox:
```
go get github.com/mitchellh/gox
gox -build-toolchain
```

Compile:
```
cd $GOPATH/src/github.com/skycoin/skycoin
G="skycoin-0.3"
gox -output="$HOME/builds/$G-{{.OS}}-{{.Arch}}/$G/address_gen" ./cmd/address_gen/
gox -output="$HOME/builds/$G-{{.OS}}-{{.Arch}}/$G/skycoin" ./cmd/skycoin/
```

Local Server API
----

Run the skycoin client then

http://127.0.0.1:6420/outputs
http://127.0.0.1:6420/blockchain/blocks?start=0&end=500
http://127.0.0.1:6420/blockchain
http://127.0.0.1:6420/connections

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
