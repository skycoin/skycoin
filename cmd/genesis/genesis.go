// Tools for creating a new blockchain
package main

import (
    "errors"
    "flag"
    "fmt"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "github.com/skycoin/skycoin/src/visor"
    "log"
    "os"
)

var (
    addressVersion = "test"
    masterKeys     = "master.keys"
    bcFile         = "blockchain.bin"
    bsFile         = "blockchain.sigs"
    showHelp       = true
    destAddress    = ""
)

func init() {
    util.DisableLogging()
}

func registerFlags() {
    flag.StringVar(&addressVersion, "address-version", addressVersion,
        "Address verson to be loaded from -keys")
    flag.StringVar(&masterKeys, "keys", masterKeys, "Master keys file")
    flag.StringVar(&bcFile, "bc", bcFile,
        "Where to write the blockchain to")
    flag.StringVar(&bsFile, "bs", bsFile,
        "Where to write the blockchain signatures to")
    flag.BoolVar(&showHelp, "help", showHelp,
        "Display help message after creating")
    flag.StringVar(&destAddress, "dest-address", destAddress,
        "Which address to send the genesis balance to")
}

// Creates a new blockchain with a single genesis block.
// Returns the visor and signed genesis block
func createGenesisVisor(masterKeys, bcFile,
    bsFile string) (*visor.Visor, visor.SignedBlock, error) {
    we := visor.MustLoadWalletEntry(masterKeys)
    c := visor.NewVisorConfig()
    c.MasterKeys = we
    c.IsMaster = true
    c.BlockSigsFile = bsFile
    c.BlockchainFile = bcFile
    c.CoinHourBurnFactor = 0
    v := visor.NewMinimalVisor(c)
    v.Wallet = visor.CreateMasterWallet(c.MasterKeys)
    return v, v.CreateGenesisBlock(), nil
}

// Transfers all the coins and hours in genesis block to an address
func transferAllToAddress(v *visor.Visor, gb visor.SignedBlock,
    dest coin.Address) (visor.SignedBlock, error) {
    sb := visor.SignedBlock{}
    if gb.Block.Head.BkSeq != uint64(0) {
        return sb, errors.New("Must be genesis block")
    }
    // Send the entire genesis block to dest
    if len(gb.Block.Body.Transactions) != 1 {
        log.Panic("Genesis block has only 1 txn")
    }
    tx := gb.Block.Body.Transactions[0]
    if len(tx.Out) != 1 {
        log.Panic("Genesis block has only 1 output")
    }
    amt := visor.NewBalance(tx.Out[0].Coins, tx.Out[0].Hours)
    tx, err := v.Spend(amt, 0, dest)
    if err != nil {
        return sb, err
    }
    // Add the tx to the unconfirmed pool so it can get picked up
    err, _ = v.RecordTxn(tx)
    if err != nil {
        return sb, err
    }
    // Put the tx in a block and commit
    sb, err = v.CreateAndExecuteBlock()
    if err != nil {
        return sb, err
    }
    return sb, nil
}

func handleError(err error) {
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func main() {
    registerFlags()
    flag.Parse()

    coin.SetAddressVersion(addressVersion)

    v, gb, err := createGenesisVisor(masterKeys, bcFile, bsFile)
    handleError(err)
    if destAddress != "" {
        addr, err := coin.DecodeBase58Address(destAddress)
        handleError(err)
        _, err = transferAllToAddress(v, gb, addr)
        handleError(err)
        fmt.Printf("Transferred genesis balance to %s\n", destAddress)
    }

    err = v.SaveBlockchain()
    handleError(err)
    fmt.Printf("Saved blockchain to %s\n", bcFile)
    err = v.SaveBlockSigs()
    handleError(err)
    fmt.Printf("Saved blockchain signatures to %s\n", bsFile)

    if showHelp {
        fmt.Println("To get the timestamp:")
        fmt.Printf("\tgo run cmd/blockchain/blockchain.go -i %s -timestamp=true\n",
            bcFile)
        fmt.Println("To get the genesis block signature:")
        fmt.Printf("\tgo run cmd/blocksigs/blocksigs.go -i %s\n", bsFile)
    }
}
