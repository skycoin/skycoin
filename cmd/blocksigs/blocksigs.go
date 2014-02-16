// Extract a signature(s) from a blockchain.sigs file
package main

import (
    "errors"
    "flag"
    "fmt"
    "github.com/skycoin/skycoin/src/util"
    "github.com/skycoin/skycoin/src/visor"
    "os"
)

func init() {
    util.DisableLogging()
}

// Fetches the saved block signature at sequence bkSeq from sigsFile
func readSignature(sigsFile string, bkSeq uint64) (string, error) {
    bs, err := visor.LoadBlockSigs(sigsFile)
    if err != nil {
        return "", err
    }
    sig, ok := bs.Sigs[0]
    if !ok {
        return "", errors.New("No block found")
    }
    return sig.Hex(), nil
}

func main() {
    sigsFile := flag.String("i", "blockchain.sigs",
        "blockchain.sigs file to read genesis signature from")
    bkSeq := flag.Uint64("b", 0, "Which block to dump signature from.")
    flag.Parse()
    if *sigsFile == "" {
        fmt.Fprintf(os.Stderr, "blocksigs file required [-i]\n")
        flag.PrintDefaults()
        return
    }
    sig, err := readSignature(*sigsFile, *bkSeq)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(sig)
    }
}
