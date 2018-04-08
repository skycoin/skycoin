
# Skycoin C client library

[![Build Status](https://travis-ci.org/skycoin/skycoin.svg)](https://travis-ci.org/skycoin/skycoin)
[![GoDoc](https://godoc.org/github.com/skycoin/skycoin?status.svg)](https://godoc.org/github.com/skycoin/skycoin)
[![Go Report Card](https://goreportcard.com/badge/github.com/skycoin/skycoin)](https://goreportcard.com/report/github.com/skycoin/skycoin)

Skycoin C client library (a.k.a libskycoin) provides access to Skycoin Core
internal and API functions for implementing third-party applications.

## API Interface

The API interface is defined in the [libskycoin header file](/include/libskycoin.h).

## Building

```sh
$ make build-libc
```

## Testing

In order to test the C client libraries follow these steps

- Install [Criterion](https://github.com/Snaipe/Criterion)
  * locally by executing `make instal-deps-libc` command
  * or by [installing Criterion system-wide](https://github.com/Snaipe/Criterion#packages)
- Run `make test-libc` command

## Binary distribution

The following files will be generated

- `include/libskycoin.h` - Platform-specific header file for including libskycoin symbols in your app code
- `build/libskycoin.a` - Static library.
- `build/libskycoin.so` - Shared library object.

In Mac OS X the linker will need extra `-framework CoreFoundation -framework Security`
options.

In GNU/Linux distributions it will be necessary to load symbols in `pthread`
library e.g. by supplying extra `-lpthread` to the linker toolchain.

