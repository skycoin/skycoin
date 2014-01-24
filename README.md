skycoin
=======

Skycoin is a next-generation post proof-of-work cryptocurrency.

It is currently in `alpha`.

No Mining
---------

Satoshi has stated that mining is Bitcoin's only flaw. Miners are the greatest threat to Bitcoin's survival. Mining introduces human factors and greed into the survivability of Bitcoin as a currency.

Problems with Mining:
- Dishonest miners are the greatest threat to Bitcoin
- gigahash.io has abused its hashing power to steal Bitcoin from gambling sites
- Every month mining power is concentrated in the hands of fewer people
- cex.io allows people to rent a 51% attack by the minute
- mining pools can form cartels to orphan blocks by honest miners
- Bitcoin puts honest miners at a disadvantage to the mining cartels
- mining cartels will rent capacity to orphan blocks from other pools, if they can
- pools have begun capping block size to drive up transaction fees
- Mining ensures majority of new crypto-coins go to people with botnets, GPU farms and ASICs not available to the public
- miners and electricity companies profit. Everyone else pays the costs.
- block rewards create incentive to sybil attack network and influence block propagation. These attacks have slowed down blockchain downloads to unacceptable levels.
- electricity used by mining is wasted
- the costs of mining is passed on to every Bitcoin holder
- Mining increases transaction fees to unsustainable levels when block rewards decrease. Satoshi has stated that this will make Bitcoin unusable for microtransactions.
- Miners driving up transaction fees will make Bitcoin transaction fees higher than the banking system Bitcoin was created to replace.
- rapid changes in hash rate and mining profitability introduce unnecessary price volatility into Bitcoin. Price volatility caused by mining discourages Merchants from transacting in Bitcoin.
- If Bitcoin reaches viability as the next global reserve currency, the incentives for nation states to monopolize and control mining are too great.
- Miners sell Bitcoins to pay equipment and electricity costs, driving down the price of Bitcoin. Every Bitcoin holder pays the cost of mining.

Skycoin was designed to eliminate mining completely and eliminate the problems it creates
- Skycoin replaces reliance on miners with reliance on mathematics
- Skycoin is more secure because it does not rely upon the good will of miners
- Skycoin transactions will be cheaper because there are no mining costs passed on to users
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

```
git clone https://github.com/skycoin/skycoin
./install.sh
```

If a library could not be found, the dependencies may need to be updated.

```
./compile/get-dependencies.sh
```

Running
-------

Make sure to follow the instructions in Setup before attempting these commands.

### Daemon

```
./run.sh
```

This runs the developer's version of the skycoin daemon.  To run
the version intended for release,

```
go run cmd/skycoind/skycoind.go
```

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
* daemon *Complete*
* util
