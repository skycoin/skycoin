# Skycoin CLI

## Environment Setting

The skycoin CLI uses environment variable to manage the settings.

## SKYCOIN_NODE_ADDR

CLI will connect to skycoin node of address: `127.0.0.1:6420` by default,
you can change the address by setting the `SKYCOIN_NODE_ADDR` env variable
with the following command:

```bash
$ export SKYCOIN_NODE_ADDR=127.0.0.1:6421
```

## SKYCOIN_WLT_DIR

The default skycoin CLI wallet dir is located in `$HOME/.skycoin-cli/wallet/`, change it by setting the 
`SKYCOIN_WLT_DIR` environment variable.

```bash
$ export SKYCOIN_WLT_DIR=$HOME/YOUR_WALLET_DIR
```
