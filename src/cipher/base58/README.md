# Fast Implementation of Base58 encoding

[![GoDoc](https://godoc.org/github.com/skycoin/skycoin/src/cipher/base58?status.svg)](https://godoc.org/github.com/skycoin/skycoin/src/cipher/base58)

Fast implementation of base58 encoding in Go.

This code is a fork of https://github.com/mr-tron/base58

Base algorithm is copied from https://github.com/trezor/trezor-crypto/blob/master/base58.c
Which was copied from an older version of libbase58 https://github.com/bitcoin/libbase58

## Performance

Other base58 golang libraries use `big.Int` which has a lot of malloc overhead and shows up as a common bottleneck when profiling.

This version removes the use of `big.Int`.

## Usage example

```go
package main

import (
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/cipher/base58"
)

func main() {
	encoded := "1QCaxc8hutpdZ62iKZsn1TCG3nh7uPZojq"
	bin, err := base58.Decode(encoded)
	if err != nil {
		fmt.Println("Decode error:", err)
		os.Exit(1)
	}

	chk := base58.Encode(bin)
	if encoded == string(chk) {
		fmt.Println("Successfully decoded then re-encoded")
	}
}
```

## base58-old

The old base58 code is retained here as a reference and used in tests to compare the output is equivalent.
