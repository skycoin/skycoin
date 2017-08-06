package cli

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/urfave/cli"
)

type UnspentOut struct {
	visor.ReadableOutput
}

type UnspentOutSet struct {
	visor.ReadableOutputSet
}

type SendAmount struct {
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
				Usage: `[send to many] use JSON string to set multiple receive addresses and coins,
				example: -m '[{"addr":"$addr1", "coins": 10}, {"addr":"$addr2", "coins": 20}]'`,
			},
			gcli.BoolFlag{
				Name:  "json,j",
				Usage: "Returns the results in JSON format.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			rawtx, err := createRawTx(c)
			if err != nil {
				errorWithHelp(c, err)
				return nil
			}

			useJson := c.Bool("json")
			if useJson {
				return printJson(struct {
					RawTx string `json:"rawtx"`
				}{
					RawTx: rawtx,
				})
			} else {
				fmt.Println(rawtx)
			}
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}

type walletOrAddress struct {
	Wallet  string
	Address string
}

func fromWalletOrAddress(c *gcli.Context) (walletOrAddress, error) {
	walletFile := c.String("f")
	addr := c.String("a")

	if addr != "" && walletFile != "" {
		// 1 1
		return walletOrAddress{}, errors.New("use either -f or -a flag")
	}

	if addr == "" {
		if walletFile == "" {
			// 0 0
			walletFile = filepath.Join(cfg.WalletDir, cfg.DefaultWalletName)
			return walletOrAddress{Wallet: walletFile}, nil
		}

		// 0 1
		// validate wallet file name
		if !strings.HasSuffix(walletFile, walletExt) {
			return walletOrAddress{}, errWalletName
		}

		if filepath.Base(walletFile) != walletFile {
			walletFile, err := filepath.Abs(walletFile)
			return walletOrAddress{Wallet: walletFile}, err
		}
		walletFile = filepath.Join(cfg.WalletDir, walletFile)
		return walletOrAddress{Wallet: walletFile}, nil
	}

	// 1 0
	if _, err := cipher.DecodeBase58Address(addr); err != nil {
		return walletOrAddress{}, fmt.Errorf("invalid from address: %s", addr)
	}
	return walletOrAddress{Address: addr}, nil
}

func getChangeAddress(wltOrAddr walletOrAddress, chgAddr string) (string, error) {
	if chgAddr == "" {
		switch {
		case wltOrAddr.Address != "":
			// use the from address as change address
			chgAddr = wltOrAddr.Address
		case wltOrAddr.Wallet != "":
			// get the default wallet's coin base address
			wlt, err := wallet.Load(wltOrAddr.Wallet)
			if err != nil {
				return "", err
			}

			if len(wlt.Entries) > 0 {
				chgAddr = wlt.Entries[0].Address.String()
			} else {
				return "", errors.New("no change address was found")
			}
		default:
			return "", errors.New("both wallet file, from address and change address are empty")
		}
	}

	// validate the address
	_, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return "", fmt.Errorf("invalid change address: %s", chgAddr)
	}

	return chgAddr, nil
}

