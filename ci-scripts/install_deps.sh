#!/usr/bin/env bash

go get github.com/mitchellh/gox
mkdir -p $GOPATH/src/github.com/skycoin
cp -rf $GOPATH/src/github.com/iketheadore/skycoin $GOPATH/src/github.com/skycoin/

# if [[ $TRAVIS_OS_NAME == 'linux' ]]; then
    # sudo apt-get install python-software-properties
    # curl -sL https://deb.nodesource.com/setup_7.x | sudo -E bash -
    # sudo apt-get install nodejs
    # sudo npm install --global electron-builder
    # sudo apt-get install --no-install-recommends -y icnsutils graphicsmagick xz-utils

    # sudo apt-get install software-properties-common
    # sudo add-apt-repository ppa:ubuntu-wine/ppa -y
    # sudo apt-get update
    # sudo apt-get install --no-install-recommends -y wine1.8
# fi