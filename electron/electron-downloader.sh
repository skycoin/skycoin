#!/usr/bin/env bash

set -e -o pipefail

. build-conf.sh

VER=$(echo "$ELN_VERSION" | cut -c2-20)

# download windows electron
electron-download --version="${VER}" --platform="win32" --arch=x64
electron-download --version="${VER}" --platform="win32" --arch=ia32

# download linux electron
electron-download --version="${VER}" --platform="linux" --arch=x64

# download mac osx electron
electron-download --version="${VER}" --platform="darwin"


