#!/usr/bin/env bash

set -x

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo "skycoin binary dir:" "$DIR"
pushd "$DIR" >/dev/null

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

go run -ldflags "${GOLDFLAGS}" cmd/skycoin/skycoin.go \
    -gui-dir="${DIR}/src/gui/static/" \
    -launch-browser=true \
    -enable-wallet-api=true \
    -enable-gui=true \
    -rpc-interface=false \
    -log-level=debug \
    $@

popd >/dev/null
