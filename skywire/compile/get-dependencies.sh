#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd "$DIR" >/dev/null
cat dependencies.txt | xargs go get -u -v
popd >/dev/null
