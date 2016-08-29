#!/usr/bin/env bash

# Compresses packaged standalone release after
# ./package-standalone-release.sh is done

. build-conf.sh

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

# Compress archives
pushd "$STL_OUTPUT" >/dev/null

# OS X
if [ -e "$OSX64_STL_ZIP" ]; then
    echo "Removing old $OSX64_STL_ZIP"
    rm "$OSX64_STL_ZIP"
fi
echo "Zipping $OSX64_STL_ZIP"
# -y preserves symlinks,
# so that the massive .framework library isn't duplicated
zip -r -y --quiet "$OSX64_STL_ZIP" "$OSX64_STL"

# Windows
if [ -e "$WIN64_STL_ZIP" ]; then
    echo "Removing old $WIN64_STL_ZIP"
    rm "$WIN64_STL_ZIP"
fi
echo "Zipping $WIN64_STL_ZIP"
zip -r --quiet "$WIN64_STL_ZIP" "$WIN64_STL"

# Linux
if [ -e "$LNX64_STL_ZIP" ]; then
    echo "Removing old $LNX64_STL_ZIP"
    rm "$LNX64_STL_ZIP"
fi
echo "Zipping $LNX64_STL_ZIP"
tar czf "$LNX64_STL_ZIP" --owner=0 --group=0 "$LNX64_STL"

popd >/dev/null

# Move to final release dir
mkdir -p "$FINAL_OUTPUT"
for var in "${OSX64_STL_ZIP}" "${WIN64_STL_ZIP}" "${LNX64_STL_ZIP}"; do
    mv "${STL_OUTPUT}/${var}" "$FINAL_OUTPUT"
done

popd >/dev/null
