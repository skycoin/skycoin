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

	gcli "github.com/urfave/cli"
)

var (
	// ErrTemporaryInsufficientBalance is returned if a wallet does not have enough balance for a spend, but will have enough after unconfirmed transactions confirm
	ErrTemporaryInsufficientBalance = errors.New("balance is not sufficient. Balance will be sufficient after unconfirmed transactions confirm")
)

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
			tx, err := createRawTxCmdHandler(c)
			if err != nil {
				errorWithHelp(c, err)
				return nil
			}

			rawTx := hex.EncodeToString(tx.Serialize())

			if c.Bool("json") {
				return printJSON(struct {
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

func createRawTxCmdHandler(c *gcli.Context) (*coin.Transaction, error) {
	rpcClient := RPCClientFromContext(c)

	wltAddr, err := fromWalletOrAddress(c)
	if err != nil {
		return nil, err
	}

	chgAddr, err := getChangeAddress(wltAddr, c.String("c"))
	if err != nil {
		return nil, err
	}

	toAddrs, err := getToAddresses(c)
	if err != nil {
		return nil, err
	}

	if err := validateSendAmounts(toAddrs); err != nil {
		return nil, err
	}

	if wltAddr.Address == "" {
		return CreateRawTxFromWallet(rpcClient, wltAddr.Wallet, chgAddr, toAddrs)
	}

	return CreateRawTxFromAddress(rpcClient, wltAddr.Address, wltAddr.Wallet, chgAddr, toAddrs)
}

func validateSendAmounts(toAddrs []SendAmount) error {
	for _, arg := range toAddrs {
		// validate to address
		_, err := cipher.DecodeBase58Address(arg.Addr)
		if err != nil {
			return ErrAddress
		}

		if arg.Coins == 0 {
			return errors.New("Cannot send 0 coins")
		}
	}

	if len(toAddrs) == 0 {
		return errors.New("No destination addresses")
	}

	return nil
}

// PUBLIC

// CreateRawTxFromWallet creates a transaction from any address or combination of addresses in a wallet
func CreateRawTxFromWallet(c *webrpc.Client, walletFile, chgAddr string, toAddrs []SendAmount) (*coin.Transaction, error) {
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

	return CreateRawTx(c, wlt, addrStrArray, chgAddr, toAddrs)
}

// CreateRawTxFromAddress creates a transaction from a specific address in a wallet
func CreateRawTxFromAddress(c *webrpc.Client, addr, walletFile, chgAddr string, toAddrs []SendAmount) (*coin.Transaction, error) {
	// check if the address is in the default wallet.
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, err
	}

	srcAddr, err := cipher.DecodeBase58Address(addr)
	if err != nil {
		return nil, ErrAddress
	}

	_, ok := wlt.GetEntry(srcAddr)
	if !ok {
		return nil, fmt.Errorf("%v address is not in wallet", addr)
	}

	// validate change address
	cAddr, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return nil, ErrAddress
	}

	_, ok = wlt.GetEntry(cAddr)
	if !ok {
		return nil, fmt.Errorf("change address %v is not in wallet", chgAddr)
	}

	return CreateRawTx(c, wlt, []string{addr}, chgAddr, toAddrs)
}

// CreateRawTx creates a transaction from a set of addresses contained in a loaded *wallet.Wallet
func CreateRawTx(c *webrpc.Client, wlt *wallet.Wallet, inAddrs []string, chgAddr string, toAddrs []SendAmount) (*coin.Transaction, error) {
	if err := validateSendAmounts(toAddrs); err != nil {
		return nil, err
	}

	// Get unspent outputs of those addresses
	unspents, err := c.GetUnspentOutputs(inAddrs)
	if err != nil {
		return nil, err
	}

	inUxs, err := unspents.Outputs.SpendableOutputs().ToUxArray()
	if err != nil {
		return nil, err
	}

	txn, err := createRawTx(unspents.Outputs, wlt, inAddrs, chgAddr, toAddrs)
	if err != nil {
		return nil, err
	}

	// filter out unspents which are not used in transaction
	var inUxsFiltered coin.UxArray
	for _, h := range txn.In {
		for _, u := range inUxs {
			if h == u.Hash() {
				inUxsFiltered = append(inUxsFiltered, u)
			}
		}
	}

	// TODO -- remove me -- reimplementation of visor.VerifySingleTxnSoftConstraints minus
	// the parts that require block head data, which is not available from the RPC API (see below)
	if err := verifyTransactionConstraints(txn, inUxsFiltered, visor.DefaultMaxBlockSize); err != nil {
		return nil, err
	}

	// TODO -- verify against soft and hard constraints
	// Need to get the head block to do verification.
	// The head block is not exposed over the JSON RPC, which webrpc.Client uses.
	// Need to remove the JSON RPC API and have the client make requests to the HTTP API.
	// Once the HTTP API is used,
	// Need to request /blockchain/metadata to get the head block time
	// This could lead to race conditions; /blockchain/metadata should return the full head, or have an API endpoint
	// just for the head, and/or include the head block in the get_outputs response
	// The head block is used for calculating inUxs's coin hours.
	// if err := visor.VerifySingleTxnSoftConstraints(txn, inUxs, visor.DefaultMaxBlockSize); err != nil {
	//     return nil, err
	// }
	// if err := visor.VerifySingleTxnHardConstraints(txn, head, inUxs); err != nil {
	// 	return nil, err
	// }

	return txn, nil
}

// TODO -- remove me -- reimplementation of visor.VerifySingleTxnSoftConstraints and HardConstraints
// minus the parts that require block head data, which is not available from the RPC API (see below)
func verifyTransactionConstraints(txn *coin.Transaction, uxIn coin.UxArray, maxSize int) error {
	// SOFT constraints:

	if txn.Size() > maxSize {
		return errors.New("Transaction size bigger than max block size")
	}

	if visor.TransactionIsLocked(uxIn) {
		return errors.New("Transaction has locked address inputs")
	}

	// Ignore transactions that do not conform to decimal restrictions
	for _, o := range txn.Out {
		if err := visor.DropletPrecisionCheck(o.Coins); err != nil {
			return err
		}
	}

	// HARD constraints:

	if err := txn.Verify(); err != nil {
		return err
	}

	// Checks whether ux inputs exist,
	// Check that signatures are allowed to spend inputs
	if err := txn.VerifyInput(uxIn); err != nil {
		return err
	}

	// Verify CoinHours do not overflow
	if _, err := txn.OutputHours(); err != nil {
		return err
	}

	// Check that no coins are created or destroyed
	// TODO -- use the correct block head, once we have it from the API
	// For now it doesn't matter, the block head is used to calculate the uxOut hours,
	// but we're not validating the hours
	uxOut := coin.CreateUnspents(coin.BlockHeader{
		BkSeq: 1,
	}, *txn)
	return coin.VerifyTransactionCoinsSpending(uxIn, uxOut)

	// TODO -- use coin.VerifyTransactionHoursSpending, once we have the head block
	// return coin.VerifyTransactionHoursSpending(head.Time(), uxIn, uxOut)
}

func createRawTx(uxouts visor.ReadableOutputSet, wlt *wallet.Wallet, inAddrs []string, chgAddr string, toAddrs []SendAmount) (*coin.Transaction, error) {
	// Calculate total required coins
	var totalCoins uint64
	for _, arg := range toAddrs {
		totalCoins += arg.Coins
	}

	spendOutputs, err := chooseSpends(uxouts, totalCoins)
	if err != nil {
		return nil, err
	}

	keys, err := getKeys(wlt, spendOutputs)
	if err != nil {
		return nil, err
	}

	txOuts, err := makeChangeOut(spendOutputs, chgAddr, toAddrs)
	if err != nil {
		return nil, err
	}

	tx := NewTransaction(spendOutputs, keys, txOuts)

	return tx, nil
}

func chooseSpends(uxouts visor.ReadableOutputSet, coins uint64) ([]wallet.UxBalance, error) {
	// Convert spendable unspent outputs to []wallet.UxBalance
	spendableOutputs, err := visor.ReadableOutputsToUxBalances(uxouts.SpendableOutputs())
	if err != nil {
		return nil, err
	}

	// Choose which unspent outputs to spend
	// Use the MinimizeUxOuts strategy, since this is most likely used by
	// application that may need to send frequently.
	// Using fewer UxOuts will leave more available for other transactions,
	// instead of waiting for confirmation.
	outs, err := wallet.ChooseSpendsMinimizeUxOuts(spendableOutputs, coins)
	if err != nil {
		// If there is not enough balance in the spendable outputs,
		// see if there is enough balance when including incoming outputs
		if err == wallet.ErrInsufficientBalance {
			expectedOutputs, otherErr := visor.ReadableOutputsToUxBalances(uxouts.ExpectedOutputs())
			if otherErr != nil {
				return nil, otherErr
			}

			if _, otherErr := wallet.ChooseSpendsMinimizeUxOuts(expectedOutputs, coins); otherErr != nil {
				return nil, err
			}

			return nil, ErrTemporaryInsufficientBalance
		}

		return nil, err
	}

	return outs, nil
}

func makeChangeOut(outs []wallet.UxBalance, chgAddr string, toAddrs []SendAmount) ([]coin.TransactionOutput, error) {
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

	outAddrs := []coin.TransactionOutput{}
	changeAmount := totalInCoins - totalOutCoins

	haveChange := changeAmount > 0
	nAddrs := uint64(len(toAddrs))
	changeHours, addrHours, totalOutHours := wallet.DistributeSpendHours(totalInHours, nAddrs, haveChange)

	if err := fee.VerifyTransactionFeeForHours(totalOutHours, totalInHours-totalOutHours); err != nil {
		return nil, err
	}

	for i, to := range toAddrs {
		// check if changeHours > 0, we do not need to cap addrHours when changeHours is zero
		// changeHours is zero when there is no change left or all the coinhours were used in fees
		// 1) if there is no change then the remaining coinhours are evenly distributed among the destination addresses
		// 2) if all the coinhours are burned in fees then all addrHours are zero by default
		if changeHours > 0 {
			// the coinhours are capped to a maximum of incoming coins for the address
			// if incoming coins < 1 then the cap is set to 1 coinhour

			spendCoinsAmt := to.Coins / 1e6
			if spendCoinsAmt == 0 {
				spendCoinsAmt = 1
			}

			// allow addrHours to be less than the incoming coins of the address but not more
			if addrHours[i] > spendCoinsAmt {
				// cap the addrHours, move the difference to changeHours
				changeHours += addrHours[i] - spendCoinsAmt
				addrHours[i] = spendCoinsAmt
			}
		}

		outAddrs = append(outAddrs, mustMakeUtxoOutput(to.Addr, to.Coins, addrHours[i]))
	}

	if haveChange {
		outAddrs = append(outAddrs, mustMakeUtxoOutput(chgAddr, changeAmount, changeHours))
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

func getKeys(wlt *wallet.Wallet, outs []wallet.UxBalance) ([]cipher.SecKey, error) {
	keys := make([]cipher.SecKey, len(outs))
	for i, o := range outs {
		entry, ok := wlt.GetEntry(o.Address)
		if !ok {
			return nil, fmt.Errorf("%v is not in wallet", o.Address.String())
		}

		keys[i] = entry.Secret
	}
	return keys, nil
}

// NewTransaction creates a transaction. The transaction should be validated against hard and soft constraints before transmission.
func NewTransaction(utxos []wallet.UxBalance, keys []cipher.SecKey, outs []coin.TransactionOutput) *coin.Transaction {
	tx := coin.Transaction{}
	for _, u := range utxos {
		tx.PushInput(u.Hash)
	}

	for _, o := range outs {
		tx.PushOutput(o.Address, o.Coins, o.Hours)
	}

	tx.SignInputs(keys)

	tx.UpdateHeader()
	return &tx
}
