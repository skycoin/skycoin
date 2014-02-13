package main

import (
    //"encoding/json"
    "flag"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    //"github.com/skycoin/skycoin/src/visor"
    //"os"
)

// Note: Address_gen generates public keys and addresses
// address, pubkey, privatekey
// -n for number of addresses
// -seed to generate public keys and private keys deterministicly. Prompt will ask
// for seed to prevent seed from being stored in .bashrc
// without -seed, addresses will be random
// -t for test network addresses

// -a for addresses only
// -p for pubkey only
// -s for secret key only
// etc, -as for address and secret key only

// -json for json output

var (
    testNetwork  = false
    //outFile      = ""
    printAddress = false
    printPubkey  = true
    printSeckey  = true
    //labelStdout  = false
    //inFile       = ""
    seed = ""
    genCount = 1
)

func registerFlags() {

    flag.BoolVar(&testNetwork, "t", testNetwork,
        "generate testnet addresses")

    flag.IntVar(&genCount, "n", genCount,
        "number of addresses to generate")

    flag.BoolVar(&printAddress, "a", printAddress,
        "print address for generated")
    flag.BoolVar(&printPubkey, "p", printPubkey,
        "print public keys for generated")
    flag.BoolVar(&printSeckey, "s", printSeckey,
        "print secret keys for generated")

    flag.StringVar(&seed, "seed", seed,
        "seed for deterministic key generation")

    //flag.StringVar(&outFile, "o", outFile,
    //    "If present, will create a new wallet entry and write to disk. "+
    //        "For safety, it will not overwrite an existing keypair")
    //flag.BoolVar(&printAddress, "print-address", printAddress,
    //    "Print the wallet entry's address")
    //flag.BoolVar(&printPublic, "print-public", printPublic,
    //    "Print the wallet entry's public key")
    //flag.BoolVar(&printSecret, "print-secret", printSecret,
    //    "Print the wallet entry's secret key")
    //flag.StringVar(&inFile, "i", inFile,
    //    "Will read a wallet entry from this file for printing info")
    //flag.BoolVar(&labelStdout, "label-output", labelStdout,
    //    "Add a label to each printed field. This is useful if you are "+
    //        "printing multiple fields")
}

func parseFlags() {
    flag.Parse()
    //if inFile != "" && outFile != "" {
    //    fmt.Printf("-i and -o are mutually exclusive\n")
    //    os.Exit(0)
    //}
    //if inFile != "" && !printPublic && !printSecret {
    //    fmt.Printf("Input file present, but not requested to print anything\n")
    //    os.Exit(0)
    //}
}

func main() {
    registerFlags()
    parseFlags()

    if seed == "" {

        for i:=0; i<genCount; i++ {
            pub, sec := coin.GenerateKeyPair()
            addr := coin.AddressFromPubKey(pub)

            str1 := fmt.Sprintf("%v ", pub.Base64())
            str2 := fmt.Sprintf("%v ", sec.Base64())
            str3 := fmt.Sprintf("%v", addr.String())

            if printPubKey == false {
                str1 = ""
            }
            if printSecKey == false {
                str2 = ""
            }
            if printAddress == false {
                str3 = ""
            }
            fmt.Printf("%s%s%s\n", pub.Base64(), sec.Base64(), addr.String(),)
        }
    }

/*
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
*/

}
