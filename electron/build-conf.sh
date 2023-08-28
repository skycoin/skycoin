#!/usr/bin/env bash
set -e -o pipefail

# These values are also in package.json and must be equal

# Get skycoin build version from package.json
APP_VERSION=`grep version package.json | sed  's/[,\", ]//g'| awk '{split($0,a,":");print a[2]}'`


# package name
PKG_NAME=`grep productName package.json | sed 's/[,\", ]//g' | awk '{split($0,s,":");print tolower(s[2])}'`

# product name
PDT_NAME=`grep productName package.json | sed 's/[,\", ]//g' | awk '{split($0,s,":");print s[2]}'`

ELN_VERSION="v1.4.13"
ELN_OUTPUT_BASE=".electron_output"
ELN_OUTPUT_DIR="${ELN_OUTPUT_BASE}/${ELN_VERSION}"

# whether do code signing for osx builds
CODE_SIGN=$CSC_IDENTITY_AUTO_DISCOVERY


if [ -n "$1" ]; then
    GOX_OSARCH="$@"
else
    GOX_OSARCH="linux/amd64 linux/arm windows/amd64 windows/386 darwin/amd64"
fi

GOX_OUTPUT_DIR=".gox_output"
GOX_GUI_OUTPUT_DIR="${GOX_OUTPUT_DIR}/gui"
GOX_DMN_OUTPUT_DIR="${GOX_OUTPUT_DIR}/daemon"
GOX_CLI_OUTPUT_DIR="${GOX_OUTPUT_DIR}/cli"
GOX_CLI_OUTPUT_NAME="${PKG_NAME}-cli"

STL_OUTPUT_DIR=".standalone_output"
DMN_OUTPUT_DIR=".daemon_output"
CLI_OUTPUT_DIR=".cli_output"

FINAL_OUTPUT_DIR="release"

GUI_DIST_DIR="../src/gui/static/dist"  # Do not append "/" to this path

# Variable suffix guide:
# _APP -- name of the OS X app
# _ELN_PLT -- directory name created by electron for its build of this platform
# _ELN -- our name for electron gui releases
# _ELN_ZIP -- our compressed name for electron gui releases
# _STL -- our name for standalone gui releases
# _STL_ZIP -- our compressed name for standalone gui releases
# _DMN -- our name for daemon releases
# _DMN_ZIP -- our compressed name for daemon releases
# _CLI -- our name for cli releases
# _CLI_ZIP -- our compressed name for cli releases

if [[ $GOX_OSARCH == *"darwin/amd64"* ]]; then
    OSX64_APP="${PDT_NAME}.app"
    OSX64_ELN_PLT="darwin-x64"
    OSX64_ELN="${PKG_NAME}-${APP_VERSION}-gui-electron-osx-darwin-x64"
    OSX64_ELN_ZIP="${OSX64_ELN}.zip"
    OSX64_STL="${PKG_NAME}-${APP_VERSION}-gui-standalone-osx-darwin-x64"
    OSX64_STL_ZIP="${OSX64_STL}.zip"
    OSX64_DMN="${PKG_NAME}-${APP_VERSION}-daemon-osx-darwin-x64"
    OSX64_DMN_ZIP="${OSX64_DMN}.zip"
    OSX64_CLI="${PKG_NAME}-${APP_VERSION}-cli-osx-darwin-x64"
    OSX64_CLI_ZIP="${OSX64_CLI}.zip"
    OSX64_OUT="mac_x64"
fi

if [[ $GOX_OSARCH == *"linux/amd64"* ]]; then
    LNX64_ELN="${PKG_NAME}-${APP_VERSION}-gui-electron-linux-x64"
    LNX64_ELN_PLT="linux-x64"
    LNX64_ELN_ZIP="${LNX64_ELN}.tar.gz"
    LNX64_STL="${PKG_NAME}-${APP_VERSION}-gui-standalone-linux-x64"
    LNX64_STL_ZIP="${LNX64_STL}.tar.gz"
    LNX64_DMN="${PKG_NAME}-${APP_VERSION}-daemon-linux-x64"
    LNX64_DMN_ZIP="${LNX64_DMN}.tar.gz"
    LNX64_CLI="${PKG_NAME}-${APP_VERSION}-cli-linux-x64"
    LNX64_CLI_ZIP="${LNX64_CLI}.tar.gz"
    LNX64_OUT="linux_x64"
fi

if [[ $GOX_OSARCH == *"windows/amd64"* ]]; then
    WIN64_ELN="${PKG_NAME}-${APP_VERSION}-gui-electron-win-x64"
    WIN64_ELN_PLT="win32-x64"
    WIN64_ELN_ZIP="${WIN64_ELN}.zip"
    WIN64_STL="${PKG_NAME}-${APP_VERSION}-gui-standalone-win-x64"
    WIN64_STL_ZIP="${WIN64_STL}.zip"
    WIN64_DMN="${PKG_NAME}-${APP_VERSION}-daemon-win-x64"
    WIN64_DMN_ZIP="${WIN64_DMN}.zip"
    WIN64_CLI="${PKG_NAME}-${APP_VERSION}-cli-win-x64"
    WIN64_CLI_ZIP="${WIN64_CLI}.zip"
    WIN64_OUT="win_x64"
fi

if [[ $GOX_OSARCH == *"windows/386"* ]]; then
    WIN32_ELN="${PKG_NAME}-${APP_VERSION}-gui-electron-win-x86"
    WIN32_ELN_PLT="win32-ia32"
    WIN32_ELN_ZIP="${WIN32_ELN}.zip"
    WIN32_STL="${PKG_NAME}-${APP_VERSION}-gui-standalone-win-x86"
    WIN32_STL_ZIP="${WIN32_STL}.zip"
    WIN32_DMN="${PKG_NAME}-${APP_VERSION}-daemon-win-x86"
    WIN32_DMN_ZIP="${WIN32_DMN}.zip"
    WIN32_CLI="${PKG_NAME}-${APP_VERSION}-cli-win-x86"
    WIN32_CLI_ZIP="${WIN32_CLI}.zip"
    WIN32_OUT="win_ia32"
fi

if [[ $GOX_OSARCH == *"linux/arm"* ]]; then
    LNX_ARM_STL="${PKG_NAME}-${APP_VERSION}-gui-standalone-linux-arm"
    LNX_ARM_STL_ZIP="${LNX_ARM_STL}.tar.gz"
    LNX_ARM_DMN="${PKG_NAME}-${APP_VERSION}-daemon-linux-arm"
    LNX_ARM_DMN_ZIP="${LNX_ARM_DMN}.tar.gz"
    LNX_ARM_CLI="${PKG_NAME}-${APP_VERSION}-cli-linux-arm"
    LNX_ARM_CLI_ZIP="${LNX_ARM_CLI}.tar.gz"
    LNX_ARM_OUT="linux_arm"
fi
