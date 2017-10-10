#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo "skycoin binary dir:" "$DIR"
pushd "$DIR" >/dev/null
go run cmd/skycoin/skycoin.go --gui-dir="${DIR}/src/gui/static/" $@
popd >/dev/null
