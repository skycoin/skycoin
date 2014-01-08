#!/usr/bin/env bash

cd compile/
./build-linux-x86_64.sh dev

if [[ $? != 0 ]]; then
    exit 1
fi

./release/skycoin_linux_x86_64/skycoin -disable-gui=false $@

