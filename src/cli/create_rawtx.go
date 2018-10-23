package cli

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/fee"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"

	"encoding/csv"
	"encoding/json"

	gcli "github.com/spf13/cobra"
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

func createRawTxCmd() *gcli.Command {
	createRawTxCmd := &gcli.Command{
		Short: "Create a raw transaction to be broadcast to the network later",
		Use:   "createRawTransaction [flags] [to address] [amount]",
		Long: fmt.Sprintf(`Note: The [amount] argument is the coins you will spend, 1 coins = 1e6 droplets.
    The default wallet (%s) will be used if no wallet and address was specified.
    
    If you are sending from a wallet the coins will be taken iteratively
    from all addresses within the wallet starting with the first address until
    the amount of the transaction is met.

    Use caution when using the "-p" command. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`, cliConfig.FullWalletPath()),
        SilenceUsage: true,
		Args: gcli.MinimumNArgs(2),
		RunE: func(c *gcli.Command, args []string) error {
			txn, err := createRawTxnCmdHandler(args)
			switch err.(type) {
			case nil:
			case WalletLoadError:
				printHelp(c)
				return err
			default:
				return err
			}

			rawTxn := hex.EncodeToString(txn.Serialize())

			if jsonOutput {
				return printJSON(struct {
					RawTx string `json:"rawtx"`
				}{
					RawTx: rawTxn,
				})
			}

			fmt.Println(rawTxn)

			return nil
		},
	}

	createRawTxCmd.Flags().StringVarP(&walletFile, "wallet-file", "f", "", "wallet file or path. If no path is specified your default wallet path will be used.")
	createRawTxCmd.Flags().StringVarP(&address, "address", "a", "", "From address")
	createRawTxCmd.Flags().StringVarP(&changeAddress, "change-address", "c", "", `Specify different change address.
By default the from address or a wallets coinbase address will be used.`)
	createRawTxCmd.Flags().StringVarP(&many, "many", "m", "", `use JSON string to set multiple receive addresses and coins,
example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'`)
	createRawTxCmd.Flags().StringVarP(&password, "password", "p", "", "Wallet password")
	createRawTxCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Returns the results in JSON format.")
	createRawTxCmd.Flags().StringVar(&csvFile, "csv-file", "", "CSV file containing addresses and amounts to send")

	return createRawTxCmd
}

type walletAddress struct {
	Wallet  string
	Address string
}

