#!/usr/bin/env bash

set -e -o pipefail

make install-deps-ui
make check-newcoin
make build-ui-travis

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
    make integration-tests-stable
    make test-ui-e2e
fi
