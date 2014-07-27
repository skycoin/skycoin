SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
RAWBIN=$1
APPNW=skycoin
APPNAME="Skycoin"
NWVERSION="v0.8.3"
BINDIR="./.bin"
CACHEDIR="./.cache"
PKGDIR=../../
RELEASEDIR="./release"
HTMLDIR="../../src/gui/static"
LICENSE=../../LICENSE

. "$CONFIG"

if [[ -z "$NWARCH" ]]; then
    NWARCH="$ARCH"
fi

# Create the node-webkit wrapped executable
function create_nw_bin() {
    echo "Building release files"
    golang-nw-pkg -app="${BINDIR}/${APP}" -name="$APPNAME" -bin="$BIN" \
        -os=$OS -arch=$NWARCH -version="$NWVERSION" -toolbar=false \
        -cacheDir="$CACHEDIR" -binDir="$BINDIR" -includesDir="$HTMLDIR"
    if [[ $? != 0 ]]; then
        echo "nw-pkg build failed"
        exit 1
    fi
}

# Compiles the core app with gox. Note: we can't really cross-compile due to
# cgo
function compile_app() {
    echo "Compiling with go"
    #gox -osarch="${OS}/${ARCH}" -output="${BINDIR}/${APP}" "$PKGDIR"
    CC="$CGOCC" GOOS="$OS" GOARCH="$ARCH" go build \
        -o "${BINDIR}/${APP}" "${PKGDIR}/cmd/${RAWBIN}"
    if [[ $? != 0 ]]; then
        echo "go compilation failed"
        exit 1
    fi
}

# Bundles files for a linux release
function create_linux_package() {
    echo "Creating linux release bundle"
    # Unzip node-webkit resources
    APPTMP=${CACHEDIR}/.unzipped
    TMPZIP=${APPTMP}/${ZIPNAME}
    rm -rf "$APPTMP"
    mkdir -p "${APPTMP}/${NWNAME}"
    tar xvzf "${CACHEDIR}/${NWNAME}.${NWEXT}" -C "$APPTMP"
    # Copy binary and dependencies
    mkdir "$TMPZIP"
    mv "${APPTMP}/${NWNAME}/nw.pak" "$TMPZIP"
    cp -R "$HTMLDIR" "$TMPZIP"
    chmod +x "${BINDIR}/${BIN}"
    cp "${BINDIR}/${BIN}" "${TMPZIP}/${BIN}"
    ln -s "$LIBUDEV" "${TMPZIP}/libudev.so.0"
    cp linux/README "$TMPZIP"
    cp "$LICENSE" "$TMPZIP"
    mkdir -p "$RELEASEDIR"
    if [[ -d "${RELEASEDIR}/${ZIPNAME}" ]]; then
        rm -rf "${RELEASEDIR}/${ZIPNAME}"
    fi
    mv "${APPTMP}/${ZIPNAME}" "$RELEASEDIR"
    mv "${BINDIR}/${APP}" "${RELEASEDIR}/${APPD}"
    mv "${BINDIR}/${BIN}.nw" "${RELEASEDIR}/${APPNW}.nw"
    echo "Created ${OS} ${ARCH} release ${RELEASEDIR}/${ZIPNAME}.tar.gz"
}

function create_osx_package() {
    echo "Creating OSX release bundle"
    APPTMP=${CACHEDIR}/.unzipped
    CONTENTS=${APPTMP}/${APPNAME}.app/Contents
    RESOURCES=${CONTENTS}/Resources

    rm -rf "$APPTMP"
    mkdir -p "$APPTMP"
    unzip "${CACHEDIR}/${NWNAME}.${NWEXT}" -d "$APPTMP"
    mv "${APPTMP}/node-webkit.app" "${APPTMP}/${APPNAME}.app"
    mv "${CONTENTS}/MacOS/node-webkit" "${CONTENTS}/MacOS/${ZIPNAME}"
    cp "${BINDIR}/${BIN}.nw" "${RESOURCES}/app.nw"
    cp -R "$HTMLDIR" "$RESOURCES"
    cp "$LICENSE" "$CONTENTS"
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
}

function create_windows_package() {
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
    cp -R "$HTMLDIR" "$RESOURCES"
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
}

function do_linux() {
    pushd "$SCRIPTDIR" >/dev/null
    compile_app
    create_nw_bin
    create_linux_package
    popd >/dev/null
}

function do_osx() {
    pushd "$SCRIPTDIR" >/dev/null
    compile_app
    create_nw_bin
    create_osx_package
    popd >/dev/null
}

function do_windows() {
    pushd "$SCRIPTDIR" >/dev/null
    compile_app
    create_nw_bin
    create_windows_package
    popd >/dev/null
}
