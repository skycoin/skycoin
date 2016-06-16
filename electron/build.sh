#!/usr/bin/env bash

# Builds an entire skycoin + electron-based GUI for release

# Implemented architectures:
#       darwin_amd64
#       windows_amd64
#       linux_amd64

. build-conf.sh

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT"
if [ $? -ne 0 ]; then
    echo "gox build failed"
    exit 1
fi

rm -r .electron_output
gulp electron
if [ $? -ne 0 ]; then
    echo "gulp electron failed"
    exit 1
fi

./package-release.sh
if [ $? -ne 0 ]; then
    echo "package-release.sh failed"
    exit 1
fi

./compress-release.sh
if [ $? -ne 0 ]; then
    echo "compress-release.sh failed"
    exit 1
fi

popd >/dev/null
