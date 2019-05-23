# Skycoin Daemon CLI Options

<!-- MarkdownTOC levels="1,2,3,4,5" autolink="true" bracket="round" -->

- [Preface](#preface)
- [Scenarios](#scenarios)
	- [Running a public API node](#running-a-public-api-node)
	- [Control which peers the node connects to](#control-which-peers-the-node-connects-to)
	- [Add Basic auth to the REST API interface](#add-basic-auth-to-the-rest-api-interface)
- [Options](#options)
	- [data-dir](#data-dir)
	- [web-interface](#web-interface)
	- [web-interface-addr](#web-interface-addr)
	- [web-interface-port](#web-interface-port)
	- [web-interface-https](#web-interface-https)
	- [web-interface-cert](#web-interface-cert)
	- [web-interface-key](#web-interface-key)
	- [web-interface-plaintext-auth](#web-interface-plaintext-auth)
	- [web-interface-username](#web-interface-username)
	- [web-interface-password](#web-interface-password)
- [Environment Variables](#environment-variables)
	- [USER_BURN_FACTOR](#userburnfactor)
	- [USER_MAX_TXN_SIZE](#usermax_txnsize)
	- [USER_MAX_DECIMALS](#usermaxdecimals)

<!-- /MarkdownTOC -->


## Preface

```
» go run cmd/skycoin/skycoin.go --help

Usage:
  -address string
    	IP Address to run application on. Leave empty to default to a public interface
  -arbitrating
    	Run node in arbitrating mode
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
  -max-txn-size-create-block uint
    	maximum size of a transaction applied when creating blocks (default 32768)
  -max-txn-size-unconfirmed uint
    	maximum size of an unconfirmed transaction (default 32768)
  -no-ping-log
    	disable "reply to ping" and "received pong" debug log messages
  -peerlist-size int
    	Max number of peers to track in peerlist (default 65535)
  -peerlist-url string
    	with -download-peerlist=true, download a peers.txt file from this url (default "https://downloads.skycoin.net/blockchain/peers.txt")
  -port int
    	Port to run application on (default 6000)
  -print-web-interface-address
    	print configured web interface address and exit
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
    	skycoind.cert file for web interface HTTPS. If not provided, will autogenerate or use skycoind.cert in -data-directory
  -web-interface-https
    	enable HTTPS for web interface
  -web-interface-key string
    	skycoind.key file for web interface HTTPS. If not provided, will autogenerate or use skycoind.key in -data-directory
  -web-interface-password string
    	password for the web interface
  -web-interface-plaintext-auth
    	allow web interface auth without https
  -web-interface-port int
    	port to serve web interface on (default 6420)
  -web-interface-username string
    	username for the web interface
Additional environment variables:
* USER_BURN_FACTOR - Set the coin hour burn factor required for user-created transactions. Must be >= 2.
* USER_MAX_TXN_SIZE - Set the maximum transaction size (in bytes) allowed for user-created transactions. Must be >= 1024.
* USER_MAX_DECIMALS - Set the maximum decimals allowed for user-created transactions. Must be <= 6.
```

## Scenarios

### Running a public API node

TODO

### Control which peers the node connects to

TODO

### Add Basic auth to the REST API interface

## Options

### data-dir

The storage location for application data. By default, the database, wallets, peers cache and other data files
will be saved into this folder.

On Linux and MacOS, this folder defaults to `$HOME/.skycoin` (`~/.skycoin`).
On Windows release builds, this folder defaults to `%HOMEPATH%\.skycoin` (`C:\Users\{user}\.skycoin`).
On Windows development builds, this folder defaults to `C:\.skycoin`. *(Note: this is a bug and will change in the future)*

### web-interface

Enable the REST API interface. By default, it serves on http://127.0.0.1:6420.

### web-interface-addr

Address to bind the REST API interface to. Default `127.0.0.1`. Use `0.0.0.0` to bind to the machine's public IP interface.

### web-interface-port

Port number for the REST API interface. Default `6420`.

### web-interface-https

Use HTTPS for the REST API interface.

### web-interface-cert

The certificate file for the HTTPS REST API. If not provided and HTTPS is enabled, the cert defaults to a file named `skycoind.cert`
in the `--data-directory`.  If this file does not exist, it will be autogenerated.

### web-interface-key

The key file for the HTTPS REST API. If not provided and HTTPS is enabled, the cert defaults to a file named `skycoind.key`
in the `--data-directory`.  If this file does not exist, it will be autogenerated.

### web-interface-plaintext-auth

If this setting is not true, the application will not run if the REST API does not have HTTPS enabled and
is configured to use a username and/or password
(that is, if at least one of `--web-interface-username` or `--web-interface-password` are set).
This is to avoid sending authorization credentials in plaintext accidentally. The user must enable this setting explicitly to
send credentials in plain text.

### web-interface-username

Optional username for the REST API. Used in `Basic` authentication.

### web-interface-password

Optional password for the REST API. Used in `Basic` authentication.

## Environment Variables

### USER_BURN_FACTOR

The coin hour burn factor is the denominator in the ratio of coinhours that must be burned by a transaction.
For example, a burn factor of 2 means 1/2 of hours must be burned. A burn factor of 10 means 1/10 of coin hours must be burned.

The coin hour burn factor can be configured with a `USER_BURN_FACTOR` envvar. It cannot be configured through the command line.

```sh
$ USER_BURN_FACTOR=999 ./run-client.sh
```

This burn factor applies to user-created transactions.

To control the burn factor in other scenarios, use `-burn-factor-unconfirmed` and `-burn-factor-create-block`.

### USER_MAX_TXN_SIZE

```sh
$ USER_MAX_TXN_SIZE=1024 ./run-client.sh
```

This maximum transaction size applies to user-created transactions.

To control the transaction size in other scenarios, use `-max-txn-size-unconfirmed` and `-max-txn-size-create-block`.

To control the max block size, use `-max-block-size`.

Transaction and block size are measured in bytes.

### USER_MAX_DECIMALS

```sh
$ USER_MAX_DECIMALS=4 ./run-client.sh
```

This maximum transaction size applies to user-created transactions.

To control the maximum decimals in other scenarios, use `-max-decimals-unconfirmed` and `-max-decimals-create-block`.
