#!/usr/bin/env bash

CONFIG=linux-x86_64-config
. build-tools.sh

compile_app
create_nw_bin
create_linux_package
