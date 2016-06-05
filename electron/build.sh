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
gulp electron
./package-release.sh

popd >/dev/null
