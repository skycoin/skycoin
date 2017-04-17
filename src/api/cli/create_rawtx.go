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

type sendToArg struct {
	Addr  string `json:"addr"`  // send to address
	Coins uint64 `json:"coins"` // send amount
}

func createRawTxCMD() gcli.Command {
	name := "createRawTransaction"
	return gcli.Command{
		Name:      name,
		Usage:     "Create a raw transaction to be broadcast to the network later",
		ArgsUsage: "[to address] [amount]",
		Description: fmt.Sprintf(`
  Note: The [amount] argument is the coins you will spend, 1 coins = 1e6 drops.
  		
		  The default wallet(%s/%s) will be 
		  used if no wallet and address was specificed. 
		

        If you are sending from a wallet the coins will be taken recursively 
        from all addresses within the wallet starting with the first address until 
        the amount of the transaction is met. 
        
        Use caution when using the "-p" command. If you have command history enabled 
        your wallet encryption password can be recovered from the history log. If you 
        do not include the "-p" option you will be prompted to enter your password 
        after you enter your command.`, cfg.WalletDir, cfg.DefaultWalletName),
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "f",
				Usage: "[wallet file or path], From wallet",
			},
			gcli.StringFlag{
				Name:  "a",
				Usage: "[address] From address",
			},
			gcli.StringFlag{
				Name: "c",
				Usage: `[changeAddress] Specify different change address. 
				By default the from address or a wallets coinbase address will be used.`,
			},
			gcli.StringFlag{
				Name: "m",
				Usage: `[send to many] use JSON string to set multiple recive addresses and coins,
				example: -m '[{"addr":"$addr1", "coins": 10}, {"addr":"$addr2", "coins": 20}]'`,
			},
			gcli.BoolFlag{
				Name:  "json,j",
				Usage: "Returns the results in JSON format.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			rawtx, err := createRawTransaction(c)
			if err != nil {
				errorWithHelp(c, err)
				return nil
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
					return errJSONMarshal
				}
				fmt.Println(string(d))
			}
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}

func createRawTransaction(c *gcli.Context) (string, error) {
	w, a, err := fromWalletOrAddress(c)
	if err != nil {
		return "", err
	}

	var chgAddr string
	chgAddr, err = getChangeAddress(w, a, c)
	if err != nil {
		return "", err
	}

	toArgs := []sendToArg{}
	m := c.String("m")
	if m != "" {
		if err := json.NewDecoder(strings.NewReader(m)).Decode(&toArgs); err != nil {
			return "", fmt.Errorf("invalid -m flag string, err:%v", err)
		}
	} else {
		toAddr, err := getToAddress(c)
		if err != nil {
			return "", err
		}

		amt, err := getAmount(c)
		if err != nil {
			return "", err
		}
		toArgs = append(toArgs, sendToArg{toAddr, amt})
	}

	if w != "" {
		return createRawTxFromWallet(w, chgAddr, toArgs...)
	}

	return createRawTxFromAddress(a, chgAddr, toArgs...)
}

func fromWalletOrAddress(c *gcli.Context) (w string, a string, err error) {
	w = c.String("f")
	a = c.String("a")

	if a != "" && w != "" {
		// 1 1
		err = errors.New("use either -f or -a flag")
		return
	}

	if a == "" {
		if w == "" {
			// 0 0
			w = filepath.Join(cfg.WalletDir, cfg.DefaultWalletName)
			return
		}

		// 0 1
		// validate wallet file name
		if !strings.HasSuffix(w, walletExt) {
			err = errWalletName
			return
		}

		if filepath.Base(w) != w {
			w, err = filepath.Abs(w)
			return
		}
		w = filepath.Join(cfg.WalletDir, w)
		return
	}

	// 1 0
	if _, err = cipher.DecodeBase58Address(a); err != nil {
		err = fmt.Errorf("invalid from address: %s", a)
	}
	return
}

func getChangeAddress(wltFile string, a string, c *gcli.Context) (string, error) {
	chgAddr := c.String("c")
	for {
		if chgAddr == "" {
			// get the default wallet's coin base address
			if a != "" {
				// use the from address as change address
				chgAddr = a
				break
			}

			if wltFile != "" {
				wlt, err := wallet.Load(wltFile)
				if err != nil {
					return "", err
				}
				if len(wlt.Entries) > 0 {
					chgAddr = wlt.Entries[0].Address.String()
					break
				}
				return "", errors.New("no change address was found")
			}
			return "", errors.New("both wallet file, from address and change address are empty")
		}
		break
	}

	// validate the address
	_, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return "", fmt.Errorf("invalid change address: %s", chgAddr)
	}

	return chgAddr, nil
}

