#!/usr/bin/env bash
set -e -o pipefail

# Copies gox-compiled skycoin binaries and compiled GUI assets
# into an electron package

. build-conf.sh

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

OSX64="${ELN_OUTPUT}/${OSX64_ELN_PLT}"
WIN64="${ELN_OUTPUT}/${WIN64_ELN_PLT}"
LNX64="${ELN_OUTPUT}/${LNX64_ELN_PLT}"

OSX64_RES="${OSX64}/${OSX64_APP}/Contents/Resources/app"
WIN64_RES="${WIN64}/resources/app"
LNX64_RES="${LNX64}/resources/app"

OSX64_SRC="${OSX64_RES}/src"
WIN64_SRC="${WIN64}/src"
LNX64_SRC="${LNX64}/src"

# Capitalize OS X .app for convention
if [ -e "${OSX64}/skycoin.app" ]; then
    mv "${OSX64}/skycoin.app" "${OSX64}/${OSX64_APP}"
fi

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

    if [  -f "$BIN" ]; then
        # Copy binary to electron app
        echo "Copying $BIN to $DESTBIN"
        cp "$BIN" "$DESTBIN"

        # Copy static resources to electron app
        echo "Copying $GUI_DIST_DIR to $DESTDIR"
        cp -R "$GUI_DIST_DIR" "$DESTDIR"

        DESTSRCS+=("$DESTSRC")
    fi
}

echo "Copying skycoin binaries"

copy_if_exists "skycoin_darwin_amd64" "$OSX64_RES" "skycoin" "$OSX64_SRC"
copy_if_exists "skycoin_windows_amd64.exe" "$WIN64_RES" "skycoin.exe" "$WIN64_SRC"
copy_if_exists "skycoin_linux_amd64" "$LNX64_RES" "skycoin" "$LNX64_SRC"

# Copy the source for reference
# tar it with filters, move it, then untar in order to do this
echo "Copying source snapshot"

./package-source.sh "${DESTSRCS[@]}"
