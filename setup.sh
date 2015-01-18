#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pushd "$DIR/compile" >/dev/null
# Install gvm
echo "Installing gvm and go1.4"
./install-gvm.sh
# Add to $GOPATH
echo "Installing skycoin to \$GOPATH"
./install-to-gopath.sh 
# Install dependencies
echo "Installing or updating skycoin dependencies"
./get-dependencies.sh
echo "Done"
echo "Do './run.sh -h' to confirm it is installed. This runs the daemon."
echo "README.md has further instructions for building and running the gui."
popd >/dev/null
