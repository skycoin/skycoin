#!/usr/bin/env bash

echo "start to build wallets..."
pushd "electron"
yarn
if [[ "$TRAVIS_OS_NAME" == "linux" ]]; then ./build.sh 'linux/amd64 linux/arm' ;fi
if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then ./build.sh 'darwin/amd64' ;fi
ls release/
popd