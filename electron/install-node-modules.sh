#!/usr/bin/env bash
set -xe -o pipefail

# installs the node modules for the skycoin electron app
# and for the electron build process

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

npm install
cd app/
./install-dependencies.sh

popd >/dev/null