func getToAddress(c *gcli.Context) (string, error) {
	if c.NArg() < 2 {
		return "", errors.New("invalid argument")
	}

	toAddr := c.Args().First()
	// validate address
	if _, err := cipher.DecodeBase58Address(toAddr); err != nil {
		return "", err
	}

	return toAddr, nil
}

func getAmount(c *gcli.Context) (uint64, error) {
	if c.NArg() < 2 {
		return 0, errors.New("invalid argument")
	}
	amount := c.Args().Get(1)
	amt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, errors.New("error amount")
	}

	return uint64(amt), nil
}

func createRawTxFromWallet(wltPath string, chgAddr string, toArgs ...sendToArg) (string, error) {
	// validate the send amount
	var err error
	for _, arg := range toArgs {
		// validate to address
		_, err = cipher.DecodeBase58Address(arg.Addr)
		if err != nil {
			return "", errAddress
		}
	}

	// check change address
	cAddr, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return "", errAddress
	}

	// check if the change address is in wallet.
	wlt, err := wallet.Load(wltPath)
	if err != nil {
		return "", err
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

	return makeTx(wlt, addrStrArray, chgAddr, toArgs...)
}

func createRawTxFromAddress(addr string, chgAddr string, toArgs ...sendToArg) (string, error) {
	var err error
	for _, arg := range toArgs {
		// validate the address
		if _, err = cipher.DecodeBase58Address(arg.Addr); err != nil {
			return "", errAddress
		}
	}

	// check if the address is in the default wallet.
	wlt, err := wallet.Load(filepath.Join(cfg.WalletDir, cfg.DefaultWalletName))
	if err != nil {
		return "", err
	}

	srcAddr, err := cipher.DecodeBase58Address(addr)
	if err != nil {
		return "", errAddress
	}

	_, ok := wlt.GetEntry(srcAddr)
	if !ok {
		return "", fmt.Errorf("%v address is not in wallet", addr)
	}

	// validate change address
	cAddr, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return "", errAddress
	}

	_, ok = wlt.GetEntry(cAddr)
	if !ok {
		return "", fmt.Errorf("change address %v is not in wallet", chgAddr)
	}

	return makeTx(wlt, []string{addr}, chgAddr, toArgs...)
}

func makeTx(wlt *wallet.Wallet, inAddrs []string, chgAddr string, toArgs ...sendToArg) (string, error) {
	// get unspent outputs of those addresses
	unspents, err := getUnspent(inAddrs)
	if err != nil {
		return "", err
	}

	spdouts := unspents.SpendableOutputs()
	spendableOuts := make([]unspentOut, len(spdouts))
	for i := range spdouts {
		spendableOuts[i] = unspentOut{spdouts[i]}
	}

	// caculate total required amount
	var totalAmt uint64
	for _, arg := range toArgs {
		totalAmt += arg.Coins
	}

	outs, err := getSufficientUnspents(spendableOuts, totalAmt)
	if err != nil {
		return "", err
	}

	keys, err := getKeys(wlt, outs)
	if err != nil {
		return "", err
	}

	txOuts, err := makeChangeOut(outs, chgAddr, toArgs...)
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

func makeChangeOut(outs []unspentOut, chgAddr string, toArgs ...sendToArg) ([]coin.TransactionOutput, error) {
	var (
		totalInAmt   uint64
		totalInHours uint64
		totalOutAmt  uint64
	)

	for _, o := range outs {
		c, err := strconv.ParseUint(o.Coins, 10, 64)
		if err != nil {
			return nil, errors.New("error coins string")
		}
		totalInAmt += c
		totalInHours += o.Hours
	}

	for _, to := range toArgs {
		totalOutAmt += to.Coins
	}

	if totalInAmt < totalOutAmt {
		return nil, errors.New("amount is not sufficient")
	}

	outAddrs := []coin.TransactionOutput{}
	chgAmt := totalInAmt - totalOutAmt*1e6
	chgHours := totalInHours / 4
	addrHours := chgHours / uint64(len(toArgs))
	if chgAmt > 0 {
		// generate a change address
		outAddrs = append(outAddrs, mustMakeUtxoOutput(chgAddr, chgAmt, chgHours/2))
	}

	for _, arg := range toArgs {
		outAddrs = append(outAddrs, mustMakeUtxoOutput(arg.Addr, arg.Coins*1e6, addrHours))
	}

	return outAddrs, nil
}

func mustMakeUtxoOutput(addr string, amount uint64, hours uint64) coin.TransactionOutput {
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
			totalAmt += coins
			outs = append(outs, us[i])

			if totalAmt >= amt {
				return outs, nil
			}
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
