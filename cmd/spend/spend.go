// Tool to create a transaction
// TODO: This tool is not functional. Do not try to use it.
package main

import (
    "flag"
    "fmt"
    "github.com/skycoin/skycoin/src/util"
    "os"
)

var (
    vc                 = visor.NewVisorConfig()
    addressVersion     = "test"
    blockchainFile     = ""
    blocksigsFile      = ""
    blockchainFileOut  = ""
    blocksigsFileOut   = ""
    transactionFileOut = ""
    isMaster           = false
    srcKeys            = ""
    srcWallet          = ""
    srcWalletEntry     = -1
    srcSecretKey       = ""
    destAddress        = ""
    spendAll           = false
    spendCoins         = uint64(0)
    spendHours         = uint64(0)
    spendFee           = uint64(0)
    burnFactor         = vc.CoinHourBurnFactor
)

func init() {
    util.DisableLogging()
}

func registerFlags() {
    flag.StringVar(&addressVersion, "address-version", addressVersion,
        "Network address version. Options are \"test\" and \"main\"")

    // Input files
    flag.StringVar(&blockchainFile, "blockchain", blockchainFile,
        "Location of the blockchain file")

    // Output files
    flag.StringVar(&transactionFileOut, "transaction-out", transactionFileOut,
        "Where to write the created transaction, if not committing to the "+
            "blockchain")

    // Ownership parameters. Some of these are mutually exclusive
    flag.StringVar(&srcWallet, "src-wallet", srcWallet,
        "Use this wallet as the spender")
    flag.IntVar(&srcWalletEntry, "src-wallet-entry", srcWalletEntry,
        "Only spend from this entry in the wallet. -1 means use all.")
    flag.StringVar(&srcSecretKey, "src-secret", srcSecretKey,
        "Spend from this secret key")

    // Spending destination
    flag.StringVar(&destAddress, "dest", destAddress,
        "Which address to send the coins to")

    // Spending amount options
    flag.Uint64Var(&spendCoins, "coins", spendCoins, "How many coins to spend")
    flag.Uint64Var(&spendHours, "hours", spendHours, "How many hours to spend")
    flag.Uint64Var(&spendFee, "fee", spendFee,
        "How much to pay in fee above the minimum fee")
    flag.Uint64Var(&burnFactor, "burn", burnFactor,
        "How many hours must be spent as a minimum fee, calculated as the "+
            "number of output hours divided by this number. 0 is no burn.")
}

func checkSrcMutex() error {
    // Handle mutually exclusive parameters
    haveCt := 0
    if srcKeys != "" {
        haveCt++
    }
    if srcWallet != "" {
        haveCt++
    }
    if srcSecretKey != "" {
        haveCt++
    }
    if haveCt != 1 {
        return errors.New("At least one of, and only one of: " +
            "[-src-keys, -src-wallet, -src-secret] can be present")
    }
    return nil
}

func checkSrcWallet() error {
    if srcWallet == "" {
        return nil
    }
    w, err := visor.LoadReadableWallet(srcWallet)
    if err != nil {
        return fmt.Errorf("Wallet \"%s\" does not exist.", srcWallet)
    }
    if srcWalletEntry < -1 || srcWalletEntry >= len(w.Entries) {
        return fmt.Errorf("Wallet entry %d does not exist in src-wallet",
            srcWalletEntry)
    }
    return nil
}

func checkDestAddress() error {
    if destAddress == "" {
        return errors.New("Destination address must be present")
    }
    _, err := coin.DecodeBase58Address(destAddress)
    if err != nil {
        return errors.New("Invalid destination address")
    }
    return nil
}

func checkSpendOptions() error {
    if spendAll {
        if spendCoins != 0 || spendHours != 0 {
            return errors.New("-spend-all is mutually exclusive with " +
                "-spend-coins and -spend-hours")
        }
    } else {
        if spendCoins == 0 {
            return errors.New("Can't spend 0 coins")
        }
    }
    return nil
}

func checkBlockchainFiles() error {
    if blockchainFile == "" {
        return errors.New("-blockchain is required")
    }
    if blocksigsFile == "" {
        return errors.New("-sigs is required")
    }
    if isMaster {
        if blocksigsFileOut == "" {
            return errors.New("-sigs-out is required if -master")
        }
        if blockchainFileOut == "" {
            return errors.New("-blockchain-out is required if -master")
        }
    } else {
        if transactionFileOut == "" {
            return errors.New("-transaction-out is required")
        }
    }
    return nil
}

func parseFlags() {
    flag.Parse()
    coin.SetAddressVersion()
    checkers := []func() error{
        checkSrcMutex,
        checkSrcWallet,
        checkDestAddress,
        checkBlockchainFiles,
    }
    for _, f := range checkers {
        if err := f(); err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            os.Exit(1)
        }
    }
}

func createTransaction() (coin.Transaction, error) {
    vc := visor.NewVisorConfig()
    vc.CoinHourBurnFactor = burnFactor
    vc.IsMaster = isMaster
    vc.MasterKeys = visor.MustLoadWalletEntry(masterKeysFile)
    vc.BlockchainFile = blockchainFile
    vc.BlockSigsFile = blocksigsFile
    v := visor.NewVisor(vc)
    addr := coin.MustDecodeBase58Address(destAddress)

    return v.Spend(amt, fee, addr)
}

func updateBlockchain(txn coin.Transaction, bcFile, sigsFile string) error {
    // Add txn to blockchain
    // Save blockchain and sigs to the file
    return nil
}

func saveTransaction(txn coin.Transaction, filename string) error {
    rtx = visor.NewReadableTransaction(&txn)
    return rtx.Save(filename)
}

func main() {
    registerFlags()
    parseFlags()
    txn, err := createTransaction()
    if err != nil {
        fmt.Fprintf(os.Stderr,
            "Failed to create transaction: %v", err)
        os.Exit(1)
    }
    if isMaster {
        if err := updateBlockchain(txn); err != nil {
            fmt.Fprintf(os.Stderr,
                "Failed to update blockchain: %v", err)
            os.Exit(1)
        } else {
            fmt.Printf("Updated blockchain")
        }
    } else {
        if err := saveTransaction(txn, transactionFileOut); err != nil {
            fmt.Fprintf(os.Stderr, "Failed to save transaction: %v",
                err)
            os.Exit(1)
            fmt.Printf("Saved transaction to \"%s\"", transactionFileOut)
        }
    }
}
