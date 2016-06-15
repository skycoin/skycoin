#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

rm -r .electron_output
rm -r .electron_cache
rm -r .gox_output
rm -r release

popd >/dev/null