func fromWalletOrAddress() (walletAddress, error) {
	wlt, err := resolveWalletPath(cliConfig, walletFile)
	if err != nil {
		return walletAddress{}, err
	}

	wltAddr := walletAddress{
		Wallet: wlt,
	}

	wltAddr.Address = address
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
				return "", WalletLoadError{err}
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

func getToAddresses(args []string) ([]SendAmount, error) {
	if csvFile != "" && many != "" {
		return nil, errors.New("-csv and -m cannot be combined")
	}

	if many != "" {
		return parseSendAmountsFromJSON(many)
	} else if csvFile != "" {
		fields, err := openCSV(csvFile)
		if err != nil {
			return nil, err
		}
		return parseSendAmountsFromCSV(fields)
	}

	toAddr := args[0]

	if _, err := cipher.DecodeBase58Address(toAddr); err != nil {
		return nil, err
	}

	amt, err := getAmount(args)
	if err != nil {
		return nil, err
	}

	return []SendAmount{{
		Addr:  toAddr,
		Coins: amt,
	}}, nil
}

func openCSV(csvFile string) ([][]string, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	return r.ReadAll()
}

func parseSendAmountsFromCSV(fields [][]string) ([]SendAmount, error) {
	var sends []SendAmount
	var errs []error
	for i, f := range fields {
		addr := f[0]

		addr = strings.TrimSpace(addr)

		if _, err := cipher.DecodeBase58Address(addr); err != nil {
			err = fmt.Errorf("[row %d] Invalid address %s: %v", i, addr, err)
			errs = append(errs, err)
			continue
		}

		coins, err := droplet.FromString(f[1])
		if err != nil {
			err = fmt.Errorf("[row %d] Invalid amount %s: %v", i, f[1], err)
			errs = append(errs, err)
			continue
		}

		sends = append(sends, SendAmount{
			Addr:  addr,
			Coins: coins,
		})
	}

	if len(errs) > 0 {
		errMsgs := make([]string, len(errs))
		for i, err := range errs {
			errMsgs[i] = err.Error()
		}

		errMsg := strings.Join(errMsgs, "\n")

		return nil, errors.New(errMsg)
	}

	return sends, nil
}

func parseSendAmountsFromJSON(m string) ([]SendAmount, error) {
	sas := []sendAmountJSON{}

	if err := json.NewDecoder(strings.NewReader(m)).Decode(&sas); err != nil {
		return nil, fmt.Errorf("invalid -m flag string, err: %v", err)
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

func getAmount(args []string) (uint64, error) {
	amount := args[1]
	amt, err := droplet.FromString(amount)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %v", err)
	}

	return amt, nil
}

// createRawTxArgs are encapsulated arguments for creating a transaction
type createRawTxArgs struct {
	WalletID      string
	Address       string
	ChangeAddress string
	SendAmounts   []SendAmount
	Password      PasswordReader
}

func parseCreateRawTxArgs(args []string) (*createRawTxArgs, error) {
	wltAddr, err := fromWalletOrAddress()
	if err != nil {
		return nil, err
	}

	chgAddr, err := getChangeAddress(wltAddr, changeAddress)
	if err != nil {
		return nil, err
	}

	toAddrs, err := getToAddresses(args)
	if err != nil {
		return nil, err
	}

	if err := validateSendAmounts(toAddrs); err != nil {
		return nil, err
	}

	pr := NewPasswordReader([]byte(password))

	return &createRawTxArgs{
		WalletID:      wltAddr.Wallet,
		Address:       wltAddr.Address,
		ChangeAddress: chgAddr,
		SendAmounts:   toAddrs,
		Password:      pr,
	}, nil
}

func createRawTxnCmdHandler(args []string) (*coin.Transaction, error) {
	parsedArgs, err := parseCreateRawTxArgs(args)
	if err != nil {
		return nil, err
	}

	if parsedArgs.Address == "" {
		return CreateRawTxFromWallet(apiClient, parsedArgs.WalletID, parsedArgs.ChangeAddress, parsedArgs.SendAmounts, parsedArgs.Password)
	}

	return CreateRawTxFromAddress(apiClient, parsedArgs.Address, parsedArgs.WalletID, parsedArgs.ChangeAddress, parsedArgs.SendAmounts, parsedArgs.Password)
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
func CreateRawTxFromWallet(c GetOutputser, walletFile, chgAddr string, toAddrs []SendAmount, pr PasswordReader) (*coin.Transaction, error) {
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

	switch pr.(type) {
	case nil:
		if wlt.IsEncrypted() {
			return nil, wallet.ErrWalletEncrypted
		}
	case PasswordFromBytes:
		p, err := pr.Password()
		if err != nil {
			return nil, err
		}

		if !wlt.IsEncrypted() && len(p) != 0 {
			return nil, wallet.ErrWalletNotEncrypted
		}
	}

	var password []byte
	if wlt.IsEncrypted() {
		var err error
		password, err = pr.Password()
		if err != nil {
			return nil, err
		}
	}

	// get all address in the wallet
	totalAddrs := wlt.GetAddresses()
	addrStrArray := make([]string, len(totalAddrs))
	for i, a := range totalAddrs {
		addrStrArray[i] = a.String()
	}

	return CreateRawTx(c, wlt, addrStrArray, chgAddr, toAddrs, password)
}

// CreateRawTxFromAddress creates a transaction from a specific address in a wallet
func CreateRawTxFromAddress(c GetOutputser, addr, walletFile, chgAddr string, toAddrs []SendAmount, pr PasswordReader) (*coin.Transaction, error) {
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

	switch pr.(type) {
	case nil:
		if wlt.IsEncrypted() {
			return nil, wallet.ErrWalletEncrypted
		}
	case PasswordFromBytes:
		p, err := pr.Password()
		if err != nil {
			return nil, err
		}

		if !wlt.IsEncrypted() && len(p) != 0 {
			return nil, wallet.ErrWalletNotEncrypted
		}
	}

	var password []byte
	if wlt.IsEncrypted() {
		var err error
		password, err = pr.Password()
		if err != nil {
			return nil, err
		}
	}

	return CreateRawTx(c, wlt, []string{addr}, chgAddr, toAddrs, password)
}

// GetOutputser implements unspent output querying
type GetOutputser interface {
	OutputsForAddresses([]string) (*readable.UnspentOutputsSummary, error)
}

// CreateRawTx creates a transaction from a set of addresses contained in a loaded *wallet.Wallet
func CreateRawTx(c GetOutputser, wlt *wallet.Wallet, inAddrs []string, chgAddr string, toAddrs []SendAmount, password []byte) (*coin.Transaction, error) {
	if err := validateSendAmounts(toAddrs); err != nil {
		return nil, err
	}

	// Get unspent outputs of those addresses
	outputs, err := c.OutputsForAddresses(inAddrs)
	if err != nil {
		return nil, err
	}

	inUxs, err := outputs.SpendableOutputs().ToUxArray()
	if err != nil {
		return nil, err
	}

	txn, err := createRawTx(outputs, wlt, chgAddr, toAddrs, password)
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

	head, err := outputs.Head.ToCoinBlockHeader()
	if err != nil {
		return nil, err
	}

	if err := visor.VerifySingleTxnSoftConstraints(*txn, head.Time, inUxsFiltered, visor.DefaultMaxBlockSize); err != nil {
		return nil, err
	}
	if err := visor.VerifySingleTxnHardConstraints(*txn, head, inUxsFiltered); err != nil {
		return nil, err
	}
	if err := visor.VerifySingleTxnUserConstraints(*txn); err != nil {
		return nil, err
	}

	return txn, nil
}

func createRawTx(uxouts *readable.UnspentOutputsSummary, wlt *wallet.Wallet, chgAddr string, toAddrs []SendAmount, password []byte) (*coin.Transaction, error) {
	// Calculate total required coins
	var totalCoins uint64
	for _, arg := range toAddrs {
		var err error
		totalCoins, err = coin.AddUint64(totalCoins, arg.Coins)
		if err != nil {
			return nil, err
		}
	}

	spendOutputs, err := chooseSpends(uxouts, totalCoins)
	if err != nil {
		return nil, err
	}

	txOuts, err := makeChangeOut(spendOutputs, chgAddr, toAddrs)
	if err != nil {
		return nil, err
	}

	f := func(w *wallet.Wallet) (*coin.Transaction, error) {
		keys, err := getKeys(w, spendOutputs)
		if err != nil {
			return nil, err
		}

		return NewTransaction(spendOutputs, keys, txOuts), nil
	}

	makeTx := func() (*coin.Transaction, error) {
		return f(wlt)
	}

	if wlt.IsEncrypted() {
		makeTx = func() (*coin.Transaction, error) {
			var tx *coin.Transaction
			if err := wlt.GuardView(password, func(w *wallet.Wallet) error {
				var err error
				tx, err = f(w)
				return err
			}); err != nil {
				return nil, err
			}

			return tx, nil
		}
	}

	return makeTx()
}

func chooseSpends(uxouts *readable.UnspentOutputsSummary, coins uint64) ([]wallet.UxBalance, error) {
	// Convert spendable unspent outputs to []wallet.UxBalance
	spendableOutputs, err := readable.OutputsToUxBalances(uxouts.SpendableOutputs())
	if err != nil {
		return nil, err
	}

	// Choose which unspent outputs to spend
	// Use the MinimizeUxOuts strategy, since this is most likely used by
	// application that may need to send frequently.
	// Using fewer UxOuts will leave more available for other transactions,
	// instead of waiting for confirmation.
	outs, err := wallet.ChooseSpendsMinimizeUxOuts(spendableOutputs, coins, 0)
	if err != nil {
		// If there is not enough balance in the spendable outputs,
		// see if there is enough balance when including incoming outputs
		if err == wallet.ErrInsufficientBalance {
			expectedOutputs, otherErr := readable.OutputsToUxBalances(uxouts.ExpectedOutputs())
			if otherErr != nil {
				return nil, otherErr
			}

			if _, otherErr := wallet.ChooseSpendsMinimizeUxOuts(expectedOutputs, coins, 0); otherErr != nil {
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
