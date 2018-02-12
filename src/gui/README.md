# APIs document

Apis service port is `6420`.

<!-- MarkdownTOC autolink="true" bracket="round" -->

- [CSRF](#csrf)
    - [Get current csrf token](#get-current-csrf-token)
- [Simple query apis](#simple-query-apis)
    - [Get node version info](#get-node-version-info)
    - [Get balance of addresses](#get-balance-of-addresses)
    - [Get unspent output set of address or hash](#get-unspent-output-set-of-address-or-hash)
- [Wallet apis](#wallet-apis)
    - [Generate wallet seed](#generate-wallet-seed)
    - [Create a wallet from seed](#create-a-wallet-from-seed)
    - [Generate new address in wallet](#generate-new-address-in-wallet)
    - [Updates wallet label](#updates-wallet-label)
    - [Get wallet balance](#get-wallet-balance)
    - [Spend coins from wallet](#spend-coins-from-wallet)
- [Transaction apis](#transaction-apis)
    - [Get unconfirmed transactions](#get-unconfirmed-transactions)
    - [Get transaction info by id](#get-transaction-info-by-id)
    - [Get raw transaction by id](#get-raw-transaction-by-id)
    - [Inject raw transaction](#inject-raw-transaction)
    - [Get transactions that are addresses related](#get-transactions-that-are-addresses-related)
- [Block apis](#block-apis)
    - [Get blochchain progress](#get-blochchain-progress)
    - [Get block by hash or seq](#get-block-by-hash-or-seq)
    - [Get blocks in specific range](#get-blocks-in-specific-range)
    - [Get last N blocks](#get-last-n-blocks)
- [Explorer apis](#explorer-apis)
    - [Get address affected transactions](#get-address-affected-transactions)
- [Uxout apis](#uxout-apis)
    - [Get uxout](#get-uxout)
    - [Get address affected uxouts](#get-address-affected-uxouts)
- [Coin supply informations](#coin-supply-informations)
- [Richlist show top N addresses by uxouts](#richlist-show-top-n-addresses-by-uxouts)
- [AddressCount show count of unique address](#addresscount-show-count-of-unique-address)

<!-- /MarkdownTOC -->

## CSRF

All `POST`, `PUT` and `DELETE` requests require a CSRF token, obtained with a `GET /csrf` call.
The token must be placed in the `X-CSRF-Token` header.  A token is only valid
for 30 seconds and it is expected that the client obtains a new CSRF token
for each request.

A request rejected for invalid or expired CSRF will respond with `403 Forbidden - invalid CSRF token`
as the response body.

### Get current csrf token

```sh
URI: /csrf
Method: GET
```

example:

```sh
curl http://127.0.0.1:6420/csrf
```

result:

```json
{
    "csrf_token": "klSgXoMOFTvEnt8KptBvHjhlFnW0OIkzyFVn4i8frDvIus9iLsFukqA9sM9Rxf3pLZHRLr82vBQxTq50vbYA8g"
}
```

## Simple query apis

### Get node version info

```sh
URI: /version
Method: GET
```

example:

```sh
curl http://127.0.0.1:6420/version
```

result:

```json
{
    "version": "0.20.0",
    "commit": "cc733e9922d85c359f5f183d3a3a6e42c73ccb16"
}
```

### Get balance of addresses

```
URI: /balance
Method: GET
Args:
    addrs: comma-separated list of addresses. must contain at least one address
```

example:

```bash
curl http://127.0.0.1:6420/balance\?addrs\=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq
```

result:

```json
{
    "confirmed": {
        "coins": 70000000,
        "hours": 28052
    },
    "predicted": {
        "coins": 9000000,
        "hours": 8385
    }
}
```

### Get unspent output set of address or hash

```sh
URI: /outputs
Method: GET
Args:
    addrs  // address list, joined with ","
    hashes // hash list, joined with ","
```

Addrs and hashes cannot be combined.

example:

```sh
curl http://127.0.0.1:6420/outputs?addrs=6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY
```

or

```sh
curl http://127.0.0.1:6420/outputs?hashes=7669ff7350d2c70a88093431a7b30d3e69dda2319dcb048aa80fa0d19e12ebe0
```

result:

```sh
{
    "head_outputs": [
        {
            "hash": "7669ff7350d2c70a88093431a7b30d3e69dda2319dcb048aa80fa0d19e12ebe0",
            "block_seq": 22,
            "time": 1494275011,
            "src_tx": "b51e1933f286c4f03d73e8966186bafb25f64053db8514327291e690ae8aafa5",
            "address": "6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY",
            "coins": "2.000000",
            "hours": 633,
            "calculated_hours": 10023
        },
    ],
    "outgoing_outputs": [],
    "incoming_outputs": []
}
```

## Wallet apis

### Generate wallet seed

```
URI: /wallet/newSeed
Method: GET
Args:
    entropy: seed entropy [optional]
             can either be 128 or 256; 128 = 12 word seed, 256 = 24 word seed
             default: 128
```

example:

```bash
curl http://127.0.0.1:6420/wallet/newSeed
```

result:

```json
{
    "seed": "helmet van actor peanut differ icon trial glare member cancel marble rack"
}
```

### Create a wallet from seed

```
URI: /wallet/create
Method: POST
Args:
    seed: wallet seed [required]
    label: wallet label [required]
    scan: the number of addresses to scan ahead for balances [optional, must be > 0]
```

example:

```bash
curl http://127.0.0.1:6420/wallet/create -d "seed=$seed&label=$label&scan=5"
```

result:

```json
{
    "meta": {
        "coin": "sky",
        "filename": "2017_05_09_d554.wlt",
        "label": "",
        "lastSeed": "4795eaf6890c0ce1d67daf87d2f85523b1d19245a7a81a38c757fc4a7e3cae3e",
        "seed": "dish slide planet night tape stick ask element title sound only typical",
        "tm": "1494315855",
        "type": "deterministic",
        "version": "0.1"
    },
    "entries": [
        {
            "address": "y2JeYS4RS8L9GYM7UKdjLRyZanKHXumFoH",
            "public_key": "0343581927c12d07582168d6092d06d0a8cefdef47541f804eae33faf027932245",
            "secret_key": "6a7215780d7adf26cd697bd5186510f0ecb9e9a1c9d1e17d7f61d703e5087620"
        }
    ]
}
```

### Generate new address in wallet

```
URI: /wallet/newAddress
Method: POST
Args:
    id: wallet file name
```

example:

```bash
curl -X POST http://127.0.0.1:6420/wallet/newAddress?id=2017_05_09_d554.wlt
```

result:

```json
{
    "addresses": [
        "TDdQmMgbEVTwLe8EAiH2AoRc4SjoEFKrHB"
    ]
}
```

### Updates wallet label

```
URI: /wallet/update
Method: POST
Args:
    id: wallet file name
    label: wallet label
```

example:

```bash
curl -X POST http://127.0.0.1:6420/wallet/update?id=$id&label=$label
```

result:

```
"success"
```

### Get wallet balance

```
URI: /wallet/balance
Method: GET
Args:
    id: wallet file name
```

example:

```bash
curl http://127.0.0.1:6420/wallet/balance?id=2017_05_09_d554.wlt
```

result:

```json
{
    "confirmed": {
        "coins": 0,
        "hours": 0
    },
    "predicted": {
        "coins": 0,
        "hours": 0
    }
}
```

### Spend coins from wallet

```
URI: /wallet/spend
Method: POST
Args:
    id: wallet id
    dst: recipient address
    coins: number of coins to send, in droplets. 1 coin equals 1e6 droplets.
Response:
    balance: new balance of the wallet
    txn: spent transaction
    error: an error that may have occured after broadcast the transaction to the network
           if this field is not empty, the spend succeeded, but the response data could not be prepared
Statuses:
    200: successful spend. NOTE: the response may include an "error" field. if this occurs, the spend succeeded
         but the response data could not be prepared. The client should NOT spend again.
    400: Invalid query params, wallet lacks enough coin hours, insufficient balance
    404: wallet does not exist
    500: other errors
```

example, send 1 coin to `2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc` from wallet `2017_05_09_ea42.wlt`:

```bash
curl -X POST \
  'http://127.0.0.1:6420/wallet/spend?id=2017_05_09_ea42.wlt&dst=2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc&coins=1000000'
```

result:

```json
{
    "balance": {
        "confirmed": {
            "coins": 61000000,
            "hours": 19667
        },
        "predicted": {
            "coins": 61000000,
            "hours": 19667
        }
    },
    "txn": {
        "length": 317,
        "type": 0,
        "txid": "89578005d8730fe1789288ee7dea036160a9bd43234fb673baa6abd91289a48b",
        "inner_hash": "cac977eee019832245724aa643ceff451b9d8b24612b2f6a58177c79e8a4c26f",
        "sigs": [
            "3f084a0c750731dd985d3137200f9b5fc3de06069e62edea0cdd3a91d88e56b95aff5104a3e797ab4d6d417861af0c343efb0fff2e5ba9e7cf88ab714e10f38101",
            "e9a8aa8860d189daf0b1dbfd2a4cc309fc0c7250fa81113aa7258f9603d19727793c1b7533131605db64752aeb9c1f4465198bb1d8dd597213d6406a0a81ed3701"
        ],
        "inputs": [
            "bb89d4ed40d0e6e3a82c12e70b01a4bc240d2cd4f252cfac88235abe61bd3ad0",
            "170d6fd7be1d722a1969cb3f7d45cdf4d978129c3433915dbaf098d4f075bbfc"
        ],
        "outputs": [
            {
                "uxid": "ec9cf2f6052bab24ec57847c72cfb377c06958a9e04a077d07b6dd5bf23ec106",
                "dst": "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
                "coins": "60.000000",
                "hours": 2458
            },
            {
                "uxid": "be40210601829ba8653bac1d6ecc4049955d97fb490a48c310fd912280422bd9",
                "dst": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc",
                "coins": "1.000000",
                "hours": 2458
            }
        ]
    },
    "error": ""
}
```

## Transaction apis

### Get unconfirmed transactions

```
URI: /pendingTxs
Method: GET
```

example:

```bash
curl http://127.0.0.1:6420/pendingTxs
```

result:

```json
[
    {
        "transaction": {
            "length": 317,
            "type": 0,
            "txid": "89578005d8730fe1789288ee7dea036160a9bd43234fb673baa6abd91289a48b",
            "inner_hash": "cac977eee019832245724aa643ceff451b9d8b24612b2f6a58177c79e8a4c26f",
            "sigs": [
                "3f084a0c750731dd985d3137200f9b5fc3de06069e62edea0cdd3a91d88e56b95aff5104a3e797ab4d6d417861af0c343efb0fff2e5ba9e7cf88ab714e10f38101",
                "e9a8aa8860d189daf0b1dbfd2a4cc309fc0c7250fa81113aa7258f9603d19727793c1b7533131605db64752aeb9c1f4465198bb1d8dd597213d6406a0a81ed3701"
            ],
            "inputs": [
                "bb89d4ed40d0e6e3a82c12e70b01a4bc240d2cd4f252cfac88235abe61bd3ad0",
                "170d6fd7be1d722a1969cb3f7d45cdf4d978129c3433915dbaf098d4f075bbfc"
            ],
            "outputs": [
                {
                    "uxid": "ec9cf2f6052bab24ec57847c72cfb377c06958a9e04a077d07b6dd5bf23ec106",
                    "dst": "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
                    "coins": "60.000000",
                    "hours": 2458
                },
                {
                    "uxid": "be40210601829ba8653bac1d6ecc4049955d97fb490a48c310fd912280422bd9",
                    "dst": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc",
                    "coins": "1.000000",
                    "hours": 2458
                }
            ]
        },
        "received": "2017-05-09T10:11:57.14303834+02:00",
        "checked": "2017-05-09T10:19:58.801315452+02:00",
        "announced": "0001-01-01T00:00:00Z",
        "is_valid": true
    }
]
```

### Get transaction info by id

```
URI: /transaction
Method: GET
Args:
    txid: transaction id
```

example:

```bash
curl http://127.0.0.1:6420/transaction?txid=a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3
```

result:

```json
{
    "status": {
        "confirmed": true,
        "unconfirmed": false,
        "height": 1,
        "block_seq": 1178,
        "unknown": false
    },
    "txn": {
        "length": 183,
        "type": 0,
        "txid": "a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3",
        "inner_hash": "075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e",
        "timestamp": 1494275231,
        "sigs": [
            "9b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00"
        ],
        "inputs": [
            "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191"
        ],
        "outputs": [
            {
                "uxid": "70fa9dfb887f9ef55beb4e960f60e4703c56f98201acecf2cad729f5d7e84690",
                "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                "coins": "8.000000",
                "hours": 931
            }
        ]
    }
}
```

### Get raw transaction by id

```
URI: /rawtx
Method: GET
```

example:

```bash
curl http://127.0.0.1:6420/rawtx?txid=a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3
```

result:

```bash
"
b700000000075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e010000009b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00010000005287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191010000000010722f061aa262381dce35193d43eceb112373c300127a0000000000a303000000000000"
```

### Inject raw transaction

```
URI: /injectTransaction
Method: POST
Content-Type: application/json
Body: {
        "rawtx":"raw transaction"
      }
```

example:

```bash
curl -X POST http://127.0.0.1:6420/injectTransaction -H 'content-type: application/json' -d '{
    "rawtx":"dc0000000008b507528697b11340f5a3fcccbff031c487bad59d26c2bdaea0cd8a0199a1720100000017f36c9d8bce784df96a2d6848f1b7a8f5c890986846b7c53489eb310090b91143c98fd233830055b5959f60030b3ca08d95f22f6b96ba8c20e548d62b342b5e0001000000ec9cf2f6052bab24ec57847c72cfb377c06958a9e04a077d07b6dd5bf23ec106020000000072116096fe2207d857d18565e848b403807cd825c044840300000000330100000000000000575e472f8c5295e8fa644e9bc5e06ec10351c65f40420f000000000066020000000000000"
}'
```

result:

```bash
"3615fc23cc12a5cb9190878a2151d1cf54129ff0cd90e5fc4f4e7debebad6868"
```

### Get transactions that are addresses related

```
URI: /transactions
Method: GET
Args:
	addrs: Comma seperated addresses [optional, returns all transactions if no address is provided]
    confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
```

To get address related confirmed transactions:

```bash
curl http://127.0.0.1:6420/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY&confirmed=1
```

To get address related unconfirmed transactions:
```bash
curl http://127.0.0.1:6420/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY&confirmed=0
```

To get all addresses related transactions:

```bash
curl http://127.0.0.1:6420/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY
```


result:

```sh
[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 10492,
            "block_seq": 1177,
            "unknown": false
        },
        "time": 1494275011,
        "txn": {
            "length": 317,
            "type": 0,
            "txid": "b09cd3a8baef6a449848f50a1b97943006ca92747d4e485d0647a3ea74550eca",
            "inner_hash": "2cb370051c92521a04ba5357e229d8ffa90d9d1741ea223b44dd60a1483ee0e5",
            "timestamp": 1494275011,
            "sigs": [
                "a55155ca15f73f0762f79c15917949a936658cff668647daf82a174eed95703a02622881f9cf6c7495536676f931b2d91d389a9e7b034232b3a1519c8da6fb8800",
                "cc7d7cbd6f31adabd9bde2c0deaa9277c0f3cf807a4ec97e11872817091dc3705841a6adb74acb625ee20ab6d3525350b8663566003276073d94c3bfe22fe48e01"
            ],
            "inputs": [
                "4f4b0078a9cd19b3395e54b3f42af6adc997f77f04e0ca54016c67c4f2384e3c",
                "36f4871646b6564b2f1ab72bd768a67579a1e0242bc68bcbcf1779bc75b3dddd"
            ],
            "outputs": [
                {
                    "uxid": "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191",
                    "dst": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 7454
                },
                {
                    "uxid": "a1268e9bd2033b49b44afa765d20876467254f51e5515626780467267a65c563",
                    "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                    "coins": "1.000000",
                    "hours": 7454
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 10491,
            "block_seq": 1178,
            "unknown": false
        },
        "time": 1494275231,
        "txn": {
            "length": 183,
            "type": 0,
            "txid": "a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3",
            "inner_hash": "075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e",
            "timestamp": 1494275231,
            "sigs": [
                "9b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00"
            ],
            "inputs": [
                "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191"
            ],
            "outputs": [
                {
                    "uxid": "70fa9dfb887f9ef55beb4e960f60e4703c56f98201acecf2cad729f5d7e84690",
                    "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                    "coins": "8.000000",
                    "hours": 931
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 8730,
            "block_seq": 2939,
            "unknown": false
        },
        "time": 1505205561,
        "txn": {
            "length": 474,
            "type": 0,
            "txid": "b45e571988bc07bd0b623c999655fa878fb9bdd24c8cd24fde179bf4b26ae7b7",
            "inner_hash": "393804eca6afadc05db80cfb9e1024ef5761231c70705c406301bad33161f8bf",
            "timestamp": 1505205561,
            "sigs": [
                "fb9dd021cdff51ab56891cca0fd1600877f6e0691136dbe3f8324c3f4f7ee5bc624ded4954c1d70d8cb776ce3454d8f195bbb252e48b0f2cd388f5a733697d9301",
                "0639e61ba87a61f10b0e0114008ddd4e7090d9397370de28da27b7852b231b8e66c36d10fe3424c9b23a41266fd2c50f169233009713b332d6a48ce9c128ccef01",
                "055afe17222aab66c48c8e08e03a406bf2b8719f5221ec54c8e678078033bcd56b66bbc46a866f2be5e3f9ca454e3fbc2021630d0430b72e18c24d02df03c03100",
                "8cf56fb96e11d49bea728cb35ba5953fbc640817fac01b82e62a959ef8d4c3105298f2a6ea127bb07552abd905a667b58f6c79717e9f05258079de08d91f10a500"
            ],
            "inputs": [
                "dea9266aa7b687f4391e92f04436407c51a834274a5a33bc8bcf3189732e82e3",
                "e811bdce52ddac0d952d2546fdca8d1ac4e0ad32f170d3d73b724fb37c802652",
                "e94ccdbc07cc62fb41140b4daa7969438c749837c0808acf20dde113bdf1876b",
                "534afc496a7aee2ec55c71d85abfc27f35d16c56506f663b24d8ee4815583b6e"
            ],
            "outputs": [
                {
                    "uxid": "732e129fc1630aba3f06d833ce0a7a25f05dae5df3e7a135b5f82e99222e8c28",
                    "dst": "2hAjmdPP9R3um9JhKczeVdJUVugY6SPJBDm",
                    "coins": "6.000000",
                    "hours": 204
                }
            ]
        }
    }
]
```

## Block apis

### Get blochchain progress

```sh
URI: /blockchain/progress
Method: GET
```

example:

```sh
curl http://127.0.0.1:6420/blockchain/progress
```

result:

```json
{
    "current": 2760,
    "highest": 2760,
    "peers": [
    {
        "address": "35.157.164.126:6000",
        "height": 2760
    },
    {
        "address": "63.142.253.76:6000",
        "height": 2760
    },
    ]
}
```

### Get block by hash or seq

```sh
URI: /block
Method: GET
Args:
    hash // get block by hash
    seq  // get block by sequence number
```

```sh
curl  http://127.0.0.1:6420/block?hash=6eafd13ab6823223b714246b32c984b56e0043412950faf17defdbb2cbf3fe30
```

or

```sh
curl http://127.0.0.1:6420/block?seq=2760
```

result:

```json
{
    {
        "header": {
            "seq": 2760,
            "block_hash": "6eafd13ab6823223b714246b32c984b56e0043412950faf17defdbb2cbf3fe30",
            "previous_block_hash": "eaccd527ef263573c29000dbfb3c782ee175153c63f42abb671588b7071e877f",
            "timestamp": 1504220821,
            "fee": 196130,
            "version": 0,
            "tx_body_hash": "825ae95b81ae0ce037cdf9f1cda138bac3f3ed41c51b09e0befb71848e0f3bfd"
        },
        "body": {
            "txns": [
                {
                    "length": 220,
                    "type": 0,
                    "txid": "825ae95b81ae0ce037cdf9f1cda138bac3f3ed41c51b09e0befb71848e0f3bfd",
                   "inner_hash": "312e5dd55e06be5f9a0ee43a00d447f2fea47a7f1fb9669ecb477d2768ab04fd",
                    "sigs": [
                            "f0d0eb337e3440af6e8f0c105037ec205f36c83770d26a9e3a0fb4b7ec1a2be64764f4e31cbaf6629933c971613d10d58e6acb592704a7d511f19836441f09fb00"
                    ],
                    "inputs": [
                            "e7594379c9a6bb111205cbfa6fac908cac1d136e207960eb0429f15fde09ac8c"
                    ],
                    "outputs": [
                        {
                            "uxid": "840d0ee483c1dc085e6518e1928c68979af61188b809fc74da9fca982e6a61ba",
                            "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                            "coins": "998.000000",
                            "hours": 35390
                        },
                        {
                            "uxid": "38177c437ff42f29dc8d682e2f7c278f2203b6b02f42b1a88f9eb6c2392a7f70",
                            "dst": "2YHKP9yH7baLvkum3U6HCBiJjnAUCLS5Z9U",
                            "coins": "2.000000",
                            "hours": 70780
                        }
                    ]
                }
            ]
        }
    }
}
```

### Get blocks in specific range

```sh
URI: /blocks
Method: GET
Args:
    start // start seq
    end // end seq
```

example:

```sh
curl http://127.0.0.1:6420/blocks?start=1&end=2
```

result:

```sh
{
    "blocks": [
        {
            "header": {
                "seq": 100,
                "block_hash": "725e76907998485d367a847b0fb49f08536c592247762279fcdbd9907fee5607",
                "previous_block_hash": "5c06896760ace71b02edab01700ff9ca8c32ef1d647e14c3e0d5fa751e47867e",
                "timestamp": 1429274636,
                "fee": 613712,
                "version": 0,
                "tx_body_hash": "9f20b52befed2cbaaa4a066de7119b7fdbff09a83d8e2a82628671f51f3f6551"
            },
            "body": {
                "txns": [
                    {
                        "length": 183,
                        "type": 0,
                        "txid": "9f20b52befed2cbaaa4a066de7119b7fdbff09a83d8e2a82628671f51f3f6551",
                        "inner_hash": "c2e60dbb6ad5095985d21391cbeb679fd0787c4a20471340d63f8de437d915df",
                        "sigs": [
                            "2fefd2da9d3b4af87c4157f87da0b1bf82e3d6c9f6427572bd768cf85900d15d36971ffa17eb3b486f7692584102a7a58d9fb3ef57fa24d9a4ab02eba811ef4f00"
                        ],
                        "inputs": [
                            "aee4af7e06c24bccc2f87b16d0708bfea68ac1b420f97914965f4a23ad9e11d6"
                        ],
                        "outputs": [
                            {
                                "uxid": "194cc596d2beda803d8142ddc455872082f84b09a5edd8085082b60d314c1e29",
                                "dst": "qxmeHkwgAMfwXyaQrwv9jq3qt228xMuoT5",
                                "coins": "23000.000000",
                                "hours": 87673
                            }
                        ]
                    }
                ]
            }
        },
        {
            "header": {
                "seq": 101,
                "block_hash": "8156057fc823589288f66c91edb60c11ff004465bcbe3a402b1328be7f0d6ce0",
                "previous_block_hash": "725e76907998485d367a847b0fb49f08536c592247762279fcdbd9907fee5607",
                "timestamp": 1429274666,
                "fee": 720335,
                "version": 0,
                "tx_body_hash": "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18"
            },
            "body": {
                "txns": [
                    {
                        "length": 183,
                        "type": 0,
                        "txid": "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18",
                        "inner_hash": "45da31b68748eafdb08ef8bf1ebd1c07c0f14fcb0d66759d6cf4642adc956d06",
                        "sigs": [
                            "09bce2c888ceceeb19999005cceb1efdee254cacb60edee118b51ffd740ff6503a8f9cbd60a16c7581bfd64f7529b649d0ecc8adbe913686da97fe8c6543189001"
                        ],
                        "inputs": [
                            "6002f3afc7054c0e1161bcf2b4c1d4d1009440751bc1fe806e0eae33291399f4"
                        ],
                        "outputs": [
                            {
                                "uxid": "f9bffdcbe252acb1c3a8a1e8c99829342ba1963860d5692eebaeb9bcfbcaf274",
                                "dst": "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
                                "coins": "27000.000000",
                                "hours": 102905
                            }
                        ]
                    }
                ]
            }
        }
    ]
}
```

### Get last N blocks

```sh
URI: /last_blocks
Method: GET
Args: num
```

example:

```sh
curl http://127.0.0.1:6420/last_blocks?num=2
```

result:

```sh
{
    "blocks": [
        {
            "header": {
                "seq": 2759,
                "block_hash": "eaccd527ef263573c29000dbfb3c782ee175153c63f42abb671588b7071e877f",
                "previous_block_hash": "ae92e2b3fa12786243c20b5eb94833dfa80919443d676839911571429aad1ba9",
                "timestamp": 1504211831,
                "fee": 332560,
                "version": 0,
                "tx_body_hash": "9c5f95902e57b303954ea760df96ff933b6df2b58b58097085ed5fa9fa8a1480"
            },
            "body": {
                "txns": [
                    {
                        "length": 317,
                        "type": 0,
                        "txid": "9c5f95902e57b303954ea760df96ff933b6df2b58b58097085ed5fa9fa8a1480",
                        "inner_hash": "9baaf1956aa0cca3e5e4e9d6c247228a99dc718ff507b9b6734bf584479463e5",
                        "sigs": [
                            "44e6a0c30b3f55974ff4dccb0f19929ae9f56b2615fce673e37918dd2abb946c2dc6ad3d05aa3b35df35387e90182eed3813d3fd02669449d8bda9a18a4735a201",
                            "323dfe3c89b8357511483f9faae13ecae23b6f8078a6a475301292799f520c440ebf650cc0795505fcd17ff4bd276c23156c04a39fe1ba23dac0f0e7c1907bee01"
                        ],
                        "inputs": [
                            "aa6a295c7197e4660c2e0c26d8dfab4f68d65c3acdb5f611d70f9781abd3c004",
                            "bdf3a268e177bbc6c4333c7d585ddf30b8fb123667255f90669956c3e61cda9c"
                        ],
                        "outputs": [
                            {
                                "uxid": "448c87cdebfa8ae92f009b961463f650bf23dfc696a381e81d0c64bafebe7847",
                                "dst": "B9UG4KLggfX9MNcVuMJPm11XXNDU5vkRcY",
                                "coins": "500.000000",
                                "hours": 55426
                            },
                            {
                                "uxid": "018b4132ad1f110619ff98074f36028cee082992feb824e5409f013cf61c048c",
                                "dst": "uTHMcHr3YSEwv3M2ne9B1KfoyVkRwyDYF9",
                                "coins": "350.000000",
                                "hours": 55426
                            }
                        ]
                    }
                ]
            }
        },
        {
            "header": {
                "seq": 2760,
                "block_hash": "6eafd13ab6823223b714246b32c984b56e0043412950faf17defdbb2cbf3fe30",
                "previous_block_hash": "eaccd527ef263573c29000dbfb3c782ee175153c63f42abb671588b7071e877f",
                "timestamp": 1504220821,
                "fee": 196130,
                "version": 0,
                "tx_body_hash": "825ae95b81ae0ce037cdf9f1cda138bac3f3ed41c51b09e0befb71848e0f3bfd"
            },
            "body": {
                "txns": [
                    {
                        "length": 220,
                        "type": 0,
                        "txid": "825ae95b81ae0ce037cdf9f1cda138bac3f3ed41c51b09e0befb71848e0f3bfd",
                        "inner_hash": "312e5dd55e06be5f9a0ee43a00d447f2fea47a7f1fb9669ecb477d2768ab04fd",
                        "sigs": [
                            "f0d0eb337e3440af6e8f0c105037ec205f36c83770d26a9e3a0fb4b7ec1a2be64764f4e31cbaf6629933c971613d10d58e6acb592704a7d511f19836441f09fb00"
                        ],
                        "inputs": [
                            "e7594379c9a6bb111205cbfa6fac908cac1d136e207960eb0429f15fde09ac8c"
                        ],
                        "outputs": [
                            {
                                "uxid": "840d0ee483c1dc085e6518e1928c68979af61188b809fc74da9fca982e6a61ba",
                                "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                                "coins": "998.000000",
                                "hours": 35390
                            },
                            {
                                "uxid": "38177c437ff42f29dc8d682e2f7c278f2203b6b02f42b1a88f9eb6c2392a7f70",
                                "dst": "2YHKP9yH7baLvkum3U6HCBiJjnAUCLS5Z9U",
                                "coins": "2.000000",
                                "hours": 70780
                            }
                        ]
                    }
                ]
            }
        }
    ]
}
```

## Explorer apis

### Get address affected transactions

```sh
URI: /explorer/address
Method: GET
Args: address
```

example:

```sh
curl http://127.0.0.1:6420/explorer/address
```

result:

```json
[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 783,
            "block_seq": 10819,
            "unknown": false
        },
        "length": 220,
        "type": 0,
        "txid": "86cdee14f1b9cc06710815f51e5a546a8a33c4179433e047ed50d17b3a7a734e",
        "inner_hash": "45ade9ec2b7618f782a869796f021486dda3856bf009dc6ee633d1840fd08a75",
        "timestamp": 1516000192,
        "sigs": [
            "ecd5d555dc13007a6ce39d7036e9e9ee6319c00f653372db2a0e64147739946370ddad9bf8a3cd187d481089a66381d59b0d0725fd1663ff8ab0eed202996a1701"
        ],
        "inputs": [
            {
                "uxid": "a1a715655c526fd4ca9a12208a7b1a4754998a47415ce2870bcdecb236a3fea0",
                "owner": "2Xdt4EUnJ9HZrc41L9DTDGPNrufxUbpUv4g",
                "coins": "6149.000000",
                "hours": "11286"
            }
        ],
        "outputs": [
            {
                "uxid": "f9bc2e30f263fd4a4c677d83d40f0cea5c9adc72ee696d8b9f9721fbc93473ac",
                "dst": "2Xdt4EUnJ9HZrc41L9DTDGPNrufxUbpUv4g",
                "coins": "6029.000000",
                "hours": 64965
            },
            {
                "uxid": "03077587d2ceb5f9b3c0680522e806dda4bf39d08c0f661740c5237ba0226105",
                "dst": "ANdw72kCg5HwVkn2fRgsHRu5g9Hoe3p93s",
                "coins": "120.000000",
                "hours": 64964
            }
        ]
    }
]
```

## Uxout apis

### Get uxout

```sh
URI: /uxout
Method: GET
Args: uxid
```

example:

```sh
curl http://127.0.0.1:6420/uxout?uxid=8b64d9b058e10472b9457fd2d05a1d89cbbbd78ce1d97b16587d43379271bed1
```

result:

```json
{
    "uxid": "8b64d9b058e10472b9457fd2d05a1d89cbbbd78ce1d97b16587d43379271bed1",
    "time": 1502870712,
    "src_block_seq": 2545,
    "src_tx": "ded9e671510ab300a4ea3ee126fe8e2d50b995021e2db4589c6fb4ac000fe7bb",
    "owner_address": "c9zyTYwgR4n89KyzknpmGaaDarUCPEs9mV",
    "coins": 2000000,
    "hours": 5039,
    "spent_block_seq": 2556,
    "spent_tx": "b51e1933f286c4f03d73e8966186bafb25f64053db8514327291e690ae8aafa5"
}
```

### Get address affected uxouts

```sh
URI: /address_uxouts
Method: GET
Args: address
```

example:

```sh
curl http://127.0.0.1:6420/address_uxouts?address=
```

result:

```json
[
    {
        "uxid": "7669ff7350d2c70a88093431a7b30d3e69dda2319dcb048aa80fa0d19e12ebe0",
        "time": 1502936862,
        "src_block_seq": 2556,
        "src_tx": "b51e1933f286c4f03d73e8966186bafb25f64053db8514327291e690ae8aafa5",
        "owner_address": "6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY",
        "coins": 2000000,
        "hours": 633,
        "spent_block_seq": 0,
        "spent_tx": "0000000000000000000000000000000000000000000000000000000000000000"
    }
]
```

## Coin supply informations

```
URI: /coinSupply
Method: GET
```

example:

```bash
curl http://127.0.0.1:6420/coinSupply
```

result:

```json
{
    "current_supply": "7187500.000000",
    "total_supply": "25000000.000000",
    "max_supply": "100000000.000000",
    "current_coinhour_supply": "23499025077",
    "total_coinhour_supply": "93679828577",
    "unlocked_distribution_addresses": [
        "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
        "2EYM4WFHe4Dgz6kjAdUkM6Etep7ruz2ia6h",
        "25aGyzypSA3T9K6rgPUv1ouR13efNPtWP5m",
        "ix44h3cojvN6nqGcdpy62X7Rw6Ahnr3Thk",
        "AYV8KEBEAPCg8a59cHgqHMqYHP9nVgQDyW",
        "2Nu5Jv5Wp3RYGJU1EkjWFFHnebxMx1GjfkF",
        "2THDupTBEo7UqB6dsVizkYUvkKq82Qn4gjf",
        "tWZ11Nvor9parjg4FkwxNVcby59WVTw2iL",
        "m2joQiJRZnj3jN6NsoKNxaxzUTijkdRoSR",
        "8yf8PAQqU2cDj8Yzgz3LgBEyDqjvCh2xR7",
        "sgB3n11ZPUYHToju6TWMpUZTUcKvQnoFMJ",
        "2UYPbDBnHUEc67e7qD4eXtQQ6zfU2cyvAvk",
        "wybwGC9rhm8ZssBuzpy5goXrAdE31MPdsj",
        "JbM25o7kY7hqJZt3WGYu9pHZFCpA9TCR6t",
        "2efrft5Lnwjtk7F1p9d7BnPd72zko2hQWNi",
        "Syzmb3MiMoiNVpqFdQ38hWgffHg86D2J4e",
        "2g3GUmTQooLrNHaRDhKtLU8rWLz36Beow7F",
        "D3phtGr9iv6238b3zYXq6VgwrzwvfRzWZQ",
        "gpqsFSuMCZmsjPc6Rtgy1FmLx424tH86My",
        "2EUF3GPEUmfocnUc1w6YPtqXVCy3UZA4rAq",
        "TtAaxB3qGz5zEAhhiGkBY9VPV7cekhvRYS",
        "2fM5gVpi7XaiMPm4i29zddTNkmrKe6TzhVZ",
        "ix3NDKgxfYYANKAb5kbmwBYXPrkAsha7uG",
        "2RkPshpFFrkuaP98GprLtgHFTGvPY5e6wCK",
        "Ak1qCDNudRxZVvcW6YDAdD9jpYNNStAVqm"
    ],
    "locked_distribution_addresses": [
        "2eZYSbzBKJ7QCL4kd5LSqV478rJQGb4UNkf",
        "KPfqM6S96WtRLMuSy4XLfVwymVqivdcDoM",
        "5B98bU1nsedGJBdRD5wLtq7Z8t8ZXio8u5",
        "2iZWk5tmBynWxj2PpAFyiZzEws9qSnG3a6n",
        "XUGdPaVnMh7jtzPe3zkrf9FKh5nztFnQU5",
        "hSNgHgewJme8uaHrEuKubHYtYSDckD6hpf",
        "2DeK765jLgnMweYrMp1NaYHfzxumfR1PaQN",
        "orrAssY5V2HuQAbW9K6WktFrGieq2m23pr",
        "4Ebf4PkG9QEnQTm4MVvaZvJV6Y9av3jhgb",
        "7Uf5xJ3GkiEKaLxC2WmJ1t6SeekJeBdJfu",
        "oz4ytDKbCqpgjW3LPc52pW2CaK2gxCcWmL",
        "2ex5Z7TufQ5Z8xv5mXe53fSQRfUr35SSo7Q",
        "WV2ap7ZubTxeDdmEZ1Xo7ufGMkekLWikJu",
        "ckCTV4r1pNuz6j2VBRHhaJN9HsCLY7muLV",
        "MXJx96ZJVSjktgeYZpVK8vn1H3xWP8ooq5",
        "wyQVmno9aBJZmQ99nDSLoYWwp7YDJCWsrH",
        "2cc9wKxCsFNRkoAQDAoHke3ZoyL1mSV14cj",
        "29k9g3F5AYfVaa1joE1PpZjBED6hQXes8Mm",
        "2XPLzz4ZLf1A9ykyTCjW5gEmVjnWa8CuatH",
        "iH7DqqojTgUn2JxmY9hgFp165Nk7wKfan9",
        "RJzzwUs3c9C8Y7NFYzNfFoqiUKeBhBfPki",
        "2W2cGyiCRM4nwmmiGPgMuGaPGeBzEm7VZPn",
        "ALJVNKYL7WGxFBSriiZuwZKWD4b7fbV1od",
        "tBaeg9zE2sgmw5ZQENaPPYd6jfwpVpGTzS",
        "2hdTw5Hk3rsgpZjvk8TyKcCZoRVXU5QVrUt",
        "A1QU6jKq8YgTP79M8fwZNHUZc7hConFKmy",
        "q9RkXoty3X1fuaypDDRUi78rWgJWYJMmpJ",
        "2Xvm6is5cAPA85xnSYXDuAqiRyoXiky5RaD",
        "4CW2CPJEzxhn2PS4JoSLoWGL5QQ7dL2eji",
        "24EG6uTzL7DHNzcwsygYGRR1nfu5kco7AZ1",
        "KghGnWw5fppTrqHSERXZf61yf7GkuQdCnV",
        "2WojewRA3LbpyXTP9ANy8CZqJMgmyNm3MDr",
        "2BsMfywmGV3M2CoDA112Rs7ZBkiMHfy9X11",
        "kK1Q4gPyYfVVMzQtAPRzL8qXMqJ67Y7tKs",
        "28J4mx8xfUtM92DbQ6i2Jmqw5J7dNivfroN",
        "gQvgyG1djgtftoCVrSZmsRxr7okD4LheKw",
        "3iFGBKapAWWzbiGFSr5ScbhrEPm6Esyvia",
        "NFW2akQH2vu7AqkQXxFz2P5vkXTWkSqrSm",
        "2MQJjLnWRp9eHh6MpCwpiUeshhtmri12mci",
        "2QjRQUMyL6iodtHP9zKmxCNYZ7k3jxtk49C",
        "USdfKy7B6oFNoauHWMmoCA7ND9rHqYw2Mf",
        "cA49et9WtptYHf6wA1F8qqVgH3kS5jJ9vK",
        "qaJT9TjcMi46sTKcgwRQU8o5Lw2Ea1gC4N",
        "22pyn5RyhqtTQu4obYjuWYRNNw4i54L8xVr",
        "22dkmukC6iH4FFLBmHne6modJZZQ3MC9BAT",
        "z6CJZfYLvmd41GRVE8HASjRcy5hqbpHZvE",
        "GEBWJ2KpRQDBTCCtvnaAJV2cYurgXS8pta",
        "oS8fbEm82cprmAeineBeDkaKd7QownDZQh",
        "rQpAs1LVQdphyj9ipEAuukAoj9kNpSP8cM",
        "6NSJKsPxmqipGAfFFhUKbkopjrvEESTX3j",
        "cuC68ycVXmD2EBzYFNYQ6akhKGrh3FGjSf",
        "bw4wtYU8toepomrhWP2p8UFYfHBbvEV425",
        "HvgNmDz5jD39Gwmi9VfDY1iYMhZUpZ8GKz",
        "SbApuZAYquWP3Q6iD51BcMBQjuApYEkRVf",
        "2Ugii5yxJgLzC59jV1vF8GK7UBZdvxwobeJ",
        "21N2iJ1qnQRiJWcEqNRxXwfNp8QcmiyhtPy",
        "9TC4RGs6AtFUsbcVWnSoCdoCpSfM66ALAc",
        "oQzn55UWG4iMcY9bTNb27aTnRdfiGHAwbD",
        "2GCdwsRpQhcf8SQcynFrMVDM26Bbj6sgv9M",
        "2NRFe7REtSmaM2qAgZeG45hC8EtVGV2QjeB",
        "25RGnhN7VojHUTvQBJA9nBT5y1qTQGULMzR",
        "26uCBDfF8E2PJU2Dzz2ysgKwv9m4BhodTz9",
        "Wkvima5cF7DDFdmJQqcdq8Syaq9DuAJJRD",
        "286hSoJYxvENFSHwG51ZbmKaochLJyq4ERQ",
        "FEGxF3HPoM2HCWHn82tyeh9o7vEQq5ySGE",
        "h38DxNxGhWGTq9p5tJnN5r4Fwnn85Krrb6",
        "2c1UU8J6Y3kL4cmQh21Tj8wkzidCiZxwdwd",
        "2bJ32KuGmjmwKyAtzWdLFpXNM6t83CCPLq5",
        "2fi8oLC9zfVVGnzzQtu3Y3rffS65Hiz6QHo",
        "TKD93RxFr2Am44TntLiJQus4qcEwTtvEEQ",
        "zMDywYdGEDtTSvWnCyc3qsYHWwj9ogws74",
        "25NbotTka7TwtbXUpSCQD8RMgHKspyDubXJ",
        "2ayCELBERubQWH5QxUr3cTxrYpidvUAzsSw",
        "RMTCwLiYDKEAiJu5ekHL1NQ8UKHi5ozCPg",
        "ejJjiCwp86ykmFr5iTJ8LxQXJ2wJPTYmkm"
    ]
}
```
## Richlist show top N addresses by uxouts

```
URI: /richlist
Method: GET
Args:
    n: top N addresses, [default 20, returns all if <= 0].
    include-distribution: include distribution addresses or not, default false.
```

example:

```bash
curl "http://127.0.0.1:6420/richlist?n=4&include-distribution=true"
```

result:

```json
[
    {
        "address": "zMDywYdGEDtTSvWnCyc3qsYHWwj9ogws74",
        "coins": "1000000.000000",
        "locked": true
    },
    {
        "address": "z6CJZfYLvmd41GRVE8HASjRcy5hqbpHZvE",
        "coins": "1000000.000000",
        "locked": true
    },
    {
        "address": "wyQVmno9aBJZmQ99nDSLoYWwp7YDJCWsrH",
        "coins": "1000000.000000",
        "locked": true
    },
    {
        "address": "tBaeg9zE2sgmw5ZQENaPPYd6jfwpVpGTzS",
        "coins": "1000000.000000",
        "locked": true
    }
]
```

## AddressCount show count of unique address

```
URI: /addresscount
Method: GET
```
example:

```bash
curl "http://127.0.0.1:6420/addresscount"
```

result:

```json
{
    "count": 10103
}
```
