# CLI

## Install

```bash
$ cd cmd/cli
$ go install
```

After `go install`, the cli command will be available in `$GOPATH/bin` folder, and if you
have added the `GOPATH/bin` to `PATH` env, you can run `cli` command directly.

## Environment Setting

The CLI uses environment variable to manage the configurations.

## RPC_ADDR

CLI will connect to skycoin node rpc address:`127.0.0.1:6430` by default,
you can change the address by setting the `RPC_ADDRESS` env variable
with the following command:

```bash
$ export RPC_ADDR=127.0.0.1:6430
```

## WALLET_DIR

The default CLI wallet dir is located in `$HOME/.skycoin/wallets/`, change it by setting the 
`WALLET_DIR` environment variable.

```bash
$ export WALLET_DIR=$HOME/YOUR_WALLET_DIR
```

## WALLET_NAME

The default CLI wallet file name is `skycoin_cli.wlt`, change it by setting the `WALLET_NAME` env.
The wallet file name must have `.wlt` extension.

```bash
$ export WALLET_NAME=YOUR_WALLET_NAME
```

## Send coin

```bash
$ cli send $recipient_address $amount
```

The above `send` command will send coins from your node's default wallet: `$HOME/.skycoin/wallets/skycoin_cli.wlt`. You can also send from the wallet
as you want, just use the `-f` option flag, example:

```bash
$ cli send -f $WALLET_PATH $recipient_address $amount
```

Use `cli send -h` to see more details about the send command.

## Check address balance

```bash
$ cli addressBalance $addr1 $addr2
```

See more details with command `cli -h`.