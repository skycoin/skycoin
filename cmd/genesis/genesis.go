// Tools for creating a new blockchain
package main

import (
    "flag"
    "fmt"
    "github.com/skycoin/skycoin/src/visor"
    "os"
)

// Creates a new blockchain with a single genesis block.
// Returns the visor and signed genesis block
func createGenesisVisor(masterKeys, bcFile, bsFile string) (*visor.Visor, visor.SignedBlock, error) {
    we, err := visor.MustLoadWalletEntry(masterKeys)
    if err != nil {
        return nil, visor.SignedBlock{}, err
    }
    c := visor.NewVisorConfig()
    c.MasterKeys = we
    c.IsMaster = true
    c.BlockSigsFile = bsFile
    c.BlockchainFile = bcFile
    v := visor.NewMinimalVisor(c)
    return v, v.CreateGenesisBlock(), nil
}

func main() {
    masterKeys := flag.String("keys", "master.keys", "Master keys file")
    bcFile := flag.String("bc", "blockchain.bin",
        "Where to write the blockchain to")
    bsFile := flag.String("bs", "blockchain.sigs",
        "Where to write the blockchain signatures to")
    help := flag.Bool("help", true, "Display help message after creating")
    flag.Parse()

    v, _, err := createGenesisVisor(*masterKeys, *bcFile, *bsFile)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    }

    if err := v.SaveBlockchain(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    } else {
        fmt.Printf("Saved blockchain to %s\n", *bcFile)
    }
    if err := v.SaveBlockSigs(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        return
    } else {
        fmt.Printf("Saved blockchain signatures to %s\n", *bsFile)
    }

    if *help {
        fmt.Println("To get the timestamp:")
        fmt.Printf("\tgo run cmd/blockchain/blockchain.go -i %s -timestamp=true\n",
            *bcFile)
        fmt.Println("To get the genesis block signature:")
        fmt.Printf("\tgo run cmd/blocksigs/blocksigs.go -i %s\n", *bsFile)
    }
}
