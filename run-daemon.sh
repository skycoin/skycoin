#!/usr/bin/env bash

# Runs skycoin in daemon mode configuration
export USER_BURN_FACTOR=0
set -x

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo "Ness binary dir:" "$DIR"
pushd "$DIR" >/dev/null

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

GORUNFLAGS=${GORUNFLAGS:-}

go run -ldflags "${GOLDFLAGS}" $GORUNFLAGS cmd/privateness/privateness.go \
    -max-default-peer-outgoing-connections=26 \
    -enable-gui=false \
    -launch-browser=false \
    -log-level=debug \
    -web-interface-https \
    $@

popd >/dev/null
