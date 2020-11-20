#!/usr/bin/env bash
set -e -o pipefail

. build-conf.sh

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

function rmnofail {
    for dir in "$@"; do
        if [ -e "$dir" ]; then
            rm -r "$dir"
            echo "removed $dir"
        fi
    done
}

rmnofail "$ELN_OUTPUT_BASE" "$STL_OUTPUT_DIR" "$GOX_OUTPUT_DIR" "$FINAL_OUTPUT_DIR"

# don't remove the electron cache by default, most of the time when we want
# to clean up build artifacts we don't want to clean this up, and downloading
# it again is slow
# rm -r .electron_cache

popd >/dev/null
