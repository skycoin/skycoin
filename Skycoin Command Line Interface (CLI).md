# Skycoin Command Line Interface (CLI)


### Purpose

The purpose of this README is to familiarize the user with installing and operating the Skycoin wallet from the command line for OSX. 

------

### Background

CLI is a command line client for interacting with a skycoin node and for offline wallet management. Using command line is an excellent option for users who can wish to develop their own user interface experience. 

------

### Installation 

####Language

Go must first be installed using ```homebrew```. 

Run the following code into the terminal to install ```homebrew```:

```/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"```

The latest version of Go can now be installed by running the following code:

```brew install go```

Mercurial and Bazaar should also be installed. Run the following code: 

```brew install mercurial bzr```

#### Enable command autocomplete

Run the following code into the terminal : 

```
$ PROG=skycoin-cli source $GOPATH/src/github.com/skycoin/skycoin/cmd/cli/autocomplete/bash_autocomplete
```

#### CLI version

Check for the current version of Skycoin CLI.

```
$ skycoin-cli version [flags]
```

```
FLAGS:
  -j, --json   Returns the results in JSON format
```



****



### Environment Setting

The following  information needs to be entered into the terminal to identify and access the user account. 

**Address**

CLI will connect to the Skycoin node REST API address `http://127.0.0.1:6420` by default. If desired, the address can be changed by setting the `RPC_ADDR` environment variable with the following command:

```
$ export RPC_ADDR=http://127.0.0.1:6420
```

*Note: `RPC_ADDR` must be in `scheme://host` format.*

**Username**

A username for authenticating requests to the Skycoin node.

```
$ export RPC_USER=...
```

**Password**

A password for authenticating requests to the Skycoin node.

```
$ export RPC_PASS=...
```

**Wallet Directory**

The default CLI wallet directory is located in `$HOME/.skycoin/wallets/`. Change it by setting the `WALLET_DIR` environment variable.

```
$ export WALLET_DIR=$HOME/YOUR_WALLET_DIR
```

**Wallet Name**

The default CLI wallet file name is `skycoin_cli.wlt`. Change it by setting the `WALLET_NAME` env. The wallet file name must have `.wlt` extension.

```
$ export WALLET_NAME=YOUR_WALLET_NAME
```

------



### Usage

After the installation, run `skycoin-cli` to see the use options:

```
$ skycoin-cli

USAGE:
  skycoin-cli [command] [flags] [arguments...]

DESCRIPTION:
    The skycoin command line interface

COMMANDS:
  addPrivateKey        Add a private key to specific wallet
  addressBalance       Check the balance of specific addresses
  addressGen           Generate skycoin or bitcoin addresses
  addressOutputs       Display outputs of specific addresses
  addressTransactions  Show detail for transaction associated with one or more specified 													 addresses
  blocks               Lists the content of a single block or a range of blocks
  broadcastTransaction Broadcast a raw transaction to the network
  checkdb              Verify the database
  createRawTransaction Create a raw transaction to be broadcast to the network later
  decodeRawTransaction Decode raw transaction
  decryptWallet        Decrypt wallet
  encryptWallet        Encrypt wallet
  fiberAddressGen      Generate addresses and seeds for a new fiber coin
  help                 Help about any command
  lastBlocks           Displays the content of the most recently N generated blocks
  listAddresses        Lists all addresses in a given wallet
  listWallets          Lists all wallets stored in the wallet directory
  richlist             Get skycoin richlist
  send                 Send skycoin from a wallet or an address to a recipient address
  showConfig           Show cli configuration
  showSeed             Show wallet seed
  status               Check the status of current skycoin node
  transaction          Show detail info of specific transaction
  verifyAddress        Verify a skycoin address
  version              List the current version of Skycoin components
  walletAddAddresses   Generate additional addresses for a wallet
  walletBalance        Check the balance of a wallet
  walletCreate         Generate a new wallet
  walletDir            Displays wallet folder address
  walletHistory        Display the transaction history of specific wallet. Requires skycoin node rpc.
  walletOutputs        Display outputs of specific wallet

FLAGS:
  -h, --help      help for skycoin-cli
      --version   version for skycoin-cli

Use "skycoin-cli [command] --help" for more information about a command.

ENVIRONMENT VARIABLES:
    RPC_ADDR: Address of RPC node. Must be in scheme://host format. Default "http://127.0.0.1:6420"
    RPC_USER: Username for RPC API, if enabled in the RPC.
    RPC_PASS: Password for RPC API, if enabled in the RPC.
    COIN: Name of the coin. Default "skycoin"
    WALLET_DIR: Directory where wallets are stored. This value is overridden by any subcommand flag specifying a wallet filename, if that filename includes a path. Default "$DATA_DIR/wallets"
    WALLET_NAME: Name of wallet file (without path). This value is overridden by any 	subcommand flag specifying a wallet filename. Default "$COIN_cli.wlt"
    DATA_DIR: Directory where everything is stored. Default "$HOME/.$COIN/"
```

