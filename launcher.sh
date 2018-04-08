#!/bin/sh
COMMAND="skycoin --data-dir $DATA_DIR --wallet-dir $WALLET_DIR $@"
VOLUME_UID=`stat -c %U $DATA_DIR`
if [[ $VOLUME_UID = "UNKNOWN" ]] ; then
    adduser -D -u $VOLUME_UID skycoin
    su skycoin $COMMAND
else
    su $USERNAME $COMMAND
fi
