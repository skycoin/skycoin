package main

import (
    "github.com/skycoin/skycoin/src/cli"
    "github.com/skycoin/skycoin/src/skycoin"
)

func main() {
    skycoin.Run(&cli.DaemonArgs)
}
