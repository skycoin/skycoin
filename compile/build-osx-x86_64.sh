#!/usr/bin/env bash

CONFIG=osx-x86_64-config
. build-tools.sh

compile_app
create_nw_bin

# Create .app and zip it
echo "Building .app package"
APPTMP=${CACHEDIR}/.unzipped
CONTENTS=${APPTMP}/${APPNAME}.app/Contents
RESOURCES=${CONTENTS}/Resources

rm -rf "$APPTMP"
mkdir -p "$APPTMP"
unzip "${CACHEDIR}/${NWNAME}.${NWEXT}" -d "$APPTMP"
mv "${APPTMP}/node-webkit.app" "${APPTMP}/${APPNAME}.app"
cp "${CACHEDIR}/${BIN}.nw" "${RESOURCES}/app.nw"
cp ../LICENSE "$CONTENTS"
cp osx/Info.plist "$CONTENTS"
# TODO -- use our own skycoin.icns file
mv "${RESOURCES}/nw.icns" "${RESOURCES}/skycoin.icns"
#cp osx/skycoin.icns "$RESOURCES"
mkdir -p "$RELEASEDIR"
if [[ -d "${RELEASEDIR}/${APPNAME}.app" ]]; then
    rm -rf "${RELEASEDIR}/${APPNAME}.app"
fi
mv "${APPTMP}/${APPNAME}.app" "$RELEASEDIR"
mv "${BINDIR}/${APP}" "${RELEASEDIR}/${APPD}"
mv "${BINDIR}/${BIN}.nw" "${RELEASEDIR}/${APPNW}.nw"

echo "Created ${OS} ${ARCH} release ${RELEASEDIR}/${ZIPNAME}.zip"
