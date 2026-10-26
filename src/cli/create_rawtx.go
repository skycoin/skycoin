package cli

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/mathutil"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/bip44wallet"
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

func createRawTxnCmd() *cobra.Command {
	createRawTxnCmd := &cobra.Command{
		Short: "Create a raw transaction that can be broadcast to the network later",
		Use:   "createRawTransaction [wallet] [to address] [amount]",
		Long: `Create a raw transaction that can be broadcast to the network later.

    Note: The [amount] argument is the coins you will spend, with decimal formatting, e.g. 1, 1.001 or 1.000000.

    The [to address] and [amount] arguments can be replaced with the --many/-m or the --csv option.

    Use caution when using the "-p" command. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`,
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			jsonOutput, err := c.Flags().GetBool("json")
			if err != nil {
				return err
			}

			txn, err := createRawTxnCmdHandler(c, args)
			switch err.(type) {
			case nil:
			case WalletLoadError:
				printHelp(c)
				return err
			default:
				return err
			}

			rawTxn, err := txn.SerializeHexWithPQ()
			if err != nil {
				return err
			}

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

	createRawTxnCmd.Flags().StringP("from-address", "a", "", "From address in wallet")
	createRawTxnCmd.Flags().StringP("change-address", "c", "", `Specify the change address.
Defaults to a new address added to the wallet: the next unused change address for bip44 wallets, which needs the node, and an input address for collection wallets.`)
	createRawTxnCmd.Flags().StringP("many", "m", "", `use JSON string to set multiple receive addresses and coins,
example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'`)
	createRawTxnCmd.Flags().StringP("password", "p", "", "Wallet password")
	createRawTxnCmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format.")
	createRawTxnCmd.Flags().String("csv", "", "CSV file containing addresses and amounts to send")

	return createRawTxnCmd
}

func createRawTxnV2Cmd() *cobra.Command {
	createRawTxnCmd := &cobra.Command{
		Short: "Create a raw transaction that can be broadcast to the network later",
		Use:   "createRawTransactionV2 [wallet] [to address] [amount]",
		Long: `Create a raw transaction that can be broadcast to the network later.

    Note: The [amount] argument is the coins you will spend, with decimal formatting, e.g. 1, 1.001 or 1.000000.

    The [to address] and [amount] arguments can be replaced with the --csv option.,

    Use caution when using the "-p" command. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`,
		SilenceUsage: true,
		Args:         cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			jsonOutput, err := c.Flags().GetBool("json")
			if err != nil {
				return err
			}

			req, err := makeWalletCreateTxnRequest(c, args)
			if err != nil {
				return err
			}

			to := make([]string, len(req.To))
			for i, r := range req.To {
				to[i] = r.Address
			}
			warnSpentAddresses(apiClient, os.Stderr, to)

			rsp, err := apiClient.WalletCreateTransaction(*req)
			if err != nil {
				return err
			}

			if jsonOutput {
				return printJSON(rsp)
			}

			fmt.Println(rsp.EncodedTransaction)

			return nil
		},
	}

	createRawTxnCmd.Flags().StringP("from-address", "a", "", "From address in wallet")
	createRawTxnCmd.Flags().StringP("change-address", "c", "", `Specify the change address.
	Defaults to a new address added to the wallet: the next unused change address for bip44 wallets, and an input address for collection wallets.`)
	createRawTxnCmd.Flags().String("csv", "", "CSV file containing addresses and amounts to send")
	createRawTxnCmd.Flags().StringP("password", "p", "", "Wallet password")
	createRawTxnCmd.Flags().BoolP("unsign", "", false, "Do not sign the transaction")
	createRawTxnCmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format.")

	createRawTxnCmd.Flags().BoolP("ignore-unconfirmed", "", false, "Ignore unconfirmed transactions")
	createRawTxnCmd.Flags().StringP("hours-selection-type", "", transaction.HoursSelectionTypeAuto, "Hours selection type")
	createRawTxnCmd.Flags().StringP("hours-selection-mode", "", transaction.HoursSelectionModeShare, "Hours selection mode")
	createRawTxnCmd.Flags().StringP("hours-selection-share-factor", "", "0.5", "Hour selection share factor")

	return createRawTxnCmd
}

