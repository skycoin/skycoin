#!/usr/bin/env bash

# Compresses packaged electron apps after
# ./package-electron-release.sh is done

. build-conf.sh

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

# Compress archives
pushd "$ELN_OUTPUT" >/dev/null

# OS X
pushd "darwin-x64" >/dev/null
if [ -e "$OSX64_ELN_ZIP" ]; then
    echo "Removing old $OSX64_ELN_ZIP"
    rm "$OSX64_ELN_ZIP"
fi
echo "Zipping $OSX64_ELN_ZIP"
# -y preserves symlinks,
# so that the massive .framework library isn't duplicated
zip -r -y --quiet "$OSX64_ELN_ZIP" "$OSX64_APP"
mv "$OSX64_ELN_ZIP" "../$OSX64_ELN_ZIP"
popd >/dev/null

# Windows
if [ -e "$WIN64_ELN_ZIP" ]; then
    echo "Removing old $WIN64_ELN_ZIP"
    rm "$WIN64_ELN_ZIP"
fi
echo "Zipping $WIN64_ELN_ZIP"
mv "win32-x64" "$WIN64_ELN"
zip -r --quiet "$WIN64_ELN_ZIP" "$WIN64_ELN"
mv "$WIN64_ELN" "win32-x64"

# Linux
if [ -e "$LNX64_ELN_ZIP" ]; then
    echo "Removing old $LNX64_ELN_ZIP"
    rm "$LNX64_ELN_ZIP"
fi
echo "Zipping $LNX64_ELN_ZIP"
mv "linux-x64" "$LNX64_ELN"
tar czf "$LNX64_ELN_ZIP" --owner=0 --group=0 "$LNX64_ELN"
mv "$LNX64_ELN" "linux-x64"

popd >/dev/null

# Move to final release dir
mkdir -p "$FINAL_OUTPUT"
mkdir -p "$FINAL_OUTPUT"
for var in "${OSX64_ELN_ZIP}" "${WIN64_ELN_ZIP}" "${LNX64_ELN_ZIP}"; do
    mv "${ELN_OUTPUT}/${var}" "$FINAL_OUTPUT"
done

popd >/dev/null
