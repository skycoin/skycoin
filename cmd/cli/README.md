# CLI

Skycoin command line interface

## Install

```bash
$ cd $GOPATH/src/github.com/skycoin/skycoin/cmd/cli
$ ./install.sh
```

### Enable command autocomplete

If you are in `bash`, run the following command:

```bash
$ PROG=skycoin-cli source $GOPATH/src/github.com/skycoin/skycoin/cmd/cli/autocomplete/bash_autocomplete
```

If you are in `zsh`, please replace the `bash_autocomplete` with `zsh_autocomplete` in the previous command.

To avoid run the command everytime when you start a new terminal session, you can copy the script into
you `~/.bashrc` or `~/.zshrc` file.

## Environment Setting

The CLI uses environment variable to manage the configurations.

### RPC_ADDR

CLI will connect to skycoin node rpc address:`127.0.0.1:6430` by default,
you can change the address by setting the `RPC_ADDR` env variable
with the following command:

```bash
$ export RPC_ADDR=127.0.0.1:6430
```

### WALLET_DIR

The default CLI wallet dir is located in `$HOME/.skycoin/wallets/`, change it by setting the
`WALLET_DIR` environment variable.

```bash
$ export WALLET_DIR=$HOME/YOUR_WALLET_DIR
```

### WALLET_NAME

The default CLI wallet file name is `skycoin_cli.wlt`, change it by setting the `WALLET_NAME` env.
The wallet file name must have `.wlt` extension.

```bash
$ export WALLET_NAME=YOUR_WALLET_NAME
```

## Usage

After the installation, you can run `skycoin-cli` to see the usage:

```bash

$ skycoin-cli

NAME:
   skycoin-cli - the skycoin command line interface

USAGE:
   skycoin-cli [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
     addPrivateKey         Add a private key to specific wallet
     blocks                Lists the content of a single block or a range of blocks
     broadcastTransaction  Broadcast a raw transaction to the network
     walletBalance         Check the balance of a wallet
     walletOutputs         Display outputs of specific wallet
     addressBalance        Check the balance of specific addresses
     addressOutputs        Display outputs of specific addresses
     createRawTransaction  Create a raw transaction to be broadcast to the network later
     generateAddresses     Generate additional addresses for a wallet
     generateWallet        Generate a new wallet
     lastBlocks            Displays the content of the most recently N generated blocks
     listAddresses         Lists all addresses in a given wallet
     listWallets           Lists all wallets stored in the default wallet directory
     send                  Send skycoin from a wallet or an address to a recipient address
     status                Check the status of current skycoin node
     transaction           Show detail info of specific transaction
     version
     walletDir             Displays wallet folder address
     walletHistory         Display the transaction history of specific wallet
     help, h               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help, can also be used to show subcommand help
   --version, -v  print the version
```

### Send coin

```bash
$ skycoin-cli send $recipient_address $amount
```

The above `send` command will send coins from your node's default wallet: `$HOME/.skycoin/wallets/skycoin_cli.wlt`. You can also send from the wallet
as you want, just use the `-f` option flag, example:

```bash
$ skycoin-cli send -f $WALLET_PATH $recipient_address $amount
```

Use `skycoin-cli send -h` to see the subcommand usage.

### Check address balance

```bash
$ skycoin-cli addressBalance 2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc
{
    "total_amount": 1,
    "addresses": [
        {
            "address": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc",
            "amount": 1
        }
    ]
}
```

### Get transaction

```bash
$ skycoin-cli transaction 824d421a25f81aa7565d042a54b3e1e8fdc58bed4eefe8f8a90748da6d77d135
{
    "transaction": {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 1,
            "block_seq": 864,
            "unknown": false
        },
        "txn": {
            "length": 220,
            "type": 0,
            "txid": "824d421a25f81aa7565d042a54b3e1e8fdc58bed4eefe8f8a90748da6d77d135",
            "inner_hash": "708a21c685041ba409b9634843003f263c7a00d99459925e319049f3e36f1163",
            "timestamp": 1492141347,
            "sigs": [
                "6772c4e1370faf7acd8fc37f6dba3bea06ad1beb1300074c33b2c6fa9b11ed7d2bcc01b7008f235aa918b9c146078dfb8e1c8ce55b0893ea5f111597b42331ba01"
            ],
            "inputs": [
                "c38c108ac3c76e5faffce0bb83153ec98bc1355a98e1a9b0f95ab1b98ef9f00e"
            ],
            "outputs": [
                {
                    "uxid": "b0586a8e731c475e87eb61ef0b845d7893cf39120a1e97cf05f78585f1a49e3c",
                    "dst": "2bfYafFtdkCRNcCyuDvsATV66GvBR9xfvjy",
                    "coins": "166365",
                    "hours": 0
                },
                {
                    "uxid": "49c64719d8df905a7fd4a1c46c2a9c358a8c5ae14befb0d57dcfe1b1d36a1758",
                    "dst": "ep3axwpJ3hWWQcACu48z9sMKUB7snXBm94",
                    "coins": "1300",
                    "hours": 0
                }
            ]
        }
    }
}
```

## Note

The `[option]` in subcommand must be set before the rest values, otherwise the `option` won't
be parsed, example:

If we want to specify a `change address` in `send` command, we can use `-c` option, if you run
the command in the following way:

```bash
$ skycoin-cli send $recipient_address $amount -c $change_address
```

The change coins won't go to the address as you wish, it will go to the
default `change address`, which can be by `from` address or the wallet's
coinbase address.

The right script should look like this:

```bash
$ skycoin-cli send -c $change_address $recipient_address $amount
```
