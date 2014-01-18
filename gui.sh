#!/usr/bin/env bash

CMD="$1"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ARCH=`uname -m`
OS=`uname -s`

if [ "$ARCH" != "x86_64" ];
then
    ARCH="x86"
fi

if [ "$OS" = "Darwin" ];
then
    OS="osx"
elif [ "$OS" = "Linux" ];
then
    OS="linux"
else
    echo "Unknown OS $OS"
    exit 0
fi

usage () {
    echo "Usage: "
    echo "./gui.sh (build|run) [args]"
    exit 0
}

pushd "$DIR/compile" >/dev/null

if [[ "$CMD" = "build" ]];
then
    ./build-${OS}-${ARCH}.sh dev
elif [[ "$CMD" = "run" ]];
then
    pushd "./release/skycoin_${OS}_${ARCH}/" >/dev/null
    ./skycoin -disable-gui=false "${@:2}"
    popd >/dev/null
else
    usage
fi

popd >/dev/null

exit $?
