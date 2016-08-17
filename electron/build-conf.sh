#!/usr/bin/env bash

# These values are also in gulpfile.js and package.json and must be equal
ELN_VERSION="v1.2.2"
ELN_OUTPUT=".electron_output/${ELN_VERSION}"
SKY_VERSION="0.1.0"

GOX_OSARCH="linux/amd64 windows/amd64 darwin/amd64"
GOX_OUTPUT=".gox_output"

GUI_DIST_DIR="../src/gui/static/dist"  # Do not append / to this path
