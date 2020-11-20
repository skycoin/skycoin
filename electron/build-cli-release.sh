#!/usr/bin/env bash
set -e -o pipefail

GOX_OSARCH="$@"

. build-conf.sh "$GOX_OSARCH"

SKIP_COMPILATION=${SKIP_COMPILATION:-0}

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

if [ $SKIP_COMPILATION -ne 1 ]; then
    CMD="skycoin-cli" CONFIG_MODE= ./gox.sh "$GOX_OSARCH" "$GOX_CLI_OUTPUT_DIR" "$GOX_CLI_OUTPUT_NAME"
fi

echo
echo "==========================="
echo "Packaging cli release"
./package-cli-release.sh "$GOX_OSARCH"

echo "------------------------------"
echo "Compressing cli release"
./compress-cli-release.sh "$GOX_OSARCH"

popd >/dev/null
