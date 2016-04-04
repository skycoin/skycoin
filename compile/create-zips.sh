#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd "$DIR" >/dev/null

pushd release/ >/dev/null
if [[ -d "skycoin_linux_x86_64" ]]; then
    if [[ -d "skycoin_linux_x86_64.tar.gz" ]]; then
        rm skycoin_linux_x86_64.tar.gz
    fi
    tar cvzf skycoin_linux_x86_64.tar.gz skycoin_linux_x86_64
fi
if [[ -d "skycoin_linux_x86" ]]; then
    if [[ -d "skycoin_linux_x86.tar.gz" ]]; then
        rm skycoin_linux_x86.tar.gz
    fi
    tar cvzf skycoin_windows_x86.tar.gz skycoin_windows_x86
fi
if [[ -d "skycoin_windows_x86" ]]; then
    if [[ -d "skycoin_windows_x86.zip" ]]; then
        rm skycoin_windows_x86.zip
    fi
    pushd skycoin_windows_x86 >/dev/null
    zip -r skycoin_windows_x86.zip *
    mv skycoin_windows_x86.zip ../
    popd >/dev/null
fi
if [[ -d "Skycoin.app" ]]; then
    if [[ -d "Skycoin.app.zip" ]]; then
        rm Skycoin.app.zip
    fi
    zip -r Skycoin.app.zip Skycoin.app
fi
popd >/dev/null

popd >/dev/null
