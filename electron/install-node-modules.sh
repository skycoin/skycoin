#!/usr/bin/env bash
set -e -o pipefail

# installs the node modules for the skycoin electron app
# NOT for the electron build process

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

pushd "$SCRIPTDIR" >/dev/null

cd src/
npm install .

popd >/dev/null
