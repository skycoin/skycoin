#!/usr/bin/env bash

# Can't cross-compile to i386, go compiler bug
./build-linux-x86_64.sh
./build-osx-x86_64.sh
./build-windows-x86_64.sh


# Bug compiling 386:

# $OS/386 fails due to a bug in the go compiler:

# ```
# ~/projects/skycoin(branch:build-cleanup*) » gox -osarch="windows/386" -output="skycoin_win_x64" ./cmd/skycoin
# Number of parallel builds: 11

# -->     windows/386: github.com/skycoin/skycoin/cmd/skycoin

# 1 errors occurred:
# --> windows/386 error: exit status 2
# Stderr: # github.com/skycoin/skycoin/src/cipher/secp256k1-go/secp256k1-go2
# run compiler with -v for register allocation sites
# ../../.gvm/pkgsets/go1.6/skycoin/src/github.com/skycoin/skycoin/src/cipher/secp256k1-go/secp256k1-go2/field.go:49: internal compiler error: out of fixed registers
# ```
