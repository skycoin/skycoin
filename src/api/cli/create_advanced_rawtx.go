package cli

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/fee"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"

	"strconv"

	gcli "github.com/urfave/cli"
)

// AdvancedSendAmount represents the coins and hours to send to an address
type AdvancedSendAmount struct {
	Addr  string
	Coins uint64
	Hours uint64
}

type advancedSendAmountJSON struct {
	Addr  string `json:"addr"`
	Coins string `json:"coins"`
	Hours string `json:"hours"`
}

type walletAddresses struct {
	Wallet  string
	Address []string
}

func fromWalletOrAddresses(c *gcli.Context) (walletAddresses, error) {
	cfg := ConfigFromContext(c)

	wlt, err := resolveWalletPath(cfg, c.String("f"))
	if err != nil {
		return walletAddresses{}, err
	}

	wltAddr := walletAddresses{
		Wallet: wlt,
	}

	addrs := c.String("a")
	if addrs == "" {
		return wltAddr, nil
	}

	wltAddr.Address = strings.Split(addrs, ",")
	for _, addr := range wltAddr.Address {
		_, err := cipher.DecodeBase58Address(addr)
		if err != nil {
			return walletAddresses{}, fmt.Errorf("invalid address: %s", addr)
		}
	}
	return wltAddr, nil
}

