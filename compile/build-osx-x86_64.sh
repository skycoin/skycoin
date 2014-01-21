#!/usr/bin/env bash

CONFIG=osx-x86_64-config
. build-tools.sh

compile_app
create_nw_bin
create_osx_package
