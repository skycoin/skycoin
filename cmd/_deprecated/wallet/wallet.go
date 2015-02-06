package main

import (
    "flag"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "github.com/skycoin/skycoin/src/visor"
    "os"
)

var (
    addressVersion = "test"
    outFile        = ""
    printPublic    = false
    printSecret    = false
    printAddress   = false
    labelStdout    = false
    inFile         = ""
    entryNumber    = -1
    populateTo     = 1
)

func init() {
    util.DisableLogging()
}

func registerFlags() {
    flag.StringVar(&addressVersion, "address-version", addressVersion,
        "Network address version. Options are \"test\" and \"main\"")

    // Writing
    flag.StringVar(&outFile, "o", outFile,
        "If present, will create a new wallet and write to disk. "+
            "For safety, it will not overwrite an existing wallet")
    flag.IntVar(&populateTo, "entries", populateTo,
        "When creating a wallet, initialize with this many entries")

    // Reading
    flag.StringVar(&inFile, "i", inFile,
        "Will read a wallet from this file for printing info")
    flag.IntVar(&entryNumber, "entry", entryNumber,
        "Which entry to print values from.  -1 means all entries")

    // Stdout controls
    flag.BoolVar(&printAddress, "print-address", printAddress,
        "Print the wallet entry's address")
    flag.BoolVar(&printPublic, "print-public", printPublic,
        "Print the wallet entry's public key")
    flag.BoolVar(&printSecret, "print-secret", printSecret,
        "Print the wallet entry's secret key")
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
    if inFile != "" && !printPublic && !printSecret && !printAddress {
        fmt.Printf("Input file present, but not requested to print anything\n")
        os.Exit(0)
    }
}

func createWallet(filename string, addressVersion string,
    populateTo int) *visor.ReadableWallet {
    coin.SetAddressVersion(addressVersion)
    w := visor.NewSimpleWallet()
    w.Populate(populateTo)
    rw := visor.NewReadableWallet(w)
    err := rw.SaveSafe(filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to write wallet to \"%s\"\n",
            filename)
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return nil
    }
    return rw
}

func printWalletFromFile(filename string, label, address, public, secret bool,
    entry int) {
    // Read wallet entry from disk
    w, err := visor.LoadReadableWallet(filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load wallet \"%s\": %v\n",
            filename, err)
        return
    }
    printWallet(w, label, address, public, secret, entry)
}

func printWallet(w *visor.ReadableWallet, label, address, public, secret bool,
    entry int) {
    if entry >= len(w.Entries) || entry < -1 {
        fmt.Fprintf(os.Stderr, "Invalid entry number %d", entry)
        os.Exit(1)
    }
    if entry == -1 {
        for _, e := range w.Entries {
            printEntry(e, label, address, public, secret)
        }
    } else {
        printEntry(w.Entries[entry], label, address, public, secret)
    }
}

func printEntry(w visor.ReadableWalletEntry, label, address, public,
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
        w := createWallet(outFile, addressVersion, populateTo)
        if w != nil {
            printWallet(w, labelStdout, printAddress, printPublic, printSecret,
                entryNumber)
        }
    }
    if inFile != "" {
        printWalletFromFile(inFile, labelStdout, printAddress, printPublic,
            printSecret, entryNumber)
    }
}