func createAdvancedRawTxCmd(cfg Config) gcli.Command {
	name := "createAdvancedRawTransaction"
	return gcli.Command{
		Name:      name,
		Usage:     "Create an advanced raw transaction to be broadcast to the network later",
		ArgsUsage: "[to address] [amount]",
		Description: `
		Note: the [amount] argument is the coins you will spend, 1 coins = 1e6 droplets.`,
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "f",
				Usage: "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.",
			},
			gcli.StringFlag{
				Name: "a",
				Usage: `[address] From address(es)
				example: -a cGPHG23We5otc5ME7EqPNKnzZe7CVm3Fjx,2E88wjND2uLGM6QCXpEDEJ29RuScXZVRLJW`,
			},
			gcli.StringFlag{
				Name:  "c",
				Usage: `[changeAddress] Specify change address, by default the from address with most coins is used`,
			},
			gcli.StringFlag{
				Name: "u",
				Usage: `[unspent address hashes] hash of unspent address(es) to use in the transaction, for multiple separate using a comma
				example: -u "7a1555d60ec1d2a8861376d4b028b321dd9c7d6b438fb628d0ec2e675d91afcb, 4e905cffbbfe0f0f4b7c62bb6c1419ad74fc923eb2d60c239d1f4cf92dce3c5e"`,
			},
			gcli.StringFlag{
				Name: "m",
				Usage: `[send to many] use JSON string to set multiple recive addresses and coins,
				example: -m '[{"addr":"$addr1", "coins": "10.2", "hours": "2"}, {"addr":"$addr2", "coins": "20"}]'`,
			},
			gcli.BoolFlag{
				Name:  "json,j",
				Usage: "Returns the results in JSON format.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			tx, err := createAdvancedRawTxCmdHandler(c)
			if err != nil {
				errorWithHelp(c, err)
				return nil
			}

			rawTx := hex.EncodeToString(tx.Serialize())

			if c.Bool("json") {
				return printJson(struct {
					RawTx string `json:"rawtx"`
				}{
					RawTx: rawTx,
				})
			}

			fmt.Println(rawTx)
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}

func createAdvancedRawTxCmdHandler(c *gcli.Context) (*coin.Transaction, error) {
	rpcClient := RpcClientFromContext(c)

	wltAddr, err := fromWalletOrAddresses(c)
	if err != nil {
		return nil, err
	}

	chgAddr, err := getChangeAddress(wltAddr, c.String("c"))
	if err != nil {
		return nil, err
	}

	toAddrs, err := advancedGetToAddresses(c)
	if err != nil {
		return nil, err
	}

	if err := validateSendAmounts(toAddrs); err != nil {
		return nil, err
	}

	// hashes of unspent outputs to spend from
	var uxOutHashes = c.String("u")
	var allowedUxOuts []string

	if uxOutHashes != "" {
		allowedUxOuts = strings.Split(uxOutHashes, ",")
	}

	if wltAddr.Address == nil {
		return AdvancedCreateRawTxFromWallet(rpcClient, wltAddr.Wallet, chgAddr, toAddrs, allowedUxOuts)
	}

	return AdvancedCreateRawTxFromAddress(rpcClient, wltAddr.Wallet, chgAddr, wltAddr.Address, toAddrs, allowedUxOuts)
}

func advancedGetToAddresses(c *gcli.Context) ([]AdvancedSendAmount, error) {
	m := c.String("m")
	if m != "" {
		var sas []advancedSendAmountJSON
		if err := json.NewDecoder(strings.NewReader(m)).Decode(&sas); err != nil {
			return nil, fmt.Errorf("invalid -m flag string, err:%v", err)
		}
		sendAmts := make([]AdvancedSendAmount, 0, len(sas))
		for _, sa := range sas {
			amt, err := droplet.FromString(sa.Coins)
			if err != nil {
				return nil, fmt.Errorf("invalid coins value in -m flag string: %v", err)
			}

			var hours uint64
			if sa.Hours != "" {
				hours, err = strconv.ParseUint(sa.Hours, 10, 0)
				if err != nil {
					return nil, fmt.Errorf("invalid coinhours in -m flag string: %v", err)
				}
			}
			sendAmts = append(sendAmts, AdvancedSendAmount{
				Addr:  sa.Addr,
				Coins: amt,
				Hours: hours,
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

	hours, err := getHours(c)
	if err != nil {
		return nil, err
	}

	// no hours given
	return []AdvancedSendAmount{{toAddr, amt, hours}}, nil
}

func AdvancedCreateRawTxFromWallet(c *webrpc.Client, walletFile, chgAddr string, toAddrs []AdvancedSendAmount, allowedUxOuts []string) (*coin.Transaction, error) {
	// check change address
	cAddr, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return nil, ErrAddress
	}

	// check if the change address is in wallet.
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, err
	}

	_, ok := wlt.GetEntry(cAddr)
	if !ok {
		return nil, fmt.Errorf("change address %v is not in wallet", chgAddr)
	}

	// get all address in the wallet
	totalAddrs := wlt.GetAddresses()
	addrStrArray := make([]string, len(totalAddrs))
	for i, a := range totalAddrs {
		addrStrArray[i] = a.String()
	}

	return AdvancedCreateRawTx(c, wlt, addrStrArray, chgAddr, toAddrs, allowedUxOuts)
}

func AdvancedCreateRawTxFromAddress(c *webrpc.Client, walletFile, chgAddr string, fromAddrs []string, toAddrs []AdvancedSendAmount, allowedUxOuts []string) (*coin.Transaction, error) {
	// check if the address is in the default wallet.
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, err
	}

	// Check that from addresses are in the wallet
	for _, addr := range fromAddrs {
		srcAddr, err := cipher.DecodeBase58Address(addr)
		if err != nil {
			return nil, ErrAddress
		}

		_, ok := wlt.GetEntry(srcAddr)
		if !ok {
			return nil, fmt.Errorf("%v address is not in wallet", addr)
		}

	}

	// validate change address
	cAddr, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return nil, ErrAddress
	}

	_, ok := wlt.GetEntry(cAddr)
	if !ok {
		return nil, fmt.Errorf("change address %v is not in wallet", chgAddr)
	}

	return AdvancedCreateRawTx(c, wlt, fromAddrs, chgAddr, toAddrs, allowedUxOuts)
}

func AdvancedCreateRawTx(c *webrpc.Client, wlt *wallet.Wallet, inAddrs []string, chgAddr string, toAddrs []AdvancedSendAmount, allowedUxOuts []string) (*coin.Transaction, error) {
	if err := validateSendAmounts(toAddrs); err != nil {
		return nil, err
	}

	// Get unspent outputs of those addresses
	// Filter using address and hashes
	var filters = map[string][]string{}
	if len(inAddrs) > 0 {
		filters["addrs"] = inAddrs
	}
	if len(allowedUxOuts) > 0 {
		filters["hashes"] = allowedUxOuts
	}

	unspents, err := c.GetUnspentOutputsWithFilters(filters)
	if err != nil {
		return nil, err
	}

	return advancedCreateRawTx(unspents.Outputs, wlt, chgAddr, toAddrs)
}

func advancedCreateRawTx(uxouts visor.ReadableOutputSet, wlt *wallet.Wallet, chgAddr string, toAddrs []AdvancedSendAmount) (*coin.Transaction, error) {
	// Calculate total required coins
	var totalCoins uint64
	for _, arg := range toAddrs {
		totalCoins += arg.Coins
	}

	outs, err := chooseSpends(uxouts, totalCoins)
	if err != nil {
		return nil, err
	}

	keys, err := getKeys(wlt, outs)
	if err != nil {
		return nil, err
	}

	txOuts, err := advancedMakeChangeOut(outs, chgAddr, toAddrs)
	if err != nil {
		return nil, err
	}

	tx := NewTransaction(outs, keys, txOuts)

	return tx, nil
}

// Total coins are calculated using unspents only
// Hence haveChange depends on unspent balances and not the total balance of the address
func advancedMakeChangeOut(outs []wallet.UxBalance, chgAddr string, toAddrs []AdvancedSendAmount) ([]coin.TransactionOutput, error) {
	var totalInCoins, totalInHours, totalOutCoins uint64

	for _, o := range outs {
		totalInCoins += o.Coins
		totalInHours += o.Hours
	}

	if totalInHours == 0 {
		return nil, fee.ErrTxnNoFee
	}

	for _, to := range toAddrs {
		totalOutCoins += to.Coins
	}

	if totalInCoins < totalOutCoins {
		return nil, wallet.ErrInsufficientBalance
	}

	var outAddrs []coin.TransactionOutput
	changeAmount := totalInCoins - totalOutCoins

	haveChange := changeAmount > 0
	nAddrs := uint64(len(toAddrs))

	// give addresses the defined hours
	addrHours := make([]uint64, nAddrs)
	for i, toAddr := range toAddrs {
		addrHours[i] = toAddr.Hours
	}

	changeHours, addrHours, totalOutHours := wallet.AdvancedDistributeSpendHours(totalInHours, nAddrs, addrHours, haveChange)
	if err := fee.VerifyTransactionFeeForHours(totalOutHours, totalInHours-totalOutHours); err != nil {
		return nil, err
	}

	if haveChange {
		outAddrs = append(outAddrs, mustMakeUtxoOutput(chgAddr, changeAmount, changeHours))
	}

	for i, to := range toAddrs {
		outAddrs = append(outAddrs, mustMakeUtxoOutput(to.Addr, to.Coins, addrHours[i]))
	}

	return outAddrs, nil
}
