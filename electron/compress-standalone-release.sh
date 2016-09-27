#!/usr/bin/env bash
set -e -o pipefail

# Compresses packaged standalone release after
# ./package-standalone-release.sh is done

. build-conf.sh

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

# Compress archives
pushd "$STL_OUTPUT" >/dev/null

FINALS=()

# OS X
if [ -e "$OSX64_STL" ]; then
    if [ -e "$OSX64_STL_ZIP" ]; then
        echo "Removing old $OSX64_STL_ZIP"
        rm "$OSX64_STL_ZIP"
    fi
    echo "Zipping $OSX64_STL_ZIP"
    # -y preserves symlinks,
    # so that the massive .framework library isn't duplicated
    zip -r -y --quiet "$OSX64_STL_ZIP" "$OSX64_STL"
    FINALS+=("$OSX64_STL_ZIP")
fi

# Windows
if [ -e "$WIN64_STL" ]; then
    if [ -e "$WIN64_STL_ZIP" ]; then
        echo "Removing old $WIN64_STL_ZIP"
        rm "$WIN64_STL_ZIP"
    fi
    echo "Zipping $WIN64_STL_ZIP"
    zip -r --quiet "$WIN64_STL_ZIP" "$WIN64_STL"
    FINALS+=("$WIN64_STL_ZIP")
fi


# Linux
if [ -e "$LNX64_STL" ]; then
    if [ -e "$LNX64_STL_ZIP" ]; then
        echo "Removing old $LNX64_STL_ZIP"
        rm "$LNX64_STL_ZIP"
    fi
    echo "Zipping $LNX64_STL_ZIP"
    tar czf "$LNX64_STL_ZIP" --owner=0 --group=0 "$LNX64_STL"
    FINALS+=("$LNX64_STL_ZIP")
fi

popd >/dev/null

# Move to final release dir
mkdir -p "$FINAL_OUTPUT"
for var in "${FINALS[@]}"; do
    mv "${STL_OUTPUT}/${var}" "$FINAL_OUTPUT"
done

popd >/dev/null
