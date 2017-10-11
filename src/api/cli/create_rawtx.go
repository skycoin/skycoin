package cli

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util/droplet"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"

	gcli "github.com/urfave/cli"
)

// UnspentOut wraps visor.ReadableOutput
type UnspentOut struct {
	visor.ReadableOutput
}

// SendAmount represents an amount to send to an address
type SendAmount struct {
	Addr  string
	Coins uint64
}

type sendAmountJSON struct {
	Addr  string `json:"addr"`
	Coins string `json:"coins"`
}

func createRawTxCmd(cfg Config) gcli.Command {
	name := "createRawTransaction"
	return gcli.Command{
		Name:      name,
		Usage:     "Create a raw transaction to be broadcast to the network later",
		ArgsUsage: "[to address] [amount]",
		Description: fmt.Sprintf(`
  Note: The [amount] argument is the coins you will spend, 1 coins = 1e6 droplets.

		  The default wallet (%s) will be
		  used if no wallet and address was specified.


        If you are sending from a wallet the coins will be taken iteratively
        from all addresses within the wallet starting with the first address until
        the amount of the transaction is met.

        Use caution when using the "-p" command. If you have command history enabled
        your wallet encryption password can be recovered from the history log. If you
        do not include the "-p" option you will be prompted to enter your password
        after you enter your command.`, cfg.FullWalletPath()),
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
				example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'`,
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

			if c.Bool("json") {
				return printJson(struct {
					RawTx string `json:"rawtx"`
				}{
					RawTx: rawtx,
				})
			}

			fmt.Println(rawtx)
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}

type walletAddress struct {
	Wallet  string
	Address string
}

func fromWalletOrAddress(c *gcli.Context) (walletAddress, error) {
	cfg := ConfigFromContext(c)

	wlt, err := resolveWalletPath(cfg, c.String("f"))
	if err != nil {
		return walletAddress{}, err
	}

	wltAddr := walletAddress{
		Wallet: wlt,
	}

	wltAddr.Address = c.String("a")
	if wltAddr.Address == "" {
		return wltAddr, nil
	}

	if _, err := cipher.DecodeBase58Address(wltAddr.Address); err != nil {
		return walletAddress{}, fmt.Errorf("invalid address: %s", wltAddr.Address)
	}

	return wltAddr, nil
}

func getChangeAddress(wltAddr walletAddress, chgAddr string) (string, error) {
	if chgAddr == "" {
		switch {
		case wltAddr.Address != "":
			// use the from address as change address
			chgAddr = wltAddr.Address
		case wltAddr.Wallet != "":
			// get the default wallet's coin base address
			wlt, err := wallet.Load(wltAddr.Wallet)
			if err != nil {
				return "", WalletLoadError(err)
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
		sas := []sendAmountJSON{}
		if err := json.NewDecoder(strings.NewReader(m)).Decode(&sas); err != nil {
			return nil, fmt.Errorf("invalid -m flag string, err:%v", err)
		}
		sendAmts := make([]SendAmount, 0, len(sas))
		for _, sa := range sas {
			amt, err := droplet.FromString(sa.Coins)
			if err != nil {
				return nil, fmt.Errorf("invalid coins value in -m flag string: %v", err)
			}

			sendAmts = append(sendAmts, SendAmount{
				Addr:  sa.Addr,
				Coins: amt,
			})
		}
		return sendAmts, nil
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
	amt, err := droplet.FromString(amount)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %v", err)
	}

	return amt, nil
}

func createRawTx(c *gcli.Context) (string, error) {
	rpcClient := RpcClientFromContext(c)

	wltAddr, err := fromWalletOrAddress(c)
	if err != nil {
		return "", err
	}

	chgAddr, err := getChangeAddress(wltAddr, c.String("c"))
	if err != nil {
		return "", err
	}

	toAddrs, err := getToAddresses(c)
	if err != nil {
		return "", err
	}

	if wltAddr.Address == "" {
		return CreateRawTxFromWallet(rpcClient, wltAddr.Wallet, chgAddr, toAddrs)
	}
	return CreateRawTxFromAddress(rpcClient, wltAddr.Address, wltAddr.Wallet, chgAddr, toAddrs)
}

// PUBLIC

// CreateRawTxFromWallet creates a transaction from any address or combination of addresses in a wallet
func CreateRawTxFromWallet(c *webrpc.Client, walletFile, chgAddr string, toAddrs []SendAmount) (string, error) {
	// validate the send amount
	for _, arg := range toAddrs {
		// validate to address
		_, err := cipher.DecodeBase58Address(arg.Addr)
		if err != nil {
			return "", ErrAddress
		}
	}

	// check change address
	cAddr, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return "", ErrAddress
	}

	// check if the change address is in wallet.
	wlt, err := wallet.Load(walletFile)
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

	return CreateRawTx(c, wlt, addrStrArray, chgAddr, toAddrs)
}

// Creates a transaction from a specific address in a wallet
func CreateRawTxFromAddress(c *webrpc.Client, addr, walletFile, chgAddr string, toAddrs []SendAmount) (string, error) {
	var err error
	for _, arg := range toAddrs {
		// validate the address
		if _, err = cipher.DecodeBase58Address(arg.Addr); err != nil {
			return "", ErrAddress
		}
	}

	// check if the address is in the default wallet.
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return "", err
	}

	srcAddr, err := cipher.DecodeBase58Address(addr)
	if err != nil {
		return "", ErrAddress
	}

	_, ok := wlt.GetEntry(srcAddr)
	if !ok {
		return "", fmt.Errorf("%v address is not in wallet", addr)
	}

	// validate change address
	cAddr, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return "", ErrAddress
	}

	_, ok = wlt.GetEntry(cAddr)
	if !ok {
		return "", fmt.Errorf("change address %v is not in wallet", chgAddr)
	}

	return CreateRawTx(c, wlt, []string{addr}, chgAddr, toAddrs)
}

