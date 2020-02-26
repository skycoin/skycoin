#!/usr/bin/env bash

set -e -o pipefail

if [[ "$TRAVIS_OS_NAME" == "osx" ]] && [[ ! "$TRAVIS_BRANCH" =~ $BUILD_BRANCH || "$TRAVIS_PULL_REQUEST" == true ]]; then
    export CSC_IDENTITY_AUTO_DISCOVERY=false;
fi

echo "start to build wallets..."
pushd "electron" >/dev/null
if [[ "$TRAVIS_OS_NAME" == "linux" ]]; then ./build.sh 'linux/amd64 linux/arm' ;fi
if [[ "$TRAVIS_OS_NAME" == "osx" ]]; then ./build.sh 'darwin/amd64' ;fi
if [[ "$TRAVIS_OS_NAME" == "windows" ]];  then ./build.sh 'windows/amd64 windows/386'; fi
ls release/
popd >/dev/null