------

### Wallet Read and Update Commands

These are the essential commands to read and update the current wallet information. 

#### Create a wallet

To create a new skycoin wallet run the following code:

```
$ skycoin-cli walletCreate [flags]
```

```
FLAGS:
  -x, --crypto-type string   The crypto type for wallet encryption, can be scrypt-chacha20poly1305 or sha256-xor (default "scrypt-chacha20poly1305")
  
  -e, --encrypt              Create encrypted wallet.
  -l, --label string         Label used to idetify your wallet.
  -m, --mnemonic             A mnemonic seed consisting of 12 dictionary words will be generated
  
  -n, --num uint             [numberOfAddresses] Number of addresses to generate
                                 By default 1 address is generated. (default 1)
                                 
  -p, --password string      Wallet password
  -r, --random               A random alpha numeric seed will be generated
  -s, --seed string          Your seed
  -f, --wallet-file string   Name of wallet. The final format will be "yourName.wlt".
        
```


  *Note: If no wallet name is specified a generic name will be selected. (default"skycoin_cli.wlt")*

#### List wallets

List wallets in the Skycoin wallet directory.

```
$ skycoin-cli listWallets
```

#### Richlist

Returns top N address (default 20) balances (based on unspent outputs). Optionally include distribution addresses (exluded by default).

```
$ skycoin-cli richlist [top N addresses (20 default)] [include distribution addresses (false default)]
FLAGS:
  -h, --help   help for richlist
```

##### Example

**Without distribution addresses**

```
$ skycoin-cli richlist 2
```

<details open="" style="box-sizing: border-box; display: block; margin-bottom: 16px; margin-top: 0px; color: rgb(36, 41, 46); font-family: -apple-system, system-ui, &quot;Segoe UI&quot;, Helvetica, Arial, sans-serif, &quot;Apple Color Emoji&quot;, &quot;Segoe UI Emoji&quot;, &quot;Segoe UI Symbol&quot;; font-size: 16px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; white-space: normal; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; background-color: rgb(255, 255, 255); text-decoration-style: initial; text-decoration-color: initial;"><summary style="box-sizing: border-box; display: list-item; cursor: pointer;">View Output</summary><div class="highlight highlight-source-json" style="box-sizing: border-box; margin-bottom: 16px;"><pre style="box-sizing: border-box; font-family: SFMono-Regular, Consolas, &quot;Liberation Mono&quot;, Menlo, Courier, monospace; font-size: 13.6px; margin-bottom: 0px; margin-top: 0px; overflow-wrap: normal; background-color: rgb(246, 248, 250); border-radius: 3px; line-height: 1.45; overflow: auto; padding: 16px; word-break: normal;">{
    <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>richlist<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: [
        {
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>address<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>zVzkqNj3Ueuzo54sbACcYBqqGBPCGAac5W<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>coins<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>2922927.299000<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>locked<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">false</span>
        },
        {
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>address<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>coins<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>675256.308000<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>locked<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">false</span>
        }
    ]
}</pre></div></details>

**Including distribution addresses**

```
$ skycoin-cli richlist 2 true
```

