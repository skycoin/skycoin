package main

import (
    "encoding/json"
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
)

func init() {
    util.DisableLogging()
}

func registerFlags() {
    flag.StringVar(&addressVersion, "address-version", addressVersion,
        "Which address verson to use, if creating with -o. "+
            "Options are \"test\" and \"main\"")
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

func createWalletEntry(filename string) *visor.ReadableWalletEntry {
    pub, sec := coin.GenerateKeyPair()
    addr := coin.AddressFromPubKey(pub)

    w := visor.WalletEntry{
        Address: addr,
        Public:  pub,
        Secret:  sec,
    }

    rw := visor.NewReadableWalletEntry(&w)

    b, err := json.Marshal(rw)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to encode wallet entry\n")
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return nil
    }

    flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
    f, err := os.OpenFile(filename, flags, 0600)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to open \"%s\" for writing\n",
            filename)
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return nil
    }
    defer f.Close()
    _, err = f.Write(b)
    if err == nil {
        fmt.Printf("Wrote wallet entry to \"%s\"\n", filename)
    } else {
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

    coin.SetAddressVersion(addressVersion)

    if outFile != "" {
        w := createWalletEntry(outFile)
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
