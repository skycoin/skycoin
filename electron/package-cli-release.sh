#!/usr/bin/env bash
set -e -o pipefail

# Builds the cli release

GOX_OSARCH="$@"

echo "In package cli release: $GOX_OSARCH"

. build-conf.sh "$GOX_OSARCH"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

DESTSRCS=()

function copy_if_exists {
    if [ -z "$1" -o -z "$2" -o -z "$3" ]; then
        echo "copy_if_exists requires 3 args"
        exit 1
    fi

    BIN="${GOX_CLI_OUTPUT_DIR}/${1}"
    DESTDIR="$2"
    DESTSRC="$3"

    if [ -f "$BIN" ]; then
        if [ -e "$DESTDIR" ]; then
            rm -r "$DESTDIR"
        fi
        mkdir -p "$DESTDIR"

        # Copy binary to app
        echo "Copying $BIN to $DESTDIR"
        cp "$BIN" "$DESTDIR"

        # Copy changelog to app
        echo "Copying CHANGELOG.md to $DESTDIR"
        cp ../CHANGELOG.md "$DESTDIR"

        # Copy cmd/skycoin-cli/README.md to app
        echo "Copying cmd/skycoin-cli/README.md to $DESTDIR"
        cp ../cmd/skycoin-cli/README.md "$DESTDIR"

        echo "Adding $DESTSRC to package-source.sh list"
        DESTSRCS+=("$DESTSRC")
    else
        echo "$BIN does not exsit"
    fi
}

echo "Copying ${GOX_CLI_OUTPUT_NAME} binaries"

# OS X
if [ ! -z "$OSX64_CLI" ]; then
    OSX64="${CLI_OUTPUT_DIR}/${OSX64_CLI}"
    OSX64_SRC="${OSX64}/src"
    copy_if_exists "${OSX64_OUT}/${GOX_CLI_OUTPUT_NAME}" "$OSX64" "$OSX64_SRC"
fi

# Linux amd64
if [ ! -z "$LNX64_CLI" ]; then
    LNX64="${CLI_OUTPUT_DIR}/${LNX64_CLI}"
    LNX64_SRC="${LNX64}/src"
    copy_if_exists "${LNX64_OUT}/${GOX_CLI_OUTPUT_NAME}" "$LNX64" "$LNX64_SRC"
fi

# Linux arm
if [ ! -z "$LNX_ARM_CLI" ]; then
    LNX_ARM="${CLI_OUTPUT_DIR}/${LNX_ARM_CLI}"
    LNX_ARM_SRC="${LNX_ARM}/src"
    copy_if_exists "${LNX_ARM_OUT}/${GOX_CLI_OUTPUT_NAME}" "$LNX_ARM" "$LNX_ARM_SRC"
fi

# Windows amd64
if [ ! -z "$WIN64_CLI" ]; then
    WIN64="${CLI_OUTPUT_DIR}/${WIN64_CLI}"
    WIN64_SRC="${WIN64}/src"
    copy_if_exists "${WIN64_OUT}/${GOX_CLI_OUTPUT_NAME}.exe" "$WIN64" "$WIN64_SRC"
fi

# Windows 386
if [ ! -z "$WIN32_CLI" ]; then
    WIN32="${CLI_OUTPUT_DIR}/${WIN32_CLI}"
    WIN32_SRC="${WIN32}/src"
    copy_if_exists "${WIN32_OUT}/${GOX_CLI_OUTPUT_NAME}.exe" "$WIN32" "$WIN32_SRC"
fi

# # Copy the source for reference
# # tar it with filters, move it, then untar in order to do this
# echo "Copying source snapshot"

# ./package-source.sh "${DESTSRCS[@]}"