<details open="" style="box-sizing: border-box; display: block; margin-bottom: 16px; margin-top: 0px; color: rgb(36, 41, 46); font-family: -apple-system, system-ui, &quot;Segoe UI&quot;, Helvetica, Arial, sans-serif, &quot;Apple Color Emoji&quot;, &quot;Segoe UI Emoji&quot;, &quot;Segoe UI Symbol&quot;; font-size: 16px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; white-space: normal; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; background-color: rgb(255, 255, 255); text-decoration-style: initial; text-decoration-color: initial;"><summary style="box-sizing: border-box; display: list-item; cursor: pointer;">View Output</summary><div class="highlight highlight-source-json" style="box-sizing: border-box; margin-bottom: 16px;"><pre style="box-sizing: border-box; font-family: SFMono-Regular, Consolas, &quot;Liberation Mono&quot;, Menlo, Courier, monospace; font-size: 13.6px; margin-bottom: 0px; margin-top: 0px; overflow-wrap: normal; background-color: rgb(246, 248, 250); border-radius: 3px; line-height: 1.45; overflow: auto; padding: 16px; word-break: normal;">{
    <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>richlist<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: [
        {
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>address<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>zVzkqNj3Ueuzo54sbACcYBqqGBPCGAac5W<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>coins<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>2922927.299000<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>locked<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">false</span>
        },
        {
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>address<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>ejJjiCwp86ykmFr5iTJ8LxQXJ2wJPTYmkm<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>coins<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>1000000.010000<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
            <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>locked<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">true</span>
        }
    ]
}</pre></div></details>

#### Add addresses to a wallet

Add new addresses to a Skycoin wallet.

```
$ skycoin-cli walletAddAddresses [flags]
```

```
FLAGS:
  -j, --json                 Returns the results in JSON format
  -n, --num uint             Number of addresses to generate (default 1)
  -p, --password string      wallet password
  -f, --wallet-file string   Generate addresses in the wallet (default 		"$HOME/.skycoin/wallets/skycoin_cli.wlt")

```

#### List wallet addresses

List addresses in a Skycoin wallet.

```
$ skycoin-cli listAddresses [walletName]
```

#### Address Count

Returns the count of all addresses that currenty have unspent outputs (coins) associated with them.

```
$ skycoin-cli addresscount
```

```
FLAGS:
  -h, --help   help for richlist
```

#### Check wallet balance

Check the balance of the Skycoin wallet.

```
$ skycoin-cli walletBalance [wallet]
```

*NOTE: Both the full wallet path or only the wallet name can be used. If no wallet is specified then the default wallet: `$HOME/.$COIN/wallets/skycoin_cli.wlt` is used.*

#### See wallet directory

Get the current Skycoin wallet directory.

```
$ skycoin-cli walletDir [flags]
```

```
FLAGS:
        -j, --json  Returns the results in JSON format.
```

#### List wallet outputs

List unspent outputs of all addresses in a wallet.

```
$ skycoin-cli walletOutputs [wallet file]
```

##### Examples

**Default wallet**

```
$ skycoin-cli walletOutputs
```

##### Specific wallet

```
$ skycoin-cli walletHistory $WALLET_NAME
```

*OR*

```
$ skycoin-cli walletHistory $WALLET_PATH
```

