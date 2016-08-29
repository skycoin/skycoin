#!/usr/bin/env bash

. build-conf.sh

SKIP_COMPILATION=${SKIP_COMPILATION:-0}

if [ -n "$1" ]; then
    GOX_OSARCH="$1"
fi

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

if [ $SKIP_COMPILATION -ne 1 ]; then
    ./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT"
    if [ $? -ne 0 ]; then
        echo "gox build failed"
        exit 1
    fi
fi

./package-standalone-release.sh
if [ $? -ne 0 ]; then
    echo "package-standalone-release.sh failed"
    exit 1
fi

./compress-standalone-release.sh
if [ $? -ne 0 ]; then
    echo "compress-standalone-release.sh failed"
    exit 1
fi

popd >/dev/null
