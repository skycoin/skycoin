skycoin
=======

Skycoin is a next-generation post proof-of-work cryptocurrency.

It is currently in `alpha`.

No Mining
---------

Satoshi has stated that mining is Bitcoin's only flaw. Miners are the greatest threat to Bitcoin's survival. Mining introduces human factors and greed into the survivability of Bitcoin as a currency.

Problems with Mining:
- Dishonest miners may undermine trust in Bitcoin
- We are close to a 51% attack.
- gigahash.io has abused its hashing power to steal Bitcoin from gambling sites
- cex.io allows people to rent a 51% attack by the minute
- mining pools will form cartels to orphan blocks mined by non-cartel members, to steal block rewards
- mining cartels will rent capacity to orphan blocks from other pools, if the reward is greater than the cost
- pools have begun capping block size to drive up transaction fees
- Every month mining power is becoming more concentrated and centralized in the hands of fewer people
- Mining ensures majority of new crypto-coins go to people with botnets, GPU farms and ASICs not available to the public.
- No one benefits from mining except miners and electricity companies
- electricity used by mining is wasted.
- the costs of mining is passed on to every Bitcoin holder
- Mining increases transaction fees to unsustainable levels when block rewards decrease. Satoshi has stated that this will make Bitcoin unusable for microtransactions. Bitcoin transaction fees may reach a level higher than the banking system Bitcoin was created to replace.
- If Bitcoin reaches viability as the next global reserve currency, the incentives for nation states to monopolize and control mining are too great.
- Miners sell Bitcoins to pay equipment and electricity costs, driving down the price of Bitcoin. Every Bitcoin holder pays the cost of mining.

Skycoin was designed to eliminate mining completely and eliminate the problems it creates. Skycoin replaces reliance on miners with reliance on mathematics.
- Skycoin is more secure because it does not rely upon the good will of miners
- Skycoin transactions will be cheaper because there are no mining costs being passed on to users
- Skycoin transactions are not subject to 51% attacks by mining cartels
- Skycoin is environmentally friendly and sustainable. Skycoin does not require twelve coal power plants to power mining

Current Status
--------------

* Peering infrastucture complete.  Clients are able to connect to each other, if configured appropriately.
* Blockchain v1 functioning
* Transactions are not networked.
* JSON RPC for wallet

Setup
-----

* Clone this repo
* Install [gvm](https://github.com/moovweb/gvm) (or hack up $GOPATH yourself)
* Install `go1.2` with gvm, or manually.  Skycoin does not work with earlier releases of go.
* `./compile/getdeps.sh` - This installs go dependencies.  Run this again if you get import errors in the future.
* Run `./compile/install-to-gopath.sh`.  This will symlink the repo directory into your `$GOPATH`.  You only need to do this if you did not clone the repo to `$GOPATH/src/github.com/skycoin/skycoin`, either via `go get` or manually.

Running
-------

### Command line 
Skycoin comes in three variations, `skycoin`, `skycoind` and `skycoindev`.  They differ in the command-line interface they expose and which values they default to.

To run skycoin with any one of these, do

```
go run cmd/$name/$name.go
```

For example, to run skycoindev,

```
go run cmd/skycoindev/skycoindev.go
``` 

For the developer's convenience,

```
./run.sh
```

will run `skycoindev`.


### GUI

To run the gui client, it must be built first.

```
./gui.sh build
```

Once it is built, you can run it with

```
./gui.sh
```

until the go source has changed and you need to rebuild.  
You do not need to rebuild if only modifying the GUI frontend code, 
located in `./static/`.

The GUI consists of a `node-webkit` binary and an `skycoin.nw` file which contains the frontend code and the skycoin binary.
When running the GUI, the `node-webkit` binary is executed, it unpacks the `skycoin.nw` file, forks skycoin which runs an http
server on `localhost:$randomport`, and the GUI's `index.html` is served from there.

If you are trying to run the skycoin GUI client on a platform that we are not targeting, you can run `node-webkit` with 
the `skycoin.nw` file produced by the build scripts in `compile/`.

Available Platforms
-------------------

The instructions for running the client apply to Linux, Windows and OSX.
Windows will need MingW.

Skycoin development is primarily done on Linux so Windows and OSX may break from time to time.

Please report any issues you have running skycoin on your system.

We will provide snapshot binary releases for Linux 32/64-bit, Windows 32-bit and OSX 32-bit once
the client is deemed ready for distribution.


Tests
-----

Skycoin tests can be run with 

```
./test.sh
```

### Test Roadmap

##### Libraries

* [pex](https://github.com/skycoin/pex) *Complete*
* [gnet](https://github.com/skycoin/gnet) *Complete*

##### Submodules

* coin
* daemon
* util
