#!/usr/bin/env bash
set -e -o pipefail

# These values are also in package.json and must be equal

# Get skycoin build version from package.json
APP_VERSION=`grep version package.json | sed  's/[,\", ]//g'| awk '{split($0,a,":");print a[2]}'`


# package name
PKG_NAME=`grep name package.json | sed 's/[,\", ]//g' | awk '{split($0,s,":");print s[2]}'`

# product name
PDT_NAME=`grep productName package.json | sed 's/[,\", ]//g' | awk '{split($0,s,":");print s[2]}'`

ELN_VERSION="v1.4.13"
ELN_OUTPUT_BASE=".electron_output"
ELN_OUTPUT="${ELN_OUTPUT_BASE}/${ELN_VERSION}"


if [ -n "$1" ]; then
    GOX_OSARCH="$1"
else
    GOX_OSARCH="linux/amd64 linux/arm windows/amd64 windows/386 darwin/amd64"
fi

# GOX_OSARCH="linux/amd64 darwin/amd64"
# GOX_OSARCH="linux/amd64 linux/arm windows/amd64 windows/386 darwin/amd64"
# GOX_OSARCH="linux/amd64"
# GOX_OSARCH="darwin/amd64"
# GOX_OSARCH="windows/amd64"
# GOX_OSARCH="windows/386"
# GOX_OSARCH="linux/arm"


GOX_OUTPUT=".gox_output"

STL_OUTPUT=".standalone_output"

FINAL_OUTPUT="release"

VERSION_FILE="./skycoin/current-skycoin.json"

GUI_DIST_DIR="../src/gui/static/dist"  # Do not append "/" to this path

# Variable suffix guide:
# _APP -- name of the OS X app
# _ELN_PLT -- directory name created by electron for its build of this platform
# _ELN -- our name for electron/gui releases
# _ELN_ZIP -- our compressed name for electron/gui releases
# _STL -- our name for standalone/non-gui releases
# _STL_ZIP -- our compressed name for standalone/non-gui releases

if [[ $GOX_OSARCH == *"darwin/amd64"* ]]; then
    OSX64_APP="${PDT_NAME}.app"
    OSX64_ELN_PLT="darwin-x64"
    OSX64_ELN="${PKG_NAME}-${APP_VERSION}-gui-osx-darwin-x64"
    OSX64_ELN_ZIP="${OSX64_ELN}.zip"
    OSX64_STL="${PKG_NAME}-${APP_VERSION}-bin-osx-darwin-x64"
    OSX64_STL_ZIP="${OSX64_STL}.zip"
    OSX64_OUT="mac_x64"
fi

if [[ $GOX_OSARCH == *"linux/amd64"* ]]; then
    LNX64_ELN="${PKG_NAME}-${APP_VERSION}-gui-linux-x64"
    LNX64_ELN_PLT="linux-x64"
    LNX64_ELN_ZIP="${LNX64_ELN}.tar.gz"
    LNX64_STL="${PKG_NAME}-${APP_VERSION}-bin-linux-x64"
    LNX64_STL_ZIP="${LNX64_STL}.tar.gz"
    LNX64_OUT="linux_x64"
fi

if [[ $GOX_OSARCH == *"windows/amd64"* ]]; then
    WIN64_ELN="${PKG_NAME}-${APP_VERSION}-gui-win-x64"
    WIN64_ELN_PLT="win32-x64"
    WIN64_ELN_ZIP="${WIN64_ELN}.zip"
    WIN64_STL="${PKG_NAME}-${APP_VERSION}-bin-win-x64"
    WIN64_STL_ZIP="${WIN64_STL}.zip"
    WIN64_OUT="win_x64"
fi

if [[ $GOX_OSARCH == *"windows/386"* ]]; then
    WIN32_ELN="${PKG_NAME}-${APP_VERSION}-gui-win-x86"
    WIN32_ELN_PLT="win32-ia32"
    WIN32_ELN_ZIP="${WIN32_ELN}.zip"
    WIN32_STL="${PKG_NAME}-${APP_VERSION}-bin-win-x86"
    WIN32_STL_ZIP="${WIN32_STL}.zip"
    WIN32_OUT="win_ia32"
fi

if [[ $GOX_OSARCH == *"linux/arm"* ]]; then
    LNX_ARM_STL="${PKG_NAME}-${APP_VERSION}-bin-linux-arm"
    LNX_ARM_STL_ZIP="${LNX_ARM_STL}.tar.gz"
    LNX_ARM_OUT="linux_arm"
fi
