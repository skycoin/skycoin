#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$DIR" >/dev/null
./compile/clean-static-libs.sh >/dev/null 2>&1
go run cmd/skycoin/skycoin.go $@

popd >/dev/null
