// Tools for interacting with a blockchain.bin
package main

import (
    "errors"
    "flag"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
    "os"
    //"log"
)

func getBlock(filename string, seq uint64) (coin.Block, error) {
    b := coin.Block{}
    bc, err := visor.LoadBlockchain(filename)
    if err != nil {
        return b, err
    }
    if uint64(len(bc.Blocks)) < seq {
        return b, errors.New("Unknown seq")
    }
    return bc.Blocks[seq], nil
}

func main() {
    bFile := flag.String("i", "blockchain.bin", "blockchain file to load")
    bkSeq := flag.Uint64("b", 0, "block sequence to dump")
    timestamp := flag.Bool("timestamp", false, "Dump only the timestamp")
    flag.Parse()

    b, err := getBlock(*bFile, *bkSeq)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
    } else {
        if *timestamp {
            fmt.Println(b.Header.Time)
        } else {
            fmt.Println(b.String())
        }
    }
}
