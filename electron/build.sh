#!/usr/bin/env bash

# Builds both the electron and standalone releases

. build-conf.sh

if [ -n "$1" ]; then
    GOX_OSARCH="$1"
fi

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

echo "Compiling with gox"

# Build with gox here and make the other scripts skip it
./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT"
if [ $? -ne 0 ]; then
    echo "gox build failed"
    exit 1
fi

echo "Building standalone release"

SKIP_COMPILATION=1 ./build-standalone-release.sh "$GOX_OSARCH"
if [ $? -ne 0 ]; then
    echo "build-standalone-release.sh failed"
    exit 1
fi

echo "Building electron release"

SKIP_COMPILATION=1 ./build-electron-release.sh "$GOX_OSARCH"
if [ $? -ne 0 ]; then
    echo "build-electron-release.sh failed"
    exit 1
fi


popd >/dev/null
