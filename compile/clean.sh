#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd "$DIR" >/dev/null

if [[ -d ".cache/" ]]; then
	rm -r .cache/
	echo "Deleted .cache/"
fi
if [[ -d ".bin/" ]]; then
	rm -r .bin/
	echo "Deleted .bin/"
fi
if [[ -d "release/" ]]; then
	rm -r release/
	echo "Deleted release/"
fi

popd >/dev/null

echo "Cleaned build byproducts"