<details open="" style="box-sizing: border-box; display: block; margin-bottom: 16px; margin-top: 0px; color: rgb(36, 41, 46); font-family: -apple-system, system-ui, &quot;Segoe UI&quot;, Helvetica, Arial, sans-serif, &quot;Apple Color Emoji&quot;, &quot;Segoe UI Emoji&quot;, &quot;Segoe UI Symbol&quot;; font-size: 16px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; white-space: normal; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; background-color: rgb(255, 255, 255); text-decoration-style: initial; text-decoration-color: initial;"><summary style="box-sizing: border-box; display: list-item; cursor: pointer;">View Output</summary><div id="details-content"><slot name="user-agent-default-slot"><div class="highlight highlight-source-json" style="box-sizing: border-box; margin-bottom: 16px;"><pre style="box-sizing: border-box; font-family: SFMono-Regular, Consolas, &quot;Liberation Mono&quot;, Menlo, Courier, monospace; font-size: 13.6px; margin-bottom: 0px; margin-top: 0px; overflow-wrap: normal; background-color: rgb(246, 248, 250); border-radius: 3px; line-height: 1.45; overflow: auto; padding: 16px; word-break: normal;">{
 <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>outputs<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: {
     <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>head_outputs<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: [
         {
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>hash<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>c51b2692aa9f296a3cd2f37b14f39c496c82f5c5ae01c54701ea60b7353f27e2<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>time<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">1523184376</span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>block_seq<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">21221</span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>src_tx<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>f3c5cfd462d95e724b7d35b1688c53f25a5f358f2eb9a6f87b63cdf31deb2bf8<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>address<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>coins<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>15.000000<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>hours<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">369</span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>calculated_hours<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">370</span>
         },
         {
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>hash<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>a0777af14223bbbd5aeb8bf3cfd6ba94c776c6eec731310caaaaee49b9feb9a5<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>time<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">1523184176</span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>block_seq<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">21220</span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>src_tx<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>4acd61d7aa7dfe20795e517d7560643d049036af9451bcbd762793bcb6a4a6de<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>address<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>coins<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>1.000000<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>hours<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">0</span>,
             <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>calculated_hours<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: <span class="pl-c1" style="box-sizing: border-box; color: rgb(0, 92, 197);">0</span>
         }
     ],
     <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>outgoing_outputs<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: [],
     <span class="pl-s" style="box-sizing: border-box; color: rgb(3, 47, 98);"><span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span>incoming_outputs<span class="pl-pds" style="box-sizing: border-box; color: rgb(3, 47, 98);">"</span></span>: []
 }
}</pre></div></slot></div></details>

#### Add Private Key

Add a private key to a Skycoin wallet.

```
$ skycoin-cli addPrivateKey [flags] [private key]
FLAGS:
  -h, --help                 help for addPrivateKey
  -p, --password string      Wallet password
  -f, --wallet-file string   wallet file or path. If no path is specified your default wallet path will be used.
```

##### Example

```
$ skycoin-cli addPrivateKey -f $WALLET_PATH $PRIVATE_KEY
$ success
```



#### Generate new addresses

Generate new Skycoin or bitcoin addresses.

```
$ skycoin-cli addressGen [flags]
FLAGS:
  -c, --coin string    Coin type. Must be skycoin or bitcoin. If bitcoin, secret keys are in Wallet Import Format instead of hex. (default "skycoin")
  -x, --encrypt        Encrypt the wallet when printing a JSON wallet
  -e, --entropy int    Entropy of the autogenerated bip39 seed, when the seed is not provided. Can be 128 or 256 (default 128)
      --hex            Use hex(sha256sum(rand(1024))) (CSPRNG-generated) as the seed if not seed is not provided
  -i, --hide-secrets   Hide the secret key and seed from the output when printing a JSON wallet file
  -l, --label string   Wallet label to use when printing or writing a wallet file
  -m, --mode string    Output mode. Options are wallet (prints a full JSON wallet), addresses (prints addresses in plain text), secrets (prints secret keys in plain text) (default "wallet")
  -n, --num int        Number of addresses to generate (default 1)
  -s, --seed string    Seed for deterministic key generation. Will use bip39 as the seed if not provided.
  -t, --strict-seed    Seed should be a valid bip39 mnemonic seed.
```

##### Examples

**Generate `n` number of skycoin addresses**

```
$ skycoin-cli addressGen --num 2
```


**Generate `n` number of bitcoin addresses**

```
$ skycoin-cli addressGen --num 2 --coin bitcoin
```


**Hide secret in output**

```
$ skycoin-cli addressGen --num 2 --hide-secrets
```



**Output only an address list**

```
$ skycoin-cli addressGen --num 2 --mode addresses
7HVmKni3ggMdtseynSkNkqoCnsH7vkS6cg
2j5QSbHgLWXA2qXZvLzJHRo6Cissxer4CSt
```

*Note: If no seed is provided with the `--seed flag` and `--hex` flag is not used then bip39 is used to generate a seed*

**Use a predefined seed value**

```
$ skycoin-cli addressGen --num 2 --seed "my super secret seed"
```



**Generate addresses with a hex (CSPRNG-generated) seed**

```
skycoin-cli addressGen --num 2 --hex
```


#### Generate distribution addresses for a new fiber coin

