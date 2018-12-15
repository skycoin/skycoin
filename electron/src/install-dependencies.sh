#!/usr/bin/env bash
npm install
./node_modules/.bin/electron-rebuild --only="node-hid"
