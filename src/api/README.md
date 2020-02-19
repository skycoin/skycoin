# REST API Documentation

API default service port is `6420`. However, if running the desktop or standalone releases from the website, the port is randomized by default.

A REST API implemented in Go is available,
see [Skycoin REST API Client Godoc](https://godoc.org/github.com/SkycoinProject/skycoin/src/api#Client).

The API has two versions, `/api/v1` and `/api/v2`.

<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3,4,5" -->

- [API Version 1](#api-version-1)
- [API Version 2](#api-version-2)
- [API Sets](#api-sets)
- [Authentication](#authentication)
- [CSRF](#csrf)
	- [Get current csrf token](#get-current-csrf-token)
- [General system checks](#general-system-checks)
	- [Health check](#health-check)
	- [Version info](#version-info)
	- [Prometheus metrics](#prometheus-metrics)
- [Simple query APIs](#simple-query-apis)
	- [Get balance of addresses](#get-balance-of-addresses)
	- [Get unspent output set of address or hash](#get-unspent-output-set-of-address-or-hash)
	- [Verify an address](#verify-an-address)
- [Wallet APIs](#wallet-apis)
	- [Get wallet](#get-wallet)
	- [Get unconfirmed transactions of a wallet](#get-unconfirmed-transactions-of-a-wallet)
	- [Get wallets](#get-wallets)
	- [Get wallet folder name](#get-wallet-folder-name)
	- [Generate wallet seed](#generate-wallet-seed)
	- [Verify wallet Seed](#verify-wallet-seed)
	- [Create wallet](#create-wallet)
	- [Generate new address in wallet](#generate-new-address-in-wallet)
    - [Scan addresses in wallet](#scan-addresses-in-wallet)
	- [Change wallet label](#change-wallet-label)
	- [Get wallet balance](#get-wallet-balance)
	- [Create transaction](#create-transaction)
	- [Sign transaction](#sign-transaction)
	- [Unload wallet](#unload-wallet)
	- [Encrypt wallet](#encrypt-wallet)
	- [Decrypt wallet](#decrypt-wallet)
	- [Get wallet seed](#get-wallet-seed)
	- [Recover encrypted wallet by seed](#recover-encrypted-wallet-by-seed)
- [Key-value storage APIs](#key-value-storage-apis)
	- [Get all storage values](#get-all-storage-values)
	- [Add value to storage](#add-value-to-storage)
	- [Remove value from storage](#remove-value-from-storage)
- [Transaction APIs](#transaction-apis)
	- [Get unconfirmed transactions](#get-unconfirmed-transactions)
	- [Create transaction from unspent outputs or addresses](#create-transaction-from-unspent-outputs-or-addresses)
	- [Get transaction info by id](#get-transaction-info-by-id)
	- [Get raw transaction by id](#get-raw-transaction-by-id)
	- [Inject raw transaction](#inject-raw-transaction)
	- [Get transactions for addresses](#get-transactions-for-addresses)
    - [Get transactions with pagination](#get-transactions-with-pagination)
	- [Resend unconfirmed transactions](#resend-unconfirmed-transactions)
	- [Verify encoded transaction](#verify-encoded-transaction)
- [Block APIs](#block-apis)
	- [Get blockchain metadata](#get-blockchain-metadata)
	- [Get blockchain progress](#get-blockchain-progress)
	- [Get block by hash or seq](#get-block-by-hash-or-seq)
	- [Get blocks in specific range](#get-blocks-in-specific-range)
	- [Get last N blocks](#get-last-n-blocks)
- [Uxout APIs](#uxout-apis)
	- [Get uxout](#get-uxout)
	- [Get historical unspent outputs for an address](#get-historical-unspent-outputs-for-an-address)
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
	- [Disconnect a peer](#disconnect-a-peer)
- [Migrating from the unversioned API](#migrating-from-the-unversioned-api)
- [Migrating from the JSONRPC API](#migrating-from-the-jsonrpc-api)
- [Migrating from /api/v1/spend](#migrating-from-apiv1spend)
- [Migration from /api/v1/explorer/address](#migration-from-apiv1exploreraddress)

<!-- /MarkdownTOC -->

## API Version 1

`/api/v1` endpoints have no standard format. Most of them accept formdata in POST requests,
but a few accept `application/json` instead. Most of them return JSON but one or two
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

`/api/v2` responses are always JSON. If there is an error, the JSON object will
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

Under some circumstances an error response body may not be valid JSON.
Any client consuming the API should accomodate this and conditionally parse JSON for non-`200` responses.

## API Sets

API endpoints are grouped into "sets" which can be toggled with the command line parameters
`-enable-api-sets`, `-disable-api-sets` and `-enable-all-api-sets`.

These API sets are:

* `READ` - All query-related endpoints, they do not modify the state of the program
* `STATUS` - A subset of `READ`, these endpoints report the application, network or blockchain status
* `TXN` - Enables `/api/v1/injectTransaction` and `/api/v1/resendUnconfirmedTxns` without enabling wallet endpoints
* `WALLET` - These endpoints operate on local wallet files
* `PROMETHEUS` - This is the `/api/v2/metrics` method exposing in Prometheus text format the default metrics for Skycoin node application
* `NET_CTRL` - The `/api/v1/network/connection/disconnect` method, intended for network administration endpoints
* `INSECURE_WALLET_SEED` - This is the `/api/v1/wallet/seed` endpoint, used to decrypt and return the seed from an encrypted wallet. It is only intended for use by the desktop client.
* `STORAGE` - This is the `/api/v2/data` endpoint, used to interact with the key-value storage.

## Authentication

Authentication can be enabled with the `-web-interface-username` and `-web-interface-password` options.
The username and password should be provided in an `Authorization: Basic` header.

Authentication can only be enabled when using HTTPS with `-web-interface-https`, unless `-web-interface-plaintext-auth` is enabled.

## CSRF

All `POST`, `PUT` and `DELETE` requests require a CSRF token, obtained with a `GET /api/v1/csrf` call.
The token must be placed in the `X-CSRF-Token` header. A token is only valid
for 30 seconds and it is expected that the client obtains a new CSRF token
for each request. Requesting a CSRF token invalidates any previous CSRF token.

A request rejected for invalid or expired CSRF will respond with `403 Forbidden - invalid CSRF token`
as the response body.

### Get current csrf token

API sets: any

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

API sets: `STATUS`, `READ`

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
        "time_since_last_block": "4m46s"
    },
    "version": {
        "version": "0.25.0",
        "commit": "8798b5ee43c7ce43b9b75d57a1a6cd2c1295cd1e",
        "branch": "develop"
    },
    "coin": "skycoin",
    "user_agent": "skycoin:0.25.0",
    "open_connections": 8,
    "outgoing_connections": 5,
    "incoming_connections": 3,
    "uptime": "6m30.629057248s",
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
    "started_at": 1542443907,
    "fiber": {
        "name": "skycoin",
        "display_name": "Skycoin",
        "ticker": "SKY",
        "coin_hours_display_name": "Coin Hours",
        "coin_hours_display_name_singular": "Coin Hour",
        "coin_hours_ticker": "SCH",
        "explorer_url": "https://explorer.skycoin.com",
        "bip44_coin": 8000
    }
}
```

### Version info

API sets: any

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
    "commit": "cc733e9922d85c359f5f183d3a3a6e42c73ccb16",
    "branch": "develop"
}
```

### Prometheus metrics

API sets: `PROMETHEUS`

```
URI: /api/v2/metrics
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v2/metrics
```

Result:

```
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 5.31e-05
go_gc_duration_seconds{quantile="0.25"} 0.000158
go_gc_duration_seconds{quantile="0.5"} 0.0001789
go_gc_duration_seconds{quantile="0.75"} 0.0002216
go_gc_duration_seconds{quantile="1"} 0.0005878
go_gc_duration_seconds_sum 0.3881053
go_gc_duration_seconds_count 1959
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 30
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 2.862168e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 4.462792584e+09
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.794588e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 4.7917586e+07
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.392064e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 2.862168e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 6.0973056e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 5.087232e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 16326
# HELP go_memstats_heap_released_bytes_total Total number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes_total counter
go_memstats_heap_released_bytes_total 0
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.6060288e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.5366276699863462e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 4.7933912e+07
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 6912
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 76000
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 180224
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 5.576912e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 792284
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 1.048576e+06
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 1.048576e+06
# HELP go_memstats_sys_bytes Number of bytes obtained by system. Sum of all system allocations.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.2284408e+07
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 36.04
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 15
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 4.9025024e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.53662761869e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 8.22317056e+08
```


## Simple query APIs

### Get balance of addresses

API sets: `READ`

```
URI: /api/v1/balance
Method: GET, POST
Args:
    addrs: comma-separated list of addresses. must contain at least one address
```

Returns the cumulative and individual balances of one or more addresses.
The `POST` method can be used if many addresses need to be queried.

Example:

```sh
curl http://127.0.0.1:6420/api/v1/balance?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq,2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6
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
        "2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6": {
            "confirmed": {
                "coins": 0,
                "hours": 0
            },
            "predicted": {
                "coins": 0,
                "hours": 0
            }
        },
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

API sets: `READ`

```
URI: /api/v1/outputs
Method: GET, POST
Args:
    addrs: address list, joined with ","
    hashes: hash list, joined with ","
```

Addrs and hashes cannot be combined.

In the response, `"head_outputs"` are outputs in the current unspent output set,
`"outgoing_outputs"` are head outputs that are being spent by an unconfirmed transaction,
and `"incoming_outputs"` are outputs that will be created by an unconfirmed transaction.

The current head block header is returned as `"head"`.

The `POST` method can be used if many addresses or hashes need to be queried.

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
    "head": {
        "seq": 58891,
        "block_hash": "d9ca9442febd8788de0a3093158943beca228017bf8c9c9b8529a382fad8d991",
        "previous_block_hash": "098ea5c6e12370c38529ef7c7c38779f83d05f707affb747022eee77332ba510",
        "timestamp": 1537580414,
        "fee": 2165,
        "version": 0,
        "tx_body_hash": "c488835c85ccb153a6d42b39aaae01c3e30d16de33de282f4b3f6fa1ccf6f7eb",
        "ux_hash": "f7d30ecb49f132283862ad58f691e8747894c9fc241cb3a864fc15bd3e2c83d3"
    },
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

API sets: `READ`

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

API sets: `WALLET`

```
URI: /api/v1/wallet
Method: GET
Args:
    id: Wallet ID [required]
```

Example ("deterministic" wallet):

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

Example ("bip44" wallet):

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
        "type": "bip44",
        "version": "0.3",
        "crypto_type": "",
        "timestamp": 1511640884,
        "encrypted": false,
        "bip44_coin": 8000,
    },
    "entries": [
        {
            "address": "2HTnQe3ZupkG6k8S81brNC3JycGV2Em71F2",
            "public_key": "0316ff74a8004adf9c71fa99808ee34c3505ee73c5cf82aa301d17817da3ca33b1",
            "child_number": 0,
            "change": 0
        },
        {
            "address": "SMnCGfpt7zVXm8BkRSFMLeMRA6LUu3Ewne",
            "public_key": "02539528248a1a2c4f0b73233491103ca83b40249dac3ae9eee9a10b9f9debd9a3",
            "child_number": 1,
            "change": 0
        },
        {
            "address": "8C5icxR9zdkYTZZTVV3cCX7QoK4EkLuK4p",
            "public_key": "0316ff74a8004adf9c71fa99808ee34c3505ee73c5cf82aa301d17817da3ca33b1",
            "child_number": 0,
            "change": 1
        }
    ]
}
```


### Get unconfirmed transactions of a wallet

API sets: `WALLET`

```
URI: /api/v1/wallet/transactions
Method: GET
Args:
    id: Wallet ID
    verbose: [bool] include verbose transaction input data
```

Returns all unconfirmed transactions for all addresses in a given wallet

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
The calculated hours are based upon the current system time, and are approximately
equal to the hours the output would have if it become confirmed immediately.


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

Example (verbose):

```sh
curl http://127.0.0.1:6420/api/v1/wallet/transactions?id=2017_11_25_e5fb.wlt&verbose=1
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
                "fee": 495076,
                "inputs": [
                    {
                        "uxid": "782a8662efb0e933cab7d3ae9429ab53c4208cf44d8cdc07c2fbd7204b6b5cad",
                        "owner": "8C5icxR9zdkYTZZTVV3cCX7QoK4EkLuK4p",
                        "coins": "997.000000",
                        "hours": 880000,
                        "calculated_hours": 990000
                    },
                    {
                        "uxid": "2f6b61a44086588c4eaa56a5dd9f1e0be2528861a6731608fcec38891b95db91",
                        "owner": "23A1EWMZopUFLCwtXMe2CU9xTCbi5Gth643",
                        "coins": "2.000000",
                        "hours": 10,
                        "calculated_hours": 152
                    }
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

API sets: `WALLET`

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

API sets: `WALLET`

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

API sets: `WALLET`

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

### Verify wallet Seed

API sets: `WALLET`

```
URI: /api/v2/wallet/seed/verify
Method: POST
Args:
    seed: seed to be verified
```

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v2/wallet/seed/verify \
 -H 'Content-type: application/json' \
 -d '{ "seed": "nut wife logic sample addict shop before tobacco crisp bleak lawsuit affair" }'
```

Result:

```json
{
    "data": {}
}
```

Example (wrong bip39 seed):

```sh
curl -X POST http://127.0.0.1:6420/api/v2/wallet/seed/verify \
 -H 'Content-type: application/json' \
 -d '{ "seed": "wrong seed" }'
```

Result:

```json
{
    "error": {
        "message": "Mnemonic must have 12, 15, 18, 21 or 24 words",
        "code": 422
    }
}
```

### Create wallet

API sets: `WALLET`

```
URI: /api/v1/wallet/create
Method: POST
Args:
    seed: wallet seed [required]
    seed-passphrase: wallet seed passphrase [optional, bip44 type wallet only]
    type: wallet type [required, one of "deterministic", "bip44" or "xpub"]
    bip44-coin: BIP44 coin type [optional, defaults to 8000 (skycoin's coin type), only valid if type is "bip44"]
    xpub: xpub key [required for xpub wallets]
    label: wallet label [required]
    scan: the number of addresses to scan ahead for balances [optional, must be > 0]
    encrypt: encrypt wallet [optional, bool value]
    password: wallet password [optional, must be provided if encrypt is true]
```

Example (deterministic):

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/create \
 -H 'Content-Type: application/x-www-form-urlencoded' \
 -d 'seed=$seed' \
 -d 'type=deterministic' \
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
        "version": "0.3",
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

Example (bip44):

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/create \
 -H 'Content-Type: application/x-www-form-urlencoded' \
 -d 'seed=$seed' \
 -d 'seed-passphrase=$seed' \
 -d 'type=bip44' \
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
        "type": "bip44",
        "version": "0.3",
        "crypto_type": "scrypt-chacha20poly1305",
        "timestamp": 1511640884,
        "encrypted": true,
        "bip44_coin": 8000,
    },
    "entries": [
        {
            "address": "y2JeYS4RS8L9GYM7UKdjLRyZanKHXumFoH",
            "public_key": "0316ff74a8004adf9c71fa99808ee34c3505ee73c5cf82aa301d17817da3ca33b1",
            "child_number": 0,
            "change": 0
        }
    ]
}
```

Example (xpub):

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/create \
 -H 'Content-Type: application/x-www-form-urlencoded' \
 -d 'type=xpub' \
 -d 'xpub=xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8' \
 -d 'label=$label' \
 -d 'scan=5'
```

Result:

```json
{
    "meta": {
        "coin": "skycoin",
        "filename": "2017_05_09_d554.wlt",
        "label": "test",
        "type": "bip44",
        "version": "0.4",
        "crypto_type": "",
        "timestamp": 1511640884,
        "encrypted": false,
        "xpub": "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8"
    },
    "entries": [
        {
            "address": "y2JeYS4RS8L9GYM7UKdjLRyZanKHXumFoH",
            "public_key": "0316ff74a8004adf9c71fa99808ee34c3505ee73c5cf82aa301d17817da3ca33b1",
            "child_number": 0
        }
    ]
}
```

### Generate new address in wallet

API sets: `WALLET`

```
URI: /api/v1/wallet/newAddress
Method: POST
Args:
    id: wallet file name
    num: the number you want to generate
    password: wallet password
```

For `bip44` type wallets, the new addresses will be generated on the `external` chain (`change=0`).

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

### Scan addresses in wallet

API sets: `WALLET`

This API would scan wallet addresses ahead to search for addresses that currently have unspent outputs (coins) associated with them.

```
URI: /api/v1/wallet/scan
Method: POST
Args:
    id: wallet file name
    num: the number you want to scan ahaed [optional(default to 20), must be > 0 if provided]
    password: wallet password
```

The return value is a list of `new` generated addresses after scanning.

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/scan \
-F id=test.wlt \
-F num=10 \
-F password="${password}"
```

Result:

```json
{
    "addresses": [
        "TDdQmMgbEVTwLe8EAiH2AoRc4SjoEFKrHB"
    ]
}
```

### Change wallet label

API sets: `WALLET`

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

API sets: `WALLET`

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

### Create transaction

API sets: `WALLET`

```
URI: /api/v1/wallet/transaction
Method: POST
Content-Type: application/json
Args: JSON body, see examples
```

Creates a transaction, returning the transaction preview and the encoded, serialized transaction.
The `encoded_transaction` can be provided to `POST /api/v1/injectTransaction` to broadcast it to the network
if the transaction is fully signed.

The request body includes:

* An optional change address
* A wallet to spend from with the optional ability to restrict which addresses or which unspent outputs in the wallet to use
* A list of destinations with address and coins specified, as well as optionally specifying hours
* A configuration for how destination hours are distributed, either manual or automatic
* Additional options

`change_address` is optional. If not provided and the wallet is a `deterministic` type
wallet, then the change address will default to an address from one of the
unspent outputs being spent as a transaction input.  If the wallet is a `bip44` type
wallet, then a new, unused change address will be created.

Example request body with manual hours selection type, unencrypted wallet and all wallet addresses may spend:

```json
{
    "hours_selection": {
        "type": "manual"
    },
    "wallet_id": "foo.wlt",
    "change_address": "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
    "to": [{
        "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
        "coins": "1.032",
        "hours": "7"
    }, {
        "address": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
        "coins": "99.2",
        "hours": "0"
    }],
    "unsigned": false,
    "ignore_unconfirmed": false
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
    "wallet_id": "foo.wlt",
    "password": "foobar",
    "addresses": ["2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc"],
    "change_address": "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
    "to": [{
        "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
        "coins": "1.032"
    }, {
        "address": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
        "coins": "99.2"
    }],
    "unsigned": false,
    "ignore_unconfirmed": false
}
```

Example request body with manual hours selection type, unencrypted wallet and spending specific unspent outputs:

```json
{
    "hours_selection": {
        "type": "manual"
    },
    "wallet_id": "foo.wlt",
    "unspents": ["519c069a0593e179f226e87b528f60aea72826ec7f99d51279dd8854889ed7e2", "4e4e41996297511a40e2ef0046bd6b7118a8362c1f4f09a288c5c3ea2f4dfb85"],
    "change_address": "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
    "to": [{
        "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
        "coins": "1.032",
        "hours": "7"
    }, {
        "address": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
        "coins": "99.2",
        "hours": "0"
    }],
    "unsigned": false,
    "ignore_unconfirmed": false
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

To control which addresses to spend from, specify `addresses`.
A subset of the unspent outputs associated with these addresses will be chosen for spending,
based upon an internal selection algorithm.

To control which unspent outputs to spend from, specify `unspents`.
A subset of these unspent outputs will be chosen for spending,
based upon an internal selection algorithm.

`addresses` and `unspents` cannot be combined.

If neither `addresses` nor `unspents` are specified,
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

`unsigned` is optional and defaults to `false`.
When `true`, the transaction will not be signed by the wallet.
An unsigned transaction will be returned.
The `"txid"` value of the `"transaction"` object will need to be updated
after signing the transaction.
The unsigned `encoded_transaction` can be sent to `POST /api/v2/wallet/transaction/sign` for signing.

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/wallet/transaction -H 'content-type: application/json' -d '{
    "hours_selection": {
        "type": "auto",
        "mode": "share",
        "share_factor": "0.5"
    },
    "wallet_id": "foo.wlt",
    "change_address": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
    "to": [{
        "address": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
        "coins": "1"
    }, {
        "address": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
        "coins": "8.99"
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


### Sign transaction

API sets: `WALLET`

```
URI: /api/v2/wallet/transaction/sign
Method: POST
Content-Type: application/json
Args: JSON body, see examples
```

Signs an unsigned transaction, returning the transaction with updated signatures and the encoded, serialized transaction.
The transaction must spendable to be signed. If the inputs of the transaction are not in the unspent pool, signing fails.

Specific transaction inputs may be signed by specifying `sign_indexes`, otherwise all transaction inputs will be signed.
`sign_indexes` is an array of positional indexes for the transaction's signature array. Indexes start at 0.

Signing an input that is already signed in the transaction is an error.

The `encoded_transaction` can be provided to `POST /api/v1/injectTransaction` to broadcast it to the network, if the transaction is fully signed.

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v2/wallet/transaction/sign -H 'content-type: application/json' -d '{
    "wallet_id": "foo.wlt",
    "password": "password",
    "encoded_transaction": "010100000097dd062820314c46da0fc18c8c6c10bfab1d5da80c30adc79bbe72e90bfab11d010000006120acebfa61ba4d3970dec5665c3c952374f5d9bbf327674a0b240de62b202b319f61182e2a262b2ca5ef5a592084299504689db5448cd64c04b1f26eb01d9100010000007068bfd0f0f914ea3682d0e5cb3231b75cb9f0776bf9013d79b998d96c93ce2b0300000000ba2a4ac4a5ce4e03a82d2240ae3661419f7081b140420f0000000000ed5600000000000000ba2a4ac4a5ce4e03a82d2240ae3661419f7081b1302d8900000000006e0d0300000000000083874350e65e84aa6e06192408951d7aaac7809e10270000000000005c64030000000000"
}'
```

Example with `sign_indexes`:

```sh
curl -X POST http://127.0.0.1:6420/api/v2/wallet/transaction/sign -H 'content-type: application/json' -d '{
    "wallet_id": "foo.wlt",
    "password": "password",
    "sign_indexes": [1, 2],
    "encoded_transaction": "010100000097dd062820314c46da0fc18c8c6c10bfab1d5da80c30adc79bbe72e90bfab11d010000006120acebfa61ba4d3970dec5665c3c952374f5d9bbf327674a0b240de62b202b319f61182e2a262b2ca5ef5a592084299504689db5448cd64c04b1f26eb01d9100010000007068bfd0f0f914ea3682d0e5cb3231b75cb9f0776bf9013d79b998d96c93ce2b0300000000ba2a4ac4a5ce4e03a82d2240ae3661419f7081b140420f0000000000ed5600000000000000ba2a4ac4a5ce4e03a82d2240ae3661419f7081b1302d8900000000006e0d0300000000000083874350e65e84aa6e06192408951d7aaac7809e10270000000000005c64030000000000"
}'
```

Result:

```json
{
    "data": {
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
}
```


### Unload wallet

API sets: `WALLET`

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

API sets: `WALLET`

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

API sets: `WALLET`

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

API sets: `INSECURE_WALLET_SEED`

```
URI: /api/v1/wallet/seed
Method: POST
Args:
    id: wallet id
    password: wallet password
```

This endpoint only works for encrypted wallets.
If the wallet is unencrypted, the seed will not be returned.

If the wallet is of type `bip44` and has a seed passphrase, it will be included
in the response. Otherwise, the seed passphrase will be missing.

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
    "seed": "your wallet seed",
    "seed_passphrase": "your optional wallet seed-passphrase"
}
```

### Recover encrypted wallet by seed

API sets: `INSECURE_WALLET_SEED`

```
URI: /api/v2/wallet/recover
Method: POST
Args:
    id: wallet id
    seed: wallet seed
    seed passphrase: wallet seed passphrase (bip44 wallets only)
    password: [optional] password to encrypt the recovered wallet with
```

Recovers an encrypted wallet by providing the wallet seed and optional seed passphrase.

Example:

```sh
curl -X POST http://127.0.0.1/api/v2/wallet/recover
 -H 'Content-Type: application/json' \
 -d '{"id":"2017_11_25_e5fb.wlt","seed":"your wallet seed","seed_passphrase":"your seed passphrase"}'
```

Result:

```json
{
    "data": {
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
}
```

## Key-value storage APIs

Endpoints interact with the key-value storage. Each request require the `type` argument to
be passed.

Currently allowed types:

* `txid`: used for transaction notes
* `client`: used for generic client data, instead of using e.g. LocalStorage in the browser

### Get all storage values

API sets: `STORAGE`

```
Method: GET
URI: /api/v2/data
Args:
    type: storage type
    key [string]: key of the specific value to get
```

If key is passed, only the specific value will be returned from the storage.
Otherwise the whole dataset will be returned.

If the key does not exist, a 404 error is returned.

Example:

```sh
curl http://127.0.0.1:6420/api/v2/data?type=txid
```

Result:

```json
{
    "data": {
        "key1": "value",
        "key2": "{\"key\":\"value\"}",
    }
}
```

Example (key):

```sh
curl http://127.0.0.1:6420/api/v2/data?type=txid&key=key1
```

Result:

```json
{
    "data": "value"
}
```

### Add value to storage

API sets: `STORAGE`

```
Method: POST
URI: /api/v2/data
Args: JSON Body, see examples
```

Sets one or more values by key. Existing values will be overwritten.

Example request body:

```json
{
    "type": "txid",
    "key": "key1",
    "val": "val1"
}
```

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v2/data -H 'Content-Type: application/json' -d '{
    "type": "txid",
    "key": "key1",
    "val": "val1"
}'
```

Result:

```json
{}
```

### Remove value from storage

API sets: `STORAGE`

```
Method: DELETE
URI: /api/v2/data
Args:
    type: storage type
    key: key of the specific value to get
```

Deletes a value by key. Returns a 404 error if the key does not exist.

Example:

```sh
curl http://127.0.0.1:6420/api/v2/data?type=txid&key=key1
```

Result:

```json
{}
```

## Transaction APIs

### Get unconfirmed transactions

API sets: `READ`

```
URI: /api/v1/pendingTxs
Method: GET
Args:
    verbose [bool] include verbose transaction input data
```

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
The calculated hours are calculated based upon the current system time, and provide an approximate
coin hour value of the output if it were to be confirmed at that instant.

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

Example (verbose):

```sh
curl http://127.0.0.1:6420/api/v1/pendingTxs?verbose=1
```

Result:

```json
[
    {
        "transaction": {
            "length": 220,
            "type": 0,
            "txid": "d455564dcf1fb666c3846cf579ff33e21c203e2923938c6563fe7fcb8573ba44",
            "inner_hash": "4e73155db8ed04a3bd2b953218efcc9122ebfbf4c55f08f50d1563e48eacf71d",
            "fee": 12855964,
            "sigs": [
                "17330c256a50e2117ddccf51f1980fc14380f0f9476432196ade3043668759847b97e1b209961458745684d9239541f79d9ca9255582864d30a540017ab84f2b01"
            ],
            "inputs": [
                {
                    "uxid": "27e7bc48ceca4d47e806a87100a8a98592b7618702e1cd479bf4c190462a6d09",
                    "owner": "23MjQipM9YsPKkYiuaBmf6m7fD54wrzHxpd",
                    "coins": "7815.000000",
                    "hours": 279089,
                    "calculated_hours": 13101146
                }
            ],
            "outputs": [
                {
                    "uxid": "4b4ebf62acbaece798d0dfc92fcea85768a2874dad8a9b8eb5454288deae468c",
                    "dst": "23MjQipM9YsPKkYiuaBmf6m7fD54wrzHxpd",
                    "coins": "586.000000",
                    "hours": 122591
                },
                {
                    "uxid": "781cfb134d5fdad48f3c937dfcfc66b169a305adc8abdfe92a0ec94c564913f2",
                    "dst": "2ehrG4VKLRuvBNWYz3U7tS75QWvzyWR89Dg",
                    "coins": "7229.000000",
                    "hours": 122591
                }
            ]
        },
        "received": "2018-06-20T14:14:52.415702671+08:00",
        "checked": "2018-08-26T19:47:45.328131142+08:00",
        "announced": "2018-08-26T19:51:47.356083569+08:00",
        "is_valid": true
    }
]
```

### Create transaction from unspent outputs or addresses

API sets: `TXN`

```
URI: /api/v2/transaction
Method: POST
Args: JSON Body, see examples
```

Creates an unsigned transaction from a pool of unspent outputs or addresses.
`addresses` and `unspents` cannot be combined, and at least one must have elements in their array.

The transaction will choose unspent outputs from the provided pool to construct a transaction
that satisfies the requested outputs in the `to` field. Not all unspent outputs will necessarily be used
in the transaction.

If `ignore_unconfirmed` is true, the transaction will not use any outputs which are being spent by an unconfirmed transaction.
If `ignore_unconfirmed` is false, the endpoint returns an error if any unspent output is spent by an unconfirmed transaction.

`change_address` is optional. If not provided then the change address will
default to an address from one of the
unspent outputs being spent as a transaction input.

Refer to `POST /api/v1/wallet/transaction` for creating a transaction from a specific wallet.

`POST /api/v2/wallet/transaction/sign` can be used to sign the transaction with a wallet,
but `POST /api/v1/wallet/transaction` can create and sign a transaction with a wallet in one operation instead.
Otherwise, sign the transaction separately from the API.

The transaction must be fully valid and spendable (except for the lack of signatures) or else an error is returned.

Example request body with manual hours selection type, spending from specific addresses, ignoring unconfirmed unspent outputs:

```json
{
    "hours_selection": {
        "type": "manual"
    },
    "addresses": ["g4XmbmVyDnkswsQTSqYRsyoh1YqydDX1wp", "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS"],
    "change_address": "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
    "to": [{
        "address": "fznGedkc87a8SsW94dBowEv6J7zLGAjT17",
        "coins": "1.032",
        "hours": "7"
    }, {
        "address": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
        "coins": "99.2",
        "hours": "0"
    }],
    "ignore_unconfirmed": false
}
```

Example request body with auto hours selection type, spending specific uxouts:

```json
{
    "hours_selection": {
        "type": "auto",
        "mode": "share",
        "share_factor": "0.5"
    },
    "unspents": ["519c069a0593e179f226e87b528f60aea72826ec7f99d51279dd8854889ed7e2", "4e4e41996297511a40e2ef0046bd6b7118a8362c1f4f09a288c5c3ea2f4dfb85"],
    "change_address": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
    "to": [{
        "address": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
        "coins": "1"
    }, {
        "address": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
        "coins": "8.99"
    }]
}
```

Example:

```sh
curl -X POST http://127.0.0.1:6420/api/v2/transaction -H 'Content-Type: application/json' -d '{
    "hours_selection": {
        "type": "auto",
        "mode": "share",
        "share_factor": "0.5"
    },
    "addresses": ["g4XmbmVyDnkswsQTSqYRsyoh1YqydDX1wp"],
    "change_address": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
    "to": [{
        "address": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
        "coins": "1"
    }, {
        "address": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
        "coins": "8.99"
    }]
}'
```

Result:

```json
{
    "data": {
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
}
```

### Get transaction info by id

API sets: `READ`

```
URI: /api/v1/transaction
Method: GET
Args:
    txid: transaction id
    verbose: [bool] include verbose transaction input data
    encoded: [bool] return the transaction as hex-encoded serialized bytes
```

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
If the transaction is confirmed, the calculated hours are the hours the transaction had in the block in which it was executed..
If the transaction is unconfirmed, the calculated hours are based upon the current system time, and are approximately
equal to the hours the output would have if it become confirmed immediately.

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
        "block_seq": 1178
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

Example (verbose):

```sh
curl http://127.0.0.1:6420/api/v1/transaction?txid=a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3&verbose=1
```

Result:

```json
{
    "status": {
        "confirmed": true,
        "unconfirmed": false,
        "height": 53107,
        "block_seq": 1178
    },
    "time": 1494275231,
    "txn": {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 53107,
            "block_seq": 1178
        },
        "timestamp": 1494275231,
        "length": 183,
        "type": 0,
        "txid": "a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3",
        "inner_hash": "075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e",
        "fee": 6523,
        "sigs": [
            "9b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00"
        ],
        "inputs": [
            {
                "uxid": "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191",
                "owner": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                "coins": "8.000000",
                "hours": 7454,
                "calculated_hours": 7454
            }
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

Example (encoded):

```sh
curl http://127.0.0.1:6420/api/v1/transaction?txid=a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3&encoded=1
```

Result:

```json
{
    "status": {
        "confirmed": true,
        "unconfirmed": false,
        "height": 53267,
        "block_seq": 1178
    },
    "time": 1494275231,
    "encoded_transaction": "b700000000075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e010000009b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00010000005287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191010000000010722f061aa262381dce35193d43eceb112373c300127a0000000000a303000000000000"
}
```

### Get raw transaction by id

API sets: `READ`

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

API sets: `TXN`, `WALLET`

```
URI: /api/v1/injectTransaction
Method: POST
Content-Type: application/json
Body: {"rawtx": "hex-encoded serialized transaction string"}
Errors:
    400 - Bad input
    500 - Other
    503 - Network unavailable (transaction failed to broadcast)
```

Broadcasts a hex-encoded, serialized transaction to the network.
Transactions are serialized with the `encoder` package.
See [`coin.Transaction.Serialize`](https://godoc.org/github.com/SkycoinProject/skycoin/src/coin#Transaction.Serialize).

If there are no available connections, the API responds with a `503 Service Unavailable` error.

Note that in some circumstances the transaction can fail to broadcast but this endpoint will still return successfully.
This can happen if the node's network has recently become unavailable but its connections have not timed out yet.

The recommended way to handle transaction injections from your system is to inject the transaction then wait
for the transaction to be confirmed. Transactions typically confirm quickly, so if it is not confirmed after some
timeout such as 1 minute, the application can continue to retry the broadcast with `/api/v1/resendUnconfirmedTxns`.
Broadcast only fails without an error if the node's peers disconnect or timeout after the broadcast was initiated,
which is a network problem that may recover, so rebroadcasting with `/api/v1/resendUnconfirmedTxns` will resolve it,
or else the network is unavailable.

`POST /api/v1/transaction` accepts an `ignore_unconfirmed` option to allow transactions to be created without waiting
for unconfirmed transactions to confirm.

Any unconfirmed transactions found in the database at startup are resent. So, if the network broadcast failed but
the transaction was saved to the database, when you restart the client, it will resend.

It is safe to retry the injection after a `503` failure.

To disable the network broadcast, add `"no_broadcast": true` to the JSON request body.
The transaction will be added to the local transaction pool but not be broadcast at the same time.
Note that transactions from the pool are periodically announced, so this transaction will still
be announced eventually if the daemon continues running with connectivity for enough time.

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

Example, without broadcasting the transaction:

```sh
curl -X POST http://127.0.0.1:6420/api/v1/injectTransaction -H 'content-type: application/json' -d '{
    "rawtx":"dc0000000008b507528697b11340f5a3fcccbff031c487bad59d26c2bdaea0cd8a0199a1720100000017f36c9d8bce784df96a2d6848f1b7a8f5c890986846b7c53489eb310090b91143c98fd233830055b5959f60030b3ca08d95f22f6b96ba8c20e548d62b342b5e0001000000ec9cf2f6052bab24ec57847c72cfb377c06958a9e04a077d07b6dd5bf23ec106020000000072116096fe2207d857d18565e848b403807cd825c044840300000000330100000000000000575e472f8c5295e8fa644e9bc5e06ec10351c65f40420f000000000066020000000000000",
    "no_broadcast": true
}'
```

Result:

```json
"3615fc23cc12a5cb9190878a2151d1cf54129ff0cd90e5fc4f4e7debebad6868"
```


### Get transactions for addresses

API sets: `READ`

```
URI: /api/v1/transactions
Method: GET, POST
Args:
    addrs: Comma separated addresses [optional, returns all transactions if no address is provided]
    confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
    verbose: [bool] include verbose transaction input data
```

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
If the transaction is confirmed, the calculated hours are the hours the transaction had in the block in which it was executed.
If the transaction is unconfirmed, the calculated hours are based upon the current system time, and are approximately
equal to the hours the output would have if it become confirmed immediately.

The `"time"` field at the top level of each object in the response array indicates either the confirmed timestamp of a confirmed
transaction or the last received timestamp of an unconfirmed transaction.

The `POST` method can be used if many addresses need to be queried.

To get confirmed transactions for one or more addresses:

```sh
curl http://127.0.0.1:6420/api/v1/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY&confirmed=1
```

To get unconfirmed transactions for one or more addresses:

```sh
curl http://127.0.0.1:6420/api/v1/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY&confirmed=0
```

To get both confirmed and unconfirmed transactions for one or more addresses:

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
            "block_seq": 1177
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
            "block_seq": 1178
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
            "block_seq": 2939
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

Example (verbose):

```sh
curl http://127.0.0.1:6420/api/v1/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT&verbose=1
```

Result:

```json
[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 53207,
            "block_seq": 1131
        },
        "time": 1494192581,
        "txn": {
            "timestamp": 1494192581,
            "length": 220,
            "type": 0,
            "txid": "b785dc57a9b53dbf0390213480dd9dffc32356fb79b82fa622a2607894bfab98",
            "inner_hash": "5279e944502d6bdaff25af7b7fb7c6e503c62ae70a01084031e1cb563afe2e2c",
            "fee": 317021,
            "sigs": [
                "f8cd208acc6674de79fa1192e5177325cda871c26707242dbd6fb9df245bf34b2fbc3dfe32e61eefa0543934556cb073bdeab6e555d7bfe6b7220f1ae575613d01"
            ],
            "inputs": [
                {
                    "uxid": "004d3ef83af64c542701b923ec5c727734de9d88837bcea37a2927a569dd3f0d",
                    "owner": "MbZvwdXHnMUZ1eUFxNDqxPEEHkkffKgq2F",
                    "coins": "904.000000",
                    "hours": 14,
                    "calculated_hours": 422693
                }
            ],
            "outputs": [
                {
                    "uxid": "4047c5cbbaf0ed927caa1391d5456d58e0857ef188f2eec8ee987a30b3f53aed",
                    "dst": "MbZvwdXHnMUZ1eUFxNDqxPEEHkkffKgq2F",
                    "coins": "903.000000",
                    "hours": 52836
                },
                {
                    "uxid": "4f4b0078a9cd19b3395e54b3f42af6adc997f77f04e0ca54016c67c4f2384e3c",
                    "dst": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "1.000000",
                    "hours": 52836
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 53206,
            "block_seq": 1132
        },
        "time": 1494192731,
        "txn": {
            "timestamp": 1494192731,
            "length": 220,
            "type": 0,
            "txid": "dc39c39bea82e5b56a1a77ce8485d9b06fda694e04ddf63af1273351c87dd077",
            "inner_hash": "b8f36a57212a68f4b3ecf9d699f286dafcdb624551e07c35a983832ffd37326c",
            "fee": 39628,
            "sigs": [
                "1005adda19efe31b5cd85caa85b4a42599263f649103fd26761f2261f3ee00460d9693c45406d782b0e04613aa412a5ef6b275c2a665a9f13167912da91777a700"
            ],
            "inputs": [
                {
                    "uxid": "4047c5cbbaf0ed927caa1391d5456d58e0857ef188f2eec8ee987a30b3f53aed",
                    "owner": "MbZvwdXHnMUZ1eUFxNDqxPEEHkkffKgq2F",
                    "coins": "903.000000",
                    "hours": 52836,
                    "calculated_hours": 52836
                }
            ],
            "outputs": [
                {
                    "uxid": "a6662ea872dabee2fae96a4561d67728d16cb3da372d4b7bbc74a18f2bc3fecf",
                    "dst": "MbZvwdXHnMUZ1eUFxNDqxPEEHkkffKgq2F",
                    "coins": "895.000000",
                    "hours": 6604
                },
                {
                    "uxid": "36f4871646b6564b2f1ab72bd768a67579a1e0242bc68bcbcf1779bc75b3dddd",
                    "dst": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 6604
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 53161,
            "block_seq": 1177
        },
        "time": 1494275011,
        "txn": {
            "timestamp": 1494275011,
            "length": 317,
            "type": 0,
            "txid": "b09cd3a8baef6a449848f50a1b97943006ca92747d4e485d0647a3ea74550eca",
            "inner_hash": "2cb370051c92521a04ba5357e229d8ffa90d9d1741ea223b44dd60a1483ee0e5",
            "fee": 44726,
            "sigs": [
                "a55155ca15f73f0762f79c15917949a936658cff668647daf82a174eed95703a02622881f9cf6c7495536676f931b2d91d389a9e7b034232b3a1519c8da6fb8800",
                "cc7d7cbd6f31adabd9bde2c0deaa9277c0f3cf807a4ec97e11872817091dc3705841a6adb74acb625ee20ab6d3525350b8663566003276073d94c3bfe22fe48e01"
            ],
            "inputs": [
                {
                    "uxid": "4f4b0078a9cd19b3395e54b3f42af6adc997f77f04e0ca54016c67c4f2384e3c",
                    "owner": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "1.000000",
                    "hours": 52836,
                    "calculated_hours": 52857
                },
                {
                    "uxid": "36f4871646b6564b2f1ab72bd768a67579a1e0242bc68bcbcf1779bc75b3dddd",
                    "owner": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 6604,
                    "calculated_hours": 6777
                }
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
            "height": 53160,
            "block_seq": 1178
        },
        "time": 1494275231,
        "txn": {
            "timestamp": 1494275231,
            "length": 183,
            "type": 0,
            "txid": "a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3",
            "inner_hash": "075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e",
            "fee": 6523,
            "sigs": [
                "9b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00"
            ],
            "inputs": [
                {
                    "uxid": "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191",
                    "owner": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 7454,
                    "calculated_hours": 7454
                }
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
]
```

### Get transactions with pagination

```
URI: /api/v2/transactions
Method: GET
Args:
    addrs: Comma separated addresses [optional, returns all transactions if no address is provided]
    confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
    verbose: [bool] include verbose transaction input data
    page: Page number [optional, default to 1, must be greater than 0]
    limit: The transactions number per page [optional, default to 10, maximum to 100]
    sort: Sort the transactions by block seq [optional, default to asc, must be 'asc' or 'desc']
``` 

This API is almost the same as the `v1` version, except that it would not return all transactions by default and has
pagination supported. If there are unconfirmed transactions, they will be appended after the confirmed transactions.

If no argument is provided, the first 10 transactions will be returned. The response would have a `page_info` field which
includes `total pages`, `page size`, and `current page`.

Example:

```sh
curl http://127.0.0.1:6420/api/v2/transaction?page=1024&limit=2
```

<details>
  <summary>View Output</summary>

```json
{
    "data": {
        "page_info": {
            "total_pages": 66530,
            "page_size": 2,
            "current_page": 1024
        },
        "txns": [
            {
                "status": {
                    "confirmed": true,
                    "unconfirmed": false,
                    "height": 128216,
                    "block_seq": 2016
                },
                "time": 1500130512,
                "txn": {
                    "timestamp": 1500130512,
                    "length": 414,
                    "type": 0,
                    "txid": "f0a3c01325f3e8f09255d49b490c804b929d668fcb70ea814e1a9868b608cfdb",
                    "inner_hash": "85a298977f5fa338b7a73359c51b83787130b4f3db4a8425a1c54e45e317499d",
                    "sigs": [
                        "25333b9a283691cb189e1d2ade7dd6eeb6a275be820ff031af9b877b56330f1546a875a528bab2e559236141a644f2248a19ee5fcc86b2271f9dc60fb296f3f701",
                        "e21fdae15af052f9b842bc062ab8a2ed42baf61fe11c60255555c0fc86b99abc659269dd907472091d392d31b3c1ad24e11176ee6a9e27da1fc57e2d8ddbd04d00",
                        "5e4aa1cfca62e0a0aac1c646c3917a96bb6d1c7b8cde2e255d01730eb9d436b446cd5a09dfec097d28f5e7038a05e7172e7d5ddfe4558b1f9e3c25367051ff4f00"
                    ],
                    "inputs": [
                        "5d83e6df94ca78079c8689e700dcabdab2de959fe9f803b36fec34b47b07d025",
                        "ba1ba491090065d943ce3990b62c5d94f363bbdf37043032d79046af3687ef4c",
                        "cdce197632464ee9c46d48cb21c959772b8bf2aa04239399353988b937b6e149"
                    ],
                    "outputs": [
                        {
                            "uxid": "d19549c470bb6d217bb8095df9ef14346ee8f86730208a4247420307fadbb0f0",
                            "dst": "WSJoAtC4XcjAxTHAFLKU6MNthhpSDX7i1z",
                            "coins": "3908.000000",
                            "hours": 1070530
                        },
                        {
                            "uxid": "1742af80ec06a3ef2123a371c6f5e82c275d881e7444f8a921818bc98032fff4",
                            "dst": "2f9JhZJ147v9D4KxnJwbj8i5iNxqeKL3xNh",
                            "coins": "50.000000",
                            "hours": 1070530
                        }
                    ]
                }
            },
            {
                "status": {
                    "confirmed": true,
                    "unconfirmed": false,
                    "height": 128215,
                    "block_seq": 2017
                },
                "time": 1500130612,
                "txn": {
                    "timestamp": 1500130612,
                    "length": 414,
                    "type": 0,
                    "txid": "7dc9ae6524abe9108fdc744f210b94274a9c9fdd3da16eaea1aa88037792c27d",
                    "inner_hash": "752142fcd1cc4b9bad972611a9e64108d91b7642e1eb6b65ac92360a0c9c6bfc",
                    "sigs": [
                        "07973d43d4782ee96af70fa0ee4c73f667b035ade8770d55524b91e6d762a73b7bb6e24358c929609dd91c0140e51fa4f55952b45638bae699e522f4009c8f0c01",
                        "6d8cd2ebabcb511b1772546c898fa456bfeffefcba69aea0f6c285cda38014c015bb664f66ce0d840c71f66403782e5b6b9fd2688a212eb2e3d275aebee5856b00",
                        "4844a246eff8b59f177f9a4a43815b1fa8a3168c18d63b50c75773ee51c0b9c047ce59431871ba9fdb7df8827c9a2b175424f22cc6cf02a784b805b50574cc7000"
                    ],
                    "inputs": [
                        "8754b0d917f6690d5e88dd0950bdcb8e09d96ffb14b76964da923f5dc8969e0c",
                        "8fe7df28494563a5b47abbe737c095da461235a2529dda2a1119a19965293c8b",
                        "421ec170519fab890e3410af5ba4cf33f71fa57d786d5d39f71b7a96ed898094"
                    ],
                    "outputs": [
                        {
                            "uxid": "fceb40fd9e8895c050fc165d861ebcbe87789eeb89809879a662fdac854bf84e",
                            "dst": "2bfYafFtdkCRNcCyuDvsATV66GvBR9xfvjy",
                            "coins": "37659.000000",
                            "hours": 93972
                        },
                        {
                            "uxid": "eb8c8677da0200be7a405f2e3497db9beaa6288734d85acc5488d573ce2b8399",
                            "dst": "fXZv5X2NXhWYShoE8jazbh5UWYVCFgUXdW",
                            "coins": "25000.000000",
                            "hours": 187945
                        }
                    ]
                }
            }
        ]
    }
}
```
</details>

### Resend unconfirmed transactions

API sets: `TXN`, `WALLET`

```
URI: /api/v1/resendUnconfirmedTxns
Method: POST
```

Example:

```sh
curl -X POST 'http://127.0.0.1:6420/api/v1/resendUnconfirmedTxns'
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

API sets: `READ`

```
URI: /api/v2/transaction/verify
Method: POST
Content-Type: application/json
Args: {"unsigned": false, "encoded_transaction": "<hex encoded serialized transaction>"}
```

If the transaction can be parsed, passes validation and has not been spent, returns `200 OK` with the decoded transaction data,
and the `"confirmed"` field will be `false`.

If the transaction is structurally valid, passes validation but has been spent, returns `422 Unprocessable Entity` with the decoded transaction data,
and the `"confirmed"` field will be `true`. The `"error"` `"message"` will be `"transaction has been spent"`.

`"unsigned"` may be specified in the request. If `true`, the transaction will report an error if it is fully signed.
It will not report an error if the transaction is missing at least one signature, and the remainder of the transaction is valid.
In the response, if the transaction has any unsigned inputs, the `"unsigned"` field will be `true`.
If the request did not specify `"unsigned"` or specified it as `false`, the response will return an error for an unsigned transaction.

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
        "unsigned": false,
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
        "unsigned": false,
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

Example of valid, unsigned transaction that has not been spent, with the unsigned parameter set to true in the request:

```sh
curl -X POST -H 'Content-Type: application/json' http://127.0.0.1:6420/api/v2/transaction/verify \
-d '{"unsigned": true, "encoded_transaction": "dc000000004fd024d60939fede67065b36adcaaeaf70fc009e3a5bbb8358940ccc8bbb2074010000007635ce932158ec06d94138adc9c9b19113fa4c2279002e6b13dcd0b65e0359f247e8666aa64d7a55378b9cc9983e252f5877a7cb2671c3568ec36579f8df1581000100000019ad5059a7fffc0369fc24b31db7e92e12a4ee2c134fb00d336d7495dec7354d02000000003f0555073e17ea6e45283f0f1115b520d0698d03a086010000000000010000000000000000b90dc595d102c48d3281b47428670210415f585200f22b0000000000ff01000000000000"}'
```

Result:

```json
{
    "data": {
        "unsigned": true,
        "confirmed": false,
        "transaction": {
            "length": 220,
            "type": 0,
            "txid": "82b5fcb182e3d70c285e59332af6b02bf11d8acc0b1407d7d82b82e9eeed94c0",
            "inner_hash": "4fd024d60939fede67065b36adcaaeaf70fc009e3a5bbb8358940ccc8bbb2074",
            "fee": "1042",
            "sigs": [
                "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
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

API sets: `STATUS`, `READ`

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
    "unconfirmed": 1
}
```

### Get blockchain progress

API sets: `STATUS`, `READ`

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

API sets: `READ`

```
URI: /api/v1/block
Method: GET
Args:
    hash: get block by hash
    seq: get block by sequence number
    verbose: [bool] return verbose transaction input data
```

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
The calculated hours are the hours the transaction had in the block in which it was executed.

Example:

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
        "tx_body_hash": "825ae95b81ae0ce037cdf9f1cda138bac3f3ed41c51b09e0befb71848e0f3bfd",
        "ux_hash": "366af6bd80cfce79ce1ef63b45fb3ae8d9a6afc92a8590f14e18220884bd9d22"
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

Example (verbose):

```sh
curl http://127.0.0.1:6420/api/v1/block?hash=6eafd13ab6823223b714246b32c984b56e0043412950faf17defdbb2cbf3fe30&verbose=1
```

or

```sh
curl http://127.0.0.1:6420/api/v1/block?seq=2760&verbose=1
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
        "tx_body_hash": "825ae95b81ae0ce037cdf9f1cda138bac3f3ed41c51b09e0befb71848e0f3bfd",
        "ux_hash": "366af6bd80cfce79ce1ef63b45fb3ae8d9a6afc92a8590f14e18220884bd9d22"
    },
    "body": {
        "txns": [
            {
                "length": 220,
                "type": 0,
                "txid": "825ae95b81ae0ce037cdf9f1cda138bac3f3ed41c51b09e0befb71848e0f3bfd",
                "inner_hash": "312e5dd55e06be5f9a0ee43a00d447f2fea47a7f1fb9669ecb477d2768ab04fd",
                "fee": 196130,
                "sigs": [
                    "f0d0eb337e3440af6e8f0c105037ec205f36c83770d26a9e3a0fb4b7ec1a2be64764f4e31cbaf6629933c971613d10d58e6acb592704a7d511f19836441f09fb00"
                ],
                "inputs": [
                    {
                        "uxid": "e7594379c9a6bb111205cbfa6fac908cac1d136e207960eb0429f15fde09ac8c",
                        "owner": "kbbzyrUKNVJsJDGFLAjVT5neVcx5SQjFx5",
                        "coins": "1000.000000",
                        "hours": 283123,
                        "calculated_hours": 302300
                    }
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

API sets: `READ`

```
URI: /api/v1/blocks
Method: GET, POST
Args:
    start: start seq
    end: end seq
    seqs: comma-separated list of block seqs
    verbose: [bool] return verbose transaction input data
```

This endpoint has two modes: range and seqs.
The `seqs` parameter cannot be combined with `start`, `end`.

If `start` and/or `end` are provided, returns blocks in the range [`start`, `end`].
Both start and end sequences are included in the returned array of blocks.

If `seqs` is provided, returns blocks matching the specified sequences.
`seqs` must not contain any duplicate values.
If a block does not exist for any of the given sequence numbers, a `404` error is returned.

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
The calculated hours are the hours the transaction had in the block in which it was executed.

Example:

```sh
curl http://127.0.0.1:6420/api/v1/blocks?start=101&end=102
```

Result:

```json
{
    "blocks": [
        {
            "header": {
                "seq": 101,
                "block_hash": "8156057fc823589288f66c91edb60c11ff004465bcbe3a402b1328be7f0d6ce0",
                "previous_block_hash": "725e76907998485d367a847b0fb49f08536c592247762279fcdbd9907fee5607",
                "timestamp": 1429274666,
                "fee": 720335,
                "version": 0,
                "tx_body_hash": "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18",
                "ux_hash": "348989599d30d3adfaaea98577963caa419ab0276279296e7d194a9cbb8cad04"
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
        },
        {
            "header": {
                "seq": 102,
                "block_hash": "311f4b83b4fdb9fd1d45648115969cf4b3aab2d1acad9e2aa735829245c525f3",
                "previous_block_hash": "8156057fc823589288f66c91edb60c11ff004465bcbe3a402b1328be7f0d6ce0",
                "timestamp": 1429274686,
                "fee": 710046,
                "version": 0,
                "tx_body_hash": "7b13cab45b52dd2df291ec97cf000bf6ea1b647d6fdf0261a7527578d8b71b9d",
                "ux_hash": "f7512b0718f392c7503f86e69175efd7835ea4c3dd3f71ff65c7ad8873a6a9e8"
            },
            "body": {
                "txns": [
                    {
                        "length": 183,
                        "type": 0,
                        "txid": "7b13cab45b52dd2df291ec97cf000bf6ea1b647d6fdf0261a7527578d8b71b9d",
                        "inner_hash": "73bfee3a7c8d4f8a68657ebcaf69a59639f762bfc1a6f4468f3ca4724bc5b9f8",
                        "sigs": [
                            "c4bcada17604a4a62baf50f929655027f2913639c27b773871f2135b72553c1959737e39d50e8349ffa5a7679de845aa6370999dbaaff4c7f9fd01260818683901"
                        ],
                        "inputs": [
                            "4e75b4bced3404590d38ca06440c275d7fd86618a84966a0a1053fb18164e898"
                        ],
                        "outputs": [
                            {
                                "uxid": "0a5603a1a5aeda575aa498cdaec5a4c893a28669dba84163eba2e90db3d9f39d",
                                "dst": "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8",
                                "coins": "26700.000000",
                                "hours": 101435
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

Example (verbose):

```sh
curl http://127.0.0.1:6420/api/v1/blocks?start=101&end=102&verbose=1
```

Result:

```json
{
    "blocks": [
        {
            "header": {
                "seq": 101,
                "block_hash": "8156057fc823589288f66c91edb60c11ff004465bcbe3a402b1328be7f0d6ce0",
                "previous_block_hash": "725e76907998485d367a847b0fb49f08536c592247762279fcdbd9907fee5607",
                "timestamp": 1429274666,
                "fee": 720335,
                "version": 0,
                "tx_body_hash": "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18",
                "ux_hash": "348989599d30d3adfaaea98577963caa419ab0276279296e7d194a9cbb8cad04"
            },
            "body": {
                "txns": [
                    {
                        "length": 183,
                        "type": 0,
                        "txid": "e8fe5290afba3933389fd5860dca2cbcc81821028be9c65d0bb7cf4e8d2c4c18",
                        "inner_hash": "45da31b68748eafdb08ef8bf1ebd1c07c0f14fcb0d66759d6cf4642adc956d06",
                        "fee": 720335,
                        "sigs": [
                            "09bce2c888ceceeb19999005cceb1efdee254cacb60edee118b51ffd740ff6503a8f9cbd60a16c7581bfd64f7529b649d0ecc8adbe913686da97fe8c6543189001"
                        ],
                        "inputs": [
                            {
                                "uxid": "6002f3afc7054c0e1161bcf2b4c1d4d1009440751bc1fe806e0eae33291399f4",
                                "owner": "2M1C5LSZ4Pvu5RWS44bCdY6or3R8grQw7ez",
                                "coins": "27000.000000",
                                "hours": 220,
                                "calculated_hours": 823240
                            }
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
        },
        {
            "header": {
                "seq": 102,
                "block_hash": "311f4b83b4fdb9fd1d45648115969cf4b3aab2d1acad9e2aa735829245c525f3",
                "previous_block_hash": "8156057fc823589288f66c91edb60c11ff004465bcbe3a402b1328be7f0d6ce0",
                "timestamp": 1429274686,
                "fee": 710046,
                "version": 0,
                "tx_body_hash": "7b13cab45b52dd2df291ec97cf000bf6ea1b647d6fdf0261a7527578d8b71b9d",
                "ux_hash": "f7512b0718f392c7503f86e69175efd7835ea4c3dd3f71ff65c7ad8873a6a9e8"
            },
            "body": {
                "txns": [
                    {
                        "length": 183,
                        "type": 0,
                        "txid": "7b13cab45b52dd2df291ec97cf000bf6ea1b647d6fdf0261a7527578d8b71b9d",
                        "inner_hash": "73bfee3a7c8d4f8a68657ebcaf69a59639f762bfc1a6f4468f3ca4724bc5b9f8",
                        "fee": 710046,
                        "sigs": [
                            "c4bcada17604a4a62baf50f929655027f2913639c27b773871f2135b72553c1959737e39d50e8349ffa5a7679de845aa6370999dbaaff4c7f9fd01260818683901"
                        ],
                        "inputs": [
                            {
                                "uxid": "4e75b4bced3404590d38ca06440c275d7fd86618a84966a0a1053fb18164e898",
                                "owner": "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8",
                                "coins": "26700.000000",
                                "hours": 54,
                                "calculated_hours": 811481
                            }
                        ],
                        "outputs": [
                            {
                                "uxid": "0a5603a1a5aeda575aa498cdaec5a4c893a28669dba84163eba2e90db3d9f39d",
                                "dst": "2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8",
                                "coins": "26700.000000",
                                "hours": 101435
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

Example (seqs):

```sh
curl http://127.0.0.1:6420/api/v1/blocks?seqs=3,5,7
```

```json
{
    "blocks": [
        {
            "header": {
                "seq": 3,
                "block_hash": "35c3ebbe6feaeeab27ac77c1712051787bdd4bbfb5cdcdebc81f8aac98a2f3f3",
                "previous_block_hash": "01723bc4dc90f1cb857a94fe5e3bb50c02e6689fd998f8147c9cae07fbfa63af",
                "timestamp": 1427927671,
                "fee": 0,
                "version": 0,
                "tx_body_hash": "a6a709e9388a4d67a47d262b11da5f804eddd9d67acc4a3e450f7a567bdc1619"
            },
            "body": {
                "txns": [
                    {
                        "length": 183,
                        "type": 0,
                        "txid": "a6a709e9388a4d67a47d262b11da5f804eddd9d67acc4a3e450f7a567bdc1619",
                        "inner_hash": "ea6adee3180c7f9d73d1e693822d5d1c2bba85067f89a873355bc771a078faa1",
                        "sigs": [
                            "ce8fd47e2044ed17998f92621e90329f673a746c802d67f639ca083705dd199f6ee346781497b44132434922879244d819694b5903093f784570c55d293ab4af01"
                        ],
                        "inputs": [
                            "af0b2c1cc882a56b6c0c06e99e7d2731413b988329a2c47a5c2aa8be589b707a"
                        ],
                        "outputs": [
                            {
                                "uxid": "9eb7954461ba0256c9054fe38c00c66e60428dccf900a62e74b9fe39310aea13",
                                "dst": "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
                                "coins": "10.000000",
                                "hours": 0
                            }
                        ]
                    }
                ]
            },
            "size": 183
        },
        {
            "header": {
                "seq": 5,
                "block_hash": "114fe60587a158428a47e0f9571d764f495912c299aa4e67fc88004cf21b0c24",
                "previous_block_hash": "415e47348a1e642cb2e31d00ee500747d3aed0336aabfff7d783ed21465251c7",
                "timestamp": 1428798821,
                "fee": 2036,
                "version": 0,
                "tx_body_hash": "0579e7727627cd9815a8a8b5e1df86124f45a4132cc0dbd00d2f110e4f409b69"
            },
            "body": {
                "txns": [
                    {
                        "length": 317,
                        "type": 0,
                        "txid": "0579e7727627cd9815a8a8b5e1df86124f45a4132cc0dbd00d2f110e4f409b69",
                        "inner_hash": "fe123ca954a82bb1ce2cc9ef9c56d6b649a4cbaf5b17394b0ffda651ed32327e",
                        "sigs": [
                            "056ed0f74367fb1370d7e98689953983d9cf34eb6669854f1645c8a16c93d85075661e7d4f6df0ce5ca8eb9852eff6a12fbac2caafee03bb8c616f847c61416800",
                            "8aaa7f320a7b01169d3217a600100cb27c55e4ce56cd3455814f56d8e4e65be746e0e20e776087af6f19361f0b898edc2123a5f9bd35d24ef8b8669ca85b142601"
                        ],
                        "inputs": [
                            "9eb7954461ba0256c9054fe38c00c66e60428dccf900a62e74b9fe39310aea13",
                            "706f82c481906108880d79372ab5c126d32ecc98cf3f7c74cf33f5fda49dcf70"
                        ],
                        "outputs": [
                            {
                                "uxid": "fa2b598d233fe434f907f858d5de812eacf50c7b3fd152c77cd6e246fe356a9e",
                                "dst": "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
                                "coins": "999890.000000",
                                "hours": 4073
                            },
                            {
                                "uxid": "dc63c680f408c4e646037966189383a5d50eda34e666c2a0c75c0c6bf13b71a1",
                                "dst": "2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF",
                                "coins": "100.000000",
                                "hours": 0
                            }
                        ]
                    }
                ]
            },
            "size": 317
        },
        {
            "header": {
                "seq": 7,
                "block_hash": "6cb71b57c998a5367101e01d48c097eccd4f5abf311c89bcca8ee213581f355f",
                "previous_block_hash": "103949030e90fcebc5d8ca1c9c59f30a31aa71911401d22a2422e4571b035701",
                "timestamp": 1428807671,
                "fee": 0,
                "version": 0,
                "tx_body_hash": "f832428481690fa918d6d29946e191f2c8c89b2388a906e0c53dceee6070a24b"
            },
            "body": {
                "txns": [
                    {
                        "length": 220,
                        "type": 0,
                        "txid": "f832428481690fa918d6d29946e191f2c8c89b2388a906e0c53dceee6070a24b",
                        "inner_hash": "f440c514779522a6387edda9b9d9835f00680fb314546efb7bc9762a17884156",
                        "sigs": [
                            "8fe96f5502270e4efa962b2aef2b81795fe26a8f0c9a494e2ae9c7e624af455c49396270ae7a25b41d439fd56dea9d556a135129122de1b1274b1e2a5d75f2ea01"
                        ],
                        "inputs": [
                            "8ff8a647e4542fab01e078ac467b2c9f2e5f7de55d77ec2711f8abc718e2c91b"
                        ],
                        "outputs": [
                            {
                                "uxid": "17090c40091d009d6a684043d3be2e9cb1dc60a664a9c2e388af1f3a7345724b",
                                "dst": "2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF",
                                "coins": "90.000000",
                                "hours": 0
                            },
                            {
                                "uxid": "f9e7a412cdff80e95ddbe1d76fcc73f967cb99d383b0659e1355c8e623f02b62",
                                "dst": "WADSeEwEQVbtUy8CfcVimyxX1KjTRkvfoK",
                                "coins": "5.000000",
                                "hours": 0
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


### Get last N blocks

API sets: `READ`

```
URI: /api/v1/last_blocks
Method: GET
Args:
    num: number of most recent blocks to return
    verbose: [bool] return verbose transaction input data
```

If verbose, the transaction inputs include the owner address, coins, hours and calculated hours.
The hours are the original hours the output was created with.
The calculated hours are the hours the transaction had in the block in which it was executed.

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

Example (verbose):

```sh
curl http://127.0.0.1:6420/api/v1/last_blocks?num=2&verbose=1
```

Result:

```json
{
    "blocks": [
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
                        "fee": 970389,
                        "sigs": [
                            "2ff7390c3b66c6b0fbb2b4c59c8e218291d4cbb82a836bb577c7264677f4a8320f6f3ad72d804e3014728baa214c223ecced8725b64be96fe3b51332ad1eda4201",
                            "9e7c715f897b3c987c00ee8c6b14e4b90bb3e4e11d003b481f82042b1795b3c75eaa3d563cd0358cdabdab77cfdbead7323323cf73e781f9c1a8cf6d9b4f8ac100",
                            "5c9748314f2fe0cd442df5ebb8f211087111d22e9463355bf9eee583d44df1bd36addb510eb470cb5dafba0732615f8533072f80ae05fc728c91ce373ada1e7b00"
                        ],
                        "inputs": [
                            {
                                "uxid": "5f634c825b2a53103758024b3cb8578b17d56d422539e23c26b91ea397161703",
                                "owner": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "25.910000",
                                "hours": 7745,
                                "calculated_hours": 17458
                            },
                            {
                                "uxid": "16ac52084ffdac2e9169b9e057d44630dec23d18cfb90b9437d28220a3dc585d",
                                "owner": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "1.000000",
                                "hours": 1915246,
                                "calculated_hours": 1915573
                            },
                            {
                                "uxid": "8d3263890d32382e182b86f8772c7685a8f253ed475c05f7d530e9296f692bc9",
                                "owner": "2Huip6Eizrq1uWYqfQEh4ymibLysJmXnWXS",
                                "coins": "0.003000",
                                "hours": 7745,
                                "calculated_hours": 7746
                            }
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
                        "fee": 485194,
                        "sigs": [
                            "af5329e77213f34446a0ff41d249fd25bc1dae913390871df359b9bd587c95a10b625a74a3477a05cc7537cb532253b12c03349ead5be066b8e0009e79462b9501"
                        ],
                        "inputs": [
                            {
                                "uxid": "fb8db3f78928aee3f5cbda8db7fc290df9e64414e8107872a1c5cf83e08e4df7",
                                "owner": "uvcDrKc8rHTjxLrU4mPN56Hyh2tR6RvCvw",
                                "coins": "26.913000",
                                "hours": 970388,
                                "calculated_hours": 970388
                            }
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

## Uxout APIs

### Get uxout

API sets: `READ`

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

### Get historical unspent outputs for an address

API sets: `READ`

```
URI: /api/v1/address_uxouts
Method: GET
Args:
    address
```

Returns the historical, spent outputs of a given address.

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

API sets: `READ`

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

API sets: `READ`

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

### Count the addresses that currently have unspent outputs (coins)

API sets: `READ`

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

API sets: `STATUS`, `READ`

```
URI: /api/v1/network/connection
Method: GET
Args:
    addr: ip:port address of a known connection
```

Connection `"state"` value can be `"pending"`, `"connected"` or `"introduced"`.

* The `"pending"` state is prior to connection establishment.
* The `"connected"` state is after connection establishment, but before the introduction handshake has completed.
* The `"introduced"` state is after the introduction handshake has completed.

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
    "connected_at": 1520675700,
    "outgoing": false,
    "state": "introduced",
    "mirror": 719118746,
    "height": 181,
    "listen_port": 6000,
    "user_agent": "skycoin:0.25.0",
    "is_trusted_peer": true,
    "unconfirmed_verify_transaction": {
        "burn_factor": 10,
        "max_transaction_size": 32768,
        "max_decimals": 3
    }
}
```

### Get a list of all connections

API sets: `STATUS`, `READ`

```
URI: /api/v1/network/connections
Method: GET
Args:
    states: [optional] comma-separated list of connection states ("pending", "connected" or "introduced"). Defaults to "connected,introduced"
    direction: [optional] "outgoing" or "incoming". If not provided, both are included.
```

Connection `"state"` value can be `"pending"`, `"connected"` or `"introduced"`.

* The `"pending"` state is prior to connection establishment.
* The `"connected"` state is after connection establishment, but before the introduction handshake has completed.
* The `"introduced"` state is after the introduction handshake has completed.

By default, both incoming and outgoing connections in the `"connected"` or `"introduced"` state are returned.

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
            "connected_at": 1520675500,
            "outgoing": false,
            "state": "introduced",
            "mirror": 1338939619,
            "listen_port": 20002,
            "height": 180,
            "user_agent": "skycoin:0.25.0",
            "is_trusted_peer": true,
            "unconfirmed_verify_transaction": {
                "burn_factor": 10,
                "max_transaction_size": 32768,
                "max_decimals": 3
            }
        },
        {
            "id": 109548,
            "address": "176.9.84.75:6000",
            "last_sent": 1520675751,
            "last_received": 1520675751,
            "connected_at": 1520675751,
            "state": "connected",
            "outgoing": true,
            "mirror": 0,
            "listen_port": 6000,
            "height": 0,
            "user_agent": "",
            "is_trusted_peer": true,
            "unconfirmed_verify_transaction": {
                "burn_factor": 0,
                "max_transaction_size": 0,
                "max_decimals": 0
            }
        },
        {
            "id": 99115,
            "address": "185.120.34.60:6000",
            "last_sent": 1520675754,
            "last_received": 1520675754,
            "connected_at": 1520673013,
            "outgoing": false,
            "state": "introduced",
            "mirror": 1931713869,
            "listen_port": 6000,
            "height": 180,
            "user_agent": "",
            "is_trusted_peer": true,
            "unconfirmed_verify_transaction": {
                "burn_factor": 0,
                "max_transaction_size": 0,
                "max_decimals": 0
            }
        }
    ]
}
```


### Get a list of all default connections

API sets: `STATUS`, `READ`

```
URI: /api/v1/network/defaultConnections
Method: GET
```

Returns addresses in the default hardcoded list of peers.

Example:

```sh
curl 'http://127.0.0.1:6420/api/v1/network/defaultConnections'
```

Result:

```json
[
    "104.237.142.206:6000",
    "118.178.135.93:6000",
    "139.162.7.132:6000",
    "172.104.85.6:6000",
    "176.58.126.224:6000",
    "47.88.33.156:6000"
]
```

### Get a list of all trusted connections

API sets: `STATUS`, `READ`

```
URI: /api/v1/network/connections/trust
Method: GET
```

Returns addresses marked as trusted in the peerlist.
This is typically equal to the list of addresses in the default hardcoded list of peers.

Example:

```sh
curl 'http://127.0.0.1:6420/api/v1/network/connections/trust'
```

Result:

```json
[
    "104.237.142.206:6000",
    "118.178.135.93:6000",
    "139.162.7.132:6000",
    "172.104.85.6:6000",
    "176.58.126.224:6000",
    "47.88.33.156:6000"
]
```

### Get a list of all connections discovered through peer exchange

API sets: `STATUS`, `READ`

```
URI: /api/v1/network/connections/exchange
Method: GET
```

Returns addresses from the peerlist that are known to have an open port.

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

### Disconnect a peer

API sets: `NET_CTRL`

```
URI: /api/v1/network/connection/disconnect
Method: POST
Args:
    id: ID of the connection

Returns 404 if the connection is not found.
```

Disconnects a peer by ID.

Example:

```sh
curl -X POST 'http://127.0.0.1:6420/api/v1/network/connection/disconnect?id=999'
```

Result:

```json
{}
```

## Migrating from the unversioned API

The unversioned API are the API endpoints without an `/api` prefix.
These endpoints are all prefixed with `/api/v1` now.

`-enable-unversioned-api` was added as an option to assist migration to `/api/v1`
but this option was removed in v0.26.0.

To migrate from the unversioned API, add `/api/v1` to all endpoints that you call
that do not have an `/api` prefix already.

For example, `/block` would become `/api/v1/block`.

## Migrating from the JSONRPC API

The JSONRPC-2.0 RPC API was deprecated in v0.25.0 and removed in v0.26.0.

Anyone still using this can follow this guide to migrate to the REST API:

* `get_status` is replaced by `/api/v1/blockchain/metadata` and `/api/v1/health`
* `get_lastblocks` is replaced by `/api/v1/last_blocks`
* `get_blocks` is replaced by `/api/v1/blocks`
* `get_outputs` is replaced by `/api/v1/outputs`
* `inject_transaction` is replaced by `/api/v1/injectTransaction`
* `get_transaction` is replaced by `/api/v1/transaction`

## Migrating from /api/v1/spend

The `POST /api/v1/spend` endpoint is deprecated and will be removed in v0.26.0.

To migrate from it, use [`POST /api/v1/wallet/transaction`](#create-transaction) followed by [`POST /api/v1/injectTransaction`](#inject-raw-transaction).
Do not create another transaction before injecting the created transaction, otherwise you might create two conflicting transactions.

`POST /api/v1/wallet/transaction` has more options for creating the transaction than the `/api/v1/spend` endpoint.
To replicate the same behavior as `/api/v1/spend`, use the following request body template:

```json
{
    "hours_selection": {
        "type": "auto",
        "mode": "share",
        "share_factor": "0.5",
    },
    "wallet": {
        "id": "$wallet_id",
        "password": "$password"
    },
    "to": [{
        "address": "$dst",
        "coins": "$coins"
    }]
}
```

You must use a string for `"coins"` instead of an integer measured in "droplets" (the smallest unit of currency in Skycoin, 1/1000000 of a skycoin).
For example, if you sent 1 Skycoin with `/api/v1/spend` you would have specified the `coins` field as `1000000`.
Now, you would specify it as `"1"`.

Some examples:

* 123.456 coins: before `123456000`, now `"123.456"`
* 0.1 coins: before `100000`, now `"0.1"`
* 1 coin: before `1000000`, now `"1"`

Extra zeros on the `"coins"` string are ok, for example `"1"` is the same as `"1.0"` or `"1.000000"`.

Only provide `"password"` if the wallet is encrypted. Note that decryption can take a few seconds, and this can impact
throughput.

The request header `Content-Type` must be `application/json`.

The response to `POST /api/v1/wallet/transaction` will include a verbose decoded transaction with details
and the hex-encoded binary transaction in the `"encoded_transaction"` field.
Use the value of `"encoded_transaction"` as the `"rawtx"` value in the request to `/api/v1/injectTransaction`.

## Migration from /api/v1/explorer/address

The `GET /api/v1/explorer/address` was deprecated in v0.25.0 and removed in v0.26.0.

To migrate from it, use [`GET /api/v1/transactions?verbose=1`](#get-transactions-for-addresses).

`/api/v1/explorer/address` accepted a single `address` query parameter. `/api/v1/transactions` uses an `addrs` query parameter and
accepts multiple addresses at once.

The response data is the same but the structure is slightly different. Compare the follow two example responses:

`/api/v1/explorer/address?address=WzPDgdfL1NzSbX96tscUNXUqtCRLjaBugC`:

```json
[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 38076,
            "block_seq": 15493
        },
        "timestamp": 1518878675,
        "length": 183,
        "type": 0,
        "txid": "6d8e2f8b436a2f38d604b3aa1196ef2176779c5e11e33fbdd09f993fe659c39f",
        "inner_hash": "8da7c64dcedeeb6aa1e0d21fb84a0028dcd68e6801f1a3cc0224fdd50682046f",
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

`/api/v1/transactions?verbose=1&addrs=WzPDgdfL1NzSbX96tscUNXUqtCRLjaBugC`:

```json
[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 57564,
            "block_seq": 7498
        },
        "time": 1514743602,
        "txn": {
            "timestamp": 1514743602,
            "length": 220,
            "type": 0,
            "txid": "df5bcef198fe6e96d496c30482730f895cabc1d55b338afe5633b0c2889d02f9",
            "inner_hash": "4677ff9b9b56485495a45693cc09f8496199929fccb52091d32f2d3cf2ee8a41",
            "fee": 69193,
            "sigs": [
                "8e1f6f621a11f737ac2031be975d4b2fc17bf9f17a0da0a2fe219ee018011ab506e2ad0367be302a8d859cc355c552313389cd0aa9fa98dc7d2085a52f11ef5a00"
            ],
            "inputs": [
                {
                    "uxid": "2374201ff29f1c024ccfc6c53160e741d06720562853ad3613c121acd8389031",
                    "owner": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                    "coins": "162768.000000",
                    "hours": 485,
                    "calculated_hours": 138385
                }
            ],
            "outputs": [
                {
                    "uxid": "63f299fc85fe6fc34d392718eee55909837c7231b6ffd93e5a9a844c4375b313",
                    "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                    "coins": "162643.000000",
                    "hours": 34596
                },
                {
                    "uxid": "349b06e5707f633fd2d8f048b687b40462d875d968b246831434fb5ab5dcac38",
                    "dst": "WzPDgdfL1NzSbX96tscUNXUqtCRLjaBugC",
                    "coins": "125.000000",
                    "hours": 34596
                }
            ]
        }
    }
]
```

The transaction data is wrapped in a `"txn"` field. A `"time"` field is present at the top level. This `"time"` field
is either the confirmation timestamp of a confirmed transaction or the last received time of an unconfirmed transaction.
