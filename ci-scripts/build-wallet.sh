#!/usr/bin/env bash

set -e -o pipefail

if [[ "$OS_NAME" == "macOS" ]] && [[ ! "${SIGN_BRANCHES[@]}" =~ "${GITHUB_REF}" || "$GITHUB_EVENT_NAME" == "pull_request" ]]; then
    export CSC_IDENTITY_AUTO_DISCOVERY=false;
fi

echo "start to build wallets..."
pushd "electron" >/dev/null
if [[ "$OS_NAME" == "macOS" ]]; then ./build.sh && ls release/; fi
popd >/dev/null
