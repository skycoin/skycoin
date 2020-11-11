# Skycoin Daemon CLI Options

<!-- MarkdownTOC levels="1,2,3,4,5" autolink="true" bracket="round" -->

- [Preface](#preface)
- [Scenarios](#scenarios)
	- [Control which API endpoints are exposed](#control-which-api-endpoints-are-exposed)
	- [Run a public API node](#run-a-public-api-node)
	- [Run a public API node with a self-signed cert](#run-a-public-api-node-with-a-self-signed-cert)
	- [Control which peers the node connects to](#control-which-peers-the-node-connects-to)
	- [Add Basic auth to the REST API interface](#add-basic-auth-to-the-rest-api-interface)
- [Options](#options)
	- [address](#address)
	- [block-publisher](#block-publisher)
	- [blockchain-public-key](#blockchain-public-key)
	- [blockchain-secret-key](#blockchain-secret-key)
	- [burn-factor-create-block](#burn-factor-create-block)
	- [burn-factor-unconfirmed](#burn-factor-unconfirmed)
	- [color-log](#color-log)
	- [connection-rate](#connection-rate)
	- [custom-peers-file](#custom-peers-file)
	- [data-dir](#data-dir)
	- [db-path](#db-path)
	- [db-read-only](#db-read-only)
	- [disable-api-sets](#disable-api-sets)
	- [disable-csp](#disable-csp)
	- [disable-csrf](#disable-csrf)
	- [disable-default-peers](#disable-default-peers)
	- [disable-header-check](#disable-header-check)
	- [disable-incoming](#disable-incoming)
	- [disable-outgoing](#disable-outgoing)
	- [disable-pex](#disable-pex)
	- [download-peerlist](#download-peerlist)
	- [enable-all-api-sets](#enable-all-api-sets)
	- [enable-api-sets](#enable-api-sets)
	- [enable-gui](#enable-gui)
	- [genesis-address](#genesis-address)
	- [genesis-signature](#genesis-signature)
	- [genesis-timestamp](#genesis-timestamp)
	- [gui-dir](#gui-dir)
	- [host-whitelist](#host-whitelist)
	- [http-prof](#http-prof)
	- [http-prof-host](#http-prof-host)
	- [launch-browser](#launch-browser)
	- [localhost-only](#localhost-only)
	- [log-level](#log-level)
	- [logtofile](#logtofile)
	- [max-block-size](#max-block-size)
	- [max-connections](#max-connections)
	- [max-decimals-create-block](#max-decimals-create-block)
	- [max-decimals-unconfirmed](#max-decimals-unconfirmed)
	- [max-default-peer-outgoing-connections](#max-default-peer-outgoing-connections)
	- [max-incoming-connections](#max-incoming-connections)
	- [max-in-msg-len](#max-in-msg-len)
	- [max-out-msg-len](#max-out-msg-len)
	- [max-outgoing-connections](#max-outgoing-connections)
	- [max-txn-size-create-block](#max-txn-size-create-block)
	- [max-txn-size-unconfirmed](#max-txn-size-unconfirmed)
	- [no-ping-log](#no-ping-log)
	- [peerlist-size](#peerlist-size)
	- [peerlist-url](#peerlist-url)
	- [port](#port)
	- [profile-cpu](#profile-cpu)
	- [profile-cpu-file](#profile-cpu-file)
	- [reset-corrupt-db](#reset-corrupt-db)
	- [storage-dir](#storage-dir)
	- [user-agent-remark](#user-agent-remark)
	- [verify-db](#verify-db)
	- [version](#version)
	- [wallet-crypto-type](#wallet-crypto-type)
	- [wallet-dir](#wallet-dir)
	- [web-interface](#web-interface)
	- [web-interface-addr](#web-interface-addr)
	- [web-interface-cert](#web-interface-cert)
	- [web-interface-https](#web-interface-https)
	- [web-interface-key](#web-interface-key)
	- [web-interface-password](#web-interface-password)
	- [web-interface-plaintext-auth](#web-interface-plaintext-auth)
	- [web-interface-port](#web-interface-port)
	- [web-interface-username](#web-interface-username)
- [Development Environment Variables](#development-environment-variables)
	- [USER_BURN_FACTOR](#userburnfactor)
	- [USER_MAX_TXN_SIZE](#usermax_txnsize)
	- [USER_MAX_DECIMALS](#usermaxdecimals)

<!-- /MarkdownTOC -->


## Preface

*Note: The defaults shown below can vary depending on the build configuration.*

```
» go run cmd/skycoin/skycoin.go --help

Usage:
  -address string
    	IP Address to run application on. Leave empty to default to a public interface
  -block-publisher
    	run the daemon as a block publisher
  -blockchain-public-key string
    	public key of the blockchain (default "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a")
  -blockchain-secret-key string
    	secret key of the blockchain
  -burn-factor-create-block uint
    	coinhour burn factor applied when creating blocks (default 10)
  -burn-factor-unconfirmed uint
    	coinhour burn factor applied to unconfirmed transactions (default 10)
  -color-log
    	Add terminal colors to log output (default true)
  -connection-rate duration
    	How often to make an outgoing connection (default 5s)
  -custom-peers-file string
    	load custom peers from a newline separate list of ip:port in a file. Note that this is different from the peers.json file in the data directory
  -data-dir string
    	directory to store app data (defaults to ~/.skycoin) (default "$HOME/.skycoin")
  -db-path string
    	path of database file (defaults to ~/.skycoin/data.db)
  -db-read-only
    	open bolt db read-only
  -disable-api-sets string
    	disable API set. Options are READ, STATUS, WALLET, TXN, PROMETHEUS, NET_CTRL, INSECURE_WALLET_SEED, STORAGE. Multiple values should be separated by comma
  -disable-csp
    	disable content-security-policy in http response
  -disable-csrf
    	disable CSRF check
  -disable-default-peers
    	disable the hardcoded default peers
  -disable-header-check
    	disables the host, origin and referer header checks.
  -disable-incoming
    	Don't allow incoming connections
  -disable-networking
    	Disable all network activity
  -disable-outgoing
    	Don't make outgoing connections
  -disable-pex
    	disable PEX peer discovery
  -download-peerlist
    	download a peers.txt from -peerlist-url (default true)
  -enable-all-api-sets
    	enable all API sets, except for deprecated or insecure sets. This option is applied before -disable-api-sets.
  -enable-api-sets string
    	enable API set. Options are READ, STATUS, WALLET, TXN, PROMETHEUS, NET_CTRL, INSECURE_WALLET_SEED, STORAGE. Multiple values should be separated by comma (default "READ,TXN")
  -enable-gui
    	Enable GUI
  -genesis-address string
    	genesis address (default "2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6")
  -genesis-signature string
    	genesis block signature (default "eb10468d10054d15f2b6f8946cd46797779aa20a7617ceb4be884189f219bc9a164e56a5b9f7bec392a804ff3740210348d73db77a37adb542a8e08d429ac92700")
  -genesis-timestamp uint
    	genesis block timestamp (default 1426562704)
  -gui-dir string
    	static content directory for the HTML interface (default "./src/gui/static/")
  -help
    	Show help
  -host-whitelist string
    	Hostnames to whitelist in the Host header check. Only applies when the web interface is bound to localhost.
  -http-prof
    	run the HTTP profiling interface
  -http-prof-host string
    	hostname to bind the HTTP profiling interface to (default "localhost:6060")
  -launch-browser
    	launch system default webbrowser at client startup
  -localhost-only
    	Run on localhost and only connect to localhost peers
  -log-level string
    	Choices are: debug, info, warn, error, fatal, panic (default "INFO")
  -logtofile
    	log to file
  -max-block-size uint
    	maximum total size of transactions in a block (default 32768)
  -max-connections int
    	Maximum number of total connections allowed (default 128)
  -max-decimals-create-block uint
    	max number of decimal places applied when creating blocks (default 3)
  -max-decimals-unconfirmed uint
    	max number of decimal places applied to unconfirmed transactions (default 3)
  -max-default-peer-outgoing-connections int
    	The maximum default peer outgoing connections allowed (default 1)
  -max-in-msg-len int
    	Maximum length of incoming wire messages (default 1048576)
  -max-out-msg-len int
    	Maximum length of outgoing wire messages (default 262144)
  -max-outgoing-connections int
    	Maximum number of outgoing connections allowed (default 8)
  -max-incoming-connections int
        Maximum number  of incoming connections allowed (default 120)
  -max-txn-size-create-block uint
    	maximum size of a transaction applied when creating blocks (default 32768)
  -max-txn-size-unconfirmed uint
    	maximum size of an unconfirmed transaction (default 32768)
  -no-ping-log
    	disable "reply to ping" and "received pong" debug log messages
  -peerlist-size int
    	Max number of peers to track in peerlist (default 65535)
  -peerlist-url string
    	with -download-peerlist=true, download a peers.txt file from this url (default "https://downloads.skycoin.com/blockchain/peers.txt")
  -port int
    	Port to run application on (default 6000)
  -profile-cpu
    	enable cpu profiling
  -profile-cpu-file string
    	where to write the cpu profile file (default "cpu.prof")
  -reset-corrupt-db
    	reset the database if corrupted, and continue running instead of exiting
  -storage-dir string
    	location of the storage data files. Defaults to ~/.skycoin/data/
  -user-agent-remark string
    	additional remark to include in the user agent sent over the wire protocol
  -verify-db
    	check the database for corruption
  -version
    	show node version
  -wallet-crypto-type string
    	wallet crypto type. Can be sha256-xor or scrypt-chacha20poly1305 (default "scrypt-chacha20poly1305")
  -wallet-dir string
    	location of the wallet files. Defaults to ~/.skycoin/wallet/
  -web-interface
    	enable the web interface (default true)
  -web-interface-addr string
    	addr to serve web interface on (default "127.0.0.1")
  -web-interface-cert string
    	skycoind.cert file for web interface HTTPS. If not provided, will autogenerate or use skycoind.cert in --data-dir
  -web-interface-https
    	enable HTTPS for web interface
  -web-interface-key string
    	skycoind.key file for web interface HTTPS. If not provided, will autogenerate or use skycoind.key in --data-dir
  -web-interface-password string
    	password for the web interface
  -web-interface-plaintext-auth
    	allow web interface auth without https
  -web-interface-port int
    	port to serve web interface on (default 6420)
  -web-interface-username string
    	username for the web interface
```

## Scenarios

### Control which API endpoints are exposed

API endpoints are grouped into "API sets", and can be toggled on or off through CLI options.

There are three options for controlling the API sets that are enabled:

* `--enable-api-sets`
* `--enable-all-api-sets`
* `--disable-api-sets`

To blacklist specific API sets, combine `--enable-all-api-sets` with `--disable-api-sets`.

To whitelist specific API sets, use `--enable-api-sets`.

Note that certain API sets must be explicitly enabled. These API sets are either deprecated or have security implications.

Read more about API sets here: https://github.com/skycoin/skycoin/blob/develop/src/api/README.md#api-sets

### Run a public API node

This example does not use HTTPS. We recommend you follow the nginx guide below
to add HTTPS to your API node.

```sh
$ go run cmd/skycoin/skycoin.go \
  --web-interface-addr=0.0.0.0 \
  --enable-api-sets=READ,TXN
```

This will expose your node's API on your server's public IP and on the default port 6420.

### Run a public API node with a self-signed cert

When you run with the HTTPS option, the daemon will use a cert and key file from the `data-dir`.
This cert and key file will be autogenerated if they are missing from the `data-dir`.

```sh
go run cmd/skycoin/skycoin.go \
  --web-interface-https \
  --web-interface-addr=0.0.0.0 \
  --enable-api-sets=READ,TXN
```

Alternatively, you can specify the cert and key files to be some other location. If specified, they will not be autogenerated.
The specified files must exist:

```sh
go run cmd/skycoin/skycoin.go \
  --web-interface-https \
  --web-interface-key=/var/local/skycoind.key \
  --web-interface-cert=/var/local/skycoind.cert \
  --enable-api-sets=READ,TXN
```

If you want to use a signed cert, we recommend handling that with nginx, which will also help you bind to port 80/443.

### Control which peers the node connects to

First, make sure the `peers.json` file in the `data-dir` is empty or does not exist.

Provide a `custom-peers-file`, which is a newline separated list of ip:port entries.

Disable the default bootstrap peers, and disable the remote peerlist bootstrap.

There is no explicit setting for max incoming connections; it is equal to the difference between `--max-connections` and `--max-outgoing-connections`,
so set these two options to the same value.

```sh
go run cmd/skycoin/skycoin.go \
  --custom-peers-file=peers-whitelist.txt \
  --disable-default-peers \
  --download-peerlist=false \
  --max-connections=8 \
  --max-outgoing-connections=8 \
  --disable-pex \
  --disable-incoming
```

### Add Basic auth to the REST API interface

This will enable `Basic` auth on the REST API interface. It will use HTTPS with an autogenerated self-signed cert.
Your client will need to be configured to accept this cert in requests.

*TODO - describe how to configure the cli client to use the self-signed cert*

```sh
$ go run cmd/skycoin/skycoin.go \
  --web-interface-https \
  --web-interface-username=abcdef \
  --web-interface-password='aCN@9xA)(CZasdmc'
```

## Options

### address

The bind interface address for the wire protocol. Binds to a public interface by default.

### block-publisher

Runs the node as a block publisher. Must set `blockchain-secret-key`.

### blockchain-public-key

The public key of the block signer

### blockchain-secret-key

The secret key of the block signer. Required for `block-publisher` mode.

### burn-factor-create-block

The coin hour burn factor applied to transactions when creating blocks.
Transactions that don't satisfy this burn factor will not be included in blocks.
Only applies when running in `block-publisher` mode.

### burn-factor-unconfirmed

The coin hour burn factor applied to unconfirmed transactions received over the network.
Transactions that don't satisfy this burn factor will not be propagated to peers.

### color-log

Use color highlighting in the log output. Disable this when logging to a file.

### connection-rate

How often an outgoing connection attempt is made.
A faster rate will establish a stable connection sooner, but if it is too fast
it can overconnect and churn connections.

### custom-peers-file

Load peers from this file into the peer database. The file format is a newline-separated list of ip:port entries.
These peers are *added* to any existing peer database; it does not restrict the peers to those in this file.

### data-dir

The storage location for application data. By default, the database, wallets, peers cache and other data files
will be saved into this folder.

On Linux and MacOS, this folder defaults to `$HOME/.skycoin` (`~/.skycoin`).
On Windows release builds, this folder defaults to `%HOMEPATH%\.skycoin` (`C:\Users\{user}\.skycoin`).
On Windows development builds, this folder defaults to `C:\.skycoin`. *(Note: this is a bug and will change in the future)*

### db-path

The path of the blockchain database file. Defaults to a file named `data.db` in `data-dir`.

### db-read-only

Open the database file in read-only mode.

### disable-api-sets

Disable one or more API sets. Possible API sets are:
`READ`, `STATUS`, `WALLET`, `TXN`, `PROMETHEUS`, `NET_CTRL`, `INSECURE_WALLET_SEED`, `STORAGE`.
Multiple values should be separated by comma. Combine with `enable-all-api-sets` to blacklist specific API sets.

Read more about API sets here: https://github.com/skycoin/skycoin/blob/develop/src/api/README.md#api-sets

### disable-csp

Disable the Content Security Policy header sent in REST API responses.

### disable-csrf

Disable the CSRF check for the REST API. The REST API requires a CSRF token for all `POST` and `DELETE` requests.
This is to protect the wallet client from certain attacks. If you do not have a hot wallet or do not expose your node
to a browser, it is safe to disable CSRF.

### disable-default-peers

Disable the default hardcoded peer list. These peers are treated differently than others; the node will always try to maintain
at least `max-default-peer-outgoing-connections` connections to peers in this list. These peers should be disabled when configuring the
node for network isolation.

### disable-header-check

As a security policy, the REST API will require certain values for the
`Host`, `Origin` and `Referer` headers in requests unless disabled by this option.

### disable-incoming

Disable all incoming connections on the wire interface.  The listener will not bind to the configured `address`.

### disable-outgoing

Don't make any outgoing connections.

### disable-pex

Don't request or accept peers over the wire.

### download-peerlist

If true, a peer list will be downloaded from `--peerlist-url`. The peer list file format is a newline-separated list of
ip:port entries. This list helps to bootstrap the initial peer database. These peers are considered "regular" peers, as opposed
to the peers from the hardcoded default peer list which are handled slightly differently.

### enable-all-api-sets

Enable all API sets except for those marked `INSECURE` or `DEPRECATED`.
Combine with `disable-api-sets` to blacklist specific API sets.
Use `enable-api-sets` in addition to `enable-all-api-sets` in order to enable specific `INSECURE` or `DEPRECATED` API sets.

Read more about API sets here: https://github.com/skycoin/skycoin/blob/develop/src/api/README.md#api-sets

### enable-api-sets

Enable one or more API sets. Possible API sets are:
`READ`, `STATUS`, `WALLET`, `TXN`, `PROMETHEUS`, `NET_CTRL`, `INSECURE_WALLET_SEED`, `STORAGE`.
Multiple values should be separated by comma.

Read more about API sets here: https://github.com/skycoin/skycoin/blob/develop/src/api/README.md#api-sets

### enable-gui

Serve the wallet GUI pages over the `web-interface-addr` and `web-interface-port` on the root path `/`.

### genesis-address

The genesis address in the genesis block.  This is used to reconstruct the genesis block, which is hardcoded in every client.

### genesis-signature

After the genesis block was created, it should have been signed by the `blockchain-secret-key`. This signature is configured here.

### genesis-timestamp

The timestamp of the genesis block. This is used to reconstruct the genesis, which is hardcoded in every client.

### gui-dir

The static content directory for the wallet GUI interface.

### host-whitelist

A comma separated list of hostnames to allow in the `Host`, `Origin` and `Referer` headers.
Use this when hosting the web interface on a domain name or proxying it through another IP address.
Or, these header checks can be disabled entirely with `disable-header-check`.

### http-prof

Enables go's http profiler interface, `pprof`.

Read more about `pprof`: https://github.com/skycoin/skycoin/wiki/Profiling-with-pprof

### http-prof-host

The interface address to bind the http profiler to.

### launch-browser

Open the web interface in the user's default browser.

### localhost-only

Bind the wire protocol `address` to localhost and only make connections to other localhost peers.

### log-level

Choose the log level verbosity.  Choices are: `debug`, `info`, `warn`, `error`, `fatal`, `panic`.

### logtofile

Write the log output to a file in `data-dir`. The logs will still be written to stdout.

### max-block-size

Maximum total size of transactions allowed when creating a new block. This value does not affect existing blocks.
The size of a transaction is the length of its byte representation in the [Skycoin binary encoding format](https://github.com/skycoin/skycoin/wiki/Skycoin-Binary-Encoding-Format).
Note that `max-block-size` is only the size limit of the transactions portion of a block; this limit does not include block metadata.
Only applies when running in `block-publisher` mode.

### max-connections

The maximum total number of connections to make over the wire protocol.

### max-decimals-create-block

The maximum number of decimal places applied to transactions when creating blocks.
Transactions that create outputs that exceed the decimal place limit will not be included in blocks.
Only applies when running in `block-publisher` mode.

### max-decimals-unconfirmed

The maximum number of decimal places applied to unconfirmed transactions received over the network.
Transactions that create outputs that exceed the decimal place limit will not be propagated to peers.

### max-default-peer-outgoing-connections

The maximum number of connections to maintain to peers from the hardcoded default peer list.
The peers in the hardcoded default peer list are maintained by Skycoin and are kept up to date with current
configurations. This value is 1 by default, to ensure at least one known stable connection is held.
More than 1 connections are not typically made, to avoid saturating the default peer connections.

### max-incoming-connections

The maximum number of incoming connections allowed.

### max-in-msg-len

Maximum length of incoming wire messages. Wire messages can include block and transaction data, so this limit should
be in accordance with `max-txn-size` and `max-block-size`.
If a peer sends a message that exceeds this limit, we disconnect from that peer.

### max-out-msg-len

Maximum length of outgoing wire messages. Wire messages can include block and transaction data, so this limit should
be in accordance with `max-txn-size` and `max-block-size`.
If we send a message that exceeds a peer's `max-in-msg-len`, they'll disconnect from us.

### max-outgoing-connections

The maximum total number of outgoing connections to make over the wire protocol.

### max-txn-size-create-block

The maximum transaction size applied to transactions when creating blocks.
The size of a transaction is the length of its byte representation in the [Skycoin binary encoding format](https://github.com/skycoin/skycoin/wiki/Skycoin-Binary-Encoding-Format).
Transactions that exceed this size will not be included in blocks.
Only applies when running in `block-publisher` mode.

### max-txn-size-unconfirmed

The maximum transaction size applied to unconfirmed transactions received over the network.
The size of a transaction is the length of its byte representation in the [Skycoin binary encoding format](https://github.com/skycoin/skycoin/wiki/Skycoin-Binary-Encoding-Format).
Transactions that exceed this size will not be propagated to peers.

### no-ping-log

Disable the "reply to ping" and "received pong" debug log messages.
These are particularly noisy, and unfortunately we only have one log level for debug,
so this option was added to disable them explicitly.

### peerlist-size

Maximum number of peers to track in the local peer database.

### peerlist-url

The URL of the remote peer list bootstrap file. Defaults to https://downloads.skycoin.com/blockchain/peers.txt.

### port

Port to bind for the wire protocol interface.

### profile-cpu

Enable the CPU profiler with `pprof`.

Read more about `pprof`: https://github.com/skycoin/skycoin/wiki/Profiling-with-pprof

### profile-cpu-file

Where to write the CPU profile data to, on exit.

### reset-corrupt-db

If the database is detected to be corrupted during startup, reset the database and continue running.
Otherwise, the application will abort if it detect a corrupted database.

The database is not always checked for corruption; it is only checked when upgrading the software and
if the upgraded version determines a corruption check is necessary.  However, if `verify-db` is enabled,
then the database is always checked for corruption.

### storage-dir

Location where the generic data storage files are saved. Defaults to a folder named `data` inside of the `data-dir`.

### user-agent-remark

An additional remark to include in the user agent that is sent in the introduction packet over the wire protocol

### verify-db

Unconditionally check the database for corruption on start.

The database is not always checked for corruption; it is only checked when upgrading the software and
if the upgraded version determines a corruption check is necessary.

### version

Print the node version and exit.

### wallet-crypto-type

Choose the encryption method for encrypted wallet data. Options are `sha256-xor` or `scrypt-chacha20poly1305`.
Do not use this option unless you know exactly what you are choosing; not every option provides meaningful encryption.

### wallet-dir

Location where the wallet files are saved. Defaults to a folder named `wallet` inside of the `data-dir`.

### web-interface

Enable the REST API interface. By default, it serves on http://127.0.0.1:6420.

### web-interface-addr

Address to bind the REST API interface to. Default `127.0.0.1`. Use `0.0.0.0` to bind to the machine's public IP interface.

### web-interface-cert

The certificate file for the HTTPS REST API. If not provided and HTTPS is enabled, the cert defaults to a file named `skycoind.cert`
in the `data-dir`. If this file does not exist, it will be autogenerated.

### web-interface-https

Use HTTPS for the REST API interface.

### web-interface-key

The key file for the HTTPS REST API. If not provided and HTTPS is enabled, the cert defaults to a file named `skycoind.key`
in the `data-dir`. If this file does not exist, it will be autogenerated.

### web-interface-password

Optional password for the REST API. Used in `Basic` authentication.

### web-interface-plaintext-auth

If this setting is not true, the application will not run if the REST API does not have HTTPS enabled and
is configured to use a username and/or password
(that is, if at least one of `web-interface-username` or `web-interface-password` are set).
This is to avoid sending authorization credentials in plaintext accidentally. The user must enable this setting explicitly to
send credentials in plain text.

### web-interface-port

Port number for the REST API interface. Default `6420`.

### web-interface-username

Optional username for the REST API. Used in `Basic` authentication.

## Development Environment Variables

These environment variables are for *development purposes only*. They are not intended
as part of the normal configuration API. That is why they are not included in the normal
configuration API (i.e. through cli `--options` or a configuration file).

### USER_BURN_FACTOR

The coin hour burn factor is the denominator in the ratio of coinhours that must be burned by a transaction.
For example, a burn factor of 2 means 1/2 of hours must be burned. A burn factor of 10 means 1/10 of coin hours must be burned.

The coin hour burn factor can be configured with a `USER_BURN_FACTOR` envvar. It cannot be configured through the command line.

```sh
$ USER_BURN_FACTOR=999 ./run-client.sh
```

This burn factor applies to user-created transactions.

To control the burn factor in other scenarios, use `burn-factor-unconfirmed` and `burn-factor-create-block`.

### USER_MAX_TXN_SIZE

```sh
$ USER_MAX_TXN_SIZE=1024 ./run-client.sh
```

This maximum transaction size applies to user-created transactions.

To control the transaction size in other scenarios, use `max-txn-size-unconfirmed` and `max-txn-size-create-block`.

To control the max block size, use `max-block-size`.

Transaction and block size are measured in bytes.

### USER_MAX_DECIMALS

```sh
$ USER_MAX_DECIMALS=4 ./run-client.sh
```

This maximum transaction size applies to user-created transactions.

To control the maximum decimals in other scenarios, use `max-decimals-unconfirmed` and `max-decimals-create-block`.
