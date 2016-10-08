#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$DIR" >/dev/null

go run cmd/mesh/mesh.go --gui-dir="${DIR}/src/mesh/gui/static/" $@

popd >/dev/null
