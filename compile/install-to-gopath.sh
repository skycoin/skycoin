#!/usr/bin/env bash

if [[ -z "$GOPATH" ]]; then
    echo "GOPATH is not set"
    exit
fi

if [[ "$GOPATH" == *:* ]]; then
    echo "This script does not work when multiple paths are in GOPATH"
    echo "Your GOPATH is ${GOPATH}"
    exit
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"
LINKPATH="${GOPATH}/src/github.com/skycoin"
LINKTO="${LINKPATH}/sync"

if [[ -d "$LINKTO" ]]; then
    POINTSTO=`readlink "$LINKTO"`
    echo "skycoin already exists in GOPATH"
    echo "${LINKTO} -> ${POINTSTO}"
    exit
fi

mkdir -p "$LINKPATH"
ln -s "$DIR" "$LINKTO"

echo "Installed symlink"
echo "${LINKTO} -> ${DIR}"
