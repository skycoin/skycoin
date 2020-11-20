#!/usr/bin/env bash
set -e -o pipefail

# Copies gox-compiled binaries and compiled GUI assets
# into an electron package

GOX_OSARCH="$@"

. build-conf.sh "$GOX_OSARCH"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

OSX64="${ELN_OUTPUT_DIR}/${OSX64_ELN_PLT}"
WIN64="${ELN_OUTPUT_DIR}/${WIN64_ELN_PLT}"
WIN32="${ELN_OUTPUT_DIR}/${WIN32_ELN_PLT}"
LNX64="${ELN_OUTPUT_DIR}/${LNX64_ELN_PLT}"

OSX64_RES="${OSX64}/${OSX64_APP}/Contents/Resources/app"
WIN64_RES="${WIN64}/resources/app"
WIN32_RES="${WIN32}/resources/app"
LNX64_RES="${LNX64}/resources/app"

OSX64_SRC="${OSX64_RES}/src"
WIN64_SRC="${WIN64}/src"
WIN32_SRC="${WIN32}/src"
LNX64_SRC="${LNX64}/src"

# Capitalize OS X .app for convention
if [ -e "${OSX64}/${PKG_NAME}.app" ]; then
    mv "${OSX64}/${PKG_NAME}.app" "${OSX64}/${OSX64_APP}"
fi

DESTSRCS=()

function copy_if_exists {
    if [ -z "$1" -o -z "$2" -o -z "$3" -o -z "$4" ]; then
        echo "copy_if_exists requires 4 args"
        exit 1
    fi

    BIN="${GOX_GUI_OUTPUT_DIR}/${1}"
    DESTDIR="$2"
    DESTBIN="${DESTDIR}/${3}"
    DESTSRC="$4"

    if [  -f "$BIN" ]; then
        # Copy binary to electron app
        echo "Copying $BIN to $DESTBIN"
        # mkdir -p $DESTBIN
        cp "$BIN" "$DESTBIN"

        # Copy static resources to electron app
        echo "Copying $GUI_DIST_DIR to $DESTDIR"
        cp -R "$GUI_DIST_DIR" "$DESTDIR"

        # Copy changelog to app
        echo "Copying CHANGELOG.md to $DESTDIR"
        cp ../CHANGELOG.md "$DESTDIR"

        DESTSRCS+=("$DESTSRC")
    else
        echo "$BIN does not exist"
    fi
}

echo "Copying ${PKG_NAME} binaries"

copy_if_exists "${PKG_NAME}_darwin_amd64" "$OSX64_RES" "${PKG_NAME}" "$OSX64_SRC"
copy_if_exists "${PKG_NAME}_windows_amd64.exe" "$WIN64_RES" "${PKG_NAME}.exe" "$WIN64_SRC"
copy_if_exists "${PKG_NAME}_windows_386.exe" "$WIN32_RES" "${PKG_NAME}.exe" "$WIN32_SRC"
copy_if_exists "${PKG_NAME}_linux_amd64" "$LNX64_RES" "${PKG_NAME}" "$LNX64_SRC"

# # Copy the source for reference
# # tar it with filters, move it, then untar in order to do this
# echo "Copying source snapshot"

# ./package-source.sh "${DESTSRCS[@]}"