```
skycoin-cli fiberAddressGen [flags]
DESCRIPTION:
    Addresses are written in a format that can be copied into fiber.toml
    for configuring distribution addresses. Addresses along with their seeds are written to a csv file,
    these seeds can be imported into the wallet to access distribution coins.

FLAGS:
  -a, --addres-file string   Output file for the generated addresses in fiber.toml format (default "addresses.txt")
  -e, --entropy int          Entropy of the autogenerated bip39 seeds. Can be 128 or 256 (default 128)
  -n, --num int              Number of addresses to generate (default 100)
  -o, --overwrite            Allow overwriting any existing addrs-file or seeds-file
  -s, --seeds-file string    Output file for the generated addresses and seeds in a csv (default "seeds.csv")
```

##### Example

```
skycoin-cli fiberAddressGen
```

#### Check address outputs

Display outputs of specific addresses, join multiple addresses with space.

```
$ skycoin-cli addressOutputs [address list]
```

##### Example

```
skycoin-cli addressOutputs tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V 29fDBQuJs2MDLymJsjyWH6rDjsyv995SrGU
```


#### Check block data

Lists the content of a single block or a range of blocks

```
$ skycoin-cli blocks [starting block or single block seq] [ending block seq]
```

##### Example

```
$ skycoin-cli blocks 41 42
```


#### Check database integrity

Checks if the given database file contains valid skycoin blockchain data If no argument is given, the default `data.db` in `$HOME/.$COIN/` will be checked.

```
$ skycoin-cli checkdb [db path]
```

##### Example

```
$ skycoin-cli checkdb $DB_PATH
```


#### Create a raw transaction

Create a raw transaction that can be broadcasted later. A raw transaction is a binary encoded hex string.

```
$ skycoin-cli createRawTransaction [flags] [to address] [amount]
FLAGS:
  -a, --address string          From address
  -c, --change-address string   Specify different change address.
                                By default the from address or a wallets coinbase address will be used.
      --csv  string         CSV file containing addresses and amounts to send
  -j, --json                    Returns the results in JSON format.
  -m, --many string             use JSON string to set multiple receive addresses and coins,
                                example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'
  -p, --password string         Wallet password
  -f, --wallet-file string      wallet file or path. If no path is specified your default wallet path will be used.
```

##### Examples

**Sending to a single address from a specified wallet**

```
$ skycoin-cli createRawTransaction -f $WALLET_PATH -a $FROM_ADDRESS $RECIPIENT_ADDRESS $AMOUNT
```

**Sending to a specific change address**

```
$ skycoin-cli createRawTransaction -f $WALLET_PATH -a $FROM_ADDRESS -c $CHANGE_ADDRESS $RECIPIENT_ADDRESS $AMOUNT
```


**Sending to multiple addresses**

```
$ skycoin-cli createRawTransaction -f $WALLET_PATH -a $FROM_ADDRESS -m '[{"addr":"$ADDR1", "coins": "$AMT1"}, {"addr":"$ADDR2", "coins": "$AMT2"}]'
```

**Sending to addresses in a CSV file**

```
$ cat <<EOF > $CSV_FILE
2Niqzo12tZ9ioZq5vwPHMVR4g7UVpp9TCmP,123.1
2UDzBKnxZf4d9pdrBJAqbtoeH641RFLYKxd,456.045
yExu4fryscnahAEMKa7XV4Wc1mY188KvGw,0.3
EOF
$ skycoin-cli createRawTransaction -f $WALLET_PATH -a $FROM_ADDRESS -csv $CSV_FILE
```


*NOTE: When sending to multiple addresses each combination of address and coins need to be unique Otherwise you get, `ERROR: Duplicate output in transaction`*

**Generate a JSON output**

```
$ skycoin-cli createRawTransaction -f $WALLET_PATH -a $FROM_ADDRESS --json $RECIPIENT_ADDRESS $AMOUNT
```


#### Decode a raw transaction

Decode a raw skycoin transaction.

```
$ skycoin-cli decodeRawTransaction [raw transaction]
```

##### Example

```
skycoin-cli decodeRawTransaction dc00000000247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af384301000000cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f000100000005e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634020000000007445b5d6fbbb1a7d70bef941fb5da234a10fcae40420f00000000000100000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec090024f400000000009805000000000000
```

