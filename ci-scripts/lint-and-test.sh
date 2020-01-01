#!/usr/bin/env bash

set -e -o pipefail

make check-newcoin
make lint
make test-386
make test-amd64
make integration-tests-stable
make lint-ui
make build-ui-travis
make test-ui
make test-ui-e2e
