#!/bin/sh
COMMAND="skycoin --data-dir $DATA_DIR --wallet-dir $WALLET_DIR $@"
VOLUME_UID=`stat -c %u $DATA_DIR`
USERNAME=`id -nu $VOLUME_UID`
if [[ $? -ne 0 ]] ; then
    adduser --uid $VOLUME_UID skycoin
    su skycoin $COMMAND
else
    su $USERNAME $COMMAND
fi