#### Broadcast a raw transaction

Broadcast a raw Skycoin transaction. Output is the transaction id.

```
$ skycoin-cli broadcastTransaction [raw transaction]
$ skycoin-cli broadcastTransaction dc00000000247bd0f0a1cf39fa51ea3eca044e4d9cbb28fff5376e90e2eb008c9fe0af384301000000cf5869cb1b21da4da98bdb5dca57b1fd5a6fcbefd37d4f1eb332b21233f92cd62e00d8e2f1c8545142eaeed8fada1158dd0e552d3be55f18dd60d7e85407ef4f000100000005e524872c838de517592c9a495d758b8ab2ec32d3e4d3fb131023a424386634020000000007445b5d6fbbb1a7d70bef941fb5da234a10fcae40420f00000000000100000000000000008001532c3a705e7e62bb0bb80630ecc21a87ec090024f400000000009805000000000000
```


------

### Send Skycoin

These are the commands to transfer Skycoin.

Make a Skycoin transaction.

```
$ skycoin-cli send [flags] [to address] [amount]
```

```
FLAGS:
  -a, --address string          From address
  -c, --change-address string   Specify different change address.
                                By default the from address or a wallets coinbase address 																	will be used.
 	--csv  string             CSV file containing addresses and amounts to send
  -j, --json                    Returns the results in JSON format.
  -m, --many string             use JSON string to set multiple receive addresses and coins,
                                example: -m '[{"addr":"$addr1", "coins": "10.2"},		{"addr":"$addr2", "coins": "20"}]'
                                
  -p, --password string         Wallet password
  -f, --wallet-file string      Wallet file or path. If no path is specified your default       																wallet path will be used.
```

#### Examples

##### Sending from the default wallet

```
$ skycoin-cli send $RECIPIENT_ADDRESS $AMOUNT
```

##### Sending from a specific wallet

```
$ skycoin-cli send -f $WALLET_PATH $RECIPIENT_ADDRESS $AMOUNT
```

##### Sending from a specific address in a wallet

```
$ skycoin-cli send -f $WALLET_PATH -a $FROM_ADDRRESS $RECIPIENT_ADDRESS $AMOUNT
```

*NOTE: If $WALLET_PATH is not specified above then the default wallet is used.*

##### Sending change to a specific change address

```
$ skycoin-cli send -f $WALLET_PATH -a $FROM_ADDRESS -c $CHANGE_ADDRESS $RECIPIENT_ADDRESS $AMOUNT
```

##### Sending to multiple addresses

```
$ skycoin-cli send -f $WALLET_PATH -a $FROM_ADDRESS -m '[{"addr":"$ADDR1", "coins": "$AMT1"}, {"addr":"$ADDR2", "coins": "$AMT2"}]'
```

##### Sending to addresses in a CSV file

```
$ cat <<EOF > $CSV_FILE
2Niqzo12tZ9ioZq5vwPHMVR4g7UVpp9TCmP,123.1
2UDzBKnxZf4d9pdrBJAqbtoeH641RFLYKxd,456.045
yExu4fryscnahAEMKa7XV4Wc1mY188KvGw,0.3
EOF
$ skycoin-cli send -f $WALLET_PATH -a $FROM_ADDRESS -csv $CSV_FILE
```

*NOTE: When sending to multiple addresses each combination of address and coins need to be unique Otherwise you get, `ERROR: Duplicate output in transaction*`

##### Generate a JSON output

```
$ skycoin-cli send -f $WALLET_PATH -a $FROM_ADDRESS --json $RECIPIENT_ADDRESS $AMOUNT
```


### Other Considerations

The `[option]` in subcommand must be set before the rest of the values, otherwise the `option` won't be parsed. For example:

To specify a `change address` in `send` command, use the `-c` option, to run the command in the following way:

```
$ skycoin-cli send $RECIPIENT_ADDRESS $AMOUNT -c $CHANGE_ADDRESS
```

The change coins will not go to the address as intended. It will go to the default `change address`, which can be by the `from`address or the wallet's coinbase address.

The correct script should look like this:

```
$ skycoin-cli send -c $CHANGE_ADDRESS $RECIPIENT_ADDRESS $AMOUNT
```

