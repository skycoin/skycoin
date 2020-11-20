#!/usr/bin/env bash

set -e -o pipefail

if [ -z "$VERSION" ]; then
	echo "VERSION must be set"
	exit 1
fi

# In alpine linux (as it does not come with curl by default)
wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $GOPATH/bin v$VERSION

golangci-lint --version
