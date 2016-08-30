#!/usr/bin/env bash

. build-conf.sh

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

rm -r "$ELN_OUTPUT_BASE"
rm -r "$STL_OUTPUT"
rm -r "$GOX_OUTPUT"
rm -r "$FINAL_OUTPUT"

# don't remove the electron cache by default, most of the time when we want
# to clean up build artifacts we don't want to clean this up, and downloading
# it again is slow
# rm -r .electron_cache

popd >/dev/null
