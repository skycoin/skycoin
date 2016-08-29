#!/usr/bin/env bash

# Copies gox-compiled skycoin binaries and compiled GUI assets
# into an electron package

. build-conf.sh

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

OSX64="${ELN_OUTPUT}/darwin-x64"
WIN64="${ELN_OUTPUT}/win32-x64"
LNX64="${ELN_OUTPUT}/linux-x64"

OSX64_RES="${OSX64}/${OSX64_APP}/Contents/Resources/app"
WIN64_RES="${WIN64}/resources/app"
LNX64_RES="${LNX64}/resources/app"

OSX64_SRC="${OSX64_RES}/src"
WIN64_SRC="${WIN64}/src"
LNX64_SRC="${LNX64}/src"

# Capitalize OS X .app for convention
if [ -e "${OSX64}/skycoin.app" ]; then
    mv "${OSX64}/skycoin.app" "${OSX64}/${OSX64_APP}"
fi

echo "Copying skycoin binaries"

# Copy binaries to electron app
cp "${GOX_OUTPUT}/skycoin_darwin_amd64" "${OSX64_RES}/skycoin"
cp "${GOX_OUTPUT}/skycoin_windows_amd64.exe" "${WIN64_RES}/skycoin.exe"
cp "${GOX_OUTPUT}/skycoin_linux_amd64" "${LNX64_RES}/skycoin"

echo "Copying static resources"

# Copy static resources to electron app
cp -R "$GUI_DIST_DIR" "$OSX64_RES"
cp -R "$GUI_DIST_DIR" "$WIN64_RES"
cp -R "$GUI_DIST_DIR" "$LNX64_RES"

# Copy the source for reference
# tar it with filters, move it, then untar in order to do this
echo "Copying source snapshot"

./package-source.sh "${OSX64_SRC}" "${WIN64_SRC}" "${LNX64_SRC}"
