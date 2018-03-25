#!/bin/sh
MY_IP=`ip addr show | grep global | awk '{print $2}' | cut -d\/ -f 1`
mkdir -p $DATA_DIR$COIN
nslookup skycoin-peer | awk '/Address/{print $3":16000"}' > $DATA_DIR$COIN/connections.txt
nslookup skycoin-master | awk '/Address/{print $3":16000"}' >> $DATA_DIR$COIN/connections.txt
sed -i '/$MY_IP/d' $DATA_DIR$COIN/connections.txt
$@
