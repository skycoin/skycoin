#!/usr/bin/env bash
set -e -o pipefail

SRC_TAR="tmp-src-snapshot.tar"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "${SCRIPTDIR}"

tar cvf "${SRC_TAR}" --owner=0 --group=0 --exclude=electron \
    --exclude=node_modules --exclude=_deprecated --exclude='.[^/\.]*' \
    "../src" "../cmd" "../run.sh" "../test.sh" "../GLOCKFILE" "../README.md" \
    >/dev/null

popd >/dev/null

function copy_source {
    echo "Copying source tree to $1"
    mkdir -p "$1"
    cp "${SRC_TAR}" "$1"
    pushd "$1"
    tar xvf "${SRC_TAR}" >/dev/null
    rm "${SRC_TAR}"
    popd >/dev/null
}

for var in "$@"; do
    copy_source "$var"
done

rm "${SRC_TAR}"
