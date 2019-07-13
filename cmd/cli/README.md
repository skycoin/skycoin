# CLI Documentation

Skycoin command line interface

The CLI command APIs can be used directly from a Go application, see [Skycoin CLI Godoc](https://godoc.org/github.com/skycoin/skycoin/src/cli).

<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3" -->

- [Install](#install)
- [Environment Setting](#environment-setting)
	- [RPC_ADDR](#rpc_addr)
	- [RPC_USER](#rpc_user)
	- [RPC_PASS](#rpc_pass)
	- [WALLET_DIR](#wallet_dir)
	- [WALLET_NAME](#wallet_name)
- [Usage](#usage)
	- [Add Private Key](#add-private-key)
	- [Check address balance](#check-address-balance)
	- [Generate new addresses](#generate-new-addresses)
	- [Generate distribution addresses for a new fiber coin](#generate-distribution-addresses-for-a-new-fiber-coin)
	- [Check address outputs](#check-address-outputs)
	- [Check block data](#check-block-data)
	- [Check database integrity](#check-database-integrity)
	- [Create a raw transaction](#create-a-raw-transaction)
	- [Decode a raw transaction](#decode-a-raw-transaction)
	- [Encode a JSON transaction](#encode-a-json-transaction)
	- [Broadcast a raw transaction](#broadcast-a-raw-transaction)
	- [Create a wallet](#create-a-wallet)
	- [Add addresses to a wallet](#add-addresses-to-a-wallet)
	- [Encrypt Wallet](#encrypt-wallet)
	- [Examples](#examples)
	- [Decrypt Wallet](#decrypt-wallet)
	- [Example](#example)
	- [Last blocks](#last-blocks)
	- [List wallet addresses](#list-wallet-addresses)
	- [List wallets](#list-wallets)
	- [Rich list](#rich-list)
	- [Send](#send)
	- [Show Seed](#show-seed)
	- [Show Config](#show-config)
	- [Status](#status)
	- [Get transaction](#get-transaction)
	- [Get address transactions](#get-address-transactions)
	- [Verify address](#verify-address)
	- [Check wallet balance](#check-wallet-balance)
	- [See wallet directory](#see-wallet-directory)
	- [List wallet transaction history](#list-wallet-transaction-history)
	- [List wallet outputs](#list-wallet-outputs)
	- [Richlist](#richlist)
	- [Address Count](#address-count)
	- [CLI version](#cli-version)
- [Note](#note)

<!-- /MarkdownTOC -->


## Install

```bash
$ cd $GOPATH/src/github.com/skycoin/skycoin/cmd/cli
$ ./install.sh
```

## Environment Setting

The CLI uses environment variable to manage the configurations.

### RPC_ADDR

CLI will connect to skycoin node REST API address `http://127.0.0.1:6420` by default.
You can change the address by setting the `RPC_ADDR` environment variable
with the following command:

```bash
$ export RPC_ADDR=http://127.0.0.1:6420
```

Note: `RPC_ADDR` must be in `scheme://host` format.

### RPC_USER

A username for authenticating requests to the skycoin node.

```bash
$ export RPC_USER=...
```

### RPC_PASS

A password for authenticating requests to the skycoin node.

```bash
$ export RPC_PASS=...
```

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

## Usage

After the installation, you can run `skycoin-cli` to see the usage:

```
$ skycoin-cli

USAGE:
  skycoin-cli [command] [flags] [arguments...]

DESCRIPTION:
    The skycoin command line interface

COMMANDS:
  addPrivateKey        Add a private key to specific wallet
  addressBalance       Check the balance of specific addresses
  addressGen           Generate skycoin or bitcoin addresses
  addressOutputs       Display outputs of specific addresses
  addressTransactions  Show detail for transaction associated with one or more specified addresses
  blocks               Lists the content of a single block or a range of blocks
  broadcastTransaction Broadcast a raw transaction to the network
  checkdb              Verify the database
  createRawTransaction Create a raw transaction to be broadcast to the network later
  decodeRawTransaction Decode raw transaction
  decryptWallet        Decrypt wallet
  encryptWallet        Encrypt wallet
  fiberAddressGen      Generate addresses and seeds for a new fiber coin
  help                 Help about any command
  lastBlocks           Displays the content of the most recently N generated blocks
  listAddresses        Lists all addresses in a given wallet
  listWallets          Lists all wallets stored in the wallet directory
  richlist             Get skycoin richlist
  send                 Send skycoin from a wallet or an address to a recipient address
  showConfig           Show cli configuration
  showSeed             Show wallet seed
  status               Check the status of current skycoin node
  transaction          Show detail info of specific transaction
  verifyAddress        Verify a skycoin address
  version              List the current version of Skycoin components
  walletAddAddresses   Generate additional addresses for a wallet
  walletBalance        Check the balance of a wallet
  walletCreate         Generate a new wallet
  walletDir            Displays wallet folder address
  walletHistory        Display the transaction history of specific wallet. Requires skycoin node rpc.
  walletOutputs        Display outputs of specific wallet

FLAGS:
  -h, --help      help for skycoin-cli
      --version   version for skycoin-cli

Use "skycoin-cli [command] --help" for more information about a command.

ENVIRONMENT VARIABLES:
    RPC_ADDR: Address of RPC node. Must be in scheme://host format. Default "http://127.0.0.1:6420"
    RPC_USER: Username for RPC API, if enabled in the RPC.
    RPC_PASS: Password for RPC API, if enabled in the RPC.
    COIN: Name of the coin. Default "skycoin"
    WALLET_DIR: Directory where wallets are stored. This value is overridden by any subcommand flag specifying a wallet filename, if that filename includes a path. Default "$DATA_DIR/wallets"
    WALLET_NAME: Name of wallet file (without path). This value is overridden by any subcommand flag specifying a wallet filename. Default "$COIN_cli.wlt"
    DATA_DIR: Directory where everything is stored. Default "$HOME/.$COIN/"
```

### Add Private Key
Add a private key to a skycoin wallet.

```bash
$ skycoin-cli addPrivateKey [flags] [private key]
```

```
FLAGS:
  -h, --help                 help for addPrivateKey
  -p, --password string      Wallet password
  -f, --wallet-file string   wallet file or path. If no path is specified your default wallet path will be used.
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
$ skycoin-cli addressGen [flags]
```

```
FLAGS:
  -c, --coin string    Coin type. Must be skycoin or bitcoin. If bitcoin, secret keys are in Wallet Import Format instead of hex. (default "skycoin")
  -x, --encrypt        Encrypt the wallet when printing a JSON wallet
  -e, --entropy int    Entropy of the autogenerated bip39 seed, when the seed is not provided. Can be 128 or 256 (default 128)
      --hex            Use hex(sha256sum(rand(1024))) (CSPRNG-generated) as the seed if not seed is not provided
  -i, --hide-secrets   Hide the secret key and seed from the output when printing a JSON wallet file
  -l, --label string   Wallet label to use when printing or writing a wallet file
  -m, --mode string    Output mode. Options are wallet (prints a full JSON wallet), addresses (prints addresses in plain text), secrets (prints secret keys in plain text) (default "wallet")
  -n, --num int        Number of addresses to generate (default 1)
  -s, --seed string    Seed for deterministic key generation. Will use bip39 as the seed if not provided.
  -t, --strict-seed    Seed should be a valid bip39 mnemonic seed.
```

#### Examples
##### Generate `n` number of skycoin addresses
```bash
$ skycoin-cli addressGen --num 2
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
$ skycoin-cli addressGen --num 2 --coin bitcoin
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
$ skycoin-cli addressGen --num 2 --hide-secrets
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
$ skycoin-cli addressGen --num 2 --mode addresses
```

```
7HVmKni3ggMdtseynSkNkqoCnsH7vkS6cg
2j5QSbHgLWXA2qXZvLzJHRo6Cissxer4CSt
```

> NOTE: If no seed is provided with the `--seed flag` and `--hex` flag is not used then bip39 is used to generate a seed

##### Use a predefined seed value
```bash
$ skycoin-cli addressGen --num 2 --seed "my super secret seed"
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
skycoin-cli addressGen --num 2 --hex
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

### Generate distribution addresses for a new fiber coin
```bash
skycoin-cli fiberAddressGen [flags]
```

```
DESCRIPTION:
    Addresses are written in a format that can be copied into fiber.toml
    for configuring distribution addresses. Addresses along with their seeds are written to a csv file,
    these seeds can be imported into the wallet to access distribution coins.

FLAGS:
  -a, --addres-file string   Output file for the generated addresses in fiber.toml format (default "addresses.txt")
  -e, --entropy int          Entropy of the autogenerated bip39 seeds. Can be 128 or 256 (default 128)
  -n, --num int              Number of addresses to generate (default 100)
  -o, --overwrite            Allow overwriting any existing addrs-file or seeds-file
  -s, --seeds-file string    Output file for the generated addresses and seeds in a csv (default "seeds.csv")
```


#### Examples
```bash
skycoin-cli fiberAddressGen
```

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
Lists the content of a single block or a range of blocks

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
                "tx_body_hash": "cf4fe76a08e3296b6f6abdb949604409be66574f211d9d14fde39103c4cfe1d6",
                "ux_hash": "d3f60f0d20aeac951aacab8d849696cac54c7057da741cfd90b63018100818d0"
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
            },
            "size": 220
        },
        {
            "header": {
                "seq": 42,
                "block_hash": "60a17e0cf411e5db7150272e597d343beaa5fbce5d61f6f647a14288262593b1",
                "previous_block_hash": "08f89cfe92be09e9848ba4d77c300908761354933f80401c107644feab1f4c9e",
                "timestamp": 1429058594,
                "fee": 292512,
                "version": 0,
                "tx_body_hash": "0e91a08561e85a36ddf44e77b9228f7d561c18c0b46d19083d4af511085b697e",
                "ux_hash": "9173768496bc49e2a34d5a7ea65d05ad6507dfdb489836e861b3c03d35efeb7a"
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
            },
            "size": 317
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
$ skycoin-cli createRawTransaction [flags] [to address] [amount]
```

```
FLAGS:
  -a, --address string          From address
  -c, --change-address string   Specify different change address.
                                By default the from address or a wallets coinbase address will be used.
      --csv  string         CSV file containing addresses and amounts to send
  -j, --json                    Returns the results in JSON format.
  -m, --many string             use JSON string to set multiple receive addresses and coins,
                                example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'
  -p, --password string         Wallet password
  -f, --wallet-file string      wallet file or path. If no path is specified your default wallet path will be used.
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

##### Sending to addresses in a CSV file
```bash
$ cat <<EOF > $CSV_FILE
2Niqzo12tZ9ioZq5vwPHMVR4g7UVpp9TCmP,123.1
2UDzBKnxZf4d9pdrBJAqbtoeH641RFLYKxd,456.045
yExu4fryscnahAEMKa7XV4Wc1mY188KvGw,0.3
EOF
$ skycoin-cli createRawTransaction -f $WALLET_PATH -a $FROM_ADDRESS -csv $CSV_FILE
```

<details>
 <summary>View Output</summary>

```
01010000000e2a5bf4964604006fea5cf8cbd705e82bebb055467f10681ef01ce5c8db654801000000d951d4e34a7b35b1b165e8302cd47e09b6433ea443a8864dc8428537dbe8b76e00ee58bb195d7de3d28935ed4fc3684f1cac5593c09c4bafb016705b7e2b3393000100000005e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634030000000007445b5d6fbbb1a7d70bef941fb5da234a10fcae40420f000000000001000000000000000056500d41a1a6f1967ffe0074bb171148667ce20d40420f00000000000100000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec09c0e1e400000000009a05000000000000
```
</details>

> NOTE: When sending to multiple addresses each combination of address and coins need to be unique
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
    "length": 220,
    "type": 0,
    "txid": "ee700309aba9b8b552f1c932a667c3701eff98e71c0e5b0e807485cea28170e5",
    "inner_hash": "247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af3843",
    "sigs": [
        "cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f00"
    ],
    "inputs": [
        "05e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634"
    ],
    "outputs": [
        {
            "uxid": "2f146924431e8c9b84a53d4d823acefb92515a264956d873ac86066c608af418",
            "dst": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
            "coins": "1.000000",
            "hours": 1
        },
        {
            "uxid": "5d69d22aff5957a18194c443557d97ec18707e4db8ee7e9a4bb8a7eef642fdff",
            "dst": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
            "coins": "16.000000",
            "hours": 1432
        }
    ]
}
```
</details>

### Encode a JSON transaction

Encode JSON Skycoin transaction.

```bash
$ skycoin-cli encodeJsonTransaction [file path or -]
```

```
FLAGS:
  -j, --json                    Returns the results in JSON format.
```

#### Example
##### Read JSON transaction from stdin
```bash
$ echo '  {
       "length": 220,
       "type": 0,
       "txid": "ee700309aba9b8b552f1c932a667c3701eff98e71c0e5b0e807485cea28170e5",
       "inner_hash": "247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af3843",
       "sigs": [
           "cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f00"
       ],
       "inputs": [
           "05e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634"
       ],
       "outputs": [
           {
               "uxid": "2f146924431e8c9b84a53d4d823acefb92515a264956d873ac86066c608af418",
               "dst": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
               "coins": "1.000000",
               "hours": 1
           },
           {
               "uxid": "5d69d22aff5957a18194c443557d97ec18707e4db8ee7e9a4bb8a7eef642fdff",
               "dst": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
               "coins": "16.000000",
               "hours": 1432
           }
       ]
   }' | skycoin-cli encodeJsonTransaction -
```

<details>
  <summary>View Output</summary>

```
dc00000000247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af384301000000cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f000100000005e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634020000000007445b5d6fbbb1a7d70bef941fb5da234a10fcae40420f00000000000100000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec090024f400000000009805000000000000
```

</details>

##### Read JSON transaction from file
```bash
$ echo '  {
       "length": 220,
       "type": 0,
       "txid": "ee700309aba9b8b552f1c932a667c3701eff98e71c0e5b0e807485cea28170e5",
       "inner_hash": "247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af3843",
       "sigs": [
           "cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f00"
       ],
       "inputs": [
           "05e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634"
       ],
       "outputs": [
           {
               "uxid": "2f146924431e8c9b84a53d4d823acefb92515a264956d873ac86066c608af418",
               "dst": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
               "coins": "1.000000",
               "hours": 1
           },
           {
               "uxid": "5d69d22aff5957a18194c443557d97ec18707e4db8ee7e9a4bb8a7eef642fdff",
               "dst": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
               "coins": "16.000000",
               "hours": 1432
           }
       ]
   }' > $FILEPATH

$ skycoin-cli encodeJsonTransaction $FILEPATH
```

<details>
  <summary>View Output</summary>

```
dc00000000247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af384301000000cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f000100000005e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634020000000007445b5d6fbbb1a7d70bef941fb5da234a10fcae40420f00000000000100000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec090024f400000000009805000000000000
```

</details>

##### Generate a JSON output
```bash
$ echo '  {
       "length": 220,
       "type": 0,
       "txid": "ee700309aba9b8b552f1c932a667c3701eff98e71c0e5b0e807485cea28170e5",
       "inner_hash": "247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af3843",
       "sigs": [
           "cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f00"
       ],
       "inputs": [
           "05e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634"
       ],
       "outputs": [
           {
               "uxid": "2f146924431e8c9b84a53d4d823acefb92515a264956d873ac86066c608af418",
               "dst": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
               "coins": "1.000000",
               "hours": 1
           },
           {
               "uxid": "5d69d22aff5957a18194c443557d97ec18707e4db8ee7e9a4bb8a7eef642fdff",
               "dst": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
               "coins": "16.000000",
               "hours": 1432
           }
       ]
   }' > $FILEPATH

$ skycoin-cli encodeJsonTransaction --json $FILEPATH
```

<details>
  <summary>View Output</summary>

```
{
    "rawtx": "dc00000000247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af384301000000cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f000100000005e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634020000000007445b5d6fbbb1a7d70bef941fb5da234a10fcae40420f00000000000100000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec090024f400000000009805000000000000"
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

### Create a wallet
Create a new skycoin wallet.

```bash
$ skycoin-cli walletCreate [flags]
```

```
FLAGS:
  -x, --crypto-type string   The crypto type for wallet encryption, can be scrypt-chacha20poly1305 or sha256-xor (default "scrypt-chacha20poly1305")
  -e, --encrypt              Create encrypted wallet.
  -h, --help                 help for walletCreate
  -l, --label string         Label used to idetify your wallet.
  -m, --mnemonic             A mnemonic seed consisting of 12 dictionary words will be generated
  -n, --num uint             Number of addresses to generate. (default 1)
  -p, --password string      Wallet password
  -r, --random               A random alpha numeric seed will be generated.
  -s, --seed string          Your seed
  -t, --type string          Wallet type. Types are "collection", "deterministic" or "bip44" (default "deterministic")
  -f, --wallet-file string   Name of wallet. The final format will be "yourName.wlt".
                             If no wallet name is specified a generic name will be selected. (default "skycoin_cli.wlt")
  -w, --wordcount uint       Number of seed words to use for mnemonic. Must be 12, 15, 18, 21 or 24 (default 12)
```

#### Examples
##### Create the default wallet
```bash
$ skycoin-cli walletCreate
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

##### Create a wallet with a random alpha numeric seed
```bash
$ skycoin-cli walletCreate -r
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

##### Create a wallet with a 12 word mnemomic seed
```bash
$ skycoin-cli walletCreate -rm
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

##### Create a wallet with a specified seed
```bash
$ skycoin-cli walletCreate -s "this is the super secret seed everyone needs but does not have"
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


##### Create more than 1 default address
```bash
$ skycoin-cli walletCreate -n 2
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

##### Create a wallet with a custom wallet name
```bash
$ skycoin-cli walletCreate -f "secret_wallet.wlt"
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

##### Create a wallet with a custom wallet label
By default the wallet label is an empty field
```bash
$ skycoin-cli walletCreate -l "cli wallet"
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

### Add addresses to a wallet
Add new addresses to a skycoin wallet.

```bash
$ skycoin-cli walletAddAddresses [flags]
```

```
FLAGS:
  -j, --json                 Returns the results in JSON format
  -n, --num uint             Number of addresses to generate (default 1)
  -p, --password string      wallet password
  -f, --wallet-file string   Generate addresses in the wallet (default "$HOME/.skycoin/wallets/skycoin_cli.wlt")
```

##### Add an address to the default wallet
```bash
$ skycoin-cli walletAddAddresses
```

<details>
 <summary>View Output</summary>

```
2mEgmYt6NZHA1erYqbAeXmGPD5gqLZ9toFv
```
</details>

##### Add an address to a specific wallet
```bash
$ skycoin-cli walletAddAddresses $WALLET_PATH
```

<details>
 <summary>View Output</summary>

```json
2cET6L4c6Bee5jucuzsTQUXFxWX76GZoDqv
```
</details>

##### Add `n` addresses
```bash
$ skycoin-cli walletAddAddresses -n 2
```

<details>
 <summary>View Output</summary>

```
2UrEV3Vyu5RJABZNukKRq25ggrrg96RUwdH,LJN5qGmLbJxLswzD3nFn3RFcmWJyZ2LGHY
```
</details>

##### Add an address to a wallet with JSON output
```bash
$ skycoin-cli walletAddAddresses --json
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

### Encrypt Wallet
Encrypt a wallet seed

```bash
$ skycoin-cli encryptWallet [flags]
```

```
FLAGS:
  -x, --crypto-type string   The crypto type for wallet encryption, can be scrypt-chacha20poly1305 or sha256-xor
  -h, --help                 help for encryptWallet
  -p, --password string      wallet password
```

### Examples
#### Encrypt wallet
```bash
$ skycoin-cli encryptWallet -p test
```

<details>
 <summary>View Output</summary>

 ```json
 {
     "meta": {
         "coin": "skycoin",
         "cryptoType": "scrypt-chacha20poly1305",
         "encrypted": "true",
         "filename": "skycoin_cli.wlt",
         "label": "",
         "lastSeed": "",
         "secrets": "dgB7Im4iOjEwNDg1NzYsInIiOjgsInAiOjEsImtleUxlbiI6MzIsInNhbHQiOiJRNVRSVHh0VFpieERpUWt0dnkzc01SYTl6U0t2aFJqVlpUUHQzeldSVGs4PSIsIm5vbmNlIjoiSUt5VG8zdWdGdFY3MWYxTiJ9LB7Cu3bvZFzsmKqToPi3bjARIRfmhL8HBUdnwLzS5Rxu4uw1tIlDDmEKUpgDWV3RvB+xDz3sHchQr5BpK72LDOwbZ6BubMHovTqC4+lx9hKc2qnDGwsymxLQJHQrQ23DkHMioSUVYNZv1/DwzJ2qI0WIOTkb+L34e9f60YV+2zF7v+C/nTS8AjMwjGYldKinPEjyDXkpxB2d4Sd3EnfUm8u76TvTKxqZpZ/tr+in/OfRsJsN7dC7rMFRZukoCJYNnWv/wgPn/NMu4DIxqF+WUQhCsCgqk6oMderdK/E/xtLJmKnbHRLH4PO/Dh4ypLXg2EzW+JBN6RpzVEXxYdvVCqmKfs7d+hnHWDmDtCLGqYyPsUa+d4PPhylruNE=",
         "seed": "",
         "tm": "1540305209",
         "type": "deterministic",
         "version": "0.2"
     },
     "entries": [
         {
             "address": "2gvvvS5jziMDQTUPB98LFipCTDjm1H723k2",
             "public_key": "032fe2ceacabc1a6acad8c93bd3493a3570fb76a9f8dc625dd200d13f96abed3e0",
             "secret_key": ""
         }
     ]
 }
 ```
</details>


#### Encrypt wallet with different crypto type
```bash
$ skycoin-cli encryptWallet -x sha256-xor -p test
```

<details>
 <summary>View Output</summary>

 ```json
 {
     "meta": {
         "coin": "skycoin",
         "cryptoType": "sha256-xor",
         "encrypted": "true",
         "filename": "skycoin_cli.wlt",
         "label": "",
         "lastSeed": "",
         "secrets": "mJ4g+/NgncOVp7gKIZqVPysmrRYKjprSuMvvpq3HLt7ajjMOheEdyU0PGtueDQADIhhTFZlQh/eaaYXF3fecS7OrGa79F+2lRRdD7Tva/MueiL9TL0ng12x0I7dXkUVsXLTl3MJK27JwS9hKedcVvnmFysJA6W3lX2aE7Qn+v6cyMbfgR8r89OHGaUZ9SPZn2HKOhhIcXt66Q/t0kVWU0XEH+G
 xUyX23ksN3scQoAshVidLAgXwpkgExEl+qjCpDNQga3MncZV+WuQxpIKodJ3l5TKoJAA0/Taz9O9Se0tIoiK2ls2m6JUayev3Id0+hkmNNSUKQ53Ni3xwjNzZXoPQAemMWpkdUSv8qNuhh7C/4gBBrZROM6ZyxmsdlWgcG0Yfrh8o505D0i4mtubkdZSGi8Djm9j1mpWTZi3VuUjtGvBAmH3Qzdma+nvORZj11QuEuCcO+
 8jmQB9bVxcTL9u4Nan2+cYijVNul93m7xWik/mSB7uIFVIJAm4kSMiJm",
         "seed": "",
         "tm": "1540305209",
         "type": "deterministic",
         "version": "0.2"
     },
     "entries": [
         {
             "address": "2gvvvS5jziMDQTUPB98LFipCTDjm1H723k2",
             "public_key": "032fe2ceacabc1a6acad8c93bd3493a3570fb76a9f8dc625dd200d13f96abed3e0",
             "secret_key": ""
         }
     ]
 }
 ```
</details>

### Decrypt Wallet
Decrypt a wallet seed

```bash
$ skycoin-cli decryptWallet [flags]
```

```
FLAGS:
  -h, --help              help for decryptWallet
  -p, --password string   wallet password
```

### Example
```bash
$ skycoin-cli decryptWallet -p test
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
         "lastSeed": "522dba68fe58c179f3467f9e799c02b25552143b250626cc03281faa28c262c0",
         "secrets": "",
         "seed": "select salute trip target blur short link suspect river ready senior bleak",
         "tm": "1540305209",
         "type": "deterministic",
         "version": "0.2"
     },
     "entries": [
         {
             "address": "2gvvvS5jziMDQTUPB98LFipCTDjm1H723k2",
             "public_key": "032fe2ceacabc1a6acad8c93bd3493a3570fb76a9f8dc625dd200d13f96abed3e0",
             "secret_key": "080bfb86463da87e06f816c4326a11b84806c9744235bb7ce7bc8d63acb4f6c2"
         }
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
                "seq": 58894,
                "block_hash": "3961bea8c4ab45d658ae42effd4caf36b81709dc52a5708fdd4c8eb1b199a1f6",
                "previous_block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
                "timestamp": 1537581604,
                "fee": 485194,
                "version": 0,
                "tx_body_hash": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
                "ux_hash": "f7d30ecb49f132283862ad58f691e8747894c9fc241cb3a864fc15bd3e2c83d3"
            },
            "body": {
                "txns": [
                    {
                        "length": 257,
                        "type": 0,
                        "txid": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
                        "inner_hash": "f7dbd09f7e9f65d87003984640f1977fb9eec95b07ef6275a1ec6261065e68d7",
                        "sigs": [
                            "af5329e77213f34446a0ff41d249fd25bc1dae913390871df359b9bd587c95a10b625a74a3477a05cc7537cb532253b12c03349ead5be066b8e0009e79462b9501"
                        ],
                        "inputs": [
                            "fb8db3f78928aee3f5cbda8db7fc290df9e64414e8107872a1c5cf83e08e4df7"
                        ],
                        "outputs": [
                            {
                                "uxid": "235811602fc96cf8b5b031edb88ee1606830aa641c06e0986681552d8728ec07",
                                "dst": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
                                "coins": "0.500000",
                                "hours": 1
                            },
                            {
                                "uxid": "873da4edc01c0b5184e1f26c4c3471dd407d08e9ab36b018ab93874e7392320b",
                                "dst": "2XBMMDMqTTYmqs2rfjEwYDz8ABd38y9B8r7",
                                "coins": "0.500000",
                                "hours": 1
                            },
                            {
                                "uxid": "42a6f0127f61e1d7bca8e9680027eddcecad772250c5634a03e56a8b1cf5a816",
                                "dst": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "25.913000",
                                "hours": 485192
                            }
                        ]
                    }
                ]
            },
            "size": 257
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
                "seq": 58892,
                "block_hash": "1f042ed976c0cb150ea6b71c9608d65b519e4bc1c507eba9f1146e443a856c2d",
                "previous_block_hash": "d9ca9442febd8788de0a3093158943beca228017bf8c9c9b8529a382fad8d991",
                "timestamp": 1537580914,
                "fee": 94694,
                "version": 0,
                "tx_body_hash": "9895f8af790e33a618004dc61f48ecc16bd642751a3fff6b05cecb8815c80942",
                "ux_hash": "bb188dcaaf28613d49b926636675dacf67a739a4e316253b1207ad674709252b"
            },
            "body": {
                "txns": [
                    {
                        "length": 1190,
                        "type": 0,
                        "txid": "9895f8af790e33a618004dc61f48ecc16bd642751a3fff6b05cecb8815c80942",
                        "inner_hash": "8bff0b7572bb49ccde4b2b313e921e5cf302a11fd9f786a2ef97a7c0ddfee261",
                        "sigs": [
                            "f843861b301eb025e58bacfb934d615f263419704b0a59f2645845344f2702fa1a7a967651f01933af4d56752c656a7e759e942b9278e228362f2ef273d4ff5200",
                            "06f15e2522e7413f25dedb5aee67ae880bd98bb7df11a1a92241d88db9bb976d2c707e77a4a3ddfd8d123ad04701fe2538ea2d0f78cfbcdc44e70fd2320a72b500",
                            "dc32fe308274f9aaa21e09046384a83b4b2c6bf800c6b9ff492af8bd3f5cd7717b245c9d460c242139034c73cd15aca9f288cb69e8ae4c33df2a807ff3b373aa00",
                            "18d83a122f0ca3629f0c82e21ea3d6fbdfd1ea07ba062ffb6647b7e2c3aa9a1d7c112dc5543435ddd0bccd163f839d9802eb344f6203372deea7402d8476679501",
                            "26a8723c1ac22dca2b61d807ca2279e341a9f5a371c4d14333f49e52b90ec87f08ab7930e5804367c1dffb01b197d976619ab26f0c8afe8837c41b0df809a23301",
                            "fe190749475cd66afcdd295b22b007c63726db0fd834acef4ecde9c41ae7d15d54c2c180c8aba5c894d3843405f6243b7ff964f974f607b38298c195d7b523c401",
                            "d58e8283a28faed377161605e252cd929959e40fd8b996f928049f2b446bc920468d1086a2bc34a8fcaedefdc471427266cc67b9770d9b0482f5f4f22729a79a00",
                            "2be852a5b589ce49f9d3678fa44a758c2e4d7372620a8241d71c41451c5244387ac337bbee5010b98fc8c38fc44619ed8a7beb2af06385a11ecb53eb0112a8a700",
                            "b6376cc54078f775da6438960e828c799c780349c8590508b4500f0e6dd9ecbc760992599d698209b078220d8aaa9db9f80091290a18cd0753efd1805515d06600",
                            "958e17753d4cedc3203b95a39d585314ac10efc00332befa81c8049b4178222d2514ba1d68385b2518d976835dee48f2bb540db0d0e728acbf59d8562cbb7baa00",
                            "d921e2aa2b1b6778a84efdc7f1d016c7aad66dfc13c0be4fee6a5f303a2c3cf465fa0d549ca5fc57d3a26832bffcaec842837905a78e8ca3fa553522d931571d01"
                        ],
                        "inputs": [
                            "c551da99c0b74b64511aaaf99536cb6d263958064890ef6c27be36e8f5a14fb8",
                            "64875d950120b16d0f0f84c708e3e48b26fb9c32f36c0fba71764dfc53e7ae05",
                            "ba50cb14fc26bcf658ace9a3b5d6e0d257fa022e80613902c693ab57a1c0924a",
                            "b388fdb6dc7c91cab7e72a4786967e18834350c1ccd149790a0a2270bdf91bf9",
                            "94f87596cb7471e2b96b7e1ddd8194d44ca4858d97ed29f83e926bcdf36601c1",
                            "dd2a4d942ba1ad4dc55f46adc31e3f62e7326b9a0d511f1faf83911af77018f1",
                            "35b82981a9648ba871c2c08604cb95f130baafa26761139c7134f5b9e4575b9d",
                            "aa0f74c067347b0178c6963d8409c6bdf7a39051641f9ba1a5d5c37d88dce7f0",
                            "5a34d07536c2964338aa89f19ab7ff857056f6ffd16e936ae13229077387afb5",
                            "dc93bb4a131cea3d3f2b523408f077779384c816a4516dfbe0817845938a26ef",
                            "53f92392b71ce79ead8452e5c31c8a404acf9770a71d4dc234f2fe54a8671495"
                        ],
                        "outputs": [
                            {
                                "uxid": "061a639996b85d2c0f19cf929a83c5abe2667a411de31fbdbd16c1da6c8e4880",
                                "dst": "2gXHek83jtEdDndgrKkEwgwZZDsHXKfNaD",
                                "coins": "87.990000",
                                "hours": 1
                            },
                            {
                                "uxid": "f210a8ed58c92094832ccb4d5a4ae7271df1df0d7176b18d5c7b149ed36a7d80",
                                "dst": "27ckSMTwxMxHanUM1VmF8BV9JuWdQd4Gd9S",
                                "coins": "0.010000",
                                "hours": 94693
                            }
                        ]
                    }
                ]
            },
            "size": 1190
        },
        {
            "header": {
                "seq": 58893,
                "block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
                "previous_block_hash": "1f042ed976c0cb150ea6b71c9608d65b519e4bc1c507eba9f1146e443a856c2d",
                "timestamp": 1537581594,
                "fee": 970389,
                "version": 0,
                "tx_body_hash": "1bea5cf1279693a0da24828c37b267c702007842b16ca5557ae497574d15aab7",
                "ux_hash": "bf35652af199779bc40cbeb339e8a782ff70673b07779e5c5621d37dfe13b42b"
            },
            "body": {
                "txns": [
                    {
                        "length": 377,
                        "type": 0,
                        "txid": "1bea5cf1279693a0da24828c37b267c702007842b16ca5557ae497574d15aab7",
                        "inner_hash": "a25232405bcef0c007bb2d7d3520f2a389e17e11125c252ab6c00168ec52c08d",
                        "sigs": [
                            "2ff7390c3b66c6b0fbb2b4c59c8e218291d4cbb82a836bb577c7264677f4a8320f6f3ad72d804e3014728baa214c223ecced8725b64be96fe3b51332ad1eda4201",
                            "9e7c715f897b3c987c00ee8c6b14e4b90bb3e4e11d003b481f82042b1795b3c75eaa3d563cd0358cdabdab77cfdbead7323323cf73e781f9c1a8cf6d9b4f8ac100",
                            "5c9748314f2fe0cd442df5ebb8f211087111d22e9463355bf9eee583d44df1bd36addb510eb470cb5dafba0732615f8533072f80ae05fc728c91ce373ada1e7b00"
                        ],
                        "inputs": [
                            "5f634c825b2a53103758024b3cb8578b17d56d422539e23c26b91ea397161703",
                            "16ac52084ffdac2e9169b9e057d44630dec23d18cfb90b9437d28220a3dc585d",
                            "8d3263890d32382e182b86f8772c7685a8f253ed475c05f7d530e9296f692bc9"
                        ],
                        "outputs": [
                            {
                                "uxid": "fb8db3f78928aee3f5cbda8db7fc290df9e64414e8107872a1c5cf83e08e4df7",
                                "dst": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "26.913000",
                                "hours": 970388
                            }
                        ]
                    }
                ]
            },
            "size": 377
        },
        {
            "header": {
                "seq": 58894,
                "block_hash": "3961bea8c4ab45d658ae42effd4caf36b81709dc52a5708fdd4c8eb1b199a1f6",
                "previous_block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
                "timestamp": 1537581604,
                "fee": 485194,
                "version": 0,
                "tx_body_hash": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
                "ux_hash": "f7d30ecb49f132283862ad58f691e8747894c9fc241cb3a864fc15bd3e2c83d3"
            },
            "body": {
                "txns": [
                    {
                        "length": 257,
                        "type": 0,
                        "txid": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
                        "inner_hash": "f7dbd09f7e9f65d87003984640f1977fb9eec95b07ef6275a1ec6261065e68d7",
                        "sigs": [
                            "af5329e77213f34446a0ff41d249fd25bc1dae913390871df359b9bd587c95a10b625a74a3477a05cc7537cb532253b12c03349ead5be066b8e0009e79462b9501"
                        ],
                        "inputs": [
                            "fb8db3f78928aee3f5cbda8db7fc290df9e64414e8107872a1c5cf83e08e4df7"
                        ],
                        "outputs": [
                            {
                                "uxid": "235811602fc96cf8b5b031edb88ee1606830aa641c06e0986681552d8728ec07",
                                "dst": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
                                "coins": "0.500000",
                                "hours": 1
                            },
                            {
                                "uxid": "873da4edc01c0b5184e1f26c4c3471dd407d08e9ab36b018ab93874e7392320b",
                                "dst": "2XBMMDMqTTYmqs2rfjEwYDz8ABd38y9B8r7",
                                "coins": "0.500000",
                                "hours": 1
                            },
                            {
                                "uxid": "42a6f0127f61e1d7bca8e9680027eddcecad772250c5634a03e56a8b1cf5a816",
                                "dst": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "25.913000",
                                "hours": 485192
                            }
                        ]
                    }
                ]
            },
            "size": 257
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

### Rich list
Returns the top N address (default 20) balances (based on unspent outputs). Optionally include distribution addresses (exluded by default).

```bash
$ skycoin-cli richlist [top N addresses] [include distribution addresses]
```

#### Example
```bash
$ skycoin-cli richlist 5 false
```

<details>
 <summary>View Output</summary>

```json
{
    "richlist": [
        {
            "address": "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8",
            "coins": "1072264.838000",
            "locked": false
        },
        {
            "address": "2fGi2jhvp6ppHg3DecguZgzqvpJj2Gd4KHW",
            "coins": "500000.000000",
            "locked": false
        },
        {
            "address": "2jNwfvZNUoRLiFzJtmnevSF6TKPfSehvrc1",
            "coins": "252297.068000",
            "locked": false
        },
        {
            "address": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
            "coins": "236884.364000",
            "locked": false
        },
        {
            "address": "2fR8BkeTRQC4R3ATNnujHsQQXcaB6m4Aqwo",
            "coins": "173571.990000",
            "locked": false
        }
    ]
}
```
</details>

### Send
Make a skycoin transaction.

```bash
$ skycoin-cli send [flags] [to address] [amount]
```

```
FLAGS:
  -a, --address string          From address
  -c, --change-address string   Specify different change address.
                                By default the from address or a wallets coinbase address will be used.
      --csv  string         CSV file containing addresses and amounts to send
  -j, --json                    Returns the results in JSON format.
  -m, --many string             use JSON string to set multiple receive addresses and coins,
                                example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'
  -p, --password string         Wallet password
  -f, --wallet-file string      wallet file or path. If no path is specified your default wallet path will be used.
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

##### Sending to addresses in a CSV file
```bash
$ cat <<EOF > $CSV_FILE
2Niqzo12tZ9ioZq5vwPHMVR4g7UVpp9TCmP,123.1
2UDzBKnxZf4d9pdrBJAqbtoeH641RFLYKxd,456.045
yExu4fryscnahAEMKa7XV4Wc1mY188KvGw,0.3
EOF
$ skycoin-cli send -f $WALLET_PATH -a $FROM_ADDRESS -csv $CSV_FILE
```

<details>
 <summary>View Output</summary>

```
txid:$TRANSACTION_ID
```
</details>

> NOTE: When sending to multiple addresses each combination of address and coins need to be unique
        Otherwise you get, `ERROR: Duplicate output in transaction`

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

### Show Seed
Show seed of a specified wallet.
The default wallet `($HOME/wallets/skycoin_cli.wlt)` will be used if no wallet was specified.
The wallet file is configured through `WALLET_NAME` env var.


```bash
$ skycoin-cli showSeed [flags]
```

```
FLAGS:
  -j, --json                 Returns the results in JSON format.
  -p, --password string      Wallet password
```

#### Example
```bash
$ skycoin-cli showSeed
```
<details>
 <summary>View Output</summary>
 ```
 eternal turtle seek nominee narrow much melody kite worth giggle shrimp horse
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
    "rpc_address": "http://127.0.0.1:6420"
}
```
</details>

### Status
#### Example
```bash
$ skycoin-cli status
```

<details>
 <summary>View Output</summary>

```json
{
    "status": {
        "blockchain": {
            "head": {
                "seq": 58894,
                "block_hash": "3961bea8c4ab45d658ae42effd4caf36b81709dc52a5708fdd4c8eb1b199a1f6",
                "previous_block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
                "timestamp": 1537581604,
                "fee": 485194,
                "version": 0,
                "tx_body_hash": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
                "ux_hash": "f7d30ecb49f132283862ad58f691e8747894c9fc241cb3a864fc15bd3e2c83d3"
            },
            "unspents": 38171,
            "unconfirmed": 1,
            "time_since_last_block": "7m44s"
        },
        "version": {
            "version": "0.25.0",
            "commit": "620405485d3276c16c0379bc3b88b588e34c45e1",
            "branch": "develop"
        },
	    "coin": "skycoin",
	    "user_agent": "skycoin:0.25.0",
        "open_connections": 8,
        "outgoing_connections": 5,
        "incoming_connections": 3,
        "uptime": "4h1m23.697072461s",
        "csrf_enabled": true,
        "csp_enabled": true,
        "wallet_api_enabled": true,
        "gui_enabled": true,
        "user_verify_transaction": {
            "burn_factor": 10,
            "max_transaction_size": 32768,
            "max_decimals": 3
        },
        "unconfirmed_verify_transaction": {
            "burn_factor": 10,
            "max_transaction_size": 32768,
            "max_decimals": 3
        },
        "started_at": 1558864387,
        "fiber": {
            "name": "skycoin",
            "display_name": "Skycoin",
            "ticker": "SKY",
            "coin_hours_display_name": "Coin Hours",
            "coin_hours_display_name_singular": "Coin Hour",
            "coin_hours_ticker": "SCH",
            "explorer_url": "https://explorer.skycoin.net"
        }
    },
    "cli_config": {
        "webrpc_address": "http://127.0.0.1:6420"
    }
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
         "block_seq": 864
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

### Get address transactions
Get transaction for one or more addresses - including listing of both inputs and outputs.

```bash
$ skycoin-cli addressTransactions [addr1 addr2 addr3]
```

#### Example
#### Single Address
```bash
$ skycoin-cli addressTransactions 21YPgFwkLxQ1e9JTCZ43G7JUyCaGRGqAsda
```

<details>
 <summary>View Output</summary>

```json
[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 66119,
            "block_seq": 21213
        },
        "time": 1523180676,
        "txn": {
            "timestamp": 1523180676,
            "length": 220,
            "type": 0,
            "txid": "8cdf82ec42e8316007ed99c0b1de1d0dfd9221c757f41fdec0b36009df74085f",
            "inner_hash": "c543f08bfe7b99a19f7bc4068a02e437ed4a043130e976551188c4d38b89ce8d",
            "fee": 726,
            "sigs": [
                "f1021744902892eb47c60f7240ce6964de3c7bf77777ce267b58df8879e208e57bd044d15a36d78bebab2897c2c61ecbbceb348cfc45152efb105960799364c401"
            ],
            "inputs": [
                {
                    "uxid": "5d69d22aff5957a18194c443557d97ec18707e4db8ee7e9a4bb8a7eef642fdff",
                    "owner": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
                    "coins": "16.000000",
                    "hours": 1432,
                    "calculated_hours": 1452
                }
            ],
            "outputs": [
                {
                    "uxid": "0020ae8da2bcc7657f3b234cbb59e0fd2486c53d7ef3f05cda6ff613587c8441",
                    "dst": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
                    "coins": "1.000000",
                    "hours": 1
                },
                {
                    "uxid": "9d79ad07a90fee10b59bea1bd6f566f0b69f6bf9a9e735c1bec4b0e5eb4b33cb",
                    "dst": "21YPgFwkLxQ1e9JTCZ43G7JUyCaGRGqAsda",
                    "coins": "15.000000",
                    "hours": 725
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 66111,
            "block_seq": 21221
        },
        "time": 1523184376,
        "txn": {
            "timestamp": 1523184376,
            "length": 183,
            "type": 0,
            "txid": "f3c5cfd462d95e724b7d35b1688c53f25a5f358f2eb9a6f87b63cdf31deb2bf8",
            "inner_hash": "8269589c228be4bc33d75f6ee5b334856e8680b7d6ec275f897406c01da8340b",
            "fee": 370,
            "sigs": [
                "33879494d644df45b5c6c7111c0e453cd42f6fe718614a9411d9fbabd57ab24749813cdf47424dcac5ed097a0de0ac7b557154d2ec93f81b12b1dfdee5138df701"
            ],
            "inputs": [
                {
                    "uxid": "9d79ad07a90fee10b59bea1bd6f566f0b69f6bf9a9e735c1bec4b0e5eb4b33cb",
                    "owner": "21YPgFwkLxQ1e9JTCZ43G7JUyCaGRGqAsda",
                    "coins": "15.000000",
                    "hours": 725,
                    "calculated_hours": 739
                }
            ],
            "outputs": [
                {
                    "uxid": "c51b2692aa9f296a3cd2f37b14f39c496c82f5c5ae01c54701ea60b7353f27e2",
                    "dst": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
                    "coins": "15.000000",
                    "hours": 369
                }
            ]
        }
    }
]
```
</details>

#### Multiple Address
```bash
$ skycoin-cli addressTransactions 21YPgFwkLxQ1e9JTCZ43G7JUyCaGRGqAsda 3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd
```

<details>
 <summary>View Output</summary>

```json
[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 66143,
            "block_seq": 21189
        },
        "time": 1523176026,
        "txn": {
            "timestamp": 1523176026,
            "length": 220,
            "type": 0,
            "txid": "ee700309aba9b8b552f1c932a667c3701eff98e71c0e5b0e807485cea28170e5",
            "inner_hash": "247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af3843",
            "fee": 1442,
            "sigs": [
                "cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f00"
            ],
            "inputs": [
                {
                    "uxid": "05e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634",
                    "owner": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
                    "coins": "17.000000",
                    "hours": 139,
                    "calculated_hours": 2875
                }
            ],
            "outputs": [
                {
                    "uxid": "2f146924431e8c9b84a53d4d823acefb92515a264956d873ac86066c608af418",
                    "dst": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
                    "coins": "1.000000",
                    "hours": 1
                },
                {
                    "uxid": "5d69d22aff5957a18194c443557d97ec18707e4db8ee7e9a4bb8a7eef642fdff",
                    "dst": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
                    "coins": "16.000000",
                    "hours": 1432
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 66142,
            "block_seq": 21190
        },
        "time": 1523176126,
        "txn": {
            "timestamp": 1523176126,
            "length": 183,
            "type": 0,
            "txid": "8c137774a2485beeaa3f8e861097ba6dffb144fb2c2f2c357c9261a324b02013",
            "inner_hash": "92da4c2d6e93a6f0a62899225a9195b95eb274f8e926b0a2ce5d259f84015014",
            "fee": 1,
            "sigs": [
                "84f9c7b5d1f88245b53d50e4e8d4fd8719089768940a4ff9d8c3d88b15c300e57f91fa07a0789bbfac8e7c77aebda83d39c6b77aa80cd70a613bf175c316b6cc00"
            ],
            "inputs": [
                {
                    "uxid": "2f146924431e8c9b84a53d4d823acefb92515a264956d873ac86066c608af418",
                    "owner": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
                    "coins": "1.000000",
                    "hours": 1,
                    "calculated_hours": 1
                }
            ],
            "outputs": [
                {
                    "uxid": "5250017c47070e011cc71c44472d5ab8e957c25c9c57fc7885e0a4301c7c014c",
                    "dst": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
                    "coins": "1.000000",
                    "hours": 0
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 66119,
            "block_seq": 21213
        },
        "time": 1523180676,
        "txn": {
            "timestamp": 1523180676,
            "length": 220,
            "type": 0,
            "txid": "8cdf82ec42e8316007ed99c0b1de1d0dfd9221c757f41fdec0b36009df74085f",
            "inner_hash": "c543f08bfe7b99a19f7bc4068a02e437ed4a043130e976551188c4d38b89ce8d",
            "fee": 726,
            "sigs": [
                "f1021744902892eb47c60f7240ce6964de3c7bf77777ce267b58df8879e208e57bd044d15a36d78bebab2897c2c61ecbbceb348cfc45152efb105960799364c401"
            ],
            "inputs": [
                {
                    "uxid": "5d69d22aff5957a18194c443557d97ec18707e4db8ee7e9a4bb8a7eef642fdff",
                    "owner": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
                    "coins": "16.000000",
                    "hours": 1432,
                    "calculated_hours": 1452
                }
            ],
            "outputs": [
                {
                    "uxid": "0020ae8da2bcc7657f3b234cbb59e0fd2486c53d7ef3f05cda6ff613587c8441",
                    "dst": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
                    "coins": "1.000000",
                    "hours": 1
                },
                {
                    "uxid": "9d79ad07a90fee10b59bea1bd6f566f0b69f6bf9a9e735c1bec4b0e5eb4b33cb",
                    "dst": "21YPgFwkLxQ1e9JTCZ43G7JUyCaGRGqAsda",
                    "coins": "15.000000",
                    "hours": 725
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 66118,
            "block_seq": 21214
        },
        "time": 1523180976,
        "txn": {
            "timestamp": 1523180976,
            "length": 220,
            "type": 0,
            "txid": "be67302e8f6f579423ba38be29de0de19815ec3c91352c6540e5f75439eb9f22",
            "inner_hash": "ef091437da13980547e33aa8647cdd1462384ec73cd57caf289e5410e3a96cf0",
            "fee": 2243,
            "sigs": [
                "d8636af89bf7f7c6aeaf32a994f8efc6e62bc25bd4e2d7b0a4deeb1e0e2888c234895f978e051985964f8b522e7d68794b90d6404809464d6c86af7153d5896e01"
            ],
            "inputs": [
                {
                    "uxid": "c981f19ff129c5746940cbf4e57383bcdf524a02055219c629e5fc4ff74067ab",
                    "owner": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
                    "coins": "3.000000",
                    "hours": 111,
                    "calculated_hours": 4486
                }
            ],
            "outputs": [
                {
                    "uxid": "ba74051563bbe6aac1836780770a66bf782a4b3a90c5ea341b43cb85a7f9d51b",
                    "dst": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
                    "coins": "1.000000",
                    "hours": 1
                },
                {
                    "uxid": "80ad81c7de66f2839b24896340890c77a79b8409abdf8e9956f5e3b65baa545b",
                    "dst": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
                    "coins": "2.000000",
                    "hours": 2242
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 66117,
            "block_seq": 21215
        },
        "time": 1523181146,
        "txn": {
            "timestamp": 1523181146,
            "length": 183,
            "type": 0,
            "txid": "5b6318a95f32487a6340f35a03cd46cba8c87d261e80ad3106a0e67d4cd4601b",
            "inner_hash": "33144c33224f1a59f75fba415a67834260e7253958d7130a0e9c0fe342ff608e",
            "fee": 1,
            "sigs": [
                "231ac8febcb4b34f6742e2c6b20690c09acffea135707fb5b6679b9cf943b9b529a06cb161e3b51d0c37e5126ce9dbf59e87eaeac511ae06d2beca5d2300611500"
            ],
            "inputs": [
                {
                    "uxid": "0020ae8da2bcc7657f3b234cbb59e0fd2486c53d7ef3f05cda6ff613587c8441",
                    "owner": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
                    "coins": "1.000000",
                    "hours": 1,
                    "calculated_hours": 1
                }
            ],
            "outputs": [
                {
                    "uxid": "2a5d9458199c977779347d160f7db4978059c70217c44f8fc34716be43b7c6f1",
                    "dst": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
                    "coins": "1.000000",
                    "hours": 0
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 66112,
            "block_seq": 21220
        },
        "time": 1523184176,
        "txn": {
            "timestamp": 1523184176,
            "length": 183,
            "type": 0,
            "txid": "4acd61d7aa7dfe20795e517d7560643d049036af9451bcbd762793bcb6a4a6de",
            "inner_hash": "c01a389f1018cf41d4ef36d550162999d82211f24f3d8b2cbf40a88edfaf690b",
            "fee": 1,
            "sigs": [
                "8ce6eff33887a8c2e31b669138163e2bcc2161782754d79c3a4c6839b4cf1fbc5a7d5e0576060d0378fbd9ee5c0f4863f949c77e7f724a4d66d75b2aed9123ae00"
            ],
            "inputs": [
                {
                    "uxid": "ba74051563bbe6aac1836780770a66bf782a4b3a90c5ea341b43cb85a7f9d51b",
                    "owner": "3vbfHxPzMuyFJvgHdAoqmFnyg6k8HiLyxd",
                    "coins": "1.000000",
                    "hours": 1,
                    "calculated_hours": 1
                }
            ],
            "outputs": [
                {
                    "uxid": "a0777af14223bbbd5aeb8bf3cfd6ba94c776c6eec731310caaaaee49b9feb9a5",
                    "dst": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
                    "coins": "1.000000",
                    "hours": 0
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 66111,
            "block_seq": 21221
        },
        "time": 1523184376,
        "txn": {
            "timestamp": 1523184376,
            "length": 183,
            "type": 0,
            "txid": "f3c5cfd462d95e724b7d35b1688c53f25a5f358f2eb9a6f87b63cdf31deb2bf8",
            "inner_hash": "8269589c228be4bc33d75f6ee5b334856e8680b7d6ec275f897406c01da8340b",
            "fee": 370,
            "sigs": [
                "33879494d644df45b5c6c7111c0e453cd42f6fe718614a9411d9fbabd57ab24749813cdf47424dcac5ed097a0de0ac7b557154d2ec93f81b12b1dfdee5138df701"
            ],
            "inputs": [
                {
                    "uxid": "9d79ad07a90fee10b59bea1bd6f566f0b69f6bf9a9e735c1bec4b0e5eb4b33cb",
                    "owner": "21YPgFwkLxQ1e9JTCZ43G7JUyCaGRGqAsda",
                    "coins": "15.000000",
                    "hours": 725,
                    "calculated_hours": 739
                }
            ],
            "outputs": [
                {
                    "uxid": "c51b2692aa9f296a3cd2f37b14f39c496c82f5c5ae01c54701ea60b7353f27e2",
                    "dst": "tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V",
                    "coins": "15.000000",
                    "hours": 369
                }
            ]
        }
    }
]
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
$ skycoin-cli walletDir [flags]
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
$ skycoin-cli walletHistory [flags]
```

```
FLAGS:
        -f value  wallet file or path. If no path is specified your default wallet path will be used.
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

### Richlist
Returns top N address (default 20) balances (based on unspent outputs). Optionally include distribution addresses (exluded by default).

```bash
$ skycoin-cli richlist [top N addresses (20 default)] [include distribution addresses (false default)]
```

```
FLAGS:
  -h, --help   help for richlist
```

#### Example
##### Without distribution addresses
```bash
$ skycoin-cli richlist 2
```
<details>
 <summary>View Output</summary>

```json
{
    "richlist": [
        {
            "address": "zVzkqNj3Ueuzo54sbACcYBqqGBPCGAac5W",
            "coins": "2922927.299000",
            "locked": false
        },
        {
            "address": "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8",
            "coins": "675256.308000",
            "locked": false
        }
    ]
}
```
</details>

##### Including distribution addresses
```bash
$ skycoin-cli richlist 2 true
```

<details>
 <summary>View Output</summary>

```json
{
    "richlist": [
        {
            "address": "zVzkqNj3Ueuzo54sbACcYBqqGBPCGAac5W",
            "coins": "2922927.299000",
            "locked": false
        },
        {
            "address": "ejJjiCwp86ykmFr5iTJ8LxQXJ2wJPTYmkm",
            "coins": "1000000.010000",
            "locked": true
        }
    ]
}
```
</details>

### Address Count
Returns the count of all addresses that currently have unspent outputs (coins) associated with them.

```bash
$ skycoin-cli addresscount
```

```
FLAGS:
  -h, --help   help for richlist
```

#### Example
```bash
$ skycoin-cli addresscount
```
<details>
 <summary>View Output</summary>

```json
12961
```
</details>


### CLI version
Get version of current skycoin cli.

```bash
$ skycoin-cli version [flags]
```

```
FLAGS:
  -j, --json   Returns the results in JSON format
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
