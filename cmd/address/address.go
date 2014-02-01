package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "os"
)

var (
    testNetwork = false
    filename    = ""
)

func registerFlags() {
    flag.BoolVar(&testNetwork, "test-network", testNetwork,
        "Create an address for the test network")
    flag.StringVar(&filename, "o", filename,
        "File to write address and keys to. "+
            "If not provided, prints to stdout. "+
            "If the file exists, it will not be overwritten.")
}

func main() {
    registerFlags()
    flag.Parse()
    pub, sec := coin.GenerateKeyPair()
    var addr coin.Address
    if testNetwork {
        addr = coin.AddressFromPubkeyTestNet(pub)
    } else {
        addr = coin.AddressFromPubKey(pub)
    }

    w := visor.WalletEntry{
        Address:   addr,
        PublicKey: pub,
        SecretKey: sec,
    }

    rw := visor.NewReadableWalletEntry(&w)

    b, err := json.Marshal(rw)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to encode wallet entry\n")
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return
    }

    if filename == "" {
        fmt.Println(string(b))
    } else {
        flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
        f, err := os.OpenFile(filename, flags, 0600)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to open \"%s\" for writing\n",
                filename)
            fmt.Fprintf(os.Stderr, "%v\n", err)
            return
        }
        defer f.Close()
        _, err = f.Write(b)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Failed to write wallet entry to \"%s\"\n",
                filename)
            fmt.Fprintf(os.Stderr, "%v\n", err)
        }
    }
}