// CreateRawTx creates a transaction from a set of addresses contained in a loaded *wallet.Wallet
func CreateRawTx(c *webrpc.Client, wlt *wallet.Wallet, inAddrs []string, chgAddr string, toAddrs []SendAmount) (string, error) {
	// get unspent outputs of those addresses
	unspents, err := c.GetUnspentOutputs(inAddrs)
	if err != nil {
		return "", err
	}

	spdouts := unspents.Outputs.SpendableOutputs()
	spendableOuts := make([]UnspentOut, len(spdouts))
	for i := range spdouts {
		spendableOuts[i] = UnspentOut{spdouts[i]}
	}

	// caculate total required amount
	var totalCoins uint64
	for _, arg := range toAddrs {
		totalCoins += arg.Coins
	}

	outs, err := getSufficientUnspents(spendableOuts, totalCoins)
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
	var totalInCoins, totalInHours, totalOutCoins uint64

	for _, o := range outs {
		c, err := droplet.FromString(o.Coins)
		if err != nil {
			return nil, errors.New("error coins string")
		}
		totalInCoins += c
		totalInHours += o.Hours
	}

	for _, to := range toAddrs {
		totalOutCoins += to.Coins
	}

	if totalInCoins < totalOutCoins {
		return nil, errors.New("amount is not sufficient")
	}

	outAddrs := []coin.TransactionOutput{}
	chgAmt := totalInCoins - totalOutCoins
	// FIXME: Why divide by 4 here?
	chgHours := totalInHours / 4
	addrHours := chgHours / uint64(len(toAddrs))
	if chgAmt > 0 {
		// generate a change address
		// FIXME: Why divide chgHours by 2 again, already divided by 4?
		outAddrs = append(outAddrs, mustMakeUtxoOutput(chgAddr, chgAmt, chgHours/2))
	}

	for _, to := range toAddrs {
		outAddrs = append(outAddrs, mustMakeUtxoOutput(to.Addr, to.Coins, addrHours))
	}

	return outAddrs, nil
}

func mustMakeUtxoOutput(addr string, coins, hours uint64) coin.TransactionOutput {
	uo := coin.TransactionOutput{}
	uo.Address = cipher.MustDecodeBase58Address(addr)
	uo.Coins = coins
	uo.Hours = hours
	return uo
}

func getKeys(wlt *wallet.Wallet, outs []UnspentOut) ([]cipher.SecKey, error) {
	keys := make([]cipher.SecKey, len(outs))
	for i, o := range outs {
		addr, err := cipher.DecodeBase58Address(o.Address)
		if err != nil {
			return nil, ErrAddress
		}
		entry, ok := wlt.GetEntry(addr)
		if !ok {
			return nil, fmt.Errorf("%v is not in wallet", o.Address)
		}

		keys[i] = entry.Secret
	}
	return keys, nil
}

func getSufficientUnspents(unspents []UnspentOut, coins uint64) ([]UnspentOut, error) {
	var totalCoins uint64
	var outs []UnspentOut

	addrOuts := make(map[string][]UnspentOut)
	for _, u := range unspents {
		addrOuts[u.Address] = append(addrOuts[u.Address], u)
	}

	for _, us := range addrOuts {
		for i, u := range us {
			coins, err := droplet.FromString(u.Coins)
			if err != nil {
				return nil, err
			}

			if coins == 0 {
				continue
			}

			totalCoins += coins
			outs = append(outs, us[i])

			if totalCoins >= coins {
				return outs, nil
			}
		}
	}

	return nil, errors.New("balance in wallet is not sufficient")
}

// NewTransaction create skycoin transaction.
func NewTransaction(utxos []UnspentOut, keys []cipher.SecKey, outs []coin.TransactionOutput) (*coin.Transaction, error) {
	tx := coin.Transaction{}
	for _, u := range utxos {
		tx.PushInput(cipher.MustSHA256FromHex(u.Hash))
	}

	for _, o := range outs {
		if err := daemon.DropletPrecisionCheck(o.Coins); err != nil {
			return nil, err
		}
		tx.PushOutput(o.Address, o.Coins, o.Hours)
	}

	tx.SignInputs(keys)
	tx.UpdateHeader()
	return &tx, nil
}