func makeWalletCreateTxnRequest(c *cobra.Command, args []string) (*api.WalletCreateTransactionRequest, error) {
	unsign, err := c.Flags().GetBool("unsign")
	if err != nil {
		return nil, err
	}

	walletFile := args[0]
	w, err := apiClient.Wallet(walletFile)
	if err != nil {
		return nil, err
	}

	wltAddr, err := fromWalletOrAddress(c, walletFile)
	if err != nil {
		return nil, err
	}

	var addrs []string
	if wltAddr.Address != "" {
		addrs = append(addrs, wltAddr.Address)
	} else {
		for _, e := range w.Entries {
			addrs = append(addrs, e.Address)
		}
	}

	ctr, err := makeCreateTransactionRequest(c, args, addrs)
	if err != nil {
		return nil, err
	}

	req := api.WalletCreateTransactionRequest{
		Unsigned:                 unsign,
		WalletID:                 w.Meta.Filename,
		CreateTransactionRequest: *ctr,
	}

	if w.Meta.Encrypted && !unsign {
		p, err := getPassword(c)
		if err != nil {
			return nil, err
		}
		defer func() {
			p = nil
		}()
		req.Password = string(p)
	}
	return &req, nil
}

func getPassword(c *cobra.Command) ([]byte, error) {
	p, err := c.Flags().GetString("password")
	if err != nil {
		return nil, err
	}
	defer func() {
		p = ""
	}()

	return NewPasswordReader([]byte(p)).Password()
}

func makeCreateTransactionRequest(c *cobra.Command, args []string, fromAddrs []string) (*api.CreateTransactionRequest, error) {
	hoursSelection, err := getHoursSelection(c)
	if err != nil {
		return nil, err
	}

	iu, err := c.Flags().GetBool("ignore-unconfirmed")
	if err != nil {
		return nil, err
	}

	var changeAddr *string
	ca, err := c.Flags().GetString("change-address")
	if err != nil {
		return nil, err
	}
	if ca != "" {
		changeAddr = &ca
	}

	to, err := getToAddressesV2(c, args[1:])
	if err != nil {
		return nil, err
	}

	return &api.CreateTransactionRequest{
		IgnoreUnconfirmed: iu,
		HoursSelection:    *hoursSelection,
		ChangeAddress:     changeAddr,
		Addresses:         fromAddrs,
		To:                to,
	}, nil
}

func getToAddressesV2(c *cobra.Command, args []string) ([]api.Receiver, error) {
	csvFile, err := c.Flags().GetString("csv")
	if err != nil {
		return nil, err
	}

	if csvFile != "" {
		fields, err := openCSV(csvFile)
		if err != nil {
			return nil, err
		}
		return parseReceiversFromCSV(fields)
	}

	if len(args) < 2 {
		return nil, fmt.Errorf("requires at least 2 arg(s), only received %d", len(args))
	}

	toAddr := args[0]
	if _, err := cipher.DecodeBase58Address(toAddr); err != nil {
		return nil, err
	}

	coins := args[1]
	if _, err := droplet.FromString(coins); err != nil {
		return nil, err
	}

	return []api.Receiver{{
		Address: toAddr,
		Coins:   coins,
	}}, nil
}
func getHoursSelection(c *cobra.Command) (*api.HoursSelection, error) {
	hst, err := c.Flags().GetString("hours-selection-type")
	if err != nil {
		return nil, err
	}

	hsm, err := c.Flags().GetString("hours-selection-mode")
	if err != nil {
		return nil, err
	}

	sf, err := c.Flags().GetString("hours-selection-share-factor")
	if err != nil {
		return nil, err
	}

	return &api.HoursSelection{
		Type:        hst,
		Mode:        hsm,
		ShareFactor: sf,
	}, nil
}

type walletAddress struct {
	Wallet  string
	Address string
}

