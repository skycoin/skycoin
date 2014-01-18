#!/usr/bin/env bash

# TODO -- update this to skycoin
REPO=
REPONAME=skycoin

# Install gvm
bash < <(curl -s https://raw.github.com/moovweb/gvm/master/binscripts/gvm-installer)

# Install go1.2
. $HOME/.gvm/scripts/gvm
gvm install go1.2
gvm use go1.2

# Install gox (this is unnecessary since we are not cross compiling)
#go get github.com/mitchellh/gox

# Grab the code
git clone $REPO "$REPONAME"
cd "${REPONAME}/compile"
