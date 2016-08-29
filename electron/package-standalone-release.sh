#!/usr/bin/env bash

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

echo "Creating target directories in $STL_OUTPUT"

# create target directories
for var in "${OSX64}" "${WIN64}" "${LNX64}"; do
    if [ -e "${var}" ]; then
        rm -r "${var}"
    fi
    mkdir -p "${var}"
done

echo "Copying skycoin binaries"

# copy binaries
cp "${GOX_OUTPUT}/skycoin_darwin_amd64" "${OSX64}/skycoin"
cp "${GOX_OUTPUT}/skycoin_windows_amd64.exe" "${WIN64}/skycoin.exe"
cp "${GOX_OUTPUT}/skycoin_linux_amd64" "${LNX64}/skycoin"

echo "Copying static resources"

# copy static resources
for var in "${OSX64}" "${WIN64}" "${LNX64}"; do
    cp -R "${GUI_DIST_DIR}" "${var}"
done

# Copy the source for reference
# tar it with filters, move it, then untar in order to do this
echo "Copying source snapshot"

./package-source.sh "${OSX64_SRC}" "${WIN64_SRC}" "${LNX64_SRC}"

