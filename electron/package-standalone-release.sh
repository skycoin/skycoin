#!/usr/bin/env bash
set -e -o pipefail

# Builds the release without electron

. build-conf.sh

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

OSX64="${STL_OUTPUT}/${OSX64_STL}"
WIN64="${STL_OUTPUT}/${WIN64_STL}"
LNX64="${STL_OUTPUT}/${LNX64_STL}"

OSX64_SRC="${OSX64}/src"
WIN64_SRC="${WIN64}/src"
LNX64_SRC="${LNX64}/src"

DESTSRCS=()

function copy_if_exists {
    if [ -z "$1" -o -z "$2" -o -z "$3" -o -z "$4" ]; then
        echo "copy_if_exists requires 4 args"
        exit 1
    fi

    BIN="${GOX_OUTPUT}/${1}"
    DESTDIR="$2"
    DESTBIN="${DESTDIR}/${3}"
    DESTSRC="$4"

    if [ -f "$BIN" ]; then
        if [ -e "$DESTDIR" ]; then
            rm -r "$DESTDIR"
        fi
        mkdir -p "$DESTDIR"

        # Copy binary to electron app
        echo "Copying $BIN to $DESTBIN"
        cp "$BIN" "$DESTBIN"

        # Copy static resources to electron app
        echo "Copying $GUI_DIST_DIR to $DESTDIR"
        cp -R "$GUI_DIST_DIR" "$DESTDIR"

        echo "Adding $DESTSRC to package-source.sh list"
        DESTSRCS+=("$DESTSRC")
    fi
}

echo "Copying skycoin binaries"

# copy binaries
copy_if_exists "skycoin_darwin_amd64" "$OSX64" "skycoin" "$OSX64_SRC"
copy_if_exists "skycoin_windows_amd64.exe" "$WIN64" "skycoin.exe" "$WIN64_SRC"
copy_if_exists "skycoin_linux_amd64" "$LNX64" "skycoin" "$LNX64_SRC"

# Copy the source for reference
# tar it with filters, move it, then untar in order to do this
echo "Copying source snapshot"

./package-source.sh "${DESTSRCS[@]}"

