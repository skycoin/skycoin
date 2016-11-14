#!/usr/bin/env bash
set -e -o pipefail

# These values are also in gulpfile.js and package.json and must be equal
SKY_VERSION="0.8.1"
ELN_VERSION="v1.2.2"
ELN_OUTPUT_BASE=".electron_output"
ELN_OUTPUT="${ELN_OUTPUT_BASE}/${ELN_VERSION}"

GOX_OSARCH="linux/amd64 windows/amd64 darwin/amd64"
GOX_OUTPUT=".gox_output"

STL_OUTPUT=".standalone_output"

FINAL_OUTPUT="release"

GUI_DIST_DIR="../src/gui/static/dist"  # Do not append "/" to this path

# Variable suffix guide:
# _APP -- name of the OS X app
# _ELN_PLT -- directory name created by electron for its build of this platform
# _ELN -- our name for electron/gui releases
# _ELN_ZIP -- our compressed name for electron/gui releases
# _STL -- our name for standalone/non-gui releases
# _STL_ZIP -- our compressed name for standalone/non-gui releases

OSX64_APP="Skycoin.app"
OSX64_ELN_PLT="darwin-x64"
OSX64_ELN="skycoin-${SKY_VERSION}-gui-osx-darwin-x64"
OSX64_ELN_ZIP="${OSX64_ELN}.zip"
OSX64_STL="skycoin-${SKY_VERSION}-bin-osx-darwin-x64"
OSX64_STL_ZIP="${OSX64_STL}.zip"

WIN64_ELN="skycoin-${SKY_VERSION}-gui-win-x64"
WIN64_ELN_PLT="win32-x64"
WIN64_ELN_ZIP="${WIN64_ELN}.zip"
WIN64_STL="skycoin-${SKY_VERSION}-bin-win-x64"
WIN64_STL_ZIP="${WIN64_STL}.zip"

LNX64_ELN="skycoin-${SKY_VERSION}-gui-linux-x64"
LNX64_ELN_PLT="linux-x64"
LNX64_ELN_ZIP="${LNX64_ELN}.tar.gz"
LNX64_STL="skycoin-${SKY_VERSION}-bin-linux-x64"
LNX64_STL_ZIP="${LNX64_STL}.tar.gz"
