#!/usr/bin/env bash

# Compresses packaged electron apps after ./package-release.sh is done

. build-conf.sh

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

FINAL_OUTPUT="release"

pushd "$SCRIPTDIR" >/dev/null

# Compress archives
pushd "$ELN_OUTPUT" >/dev/null

# OS X
pushd "darwin-x64" >/dev/null
OSX_ZIP="skycoin-$SKY_VERSION-osx-darwin-x64.zip"
if [ -e "$OSX_ZIP" ]; then
    echo "Removing old $OSX_ZIP"
    rm "$OSX_ZIP"
fi
echo "Zipping $OSX_ZIP"
# -y preserves symlinks,
# so that the massive .framework library isn't duplicated
zip -r -y --quiet "$OSX_ZIP" "Skycoin.app"
mv "$OSX_ZIP" "../$OSX_ZIP"
popd >/dev/null

# Windows
WIN64_PRE="skycoin-$SKY_VERSION-win-x64"
WIN64_ZIP="${WIN64_PRE}.zip"

if [ -e "$WIN64_ZIP" ]; then
    echo "Removing old $WIN64_ZIP"
    rm "$WIN64_ZIP"
fi
echo "Zipping $WIN64_ZIP"
mv "win32-x64" "$WIN64_PRE"
zip -r --quiet "$WIN64_ZIP" "$WIN64_PRE"
mv "$WIN64_PRE" "win32-x64"

# Linux
LNX64_PRE="skycoin-$SKY_VERSION-linux-x64"
LNX64_ZIP="${LNX64_PRE}.tar.gz"

if [ -e "$LNX64_ZIP" ]; then
    echo "Removing old $LNX64_ZIP"
    rm "$LNX64_ZIP"
fi
echo "Zipping $LNX64_ZIP"
mv "linux-x64" "$LNX64_PRE"
tar czf "$LNX64_ZIP" --owner=0 --group=0 "$LNX64_PRE"
mv "$LNX64_PRE" "linux-x64"

popd >/dev/null

# Move to final release dir
mkdir -p "$FINAL_OUTPUT"
mv "$ELN_OUTPUT/"*.zip "$FINAL_OUTPUT"
mv "$ELN_OUTPUT/"*.tar.gz "$FINAL_OUTPUT"

popd >/dev/null
