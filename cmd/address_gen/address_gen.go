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
    PrintAddress = false
    PrintPubKey  = true
    PrintSeckey  = true
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

    flag.BoolVar(&PrintAddress, "a", PrintAddress,
        "print address for generated")
    flag.BoolVar(&PrintPubKey, "p", PrintPubKey,
        "print public keys for generated")
    flag.BoolVar(&PrintSeckey, "s", PrintSeckey,
        "print secret keys for generated")

    flag.StringVar(&seed, "seed", seed,
        "seed for deterministic key generation")

    //flag.StringVar(&outFile, "o", outFile,
    //    "If present, will create a new wallet entry and write to disk. "+
    //        "For safety, it will not overwrite an existing keypair")
    //flag.BoolVar(&PrintAddress, "print-address", PrintAddress,
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


func tstring(pub coin.Pubkey, sec coin.Seckey) string {

    addr := coin.AddressFromPubKey(pub)

    str1 := fmt.Sprintf("%v ", pub.Base64())
    str2 := fmt.Sprintf("%v ", sec.Base64())
    str3 := fmt.Sprintf("%v", addr.String())

    if PrintPubKey == false {
        str1 = ""
    }
    if PrintSeckey == false {
        str2 = ""
    }
    if PrintAddress == false {
        str3 = ""
    }

    return fmt.Sprintf("%s%s%s\n", str1,str2,str3)
}

func main() {
    registerFlags()
    parseFlags()

    if seed == "" {

        for i:=0; i<genCount; i++ {
            pub, sec := coin.GenerateKeyPair()


        }
    }

    if seed != "" {
        if n != 1 {
            log.Panic("multiple deterministic addresses not implemented yet")
        }


    }

/*
    if outFile != "" {
        w := createWalletEntry(outFile, testNetwork)
        if w != nil {
            printWalletEntry(w, labelStdout, PrintAddress, printPublic,
                printSecret)
        }
    }
    if inFile != "" {
        printWalletEntryFromFile(inFile, labelStdout, PrintAddress,
            printPublic, printSecret)
    }
*/

}
