#!/usr/bin/env bash
set -e -o pipefail

# Builds an entire skycoin + electron-based GUI for release

# Implemented architectures:
#       darwin/amd64
#       windows/amd64
#       linux/amd64
#
# By default builds all architectures.
# A single arch can be built by specifying it using gox's arch names

. build-conf.sh

SKIP_COMPILATION=${SKIP_COMPILATION:-0}

GULP_PLATFORM=""
if [ -n "$1" ]; then
    GOX_OSARCH="$1"
    case "$1" in
    "linux/amd64")
        GULP_PLATFORM="linux-x64"
        ;;
    "windows/amd64")
        GULP_PLATFORM="win32-x64"
        ;;
    "darwin/amd64")
        GULP_PLATFORM="darwin-x64"
        ;;
    esac
fi

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

if [ $SKIP_COMPILATION -ne 1 ]; then
    ./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT"
fi

if [ -e "$ELN_OUTPUT" ]; then
    rm -r "$ELN_OUTPUT"
fi

if [ -n "$GULP_PLATFORM" ]; then
    gulp electron --platform "$GULP_PLATFORM"
else
    gulp electron
fi

echo "--------------------------"
echo "Packaging electron release"
./package-electron-release.sh

echo "----------------------------"
echo "Compressing electron release"
./compress-electron-release.sh

popd >/dev/null
