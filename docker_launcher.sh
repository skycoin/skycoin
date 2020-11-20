#!/bin/sh
COMMAND="skycoin --data-dir $DATA_DIR --wallet-dir $WALLET_DIR $@"

adduser -D -u 10000 skycoin

if [[ \! -d $DATA_DIR ]]; then
    mkdir -p $DATA_DIR
fi
if [[ \! -d $WALLET_DIR ]]; then
    mkdir -p $WALLET_DIR
fi

chown -R skycoin:skycoin $( realpath $DATA_DIR )
chown -R skycoin:skycoin $( realpath $WALLET_DIR )

su skycoin -c "$COMMAND"
