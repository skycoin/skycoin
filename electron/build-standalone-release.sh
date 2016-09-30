#!/usr/bin/env bash
set -e -o pipefail

. build-conf.sh

SKIP_COMPILATION=${SKIP_COMPILATION:-0}

if [ -n "$1" ]; then
    GOX_OSARCH="$1"
fi

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

if [ $SKIP_COMPILATION -ne 1 ]; then
    ./gox.sh "$GOX_OSARCH" "$GOX_OUTPUT"
fi

echo "----------------------------"
echo "Packaging standalone release"
./package-standalone-release.sh

echo "------------------------------"
echo "Compressing standalone release"
./compress-standalone-release.sh

popd >/dev/null
