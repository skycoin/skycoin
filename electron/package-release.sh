#!/usr/bin/env bash

# Copies gox-compiled skycoin binaries and compiled GUI assets
# into an electron package, then compresses them

. build-conf.sh

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

OSX64="${ELN_OUTPUT}/darwin-x64"
WIN64="${ELN_OUTPUT}/win32-x64"
LNX64="${ELN_OUTPUT}/linux-x64"

OSX64_RES="${OSX64}/skycoin.app/Contents/Resources/app"
WIN64_RES="${WIN64}/resources/app"
LNX64_RES="${LNX64}/resources/app"

FINAL_OUTPUT="release"

# Copy binaries to electron app
cp "${GOX_OUTPUT}/skycoin_darwin_amd64" "${OSX64_RES}/skycoin"
cp "${GOX_OUTPUT}/skycoin_windows_amd64.exe" "${WIN64_RES}/skycoin.exe"
cp "${GOX_OUTPUT}/skycoin_linux_amd64" "${LNX64_RES}/skycoin"

# Copy static resources to electron app
cp -R "$GUI_DIST_DIR" "$OSX64_RES"
cp -R "$GUI_DIST_DIR" "$WIN64_RES"
cp -R "$GUI_DIST_DIR" "$LNX64_RES"

# Capitalize OS X .app for convention
mv "${OSX64}/skycoin.app" "${OSX64}/Skycoin.app"

# Compress archives

# OS X
pushd "$OSX64" >/dev/null
OSX_ZIP="Skycoin-$SKY_VERSION.zip"
if [ -e "$OSX_ZIP" ]; then
    echo "Removing old $OSX_ZIP"
    rm "$OSX_ZIP"
fi
echo "Zipping $OSX_ZIP"
zip -r --quiet "$OSX_ZIP" "Skycoin.app"
mv "$OSX_ZIP" "../$OSX_ZIP"
popd >/dev/null

# Windows, linux
pushd "$ELN_OUTPUT" >/dev/null

WIN64_PRE="skycoin-$SKY_VERSION-win-x64"
WIN64_ZIP="${WIN64_PRE}.zip"

if [ -e "$WIN64_ZIP" ]; then
    echo "Removing old $WIN64_ZIP"
    rm "$WIN64_ZIP"
fi
echo "Zipping $WIN64_ZIP"
mv "win32-x64" "$WIN64_PRE"
zip -r --quiet "$WIN64_ZIP" "$WIN64_PRE"

LNX64_PRE="skycoin-$SKY_VERSION-linux-x64"
LNX64_ZIP="${LNX64_PRE}.tar.gz"

if [ -e "$LNX64_ZIP" ]; then
    echo "Removing old $LNX64_ZIP"
    rm "$LNX64_ZIP"
fi
echo "Zipping $LNX64_ZIP"
mv "linux-x64" "$LNX64_PRE"
tar czf "$LNX64_ZIP" "$LNX64_PRE"

popd >/dev/null

# Move to final release dir
mkdir -p "$FINAL_OUTPUT"
mv "$ELN_OUTPUT/"*.zip "$FINAL_OUTPUT"
mv "$ELN_OUTPUT/"*.tar.gz "$FINAL_OUTPUT"

popd >/dev/null
