ENABLE_GOX=${ENABLE_GOX:-1}
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
RAWBIN="skycoin"
APPNAME="Skycoin"
BINDIR=".bin/"
CACHEDIR=".cache/"
PKGDIR="../"
RELEASEDIR="release/"
HTMLDIRSRC="../src/gui/static/dist/"

CGOCC="gcc"

. "$CONFIG"

if [[ -z "${OS}" ]]; then
    echo "Configuration variable missing: OS"
    exit 1
fi
if [[ -z "${ARCH}" ]]; then
    echo "Configuration variable missing: ARCH"
    exit 1
fi
if [[ -z "${APP}" ]]; then
    echo "Configuration variable missing: APP"
    exit 1
fi
if [[ -z "${APPD}" ]]; then
    echo "Configuration variable missing: APPD"
    exit 1
fi
if [[ -z "${BIN}" ]]; then
    echo "Configuration variable missing: BIN"
    exit 1
fi
if [[ -z "${ZIPNAME}" ]]; then
    echo "Configuration variable missing: ZIPNAME"
    exit 1
fi

# Scratch space directory for assembling the folder to be zipped
SCRATCHDIR="${CACHEDIR}/.unzipped/"
# Name of the folder to zip
DIRTOZIP="${SCRATCHDIR}/${ZIPNAME}"
HTMLDIRDST="${DIRTOZIP}/src/gui/static/"


# Compiles the core app with gox. Note: we can't really cross-compile due to
# cgo
function compile_app_go() {
    echo "Compiling with go"
    CC="$CGOCC" GOOS="$OS" GOARCH="$ARCH" go build \
        -o "${BINDIR}/${APP}" "${PKGDIR}/cmd/${RAWBIN}"
    if [[ $? != 0 ]]; then
        echo "go compilation failed"
        exit 1
    fi
}

function compile_app_gox() {
	echo "Compiling with gox"
	gox -arch="$ARCH" -os="$OS" -output="${BINDIR}/${APP}" "${PKGDIR}/cmd/${RAWBIN}"
	if [[ $? -ne 0 ]]; then
		echo "gox compilation of $OS $ARCH failed"
		exit 1
	fi
}

function compile_app() {
	if [[ $ENABLE_GOX -eq 1 ]]; then
		compile_app_gox
	else
		compile_app_go
	fi
}

# Bundles files for a linux release
function create_linux_package() {
    echo "Creating linux release bundle"
    # Reset the scratch space
    rm -rf "$SCRATCHDIR"
    mkdir -p "$SCRATCHDIR"

    # Setup the scratch space with the target folder
    mkdir "$DIRTOZIP"

    # Copy GUI static resources
    echo
    mkdir -p "$HTMLDIRDST"
    cp -R "$HTMLDIRSRC" "$HTMLDIRDST"

    # Make executable and copy the client binary and README
    chmod +x "${BINDIR}/${APP}"
    cp "${BINDIR}/${APP}" "${DIRTOZIP}/${BIN}"
    cp linux/README "$DIRTOZIP"

    # Zip the target folder. Need to cd to the target directory so that the
    # tarball has the correct folder structure
    pushd "$SCRATCHDIR" >/dev/null
    tar cvzf "${ZIPNAME}.tar.gz" "${ZIPNAME}"
    popd >/dev/null

    # Setup the release dir, final results stored here
    mkdir -p "$RELEASEDIR"
    if [[ -d "${RELEASEDIR}/${ZIPNAME}" ]]; then
        rm -rf "${RELEASEDIR}/${ZIPNAME}"
    fi

    # Copy the .tar.gz, unzipped folder, and skycoind binary to the release dir
    mv "${DIRTOZIP}" "$RELEASEDIR"
    mv "${DIRTOZIP}.tar.gz" "$RELEASEDIR"
    ls "${BINDIR}"
    mv "${BINDIR}/${APP}" "${RELEASEDIR}/${APPD}"

    echo "Created ${OS} ${ARCH} release ${RELEASEDIR}/${ZIPNAME}.tar.gz"
}

function create_osx_package() {
    echo "Creating OSX release bundle"

    # TODO -- remove once the client binary is copied into the .app properly
    echo "Not fully implemented, aborting"
    exit 1

    CONTENTS=${SCRATCHDIR}/${APPNAME}.app/Contents
    RESOURCES=${CONTENTS}/Resources

    # Reset the scratch space
    rm -rf "$SCRATCHDIR"
    mkdir -p "$SCRATCHDIR"

    # Setup the target folder
    mkdir -p "$RESOURCES"

    # Copy static resources into the .app
    cp -R "$HTMLDIR" "$RESOURCES"
    cp osx/Info.plist "$CONTENTS"
    # TODO -- use our own skycoin.icns file
    #cp osx/skycoin.icns "$RESOURCES"

    # TODO -- package the binary properly for use in a .app here
    # cp $BIN ${APPNAME}.app/${SOMEWHERE}

    # zip the .app
    pushd "$SCRATCHDIR" >/dev/null
    zip -r "${ZIPNAME}.zip" "${APPNAME}.app"
    popd >/dev/null

    # Setup the release dir, final results stored here
    mkdir -p "$RELEASEDIR"
    if [[ -d "${RELEASEDIR}/${APPNAME}.app" ]]; then
        rm -rf "${RELEASEDIR}/${APPNAME}.app"
    fi

    # Copy the zipped .app, the .app and the skycoind binary to the release dir
    mv "${SCRATCHDIR}/${APPNAME}.app" "$RELEASEDIR"
    mv "${SCRATCHDIR}/${ZIPNAME}.zip" "$RELEASEDIR"
    mv "${BINDIR}/${APP}" "${RELEASEDIR}/${APPD}"

    echo "Created ${OS} ${ARCH} release ${RELEASEDIR}/${ZIPNAME}.zip"
}

function create_windows_package() {
    echo "Building windows release bundle"

    # Reset the scratch space
    rm -rf "$SCRATCHDIR"
    mkdir -p "$SCRATCHDIR"

    # Setup the scratch space with the target folder
    mkdir "$DIRTOZIP"

    # Copy GUI static resources
    cp -R "$HTMLDIR" "$DIRTOZIP"

    # Copy client binary
    ls ${BINDIR}
    cp "${BINDIR}/${APP}.exe" "${DIRTOZIP}/${BIN}.exe"

    # Create zip
    pushd "$SCRATCHDIR" >/dev/null
    zip -r "${ZIPNAME}.zip" "$ZIPNAME"
    popd >/dev/null

    # Setup release dir
    mkdir -p "$RELEASEDIR"
    if [[ -d "${RELEASEDIR}/${ZIPNAME}" ]]; then
        rm -rf "${RELEASEDIR}/${ZIPNAME}"
    fi

    # Copy the zip, unzipped folder and skycoind binary to the release dir
    mv "${DIRTOZIP}" "$RELEASEDIR"
    mv "${DIRTOZIP}.zip" "$RELEASEDIR"
    mv "${BINDIR}/${APP}.exe" "${RELEASEDIR}/${APPD}.exe"

    echo "Created ${OS} ${ARCH} release ${RELEASEDIR}/${ZIPNAME}.zip"
}

function do_linux() {
    pushd "$SCRIPTDIR" >/dev/null
    compile_app
    create_linux_package
    popd >/dev/null
}

function do_osx() {
    pushd "$SCRIPTDIR" >/dev/null
    compile_app
    create_osx_package
    popd >/dev/null
}

function do_windows() {
    pushd "$SCRIPTDIR" >/dev/null
    compile_app
    create_windows_package
    popd >/dev/null
}
