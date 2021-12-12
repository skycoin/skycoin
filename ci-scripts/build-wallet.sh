#!/usr/bin/env bash

set -e -o pipefail

echo "start to build wallets..."
pushd "electron" >/dev/null
./build.sh
ls release/
popd >/dev/null
