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
CMD="${CMD:-$PKG_NAME}"  # name of folder in ../cmd to build. defaults to PKG_NAME which is the name of the coin

OSARCH="$1"
OUTPUT_DIR="$2"
BIN_NAME="${3:-$CMD}"

CONFIG_MODE=${CONFIG_MODE:-}

if [ -z "$OSARCH" ]; then
    echo "$USAGE"
    exit 1
fi

case "$OUTPUT_DIR" in
"")
    ;;
*/)
    ;;
*)
    OUTPUT_DIR+="/"
    ;;
esac

pushd "$SCRIPTDIR" >/dev/null

if [ -n "$OUTPUT_DIR" ]; then
    mkdir -p "$OUTPUT_DIR"
fi

COMMIT=`git rev-parse HEAD`

CLI_IMPORT_PATH=`go list ../src/cli`

gox -osarch="$OSARCH" \
    -gcflags="-trimpath=${HOME}" \
    -asmflags="-trimpath=${HOME}" \
    -ldflags="-X main.Version=${APP_VERSION} -X main.Commit=${COMMIT} -X main.ConfigMode=${CONFIG_MODE} -X ${CLI_IMPORT_PATH}.Version=${APP_VERSION}" \
    -output="${OUTPUT_DIR}{{.Dir}}_{{.OS}}_{{.Arch}}" \
    "${CMDDIR}/${CMD}"

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
            OUT="${OUTPUT_DIR}${WIN32_OUT}"
            echo "mkdir $OUT"
            mkdir -p "$OUT"
            mv "${OUTPUT_DIR}${CMD}_${s[0]}_${s[1]}.exe" "${OUT}/${BIN_NAME}.exe"
        else
            OUT="${OUTPUT_DIR}${WIN64_OUT}"
            mkdir -p "${OUT}"
            mv "${OUTPUT_DIR}${CMD}_${s[0]}_${s[1]}.exe" "${OUT}/${BIN_NAME}.exe"
        fi
        ;;
    "darwin")
        OUT="${OUTPUT_DIR}${OSX64_OUT}"
        echo "mkdir ${OUT}"
        mkdir -p "${OUT}"
        mv "${OUTPUT_DIR}${CMD}_${s[0]}_${s[1]}" "${OUT}/${BIN_NAME}"
        ;;
    "linux")
        if [ "${s[1]}" = "amd64" ]; then
            OUT="${OUTPUT_DIR}${LNX64_OUT}"
            echo "mkdir ${OUT}"
            mkdir -p "${OUT}"
            mv "${OUTPUT_DIR}${CMD}_${s[0]}_${s[1]}" "${OUT}/${BIN_NAME}"
        elif [ "${s[1]}" = "arm" ]; then
            OUT="${OUTPUT_DIR}${LNX_ARM_OUT}"
            echo "mkdir ${OUT}"
            mkdir -p "${OUT}"
            mv "${OUTPUT_DIR}${CMD}_${s[0]}_${s[1]}" "${OUT}/${BIN_NAME}"
        fi
        ;;
    esac
done

popd >/dev/null
