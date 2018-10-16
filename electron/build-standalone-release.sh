#!/usr/bin/env bash
set -e -o pipefail

GOX_OSARCH="$@"

. build-conf.sh "$GOX_OSARCH"

SKIP_COMPILATION=${SKIP_COMPILATION:-0}

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

if [ $SKIP_COMPILATION -ne 1 ]; then
    CONFIG_MODE=STANDALONE_CLIENT ./gox.sh "$GOX_OSARCH" "$GOX_GUI_OUTPUT_DIR"
fi

echo
echo "==========================="
echo "Packaging standalone release"
./package-standalone-release.sh "$GOX_OSARCH"

echo "------------------------------"
echo "Compressing standalone release"
./compress-standalone-release.sh "$GOX_OSARCH"

popd >/dev/null
