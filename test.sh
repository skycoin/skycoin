#!/usr/bin/env bash

GOCMD="go test -v"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

declare -a libs=(./src/lib/secp256k1-go)
declare -a pkgs=(./src/cli ./src/gui ./src/util ./src/coin ./src/daemon ./src/skycoin)
declare -a cmds=(./cmd/skycoin ./cmd/skycoind ./cmd/skycoindev ./cmd/blockchain)

pushd "$DIR" >/dev/null

for i in "${pkgs[@]}" 
do
    $GOCMD "$i"
done

for i in "${libs[@]}" 
do
    $GOCMD "$i"
done

for i in "${cmds[@]}"
do
    $GOCMD "$i"
done

popd >/dev/null
