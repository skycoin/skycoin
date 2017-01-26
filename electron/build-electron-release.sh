#!/usr/bin/env bash
set -e -o pipefail

# Builds an entire skycoin + electron-based GUI for release

# Implemented architectures:
#       darwin/amd64
#       windows/amd64
#       windows/386
#       linux/amd64
#
# By default builds all architectures.
# A single arch can be built by specifying it using gox's arch names

. build-conf.sh

SKIP_COMPILATION=${SKIP_COMPILATION:-0}

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

if [ $SKIP_COMPILATION -ne 1 ]; then
    ./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT"
fi

if [ -e "$ELN_OUTPUT" ]; then
    rm -r "$ELN_OUTPUT"
fi

if [ ! -z "$WIN64_ELN" ] && [ ! -z "$WIN32_ELN" ]; then
    npm run dist-win
fi

if [ ! -z "$LNX64_ELN" ]; then
    npm run dist-linux
fi

if [ ! -z "$OSX64_ELN" ]; then
    if [[ "$OSTYPE" == "darwin"* ]]; then
        npm run dist-mac
    elif [[ "$OSTYPE" == "linux"* ]]; then
        npm run pack-mac
    else
        echo "Can not run build script in $OSTYPE"
    fi
fi

pushd "$FINAL_OUTPUT" >/dev/null
if [ -e "mac" ]; then
    pushd "mac" >/dev/null
    if [ -e "Skycoin-${SKY_VERSION}.dmg" ]; then
        mv "Skycoin-${SKY_VERSION}.dmg" "../skycoin-${SKY_VERSION}-gui-osx-x64.dmg"
    elif [ -e "Skycoin.app" ]; then
        tar czf "../skycoin-${SKY_VERSION}-gui-osx-x64.zip" --owner=0 --group=0 "Skycoin.app"
    fi
    popd >/dev/null
    rm -rf "mac"
fi

IMG="skycoin-${SKY_VERSION}-x86_64.AppImage"
DEST_IMG="skycoin-${SKY_VERSION}-gui-linux-x64.AppImage"
if [ -e $IMG ]; then
    mv "$IMG" "$DEST_IMG"
    chmod +x "$DEST_IMG"
fi

EXE="Skycoin Setup ${SKY_VERSION}.exe"
if [ -e "$EXE" ]; then
    mv "$EXE" "skycoin-${SKY_VERSION}-gui-win-setup.exe"
fi

# clean unpacked folders
rm -rf *-unpacked

popd >/dev/null
popd >/dev/null