func getToAddresses(c *gcli.Context) ([]SendAmount, error) {
	m := c.String("m")
	if m != "" {
		toAddrs := []SendAmount{}
		if err := json.NewDecoder(strings.NewReader(m)).Decode(&toAddrs); err != nil {
			return nil, fmt.Errorf("invalid -m flag string, err:%v", err)
		}
		return toAddrs, nil
	}

	if c.NArg() < 2 {
		return nil, errors.New("invalid argument")
	}

	toAddr := c.Args().First()
	// validate address
	if _, err := cipher.DecodeBase58Address(toAddr); err != nil {
		return nil, err
	}

	amt, err := getAmount(c)
	if err != nil {
		return nil, err
	}
	return []SendAmount{{toAddr, amt}}, nil
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

func createRawTx(c *gcli.Context) (string, error) {
	wltOrAddr, err := fromWalletOrAddress(c)
	if err != nil {
		return "", err
	}

	chgAddr, err := getChangeAddress(wltOrAddr, c.String("c"))
	if err != nil {
		return "", err
	}

	toAddrs, err := getToAddresses(c)
	if err != nil {
		return "", err
	}

	if wltOrAddr.Wallet != "" {
		return CreateRawTxFromWallet(wltOrAddr.Wallet, chgAddr, toAddrs)
	}

	return CreateRawTxFromAddress(wltOrAddr.Address, chgAddr, toAddrs)
}

// PUBLIC

func CreateRawTxFromWallet(wltPath, chgAddr string, toAddrs []SendAmount) (string, error) {
	// validate the send amount
	for _, arg := range toAddrs {
		// validate to address
		_, err := cipher.DecodeBase58Address(arg.Addr)
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

	return makeTx(wlt, addrStrArray, chgAddr, toAddrs)
}

func CreateRawTxFromAddress(addr, chgAddr string, toAddrs []SendAmount) (string, error) {
	var err error
	for _, arg := range toAddrs {
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

	return makeTx(wlt, []string{addr}, chgAddr, toAddrs)
}

func makeTx(wlt *wallet.Wallet, inAddrs []string, chgAddr string, toAddrs []SendAmount) (string, error) {
	// get unspent outputs of those addresses
	unspents, err := GetUnspent(inAddrs)
	if err != nil {
		return "", err
	}

	spdouts := unspents.SpendableOutputs()
	spendableOuts := make([]UnspentOut, len(spdouts))
	for i := range spdouts {
		spendableOuts[i] = UnspentOut{spdouts[i]}
	}

	// caculate total required amount
	var totalAmt uint64
	for _, arg := range toAddrs {
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

	txOuts, err := makeChangeOut(outs, chgAddr, toAddrs)
	if err != nil {
		return "", err
	}

	tx, err := NewTransaction(outs, keys, txOuts)
	if err != nil {
		return "", err
	}

	d := tx.Serialize()
	return hex.EncodeToString(d), nil
}

func makeChangeOut(outs []UnspentOut, chgAddr string, toAddrs []SendAmount) ([]coin.TransactionOutput, error) {
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

	for _, to := range toAddrs {
		totalOutAmt += to.Coins
	}

	if totalInAmt < totalOutAmt {
		return nil, errors.New("amount is not sufficient")
	}

	outAddrs := []coin.TransactionOutput{}
	chgAmt := totalInAmt - totalOutAmt*1e6
	chgHours := totalInHours / 4
	addrHours := chgHours / uint64(len(toAddrs))
	if chgAmt > 0 {
		// generate a change address
		outAddrs = append(outAddrs, mustMakeUtxoOutput(chgAddr, chgAmt, chgHours/2))
	}

	for _, arg := range toAddrs {
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

func getKeys(wlt *wallet.Wallet, outs []UnspentOut) ([]cipher.SecKey, error) {
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

func getSufficientUnspents(unspents []UnspentOut, amt uint64) ([]UnspentOut, error) {
	var (
		totalAmt uint64
		outs     []UnspentOut
	)

	addrOuts := make(map[string][]UnspentOut)
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
func NewTransaction(utxos []UnspentOut, keys []cipher.SecKey, outs []coin.TransactionOutput) (*coin.Transaction, error) {
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

func GetUnspent(addrs []string) (UnspentOutSet, error) {
	req, err := webrpc.NewRequest("get_outputs", addrs, "1")
	if err != nil {
		return UnspentOutSet{}, fmt.Errorf("create webrpc request failed:%v", err)
	}

	rsp, err := webrpc.Do(req, cfg.RPCAddress)
	if err != nil {
		return UnspentOutSet{}, fmt.Errorf("do rpc request failed:%v", err)
	}

	if rsp.Error != nil {
		return UnspentOutSet{}, fmt.Errorf("rpc request failed, %+v", *rsp.Error)
	}

	var rlt webrpc.OutputsResult
	if err := json.NewDecoder(bytes.NewBuffer(rsp.Result)).Decode(&rlt); err != nil {
		return UnspentOutSet{}, errJSONUnmarshal
	}

	return UnspentOutSet{rlt.Outputs}, nil
}
