#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"


pushd "$DIR" >/dev/null
pushd .. >/dev/null
echo "Generating new master keypair"
go run cmd/address/address.go -o master.keys
echo "Creating genesis block"
go run cmd/genesis/genesis.go -keys master.keys -help=false
echo "=================="
echo "=================="
echo "Timestamp:"
go run cmd/blockchain/blockchain.go -i blockchain.bin -timestamp=true
echo "Signature:"
go run cmd/blocksigs/blocksigs.go -i blockchain.sigs
popd >/dev/null
popd >/dev/null
