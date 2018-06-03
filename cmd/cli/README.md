# CLI Documentation

Skycoin command line interface

The CLI command APIs can be used directly from a Go application, see [Skycoin CLI Godoc](https://godoc.org/github.com/skycoin/skycoin/src/cli).

<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3" -->

- [Install](#install)
    - [Enable command autocomplete](#enable-command-autocomplete)
- [Environment Setting](#environment-setting)
    - [RPC_ADDR](#rpc_addr)
    - [WALLET_DIR](#wallet_dir)
    - [WALLET_NAME](#wallet_name)
    - [USE_CSRF](#use_csrf)
- [Usage](#usage)
    - [Add Private Key](#add-private-key)
    - [Check address balance](#check-address-balance)
    - [Generate new addresses](#generate-new-addresses)
    - [Check address outputs](#check-address-outputs)
    - [Check block data](#check-block-data)
    - [Check database integrity](#check-database-integrity)
    - [Create a raw transaction](#create-a-raw-transaction)
    - [Decode a raw transaction](#decode-a-raw-transaction)
    - [Broadcast a raw transaction](#broadcast-a-raw-transaction)
    - [Generate a wallet](#generate-a-wallet)
    - [Generate addresses for a wallet](#generate-addresses-for-a-wallet)
    - [Last blocks](#last-blocks)
    - [List wallet addresses](#list-wallet-addresses)
    - [List wallets](#list-wallets)
    - [Send](#send)
    - [Show Config](#show-config)
    - [Status](#status)
    - [Get transaction](#get-transaction)
    - [Verify address](#verify-address)
    - [Check wallet balance](#check-wallet-balance)
    - [See wallet directory](#see-wallet-directory)
    - [List wallet transaction history](#list-wallet-transaction-history)
    - [List wallet outputs](#list-wallet-outputs)
    - [CLI version](#cli-version)
- [Note](#note)

<!-- /MarkdownTOC -->


## Install

```bash
$ cd $GOPATH/src/github.com/skycoin/skycoin/cmd/cli
$ ./install.sh
```

### Enable command autocomplete

If you are in `bash`, run the following command:

```bash
$ PROG=skycoin-cli source $GOPATH/src/github.com/skycoin/skycoin/cmd/cli/autocomplete/bash_autocomplete
```

If you are in `zsh`, please replace the `bash_autocomplete` with `zsh_autocomplete` in the previous command.

To avoid run the command everytime when you start a new terminal session, you can copy the script into
you `~/.bashrc` or `~/.zshrc` file.

## Environment Setting

The CLI uses environment variable to manage the configurations.

### RPC_ADDR

CLI will connect to skycoin node rpc address `http://127.0.0.1:6420` by default.
You can change the address by setting the `RPC_ADDR` environment variable
with the following command:

```bash
$ export RPC_ADDR=http://127.0.0.1:6420
```

Note: `RPC_ADDR` must be in `scheme://host` format.

### WALLET_DIR

The default CLI wallet dir is located in `$HOME/.skycoin/wallets/`, change it by setting the
`WALLET_DIR` environment variable.

```bash
$ export WALLET_DIR=$HOME/YOUR_WALLET_DIR
```

### WALLET_NAME

The default CLI wallet file name is `skycoin_cli.wlt`, change it by setting the `WALLET_NAME` env.
The wallet file name must have `.wlt` extension.

```bash
$ export WALLET_NAME=YOUR_WALLET_NAME
```

### USE_CSRF

If the remote node to communicate with is not run with `-csrf-disabled`, set this variable.
CSRF is enabled by default on nodes.

```bash
$ export USE_CSRF=1
```

## Usage

After the installation, you can run `skycoin-cli` to see the usage:

```
$ skycoin-cli

NAME:
   skycoin-cli - the skycoin command line interface

USAGE:
   skycoin-cli [global options] command [command options] [arguments...]

VERSION:
   0.23.1-rc2

COMMANDS:
     addPrivateKey         Add a private key to specific wallet
     addressBalance        Check the balance of specific addresses
     addressGen            Generate skycoin or bitcoin addresses
     addressOutputs        Display outputs of specific addresses
     blocks                Lists the content of a single block or a range of blocks
     broadcastTransaction  Broadcast a raw transaction to the network
     checkdb               Verify the database
     createRawTransaction  Create a raw transaction to be broadcast to the network later
     decodeRawTransaction  Decode raw transaction
     generateAddresses     Generate additional addresses for a wallet
     generateWallet        Generate a new wallet
     lastBlocks            Displays the content of the most recently N generated blocks
     listAddresses         Lists all addresses in a given wallet
     listWallets           Lists all wallets stored in the wallet directory
     send                  Send skycoin from a wallet or an address to a recipient address
     showConfig            show cli configuration
     status                Check the status of current skycoin node
     transaction           Show detail info of specific transaction
     verifyAddress         Verify a skycoin address
     version
     walletBalance         Check the balance of a wallet
     walletDir             Displays wallet folder address
     walletHistory         Display the transaction history of specific wallet. Requires skycoin node rpc.
     walletOutputs         Display outputs of specific wallet
     help, h               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
ENVIRONMENT VARIABLES:
    RPC_ADDR: Address of RPC node. Must be in scheme://host format. Default "http://127.0.0.1:6420"
    COIN: Name of the coin. Default "skycoin"
    USE_CSRF: Set to 1 or true if the remote node has CSRF enabled. Default false (unset)
    WALLET_DIR: Directory where wallets are stored. This value is overriden by any subcommand flag specifying a wallet filename, if that filename includes a path. Default "$HOME/.$COIN/wallets"
    WALLET_NAME: Name of wallet file (without path). This value is overriden by any subcommand flag specifying a wallet filename. Default "$COIN_cli.wlt"
```

### Add Private Key
Add a private key to a skycoin wallet.

```bash
$ skycoin-cli addPrivateKey [command options] [private key]
```

```
OPTIONS:
    -f value [wallet file or path] private key will be added to this wallet
    if not specified then default wallet ($HOME/.skycoin/wallets//wallets/skycoin_cli.wlt)
    will be used
```

#### Example
```bash
$ skycoin-cli addPrivateKey -f $WALLET_PATH $PRIVATE_KEY
```

```
$ success
```

### Check address balance
Check balance of specific addresses, join multiple addresses with space.

```bash
$ skycoin-cli addressBalance [addresses]
```

#### Example
```bash
$ skycoin-cli addressBalance 2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc 2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv
```
<details>
 <summary>View Output</summary>

```json
{
 "confirmed": {
     "coins": "324951.932000",
     "hours": "166600293"
 },
 "spendable": {
     "coins": "324951.932000",
     "hours": "166600293"
 },
 "expected": {
     "coins": "324951.932000",
     "hours": "166600293"
 },
 "addresses": [
     {
         "confirmed": {
             "coins": "2.000000",
             "hours": "1158"
         },
         "spendable": {
             "coins": "2.000000",
             "hours": "1158"
         },
         "expected": {
             "coins": "2.000000",
             "hours": "1158"
         },
         "address": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc"
     },
     {
         "confirmed": {
             "coins": "324949.932000",
             "hours": "166599135"
         },
         "spendable": {
             "coins": "324949.932000",
             "hours": "166599135"
         },
         "expected": {
             "coins": "324949.932000",
             "hours": "166599135"
         },
         "address": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"
     }
 ]
}
```
</details>

### Generate new addresses
Generate new skycoin or bitcoin addresses.

```bash
$ skycoin-cli addressGen [command options] [arguments...]
```

```
OPTIONS:
        --count value, -c value  Number of addresses to generate (default: 1)
        --hide-secret, -s        Hide the secret key from the output
        --bitcoin, -b            Output the addresses as bitcoin addresses instead of skycoin addresses
        --hex, -x                Use hex(sha256sum(rand(1024))) (CSPRNG-generated) as the seed if no seed is not provided
        --only-addr, --oa        Only show generated address list, hide seed, secret key and public key
        --seed value             Seed for deterministic key generation. If `--hex` is not defined will use bip39 to generate a seed if not provided.
```

#### Examples
##### Generate `n` number of skycoin addresses
```bash
$ skycoin-cli addressGen --count 2
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "skycoin",
     "seed": "genius canyon asset swallow picture torch awkward radar nest bunker walnut garage"
 },
 "entries": [
     {
         "address": "2KC8vSgNXVZ6GMYY8edUqvRsbupZVwWKETe",
         "public_key": "0215682c10f6293bf52c543787613e898f4e4af70aa87eac8848b0535d8b8c05f1",
         "secret_key": "31d6f561dad49f2c60c02a97395a2f3df67bb9e092806356ddbb952556c96e82"
     },
     {
         "address": "yzkv7v2T4fbQmZKdiLcq8tAHEVrVbrvGvh",
         "public_key": "02dc8409077376bc8a834185739133f47805764f061103c88a4b5b0d2809b310b7",
         "secret_key": "3ba5855ad3b1ec7e02918d5329dc5425690a93331514370d739f58556236c1ce"
     }
 ]
}
```
</details>


##### Generate `n` number of bitcoin addresses
```bash
$ skycoin-cli addressGen --count 2 --bitcoin
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "bitcoin",
     "seed": "sun resemble crisp clerk square initial like urge abstract picnic impulse robot"
 },
 "entries": [
     {
         "address": "15FW6YEZuwr68h54DrXD3Tv1Cc1aKHoynF",
         "public_key": "024529a5b1f1c110dd122598052115717a7a042b4831acbf234fe01968f82d1312",
         "secret_key": "L3crKysGdHmKQ2j27wfgew4timWiUrPxwUi8FNE75S872C1K9pns"
     },
     {
         "address": "1EBKC7u29ea1jPtEoC5LLMcXwpZBTmFxhs",
         "public_key": "03faaff073f752cafccb3f639b2174e6c48b04a24cbdefafbfdbda2f54ba5e81a9",
         "secret_key": "KxomKUvagGTviuxAr9HNRfXEaim4evvFJVGmuk2LYA5ZLXznvX6k"
     }
 ]
}
```
</details>

##### Hide secret in output
```bash
$ skycoin-cli addressGen --count 2 --hide-secret
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "skycoin",
     "seed": "walnut wise pluck sniff weird enable document special soul era mercy you"
 },
 "entries": [
     {
         "address": "27ohsY7Hg5dEDySUg17gStEQRzLFxE8mVrU",
         "public_key": "02ead2834f41f91dd3847924f6257b2b852390708cd2c91db80f017fd21f9768af",
         "secret_key": ""
     },
     {
         "address": "2FHuME9U7CEN3vWwwRzJAgP4K2JKPfoRxzp",
         "public_key": "03af027c379b380c009cbefc3b251e7b42af9753125a8d9ef0a50249e97060c673",
         "secret_key": ""
     }
 ]
}
```
</details>

##### Output only an address list
```bash
$ skycoin-cli addressGen --count 2 --only-addr
```

```
7HVmKni3ggMdtseynSkNkqoCnsH7vkS6cg
2j5QSbHgLWXA2qXZvLzJHRo6Cissxer4CSt
```

> NOTE: If no seed is provided with the `--seed flag` and `--hex` flag is not used then bip39 is used to generate a seed

##### Use a predefined seed value
```bash
$ skycoin-cli addressGen --count 2 --seed "my super secret seed"
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "skycoin",
     "seed": "my super secret seed"
 },
 "entries": [
     {
         "address": "NMwevFV1NhPSp8H4QHPQyDuKCkLjsQ5zRb",
         "public_key": "03a0571ef3ac621aa1fe950753110187bd096a32cc791f227406badbc21676743c",
         "secret_key": "0977909e18ef6b3dbc522e79c26c3113ab6d4ae7a9f4a848dcd49e5b4634a77c"
     },
     {
         "address": "a1ec9zaxj5ndhteyJeocdhYFLHfvm86UPG",
         "public_key": "021990611d33bdc3ca70da07b1e4f8a1928a6cb369fb343d4b2ce0c0b123273387",
         "secret_key": "e08dd4de4920edc1ae5aa2260167657e64a5ff146b90d21fb1a39294c94c940c"
     }
 ]
}
```
</details>

##### Generate addresses with a hex (CSPRNG-generated) seed
```bash
skycoin-cli addressGen --count 2 --hex
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "skycoin",
     "seed": "d5fa95cc3bd265c9ef99e7c2b300f0ede75375fbb76b2329bd5877631c315068"
 },
 "entries": [
     {
         "address": "2URMzEQ2A1xuf3GGN6Tr8tucCzJKYVpj9Fz",
         "public_key": "03d5e38c96829dbc8b873822ba29ebb7cf5c8c32c386348f032d523f0640c31640",
         "secret_key": "e0eccbd416b9fb8a109426681e890362a24491072edd8275a29b1878874fc9b4"
     },
     {
         "address": "2Qct8BmaCvMfPUnMCtTj5sMLNNoLbshAxoe",
         "public_key": "0241404173e29d9ae4a16c6496baff72cfd94fca16c727b7b1192bdeb736ced278",
         "secret_key": "972c0596a442d495fda1bba055df6334aa0121376248f19194ddc602368bda30"
     }
 ]
}
```
</details>


### Check address outputs
Display outputs of specific addresses, join multiple addresses with space.

```bash
$ skycoin-cli addressOutputs [address list]
```

#### Example
```bash
skycoin-cli addressOutputs tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V 29fDBQuJs2MDLymJsjyWH6rDjsyv995SrGU
```

<details>
 <summary>View Output</summary>

```json
{
 "outputs": {
     "head_outputs": [
         {
             "hash": "f5f838edf75b68882cacb7fa071538bcf800515d5a7f42e3a8c5e6d681879a82",
             "time": 1522603686,
             "block_seq": 20256,
             "src_tx": "cd0725e9cfc23cfed279aeda70b765238d0cc282406c48f811c6ad2874593f03",
             "address": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
             "coins": "0.500000",
             "hours": 43,
             "calculated_hours": 122
         },
         {
             "hash": "b2182ec3481f7af9884e4839016a145c07b66fce68c3b9ff04d897d1f1db5717",
             "time": 1522603586,
             "block_seq": 20255,
             "src_tx": "48b385567796725212ed8195a9437b15d5cd82186205b9d8fd027fa75f98060e",
             "address": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
             "coins": "0.500000",
             "hours": 173,
             "calculated_hours": 252
         },
         {
             "hash": "86ba8131756b4db5e163d38aa119ff07e4bd3fe05bbf3c28cef8471648d77080",
             "time": 1517145975,
             "block_seq": 12675,
             "src_tx": "ad191f910e5508e0b0e0ab24ba815e784a1a2b63ca21043e7746bebf25106742",
             "address": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
             "coins": "1.000000",
             "hours": 1232,
             "calculated_hours": 2906
         }
     ],
     "outgoing_outputs": [],
     "incoming_outputs": []
 }
}
```
</details>

### Check block data
Get block data of a range of blocks.

```bash
$ skycoin-cli blocks [starting block or single block seq] [ending block seq]
```

#### Example
```bash
$ skycoin-cli blocks 41 42
```

<details>
 <summary>View Output</summary>

```json
{
 "blocks": [
     {
         "header": {
             "seq": 41,
             "block_hash": "08f89cfe92be09e9848ba4d77c300908761354933f80401c107644feab1f4c9e",
             "previous_block_hash": "fad5aca57144cbc86ad916492e814ec84c825d9870a86beac81980de30b0ae60",
             "timestamp": 1429058524,
             "fee": 4561,
             "version": 0,
             "tx_body_hash": "cf4fe76a08e3296b6f6abdb949604409be66574f211d9d14fde39103c4cfe1d6"
         },
         "body": {
             "txns": [
                 {
                     "length": 220,
                     "type": 0,
                     "txid": "cf4fe76a08e3296b6f6abdb949604409be66574f211d9d14fde39103c4cfe1d6",
                     "inner_hash": "2f5942207104d52dbd6191684b2a97392e616b7fa51dde314dbddd58d34b8027",
                     "sigs": [
                         "b2b8c8ec1e1dfdeac4690e88d4ef9fcc4b52fcb771153f391cbcb58d651505a94c6263b6dc15a948c0396c0d8be20d9e0d1993b494bd9189c778d3673363bfc401"
                     ],
                     "inputs": [
                         "c65a9e6aa33244958e9595e9eceed678f9f17761753bf77000c5474f7696da53"
                     ],
                     "outputs": [
                         {
                             "uxid": "195f5e50b4eed1ec7ff968feca90356285437adc8ccfcf6623b55a4eebf7bbb5",
                             "dst": "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
                             "coins": "969790.000000",
                             "hours": 760
                         },
                         {
                             "uxid": "6bbf13da052e1baade111ae8bb85548732532c8f5286eba8345d436d315d1c93",
                             "dst": "qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5",
                             "coins": "9000.000000",
                             "hours": 760
                         }
                     ]
                 }
             ]
         }
     },
     {
         "header": {
             "seq": 42,
             "block_hash": "60a17e0cf411e5db7150272e597d343beaa5fbce5d61f6f647a14288262593b1",
             "previous_block_hash": "08f89cfe92be09e9848ba4d77c300908761354933f80401c107644feab1f4c9e",
             "timestamp": 1429058594,
             "fee": 292512,
             "version": 0,
             "tx_body_hash": "0e91a08561e85a36ddf44e77b9228f7d561c18c0b46d19083d4af511085b697e"
         },
         "body": {
             "txns": [
                 {
                     "length": 317,
                     "type": 0,
                     "txid": "0e91a08561e85a36ddf44e77b9228f7d561c18c0b46d19083d4af511085b697e",
                     "inner_hash": "d78230e22b358d7cc8d491adb3c0ec1e77a5170602a4ec92d700c4b4bb101f98",
                     "sigs": [
                         "17ba9c495e4d396a37eaf062e1806a13b3bdc91a83151c2455cf948a7e6d91882dc02ec6443970517f0f7daf59ce9b89658a17f5d51c0cbc18056811d0f3006501",
                         "e4e8f28801fe461cc8097b29cfe1307739bdfbdd6b20c31e04eef89aede641a6407fa0c41b0ad5ef167e3255e1916c0bbd358ffd70f34dc7944ffe67514bc5f501"
                     ],
                     "inputs": [
                         "f48432d381a10abecbd1357d81705ea922246e92170fe405d1a4a35c5ceef6a4",
                         "6bbf13da052e1baade111ae8bb85548732532c8f5286eba8345d436d315d1c93"
                     ],
                     "outputs": [
                         {
                             "uxid": "19efa2bd8c59623a092612c511fb66333e2049a57d546269c19255852056fead",
                             "dst": "qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5",
                             "coins": "9000.000000",
                             "hours": 48752
                         },
                         {
                             "uxid": "9953e00abe05db134510693a44b8928ca9b29d0009b38d9c4f8dcdedee7edc35",
                             "dst": "4EHiTjCsxQmt4wRy5yJxBMcxsM5yGqtuqu",
                             "coins": "1000.000000",
                             "hours": 48752
                         }
                     ]
                 }
             ]
         }
     }
 ]
}
```
</details>

### Check database integrity
Checks if the given database file contains valid skycoin blockchain data
If no argument is given, the default `data.db` in `$HOME/.$COIN/` will be checked.

```bash
$ skycoin-cli checkdb [db path]
```

#### Example
```bash
$ skycoin-cli checkdb $DB_PATH
```

<details>
 <summary>View Output</summary>

```
check db success
```
</details>

### Create a raw transaction
Create a raw transaction that can be broadcasted later.
A raw transaction is a binary encoded hex string.

```bash
$ skycoin-cli createRawTransaction [command options] [to address] [amount]
```

```
OPTIONS:
        -f value    [wallet file or path], From wallet
        -a value    [address] From address
        -c value    [changeAddress] Specify different change address.
                          By default the from address or a wallets coinbase address will be used.
        -m value    [send to many] use JSON string to set multiple receive addresses and coins,
                          example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'
        --json, -j  Returns the results in JSON format.
```

#### Examples
##### Sending to a single address from a specified wallet
```bash
$ skycoin-cli createRawTransaction -f $WALLET_PATH -a $FROM_ADDRESS $RECIPIENT_ADDRESS $AMOUNT
```

##### Sending to a specific change address

```bash
$ skycoin-cli createRawTransaction -f $WALLET_PATH -a $FROM_ADDRESS -c $CHANGE_ADDRESS $RECIPIENT_ADDRESS $AMOUNT
```

<details>
 <summary>View Output</summary>

```
dc00000000c7425e5a49fce496d78ea9b04fc47e4126b91f675b00c16b3a7515c1555c252001000000115112dbb438b423dccd5f1afb7bce3d0cd4b87b57fd9fd3e5a26ee24e05fb696f0c7f3d6a84eafd80e051117162d790fa0e57c01a0e570b8ac0ae5faa5bf782000100000005e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634020000000007445b5d6fbbb1a7d70bef941fb5da234a10fcae40420f000000000001000000000000000056500d41a1a6f1967ffe0074bb171148667ce20d0024f400000000009a05000000000000
```
</details>

##### Sending to multiple addresses
```bash
$ skycoin-cli createRawTransaction -f $WALLET_PATH -a $FROM_ADDRESS -m '[{"addr":"$ADDR1", "coins": "$AMT1"}, {"addr":"$ADDR2", "coins": "$AMT2"}]'
```

<details>
 <summary>View Output</summary>

```
01010000000e2a5bf4964604006fea5cf8cbd705e82bebb055467f10681ef01ce5c8db654801000000d951d4e34a7b35b1b165e8302cd47e09b6433ea443a8864dc8428537dbe8b76e00ee58bb195d7de3d28935ed4fc3684f1cac5593c09c4bafb016705b7e2b3393000100000005e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634030000000007445b5d6fbbb1a7d70bef941fb5da234a10fcae40420f000000000001000000000000000056500d41a1a6f1967ffe0074bb171148667ce20d40420f00000000000100000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec09c0e1e400000000009a05000000000000
```
</details>

> NOTE: When sending to multiple addresses all the receiving addresses need to be different
        Otherwise you get, `ERROR: Duplicate output in transaction`


##### Generate a JSON output
```bash
$ skycoin-cli createRawTransaction -f $WALLET_PATH -a $FROM_ADDRESS --json $RECIPIENT_ADDRESS $AMOUNT
```

<details>
 <summary>View Output</summary>

```json
{
 "rawtx": "dc00000000c7425e5a49fce496d78ea9b04fc47e4126b91f675b00c16b3a7515c1555c252001000000115112dbb438b423dccd5f1afb7bce3d0cd4b87b57fd9fd3e5a26ee24e05fb696f0c7f3d6a84eafd80e051117162d790fa0e57c01a0e570b8ac0ae5faa5bf782000100000005e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634020000000007445b5d6fbbb1a7d70bef941fb5da234a10fcae40420f000000000001000000000000000056500d41a1a6f1967ffe0074bb171148667ce20d0024f400000000009a05000000000000"
}
```
</details>

### Decode a raw transaction
```bash
$ skycoin-cli decodeRawTransaction [raw transaction]
```

Decode a raw skycoin transaction.

#### Example

```bash
skycoin-cli decodeRawTransaction dc00000000247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af384301000000cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f000100000005e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634020000000007445b5d6fbbb1a7d70bef941fb5da234a10fcae40420f00000000000100000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec090024f400000000009805000000000000
```

<details>
 <summary>View Output</summary>

```json
{
  "hash": "ee700309aba9b8b552f1c932a667c3701eff98e71c0e5b0e807485cea28170e5",
  "inner_hash": "247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af3843",
  "sigs": [
    "cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f00"
  ],
  "in": [
    "05e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634"
  ],
  "out": [
    {
      "hash": "2cb770d7c045954e9195b312e5409d0070c15361da7148879fb8658b766fae90",
      "src_tx": "247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af3843",
      "address": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
      "coins": "1.000000",
      "hours": 1
    },
    {
      "hash": "0de690eeec960274539c2ad35b57d7c0492a268a5f17ab54e5e24f3d6e14bc72",
      "src_tx": "247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af3843",
      "address": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
      "coins": "16.000000",
      "hours": 1432
    }
  ]
}
```
</details>


### Broadcast a raw transaction
Broadcast a raw skycoin transaction.
Output is the transaction id.

```bash
$ skycoin-cli broadcastTransaction [raw transaction]
```

```bash
$ skycoin-cli broadcastTransaction dc00000000247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af384301000000cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f000100000005e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634020000000007445b5d6fbbb1a7d70bef941fb5da234a10fcae40420f00000000000100000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec090024f400000000009805000000000000
```
<details>
 <summary>View Output</summary>

```
ee700309aba9b8b552f1c932a667c3701eff98e71c0e5b0e807485cea28170e5
```
</details>

### Generate a wallet
Generate a new skycoin wallet.

```bash
$ skycoin-cli generateWallet [command options]
```

```
OPTIONS:
        -r        A random alpha numeric seed will be generated for you
        --rd      A random seed consisting of 12 dictionary words will be generated for you (default)
        -s value  Your seed
        -n value  [numberOfAddresses] Number of addresses to generate
                            By default 1 address is generated. (default: 1)
        -f value  [walletName] Name of wallet. The final format will be "yourName.wlt".
                             If no wallet name is specified a generic name will be selected. (default: "skycoin_cli.wlt")
        -l value  [label] Label used to idetify your wallet.
```

#### Examples
##### Generate the default wallet
```bash
$ skycoin-cli generateWallet
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "skycoin",
     "cryptoType": "",
     "encrypted": "false",
     "filename": "skycoin_cli.wlt",
     "label": "",
     "lastSeed": "781576ec74bfa2cc9eb06f8613b96db9be21438b9cd6b6ded09df3bc5b9da279",
     "secrets": "",
     "seed": "foster blossom glare cube reopen october refuse about journey arrange music alone",
     "tm": "1523176366",
     "type": "deterministic",
     "version": "0.2"
 },
 "entries": [
     {
         "address": "FSkC7V5rAkCnNrtCe1HBnD2iTm1J8jn63V",
         "public_key": "03a16c8e9ea86ea2358364757431b84cc388b34be776bb6a23ed2b83731957d33a",
         "secret_key": "3938826649631e2abc1c47c050d0fea5ceac7c45e3fa6cd3ddf1621bdd519150"
     }
 ]
}
```
</details>

> NOTE: If a wallet with the same name already exists then the cli exits with an error.

##### Generate a wallet with a random alpha numeric seed
```bash
$ skycoin-cli generateWallet -r
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "skycoin",
     "cryptoType": "",
     "encrypted": "false",
     "filename": "skycoin_cli.wlt",
     "label": "",
     "lastSeed": "fdaf0729903fbd5962301f16531a1da102bf0875b4a636cb43ce24b967b932ac",
     "secrets": "",
     "seed": "8af187f04c306538544a1c2c4d0a51e9220bd17fc2fcb3fd72ba2ca3ce7aa212",
     "tm": "1523177044",
     "type": "deterministic",
     "version": "0.2"
 },
 "entries": [
     {
         "address": "9YogvtjYgeLn3gQX2wzsXDpZn7LuoArdzZ",
         "public_key": "022b4bd33f0ad037756ae19f8dfab935fed1118980b4067b4a6b7f03333ba5ccae",
         "secret_key": "b4cf1731be9f930ba3a67179eed5dca5af2adee1ce4df96383923f775bf575c0"
     }
 ]
}
```
</details>

##### Generate a wallet with a 12 word mnemomic seed
```bash
$ skycoin-cli generateWallet -rd
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "skycoin",
     "cryptoType": "",
     "encrypted": "false",
     "filename": "skycoin_cli.wlt",
     "label": "",
     "lastSeed": "f219c2e902940f27ea735d866a495372debcbd01da287a2ec1226d0eb43b9890",
     "secrets": "",
     "seed": "motor cross wrap intact soup critic club allow track come dizzy cool",
     "tm": "1523177162",
     "type": "deterministic",
     "version": "0.2"
 },
 "entries": [
     {
         "address": "E9p6Eck7Q6bYBnEkCdB3vCDf3YYkQxCHwv",
         "public_key": "02c41e7b03a6a848a417d7d270b9d83c4d9534c2cd5eace8046c67d012b920f1db",
         "secret_key": "41b6aa1780f425dac942c8bd1570248ebfca24778e866705a6573b17ead57a4d"
     }
 ]
}
```
</details>

##### Generate a wallet with a specified seed
```bash
$ skycoin-cli generateWallet -s "this is the super secret seed everyone needs but does not have"
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "skycoin",
     "cryptoType": "",
     "encrypted": "false",
     "filename": "skycoin_cli.wlt",
     "label": "",
     "lastSeed": "c34a83b473ea4d2f9dc394d0b9c1c0d4578012252b842ef1bfce9950cfe50b06",
     "secrets": "",
     "seed": "this is the super secret everyone needs but does not have",
     "tm": "1523178336",
     "type": "deterministic",
     "version": "0.2"
 },
 "entries": [
     {
         "address": "NEQVmBJPidzo3SfDRJHNDMHL7VbqNa7Cku",
         "public_key": "0348400c3c1a733a6e25c77f1ffea64c887bc9344a0366821ef07b9b3abadcaf10",
         "secret_key": "42e3906d86ca25eb408d2af90b0810d7831b7d777e756021b607bca6538952eb"
     }
 ]
}
```
</details>


##### Generate more than 1 default address
```bash
$ skycoin-cli generateWallet -n 2
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "skycoin",
     "cryptoType": "",
     "encrypted": "false",
     "filename": "skycoin_cli.wlt",
     "label": "",
     "lastSeed": "861a8989e6c85fb69cf5968586fe9d5a1e26936ab122c5d542bf78fb35e0d247",
     "secrets": "",
     "seed": "cause custom canal kitchen short cement round cat shine renew pair crowd",
     "tm": "1523178418",
     "type": "deterministic",
     "version": "0.2"
 },
 "entries": [
     {
         "address": "2accTtyD7tqzLh7c62BE9zjGiyEpoMyQ3bb",
         "public_key": "027c30928161755c913e1b3db208f95a66be0f550b9620cefd44902b5354365b73",
         "secret_key": "89b2f71fb773a00480637fd83c93e27499fd5e55f69a6e2b58f0847c3ce5040c"
     },
     {
         "address": "goyx9VE3q73zAWntmwwyaUoTZhtTyG4vt",
         "public_key": "025c0b06471b865cb5eab23f9a9dc0a992fe70d0576eb400aa4978ddd0a2124b95",
         "secret_key": "75deabceedb9b09a109f5d982fba13a56622d93916a8ef81ddccca69fcc9d7e3"
     }
 ]
}
```
</details>

##### Generate wallet with a custom wallet name
```bash
$ skycoin-cli generateWallet -f "secret_wallet.wlt"
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "skycoin",
     "cryptoType": "",
     "encrypted": "false",
     "filename": "secret_wallet.wlt",
     "label": "",
     "lastSeed": "02a240aa6d0dcc8262756bef2ba1b1ffbf5f5665b8d6b6863a4c833c9b5ae8e3",
     "secrets": "",
     "seed": "bundle model dice age profit child ribbon below tide load grocery leave",
     "tm": "1523178575",
     "type": "deterministic",
     "version": "0.2"
 },
 "entries": [
     {
         "address": "23ycmCLQGzjRu6orabHCvPhwJspVWq8HHVE",
         "public_key": "03d7c80bc37912edc0450aa782c88e1a2bb83199c3884c42e624a3ba346636c2bd",
         "secret_key": "bf2237b0b0fd036fe6ee6a92bd5fee6034f4c22d6c3607d63996ff4ae741238c"
     }
 ]
}
```
</details>

> NOTE: The wallet name needs to end with `.wlt` and it should not be a path.

##### Generate wallet with a custom wallet label
By default the wallet label is an empty field
```bash
$ skycoin-cli generateWallet -l "cli wallet"
```

<details>
 <summary>View Output</summary>

```json
{
 "meta": {
     "coin": "skycoin",
     "cryptoType": "",
     "encrypted": "false",
     "filename": "skycoin_cli.wlt",
     "label": "cli wallet",
     "lastSeed": "b3b3c13419a8343f8845a8de30543fa33680e25251a3a1bda3e49346f1d640f9",
     "secrets": "",
     "seed": "offer spoil crane trial submit kite venture edit repair mushroom fetch bounce",
     "tm": "1523178769",
     "type": "deterministic",
     "version": "0.2"
 },
 "entries": [
     {
         "address": "21YPgFwkLxQ1e9JTCZ43G7JUyCaGRGqAsda",
         "public_key": "03784cf30195259e4bf89e15d343417d38ecd05b2f61fd2b2f71020ad7b1de3577",
         "secret_key": "8f6f2e3b63310f94c1440ba230eb170dbc1ffd2ad355274c05b169c290216a3c"
     }
 ]
}
```
</details>

### Generate addresses for a wallet
Generate new addresses for a skycoin wallet.

```bash
$ skycoin-cli generateAddresses [command options]
```

```
OPTIONS:
        -n value    [numberOfAddresses]  Number of addresses to generate (default: 1)
        -f value    [wallet file or path] Generate addresses in the wallet (default: $HOME/.skycoin/wallets//wallets/skycoin_cli.wlt)
        --json, -j  Returns the results in JSON format
```

##### Generate an address for the default wallet
```bash
$ skycoin-cli generateAddresses
```

<details>
 <summary>View Output</summary>

```
2mEgmYt6NZHA1erYqbAeXmGPD5gqLZ9toFv
```
</details>

##### Generate an address for a specific wallet
```bash
$ skycoin-cli generateAddresses $WALLET_PATH
```

<details>
 <summary>View Output</summary>

```json
2cET6L4c6Bee5jucuzsTQUXFxWX76GZoDqv
```
</details>

##### Generate `n` addresses
```bash
$ skycoin-cli generateAddresses -n 2
```

<details>
 <summary>View Output</summary>

```
2UrEV3Vyu5RJABZNukKRq25ggrrg96RUwdH,LJN5qGmLbJxLswzD3nFn3RFcmWJyZ2LGHY
```
</details>

##### Generate a JSON output
```bash
$ skycoin-cli generateAddresses --json
```

<details>
 <summary>View Output</summary>

```json
{
 "addresses": [
     "QuLaPirJNUkBpMoe5tzzY7j6nJ5maUVJF1"
 ]
}
```
</details>

### Last blocks
Show the last `n` skycoin blocks.
By default the last block is shown.

```bash
$  skycoin-cli lastBlocks [numberOfBlocks]
```

#### Examples
##### Get the last block
```bash
$ skycoin-cli lastBlocks
```

<details>
 <summary>View Output</summary>

```json
{
 "blocks": [
     {
         "header": {
             "seq": 21202,
             "block_hash": "07140b536a5d0c3fcecec6cd61d2d4628a6fe6d7f7933223365269ea78d47c06",
             "previous_block_hash": "baf74fae7ba8e304236e84e4cbc24b810a3827512bfdb82f892ffcb8682f9d16",
             "timestamp": 1523179526,
             "fee": 38,
             "version": 0,
             "tx_body_hash": "31b4f132374b4eb6c31c1c2f51a07a336a5d290b859d4c32a4982f325f124023"
         },
         "body": {
             "txns": [
                 {
                     "length": 220,
                     "type": 0,
                     "txid": "31b4f132374b4eb6c31c1c2f51a07a336a5d290b859d4c32a4982f325f124023",
                     "inner_hash": "78b51146a1822c777f5de05fd9c78bb457bcbaeec13d9c927f548cc4019a467c",
                     "sigs": [
                         "22e13afd39ac22ed452b716c721c82374b4d6d69fe2db8a70ea74f9fa37487761039365a5de0bf741e7c4da171b2164ddf8bd4326b9845784719645fa16a8ec401"
                     ],
                     "inputs": [
                         "1714a29f821bc955857132a495e64c458794564108245bb605d2eb418edb8b54"
                     ],
                     "outputs": [
                         {
                             "uxid": "8d2fe35c36f69866ea91ff468464ddb6e0bee7fd145df7319fe192245c0dd646",
                             "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                             "coins": "6458.000000",
                             "hours": 19
                         },
                         {
                             "uxid": "4d0a689781cd6aaf61d33f8f0bf87f16210c3e5b45f83a789f454530c2d64b21",
                             "dst": "i66Gax2z4cMiTYGHjRjmKpooWS2Rhnmoeu",
                             "coins": "1.000000",
                             "hours": 19
                         }
                     ]
                 }
             ]
         }
     }
 ]
}
```
</details>

##### Get the last `n` blocks
```bash
$ skycoin-cli lastBlocks 3
```

<details>
 <summary>View Output</summary>

```json
{
 "blocks": [
     {
         "header": {
             "seq": 21202,
             "block_hash": "07140b536a5d0c3fcecec6cd61d2d4628a6fe6d7f7933223365269ea78d47c06",
             "previous_block_hash": "baf74fae7ba8e304236e84e4cbc24b810a3827512bfdb82f892ffcb8682f9d16",
             "timestamp": 1523179526,
             "fee": 38,
             "version": 0,
             "tx_body_hash": "31b4f132374b4eb6c31c1c2f51a07a336a5d290b859d4c32a4982f325f124023"
         },
         "body": {
             "txns": [
                 {
                     "length": 220,
                     "type": 0,
                     "txid": "31b4f132374b4eb6c31c1c2f51a07a336a5d290b859d4c32a4982f325f124023",
                     "inner_hash": "78b51146a1822c777f5de05fd9c78bb457bcbaeec13d9c927f548cc4019a467c",
                     "sigs": [
                         "22e13afd39ac22ed452b716c721c82374b4d6d69fe2db8a70ea74f9fa37487761039365a5de0bf741e7c4da171b2164ddf8bd4326b9845784719645fa16a8ec401"
                     ],
                     "inputs": [
                         "1714a29f821bc955857132a495e64c458794564108245bb605d2eb418edb8b54"
                     ],
                     "outputs": [
                         {
                             "uxid": "8d2fe35c36f69866ea91ff468464ddb6e0bee7fd145df7319fe192245c0dd646",
                             "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                             "coins": "6458.000000",
                             "hours": 19
                         },
                         {
                             "uxid": "4d0a689781cd6aaf61d33f8f0bf87f16210c3e5b45f83a789f454530c2d64b21",
                             "dst": "i66Gax2z4cMiTYGHjRjmKpooWS2Rhnmoeu",
                             "coins": "1.000000",
                             "hours": 19
                         }
                     ]
                 }
             ]
         }
     },
     {
         "header": {
             "seq": 21203,
             "block_hash": "5b1f652bd5639e7cd2f9e5af4c1f86685b3f8104bb2039e2d6e56ec90003b086",
             "previous_block_hash": "07140b536a5d0c3fcecec6cd61d2d4628a6fe6d7f7933223365269ea78d47c06",
             "timestamp": 1523179676,
             "fee": 10,
             "version": 0,
             "tx_body_hash": "3c2881347dcd6eb767f447720fc337dc75c8698ff1d25e6ebc490befdb4123fa"
         },
         "body": {
             "txns": [
                 {
                     "length": 220,
                     "type": 0,
                     "txid": "3c2881347dcd6eb767f447720fc337dc75c8698ff1d25e6ebc490befdb4123fa",
                     "inner_hash": "ec69bb1f2c2f1acffb95a7623a8f85b0797c48fec251b64622f52b0cb71f1e9f",
                     "sigs": [
                         "9287f01192d7b8afe9cd7b15c2859bb13032b96b24d791eacadd612c14a0d30e7f1a503bcb01eb5f4d2ccd2400c320a82fd2825a410309f00c41086debf639d701"
                     ],
                     "inputs": [
                         "8d2fe35c36f69866ea91ff468464ddb6e0bee7fd145df7319fe192245c0dd646"
                     ],
                     "outputs": [
                         {
                             "uxid": "814922f2428e6a5c8acdd364785309cc9a446a10d80901baf0600c5de3b30171",
                             "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                             "coins": "6457.000000",
                             "hours": 5
                         },
                         {
                             "uxid": "70b6b19d632e26b2a2ccae4dce6fc4093c5e02cfde03da6f893e8db14c7afa44",
                             "dst": "i66Gax2z4cMiTYGHjRjmKpooWS2Rhnmoeu",
                             "coins": "1.000000",
                             "hours": 4
                         }
                     ]
                 }
             ]
         }
     },
     {
         "header": {
             "seq": 21204,
             "block_hash": "bb3dd979527f6133e8a5b9bb9792d2fa964daf6158834a02034719f0e943d049",
             "previous_block_hash": "5b1f652bd5639e7cd2f9e5af4c1f86685b3f8104bb2039e2d6e56ec90003b086",
             "timestamp": 1523179736,
             "fee": 3,
             "version": 0,
             "tx_body_hash": "49f0f5813ecfee45556f769c1971cf8b4300ec4ce9f05deedb57b8048e2bd0f9"
         },
         "body": {
             "txns": [
                 {
                     "length": 220,
                     "type": 0,
                     "txid": "49f0f5813ecfee45556f769c1971cf8b4300ec4ce9f05deedb57b8048e2bd0f9",
                     "inner_hash": "11a949890d63993b98a92418ba2682795ebad30583e4c16a108939a6945a76a9",
                     "sigs": [
                         "9d23748dac110ab464d9cbc2de472fdf00da804a8da440d592667f11fdf7f4937620ddcc265f75751ca2566d3a39df0b04f7ab638ecff5f32f5d2699e627820700"
                     ],
                     "inputs": [
                         "814922f2428e6a5c8acdd364785309cc9a446a10d80901baf0600c5de3b30171"
                     ],
                     "outputs": [
                         {
                             "uxid": "3ebedfac8cef4d3fc27c48060277a8b4880c047ed77a72ffdab40fdb274c5a93",
                             "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                             "coins": "6456.000000",
                             "hours": 1
                         },
                         {
                             "uxid": "78500a48b430d6c6e6b744bdb193df7a4f90c60ec0bf6a7bebf6b6e34f49e6f2",
                             "dst": "2ESrbYizrvrZNK5zFoHsj69a7fTVXxTTpMi",
                             "coins": "1.000000",
                             "hours": 1
                         }
                     ]
                 }
             ]
         }
     }
 ]
}
```
</details>


### List wallet addresses
List addresses in a skycoin wallet.

```bash
$ skycoin-cli listAddresses [walletName]
```

If no `walletName` is given then default wallet ($HOME/.skycoin/wallets/skycoin_cli.wlt) is used.

> NOTE: The wallet name `skycoin_cli.wlt` or full path `$HOME/.skycoin/wallets/skycoin_cli.wlt` can be used.
        When only the wallet name is given then the default wallet dir, $HOME/.$COIN/wallets is used.

#### Examples
##### List addresses of default wallet
```bash
$ skycoin-cli listAddresses
```

<details>
 <summary>View Output</summary>

```json
{
 "addresses": [
     "21YPgFwkLxQ1e9JTCZ43G7JUyCaGRGqAsda",
     "2mEgmYt6NZHA1erYqbAeXmGPD5gqLZ9toFv",
     "2cET6L4c6Bee5jucuzsTQUXFxWX76GZoDqv",
     "2UrEV3Vyu5RJABZNukKRq25ggrrg96RUwdH",
     "LJN5qGmLbJxLswzD3nFn3RFcmWJyZ2LGHY",
     "QuLaPirJNUkBpMoe5tzzY7j6nJ5maUVJF1"
 ]
}
```
</details>

##### List addresses of a specific wallet
```bash
$ skycoin-cli listAddresses $WALLET_NAME or $WALLET_PATH
```

<details>
 <summary>View Output</summary>

```json
{
 "addresses": [
     "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
     "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
     "bjN9ckj5HRvgDfcvKNboM8cvohJzy9oXJV"
 ]
}
```
</details>

### List wallets
List wallets in the skycoin wallet directory.

```bash
$ skycoin-cli listWallets
```

#### Example
```bash
$ skycoin-cli listWallets
```

<details>
 <summary>View Output</summary>

```json
{
 "wallets": [
     {
         "name": "2018_02_04_45bc.wlt",
         "label": "Your Wallet",
         "address_num": 60
     },
     {
         "name": "2018_03_22_6e61.wlt",
         "label": "craptopia",
         "address_num": 3
     },
     {
         "name": "2018_04_01_198c.wlt",
         "label": "wings",
         "address_num": 2
     },
     {
         "name": "secret_wallet.wlt",
         "label": "",
         "address_num": 1
     },
     {
         "name": "skycoin_cli.wlt",
         "label": "cli wallet",
         "address_num": 6
     }
 ]
}
```
</details>

### Send
Make a skycoin transaction.

```bash
$ skycoin-cli send [command options] [to address] [amount]
```

```
OPTIONS:
        -f value    [wallet file or path] From wallet. If no path is specified your default wallet
                    (`$HOME/.skycoin/wallets/skycoin_cli.wlt`) path will be used.
        -a value    [address] From address
        -c value    [changeAddress] Specify change address, by default the from address or
                          the wallet's coinbase address will be used
        -m value    [send to many] use JSON string to set multiple recive addresses and coins,
                          example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'
        --json, -j  Returns the results in JSON format.
```

#### Examples
##### Sending from the default wallet
```bash
$ skycoin-cli send $RECIPIENT_ADDRESS $AMOUNT
```

##### Sending from a specific wallet
```bash
$ skycoin-cli send -f $WALLET_PATH $RECIPIENT_ADDRESS $AMOUNT
```

##### Sending from a specific address in a wallet
```bash
$ skycoin-cli send -f $WALLET_PATH -a $FROM_ADDRRESS $RECIPIENT_ADDRESS $AMOUNT
```

> NOTE: If $WALLET_PATH is not specified above then the default wallet is used.

##### Sending change to a specific change address
```bash
$ skycoin-cli send -f $WALLET_PATH -a $FROM_ADDRESS -c $CHANGE_ADDRESS $RECIPIENT_ADDRESS $AMOUNT
```

##### Sending to multiple addresses
```bash
$ skycoin-cli send -f $WALLET_PATH -a $FROM_ADDRESS -m '[{"addr":"$ADDR1", "coins": "$AMT1"}, {"addr":"$ADDR2", "coins": "$AMT2"}]'
```

<details>
 <summary>View Output</summary>

```
txid:$TRANSACTION_ID
```
</details>

##### Generate a JSON output
```bash
$ skycoin-cli send -f $WALLET_PATH -a $FROM_ADDRESS --json $RECIPIENT_ADDRESS $AMOUNT
```

<details>
 <summary>View Output</summary>

```json
{
 "txid": "$TRANSACTION_ID"
}
```
</details>

### Show Config
Show the CLI tool's local configuration.

#### Example
```bash
$ skycoin-cli showConfig
```

<details>
 <summary>View Output</summary>

```json
{
    "wallet_directory": "/home/user/.skycoin/wallets",
    "wallet_name": "skycoin_cli.wlt",
    "data_directory": "/home/user/.skycoin",
    "coin": "skycoin",
    "rpc_address": "http://127.0.0.1:6420",
    "use_csrf": false
}
```

### Status
#### Example
```bash
$ skycoin-cli status
```

<details>
 <summary>View Output</summary>

```json
{
 "running": true,
 "num_of_blocks": 21210,
 "hash_of_last_block": "d5797705bfc0ac7956f3eeaa083aec4e89a6b27ada7499c5a53dad2fda84c5f9",
 "time_since_last_block": "18446744073709551591s",
 "webrpc_address": "127.0.0.1:6420"
}
```
</details>

### Get transaction
Get transaction data from a `txid`.

```bash
$ skycoin-cli transaction [transaction id]
```

#### Example
```bash
$ skycoin-cli transaction 824d421a25f81aa7565d042a54b3e1e8fdc58bed4eefe8f8a90748da6d77d135
```

<details>
 <summary>View Output</summary>

```json
{
 "transaction": {
     "status": {
         "confirmed": true,
         "unconfirmed": false,
         "height": 1,
         "block_seq": 864,
         "unknown": false
     },
     "txn": {
         "length": 220,
         "type": 0,
         "txid": "824d421a25f81aa7565d042a54b3e1e8fdc58bed4eefe8f8a90748da6d77d135",
         "inner_hash": "708a21c685041ba409b9634843003f263c7a00d99459925e319049f3e36f1163",
         "timestamp": 1492141347,
         "sigs": [
             "6772c4e1370faf7acd8fc37f6dba3bea06ad1beb1300074c33b2c6fa9b11ed7d2bcc01b7008f235aa918b9c146078dfb8e1c8ce55b0893ea5f111597b42331ba01"
         ],
         "inputs": [
             "c38c108ac3c76e5faffce0bb83153ec98bc1355a98e1a9b0f95ab1b98ef9f00e"
         ],
         "outputs": [
             {
                 "uxid": "b0586a8e731c475e87eb61ef0b845d7893cf39120a1e97cf05f78585f1a49e3c",
                 "dst": "2bfYafFtdkCRNcCyuDvsATV66GvBR9xfvjy",
                 "coins": "166365",
                 "hours": 0
             },
             {
                 "uxid": "49c64719d8df905a7fd4a1c46c2a9c358a8c5ae14befb0d57dcfe1b1d36a1758",
                 "dst": "ep3axwpJ3hWWQcACu48z9sMKUB7snXBm94",
                 "coins": "1300",
                 "hours": 0
             }
         ]
     }
 }
}
```
</details>

### Verify address
Verify whether a given address is a valid skycoin addres or not.

```bash
$  skycoin-cli verifyAddress [skycoin address]
```

#### Example
##### Valid addresss

```bash
$ skycoin-cli verifyAddress 21YPgFwkLxQ1e9JTCZ43G7JUyCaGRGqAsda
```

```
No Output
```

##### Invalid Address
###### Invalid checksum
```bash
$ skycoin-cli verifyAddress 21YPgFwkLxQ1e9JTCZ43G7JUyCaGRGqAsdx
```

<details>
 <summary>View Output</summary>

```
Invalid checksum
```
</details>

###### Invalid address length
```bash
$ skycoin-cli verifyAddress 21YPg
```

<details>
 <summary>View Output</summary>

```
Invalid address length
```
</details>


### Check wallet balance
Check the wallet a skycoin wallet.

```bash
$ skycoin-cli walletBalance [wallet]
```

> NOTE: Both the full wallet path or only the wallet name can be used.
        If no wallet is specified then the default wallet: `$HOME/.$COIN/wallets/skycoin_cli.wlt` is used.

#### Example
##### Balance of default wallet
```bash
$ skycoin-cli walletBalance
```

<details>
 <summary>View Output</summary>

```json
{
    "confirmed": {
        "coins": "123.000000",
        "hours": "456"
    },
    "spendable": {
        "coins": "123.000000",
        "hours": "456"
    },
    "expected": {
        "coins": "123.000000",
        "hours": "456"
    },
    "addresses": [
        {
            "confirmed": {
                "coins": "123.000000",
                "hours": "456"
            },
            "spendable": {
                "coins": "123.000000",
                "hours": "456"
            },
            "expected": {
                "coins": "123.000000",
                "hours": "456"
            },
            "address": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc"
        }, {
            "confirmed": {
                "coins": "0.000000",
                "hours": "0"
            },
            "spendable": {
                "coins": "0.000000",
                "hours": "0"
            },
            "expected": {
                "coins": "0.000000",
                "hours": "0"
            },
            "address": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"
        }
    ]
}
```
</details>


##### Balance of a specific wallet
```bash
$ skycoin-cli walletBalance 2018_04_01_198c.wlt
```
*OR*

```bash
$ skycoin-cli walletBalance ~/.skycoin/wallets/2018_04_01_198c.wlt
```

<details>
 <summary>View Output</summary>

```json
{
 "confirmed": {
     "coins": "31.000000",
     "hours": "25255"
 },
 "spendable": {
     "coins": "31.000000",
     "hours": "25255"
 },
 "expected": {
     "coins": "31.000000",
     "hours": "25255"
 },
 "addresses": [
     {
         "confirmed": {
             "coins": "0.000000",
             "hours": "0"
         },
         "spendable": {
             "coins": "0.000000",
             "hours": "0"
         },
         "expected": {
             "coins": "0.000000",
             "hours": "0"
         },
         "address": "29fDBQuJs2MDLymJsjyWH6rDjsyv995SrGU"
     },
     {
         "confirmed": {
             "coins": "31.000000",
             "hours": "25255"
         },
         "spendable": {
             "coins": "31.000000",
             "hours": "25255"
         },
         "expected": {
             "coins": "31.000000",
             "hours": "25255"
         },
         "address": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V"
     }
 ]
}
```
</details>

### See wallet directory
Get the current skycoin wallet directory.

```bash
$ skycoin-cli walletDir [command options]
```

```
OPTIONS:
        -j, --json  Returns the results in JSON format.
```

#### Examples
##### Text output
```bash
$ skycoin-cli walletDir
```
<details>
 <summary>View Output</summary>

```
$WALLET_DIR
```
</details>

##### JSON output
```bash
$ skycoin-cli walletDir --json
```

<details>
 <summary>View Output</summary>

```json
{
 "walletDir": "$WALLET_DIR"
}
```
</details>

### List wallet transaction history
Show all previous transactions made by the addresses in a wallet.

```bash
$ skycoin-cli walletHistory [command options]
```

```
OPTIONS:
        -f value  [wallet file or path] From wallet. If no path is specified your default wallet path will be used.
```

#### Examples
##### Default wallet
```bash
$ skycoin-cli walletHistory
```

##### Specific wallet
```bash
$ skycoin-cli walletHistory -f $WALLET_NAME
```
*OR*
```bash
$ skycoin-cli walletHistory -f $WALLET_PATH
```

<details>
 <summary>View Output</summary>

```json
[
 {
     "txid": "d1ded06a49b7588b897a2186bbe76de7ee93f49084ad35e1a7f47cbf6cd3a7fa",
     "address": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
     "amount": "1.000000",
     "timestamp": "2018-01-28T13:11:15Z",
     "status": 1
 },
 {
     "txid": "ad191f910e5508e0b0e0ab24ba815e784a1a2b63ca21043e7746bebf25106742",
     "address": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
     "amount": "1.000000",
     "timestamp": "2018-01-28T13:26:15Z",
     "status": 1
 }
]
```
</details>

### List wallet outputs
List unspent outputs of all addresses in a wallet.

```bash
$ skycoin-cli walletOutputs [wallet file]
```

#### Examples
##### Default wallet
```bash
$ skycoin-cli walletOutputs
```

##### Specific wallet
```bash
$ skycoin-cli walletHistory $WALLET_NAME
```
*OR*
```bash
$ skycoin-cli walletHistory $WALLET_PATH
```

<details>
 <summary>View Output</summary>

```json
{
 "outputs": {
     "head_outputs": [
         {
             "hash": "c51b2692aa9f296a3cd2f37b14f39c496c82f5c5ae01c54701ea60b7353f27e2",
             "time": 1523184376,
             "block_seq": 21221,
             "src_tx": "f3c5cfd462d95e724b7d35b1688c53f25a5f358f2eb9a6f87b63cdf31deb2bf8",
             "address": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
             "coins": "15.000000",
             "hours": 369,
             "calculated_hours": 370
         },
         {
             "hash": "a0777af14223bbbd5aeb8bf3cfd6ba94c776c6eec731310caaaaee49b9feb9a5",
             "time": 1523184176,
             "block_seq": 21220,
             "src_tx": "4acd61d7aa7dfe20795e517d7560643d049036af9451bcbd762793bcb6a4a6de",
             "address": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
             "coins": "1.000000",
             "hours": 0,
             "calculated_hours": 0
         }
     ],
     "outgoing_outputs": [],
     "incoming_outputs": []
 }
}
```
</details>

### CLI version
Get version of current skycoin cli.

```bash
$ skycoin-cli version [command options]
```

```
OPTIONS:
        --json, -j  Returns the results in JSON format
```

#### Examples
##### Text output
```bash
$ skycoin-cli version
```

<details>
 <summary>View Output</summary>

```
skycoin:0.23.0
cli:0.23.0
rpc:0.23.0
wallet:0.23.0
```
</details>

##### JSON output
```bash
$ skycoin-cli version --json
```

<details>
 <summary>View Output</summary>

```json
{
 "skycoin": "0.23.0",
 "cli": "0.23.0",
 "rpc": "0.23.0",
 "wallet": "0.23.0"
}
```
</details>

## Note

The `[option]` in subcommand must be set before the rest of the values, otherwise the `option` won't
be parsed. For example:

If we want to specify a `change address` in `send` command, we can use `-c` option, if you run
the command in the following way:

```bash
$ skycoin-cli send $RECIPIENT_ADDRESS $AMOUNT -c $CHANGE_ADDRESS
```

The change coins won't go to the address as you wish, it will go to the
default `change address`, which can be by `from` address or the wallet's
coinbase address.

The right script should look like this:

```bash
$ skycoin-cli send -c $CHANGE_ADDRESS $RECIPIENT_ADDRESS $AMOUNT
```
