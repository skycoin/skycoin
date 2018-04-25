#!/usr/bin/env bash

set -x
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo "skycoin binary dir:" "$DIR"
pushd "$DIR" >/dev/null

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"
PORT=$(( ( (($RANDOM+$RANDOM)%64512)+1024 ) ))
IS_AVAILABLE=$(netstat -an | grep $PORT)

while [[ ! $IS_AVAILABLE = "" ]]; do
	PORT=$(( ( (($RANDOM+$RANDOM)%64512)+1024 ) ))
	IS_AVAILABLE=$(netstat -an | grep $PORT)
done


go run -ldflags "${GOLDFLAGS}" cmd/skycoin/skycoin.go \
    -gui-dir="${DIR}/src/gui/static/" \
    -launch-browser=true \
    -enable-wallet-api=true \
    -rpc-interface=false \
    -log-level=debug \
    -web-interface-port $PORT   \
    $@

popd >/dev/null
