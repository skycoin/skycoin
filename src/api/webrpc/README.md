# Webrpc

This is a description about skycoin webrpc, which implemented the [json-rpc 2.0](http://www.jsonrpc.org/specification) protocol.
The rpc service entry point is /webrpc, and only accept the HTTP `POST` requests.

## Get Status

Get status of rpc server.

request:

```json
{
    "id": "1",
    "jsonrpc": "2.0",
    "method": "get_status"
}
```

## Get last blocks

Get last `N` blocks.

request:

```json
{
    "id": "1",
    "jsonrpc": "2.0",
    "method": "get_lastblocks",
    "params": [3]
}
```

The params must be an array with one integer value.

## Get blocks

Get blocks in specific range, inclusive.

request:

```json
{
    "id": "1",
    "jsonrpc": "2.0",
    "method": "get_blocks",
    "params": [2, 10]
}
```

The params must be an array with two integer values.

## Get blocks by sequence number

Get blocks at specific sequence numbers.

request:

```json
{
    "id": "1",
    "jsonrpc": "2.0",
    "method": "get_blocks",
    "params": [133, 401, 212]
}
```

The params must be an array of integer values.

## Get outputs

Get unspent outputs of specific addresses.

request:

```json
{
    "id": "1",
    "jsonrpc": "2.0",
    "method": "get_outputs",
    "params": ["fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4C", "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B"]
}
```

The params must be an array of strings.

## Inject transaction

Broadcast raw transaction.

request:

```json
{
    "id": "1",
    "jsonrpc": "2.0",
    "method": "inject_transaction",
    "params": ["dc0000000010e05181fd4023f865a84359bf72a304e687b6f00e42f93ad9a4b8ee5a64aabc01000000dcb5b236eecd97a36c7d0a0b8ed68bb5df6274433a51fddf911f02f3926d20bf6eaabdc21529b7696f498545b06cc7e69f2f08b4dc5fa823c5b3f03da06794a300010000006d8a9c89177ce5e9d3b4b59fff67c00f0471fdebdfbb368377841b03fc7d688b02000000005771eeda2e253697cf5368f16fe05210d5cd319040420f0000000000af010000000000000060dfa95881cdc827b45a6d49b11dbc152ecd4de600093d0000000000af01000000000000"]
}
```

The params must be an array with one raw transaction string.

## Get transaction

Get transaction verbose info of specific transaction id.

request:

```json
{
    "id": "1",
    "jsonrpc": "2.0",
    "method": "get_transaction",
    "params": ["bdc4a85a3e9d17a8fe00aa7430d0347c7f1dd6480a16da7147b6e43905057d43"]
}
```

The params must be an array with one txid string.
