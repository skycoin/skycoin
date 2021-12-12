#!/usr/bin/env bash

# Runs the node with configuration necessary for running the live integration tests,
# with coverage enabled

set -exu

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd "$DIR/.." >/dev/null

COIN=${COIN:-skycoin}
COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"
BINARY="${COIN}-live.test"
CMDPKG=$(go list ./cmd/${COIN})
COVERPKG=$(dirname $(dirname ${CMDPKG}))

COVERAGEFILE="coverage/${BINARY}.coverage.out"
if [ -f "${COVERAGEFILE}" ]; then
    rm "${COVERAGEFILE}"
fi

go test -c -ldflags "-X ${CMDPKG}.Commit=$COMMIT -X ${CMDPKG}.Branch=${BRANCH}" -tags testrunmain -o "${BINARY}" -coverpkg="${COVERPKG}/..." ./cmd/${COIN}/

./${BINARY} \
    -gui-dir="${DIR}/src/gui/static/" \
    -launch-browser=false \
    -enable-all-api-sets=true \
    -enable-api-sets=INSECURE_WALLET_SEED \
    -test.run "^TestRunMain$" \
    -test.coverprofile="${COVERAGEFILE}" \
    $@

popd >/dev/null
