# CLI

## Environment Setting

The CLI uses environment variable to manage the configuration.

## RPC_ADDRESS

CLI will connect to skycoin node rpc address:`127.0.0.1:6422` by default,
you can change the address by setting the `RPC_ADDRESS` env variable
with the following command:

```bash
$ export RPC_ADDRESS=127.0.0.1:6422
```

## WALLET_DIR

The default CLI wallet dir is located in `$HOME/.skycoin/wallets/`, change it by setting the 
`WALLET_DIR` environment variable.

```bash
$ export WALLET_DIR=$HOME/YOUR_WALLET_DIR
```

## WALLET_NAME

The default CLI wallet file name is `skycoin_cli.wlt`, chaing it by setting the `WALLET_NAME` env.
The wallet file name must have `.wlt` extension.

```bash
$ export WALLET_NAME=YOUR_WALLET_NAME
```


