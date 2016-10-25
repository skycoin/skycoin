package cli

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:      "createRawTransaction",
		ArgsUsage: "Create a raw transaction to be broadcast to the network later or to a remote server via RPC.",
		Usage:     "[option] [from wallet or address] [to address] [amount]",
		Description: `
        If you are sending from a wallet the coins will be taken recursively 
        from all addresses within the wallet starting with the first address until 
        the amount of the transaction is met. 
        
        Use caution when using the “-p” command. If you have command history enabled 
        your wallet encryption password can be recovered from the history log. If you 
        do not include the “-p” option you will be prompted to enter your password 
        after you enter your command.`,
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "w",
				Usage: "[wallet file or path], From wallet. If no path is specified your default wallet path will be used.",
			},
			gcli.StringFlag{
				Name:  "a",
				Usage: "[address] From address.",
			},
			gcli.StringFlag{
				Name:  "c",
				Usage: "[changeAddress] Specify different change address. By default the from address or a wallets coinbase address will be used.",
			},
			gcli.StringFlag{
				Name:  "p",
				Usage: "[password] Password for address or wallet.",
			},
			gcli.BoolFlag{
				Name:  "j,json",
				Usage: "Returns the results in JSON format.",
			},
		},
		Action: func(c *gcli.Context) error {
			rawtx, err := createRawTransaction(c)
			if err != nil {
				return err
			}

			j := c.Bool("json")
			if !j {
				fmt.Println(rawtx)
			} else {
				var jsn = struct {
					RawTx string `json:"rawtx"`
				}{rawtx}
				d, err := json.MarshalIndent(jsn, "", "    ")
				if err != nil {
					return err
				}
				fmt.Println(string(d))
			}
			return nil
		},
	}
	Commands = append(Commands, cmd)
}

func createRawTransaction(c *gcli.Context) (string, error) {
	w, a, err := fromWalletOrAddress(c)
	if err != nil {
		return "", err
	}

	var chgAddr string
	if w == "" {
		chgAddr, err = getChangeAddress(filepath.Join(walletDir, defaultWalletName), c)
	} else {
		chgAddr, err = getChangeAddress(w, c)
	}
	if err != nil {
		return "", err
	}

	toAddr, err := getToAddress(c)
	if err != nil {
		return "", err
	}

	amt, err := getAmount(c)
	if err != nil {
		return "", err
	}

	if w != "" {
		return createRawTxFromWallet(w, chgAddr, toAddr, amt)
	}

	return createRawTxFromAddress(a, chgAddr, toAddr, amt)
}

func fromWalletOrAddress(c *gcli.Context) (w string, a string, err error) {
	w = c.String("w")
	a = c.String("a")

	if a != "" && w != "" {
		// 1 1
		err = errors.New("use either -w or -a")
		return
	}

	if a == "" {
		if w == "" {
			// 0 0
			w = filepath.Join(walletDir, defaultWalletName)
			return
		}

		// 0 1
		// validate wallet file name
		if !strings.HasSuffix(w, walletExt) {
			err = errWalletName
			return
		}

		if filepath.Base(w) != w {
			w, err = filepath.Abs(filepath.Dir(w))
			return
		}
		w = filepath.Join(walletDir, w)
		return
	}
	// 1 0
	return
}

func getChangeAddress(wltFile string, c *gcli.Context) (string, error) {
	chgAddr := c.String("c")
	if chgAddr == "" {
		// get the default wallet's coin base address
		wlt, err := wallet.Load(wltFile)
		if err != nil {
			return "", errLoadWallet
		}
		return wlt.Entries[0].Address.String(), nil
	}
	// validate the address
	_, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return "", errors.New("error address")
	}

	return chgAddr, nil
}

func getToAddress(c *gcli.Context) (string, error) {
	if c.NArg() < 2 {
		return "", errors.New("error argument")
	}

	toAddr := c.Args().Get(c.NArg() - 2)
	return toAddr, nil
}

func getAmount(c *gcli.Context) (uint64, error) {
	if c.NArg() < 2 {
		return 0, errors.New("error argument")
	}
	amount := c.Args().Get(c.NArg() - 1)
	amt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return 0, errors.New("error amount")
	}
	if (amt % 1e6) != 0 {
		return 0, errors.New("skycoin coins must be multiple of 1e6")
	}
	return amt, nil
}

func createRawTxFromWallet(wltPath string, chgAddr string, toAddr string, amt uint64) (string, error) {
	// validate the amt
	if (amt % 1e6) != 0 {
		return "", errors.New("skycoin coins must be multiple of 1e6")
	}

	// check if the change address is in wallet.
	wlt, err := wallet.Load(wltPath)
	if err != nil {
		return "", errLoadWallet
	}

	// check change address
	cAddr, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return "", errAddress
	}
	_, ok := wlt.GetEntry(cAddr)
	if !ok {
		return "", fmt.Errorf("change address %v is not in wallet", chgAddr)
	}

	// get all address in the wallet
	totalAddrs := wlt.GetAddresses()
	addrStrArray := make([]string, len(totalAddrs))
	for i, a := range totalAddrs {
		addrStrArray[i] = a.String()
	}

	return makeTx(addrStrArray, chgAddr, toAddr, amt, wlt)
}

