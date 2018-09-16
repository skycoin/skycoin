#!/usr/bin/env bash

# Runs the node with configuration necessary for running the live integration tests

set -xu

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd "$DIR/.." >/dev/null

COIN=${COIN:-skycoin}
COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

go run -ldflags "${GOLDFLAGS}" cmd/${COIN}/${COIN}.go \
    -gui-dir="${DIR}/src/gui/static/" \
    -launch-browser=false \
    -enable-api-set=ALL \
    $@

popd >/dev/null
