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

# Coin supply informations

```bash
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
    "current_supply": 5847530,
    "total_supply": 30000000,
    "max_supply": 100000000,
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
        "Ak1qCDNudRxZVvcW6YDAdD9jpYNNStAVqm",
        "2eZYSbzBKJ7QCL4kd5LSqV478rJQGb4UNkf",
        "KPfqM6S96WtRLMuSy4XLfVwymVqivdcDoM",
        "5B98bU1nsedGJBdRD5wLtq7Z8t8ZXio8u5",
        "2iZWk5tmBynWxj2PpAFyiZzEws9qSnG3a6n",
        "XUGdPaVnMh7jtzPe3zkrf9FKh5nztFnQU5"
    ],
    "locked_distribution_addresses": [
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
