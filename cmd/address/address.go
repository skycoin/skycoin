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

func createWalletEntry(filename string) (*visor.ReadableWalletEntry, error) {
    pub, sec := coin.GenerateKeyPair()
    addr := coin.AddressFromPubKey(pub)

    w := visor.WalletEntry{
        Address: addr,
        Public:  pub,
        Secret:  sec,
    }

    rw := visor.NewReadableWalletEntry(&w)
    if err := rw.Save(filename); err == nil {
        fmt.Printf("Wrote wallet entry to \"%s\"\n", filename)
        return &rw, nil
    } else {
        fmt.Fprintf(os.Stderr, "Failed to write wallet entry to \"%s\"\n",
            filename)
        return nil, err
    }
}

func printWalletEntryFromFile(filename string, label, address, public,
    secret bool) error {
    // Read wallet entry from disk
    w, err := visor.LoadReadableWalletEntry(filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load wallet entry \"%s\": %v\n",
            filename, err)
        return err
    }
    printWalletEntry(&w, label, address, public, secret)
    return nil
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

func handleError(err error) {
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func main() {
    registerFlags()
    parseFlags()

    coin.SetAddressVersion(addressVersion)

    if outFile != "" {
        w, err := createWalletEntry(outFile)
        handleError(err)
        printWalletEntry(w, labelStdout, printAddress, printPublic,
            printSecret)
    }
    if inFile != "" {
        err := printWalletEntryFromFile(inFile, labelStdout, printAddress,
            printPublic, printSecret)
        handleError(err)
    }
}
