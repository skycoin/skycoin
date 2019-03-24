#!/usr/bin/env bash

set -e -o pipefail

npm install
./node_modules/.bin/electron-rebuild --only="node-hid"
cd node_modules/hardware-wallet-js
npm run make-protobuf-files
cd ../..
