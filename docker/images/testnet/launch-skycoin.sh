#!/bin/sh
mkdir -p $DATA_DIR$COIN
nslookup skycoin-peer | awk '/Address/{print $3":16000"}' > $DATA_DIR$COIN/connections.txt
nslookup skycoin-master | awk '/Address/{print $3":16000"}' >> $DATA_DIR$COIN/connections.txt
$@
