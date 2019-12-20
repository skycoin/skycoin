#!/usr/bin/env bash
set -e -o pipefail

# Builds an entire electron-based GUI for release

# Implemented architectures:
#       darwin/amd64
#       windows/amd64
#       windows/386
#       linux/amd64
#
# By default builds all architectures.
# A single arch can be built by specifying it using gox's arch names

GOX_OSARCH="$@"

. build-conf.sh "$GOX_OSARCH"

SKIP_COMPILATION=${SKIP_COMPILATION:-0}

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

if [ $SKIP_COMPILATION -ne 1 ]; then
    CONFIG_MODE=STANDALONE_CLIENT ./gox.sh "$GOX_OSARCH" "$GOX_GUI_OUTPUT_DIR"
fi

if [ -e "$ELN_OUTPUT_DIR" ]; then
    rm -r "$ELN_OUTPUT_DIR"
fi

if [ ! -z "$WIN64_ELN" ] && [ ! -z "$WIN32_ELN" ]; then
    npm run dist-win
elif [ ! -z "$WIN64_ELN" ]; then
    npm run dist-win64
elif [ ! -z "$WIN32_ELN" ]; then
    npm run dist-win32
fi

if [ ! -z "$LNX64_ELN" ]; then
    npm run dist-linux
fi

if [ ! -z "$OSX64_ELN" ]; then
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "run dist-mac"
        npm run dist-mac
    else
        echo "Can not run build script in $OSTYPE"
    fi
fi

pushd "$FINAL_OUTPUT_DIR" >/dev/null
if [ -e "mac" ]; then
    pushd "mac" >/dev/null
    if [ -e "${PDT_NAME}-${APP_VERSION}.dmg" ]; then
        mv "${PDT_NAME}-${APP_VERSION}.dmg" "../${PKG_NAME}-${APP_VERSION}-gui-electron-osx-x64.dmg"
    elif [ -e "${PDT_NAME}.app" ]; then
        rm -rf "${PDT_NAME}.app"
    fi
    popd >/dev/null
    rm -rf "mac"
fi

IMG="${PKG_NAME}-${APP_VERSION}-x86_64.AppImage"
DEST_IMG="${PKG_NAME}-${APP_VERSION}-gui-electron-linux-x64.AppImage"
DEST_IMG_ZIP="${DEST_IMG}.tar.gz"
if [ -e $IMG ]; then
    mv "$IMG" "$DEST_IMG"
    chmod +x "$DEST_IMG"

    if [ -e $DEST_IMG_ZIP ]; then
        echo  "Removing old $DEST_IMG_ZIP"
        rm "$DEST_IMG_ZIP"
    fi

    echo "Zipping $DEST_IMG_ZIP"
    if [[ "$OSTYPE" == "linux"* ]]; then
        tar czf "$DEST_IMG_ZIP" --owner=0 --group=0 "$DEST_IMG"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        tar czf "$DEST_IMG_ZIP" "$DEST_IMG"
    fi
    rm "$DEST_IMG"
fi

EXE="${PDT_NAME} Setup ${APP_VERSION}.exe"
if [ -e "$EXE" ]; then
    if [ ! -z $WIN32_ELN ] && [ ! -z $WIN64_ELN ]; then
        mv "$EXE" "${PKG_NAME}-${APP_VERSION}-gui-electron-win-setup.exe"
    elif [ ! -z $WIN32_ELN ]; then
        mv "$EXE" "${WIN32_ELN}.exe"
    elif [ ! -z $WIN64_ELN ]; then
        mv "$EXE" "${WIN64_ELN}.exe"
    fi
fi

# rename dmg file name
DMG="${PKG_NAME}-${APP_VERSION}.dmg"
if [ -e "$DMG" ]; then
    mv "$DMG" "${PKG_NAME}-${APP_VERSION}-gui-electron-osx.dmg"
fi

# delete app zip file
MZIP="${PKG_NAME}-${APP_VERSION}-mac.zip"
if [ -e "$MZIP" ]; then
    rm "$MZIP"
fi

# delete github and latest-mac.yml
if [ -d "github" ]; then rm -rf github ;fi
if [ -e "latest-mac.yml" ]; then rm latest-mac.yml ;fi

# clean unpacked folders
rm -rf *-unpacked

# delete blockmap and electron-builder.yaml
rm -f *.blockmap
rm -f *.yaml

popd >/dev/null

popd >/dev/null
