. "$CONFIG"

APPNW=skycoin
RAWBIN=$1

if [[ -z "$NWARCH" ]]; then
    NWARCH="$ARCH"
fi

# Create the node-webkit wrapped executable
function create_nw_bin() {
    echo "Building release files"
    golang-nw-pkg -app="${BINDIR}/${APP}" -name="$APPNAME" -bin="$BIN" \
        -os=$OS -arch=$NWARCH -version="$NWVERSION" -toolbar=false \
        -cacheDir="$CACHEDIR" -binDir="$BINDIR" -includesDir=../static
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
    cp -R ../static "$TMPZIP"
    chmod +x "${BINDIR}/${BIN}"
    cp "${BINDIR}/${BIN}" "${TMPZIP}/${BIN}"
    ln -s "$LIBUDEV" "${TMPZIP}/libudev.so.0"
    cp linux/README "$TMPZIP"
    cp ../LICENSE "$TMPZIP"
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
}
