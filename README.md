skycoin
=======

Skycoin is a next-generation post proof-of-work cryptocurrency.

Installation
------------

*For detailed installation instructions, see [Installing Skycoin](../../wiki/Installation)*

```
./setup.sh
./run.sh -h
```

*Running Wallet

```
./run.sh -web-interface=true
Goto http://127.0.0.1:6402
```


Skycoin
-------

Skycoin is inspired by Bitcoin but emphasizes
- usability
- security
- simplicity

Skycoin Improves on Bitcoin in several ways
- easy to use wallet designed for the mainstream
- 15 second transactions. Skycoin transactions are forty times faster than Bitcoin
- better name and branding
- clean, easy to use API for merchants and developers
- increased transaction privacy
- native support for off-block chain transactions for instant payments at point-of-sale

Security Improvements:
- immunity from the Transaction Malleability attack
- immunity to the 51% attack
- ECDSA signature compression support for smaller/faster transactions
- no hash collisions or duplicate coinbase outputs
- HTTPS encryption

Skycoin is the most secure and most usable cryptocoin that has ever been designed.

No Mining
---------

Satoshi has stated that mining is Bitcoin's only flaw. Miners are the greatest threat to Bitcoin's survival. Mining introduces social factors and greed into the survivability of Bitcoin as a currency.

Problems with Mining:
- Dishonest miners are the greatest threat to Bitcoin
- Bitcoin puts honest miners at a disadvantage to the mining cartels
- cex.io allows people to rent a 51% attack by the minute
- gigahash.io has abused its hashing power to steal Bitcoin from gambling sites
- pools have begun capping block size to drive up transaction fees
- Every month mining power is concentrated in the hands of fewer people
- the costs of mining are passed on to every Bitcoin holder
- mining pools can form cartels to orphan blocks by honest miners
- mining cartels can rent capacity to orphan blocks from other pools
- Mining ensures majority of new crypto-coins go to people with botnets, GPU farms and ASICs not available to the public
- miners and electricity companies profit. Everyone else pays the costs.
- block rewards create incentive to sybil attack network and influence block propagation. These attacks have slowed down blockchain downloads to unacceptable levels.
- electricity used by mining is wasted
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

Distributed Consensus
---------------------

Skycoin uses a new distributed consensus algorithm called Obelisk to replace miners.

Key features of Obelisk
- first provably secure distributed consensus algorithm.
- extremely simple
- easy to model
- more ecologically friendly than mining

The Obelisk white paper draft will be released with the first Skycoin client.

Skycoin Project
---------------

Skycoin is more than a coin. Skycoin is at the core of number of projects that will be announced over the next year.

Current Status
--------------

Skycoin was announced in December with first Public Release scheduled by March. 

* Crypto is done
* Blockchain is done
* Networking is done
* Wallet is working with on going improvements

* Working on Obelisk
* Working on Darknet

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
* util *Complete*