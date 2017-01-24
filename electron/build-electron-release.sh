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

# GULP_PLATFORM=""
# if [ -n "$1" ]; then
#     GOX_OSARCH="$1"
#     case "$1" in
#     "linux/amd64")
#         GULP_PLATFORM="linux-x64"
#         ;;
#     "linux/arm")
#         GULP_PLATFORM="linux-arm"
#         ;;
#     "windows/amd64")
#         GULP_PLATFORM="win32-x64"
#         ;;
#     "windows/386")
#         GULP_PLATFORM="win32-ia32"
#         ;;
#     "darwin/amd64")
#         GULP_PLATFORM="darwin-x64"
#         ;;
#     esac
# fi

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

if [ $SKIP_COMPILATION -ne 1 ]; then
    ./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT"
fi

if [ -e "$ELN_OUTPUT" ]; then
    rm -r "$ELN_OUTPUT"
fi

npm run dist-win
npm run dist-l

if [[ "$OSTYPE" == "darwin"* ]]; then
    npm run dist-m
fi

# if [ -n "$GULP_PLATFORM" ]; then
#     gulp electron --platform "$GULP_PLATFORM"
# else
#     gulp electron
# fi

pushd "$FINAL_OUTPUT" >/dev/null
if [ -e "mac" ]; then
    pushd "mac" >/dev/null
    mv "Skycoin-${SKY_VERSION}.dmg" "../skycoin-${SKY_VERSION}-gui-osx-x64.dmg"
    popd >/dev/null
    rm -rf "mac"
fi

IMG="skycoin-${SKY_VERSION}-x86_64.AppImage"
DEST_IMG="skycoin-${SKY_VERSION}-gui-linux-x64.AppImage"
if [ -e $IMG ]; then
    mv "$IMG" "$DEST_IMG"
    chmod +x "$DEST_IMG"
fi

# mv "skycoin-${SKY_VERSION}-x86_64.AppImage" "skycoin-${SKY_VERSION}-gui-linux-x64.AppImage"
# chmod +x "skycoin-${SKY_VERSION}-gui-linux-x64.AppImage"
EXE="Skycoin Setup ${SKY_VERSION}.exe"
if [ -e "$EXE" ]; then
    mv "$EXE" "skycoin-${SKY_VERSION}-gui-win-installer.exe"
fi

# mv "Skycoin Setup ${SKY_VERSION}.exe" "skycoin-${SKY_VERSION}-gui-win-insaller.exe"

# clean unpacked folders
rm -rf *-unpacked

popd >/dev/null


# echo "--------------------------"
# echo "Packaging electron release"
# ./package-electron-release.sh

# echo "----------------------------"
# echo "Compressing electron release"
# ./compress-electron-release.sh

popd >/dev/null
