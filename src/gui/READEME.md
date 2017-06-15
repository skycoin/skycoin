# APIs document

Wallet apis service port is `6420`.

## Generate wallet seed

```bash
URI: /wallet/newSeed
Method: GET
```

example:

```bash
curl http://127.0.0.1:6420/wallet/newSeed
```

result:

```bash
{
seed: "helmet van actor peanut differ icon trial glare member cancel marble rack"
}
```

## Create wallet

```bash
URI: /wallet/create
Method: POST
Arguments:
    seed [optional]
```

example:

```bash
curl -X POST http://127.0.0.1:6420/wallet/create
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

## Generate new address in wallet

```bash
URI: /wallet/newAddress
Method: POST
Arguments:
    id: wallet file name
```

example:

```bash
curl -X POST http://127.0.0.1:6420/wallet/newAddress?id=2017_05_09_d554.wlt
```

result:

```json
{
    "address": "TDdQmMgbEVTwLe8EAiH2AoRc4SjoEFKrHB"
}
```

## Get wallet balance

```bash
URI: /wallet/balance
Method: GET
Arguments:
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

## Spend coins from wallet

```bash
URI: /wallet/spend
Arguments:
      id: wallet id
     dst: recipient address
   coins: send coin number, unit is drops, 1 shellcoin = 1e6 drops
```

example:

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
                "coins": "60",
                "hours": 2458
            },
            {
                "uxid": "be40210601829ba8653bac1d6ecc4049955d97fb490a48c310fd912280422bd9",
                "dst": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc",
                "coins": "1",
                "hours": 2458
            }
        ]
    },
    "error": ""
}
```

## Get balance of addresses

```bash
URI: /balance
Method: GET
Arguments:
    addrs: addresses
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


## Get unconfirmed transactions

```bash
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
                    "coins": "60",
                    "hours": 2458
                },
                {
                    "uxid": "be40210601829ba8653bac1d6ecc4049955d97fb490a48c310fd912280422bd9",
                    "dst": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc",
                    "coins": "1",
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

## Get transaction info by id

```bash
URI: /transaction
Method: GET
Arguments:
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
                "coins": "8",
                "hours": 931
            }
        ]
    }
}
```

## Get raw transaction by id

```bash
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

## Inject raw transaction

```bash
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
