#!/usr/bin/env bash
set -e -o pipefail

# Builds both the electron and standalone releases

. build-conf.sh "$1"

if [ -n "$1" ]; then
    GOX_OSARCH="$1"
fi

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

echo "Compiling with gox"
pwd
# Build with gox here and make the other scripts skip it
./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT"

echo "Installing node modules"
./install-node-modules.sh

echo
echo "==========================="
echo "Building standalone release"

SKIP_COMPILATION=1 ./build-standalone-release.sh "$GOX_OSARCH"

echo
echo "==========================="
echo "Building electron release"

SKIP_COMPILATION=1 ./build-electron-release.sh "$GOX_OSARCH"

popd >/dev/null
