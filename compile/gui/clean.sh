#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd "$DIR" >/dev/null

rm -rf .cache/
rm -rf .bin/
rm -rf release/

popd /dev/null

echo "Cleaned build byproducts"
