#!/usr/bin/env bash

set -x

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo "skycoin binary dir:" "$DIR"
pushd "$DIR" >/dev/null

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

GORUNFLAGS=${GORUNFLAGS:-}

go run -ldflags "${GOLDFLAGS}" $GORUNFLAGS cmd/skycoin/skycoin.go \
    -gui-dir="${DIR}/src/gui/static/" \
    -launch-browser=true \
    -enable-api-set=READ,WALLET,STATUS \
    -enable-gui=true \
    -rpc-interface=false \
    -log-level=debug \
    $@

popd >/dev/null
