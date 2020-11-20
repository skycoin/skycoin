#!/usr/bin/env bash
set -e -o pipefail

SRC_TAR="tmp-src-snapshot.tar"

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "${SCRIPTDIR}"

if [[ "$OSTYPE" == "linux"* ]]; then
    tar -C .. -cvPf "${SRC_TAR}" --owner=0 --group=0 --exclude=electron \
        --exclude=node_modules --exclude=_deprecated --exclude='.*' \
        src cmd run-client.sh run-daemon.sh README.md INSTALLATION.md CHANGELOG.md INTEGRATION.md \
        >/dev/null
elif [[ "$OSTYPE" == "darwin"* ]]; then
    tar -C .. -cvf "${SRC_TAR}" --exclude=electron \
        --exclude=node_modules --exclude=_deprecated --exclude='.*' \
        src cmd run-client.sh run-daemon.sh README.md INSTALLATION.md CHANGELOG.md INTEGRATION.md \
        >/dev/null
elif [[ "$OSTYPE" == "msys"* ]]; then
    tar -C .. -cvPf "${SRC_TAR}" --owner=0 --group=0 --exclude=electron \
        --exclude=node_modules --exclude=_deprecated --exclude='.*' \
        src cmd run-client.sh run-daemon.sh README.md INSTALLATION.md CHANGELOG.md INTEGRATION.md \
        >/dev/null
fi

popd >/dev/null

function copy_source {
    echo "Copying source tree to $1"
    mkdir -p "$1"
    cp "${SRC_TAR}" "$1"
    pushd "$1"
    tar xvPf "${SRC_TAR}" >/dev/null
    rm "${SRC_TAR}"
    popd >/dev/null
}

for var in "$@"; do
    copy_source "$var"
done

rm "${SRC_TAR}"
