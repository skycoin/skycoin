#!/usr/bin/env bash

set -e -o pipefail

npm install

# Code for compiling the hw wallet js library native dependencies. Here only for precaution, should be deleted soon.
# ./node_modules/.bin/electron-rebuild --only="node-hid"
# cd node_modules/hardware-wallet-js
# npm run make-protobuf-files
# cd ../..