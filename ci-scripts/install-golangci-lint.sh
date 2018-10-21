#!/usr/bin/env bash

if [ -z "$VERSION" ]; then
	echo "VERSION must be set"
	exit 1
fi

# binary will be $GOPATH/bin/golangci-lint
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $GOPATH/bin v$VERSION

# or install it into ./bin/
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s vX.Y.Z

# In alpine linux (as it does not come with curl by default)
wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s vX.Y.Z

golangci-lint --version
