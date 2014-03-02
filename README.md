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

Several emerging problems are throwing the long term survival of Bitcoin into doubt. Bitcoin was a proof-of-concept and an experiment. Satoshi never believed that it would reach as far as it has. Bitcoin was a seed of an idea and never intended as its final form. Skycoin is a coin which addresses the core challenges Bitcoin is facing.

Skycoin is inspired by Bitcoin but emphasizes
- usability
- security
- simplicity

Skycoin Improves on Bitcoin in several ways
- does not rely on miners for blockchain consensus
- easy to use wallet designed for the mainstream
- more secure than Bitcoin (51% attack immune, no double spending)
- 15 second transactions. Skycoin transactions are forty times faster than Bitcoin
- better name and branding
- clean, easy to use API for merchants and developers
- increased transaction privacy
- instant payments at point-of-sale with off-block chain transactions
- snapshot support. Devices will never need more than 1 GB of storage or RAM to run full Skycoin client

Security Improvements:
- immunity from Transaction Malleability attack
- no 51% attack and no double spending
- strict transaction propagation rules, increased DDoS resistance
- ECDSA signature compression support for smaller/faster transactions
- no hash collisions or duplicate coinbase outputs
- HTTPS encryption

Skycoin is the most secure and most usable cryptocoin that has ever been designed.

Distributed Consensus
---------------------

Bitcoin miners currently receive $34 in mining reward subsidies and $0.10 in fees, per Bitcoin transaction. Every Bitcoin user is paying the cost for mining. The cost of mining drags down the price of Bitcoin. We need alternatives to mining for securing blockchains.

Mathematical modeling of mining rewards suggests a difficulty peak and declining mining profitability by 2040. Without sufficient price deflation, transaction volume growth and increased transaction fees the difficulty peak leaves the resources securing the Bitcoin network at levels which Bitcoin open to attack and control by governments and corporations.

Reliance on mining is the primary existential threat facing Bitcoin over the next decade. The cost to disrupt the Bitcoin network stands at 50 million dollars in ASICs. The capactity to disrupt or shutdown the Bitcoin network is therefore within the capacity of governments and corporations.

If Bitcoin becomes a "threat", it will not survive technical attacks on the blockchain. The survivability of the Bitcoin network depends upon Bitcoin mining reaching a level of resource consumption and cost, beyond the reach of governments and finantial institions, while passing this cost onto the Bitcoin user base. 

Skycoin eliminates the cost of resource consumption for mining and replaces the dependence on miners with a distributed consensus algorithm called Obelisk. Skycoin achieves a higher level of blockchain security than Bitcoin, with lower electricity and hardware costs.

Key features of Obelisk
- first provably secure distributed consensus algorithm.
- extremely simple
- easy to model
- more ecologically friendly than mining

The Obelisk white paper draft will be released with the first Skycoin client.

No Mining
---------

Satoshi has stated that mining is Bitcoin's only flaw. Miners are the greatest threat to Bitcoin's survival. Mining introduces social factors and greed into the survivability of Bitcoin as a currency.

Problems with Mining:
- Dishonest miners are the greatest threat to Bitcoin
- the costs of mining are passed on to every Bitcoin user
- Bitcoin puts honest miners at a disadvantage to the mining cartels
- cex.io allows people to rent a 51% attack by the minute
- gigahash.io has abused its hashing power to steal Bitcoin from gambling sites
- pools have begun capping block size to drive up transaction fees
- Every month mining power is concentrated in the hands of fewer people
- mining pools can form cartels to orphan blocks by honest miners
- mining cartels can rent capacity to orphan blocks from other pools
- Mining ensures majority of new crypto-coins go to people with botnets, GPU farms and ASICs not available to the public
- miners and electricity companies profit. Everyone else pays the costs.
- block rewards create incentive to sybil attack network and influence block propagation. These attacks have slowed down blockchain downloads to unacceptable levels.
- electricity used by mining is wasted
- Mining increases transaction fees to unsustainable levels when block rewards decrease. Satoshi has stated that this will make Bitcoin unusable for microtransactions.
- Miners driving up transaction fees will make Bitcoin transaction fees higher than the banking system Bitcoin was created to replace.
- rapid changes in hash rate and mining profitability introduce unnecessary price volatility into Bitcoin. Price volatility caused by mining discourages Merchants from transacting in Bitcoin.
- Miners sell Bitcoins to pay equipment and electricity costs, driving down the price of Bitcoin. Every Bitcoin holder pays the cost of mining.
- If Bitcoin reaches viability as the next global reserve currency, the incentives for nation states to monopolize and control mining are too great.

Skycoin was designed to eliminate mining completely and eliminate the problems it creates
- Skycoin replaces reliance on miners with reliance on mathematics
- Skycoin is more secure because it does not rely upon the good will of miners
- Skycoin transactions will be cheaper because there are no mining costs passed on to users
- Skycoin transactions are not subject to 51% attacks by mining cartels
- Skycoin is environmentally friendly and sustainable. Skycoin does not use twelve coal power plants of electricity to be secure


Skycoin Project
---------------

Skycoin is more than a coin. Skycoin is at the core of number of projects that will be announced over the next year.

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