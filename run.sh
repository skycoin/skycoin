#!/usr/bin/env bash

set -x
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo "skycoin binary dir:" "$DIR"
pushd "$DIR" >/dev/null

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"
PORT=$(shuf -n 1 -i 1024-65535)
# PORT=$(tr -dc 0-9 < /dev/urandom | head -c 4 | xargs)
IS_READY=$(fuser -n tcp $PORT)

while [[ ! $IS_READY = "" ]]; do
	# PORT=$(tr -dc 0-9 < /dev/urandom | head -c 4 | xargs)
	PORT=$(shuf -n 1 -i 1024-65535)
	IS_READY=$(fuser -n tcp $PORT)
done


go run -ldflags "${GOLDFLAGS}" cmd/skycoin/skycoin.go \
    -gui-dir="${DIR}/src/gui/static/" \
    -launch-browser=true \
    -enable-wallet-api=true \
    -rpc-interface=false \
    -web-interface-port $PORT   \
    $@

popd >/dev/null
