#!/usr/bin/env bash
set -e -o pipefail

# Compresses packaged cli release after
# ./package-cli-release.sh is done

GOX_OSARCH="$@"

. build-conf.sh "$GOX_OSARCH"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

# Compress archives
pushd "$CLI_OUTPUT_DIR" >/dev/null

FINALS=()

# OS X
if [ -e "$OSX64_CLI" ]; then
    if [ -e "$OSX64_CLI_ZIP" ]; then
        echo "Removing old $OSX64_CLI_ZIP"
        rm "$OSX64_CLI_ZIP"
    fi
    echo "Zipping $OSX64_CLI_ZIP"
    # -y preserves symlinks,
    # so that the massive .framework library isn't duplicated
    zip -r -y --quiet "$OSX64_CLI_ZIP" "$OSX64_CLI"
    FINALS+=("$OSX64_CLI_ZIP")
fi

# Windows 64bit
if [ -e "$WIN64_CLI" ]; then
    if [ -e "$WIN64_CLI_ZIP" ]; then
        echo "Removing old $WIN64_CLI_ZIP"
        rm "$WIN64_CLI_ZIP"
    fi
    echo "Zipping $WIN64_CLI_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        zip -r --quiet -X "$WIN64_CLI_ZIP"  "$WIN64_CLI"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        zip -r --quiet "$WIN64_CLI_ZIP" "$WIN64_CLI"
    elif [[ "$OSTYPE" == "msys"* ]]; then
        7z a "$WIN64_CLI_ZIP" "$WIN64_CLI"
    fi
    FINALS+=("$WIN64_CLI_ZIP")
fi

# Windows 32bit
if [ -e "$WIN32_CLI" ]; then
    if [ -e "$WIN32_CLI_ZIP" ]; then
        echo "Removing old $WIN32_CLI_ZIP"
        rm "$WIN32_CLI_ZIP"
    fi
    echo "Zipping $WIN32_CLI_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        zip -r --quiet -X "$WIN32_CLI_ZIP"  "$WIN32_CLI"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        zip -r --quiet "$WIN32_CLI_ZIP" "$WIN32_CLI"
    elif [[ "$OSTYPE" == "msys"* ]]; then
        7z a "$WIN32_CLI_ZIP" "$WIN32_CLI"
    fi
    FINALS+=("$WIN32_CLI_ZIP")
fi

# Linux
if [ -e "$LNX64_CLI" ]; then
    if [ -e "$LNX64_CLI_ZIP" ]; then
        echo "Removing old $LNX64_CLI_ZIP"
        rm "$LNX64_CLI_ZIP"
    fi
    echo "Zipping $LNX64_CLI_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        tar czf "$LNX64_CLI_ZIP" --owner=0 --group=0 "$LNX64_CLI"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        tar czf "$LNX64_CLI_ZIP"  "$LNX64_CLI"
    fi
    FINALS+=("$LNX64_CLI_ZIP")
fi

# Linux arm
if [ -e "$LNX_ARM_CLI" ]; then
    if [ -e "$LNX_ARM_CLI_ZIP" ]; then
        echo "Removing old $LNX_ARM_CLI_ZIP"
        rm "$LNX_ARM_CLI_ZIP"
    fi
    echo "Zipping $LNX_ARM_CLI_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        tar czf "$LNX_ARM_CLI_ZIP" --owner=0 --group=0 "$LNX_ARM_CLI"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        tar czf "$LNX_ARM_CLI_ZIP"  "$LNX_ARM_CLI"
    fi
    FINALS+=("$LNX_ARM_CLI_ZIP")
fi

popd >/dev/null

# Move to final release dir
mkdir -p "$FINAL_OUTPUT_DIR"
for var in "${FINALS[@]}"; do
    mv "${CLI_OUTPUT_DIR}/${var}" "$FINAL_OUTPUT_DIR"
done

popd >/dev/null
