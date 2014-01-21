#!/usr/bin/env bash

CONFIG=osx-x86-config
. build-tools.sh

compile_app
create_nw_bin
create_osx_package
