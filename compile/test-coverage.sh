#!/usr/bin/env bash

# Tests a single module
# Use it like this:
#   cd src/testme
#   SKYCOINPATH=/path/to/skycoinrepo ../../compile/test-coverage.sh
# Coverage will open up in html
# You don't need to do this if you are not symlinking the repo into $GOPATH

go test -cover -coverprofile=coverage.out
sed -i 's|_'${SKYCOINPATH}'|github.com/skycoin|g' coverage.out
go tool cover -html=coverage.out
rm coverage.out
