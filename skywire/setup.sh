#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd "$DIR/compile" >/dev/null
# Install gvm
echo "Installing gvm and go1.2"
./install-gvm.sh
# Add to $GOPATH
echo "Installing skycoin to \$GOPATH"
./install-to-gopath.sh 
# Install dependencies
echo "Installing or updating skycoin dependencies"
./get-dependencies.sh
echo "Done"
popd >/dev/null
