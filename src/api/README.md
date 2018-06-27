# REST API Documentation

API default service port is `6420`.  However, if running the desktop or standalone releases from the website, the port is randomized by default.

A REST API implemented in Go is available, see [Skycoin REST API Client Godoc](https://godoc.org/github.com/skycoin/skycoin/src/api#Client).

The API has two versions, `/api/v1` and `/api/v2`.
Previously, there was no `/api/vx` prefix.
Starting in application version 0.24.0, the existing endpoints from v0.23.0
are now prefixed with `/api/v1`. To retain the old endpoints, run the application
with `-enable-unversioned-api`.

## API Version 1

`/api/v1` endpoints have no standard format. Most of them accept formdata in POST requests,
but a few accept `application/json` instead.  Most of them return JSON but one or two
return a plaintext string.

All endpoints will set an appropriate HTTP status code, using `200` for success and codes greater than or equal to `400` for error.

`/api/v1` endpoints guarantee backwards compatibility.

## API Version 2

*Note: API Version 2 is under development, and not stable. The guidelines here are subject to change.*

`/api/v2` endpoints have a standard format.

All `/api/v2` `POST` endpoints accept only `application/json` and return `application/json`.

All `/api/v2` `GET` requires accept data in the query string.
In the future we may have choose to have `GET` requests also accept `POST` with a JSON body,
to support requests with a large query body, such as when requesting data for a large number
of addresses or transactions.

`/api/v2` responses are always JSON.  If there is an error, the JSON object will
look like this:

```json
{
    "error": {
        "code": 400,
        "message": "bad arguments",
    }
}
```

Response data will be included in a `"data"` field, which will always be a JSON object (not an array).

Some endpoints may return both `"error"` and `"data"`. This will be noted in the documentation for that endpoint.

All responses will set an appropriate HTTP status code indicating an error, and it will be equal to the value of `response["error"]["code"]`.

Since `/api/v2` is still under development, there are no guarantees for backwards compatibility.
However, any changes to the API will be recorded in the [changelog](../../CHANGELOG.md).

<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3,4,5" -->

- [CSRF](#csrf)
    - [Get current csrf token](#get-current-csrf-token)
- [General system checks](#general-system-checks)
    - [Health check](#health-check)
- [Simple query APIs](#simple-query-apis)
    - [Get node version info](#get-node-version-info)
    - [Get balance of addresses](#get-balance-of-addresses)
    - [Get unspent output set of address or hash](#get-unspent-output-set-of-address-or-hash)
    - [Verify an address](#verify-an-address)
- [Wallet APIs](#wallet-apis)
    - [Get wallet](#get-wallet)
    - [Get wallet transactions](#get-wallet-transactions)
    - [Get wallets](#get-wallets)
    - [Get wallet folder name](#get-wallet-folder-name)
    - [Generate wallet seed](#generate-wallet-seed)
    - [Create a wallet from seed](#create-a-wallet-from-seed)
    - [Generate new address in wallet](#generate-new-address-in-wallet)
    - [Updates wallet label](#updates-wallet-label)
    - [Get wallet balance](#get-wallet-balance)
    - [Spend coins from wallet](#spend-coins-from-wallet)
    - [Create transaction](#create-transaction)
    - [Unload wallet](#unload-wallet)
    - [Encrypt wallet](#encrypt-wallet)
    - [Decrypt wallet](#decrypt-wallet)
    - [Get wallet seed](#get-wallet-seed)
- [Transaction APIs](#transaction-apis)
    - [Get unconfirmed transactions](#get-unconfirmed-transactions)
    - [Get transaction info by id](#get-transaction-info-by-id)
    - [Get raw transaction by id](#get-raw-transaction-by-id)
    - [Inject raw transaction](#inject-raw-transaction)
    - [Get transactions that are addresses related](#get-transactions-that-are-addresses-related)
    - [Resend unconfirmed transactions](#resend-unconfirmed-transactions)
    - [Verify encoded transaction](#verify-encoded-transaction)
- [Block APIs](#block-apis)
    - [Get blockchain metadata](#get-blockchain-metadata)
    - [Get blockchain progress](#get-blockchain-progress)
    - [Get block by hash or seq](#get-block-by-hash-or-seq)
    - [Get blocks in specific range](#get-blocks-in-specific-range)
    - [Get last N blocks](#get-last-n-blocks)
- [Explorer APIs](#explorer-apis)
    - [Get address affected transactions](#get-address-affected-transactions)
- [Uxout APIs](#uxout-apis)
    - [Get uxout](#get-uxout)
    - [Get address affected uxouts](#get-address-affected-uxouts)
- [Coin supply related information](#coin-supply-related-information)
    - [Coin supply](#coin-supply)
    - [Richlist show top N addresses by uxouts](#richlist-show-top-n-addresses-by-uxouts)
    - [Count unique addresses](#count-unique-addresses)
- [Network status](#network-status)
    - [Get information for a specific connection](#get-information-for-a-specific-connection)
    - [Get a list of all connections](#get-a-list-of-all-connections)
    - [Get a list of all default connections](#get-a-list-of-all-default-connections)
    - [Get a list of all trusted connections](#get-a-list-of-all-trusted-connections)
    - [Get a list of all connections discovered through peer exchange](#get-a-list-of-all-connections-discovered-through-peer-exchange)

<!-- /MarkdownTOC -->

## CSRF

All `POST`, `PUT` and `DELETE` requests require a CSRF token, obtained with a `GET /api/v1/csrf` call.
The token must be placed in the `X-CSRF-Token` header. A token is only valid
for 30 seconds and it is expected that the client obtains a new CSRF token
for each request. Requesting a CSRF token invalidates any previous CSRF token.

A request rejected for invalid or expired CSRF will respond with `403 Forbidden - invalid CSRF token`
as the response body.

### Get current csrf token

```
URI: /api/v1/csrf
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/csrf
```

Result:

```json
{
    "csrf_token": "klSgXoMOFTvEnt8KptBvHjhlFnW0OIkzyFVn4i8frDvIus9iLsFukqA9sM9Rxf3pLZHRLr82vBQxTq50vbYA8g"
}
```

## General system checks

### Health check

```
URI: /api/v1/health
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/health
```

Response:

```json
{
    "blockchain": {
        "head": {
            "seq": 21175,
            "block_hash": "8a3e0aac619551ae009cfb28c2b36bb1300925f74da770d1512072314f6a4c80",
            "previous_block_hash": "001eb7911b6a6ab7c75feb88726dd2bc8b87133aebc82201c4404537eb74f7ac",
            "timestamp": 1523168686,
            "fee": 2,
            "version": 0,
            "tx_body_hash": "36be8d70d1e9f70b340ea7ecf0b247c27086bad10568044c1196fe150f6cea1b"
        },
        "unspents": 14750,
        "unconfirmed": 0,
        "time_since_last_block": "12m6s"
    },
    "version": {
        "version": "0.23.0",
        "commit": "f61b4319c2f146a5ad86f7cbda26a1ba6a09998d",
        "branch": "develop"
    },
    "open_connections": 30,
    "uptime": "13.686460853s"
}
```

## Simple query APIs

### Get node version info

```
URI: /api/v1/version
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/version
```

Result:

```json
{
    "version": "0.20.0",
    "commit": "cc733e9922d85c359f5f183d3a3a6e42c73ccb16"
}
```

### Get balance of addresses

```
URI: /api/v1/balance
Method: GET
Args:
    addrs: comma-separated list of addresses. must contain at least one address
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/balance\?addrs\=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq
```

Result:

```json
{
    "confirmed": {
        "coins": 21000000,
        "hours": 142744
    },
    "predicted": {
        "coins": 21000000,
        "hours": 142744
    },
    "addresses": {
        "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD": {
            "confirmed": {
                "coins": 9000000,
                "hours": 88075
            },
            "predicted": {
                "coins": 9000000,
                "hours": 88075
            }
        },
        "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq": {
            "confirmed": {
                "coins": 12000000,
                "hours": 54669
            },
            "predicted": {
                "coins": 12000000,
                "hours": 54669
            }
        }
    }
}
```

### Get unspent output set of address or hash

```
URI: /api/v1/outputs
Method: GET
Args:
    addrs: address list, joined with ","
    hashes: hash list, joined with ","
```

Addrs and hashes cannot be combined.

Example:

```sh
curl http://127.0.0.1:6420/api/v1/outputs?addrs=6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY
```

or

```sh
curl http://127.0.0.1:6420/api/v1/outputs?hashes=7669ff7350d2c70a88093431a7b30d3e69dda2319dcb048aa80fa0d19e12ebe0
```

Result:

```json
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

### Verify an address

```
URI: /api/v2/address/verify
Method: POST
Content-Type: application/json
Args: {"address": "<address>"}
```

Parses and validates a Skycoin address. Returns the address version in the response.

Error responses:

* `400 Bad Request`: The request body is not valid JSON or the address is missing from the request body
* `422 Unprocessable Entity`: The address is invalid

Example for a valid address:

```sh
curl -X POST http://127.0.0.1:6420/api/v2/address/verify \
 -H 'Content-Type: application/json' \
 -d '{"address":"2HTnQe3ZupkG6k8S81brNC3JycGV2Em71F2"}'
```

Result:

```json
{
    "data": {
        "version": 0,
    }
}
```

Example for an invalid address:

```sh
curl -X POST http://127.0.0.1:6420/api/v2/address/verify \
 -H 'Content-Type: application/json' \
 -d '{"address":"2aTnQe3ZupkG6k8S81brNC3JycGV2Em71F2"}'
```

Result:

```json
{
    "error": {
        "message": "Invalid checksum",
        "code": 422
    }
}
```

## Wallet APIs

### Get wallet

```
URI: /api/v1/wallet
Method: GET
Args:
    id: Wallet ID [required]
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/wallet?id=2017_11_25_e5fb.wlt
```

Result:

```json
{
    "meta": {
        "coin": "skycoin",
        "filename": "2017_11_25_e5fb.wlt",
        "label": "test",
        "type": "deterministic",
        "version": "0.2",
        "crypto_type": "",
        "timestamp": 1511640884,
        "encrypted": false
    },
    "entries": [
        {
            "address": "2HTnQe3ZupkG6k8S81brNC3JycGV2Em71F2",
            "public_key": "0316ff74a8004adf9c71fa99808ee34c3505ee73c5cf82aa301d17817da3ca33b1"
        },
        {
            "address": "SMnCGfpt7zVXm8BkRSFMLeMRA6LUu3Ewne",
            "public_key": "02539528248a1a2c4f0b73233491103ca83b40249dac3ae9eee9a10b9f9debd9a3"
        }
    ]
}
```

### Get wallet transactions

```
URI: /api/v1/wallet/transactions
Method: GET
Args:
	id: Wallet ID
```

Returns all pending transaction for all addresses by selected Wallet

Example:

```sh
curl http://127.0.0.1:6420/api/v1/wallet/transactions?id=2017_11_25_e5fb.wlt
```

Result:

```json
{
    "transactions": [
        {
            "transaction": {
                "length": 317,
                "type": 0,
                "txid": "76ecbabc53ea2a3be46983058433dda6a3cf7ea0b86ba14d90b932fa97385de7",
                "inner_hash": "5d55837bb0cbda9c9323ff9aafd7c3d31d0d38638346172fbe2d9078ebaa892a",
                "sigs": [
                    "464b7724302178c1cfeacadaaf3556a3b7e5259adf51919476c3acc695747ed244b5ce2187ce7bedb6ad65c71f7f7ff3fa6805e64fe5da3aaa00ad563c7424f600",
                    "1155537b0391d4a6ee5eac07dee5798e953dca3a7c30643403dd2d326582c7d35080a16dc22644782ce1087bfc3bd06c2bf68e9a98e3989d90831646a9be2c9101"
                ],
                "inputs": [
                    "782a8662efb0e933cab7d3ae9429ab53c4208cf44d8cdc07c2fbd7204b6b5cad",
                    "2f6b61a44086588c4eaa56a5dd9f1e0be2528861a6731608fcec38891b95db91"
                ],
                "outputs": [
                    {
                        "uxid": "bd302ef776efa8548183b89f21e90649f21b90fe2d2e90ecc1b880f2d995f226",
                        "dst": "2UXZTg4ZHF6715b6tRhtaqceuQQ3G79GiZg",
                        "coins": "998.000000",
                        "hours": 247538
                    },
                    {
                        "uxid": "31058b6bfb30bfd441aec00929e75782bce47c8a75787ba519dbb268f89d2c4b",
                        "dst": "2awsJ2CR5H6QXCF2hwDjcvcAH9SgyfxCxgz",
                        "coins": "1.000000",
                        "hours": 247538
                    }
                ]
            },
            "received": "2018-03-16T18:03:57.139109904+05:30",
            "checked": "2018-03-16T18:03:57.139109904+05:30",
            "announced": "0001-01-01T00:00:00Z",
            "is_valid": true
        }
    ]
}
```

### Get wallets

```
URI: /api/v1/wallets
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/wallets
```

Result:

```json
[
    {
        "meta": {
            "coin": "skycoin",
            "filename": "2017_11_25_e5fb.wlt",
            "label": "test",
            "type": "deterministic",
            "version": "0.2",
            "crypto_type": "",
            "timestamp": 1511640884,
            "encrypted": false
        },
        "entries": [
            {
                "address": "8C5icxR9zdkYTZZTVV3cCX7QoK4EkLuK4p",
                "public_key": "0316ff74a8004adf9c71fa99808ee34c3505ee73c5cf82aa301d17817da3ca33b1"
            },
            {
                "address": "23A1EWMZopUFLCwtXMe2CU9xTCbi5Gth643",
                "public_key": "02539528248a1a2c4f0b73233491103ca83b40249dac3ae9eee9a10b9f9debd9a3"
            }
        ]
    }
]
```

### Get wallet folder name

```
URI: /api/v1/wallets/folderName
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/wallets/folderName
```

Result:

```json
{
    "address": "/Users/user/.skycoin/wallets"
}
```

### Generate wallet seed

```
URI: /api/v1/wallet/newSeed
Method: GET
Args:
    entropy: seed entropy [optional]
             can either be 128 or 256; 128 = 12 word seed, 256 = 24 word seed
             default: 128
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/wallet/newSeed
```

Result:

```json
{
    "seed": "helmet van actor peanut differ icon trial glare member cancel marble rack"
}
```

### Create a wallet from seed

```
URI: /api/v1/wallet/create
Method: POST
Args:
    seed: wallet seed [required]
    label: wallet label [required]
    scan: the number of addresses to scan ahead for balances [optional, must be > 0]
    encrypt: encrypt wallet [optional, bool value]
    password: wallet password[optional, must be provided if encrypt is true]
```

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/create \
 -H 'Content-Type: application/x-www-form-urlencoded' \
 -d 'seed=$seed' \
 -d 'label=$label' \
 -d 'scan=5' \
 -d 'password=$password'
```

Result:

```json
{
    "meta": {
        "coin": "skycoin",
        "filename": "2017_05_09_d554.wlt",
        "label": "test",
        "type": "deterministic",
        "version": "0.2",
        "crypto_type": "",
        "timestamp": 1511640884,
        "encrypted": false
    },
    "entries": [
        {
            "address": "y2JeYS4RS8L9GYM7UKdjLRyZanKHXumFoH",
            "public_key": "0316ff74a8004adf9c71fa99808ee34c3505ee73c5cf82aa301d17817da3ca33b1"
        }
    ]
}
```

### Generate new address in wallet

```
URI: /api/v1/wallet/newAddress
Method: POST
Args:
    id: wallet file name
    num: the number you want to generate
    password: wallet password
```

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/newAddress \
 -H 'Content-Type: x-www-form-urlencoded' \
 -d 'id=2017_05_09_d554.wlt' \
 -d 'num=2' \
 -d 'password=$password'
```

Result:

```json
{
    "addresses": [
        "TDdQmMgbEVTwLe8EAiH2AoRc4SjoEFKrHB"
    ]
}
```

### Updates wallet label

```
URI: /api/v1/wallet/update
Method: POST
Args:
    id: wallet file name
    label: wallet label
```

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/update \
 -H 'Content-Type: application/x-www-form-urlencoded' \
 -d 'id=$id' \
 -d 'label=$label'
```

Result:

```json
"success"
```

### Get wallet balance

```
URI: /api/v1/wallet/balance
Method: GET
Args:
    id: wallet file name
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/wallet/balance?id=2018_03_07_3088.wlt
```

Result:

```json
{
    "confirmed": {
        "coins": 210400000,
        "hours": 1873147
    },
    "predicted": {
        "coins": 210400000,
        "hours": 1873147
    },
    "addresses": {
        "AXrFisGovRhRHipsbGahs4u2hXX7pDRT5p": {
            "confirmed": {
                "coins": 1250000,
                "hours": 941185
            },
            "predicted": {
                "coins": 1250000,
                "hours": 941185
            }
        },
        "AtNorKBpCgkSRL7zES7aAQyNjqjqPp2QJU": {
            "confirmed": {
                "coins": 1150000,
                "hours": 61534
            },
            "predicted": {
                "coins": 1150000,
                "hours": 61534
            }
        },
        "VUv9ehMZWmDvwWV36BQ3eL1ujb4MQ5TGyK": {
            "confirmed": {
                "coins": 208000000,
                "hours": 870428
            },
            "predicted": {
                "coins": 208000000,
                "hours": 870428
            }
        },
        "j4mbF1fTe8jgXbrRARZSBjDpD1hMGSe1E4": {
            "confirmed": {
                "coins": 0,
                "hours": 0
            },
            "predicted": {
                "coins": 0,
                "hours": 0
            }
        },
        "uyqBPcRCWucHXs18e9VZyNEeuNsD5tFDhy": {
            "confirmed": {
                "coins": 0,
                "hours": 0
            },
            "predicted": {
                "coins": 0,
                "hours": 0
            }
        }
    }
}
```

### Spend coins from wallet

```
URI: /api/v1/wallet/spend
Method: POST
Args:
    id: wallet id
    dst: recipient address
    coins: number of coins to send, in droplets. 1 coin equals 1e6 droplets.
    password: wallet password.
Response:
    balance: new balance of the wallet
    txn: spent transaction
    error: an error that may have occured after broadcast the transaction to the network
           if this field is not empty, the spend succeeded, but the response data could not be prepared
Statuses:
    200: successful spend. NOTE: the response may include an "error" field. if this occurs, the spend succeeded
         but the response data could not be prepared. The client should NOT spend again.
    400: Invalid query params, wallet lacks enough coin hours, insufficient balance
    403: Wallet api disabled
    404: wallet does not exist
    500: other errors
```


**This endpoint is deprecated, use [POST /wallet/transaction](#create-transaction)**

Example, send 1 coin to `2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc` from wallet `2017_05_09_ea42.wlt`:

```sh
curl -X POST  http://127.0.0.1:6420/api/v1/wallet/spend \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'id=2017_05_09_ea42.wlt' \
  -d 'dst=2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc' \
  -d 'coins=1000000'
  -d 'password=$password'
```

Result:

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

### Create transaction

```
URI: /api/v1/wallet/transaction
Method: POST
Content-Type: application/json
Args: JSON body, see examples
```

Creates a transaction, returning the transaction preview and the encoded, serialized transaction.
The `encoded_transaction` can be provided to `POST /api/v1/injectTransaction` to broadcast it to the network.

The request body includes:

* An optional change address
* A wallet to spend from with the optional ability to restrict which addresses or which unspent outputs in the wallet to use
* A list of destinations with address and coins specified, as well as optionally specifying hours
* A configuration for how destination hours are distributed, either manual or automatic
* Additional options

Example request body with manual hours selection type, unencrypted wallet and all wallet addresses may spend:

```json
{
    "ignore_unconfirmed": false,
    "hours_selection": {
        "type": "manual"
    },
    "wallet": {
        "id": "foo.wlt"
    },
    "change_address": "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
    "to": [{
        "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
        "coins": "1.032",
        "hours": 7
    }, {
        "address": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
        "coins": "99.2",
        "hours": 0
    }]
}
```

Example request body with auto hours selection type, encrypted wallet, specified spending addresses:

```json
{
    "hours_selection": {
        "type": "auto",
        "mode": "share",
        "share_factor": "0.5"
    },
    "wallet": {
        "id": "foo.wlt",
        "addresses": ["2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc"],
        "password": "foobar",
    },
    "change_address": "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
    "to": [{
        "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
        "coins": "1.032"
    }, {
        "address": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
        "coins": "99.2"
    }]
}
```

Example request body with manual hours selection type, unencrypted wallet and spending specific unspent outputs:

```json
{
    "hours_selection": {
        "type": "manual"
    },
    "wallet": {
        "id": "foo.wlt",
        "unspents": ["519c069a0593e179f226e87b528f60aea72826ec7f99d51279dd8854889ed7e2", "4e4e41996297511a40e2ef0046bd6b7118a8362c1f4f09a288c5c3ea2f4dfb85"]
    },
    "change_address": "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
    "to": [{
        "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
        "coins": "1.032",
        "hours": 7
    }, {
        "address": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
        "coins": "99.2",
        "hours": 0
    }]
}
```


The `hours_selection` field has two types: `manual` or `auto`.

If `manual`, all destination hours must be specified.

If `auto`, the `mode` field must be set. The only valid value for `mode` is `"share"`.
For the `"share"` mode, `share_factor` must also be set. This must be a decimal value greater than or equal to 0 and less than or equal to 1.
In the auto share mode, the remaining hours after the fee are shared between the destination addresses as a whole,
and the change address. Amongst the destination addresses, the shared hours are distributed proportionally.

When using the `auto` `"share"` `mode`, if there are remaining coin hours as change,
but no coins are available as change from the wallet (which are needed to retain the coin hours as change),
the `share_factor` will switch to `1.0` so that extra coin hours are distributed to the outputs
instead of being burned as an additional fee.
For the `manual` mode, if there are leftover coin hours but no coins to make change with,
the leftover coin hours will be burned in addition to the required fee.

All objects in `to` must be unique; a single transaction cannot create multiple outputs with the same `address`, `coins` and `hours`.

For example, this is a valid value for `to`, if `hours_selection.type` is `"manual"`:

```json
[{
    "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
    "coins": "1.2",
    "hours": "1"
}, {
    "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
    "coins": "1.2",
    "hours": "2"
}]
```

But this is an invalid value for `to`, if `hours_selection.type` is `"manual"`:

```json
[{
    "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
    "coins": "1.2",
    "hours": "1"
}, {
    "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
    "coins": "1.2",
    "hours": "1"
}]
```

And this is a valid value for `to`, if `hours_selection.type` is `"auto"`:

```json
[{
    "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
    "coins": "1.2"
}, {
    "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
    "coins": "1.201"
}]
```

But this is an invalid value for `to`, if `hours_selection.type` is `"auto"`:

```json
[{
    "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
    "coins": "1.2"
}, {
    "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
    "coins": "1.2"
}]
```

To control which addresses to spend from, specify `wallet.addresses`.
A subset of the unspent outputs associated with these addresses will be chosen for spending,
based upon an internal selection algorithm.

To control which unspent outputs to spend from, specify `wallet.unspents`.
A subset of these unspent outputs will be chosen for spending,
based upon an internal selection algorithm.

`wallet.addresses` and `wallets.uxouts` cannot be combined.

If neither `wallet.addresses` nor `wallet.unspents` are specified,
then all outputs associated with all addresses in the wallet may be chosen from to spend with.

`change_address` is optional.
If set, it is not required to be an address in the wallet.
If not set, it will default to one of the addresses associated with the unspent outputs being spent in the transaction.

`ignore_unconfirmed` is optional and defaults to `false`.
When `false`, the API will return an error if any of the unspent outputs
associated with the wallet addresses or the wallet outputs appear as spent in
a transaction in the unconfirmed transaction pool.
When `true`, the API will ignore unspent outputs that appear as spent in
a transaction in the unconfirmed transaction pool when building the transaction,
but not return an error.

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/transaction -H 'content-type: application/json' -d '{
    "hours_selection": {
        "type": "auto",
        "mode": "share",
        "share_factor": "0.5"
    },
    "wallet": {
        "id": "foo.wlt"
    },
    "change_address": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
    "to": [{
        "address": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
        "coins": "1",
    }, {
        "address": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
        "coins": "8.99",
    }]
}'
```

Result:

```json
{
    "transaction": {
        "length": 257,
        "type": 0,
        "txid": "5f060918d2da468a784ff440fbba80674c829caca355a27ae067f465d0a5e43e",
        "inner_hash": "97dd062820314c46da0fc18c8c6c10bfab1d5da80c30adc79bbe72e90bfab11d",
        "fee": "437691",
        "sigs": [
            "6120acebfa61ba4d3970dec5665c3c952374f5d9bbf327674a0b240de62b202b319f61182e2a262b2ca5ef5a592084299504689db5448cd64c04b1f26eb01d9100"
        ],
        "inputs": [
            {
                "uxid": "7068bfd0f0f914ea3682d0e5cb3231b75cb9f0776bf9013d79b998d96c93ce2b",
                "address": "g4XmbmVyDnkswsQTSqYRsyoh1YqydDX1wp",
                "coins": "10.000000",
                "hours": "853667",
                "calculated_hours": "862290",
                "timestamp": 1524242826,
                "block": 23575,
                "txid": "ccfbb51e94cb58a619a82502bc986fb028f632df299ce189c2ff2932574a03e7"
            }
        ],
        "outputs": [
            {
                "uxid": "519c069a0593e179f226e87b528f60aea72826ec7f99d51279dd8854889ed7e2",
                "address": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
                "coins": "1.000000",
                "hours": "22253"
            },
            {
                "uxid": "4e4e41996297511a40e2ef0046bd6b7118a8362c1f4f09a288c5c3ea2f4dfb85",
                "address": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
                "coins": "8.990000",
                "hours": "200046"
            },
            {
                "uxid": "fdeb3f77408f39e50a8e3b6803ce2347aac2eba8118c494424f9fa4959bab507",
                "address": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                "coins": "0.010000",
                "hours": "222300"
            }
        ]
    },
    "encoded_transaction": "010100000097dd062820314c46da0fc18c8c6c10bfab1d5da80c30adc79bbe72e90bfab11d010000006120acebfa61ba4d3970dec5665c3c952374f5d9bbf327674a0b240de62b202b319f61182e2a262b2ca5ef5a592084299504689db5448cd64c04b1f26eb01d9100010000007068bfd0f0f914ea3682d0e5cb3231b75cb9f0776bf9013d79b998d96c93ce2b0300000000ba2a4ac4a5ce4e03a82d2240ae3661419f7081b140420f0000000000ed5600000000000000ba2a4ac4a5ce4e03a82d2240ae3661419f7081b1302d8900000000006e0d0300000000000083874350e65e84aa6e06192408951d7aaac7809e10270000000000005c64030000000000"
}
```

### Unload wallet

```
URI: /api/v1/wallet/unload
Method: POST
Args:
    id: wallet file name
```

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/unload \
 -H 'Content-Type: x-www-form-urlencoded' \
 -d 'id=2017_05_09_d554.wlt'
```

### Encrypt wallet

```
URI: /api/v1/wallet/encrypt
Method: POST
Args:
    id: wallet id
    password: wallet password
```

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/encrypt \
 -H 'Content-Type: application/x-www-form-urlencoded' \
 -d 'id=test.wlt' \
 -d 'password=$password'
```

Result:

```json
{
    "meta": {
        "coin": "skycoin",
        "filename": "test.wlt",
        "label": "test",
        "type": "deterministic",
        "version": "0.2",
        "crypto_type": "scrypt-chacha20poly1305",
        "timestamp": 1521083044,
        "encrypted": true
    },
    "entries": [
        {
            "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
            "public_key": "0316ff74a8004adf9c71fa99808ee34c3505ee73c5cf82aa301d17817da3ca33b1"
        }
    ]
}
```

### Decrypt wallet

```
URI: /api/v1/wallet/decrypt
Method: POST
Args:
    id: wallet id
    password: wallet password
```

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/decrypt \
 -H 'Content-Type: application/x-www-form-urlencoded' \
 -d 'id=test.wlt' \
 -d 'password=$password'
```

Result:

```json
{
    "meta": {
        "coin": "skycoin",
        "filename": "test.wlt",
        "label": "test",
        "type": "deterministic",
        "version": "0.2",
        "crypto_type": "",
        "timestamp": 1521083044,
        "encrypted": false
    },
    "entries": [
        {
            "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
            "public_key": "032a1218cbafc8a93233f363c19c667cf02d42fa5a8a07c0d6feca79e82d72753d"
        }
    ]
}
```

### Get wallet seed

This endpoint is supported only when `-enable-seed-api` option is enabled and the wallet is encrypted.

```
URI: /api/v1/wallet/seed
Method: POST
Args:
    id: wallet id
    password: wallet password
```

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/seed \
 -H 'Content-type: application/x-www-form-urlencoded' \
 -d 'id=test.wlt' \
 -d 'password=$password'
```

Result:

```json
{
    "seed": "your wallet seed"
}
```

## Transaction APIs

### Get unconfirmed transactions

```
URI: /api/v1/pendingTxs
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/pendingTxs
```

Result:

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
URI: /api/v1/transaction
Method: GET
Args:
    txid: transaction id
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/transaction?txid=a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3
```

Result:

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
URI: /api/v1/rawtx
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/rawtx?txid=a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3
```

Result:

```json
"b700000000075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e010000009b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00010000005287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191010000000010722f061aa262381dce35193d43eceb112373c300127a0000000000a303000000000000"
```

### Inject raw transaction

```
URI: /api/v1/injectTransaction
Method: POST
Content-Type: application/json
Body: {"rawtx": "raw transaction"}
```

Broadcasts an encoded transaction to the network.

If there are no available connections, the API responds with a 503 Service Unavailable error.

Note that in some circumstances the transaction can fail to broadcast but this endpoint will still return successfully.
This can happen if the node's network has recently become unavailable but its connections have not timed out yet.

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/injectTransaction -H 'content-type: application/json' -d '{
    "rawtx":"dc0000000008b507528697b11340f5a3fcccbff031c487bad59d26c2bdaea0cd8a0199a1720100000017f36c9d8bce784df96a2d6848f1b7a8f5c890986846b7c53489eb310090b91143c98fd233830055b5959f60030b3ca08d95f22f6b96ba8c20e548d62b342b5e0001000000ec9cf2f6052bab24ec57847c72cfb377c06958a9e04a077d07b6dd5bf23ec106020000000072116096fe2207d857d18565e848b403807cd825c044840300000000330100000000000000575e472f8c5295e8fa644e9bc5e06ec10351c65f40420f000000000066020000000000000"
}'
```

Result:

```json
"3615fc23cc12a5cb9190878a2151d1cf54129ff0cd90e5fc4f4e7debebad6868"
```

### Get transactions that are addresses related

```
URI: /api/v1/transactions
Method: GET
Args:
	addrs: Comma seperated addresses [optional, returns all transactions if no address is provided]
    confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
```

To get address related confirmed transactions:

```sh
curl http://127.0.0.1:6420/api/v1/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY&confirmed=1
```

To get address related unconfirmed transactions:
```sh
curl http://127.0.0.1:6420/api/v1/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY&confirmed=0
```

To get all addresses related transactions:

```sh
curl http://127.0.0.1:6420/api/v1/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY
```


Result:

```json
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

### Resend unconfirmed transactions

```
URI: /api/v1/resendUnconfirmedTxns
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/resendUnconfirmedTxns
```

Result:

```json
{
    "txids":[
        "b45e571988bc07bd0b623c999655fa878fb9bdd24c8cd24fde179bf4b26ae7b7",
        "a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3"
    ]
}
```

### Verify encoded transaction

```
URI: /api/v2/transaction/verify
Method: POST
Content-Type: application/json
Args: {"encoded_transaction": "<hex encoded serialized transaction>"}
```

If the transaction can be parsed, passes validation and has not been spent, returns `200 OK` with the decoded transaction data,
and the `"confirmed"` field will be `false`.

If the transaction is structurally valid, passes validation but has been spent, returns `422 Unprocessable Entity` with the decoded transaction data,
and the `"confirmed"` field will be `true`.  The `"error"` `"message"` will be `"transaction has been spent"`.

If the transaction can be parsed but does not pass validation, returns `422 Unprocessable Entity` with the decoded transaction data.
The `"error"` object will be included in the response with the reason why.
If the transaction's inputs cannot be found in the unspent pool nor in the historical archive of unspents,
the transaction `"inputs"` metadata will be absent and only `"uxid"` will be present.

If the transaction can not be parsed, returns `400 Bad Request` and the `"error"` object will be included in the response with the reason why.

Example of valid transaction that has not been spent:

```sh
curl -X POST -H 'Content-Type: application/json' http://127.0.0.1:6420/api/v2/transaction/verify \
-d '{"encoded_transaction": "dc000000004fd024d60939fede67065b36adcaaeaf70fc009e3a5bbb8358940ccc8bbb2074010000007635ce932158ec06d94138adc9c9b19113fa4c2279002e6b13dcd0b65e0359f247e8666aa64d7a55378b9cc9983e252f5877a7cb2671c3568ec36579f8df1581000100000019ad5059a7fffc0369fc24b31db7e92e12a4ee2c134fb00d336d7495dec7354d02000000003f0555073e17ea6e45283f0f1115b520d0698d03a086010000000000010000000000000000b90dc595d102c48d3281b47428670210415f585200f22b0000000000ff01000000000000"}'
```

Result:

```json
{
    "data": {
        "confirmed": false,
        "transaction": {
            "length": 220,
            "type": 0,
            "txid": "82b5fcb182e3d70c285e59332af6b02bf11d8acc0b1407d7d82b82e9eeed94c0",
            "inner_hash": "4fd024d60939fede67065b36adcaaeaf70fc009e3a5bbb8358940ccc8bbb2074",
            "fee": "1042",
            "sigs": [
                "7635ce932158ec06d94138adc9c9b19113fa4c2279002e6b13dcd0b65e0359f247e8666aa64d7a55378b9cc9983e252f5877a7cb2671c3568ec36579f8df158100"
            ],
            "inputs": [
                {
                    "uxid": "19ad5059a7fffc0369fc24b31db7e92e12a4ee2c134fb00d336d7495dec7354d",
                    "address": "2HTnQe3ZupkG6k8S81brNC3JycGV2Em71F2",
                    "coins": "2.980000",
                    "hours": "985",
                    "calculated_hours": "1554",
                    "timestamp": 1527080354,
                    "block": 30074,
                    "txid": "94204347ef52d90b3c5d6c31a3fced56ae3f74fd8f1f5576931aeb60847f0e59"
                }
            ],
            "outputs": [
                {
                    "uxid": "b0911a5fc4dfe4524cdb82f6db9c705f4849af42fcd487a3c4abb2d17573d234",
                    "address": "SMnCGfpt7zVXm8BkRSFMLeMRA6LUu3Ewne",
                    "coins": "0.100000",
                    "hours": "1"
                },
                {
                    "uxid": "a492e6b85a434866be40da7e287bfcf14efce9803ff2fcd9d865c4046e81712a",
                    "address": "2HTnQe3ZupkG6k8S81brNC3JycGV2Em71F2",
                    "coins": "2.880000",
                    "hours": "511"
                }
            ]
        }
    }
}
```

Example of valid transaction that *has* been spent:

```sh
curl -X POST -H 'Content-Type: application/json' http://127.0.0.1:6420/api/v2/transaction/verify \
-d '{"encoded_transaction": "dc000000004fd024d60939fede67065b36adcaaeaf70fc009e3a5bbb8358940ccc8bbb2074010000007635ce932158ec06d94138adc9c9b19113fa4c2279002e6b13dcd0b65e0359f247e8666aa64d7a55378b9cc9983e252f5877a7cb2671c3568ec36579f8df1581000100000019ad5059a7fffc0369fc24b31db7e92e12a4ee2c134fb00d336d7495dec7354d02000000003f0555073e17ea6e45283f0f1115b520d0698d03a086010000000000010000000000000000b90dc595d102c48d3281b47428670210415f585200f22b0000000000ff01000000000000"}'
```

Result:

```json
{
    "error": {
        "message": "transaction has been spent",
        "code": 422
    },
    "data": {
        "confirmed": true,
        "transaction": {
            "length": 220,
            "type": 0,
            "txid": "82b5fcb182e3d70c285e59332af6b02bf11d8acc0b1407d7d82b82e9eeed94c0",
            "inner_hash": "4fd024d60939fede67065b36adcaaeaf70fc009e3a5bbb8358940ccc8bbb2074",
            "fee": "1042",
            "sigs": [
                "7635ce932158ec06d94138adc9c9b19113fa4c2279002e6b13dcd0b65e0359f247e8666aa64d7a55378b9cc9983e252f5877a7cb2671c3568ec36579f8df158100"
            ],
            "inputs": [
                {
                    "uxid": "19ad5059a7fffc0369fc24b31db7e92e12a4ee2c134fb00d336d7495dec7354d",
                    "address": "2HTnQe3ZupkG6k8S81brNC3JycGV2Em71F2",
                    "coins": "2.980000",
                    "hours": "985",
                    "calculated_hours": "1554",
                    "timestamp": 1527080354,
                    "block": 30074,
                    "txid": "94204347ef52d90b3c5d6c31a3fced56ae3f74fd8f1f5576931aeb60847f0e59"
                }
            ],
            "outputs": [
                {
                    "uxid": "b0911a5fc4dfe4524cdb82f6db9c705f4849af42fcd487a3c4abb2d17573d234",
                    "address": "SMnCGfpt7zVXm8BkRSFMLeMRA6LUu3Ewne",
                    "coins": "0.100000",
                    "hours": "1"
                },
                {
                    "uxid": "a492e6b85a434866be40da7e287bfcf14efce9803ff2fcd9d865c4046e81712a",
                    "address": "2HTnQe3ZupkG6k8S81brNC3JycGV2Em71F2",
                    "coins": "2.880000",
                    "hours": "511"
                }
            ]
        }
    }
}
```


## Block APIs

### Get blockchain metadata

```
URI: /api/v1/blockchain/metadata
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/blockchain/metadata
```

Result:

```json
{
    "head": {
        "seq": 17936,
        "block_hash": "b91663fa8ff14aab529cd7bfd48bde5bd86e3c2db154d601528801ee0d064d19",
        "previous_block_hash": "b57d3b644898f95c9f7a9281e786a0ae2a567e9dc573654363ffafaa41ab4caf",
        "timestamp": 1520967639,
        "fee": 61662,
        "version": 0,
        "tx_body_hash": "f0e8440f30acf01def3acaa9a88ea91f1fbaea19c0df003726edfe5bd1c7b51d"
    },
    "unspents": 12704,
    "unconfirmed": 0
}
```

### Get blockchain progress

```
URI: /api/v1/blockchain/progress
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/blockchain/progress
```

Result:

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

```
URI: /api/v1/block
Method: GET
Args:
    hash: get block by hash
    seq: get block by sequence number
```

```sh
curl http://127.0.0.1:6420/api/v1/block?hash=6eafd13ab6823223b714246b32c984b56e0043412950faf17defdbb2cbf3fe30
```

or

```sh
curl http://127.0.0.1:6420/api/v1/block?seq=2760
```

Result:

```json
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
    },
    "size": 220
}
```

### Get blocks in specific range

```
URI: /api/v1/blocks
Method: GET
Args:
    start: start seq
    end: end seq
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/blocks?start=1&end=2
```

Result:

```json
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
            },
            "size": 183
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
            },
            "size": 183
        }
    ]
}
```

### Get last N blocks

```
URI: /api/v1/last_blocks
Method: GET
Args:
    num: number of most recent blocks to return
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/last_blocks?num=2
```

Result:

```json
{
    "blocks": [
        {
            "header": {
                "seq": 21182,
                "block_hash": "a9045e524ff3bef82955198f274a5538ccec3958b3045c396a4b2a591fa1d99c",
                "previous_block_hash": "819b3f83afef7be9c37ab7819e193ad3a55439fb3cda52cb8691aab62bfc3936",
                "timestamp": 1523174576,
                "fee": 34572,
                "version": 0,
                "tx_body_hash": "99548cd7cc0091ce2d324647c20d22616e44424961f098c7fc81be9e9dc90c62"
            },
            "body": {
                "txns": [
                    {
                        "length": 220,
                        "type": 0,
                        "txid": "99548cd7cc0091ce2d324647c20d22616e44424961f098c7fc81be9e9dc90c62",
                        "inner_hash": "b283e783a3055c5b9e89449434ebd4f63c88450d48bc0eccd274c39c347a498a",
                        "sigs": [
                            "c0e8e1b6252cc9a5e0de3ded8a63f291781cc5866956b4409e1afa165b56245e5adb87dd29e10c73038389a056667ebba76545e0e931d261444da5e09cf4b7d901"
                        ],
                        "inputs": [
                            "0973f15386bbf39ad530c704bcfa3e768ec6515a3798455feaa89a7f00beb276"
                        ],
                        "outputs": [
                            {
                                "uxid": "7c061ba81dedcf9046d6f964efea99860dfa0b07a2de80dd24f658b8c12d2ffc",
                                "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                                "coins": "648.000000",
                                "hours": 2
                            },
                            {
                                "uxid": "89835a170bf7b1d5067c77b84c15ddc9b7d62da213eef34d8c4dbe40ab2b3158",
                                "dst": "2QzAjmkm9UZoYodZLDrDUjEuFCjU3uoNoEm",
                                "coins": "283.000000",
                                "hours": 2
                            }
                        ]
                    }
                ]
            },
            "size": 220
        },
        {
            "header": {
                "seq": 21183,
                "block_hash": "96a2f810e56545bf819ad76123609e16e4c82b5d14911b0b18368e4347e8b1b5",
                "previous_block_hash": "a9045e524ff3bef82955198f274a5538ccec3958b3045c396a4b2a591fa1d99c",
                "timestamp": 1523174636,
                "fee": 108118,
                "version": 0,
                "tx_body_hash": "c3052a92828873a0594a95a44051a9962f3794b0bcd4948d2d616496f1e06cd7"
            },
            "body": {
                "txns": [
                    {
                        "length": 220,
                        "type": 0,
                        "txid": "c3052a92828873a0594a95a44051a9962f3794b0bcd4948d2d616496f1e06cd7",
                        "inner_hash": "a81ae47fdf0a09933576da113ed098b1f8e62b7fdddf318fd5cc70420bad2228",
                        "sigs": [
                            "6bb47b0ac89cc3b0bf7b276c98dcf553f6f8980ced386cccca927d439f7cfb9c56cf6019a299f917ce62c4ab3277b7b624a351808513400d1f240f1b900d65b701"
                        ],
                        "inputs": [
                            "898d42774d3cd9910691630fb7c22f786d3ee99a6616bc5bc1b8b813e1499b06"
                        ],
                        "outputs": [
                            {
                                "uxid": "96f14e7e9d024aa49ad7e3e7b55d98bef8601ce9696f09581391eb07cda42ffe",
                                "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                                "coins": "921.000000",
                                "hours": 98
                            },
                            {
                                "uxid": "5cd59f80656e92961f824a9c65b30f61badd9265069f0a504f5c73a2bea54990",
                                "dst": "CP7tbttW82zNdygJ1UBFhzbhu9bbz8Rcez",
                                "coins": "10.000000",
                                "hours": 98
                            }
                        ]
                    }
                ]
            },
            "size": 220
        }
    ]
}
```

## Explorer APIs

### Get address affected transactions

```
URI: /api/v1/explorer/address
Method: GET
Args:
    address
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/explorer/address?address=2NfNKsaGJEndpSajJ6TsKJfsdDjW2gFsjXg
```

Result:

```json
[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 12639,
            "block_seq": 15493,
            "unknown": false
        },
        "length": 183,
        "type": 0,
        "txid": "6d8e2f8b436a2f38d604b3aa1196ef2176779c5e11e33fbdd09f993fe659c39f",
        "inner_hash": "8da7c64dcedeeb6aa1e0d21fb84a0028dcd68e6801f1a3cc0224fdd50682046f",
        "timestamp": 1518878675,
        "size": 183,
        "fee": 126249,
        "sigs": [
            "c60e43980497daad59b4c72a2eac053b1584f960c57a5e6ac8337118dccfcee4045da3f60d9be674867862a13fdd87af90f4b85cbf39913bde13674e0a039b7800"
        ],
        "inputs": [
            {
                "uxid": "349b06e5707f633fd2d8f048b687b40462d875d968b246831434fb5ab5dcac38",
                "owner": "WzPDgdfL1NzSbX96tscUNXUqtCRLjaBugC",
                "coins": "125.000000",
                "hours": 34596,
                "calculated_hours": 178174
            }
        ],
        "outputs": [
            {
                "uxid": "5b4a79c7de2e9099e083bbc8096619ae76ba6fbe34875c61bbe2d3bfa6b18b99",
                "dst": "2NfNKsaGJEndpSajJ6TsKJfsdDjW2gFsjXg",
                "coins": "125.000000",
                "hours": 51925
            }
        ]
    }
]
```

## Uxout APIs

### Get uxout

```
URI: /api/v1/uxout
Method: GET
Args:
    uxid
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/uxout?uxid=8b64d9b058e10472b9457fd2d05a1d89cbbbd78ce1d97b16587d43379271bed1
```

Result:

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

```
URI: /api/v1/address_uxouts
Method: GET
Args:
    address
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/address_uxouts?address=6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY
```

Result:

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

## Coin supply related information

### Coin supply

```
URI: /api/v1/coinSupply
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v1/coinSupply
```

Result:

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

### Richlist show top N addresses by uxouts

```
URI: /api/v1/richlist
Method: GET
Args:
    n: top N addresses, [default 20, returns all if <= 0].
    include-distribution: include distribution addresses or not, default false.
```

Example:

```sh
curl "http://127.0.0.1:6420/api/v1/richlist?n=4&include-distribution=true"
```

Result:

```json
{
    "richlist": [
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
}
```

### Count unique addresses

```
URI: /api/v1/addresscount
Method: GET
```

Example:

```sh
curl "http://127.0.0.1:6420/api/v1/addresscount"
```

Result:

```json
{
    "count": 10103
}
```

## Network status

### Get information for a specific connection

```
URI: /api/v1/network/connection
Method: GET
Args:
    addr: ip:port address of a known connection
```

Example:

```sh
curl 'http://127.0.0.1:6420/api/v1/network/connection?addr=176.9.84.75:6000'
```

Result:

```json
{
    "id": 109548,
    "address": "176.9.84.75:6000",
    "last_sent": 1520675817,
    "last_received": 1520675817,
    "outgoing": false,
    "introduced": true,
    "mirror": 719118746,
    "height": 181,
    "listen_port": 6000
}
```

### Get a list of all connections

```
URI: /api/v1/network/connections
Method: GET
```

Example:

```sh
curl 'http://127.0.0.1:6420/api/v1/network/connections'
```

Result:

```json
{
    "connections": [
        {
            "id": 99107,
            "address": "139.162.161.41:20002",
            "last_sent": 1520675750,
            "last_received": 1520675750,
            "outgoing": false,
            "introduced": true,
            "mirror": 1338939619,
            "height": 180,
            "listen_port": 20002
        },
        {
            "id": 109548,
            "address": "176.9.84.75:6000",
            "last_sent": 1520675751,
            "last_received": 1520675751,
            "outgoing": false,
            "introduced": true,
            "mirror": 719118746,
            "height": 182,
            "listen_port": 6000
        },
        {
            "id": 99115,
            "address": "185.120.34.60:6000",
            "last_sent": 1520675754,
            "last_received": 1520675754,
            "outgoing": false,
            "introduced": true,
            "mirror": 1931713869,
            "height": 180,
            "listen_port": 6000
        }
    ]
}
```


### Get a list of all default connections

```
URI: /api/v1/network/defaultConnections
Method: GET
```

Example:

```sh
curl 'http://127.0.0.1:6420/api/v1/network/defaultConnections'
```

Result:

```json
[
    "104.237.142.206:6000",
    "118.178.135.93:6000",
    "120.77.69.188:6000",
    "121.41.103.148:6000",
    "139.162.7.132:6000",
    "172.104.85.6:6000",
    "176.58.126.224:6000",
    "47.88.33.156:6000"
]
```

### Get a list of all trusted connections

```
URI: /api/v1/network/connections/trust
Method: GET
```

Example:

```sh
curl 'http://127.0.0.1:6420/api/v1/network/connections/trust'
```

Result:

```json
[
    "104.237.142.206:6000",
    "118.178.135.93:6000",
    "120.77.69.188:6000",
    "121.41.103.148:6000",
    "139.162.7.132:6000",
    "172.104.85.6:6000",
    "176.58.126.224:6000",
    "47.88.33.156:6000"
]
```

### Get a list of all connections discovered through peer exchange

```
URI: /api/v1/network/connections/exchange
Method: GET
```

Example:

```sh
curl 'http://127.0.0.1:6420/api/v1/network/connections/exchange'
```

Result:

```json
[
    "104.237.142.206:6000",
    "116.62.220.158:7200",
    "118.237.210.163:6000",
    "120.77.69.188:6000",
    "121.41.103.148:6000",
    "121.41.103.148:7200",
    "139.162.161.41:20000",
    "139.162.161.41:20001",
    "139.162.161.41:20002",
    "139.162.33.154:6000",
    "139.162.7.132:6000",
    "155.94.137.34:6000",
    "164.132.108.92:6000",
    "165.227.199.63:6000",
    "172.104.145.6:6000",
    "172.104.52.230:7200",
    "172.104.85.6:6000",
    "173.212.205.184:6000",
    "173.249.30.221:6000",
    "176.58.126.224:6000",
    "176.9.84.75:6000",
    "185.120.34.60:6000",
    "35.201.160.163:6000",
    "47.88.33.156:6000"
]
```
