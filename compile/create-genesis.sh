#!/usr/bin/env bash

BCFILE=blockchain.bin
BSFILE=blockchain.sigs
MASTERKEYS=master.keys
BENEFACTORWALLET=benefactor-wallet.json

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
echo "Generating new master genesis keypair"
# TODO -- remove address.go, use wallet.go
go run cmd/address/address.go -o="$MASTERKEYS" -address-version="$ADDRV"
check_status
echo "Generating benefactor wallet"
BENEFACTOR=`go run cmd/wallet/wallet.go -entries=1 -o="$BENEFACTORWALLET" \
            -print-address=true -address-version="$ADDRV" -entry=0`
check_status
if [ -z "$BENEFACTOR" ]; then
    echo "Failed to extract benefactor address."
    echo "cmd/wallet/wallet.go may be broken."
    exit 1
fi
echo "Creating genesis block and transferring balance to benefactor"
go run cmd/genesis/genesis.go -keys="$MASTERKEYS" -help=false \
                              -address-version="$ADDRV" \
                              -dest-address="$BENEFACTOR"
check_status
echo "====================="
echo "NOTICE: Make sure to copy ${BCFILE}, ${BSFILE}, and ${MASTERKEYS} to \
the -data-directory used by the master node."
echo "Save ${BENEFACTORWALLET} somewhere safe.  This controls the genesis \
balance. Do not distribute it and do not use it in the master chain."
echo "====================="
echo "CLI flags for non-master client:"
echo "====================="
echo -n "-master-public-key="
go run cmd/address/address.go -i master.keys -print-public=true
check_status
echo -n "-genesis-timestamp="
go run cmd/blockchain/blockchain.go -i blockchain.bin -timestamp=true -b=0
check_status
echo -n "-genesis-signature="
go run cmd/blocksigs/blocksigs.go -i blockchain.sigs -b=0
check_status
popd >/dev/null
popd >/dev/null
