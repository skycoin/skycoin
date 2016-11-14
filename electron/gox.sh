#!/usr/bin/env bash
set -e -o pipefail

USAGE="./gox.sh \"osarch\" [output directory]

Builds gox with the osarch string (see 'gox --help' for specifications)

Optionally specify an output directory for the build files. Will be created
if it does not exist.  Defaults to the working directory.
"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

CMDDIR="../cmd"  # relative to compile/electron/
CMD="skycoin"

OSARCH="$1"
OUTPUT="$2"

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

gox -osarch="$OSARCH" \
    -output="${OUTPUT}{{.Dir}}_{{.OS}}_{{.Arch}}" \
    "${CMDDIR}/${CMD}"

popd >/dev/null
