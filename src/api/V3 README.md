# REST API specification v3

## API Version 2

`/api/v3` endpoints have a standard format.

All `/api/v3` `POST` endpoints accept only `application/json` and return `application/json`.

All `/api/v3` `GET` requires accept data in JSON body.
In the future we may have choose to have `GET` requests also accept `POST` with a JSON body,
to support requests with a large query body, such as when requesting data for a large number
of addresses or transactions.

`/api/v3` responses are always JSON. If there is an error, the JSON object will
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
URI: /api/v3/csrf
Method: GET
```

Example:

```sh
curl http://127.0.0.1:6420/api/v3/csrf
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
curl http://127.0.0.1:6420/api/v3/health
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
    "unversioned_api_enabled": false,
    "json_rpc_enabled": false,
    "user_verify_transaction": {
        "burn_factor": 2,
        "max_transaction_size": 32768,
        "max_decimals": 3
    },
    "unconfirmed_verify_transaction": {
        "burn_factor": 2,
        "max_transaction_size": 32768,
        "max_decimals": 3
    },
    "started_at": 1542443907
}
```
