#!/usr/bin/env bash
set -e -o pipefail

# Compresses packaged standalone release after
# ./package-standalone-release.sh is done

GOX_OSARCH="$@"

. build-conf.sh "$GOX_OSARCH"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

# Compress archives
pushd "$STL_OUTPUT_DIR" >/dev/null

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

# Windows 64bit
if [ -e "$WIN64_STL" ]; then
    if [ -e "$WIN64_STL_ZIP" ]; then
        echo "Removing old $WIN64_STL_ZIP"
        rm "$WIN64_STL_ZIP"
    fi
    echo "Zipping $WIN64_STL_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        zip -r --quiet -X "$WIN64_STL_ZIP"  "$WIN64_STL"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        zip -r --quiet "$WIN64_STL_ZIP" "$WIN64_STL"
    elif [[ "$OSTYPE" == "msys"* ]]; then
        7z a "$WIN64_STL_ZIP" "$WIN64_STL"
    fi
    FINALS+=("$WIN64_STL_ZIP")
fi

# Windows 32bit
if [ -e "$WIN32_STL" ]; then
    if [ -e "$WIN32_STL_ZIP" ]; then
        echo "Removing old $WIN32_STL_ZIP"
        rm "$WIN32_STL_ZIP"
    fi
    echo "Zipping $WIN32_STL_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        zip -r --quiet -X "$WIN32_STL_ZIP"  "$WIN32_STL"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        zip -r --quiet "$WIN32_STL_ZIP" "$WIN32_STL"
    elif [[ "$OSTYPE" == "msys"* ]]; then
        7z a "$WIN32_STL_ZIP" "$WIN32_STL"
    fi
    FINALS+=("$WIN32_STL_ZIP")
fi

# Linux
if [ -e "$LNX64_STL" ]; then
    if [ -e "$LNX64_STL_ZIP" ]; then
        echo "Removing old $LNX64_STL_ZIP"
        rm "$LNX64_STL_ZIP"
    fi
    echo "Zipping $LNX64_STL_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        tar czf "$LNX64_STL_ZIP" --owner=0 --group=0 "$LNX64_STL"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        tar czf "$LNX64_STL_ZIP"  "$LNX64_STL"
    fi
    FINALS+=("$LNX64_STL_ZIP")
fi

# Linux arm
if [ -e "$LNX_ARM_STL" ]; then
    if [ -e "$LNX_ARM_STL_ZIP" ]; then
        echo "Removing old $LNX_ARM_STL_ZIP"
        rm "$LNX_ARM_STL_ZIP"
    fi
    echo "Zipping $LNX_ARM_STL_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        tar czf "$LNX_ARM_STL_ZIP" --owner=0 --group=0 "$LNX_ARM_STL"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        tar czf "$LNX_ARM_STL_ZIP"  "$LNX_ARM_STL"
    fi
    FINALS+=("$LNX_ARM_STL_ZIP")
fi

popd >/dev/null

# Move to final release dir
mkdir -p "$FINAL_OUTPUT_DIR"
for var in "${FINALS[@]}"; do
    mv "${STL_OUTPUT_DIR}/${var}" "$FINAL_OUTPUT_DIR"
done

popd >/dev/null
