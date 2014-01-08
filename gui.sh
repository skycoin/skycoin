#!/usr/bin/env bash

CMD="$1"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

usage () {
    echo "Usage: "
    echo "./gui.sh (build|run) [args]"
    exit 0
}

pushd "$DIR/compile" >/dev/null

if [[ "$CMD" = "build" ]];
then
    ./build-linux-x86_64.sh dev
elif [[ "$CMD" = "run" ]];
then
    ./release/skycoin_linux_x86_64/skycoin -disable-gui=false "${@:2}"
else
    usage
fi

popd >/dev/null

exit $?
