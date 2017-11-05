#!/usr/bin/env bash
set -e -o pipefail

 if [ -n "$1" ]; then
    GOX_OSARCH="$1"
fi

. build-conf.sh "$GOX_OSARCH"

SKIP_COMPILATION=${SKIP_COMPILATION:-0}

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

if [ $SKIP_COMPILATION -ne 1 ]; then
    ./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT"
fi

echo
echo "==========================="
echo "Stamping the release with proper version"
./version-control.sh

echo "----------------------------"
echo "Packaging standalone release"
./package-standalone-release.sh "$GOX_OSARCH"

echo "------------------------------"
echo "Compressing standalone release"
./compress-standalone-release.sh "$GOX_OSARCH"

popd >/dev/null
