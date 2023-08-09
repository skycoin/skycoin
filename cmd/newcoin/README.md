
[2021# Fiber Coin Creator CLI Documention
This tool can be used to create a new fiber coin easily from a config file.

## Usage

```
$ newcoin

NAME:
   newcoin - newcoin is a helper tool for creating new fiber coins

USAGE:
   newcoin [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
     createcoin  Create a new coin from a template file
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

#### Satisfy Dependencies

requires:
* `git`
* `go`

git clone the skycoin source code

```
mkdir -p $HOME/go/src/github.com/skycoin
cd $HOME/go/src/github.com/skycoin
git clone https://github.com/skycoin/skycoin
cd skycoin
```

NOTE: run all commands from within that directory NOTE: the default branch should be ‘develop’ - please make sure to use develop


### Create New Coin
When using the `newcoin` command, you should run it from the `$GOPATH/src/github.com/skycoin/skycoin` folder to utilise the built in default templates.

```bash
$ cd $GOPATH/src/github.com/skycoin/skycoin
$ go run cmd/newcoin/newcoin.go createcoin [command options]
```

```
OPTIONS:
   --coin value                             name of the coin to create (default: "skycoin")
   --template-dir value, --td value         template directory path (default: "./template")
   --coin-template-file value, --ct value   coin template file (default: "coin.template")
   --visor-template-file value, --vt value  visor template file (default: "visor_parameters.template")
   --config-dir value, --cd value           config directory path (default: "./")
   --config-file value, --cf value          config file path (default: "fiber.toml")
```

#### Example
Create a test coin using application defaults.

```bash
$ cd $GOPATH/src/github.com/skycoin/skycoin
$ go run cmd/newcoin/newcoin.go --coin testcoin
```

This will create a new directory, `testcoin`, in `cmd` folder and a `testcoin.go` file inside that folder.
It will also use the built-in defaul options (specified above) and draw template configuration from `$GOPATH/src/github.com/skycoin/skycoin/template`

This file can be used to run a "testcoin" node.


#### Coin Templates

A coin template is used by the newcoin command to generate a Go source file that is used to run a fibercoin. These templates are located in the `template/` directory.

The coin.template file is used to generate the `cmd/mycoin/mycoin.go` source file, while the file params.template is used to generate the file `src/params/params.go` source file. The former is used to run the peer and publisher nodes for the fibercoin, while the latter is used to hold multiple configuration parameters for the fibercoin nodes.
Genesis Address and Genesis Private and Public Keys

In order to initialize a new fibercoin, secret and public keys need to be generated to create the genesis transaction. Generating these keys is achieved by running the following command:

```
go run cmd/skycoin-cli/skycoin-cli.go addressGen
```

The output of this command will be similar to the one below:
```

{
    "meta": {
        "coin": "mycoin",
        "cryptoType": "",
        "encrypted": "false",
        "filename": "2023_08_08_8b30.wlt",
        "label": "",
        "lastSeed": "13017d89a3a23193107709d06e81db4d3c787ab1044fd15189424efe8cdd128a",
        "secrets": "",
        "seed": "marble face march turtle essence motion expand lift honey hole bronze test",
        "seedPassphrase": "",
        "tm": "1691498129",
        "type": "deterministic",
        "version": "0.4",
        "xpub": ""
    },
    "entries": [
        {
            "address": "mhypoFwrNE4woCAfvBr4JhmY3pw7u746hr",
            "public_key": "032e123ad2d33f3bad388cd914bd250580c31803d0cff822b9b88ff4ce204a1acc",
            "secret_key": "926abeadce6e4283dae8782af2b1e7549267d29d8d7517d48a4d9b0d9740e22b"
        }
    ]
}

```

The bits of interest from this output are the values of the JSON keys address, public_key and secret_key.

These values are used for editing the file fiber.toml, with the exception of secret_key. At the moment, the modification of this file needs to be done manually, but this process should be performed automatically in later versions of newcoin. The value of the secret key must be kept secret, as the name implies, as this key could be used to sign transactions by anyone who posseses it.

#### fiber.toml Configuration File

fiber.toml is used to set parameters that are used during the initialization and operation of a fibercoin. The file already contains some values that can be considered as default, such as the genesis_timestamp or max_block_size, but other fields need to be set up with different values for every fibercoin. The following fields need to be updated:

* blockchain_pubkey_str
* genesis_address_str
* genesis_signature_str

 The values of the first two fields are updated with the values obtained by following the instructions in the section Genesis Transaction, while the last one is automatically generated and added to fiber.toml by initializing a blockchain

Other fields that can be of interest in this file are:

* create_block_max_transaction_size
* max_block_size
* unconfirmed_max_transaction_size.

These fields control how large a fibercoin's transactions can be. The default is set to be 5 Mb for all of these parameters.

Lastly, any field related to the configuration of a cryptocurrency can be changed to alter the parameters of the fibercoin blockchain:

*    genesis_coin_volume
*    create_block_burn_factor
*    unconfirmed_burn_factor
*    max_coin_supply
*    user_burn_factor

#### Initializing a fibercoin

In order to initialize a new fibercoin, `newcoin` needs to create the `mycoin` command (located in `cmd/mycoin`) using the parameters defined in `./fiber.toml`.

The workflow is as follows:

* `newcoin` is run in order to create `mycoin`

```
go run cmd/newcoin/newcoin.go createcoin --coin mycoin
```

* Create the genesis address and keys
```
skycoin-cli addressGen > genesis.json
```
the output is saved to a file called `genesis.json` (later delete this file or save offline!)

Note: when troubleshooting or attempting this process multiple times in succession, start again with the next step

* `mycoin` is run to initialize the blockchain. Read carefully:

`$KEY` envs are obtained from `genesis.json` which was created in the previous step

you can set them by using `export`
```
export SEC_KEY=<substituite-the-blockchain-secret-key-string-here>
export GEN_ADD=<substituite-the-genesis-address-here>
export PUB_KEY=<substituite-the-blockchain-public-key-here>
```

Or simply copy them into place of the following command
```
go run cmd/mycoin/mycoin.go --block-publisher=true --blockchain-secret-key=$SEC_KEY --blockchain-public-key=$PUB_KEY --genesis-address=$GEN_ADD
```

* run the above _for a few moments_ until you start getting `“ERROR”` messages;  then `ctrl+c` to stop the process.

The output from the above command:
```
$ go run cmd/mycoin/mycoin.go --block-publisher=true --blockchain-secret-key=$SEC_KEY --blockchain-public-key=$PUB_KEY --genesis-address=$GEN_ADD
[2021-02-19T11:59:21-06:00] INFO [main]: App version: 0.27.1
[2021-02-19T11:59:21-06:00] INFO [main]: OS: linux
[2021-02-19T11:59:21-06:00] INFO [main]: Arch: amd64
[2021-02-19T11:59:21-06:00] INFO [main]: Opening database /home/user/.mycoin/data.db
[2021-02-19T11:59:21-06:00] INFO [main]: DB version: 0.27.1
[2021-02-19T11:59:21-06:00] INFO [main]: Coinhour burn factor for user transactions is 10
[2021-02-19T11:59:21-06:00] INFO [main]: Max transaction size for user transactions is 32768
[2021-02-19T11:59:21-06:00] INFO [main]: Max decimals for user transactions is 3
[2021-02-19T11:59:21-06:00] INFO [main]: wallet.NewService
[2021-02-19T11:59:21-06:00] INFO [main]: visor.New
[2021-02-19T11:59:21-06:00] INFO [visor]: Creating new visor
[2021-02-19T11:59:21-06:00] INFO [visor]: Visor running in block publisher mode
[2021-02-19T11:59:21-06:00] INFO [visor]: Coinhour burn factor for unconfirmed transactions is 10
[2021-02-19T11:59:21-06:00] INFO [visor]: Max transaction size for unconfirmed transactions is 32768
[2021-02-19T11:59:21-06:00] INFO [visor]: Max decimals for unconfirmed transactions is 3
[2021-02-19T11:59:21-06:00] INFO [visor]: Coinhour burn factor for transactions when creating blocks is 10
[2021-02-19T11:59:21-06:00] INFO [visor]: Max transaction size for transactions when creating blocks is 32768
[2021-02-19T11:59:21-06:00] INFO [visor]: Max decimals for transactions when creating blocks is 3
[2021-02-19T11:59:21-06:00] INFO [visor]: Max block size is 32768
[2021-02-19T11:59:21-06:00] INFO [blockdb]: Unspents.MaybeBuildIndexes
[2021-02-19T11:59:21-06:00] INFO [blockdb]: Rebuilding unspent_pool_addr_index (addrHeightIndexExists=false, addrIndexHeight=0, headSeq=0)
[2021-02-19T11:59:21-06:00] INFO [blockdb]: Building unspent address index
[2021-02-19T11:59:21-06:00] INFO [blockdb]: No unspents to index
[2021-02-19T11:59:21-06:00] INFO [visor]: Visor initHistory
[2021-02-19T11:59:21-06:00] INFO [visor]: Resetting historyDB
[2021-02-19T11:59:21-06:00] INFO [visor]: Visor parseHistoryTo
[2021-02-19T11:59:21-06:00] INFO [visor]: Unconfirmed transaction pool size: 0
[2021-02-19T11:59:21-06:00] INFO [main]: daemon.New
[2021-02-19T11:59:21-06:00] INFO [main]: kvstorage.NewManager
[2021-02-19T11:59:21-06:00] INFO [kvstorage]: Creating new KVStorage manager
[2021-02-19T11:59:21-06:00] INFO [kvstorage]: KVStorage is disabled
[2021-02-19T11:59:21-06:00] INFO [main]: api.NewGateway
[2021-02-19T11:59:21-06:00] WARN [api]: HTTPS not in use!
[2021-02-19T11:59:21-06:00] INFO [pex]: Trying to download peers list url="https://downloads.mycoin.net/blockchain/peers.txt"
[2021-02-19T11:59:21-06:00] INFO [main:CRITICAL]: Full address: http://127.0.0.1:6419
[2021-02-19T11:59:21-06:00] INFO [main]: visor.Init
[2021-02-19T11:59:21-06:00] INFO [visor]: Visor init
[2021-02-19T11:59:21-06:00] INFO [visor]: Visor maybeCreateGenesisBlock
[2021-02-19T11:59:21-06:00] INFO [visor]: Create genesis block
[2021-02-19T11:59:21-06:00] INFO [visor]: Genesis block signature=86f410b1c507e75cd9e05bfecb735c8fad366579afdd714ed0f27730369768f262ed21940c814dd522e6ce9535cebd95e05dbfbb1eb48fcd0dd91d7240f7ec1401
[2021-02-19T11:59:21-06:00] INFO [visor]: Removed 0 invalid txns from pool
[2021-02-19T11:59:21-06:00] INFO [main]: webInterface.Serve
[2021-02-19T11:59:21-06:00] INFO [api]: Starting web interface on 127.0.0.1:6419
[2021-02-19T11:59:21-06:00] INFO [main]: daemon.Run
[2021-02-19T11:59:21-06:00] INFO [daemon]: Daemon UserAgent is mycoin:0.27.1
[2021-02-19T11:59:21-06:00] INFO [daemon]: Daemon unconfirmed BurnFactor is 10
[2021-02-19T11:59:21-06:00] INFO [daemon]: Daemon unconfirmed MaxTransactionSize is 32768
[2021-02-19T11:59:21-06:00] INFO [daemon]: Daemon unconfirmed MaxDropletPrecision is 3
[2021-02-19T11:59:21-06:00] INFO [pex]: Pex.Run started
[2021-02-19T11:59:21-06:00] INFO [daemon]: daemon.Pool listening on port 6001
[2021-02-19T11:59:21-06:00] INFO [gnet]: Listening for connections on :6001...
[2021-02-19T11:59:21-06:00] WARN [daemon]: maybeConnectToTrustedPeer: connectToPeer failed addr="70.121.6.216:6001" error="Already connected to this peer"
[2021-02-19T11:59:21-06:00] ERROR [daemon:CRITICAL]: maybeConnectToTrustedPeer error="Could not connect to any trusted peer"
[2021-02-19T11:59:21-06:00] ERROR [pex]: waiting to retry downloadText error="Get "https://downloads.mycoin.net/blockchain/peers.txt": dial tcp: lookup downloads.mycoin.net: no such host" waitTime=719.600155ms
[2021-02-19T11:59:22-06:00] INFO [pex]: Trying to download peers list url="https://downloads.mycoin.net/blockchain/peers.txt"
[2021-02-19T11:59:22-06:00] ERROR [pex]: waiting to retry downloadText error="Get "https://downloads.mycoin.net/blockchain/peers.txt": dial tcp: lookup downloads.mycoin.net: no such host" waitTime=436.089288ms
[2021-02-19T11:59:23-06:00] INFO [pex]: Trying to download peers list url="https://downloads.mycoin.net/blockchain/peers.txt"
[2021-02-19T11:59:23-06:00] ERROR [pex]: waiting to retry downloadText error="Get "https://downloads.mycoin.net/blockchain/peers.txt": dial tcp: lookup downloads.mycoin.net: no such host" waitTime=966.087607ms
[2021-02-19T11:59:23-06:00] INFO [pex]: Trying to download peers list url="https://downloads.mycoin.net/blockchain/peers.txt"
[2021-02-19T11:59:23-06:00] ERROR [pex]: waiting to retry downloadText error="Get "https://downloads.mycoin.net/blockchain/peers.txt": dial tcp: lookup downloads.mycoin.net: no such host" waitTime=1.534681399s
[2021-02-19T11:59:25-06:00] INFO [pex]: Trying to download peers list url="https://downloads.mycoin.net/blockchain/peers.txt"
[2021-02-19T11:59:25-06:00] ERROR [pex]: waiting to retry downloadText error="Get "https://downloads.mycoin.net/blockchain/peers.txt": dial tcp: lookup downloads.mycoin.net: no such host" waitTime=1.977973304s
^C[2021-02-19T11:59:26-06:00] INFO [main]: Shutting down...
[2021-02-19T11:59:26-06:00] INFO [main]: Closing web interface
[2021-02-19T11:59:26-06:00] INFO [api]: Shutting down web interface
[2021-02-19T11:59:26-06:00] INFO [api]: Web interface closed
[2021-02-19T11:59:26-06:00] ERROR [main]: webInterface.Serve failed error="accept tcp 127.0.0.1:6419: use of closed network connection"
[2021-02-19T11:59:26-06:00] INFO [api]: Web interface shut down
[2021-02-19T11:59:26-06:00] INFO [main]: Closing daemon
[2021-02-19T11:59:26-06:00] INFO [daemon]: Stopping the daemon run loop
[2021-02-19T11:59:26-06:00] INFO [daemon]: Shutting down Pool
[2021-02-19T11:59:26-06:00] INFO [gnet]: ConnectionPool.Shutdown called
[2021-02-19T11:59:26-06:00] INFO [gnet]: ConnectionPool.Shutdown closed pool.quit
[2021-02-19T11:59:26-06:00] INFO [daemon]: Daemon closed
[2021-02-19T11:59:26-06:00] INFO [gnet]: ConnectionPool.Shutdown waiting for strandDone
[2021-02-19T11:59:26-06:00] INFO [gnet]: ConnectionPool.Shutdown closing the listener
[2021-02-19T11:59:26-06:00] INFO [gnet]: Connection pool closed
[2021-02-19T11:59:26-06:00] INFO [gnet]: ConnectionPool.Shutdown disconnecting all connections
[2021-02-19T11:59:26-06:00] INFO [gnet]: ConnectionPool.Shutdown waiting for done
[2021-02-19T11:59:26-06:00] INFO [daemon]: Shutting down Pex
[2021-02-19T11:59:26-06:00] INFO [pex]: Shutting down pex
[2021-02-19T11:59:26-06:00] INFO [pex]: Save peerlist
[2021-02-19T11:59:26-06:00] INFO [pex]: Pex.Run stopped
[2021-02-19T11:59:26-06:00] INFO [pex]: Pex shutdown
[2021-02-19T11:59:26-06:00] INFO [daemon]: Daemon shutdown complete
[2021-02-19T11:59:26-06:00] INFO [main]: Waiting for goroutines to finish
[2021-02-19T11:59:26-06:00] INFO [main]: Closing database
[2021-02-19T11:59:26-06:00] INFO [main]: Goodbye
```

Now, about halfway along in the above command output, the genesis signature of the blockchain was generated.

__Copy the genesis signature string__

* add the genesis signature from above, and the address and pubkey fields from genesis.json to fiber.toml

__NOTE: if you add the secret key to fiber.toml it will show up in the help menu of the created .go file or binary. Don’t do that for a distributable binary. Run the wallet in publisher mode and explicitly specify the secret key string for a block publisher as an env and clear your terminal history__

__note: the timestamp cannot be blank. Choose a recent time.__


* use `newcoin` to regenerate the `cmd/mycoin/mycoin.go` file

We could have actually used the skycoin wallet binary from the start, howeverit's more complex to do that and it may present a risk to the local copy of the existing skycoin blockchain in `~/.skycoin`

Also, it's not possible to change the display name, etc. with just flags. But it will work still.

The new blockchain will have hence (in this example) been created in `~/.mycoin`

```
skycoin-newcoin createcoin --coin mycoin
```

* Check the help menu of the created mycoin.go to see the defaults such as genesis address, signature, etc.
```
go run cmd/mycoin/mycoin.go --help
```
__if this is intended to produce a distributable binary, the blockchain secret key must not be included!__

Remember, your blockchain currently only exists locally. `~/.mycoin` will need to be present in the production deployment environment.

#### DistributeGenesis

The exact process for this is pending improvement

* Run a block publisher node
```
go run cmd/mycoin/mycoin.go --block-publisher=true --blockchain-secret-key=$SECRET_KEY
```

* Run a peer node

To run a peer node on the same machine as the publisher, change the port and data directory
```
go run cmd/tesla/tesla.go -launch-browser=true -enable-all-api-sets=true -enable-gui=true -log-level=debug --port=5998 -data-dir=$HOME/.tesla1 -web-interface-port=6418
```

at this point, one can optionally go into the web interface of the wallet and load the genesis wallet seed. The entire coin volume should be present in the wallet.

* Create the first transaction using `skycoin-cli`

Export the necessary environmental variables to make skycoin-cli work with your chain
```
export RPC_ADDR="http://127.0.0.1:6418"
export COIN="mycoin"
export DATA_DIR=$HOME/.mycoin1/"
```

At this point it would be typical to use:

```
skycoin-cli distributeGenesis $SEC_KEY
```

This currently does not work.
The genesis address creates a transaction to itself, showing that coins were sent and coins were received; but the coins simply vanish.

A workaround at this point is to create a raw transaction:

```
skycoin-cli createRawTransaction '/home/user/.mycoin1/wallets/genesis.wlt' $GEN_ADD 100000000 -a $GEN_ADD
```

#### TROUBLESHOOTING

re-initialize the blockchain:

* stop any running instances and remove everything except wallets dir from ~/.mycoin & ~/.mycoin1

* restart the publisher briefly, then kill it

* copy the signature to fiber.toml

* update `cmd/mycoin/mycoin.go` file with `newcoin`

* restart the publisher

* restart the peer

#### Deployment Considerations - Network

A typical production deployment for a fibercoin, following the example of skycoin, has the following subdomains and endpoints:

* node.skycoin.com - mobile wallet node
* explorer.skycoin.com - blockchain explorer
* downloads.skycoin.com/blockchain/peers.txt - peers.txt url

It is necessary to forward the port (referenced in peers.txt) to the machine hosting the node endpoint.

LAN IP addresses can be set initially for testing purposes.

#### PRODUCTION DEPLOYMENT

In this example, a publisher node is run on port 6418, the node for the mobile wallets is run on port 6419 (the default port) and any other instance for testing purposes, is run on the next lower port. Specify the web interface port similarly to avoid conflicts.

A copy of the data dir should be made and specified with the -data-dir flag for use with each additional instance.

An HTTP server such as caddy-server may be used to reverse proxy port 6419 to the desired subdomain, for example node.mycoin.net, with the following lines in a Caddyfile:
```
node.mycoin.net { reverse_proxy 127.0.0.1:6419 }
```
The node used by mobile wallets should be run with the following flags specified:

```
-enable-all-api-sets=true -log-level=debug -disable-csrf -host-whitelist node.mycoin.net
```
The subdomain can now be tested with the mobile wallet, by adding node.mycoin.net as the nodeURL in the settings page of the skycoin mobile wallet.

Running the skycoin explorer with your fibercoin is as simple as changing two lines in explorer.go.

The explorer.go file should similarly be changed to increment the port it serves on, and change the port it uses to connect to skycoin / fibercoins to the one you are using. The port which is served by the explorer can then be reverse-proxied to another subdomain such as:
```
explorer.mycoin.net { reverse_proxy 127.0.0.1:8003 }
```