func createRawTxFromAddress(addr string, chgAddr string, toAddr string, amt uint64) (string, error) {
	if (amt % 1e6) != 0 {
		return "", errors.New("skycoin coins must be multiple of 1e6")
	}

	// check if the address is in the default wallet.
	wlt, err := wallet.Load(filepath.Join(walletDir, defaultWalletName))
	if err != nil {
		return "", errLoadWallet
	}
	srcAddr, err := cipher.DecodeBase58Address(addr)
	if err != nil {
		return "", errAddress
	}

	// validate address
	_, ok := wlt.GetEntry(srcAddr)
	if !ok {
		return "", fmt.Errorf("%v address is not in wallet", addr)
	}

	// check change address
	cAddr, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return "", errAddress
	}
	_, ok = wlt.GetEntry(cAddr)
	if !ok {
		return "", fmt.Errorf("change address %v is not in wallet", chgAddr)
	}

	return makeTx([]string{addr}, chgAddr, toAddr, amt, wlt)
}

func makeTx(inAddrs []string, chgAddr string, toAddr string, amt uint64, wlt *wallet.Wallet) (string, error) {
	// get unspent outputs of those addresses
	unspents, err := getUnspent(inAddrs)
	if err != nil {
		return "", err
	}

	outs, err := getSufficientUnspents(unspents, amt)
	if err != nil {
		return "", err
	}

	keys, err := getKeys(wlt, outs)
	if err != nil {
		return "", err
	}

	txOuts, err := makeChangeOut(outs, amt, chgAddr, toAddr)
	if err != nil {
		return "", err
	}

	tx, err := newTransaction(outs, keys, txOuts)
	if err != nil {
		return "", err
	}

	d := tx.Serialize()
	return hex.EncodeToString(d), nil
}

func makeChangeOut(outs []unspentOut, amt uint64, chgAddr string, toAddr string) ([]coin.TransactionOutput, error) {
	var (
		totalAmt   uint64
		totalHours uint64
	)

	for _, o := range outs {
		c, err := strconv.ParseUint(o.Coins, 10, 64)
		if err != nil {
			return nil, errors.New("error coins string")
		}
		totalAmt += c
		totalHours += o.Hours
	}

	if totalAmt < amt {
		return nil, errors.New("amount is not sufficient")
	}

	outAddrs := []coin.TransactionOutput{}
	chgAmt := totalAmt - amt
	chgHours := totalHours / 4
	if chgAmt > 0 {
		// generate a change address
		outAddrs = append(outAddrs,
			makeUtxoOutput(toAddr, amt, chgHours/2),
			makeUtxoOutput(chgAddr, chgAmt, chgHours/2))
	} else {
		outAddrs = append(outAddrs, makeUtxoOutput(toAddr, amt, chgHours/2))
	}
	return outAddrs, nil
}

func makeUtxoOutput(addr string, amount uint64, hours uint64) coin.TransactionOutput {
	uo := coin.TransactionOutput{}
	uo.Address = cipher.MustDecodeBase58Address(addr)
	uo.Coins = amount
	uo.Hours = hours
	return uo
}

func getKeys(wlt *wallet.Wallet, outs []unspentOut) ([]cipher.SecKey, error) {
	keys := make([]cipher.SecKey, len(outs))
	for i, o := range outs {
		addr, err := cipher.DecodeBase58Address(o.Address)
		if err != nil {
			return nil, errAddress
		}
		entry, ok := wlt.GetEntry(addr)
		if !ok {
			return nil, fmt.Errorf("%v is not in wallet", o.Address)
		}

		keys[i] = entry.Secret
	}
	return keys, nil
}

func getSufficientUnspents(unspents []unspentOut, amt uint64) ([]unspentOut, error) {
	var (
		totalAmt uint64
		outs     []unspentOut
	)

	addrOuts := make(map[string][]unspentOut)
	for _, u := range unspents {
		addrOuts[u.Address] = append(addrOuts[u.Address], u)
	}

	for _, us := range addrOuts {
		var tmpAmt uint64
		for i, u := range us {
			coins, err := strconv.ParseUint(u.Coins, 10, 64)
			if err != nil {
				return nil, errors.New("error coins string")
			}
			if coins == 0 {
				continue
			}
			tmpAmt = (coins * 1e6)
			us[i].Coins = strconv.FormatUint(tmpAmt, 10)
			totalAmt += tmpAmt
			outs = append(outs, us[i])
		}

		if totalAmt >= amt {
			return outs, nil
		}
	}

	return nil, errors.New("balance in wallet is not sufficient")
}

// NewTransaction create skycoin transaction.
func newTransaction(utxos []unspentOut, keys []cipher.SecKey, outs []coin.TransactionOutput) (*coin.Transaction, error) {
	tx := coin.Transaction{}
	// keys := make([]cipher.SecKey, len(utxos))
	for _, u := range utxos {
		tx.PushInput(cipher.MustSHA256FromHex(u.Hash))
	}

	for _, o := range outs {
		if (o.Coins % 1e6) != 0 {
			return nil, errors.New("skycoin coins must be multiple of 1e6")
		}
		tx.PushOutput(o.Address, o.Coins, o.Hours)
	}
	// tx.Verify()

	tx.SignInputs(keys)
	tx.UpdateHeader()
	return &tx, nil
}
