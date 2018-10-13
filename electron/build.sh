#!/usr/bin/env bash
set -e -o pipefail

# Builds both the electron and standalone releases

GOX_OSARCH="$@"

. build-conf.sh "$GOX_OSARCH"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

echo "Compiling with gox"
pwd
# Build the client mode with gox here so that the standalone and electron releases don't need to compile twice
CONFIG_MODE=STANDALONE_CLIENT ./gox.sh "$GOX_OSARCH" "$GOX_GUI_OUTPUT_DIR"

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

echo
echo "==========================="
echo "Building daemon release"

./build-daemon-release.sh "$GOX_OSARCH"

echo
echo "==========================="
echo "Building cli release"

./build-cli-release.sh "$GOX_OSARCH"

popd >/dev/null
