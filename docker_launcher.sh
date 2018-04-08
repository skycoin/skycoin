#!/bin/sh
COMMAND="skycoin --data-dir $DATA_DIR --wallet-dir $WALLET_DIR $@"
VOLUME_USERNAME=`stat -c %U $DATA_DIR`
VOLUME_UID=`stat -c %u $DATA_DIR`
if [[ $VOLUME_USERNAME = "UNKNOWN" ]] ; then
    adduser -D -u $VOLUME_UID skycoin
    su skycoin $COMMAND
else
    su $VOLUME_USERNAME $COMMAND
fi
