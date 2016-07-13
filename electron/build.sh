#!/usr/bin/env bash

# Builds an entire skycoin + electron-based GUI for release

# Implemented architectures:
#       darwin/amd64
#       windows/amd64
#       linux/amd64
#
# By default builds all architectures.
# A single arch can be built by specifying it using gox's arch names

. build-conf.sh

GULP_PLATFORM=""
ARCH_RESTRICTION=""
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
    *)
        echo "Unknown build arch $1"
        exit 1
        ;;
    esac
fi

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT"
if [ $? -ne 0 ]; then
    echo "gox build failed"
    exit 1
fi

rm -r .electron_output
if [ -n "$GULP_PLATFORM" ]; then
    gulp electron --platform "$GULP_PLATFORM"
else
    gulp electron
fi
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
