#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"


pushd "$DIR" >/dev/null
pushd .. >/dev/null
echo "Generating new master keypair"
go run cmd/address/address.go -o master.keys
echo "Creating genesis block"
go run cmd/genesis/genesis.go -keys master.keys -help=false
echo "====================="
echo "NOTICE: Make sure to copy the blockchain files and master keys file to the -data-directory used by the master node"
echo "====================="
echo "CLI flags for non-master client:"
echo "====================="
echo -n "-master-public-key="
go run cmd/address/address.go -i master.keys -print-public=true
echo -n "-genesis-timestamp="
go run cmd/blockchain/blockchain.go -i blockchain.bin -timestamp=true
echo -n "-genesis-signature="
go run cmd/blocksigs/blocksigs.go -i blockchain.sigs
popd >/dev/null
popd >/dev/null
