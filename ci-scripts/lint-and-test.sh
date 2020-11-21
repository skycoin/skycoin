#!/usr/bin/env bash

set -e -o pipefail

make install-deps-ui
make check-newcoin

if [[ ${TEST_SUIT} == "units" ]]; then
    echo "Do unit tests"
    make install-linters
    make lint
    make lint-ui
    make test-386
    make test-amd64
    make test-ui
elif [[ ${TEST_SUIT} == "integrations" ]]; then
    echo "Do integration tests"
    make build-ui-travis
    make test-ui-e2e
    make integration-test-stable
    make integration-test-stable-disable-wallet-api
    make integration-test-stable-enable-seed-api
    make integration-test-stable-disable-gui
elif [[ ${TEST_SUIT} == "integrations/disable-csrf" ]]; then
    echo "Do integration/disable-csrf tests"
    make integration-test-stable-disable-csrf
elif [[ ${TEST_SUIT} == "integrations/auth" ]]; then
    echo "Do integration/auth tests"
    make integration-test-stable-auth
fi
