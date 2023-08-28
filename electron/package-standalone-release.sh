#!/usr/bin/env bash
set -e -o pipefail

# Builds the release without electron

GOX_OSARCH="$@"

echo "In package standalone release: $GOX_OSARCH"

. build-conf.sh "$GOX_OSARCH"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

DESTSRCS=()

function copy_if_exists {
    if [ -z "$1" -o -z "$2" -o -z "$3" ]; then
        echo "copy_if_exists requires 3 args"
        exit 1
    fi

    BIN="${GOX_GUI_OUTPUT_DIR}/${1}"
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

        # Copy changelog to app
        echo "Copying CHANGELOG.md to $DESTDIR"
        cp ../CHANGELOG.md "$DESTDIR"

        echo "Adding $DESTSRC to package-source.sh list"
        DESTSRCS+=("$DESTSRC")
    else
        echo "$BIN does not exsit"
    fi
}

function codesign_if_exists {
     if [ -z "$1" ]; then
        echo "codesign_if_exists requires binary path"
        exit 1
    fi

    BIN="${GOX_GUI_OUTPUT_DIR}/${1}"

    if [ -f "$BIN" ]; then
        if $CODE_SIGN; then
            echo "signing standalone binary"
            codesign --force --sign "Developer ID Application: yunfei mao" "${BIN}"
        fi
    else
        echo "$BIN does not exsit"
    fi
}

echo "Copying ${PKG_NAME} binaries"

# OS X
if [ ! -z "$OSX64_STL" ]; then
    OSX64="${STL_OUTPUT_DIR}/${OSX64_STL}"
    OSX64_SRC="${OSX64}/src"
    codesign_if_exists "${OSX64_OUT}/${PKG_NAME}"
    copy_if_exists "${OSX64_OUT}/${PKG_NAME}" "$OSX64" "$OSX64_SRC"
fi

# Linux amd64
if [ ! -z "$LNX64_STL" ]; then
    LNX64="${STL_OUTPUT_DIR}/${LNX64_STL}"
    LNX64_SRC="${LNX64}/src"
    copy_if_exists "${LNX64_OUT}/${PKG_NAME}" "$LNX64" "$LNX64_SRC"
fi

# Linux arm
if [ ! -z "$LNX_ARM_STL" ]; then
    LNX_ARM="${STL_OUTPUT_DIR}/${LNX_ARM_STL}"
    LNX_ARM_SRC="${LNX_ARM}/src"
    copy_if_exists "${LNX_ARM_OUT}/${PKG_NAME}" "$LNX_ARM" "$LNX_ARM_SRC"
fi

# Windows amd64
if [ ! -z "$WIN64_STL" ]; then
    WIN64="${STL_OUTPUT_DIR}/${WIN64_STL}"
    WIN64_SRC="${WIN64}/src"
    copy_if_exists "${WIN64_OUT}/${PKG_NAME}.exe" "$WIN64" "$WIN64_SRC"
fi

# Windows 386
if [ ! -z "$WIN32_STL" ]; then
    WIN32="${STL_OUTPUT_DIR}/${WIN32_STL}"
    WIN32_SRC="${WIN32}/src"
    copy_if_exists "${WIN32_OUT}/${PKG_NAME}.exe" "$WIN32" "$WIN32_SRC"
fi

# # Copy the source for reference
# # tar it with filters, move it, then untar in order to do this
# echo "Copying source snapshot"

# ./package-source.sh "${DESTSRCS[@]}"
