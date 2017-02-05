#!/usr/bin/env bash
set -e -o pipefail

. build-conf.sh

echo "$VERSION_FILE ---- version file"

APP_VERSION=`grep version package.json | sed  's/[,]//g'`
# versionData='{ "version":"0.12.1" }';

echo "versionData='{ $APP_VERSION }';" > skycoin/current-skycoin.json
