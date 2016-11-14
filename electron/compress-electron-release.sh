#!/usr/bin/env bash
set -e -o pipefail

# Compresses packaged electron apps after
# ./package-electron-release.sh is done

. build-conf.sh

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

# Compress archives
pushd "$ELN_OUTPUT" >/dev/null

FINALS=()

# OS X
if [ -e "$OSX64_ELN_PLT" ]; then
    pushd "$OSX64_ELN_PLT" >/dev/null

    if [ -e "$OSX64_ELN_ZIP" ]; then
        echo "Removing old $OSX64_ELN_ZIP"
        rm "$OSX64_ELN_ZIP"
    fi
    echo "Zipping $OSX64_ELN_ZIP"
    # -y preserves symlinks,
    # so that the massive .framework library isn't duplicated
    zip -r -y --quiet "$OSX64_ELN_ZIP" "$OSX64_APP"
    mv "$OSX64_ELN_ZIP" "../$OSX64_ELN_ZIP"
    FINALS+=("$OSX64_ELN_ZIP")

    popd >/dev/null
fi

# Windows
if [ -e "$WIN64_ELN_PLT" ]; then
    if [ -e "$WIN64_ELN_ZIP" ]; then
        echo "Removing old $WIN64_ELN_ZIP"
        rm "$WIN64_ELN_ZIP"
    fi
    echo "Zipping $WIN64_ELN_ZIP"
    mv "$WIN64_ELN_PLT" "$WIN64_ELN"
    zip -r --quiet "$WIN64_ELN_ZIP" "$WIN64_ELN"
    mv "$WIN64_ELN" "$WIN64_ELN_PLT"
    FINALS+=("$WIN64_ELN_ZIP")
fi

# Linux
if [ -e "$LNX64_ELN_PLT" ]; then
    if [ -e "$LNX64_ELN_ZIP" ]; then
        echo "Removing old $LNX64_ELN_ZIP"
        rm "$LNX64_ELN_ZIP"
    fi
    echo "Zipping $LNX64_ELN_ZIP"
    mv "$LNX64_ELN_PLT" "$LNX64_ELN"
    tar czf "$LNX64_ELN_ZIP" --owner=0 --group=0 "$LNX64_ELN"
    mv "$LNX64_ELN" "$LNX64_ELN_PLT"
    FINALS+=("$LNX64_ELN_ZIP")
fi

popd >/dev/null

# Move to final release dir
mkdir -p "$FINAL_OUTPUT"
for var in "${FINALS[@]}"; do
    mv "${ELN_OUTPUT}/${var}" "$FINAL_OUTPUT"
done

popd >/dev/null