func fromWalletOrAddress(c *cobra.Command, walletFile string) (walletAddress, error) {
	address, err := c.Flags().GetString("from-address")
	if err != nil {
		return walletAddress{}, err
	}

	wltAddr := walletAddress{
		Wallet: walletFile,
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

// getChangeAddress validates chgAddr. An empty one gets the next unused change
// address of a bip44 wallet, or stays empty, and the transaction then gets a new
// change address from the wallet (see newChangeAddress).
func getChangeAddress(wltAddr walletAddress, chgAddr string, tf wallet.TransactionsFinder) (string, error) {
	if chgAddr == "" && wltAddr.Wallet != "" {
		a, err := bip44ChangeAddress(wltAddr.Wallet, tf)
		if err != nil {
			return "", err
		}
		chgAddr = a
	}

	if chgAddr == "" {
		if wltAddr.Address == "" && wltAddr.Wallet == "" {
			return "", errors.New("wallet file, from address and change address are empty")
		}
		return "", nil
	}

	// validate the address
	_, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return "", fmt.Errorf("invalid change address: %s", chgAddr)
	}

	return chgAddr, nil
}

func getToAddresses(c *cobra.Command, args []string) ([]SendAmount, error) {
	csvFile, err := c.Flags().GetString("csv")
	if err != nil {
		return nil, err
	}
	many, err := c.Flags().GetString("many")
	if err != nil {
		return nil, err
	}

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

	if len(args) < 2 {
		return nil, fmt.Errorf("requires at least 2 arg(s), only received %d", len(args))
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
	f, err := os.Open(csvFile) //nolint:gosec
	if err != nil {
		return nil, err
	}
	defer f.Close() //nolint:errcheck

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

func parseReceiversFromCSV(fields [][]string) ([]api.Receiver, error) {
	var sends []api.Receiver
	var errs []error
	for i, f := range fields {
		addr := f[0]

		addr = strings.TrimSpace(addr)

		if _, err := cipher.DecodeBase58Address(addr); err != nil {
			err = fmt.Errorf("[row %d] Invalid address %s: %v", i, addr, err)
			errs = append(errs, err)
			continue
		}

		_, err := droplet.FromString(f[1])
		if err != nil {
			err = fmt.Errorf("[row %d] Invalid amount %s: %v", i, f[1], err)
			errs = append(errs, err)
			continue
		}

		sends = append(sends, api.Receiver{
			Address: addr,
			Coins:   f[1],
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

// createRawTxnArgs are encapsulated arguments for creating a transaction
type createRawTxnArgs struct {
	WalletID      string
	Address       string
	ChangeAddress string
	SendAmounts   []SendAmount
	Password      PasswordReader
}

func parseCreateRawTxnArgs(c *cobra.Command, args []string) (*createRawTxnArgs, error) {
	wltAddr, err := fromWalletOrAddress(c, args[0])
	if err != nil {
		return nil, err
	}

	changeAddress, err := c.Flags().GetString("change-address")
	if err != nil {
		return nil, err
	}
	chgAddr, err := getChangeAddress(wltAddr, changeAddress, apiTransactionsFinder{apiClient})
	if err != nil {
		return nil, err
	}

	toAddrs, err := getToAddresses(c, args[1:])
	if err != nil {
		return nil, err
	}
	if err := validateSendAmounts(toAddrs); err != nil {
		return nil, err
	}

	password, err := c.Flags().GetString("password")
	if err != nil {
		return nil, err
	}
	pr := NewPasswordReader([]byte(password))

	return &createRawTxnArgs{
		WalletID:      wltAddr.Wallet,
		Address:       wltAddr.Address,
		ChangeAddress: chgAddr,
		SendAmounts:   toAddrs,
		Password:      pr,
	}, nil
}

func createRawTxnCmdHandler(c *cobra.Command, args []string) (*coin.Transaction, error) {
	parsedArgs, err := parseCreateRawTxnArgs(c, args)
	if err != nil {
		return nil, err
	}

	to := make([]string, len(parsedArgs.SendAmounts))
	for i, s := range parsedArgs.SendAmounts {
		to[i] = s.Addr
	}
	warnSpentAddresses(apiClient, os.Stderr, to)

	// TODO -- load distribution params from config? Need to allow fiber chains to be used easily
	// There's too many distribution parameters to put them in command line, but we could read them from a file.
	// We could also have multiple hardcoded known distribution parameters for fiber coins, in the source,
	// but this wouldn't work for new fiber coins that hadn't been hardcoded yet.
	if parsedArgs.Address == "" {
		return CreateRawTxnFromWallet(apiClient, parsedArgs.WalletID,
			parsedArgs.ChangeAddress, parsedArgs.SendAmounts,
			parsedArgs.Password, params.MainNetDistribution)
	}

	return CreateRawTxnFromAddress(apiClient, parsedArgs.Address,
		parsedArgs.WalletID, parsedArgs.ChangeAddress, parsedArgs.SendAmounts,
		parsedArgs.Password, params.MainNetDistribution)
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

// CreateRawTxnFromWallet creates a transaction from any address or combination of addresses in a wallet
// An empty chgAddr sends the change to a new address added to the wallet file.
func CreateRawTxnFromWallet(c GetOutputser, walletFile, chgAddr string, toAddrs []SendAmount, pr PasswordReader, distParams params.Distribution) (*coin.Transaction, error) {
	if chgAddr != "" {
		if _, err := cipher.DecodeBase58Address(chgAddr); err != nil {
			return nil, ErrAddress
		}
	}

	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, err
	}

	if err := checkChangeAddress(wlt, chgAddr); err != nil {
		return nil, err
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
	totalAddrs, err := wlt.GetAddresses()
	if err != nil {
		return nil, err
	}

	addrStrArray := make([]string, len(totalAddrs))
	for i, a := range totalAddrs {
		addrStrArray[i] = a.String()
	}

	return createRawTxnWithChange(c, wlt, walletFile, addrStrArray, chgAddr, toAddrs, password, distParams)
}

// CreateRawTxnFromAddress creates a transaction from a specific address in a wallet
func CreateRawTxnFromAddress(c GetOutputser, addr, walletFile, chgAddr string, toAddrs []SendAmount, pr PasswordReader, distParams params.Distribution) (*coin.Transaction, error) {
	// check if the address is in the default wallet.
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, err
	}

	srcAddr, err := cipher.DecodeBase58Address(addr)
	if err != nil {
		return nil, ErrAddress
	}

	if _, err := wlt.GetEntry(srcAddr); err != nil {
		if err == wallet.ErrEntryNotFound {
			return nil, fmt.Errorf("%v address is not in wallet", addr)
		}
		return nil, err
	}

	if err := checkChangeAddress(wlt, chgAddr); err != nil {
		return nil, err
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

	return createRawTxnWithChange(c, wlt, walletFile, []string{addr}, chgAddr, toAddrs, password, distParams)
}

// checkChangeAddress checks that a named change address belongs to wlt.
func checkChangeAddress(wlt wallet.Wallet, chgAddr string) error {
	if chgAddr == "" {
		return nil
	}
	cAddr, err := cipher.DecodeBase58Address(chgAddr)
	if err != nil {
		return ErrAddress
	}
	if _, err := wlt.GetEntry(cAddr); err != nil {
		if err == wallet.ErrEntryNotFound {
			return fmt.Errorf("change address %v is not in wallet", chgAddr)
		}
		return err
	}
	return nil
}

// createRawTxnWithChange is CreateRawTxn that, given no change address, adds a new
// one to the wallet and saves the wallet file once the transaction is made.
func createRawTxnWithChange(c GetOutputser, wlt wallet.Wallet, walletFile string, inAddrs []string, chgAddr string, toAddrs []SendAmount, password []byte, distParams params.Distribution) (*coin.Transaction, error) {
	added := false
	if chgAddr == "" {
		var err error
		if chgAddr, added, err = newChangeAddress(wlt, password, inAddrs); err != nil {
			return nil, err
		}
	}

	txn, err := CreateRawTxn(c, wlt, inAddrs, chgAddr, toAddrs, password, distParams)
	if err != nil || !added {
		return txn, err
	}
	if err := saveWalletFile(wlt, walletFile); err != nil {
		return nil, err
	}
	return txn, nil
}

// saveWalletFile writes wlt back to the file it was loaded from.
func saveWalletFile(wlt wallet.Wallet, walletFile string) error {
	data, err := wlt.Serialize()
	if err != nil {
		return err
	}
	if err := file.SaveBinary(walletFile, data, 0600); err != nil {
		return fmt.Errorf("save the new address to %s: %v", walletFile, err)
	}
	return nil
}

// newChangeAddress adds an unused address to wlt for the change, so the change does
// not return to an address whose public key this spend reveals. A collection wallet
// cannot make one and keeps the old default, the first of inAddrs.
func newChangeAddress(wlt wallet.Wallet, password []byte, inAddrs []string) (string, bool, error) {
	opts := []wallet.Option{wallet.OptionGenerateN(1)}
	switch wlt.Type() {
	case wallet.WalletTypeCollection:
		if len(inAddrs) == 0 {
			return "", false, errors.New("no change address was found")
		}
		return inAddrs[0], false, nil
	case wallet.WalletTypeBip44:
		opts = append(opts, wallet.OptionChange())
	}

	var addrs []cipher.Addresser
	generate := func(w wallet.Wallet) error {
		var err error
		addrs, err = w.GenerateAddresses(opts...)
		return err
	}

	var err error
	if wlt.IsEncrypted() && wlt.Type() != wallet.WalletTypeBip44 {
		err = wallet.GuardUpdate(wlt, password, generate)
	} else {
		err = generate(wlt)
	}
	if err != nil {
		return "", false, err
	}
	return addrs[0].String(), true, nil
}

// GetOutputser implements unspent output querying
type GetOutputser interface {
	OutputsForAddresses([]string) (*readable.UnspentOutputsSummary, error)
}

// CreateRawTxn creates a transaction from a set of addresses contained in a loaded wallet.Wallet
func CreateRawTxn(c GetOutputser, wlt wallet.Wallet, inAddrs []string, chgAddr string, toAddrs []SendAmount, password []byte, distParams params.Distribution) (*coin.Transaction, error) {
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

	txn, err := createRawTxn(outputs, wlt, chgAddr, toAddrs, password)
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

	if err := transaction.VerifySingleTxnSoftConstraints(*txn, head.Time, inUxsFiltered, distParams, params.UserVerifyTxn); err != nil {
		return nil, err
	}
	if err := transaction.VerifySingleTxnHardConstraints(*txn, head, inUxsFiltered, transaction.TxnSigned); err != nil {
		return nil, err
	}
	if err := transaction.VerifySingleTxnUserConstraints(*txn); err != nil {
		return nil, err
	}

	return txn, nil
}

func createRawTxn(uxouts *readable.UnspentOutputsSummary, wlt wallet.Wallet, chgAddr string, toAddrs []SendAmount, password []byte) (*coin.Transaction, error) {
	// Calculate total required coins
	var totalCoins uint64
	for _, arg := range toAddrs {
		var err error
		totalCoins, err = mathutil.AddUint64(totalCoins, arg.Coins)
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

	f := func(w wallet.Wallet) (*coin.Transaction, error) {
		entries, err := getEntries(w, spendOutputs)
		if err != nil {
			return nil, err
		}

		return NewTransaction(spendOutputs, entries, txOuts)
	}

	makeTxn := func() (*coin.Transaction, error) {
		return f(wlt)
	}

	if wlt.IsEncrypted() {
		makeTxn = func() (*coin.Transaction, error) {
			var tx *coin.Transaction
			if err := wallet.GuardView(wlt, password, func(w wallet.Wallet) error {
				var err error
				tx, err = f(w)
				return err
			}); err != nil {
				return nil, err
			}

			return tx, nil
		}
	}

	return makeTxn()
}

func chooseSpends(uxouts *readable.UnspentOutputsSummary, coins uint64) ([]transaction.UxBalance, error) {
	// Convert spendable unspent outputs to []transaction.UxBalance
	spendableOutputs, err := readable.OutputsToUxBalances(uxouts.SpendableOutputs())
	if err != nil {
		return nil, err
	}

	// Choose which unspent outputs to spend
	// Use the MinimizeUxOuts strategy, since this is most likely used by
	// application that may need to send frequently.
	// Using fewer UxOuts will leave more available for other transactions,
	// instead of waiting for confirmation.
	outs, err := transaction.ChooseSpendsMinimizeUxOuts(spendableOutputs, coins, 0)
	if err != nil {
		// If there is not enough balance in the spendable outputs,
		// see if there is enough balance when including incoming outputs
		if err == transaction.ErrInsufficientBalance {
			expectedOutputs, otherErr := readable.OutputsToUxBalances(uxouts.ExpectedOutputs())
			if otherErr != nil {
				return nil, otherErr
			}

			if _, otherErr := transaction.ChooseSpendsMinimizeUxOuts(expectedOutputs, coins, 0); otherErr != nil {
				return nil, err
			}

			return nil, ErrTemporaryInsufficientBalance
		}

		return nil, err
	}

	return outs, nil
}

func makeChangeOut(outs []transaction.UxBalance, chgAddr string, toAddrs []SendAmount) ([]coin.TransactionOutput, error) {
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
		return nil, transaction.ErrInsufficientBalance
	}

	outAddrs := []coin.TransactionOutput{}
	changeAmount := totalInCoins - totalOutCoins

	haveChange := changeAmount > 0
	nAddrs := uint64(len(toAddrs))
	changeHours, addrHours, totalOutHours := transaction.DistributeSpendHours(totalInHours, nAddrs, haveChange)

	if err := fee.VerifyTransactionFeeForHours(totalOutHours, totalInHours-totalOutHours, params.UserVerifyTxn.BurnFactor); err != nil {
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

func getEntries(wlt wallet.Wallet, outs []transaction.UxBalance) ([]wallet.Entry, error) {
	entries := make([]wallet.Entry, len(outs))
	for i, o := range outs {
		entry, err := wlt.GetEntry(o.Address)
		if err != nil {
			if err == wallet.ErrEntryNotFound {
				return nil, fmt.Errorf("%v is not in wallet", o.Address.String())
			}
			return nil, err
		}
		entries[i] = entry
	}
	return entries, nil
}

// NewTransaction creates a transaction signed by the wallet entries that own utxos.
// The transaction should be validated against hard and soft constraints before transmission.
func NewTransaction(utxos []transaction.UxBalance, entries []wallet.Entry, outs []coin.TransactionOutput) (*coin.Transaction, error) {
	txn := coin.Transaction{}
	for _, u := range utxos {
		if err := txn.PushInput(u.Hash); err != nil {
			return nil, err
		}
	}

	for _, o := range outs {
		if err := txn.PushOutput(o.Address, o.Coins, o.Hours); err != nil {
			return nil, err
		}
	}

	txn.Sigs = make([]cipher.Sig, len(txn.In))
	if err := wallet.CheckSignedSize(&txn, utxos, params.UserVerifyTxn.MaxTransactionSize); err != nil {
		return nil, err
	}

	txn.InnerHash = txn.HashInner()
	for i, e := range entries {
		if err := wallet.SignInput(&txn, e, i); err != nil {
			return nil, err
		}
	}

	err := txn.UpdateHeader()
	if err != nil {
		return nil, err
	}

	return &txn, nil
}

// apiTransactionsFinder reports address activity using the node's API.
type apiTransactionsFinder struct {
	c *api.Client
}

// AddressesActivity implements wallet.TransactionsFinder
func (f apiTransactionsFinder) AddressesActivity(addrs []cipher.Addresser) ([]bool, error) {
	strs := make([]string, len(addrs))
	for i, a := range addrs {
		strs[i] = a.String()
	}
	txns, err := f.c.Transactions(strs)
	if err != nil {
		return nil, fmt.Errorf("bip44 change address lookup needs the node: %w", err)
	}
	active := make([]bool, len(addrs))
	for i, s := range strs {
		for _, t := range txns {
			if txnTouches(t, s) {
				active[i] = true
				break
			}
		}
	}
	return active, nil
}

func txnTouches(t readable.TransactionWithStatus, addr string) bool {
	for _, o := range t.Transaction.Out {
		if o.Address == addr {
			return true
		}
	}
	return false
}

// bip44ChangeAddress returns the next unused change address of a bip44
// wallet file, or "" for other wallet types. A new entry is saved.
func bip44ChangeAddress(path string, tf wallet.TransactionsFinder) (string, error) {
	wlt, err := wallet.Load(path)
	if err != nil {
		return "", WalletLoadError{err}
	}
	bw, ok := wlt.(*bip44wallet.Wallet)
	if !ok {
		return "", nil
	}
	before, err := bw.GetEntries(wallet.OptionChange())
	if err != nil {
		return "", err
	}
	addr, err := bw.PeekChangeAddress(tf)
	if err != nil {
		return "", err
	}
	after, err := bw.GetEntries(wallet.OptionChange())
	if err != nil {
		return "", err
	}
	if len(after) != len(before) {
		if err := wallet.Save(bw, filepath.Dir(path)); err != nil {
			return "", err
		}
	}
	return addr.String(), nil
}
