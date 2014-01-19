#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd "$DIR" >/dev/null

deps=($(grep -r "github.com" ../src/ ../main.go | cut -d ":" -f 2 | \
    sed -e "s/[ \t]*//g" | sed -e "s/\"//g" | \
    sed -e "s/import//" | sort | uniq))

for i in "${deps[@]}"
do
    go get -u -v "$i"
done

popd >/dev/null

# Compilation dependencies
go get -u -v github.com/lonnc/golang-nw
go get -u -v github.com/lonnc/golang-nw/cmd/golang-nw-pkg
