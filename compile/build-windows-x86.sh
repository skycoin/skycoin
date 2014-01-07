#!/usr/bin/env bash

CONFIG=windows-x86-config
. build-tools.sh

compile_app
create_nw_bin

echo "Building zip release"
RELEASEDIR=${RELEASEDIR}
APPTMP=${CACHEDIR}/.unzipped

# Unzip node-webkit resources
rm -rf "$APPTMP"
mkdir -p "${APPTMP}/${NWNAME}"
unzip "${CACHEDIR}/${NWNAME}.${NWEXT}" -d "${APPTMP}/${NWNAME}"
# Copy binary and dependencies
mkdir "${APPTMP}/${ZIPNAME}"
mv "${APPTMP}/${NWNAME}/icudt.dll" "${APPTMP}/${ZIPNAME}"
mv "${APPTMP}/${NWNAME}/nw.pak" "${APPTMP}/${ZIPNAME}"
cp "${BINDIR}/${BIN}" "${APPTMP}/${ZIPNAME}"
cp ../LICENSE "${APPTMP}/${ZIPNAME}"
# Create zip and place it in $RELEASEDIR
pushd "${APPTMP}/${ZIPNAME}"
zip "${ZIPNAME}.zip" *
popd
mkdir -p "$RELEASEDIR"
mv "${APPTMP}/${ZIPNAME}/${ZIPNAME}.zip" "$RELEASEDIR"
if [[ -d "${RELEASEDIR}/${ZIPNAME}" ]]; then
    rm -rf "${RELEASEDIR}/${ZIPNAME}"
fi
mv "${APPTMP}/${ZIPNAME}" "$RELEASEDIR"
mv "${BINDIR}/${APP}" "${RELEASEDIR}/${APPD}"
mv "${BINDIR}/${BIN}.nw" "${RELEASEDIR}/${APPNW}.nw"

echo "Created ${OS} ${ARCH} release ${RELEASEDIR}/${ZIPNAME}.zip"
