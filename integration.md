# Skycoin Exchange Integration

A Skycoin node offers multiple interfaces:

* REST API on port 6420
* JSON-RPC 2.0 API on port 6430 **[deprecated]**

A CLI tool is provided in `cmd/cli/cli.go`. This tool communicates over the JSON-RPC 2.0 API.

*Note*: Do not interface with the JSON-RPC 2.0 API directly, it is deprecated and will be removed in a future version.

The API interfaces do not support authentication or encryption so they should only be used over localhost.

## API Documentation

### Wallet REST API

[Wallet REST API](src/gui/README.md).

### Skycoin command line interface

[CLI command API](cmd/cli/README.md).

## Implementation guidelines

### Scanning deposits

TODO

### Sending coins

#### General principles

After each spend, wait for the transaction to confirm before trying to spend again.

For higher throughput, combine multiple spends into one transaction.

Skycoin uses "coin hours" to ratelimit transactions.
The total number of coinhours in a transaction's outputs must be 50% or less than the number of coinhours in a transaction's inputs,
or else the transaction is invalid and will not be accepted.  A transaction must have at least 1 input with at least 1 coin hour.
Sending too many transactions in quick succession will use up all available coinhours.
Coinhours are earned at a rate of 1 coinhour per coin per hour, calculated per second.
This means that 3600 coins will earn 1 coinhour per second.
However, coinhours are only updated when a new block is published to the blockchain.
New blocks are published every 10 seconds, but only if there are pending transactions in the network.

To avoid running out of coinhours in situations where the application may frequently send,
the sender should batch sends into a single transaction and send them on a
30 second to 1 minute interval.

There are other strategies to minimize the likelihood of running out of coinhours, such
as splitting up balances into many unspent outputs and having a large balance which generates
coinhours quickly.

#### Using the CLI

When sending coins from the CLI tool, a wallet file local to the caller is used.
The CLI tool allows you to specify the wallet file on disk to use for operations.

See [CLI command API](cmd/cli/README.md) for documentation of the CLI interface.

To perform a send, the preferred method follows these steps in a loop:

* `skycoin-cli createRawTransaction -m '[{"addr:"$addr1,"coins:"$coins1"}, ...]` - `-m` flag is send-to-many
* `skycoin-cli broadcastTransaction` - returns `txid`
* `skycoin-cli transaction $txid` - repeat this command until `"status"` is `"confirmed"`

That is, create a raw transaction, broadcast it, and wait for it to confirm.

#### Using the REST API

When sending coins via the REST API, a wallet file local to the skycoin node is used.
The wallet file is specified by wallet ID, and all wallet files are in the
configured data directory (which is `$HOME/.skycoin/wallets` by default).

#### Using skycoin as a library in a Go application

If your application is written in Go, you can interface with the CLI library
directly, see [Skycoin CLI Godoc](https://godoc.org/github.com/skycoin/skycoin/src/api/cli).

A REST API client is also available: [Skycoin REST API Client Godoc](https://godoc.org/github.com/skycoin/skycoin/src/gui#Client).

### Verifying addresses

#### Using the CLI

```sh
skycoin-cli verifyAddress $addr
```

#### Using the REST API

Not supported.

#### Using skycoin as a library in a Go application

https://godoc.org/github.com/skycoin/skycoin/src/cipher#DecodeBase58Address

```go
if _, err := cipher.DecodeBase58Address(address); err != nil {
    fmt.Println("Invalid address:", err)
    return
}
```

#### Using skycoin as a library in other applications

Some methods of the skycoin go library are available through a C wrapper, `libskycoin`.

To build the library,

```sh
make build-libc
```

