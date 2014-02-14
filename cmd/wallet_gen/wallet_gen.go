package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
    "os"
)

var (
    testNetwork  = false
    outFile      = ""
    printPublic  = false
    printSecret  = false
    printAddress = false
    labelStdout  = false
    inFile       = ""
)

func registerFlags() {
    flag.BoolVar(&testNetwork, "test-network", testNetwork,
        "Use test network for address verson, if creating with -o")
    flag.StringVar(&outFile, "o", outFile,
        "If present, will create a new wallet entry and write to disk. "+
            "For safety, it will not overwrite an existing keypair")
    flag.BoolVar(&printAddress, "print-address", printAddress,
        "Print the wallet entry's address")
    flag.BoolVar(&printPublic, "print-public", printPublic,
        "Print the wallet entry's public key")
    flag.BoolVar(&printSecret, "print-secret", printSecret,
        "Print the wallet entry's secret key")
    flag.StringVar(&inFile, "i", inFile,
        "Will read a wallet entry from this file for printing info")
    flag.BoolVar(&labelStdout, "label-output", labelStdout,
        "Add a label to each printed field. This is useful if you are "+
            "printing multiple fields")
}

func parseFlags() {
    flag.Parse()
    if inFile != "" && outFile != "" {
        fmt.Printf("-i and -o are mutually exclusive\n")
        os.Exit(0)
    }
    if inFile != "" && !printPublic && !printSecret {
        fmt.Printf("Input file present, but not requested to print anything\n")
        os.Exit(0)
    }
}

func createWalletEntry(filename string, testNetwork bool) *visor.ReadableWalletEntry {
    pub, sec := coin.GenerateKeyPair()
    var addr coin.Address
    if testNetwork {
        addr = coin.AddressFromPubKey(pub)
    } else {
        addr = coin.AddressFromPubKey(pub)
    }

    w := visor.WalletEntry{
        Address: addr,
        Public:  pub,
        Secret:  sec,
    }

    rw := visor.NewReadableWalletEntry(&w)

    err := rw.Save(filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to write wallet entry to \"%s\"\n",
            filename)
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return nil
    }

    return &rw
}

func printWalletEntryFromFile(filename string, label, address, public,
    secret bool) {
    // Read wallet entry from disk
    w, err := visor.LoadReadableWalletEntry(filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load wallet entry \"%s\": %v\n",
            filename, err)
        return
    }
    printWalletEntry(&w, label, address, public, secret)
}

func printWalletEntry(w *visor.ReadableWalletEntry, label, address, public,
    secret bool) {
    if public {
        if label {
            fmt.Printf("Public: ")
        }
        fmt.Printf("%s\n", w.Public)
    }
    if address {
        if label {
            fmt.Printf("Address: ")
        }
        fmt.Printf("%s\n", w.Address)
    }
    if secret {
        if label {
            fmt.Printf("Secret: ")
        }
        fmt.Printf("%s\n", w.Secret)
    }
}

func main() {
    registerFlags()
    parseFlags()

    if outFile != "" {
        w := createWalletEntry(outFile, testNetwork)
        if w != nil {
            printWalletEntry(w, labelStdout, printAddress, printPublic,
                printSecret)
        }
    }
    if inFile != "" {
        printWalletEntryFromFile(inFile, labelStdout, printAddress,
            printPublic, printSecret)
    }
}
