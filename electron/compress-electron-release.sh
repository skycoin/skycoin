#!/usr/bin/env bash
set -e -o pipefail

# Compresses packaged electron apps after
# ./package-electron-release.sh is done

GOX_OSARCH="$@"

. build-conf.sh "$GOX_OSARCH"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

# Compress archives
pushd "$ELN_OUTPUT_DIR" >/dev/null

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

# Windows 64bit
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

# Windows 32bit
if [ -e "$WIN32_ELN_PLT" ]; then
    if [ -e "$WIN32_ELN_ZIP" ]; then
        echo "Removing old $WIN32_ELN_ZIP"
        rm "$WIN32_ELN_ZIP"
    fi
    echo "Zipping $WIN32_ELN_ZIP"
    mv "$WIN32_ELN_PLT" "$WIN32_ELN"
    zip -r --quiet "$WIN32_ELN_ZIP" "$WIN32_ELN"
    mv "$WIN32_ELN" "$WIN32_ELN_PLT"
    FINALS+=("$WIN32_ELN_ZIP")
fi

# Linux
if [ -e "$LNX64_ELN_PLT" ]; then
    if [ -e "$LNX64_ELN_ZIP" ]; then
        echo "Removing old $LNX64_ELN_ZIP"
        rm "$LNX64_ELN_ZIP"
    fi
    echo "Zipping $LNX64_ELN_ZIP"
    mv "$LNX64_ELN_PLT" "$LNX64_ELN"
    if [[ "$OSTYPE" == "linux"* ]]; then
        tar czf "$LNX64_ELN_ZIP" --owner=0 --group=0 "$LNX64_ELN"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        tar czf "$LNX64_ELN_ZIP"  "$LNX64_ELN"
    fi
    mv "$LNX64_ELN" "$LNX64_ELN_PLT"
    FINALS+=("$LNX64_ELN_ZIP")
fi

popd >/dev/null

# Move to final release dir
mkdir -p "$FINAL_OUTPUT_DIR"
for var in "${FINALS[@]}"; do
    mv "${ELN_OUTPUT_DIR}/${var}" "$FINAL_OUTPUT_DIR"
done

popd >/dev/null
