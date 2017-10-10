#!/usr/bin/env bash
set -e -o pipefail

. build-conf.sh

USAGE="./gox.sh \"osarch\" [output directory] [with builder]

Builds gox with the osarch string (see 'gox --help' for specifications)

Optionally specify an output directory for the build files. Will be created
if it does not exist.  Defaults to the working directory.

"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

CMDDIR="../cmd"  # relative to compile/electron/
CMD="${PKG_NAME}"

OSARCH="$1"
OUTPUT="$2"

WITH_BUILDER=$3
WITH_BUILDER=${WITH_BUILDER:-1}

if [ -z "$OSARCH" ]; then
    echo "$USAGE"
    exit 1
fi

case "$OUTPUT" in
"")
    ;;
*/)
    ;;
*)
    OUTPUT+="/"
    ;;
esac

pushd "$SCRIPTDIR" >/dev/null

if [ -n "$OUTPUT" ]; then
    mkdir -p "$OUTPUT"
fi

COMMIT=`git rev-parse HEAD`

gox -osarch="$OSARCH" \
    -ldflags="-X main.Version=${APP_VERSION} -X main.Commit=${COMMIT}" \
    -output="${OUTPUT}{{.Dir}}_{{.OS}}_{{.Arch}}" \
    "${CMDDIR}/${CMD}"

if [ "$WITH_BUILDER" = "1" ]; then
    # move the executable files into ${os}_${arch} folders, electron-builder will pack
    # the file into corresponding packages.

    platforms=$(echo $OSARCH | tr ";" "\n")

    for plt in $platforms
    do
        set -- "$plt"
        IFS="/"; declare -a s=($*)
        case "${s[0]}" in
        "windows")
            if [ "${s[1]}" = "386" ]; then
                OUT="${OUTPUT}${WIN32_OUT}"
                echo "mkdir $OUT"
                mkdir -p "$OUT"
                mv "${OUTPUT}${PKG_NAME}_${s[0]}_${s[1]}.exe" "${OUT}/${PKG_NAME}.exe"
            else
                OUT="${OUTPUT}${WIN64_OUT}"
                mkdir -p "${OUT}"
                mv "${OUTPUT}${PKG_NAME}_${s[0]}_${s[1]}.exe" "${OUT}/${PKG_NAME}.exe"
            fi
            ;;
        "darwin")
            OUT="${OUTPUT}${OSX64_OUT}"
            echo "mkdir ${OUT}"
            mkdir -p "${OUT}"
            mv "${OUTPUT}${PKG_NAME}_${s[0]}_${s[1]}" "${OUT}/${PKG_NAME}"
            ;;
        "linux")
            if [ "${s[1]}" = "amd64" ]; then
                OUT="${OUTPUT}${LNX64_OUT}"
                echo "mkdir ${OUT}"
                mkdir -p "${OUT}"
                mv "${OUTPUT}${PKG_NAME}_${s[0]}_${s[1]}" "${OUT}/${PKG_NAME}"
            elif [ "${s[1]}" = "arm" ]; then
                OUT="${OUTPUT}${LNX_ARM_OUT}"
                echo "mkdir ${OUT}"
                mkdir -p "${OUT}"
                mv "${OUTPUT}${PKG_NAME}_${s[0]}_${s[1]}" "${OUT}/${PKG_NAME}"
            fi
            ;;
        esac
    done
fi

popd >/dev/null
