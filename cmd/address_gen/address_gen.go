package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/cipher"
	//"log"
	//"github.com/skycoin/skycoin/src/visor"
	//"os"
)

// TODO: Make this print JSON! Needs labels and printy printing

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
	testNetwork = false
	//outFile      = ""
	PrintAddress = true
	PrintPubKey  = true
	PrintSeckey  = true

	BitcoinAddress = false
	//labelStdout  = false
	//inFile       = ""
	seed     = ""
	genCount = 1
)

func registerFlags() {

	flag.IntVar(&genCount, "n", genCount,
		"number of addresses to generate")

	flag.BoolVar(&PrintAddress, "a", PrintAddress,
		"print address for generated")
	flag.BoolVar(&PrintPubKey, "p", PrintPubKey,
		"print public keys for generated")
	flag.BoolVar(&PrintSeckey, "s", PrintSeckey,
		"print secret keys for generated")

	flag.BoolVar(&BitcoinAddress, "b", BitcoinAddress,
		"print seckey address as bitcoin address")

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

func getReadableWalletEntry(pub cipher.PubKey, sec cipher.SecKey) wallet.ReadableWalletEntry {

	var addr_str string // cipher.Address
	if BitcoinAddress == false {
		addr := cipher.AddressFromPubKey(pub)
		addr_str = addr.String()
	} else {
		addr_str = cipher.BitcoinAddressFromPubkey(pub)
	}

	str1 := fmt.Sprintf("%v ", pub.Hex())
	str2 := fmt.Sprintf("%v ", sec.Hex())
	str3 := fmt.Sprintf("%v", addr_str)

	if PrintPubKey == false {
		str1 = ""
	}
	if PrintSeckey == false {
		str2 = ""
	}
	if PrintAddress == false {
		str3 = ""
	}
	return wallet.ReadableWalletEntry{
		Address: str3,
                Public:  str1,
                Secret:  str2,
        }
}

func main() {
	registerFlags()
	parseFlags()
	meta := map[string]string{"coin": "sky"}
	entries := make([]wallet.ReadableWalletEntry, genCount)
	rw := wallet.ReadableWallet{
		Meta : meta,
		Entries : entries,
	}
	if seed == "" {
		meta["type"] = "simple"
		for i := 0; i < genCount; i++ {
			pub, sec := cipher.GenerateKeyPair()
                        entries[i] = getReadableWalletEntry(pub,sec)
		}
	}

	if seed != "" {
		meta["seed"] = seed
		meta["type"] = "deterministic"
		seckeys := cipher.GenerateDeterministicKeyPairs([]byte(seed), genCount)
		i := 0
		for _, sec := range seckeys {
			pub := cipher.PubKeyFromSecKey(sec)
			entries[i] = getReadableWalletEntry(pub,sec)
			i++
		}
		//pub, sec := cipher.GenerateDeterministicKeyPair([]byte(seed))
	}
	output, err := json.MarshalIndent(rw,  "", "    ")
	if err != nil {
		fmt.Printf("Error formating wallet to JSON. Error : %s\n", err.Error())
		return
	}
	fmt.Printf("%s\n",string(output))
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
