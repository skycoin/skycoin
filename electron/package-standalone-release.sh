#!/usr/bin/env bash
set -e -o pipefail

# Builds the release without electron

if [ -n "$1" ]; then
    GOX_OSARCH="$2"
fi

echo "In package standalone release: $GOX_OSARCH"

. build-conf.sh "$GOX_OSARCH"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

OSX64="${STL_OUTPUT}/${OSX64_STL}"
WIN64="${STL_OUTPUT}/${WIN64_STL}"
WIN32="${STL_OUTPUT}/${WIN32_STL}"
LNX64="${STL_OUTPUT}/${LNX64_STL}"
LNX_ARM="${STL_OUTPUT}/${LNX_ARM_STL}"

OSX64_SRC="${OSX64}/src"
WIN64_SRC="${WIN64}/src"
WIN32_SRC="${WIN32}/src"
LNX64_SRC="${LNX64}/src"
LNX_ARM_SRC="${LNX_ARM}/src"

DESTSRCS=()

function copy_if_exists {
    if [ -z "$1" -o -z "$2" -o -z "$3" ]; then
        echo "copy_if_exists requires 3 args"
        exit 1
    fi

    BIN="${GOX_OUTPUT}/${1}"
    DESTDIR="$2"
    DESTSRC="$3"

    if [ -f "$BIN" ]; then
        if [ -e "$DESTDIR" ]; then
            rm -r "$DESTDIR"
        fi
        mkdir -p "$DESTDIR"

        # Copy binary to electron app
        echo "Copying $BIN to $DESTDIR"
        cp "$BIN" "$DESTDIR"

        # Copy static resources to electron app
        echo "Copying $GUI_DIST_DIR to ${DESTDIR}/src/gui/static"
        mkdir -p "${DESTDIR}/src/gui/static"
        cp -R "$GUI_DIST_DIR" "${DESTDIR}/src/gui/static"

        echo "Adding $DESTSRC to package-source.sh list"
        DESTSRCS+=("$DESTSRC")
    else
        echo "$BIN does not exsit"
    fi
}

echo "Copying ${PKG_NAME} binaries"

# copy binaries
copy_if_exists "${OSX64_OUT}/${PKG_NAME}" "$OSX64" "$OSX64_SRC"
copy_if_exists "${WIN64_OUT}/${PKG_NAME}.exe" "$WIN64" "$WIN64_SRC"
copy_if_exists "${WIN32_OUT}/${PKG_NAME}.exe" "$WIN32" "$WIN32_SRC"
copy_if_exists "${LNX64_OUT}/${PKG_NAME}" "$LNX64" "$LNX64_SRC"
copy_if_exists "${LNX_ARM_OUT}/${PKG_NAME}" "$LNX_ARM" "$LNX_ARM_SRC"

# Copy the source for reference
# tar it with filters, move it, then untar in order to do this
echo "Copying source snapshot"

./package-source.sh "${DESTSRCS[@]}"

