. "$CONFIG"

APPNW=skycoin
TAGS=$1

# Create the node-webkit wrapped executable
function create_nw_bin() {
    echo "Building release files"
    golang-nw-pkg -app="${BINDIR}/${APP}" -name="$APPNAME" -bin="$BIN" \
        -os=$OS -arch=$ARCH -version="$NWVERSION" -toolbar=false \
        -cacheDir="$CACHEDIR" -binDir="$BINDIR" -includesDir=../static
    if [[ $? != 0 ]]; then
        echo "nw-pkg build failed"
        exit 1
    fi
}

# Compiles the core app with gox. Note: we can't really cross-compile due to
# cgo
function compile_app() {
    echo "Comiling with go"
    #gox -osarch="${OS}/${ARCH}" -output="${BINDIR}/${APP}" "$PKGDIR"
    CC="$CGOCC" GOOS="$OS" GOARCH="$ARCH" go build -tags "$TAGS" \
        -o "${BINDIR}/${APP}" "$PKGDIR"
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
