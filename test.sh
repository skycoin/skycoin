#!/usr/bin/env bash

GOCMD="go test -v"

declare -a libs=(./src/lib/secp256k1-go)
declare -a pkgs=(. ./src/cli ./src/gui ./src/util ./src/coin ./src/daemon ../src/skycoin)

for i in "${pkgs[@]}" 
do
    $GOCMD "$i"
done

for i in "${libs[@]}" 
do
    $GOCMD "$i"
done
