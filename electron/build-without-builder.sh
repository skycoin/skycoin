#!/usr/bin/env bash
set -e -o pipefail

# Builds both the electron and standalone releases

. build-conf.sh

if [ -n "$1" ]; then
    GOX_OSARCH="$1"
fi

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

WITH_BUILDER=0

pushd "$SCRIPTDIR" >/dev/null

echo "Compiling with gox"

# Build with gox here and make the other scripts skip it
./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT" "$WITH_BUILDER"

echo "Installing node modules"
./install-node-modules.sh

echo
echo "==========================="
echo "Building standalone release"

SKIP_COMPILATION=1 ./build-standalone-release.sh "$GOX_OSARCH" "$WITH_BUILDER"

echo
echo "==========================="
echo "Building electron release"

SKIP_COMPILATION=1 ./build-electron-release.sh "$GOX_OSARCH" "$WITH_BUILDER"

popd >/dev/null
