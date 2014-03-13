#!/usr/bin/env bash

CMD="$1"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ARCH=`uname -m`
OS=`uname -s`


pushd "$DIR" >/dev/null
./compile/clean-static-libs.sh >/dev/null 2>&1
popd >/dev/null

if [[ "$ARCH" != "x86_64" ]]; then
    ARCH="x86"
fi

if [[ "$OS" = "Darwin" ]]; then
    OS="osx"
    BINDIR="./release/Skycoin.app"
elif [[ "$OS" = "Linux" ]]; then
    OS="linux"
    BINDIR="./release/skycoin_${OS}_${ARCH}/"
else
    echo "Unknown OS $OS"
    exit 0
fi

usage () {
    echo "Usage: "
    echo "./gui.sh (build|run) [args]"
    exit 0
}

pushd "$DIR/compile/gui" >/dev/null

if [[ "$CMD" = "build" ]]; then
    ./build-${OS}-${ARCH}.sh skycoindev
elif [[ "$CMD" = "clean" ]]; then
    ./clean.sh
elif [[ "$CMD" = "run" || "$CMD" = "" ]]; then
    if [[ -d "$BINDIR" ]]; then
        if [[ "$OS" = "osx" ]]; then
            pushd "./release" >/dev/null
            open Skycoin.app
            popd >/dev/null
        elif [[ "$OS" = "linux" ]]; then
          pushd "$BINDIR" >/dev/null
          ./skycoin -disable-gui=false -color-log=false "${@:2}"
          popd >/dev/null
        else
            echo "Unhandled OS $OS"
            exit 0
        fi
    else
        echo "Do \"./gui.sh build\" first"
    fi
else
    usage
fi

popd >/dev/null

exit $?
