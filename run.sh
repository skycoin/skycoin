#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$DIR" >/dev/null
go run cmd/skycoindev/skycoindev.go $@
popd >/dev/null
