#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

ADDRV=$1
if [ -z "$ADDRV" ]; then
    ADDRV=test
fi

function check_status {
    if [ $? -ne 0 ]; then
        popd >/dev/null
        popd >/dev/null
        exit 1
    fi
}

pushd "$DIR" >/dev/null
pushd .. >/dev/null
echo "Generating new master keypair"
go run cmd/address/address.go -o master.keys -address-version="$ADDRV"
check_status
echo "Creating genesis block"
go run cmd/genesis/genesis.go -keys master.keys -help=false -address-version="$ADDRV"
check_status
echo "====================="
echo "NOTICE: Make sure to copy the blockchain files and master keys file to the -data-directory used by the master node"
echo "====================="
echo "CLI flags for non-master client:"
echo "====================="
echo -n "-master-public-key="
go run cmd/address/address.go -i master.keys -print-public=true
check_status
echo -n "-genesis-timestamp="
go run cmd/blockchain/blockchain.go -i blockchain.bin -timestamp=true
check_status
echo -n "-genesis-signature="
go run cmd/blocksigs/blocksigs.go -i blockchain.sigs
check_status
popd >/dev/null
popd >/dev/null
