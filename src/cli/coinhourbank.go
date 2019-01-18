package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
)

func coinhourBalanceCmd() *cobra.Command {
	coinhourBalanceCmd := &cobra.Command{
		Use:          "coinhourBalance",
		Short:        "Get balance of coinhour bank account",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			bankClient, err := getCoinhourBankClient(c)
			if err != nil {
				return err
			}

			address, err := c.Flags().GetString("address")
			if err != nil {
				return err
			}

			if _, err := cipher.DecodeBase58Address(address); err != nil {
				return err
			}

			balance, err := bankClient.Balance(address)
			if err != nil {
				return err
			}

			fmt.Printf("%s balance: %v\n", address, balance)

			return nil
		},
	}

	coinhourBalanceCmd.Flags().StringP("address", "a", "", "wallet address to take coinhours from")
	coinhourBalanceCmd.Flags().StringP("bankURL", "b", "http://localhost:8081", "coinhour bank backend url")

	return coinhourBalanceCmd
}

func depositCoinhoursCmd() *cobra.Command {
	depositCoinhoursCmd := &cobra.Command{
		Use:   "depositCoinhours [hours amount]",
		Short: "Sends coinhours to a coinhour bank account.",
		Long: `Deposits coinhours into a coinhour bank account which a skycoin address you want to deposit hours into.
		Once hours are into coinhour bank they can be transferred to other addresses without paying transaction fee.`,
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			bankClient, err := getCoinhourBankClient(c)
			if err != nil {
				return err
			}

			wlt, err := getWallet(c)
			if err != nil {
				return err
			}
			defer wlt.Erase()

			address, _ := c.Flags().GetString("address") // nolint: errcheck

			coinhours, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			return bankClient.DepositHours(coinhours, address, wlt)
		},
	}

	depositCoinhoursCmd.Flags().StringP("wallet-file", "f", "", "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.")
	depositCoinhoursCmd.Flags().StringP("password", "p", "", "wallet password")
	depositCoinhoursCmd.Flags().StringP("address", "a", "", "wallet address to take coinhours from")
	depositCoinhoursCmd.Flags().StringP("nodeURL", "n", "http://localhost:6420", "skycoin node url")
	depositCoinhoursCmd.Flags().StringP("bankURL", "b", "http://localhost:8081", "coinhour bank backend url")

	return depositCoinhoursCmd
}

func transferCoinhoursCmd() *cobra.Command {
	transferCoinhoursCmd := &cobra.Command{
		Use:          "transferCoinhours [destination address] [hours amount]",
		Short:        "Transfer coinhours from one coinhour bank account to another",
		Long:         `Transferring coinhours from one account to another does not require any transaction fee.`,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			bankClient, err := getCoinhourBankClient(c)
			if err != nil {
				return err
			}

			address, err := c.Flags().GetString("address")
			if err != nil {
				return err
			}
			if _, err := cipher.DecodeBase58Address(address); err != nil {
				return err
			}

			coinhours, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			return bankClient.TransferHours(address, args[0], coinhours)
		},
	}

	transferCoinhoursCmd.Flags().StringP("bankURL", "b", "http://localhost:8081", "coinhour bank backend url")

	return transferCoinhoursCmd
}

func withdrawCoinhoursCmd() *cobra.Command {
	withdrawCoinhoursCmd := &cobra.Command{
		Use:          "withdrawCoinhours [hours amount]",
		Short:        "Withdraws coinhours from coinhour bank.",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			bankClient, err := getCoinhourBankClient(c)
			if err != nil {
				return err
			}

			wlt, err := getWallet(c)
			if err != nil {
				return err
			}

			address, _ := c.Flags().GetString("address") // nolint: errcheck

			coinhours, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			return bankClient.WithdrawHours(coinhours, address, wlt)
		},
	}

	withdrawCoinhoursCmd.Flags().StringP("wallet-file", "f", "", "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.")
	withdrawCoinhoursCmd.Flags().StringP("password", "p", "", "wallet password")
	withdrawCoinhoursCmd.Flags().StringP("address", "a", "", "wallet address to take coinhours from")
	withdrawCoinhoursCmd.Flags().StringP("nodeURL", "n", "http://localhost:6420", "skycoin node url")
	withdrawCoinhoursCmd.Flags().StringP("bankURL", "b", "http://localhost:8081", "coinhour bank backend url")

	return withdrawCoinhoursCmd
}

// coinhour bank client

// HourBankClient provides simplified access to the coin hour bank
type HourBankClient struct {
	bankURL string
}

// unspentOutput is UxOuts of transaction
type unspentOutput struct {
	Hash    string // Hash of unspent output
	Address string // Address of receiver
	Coins   uint64 // Number of coins
	Hours   uint64 // Coin hours
}

// NewHourBankClient creates a new instance of HourBankClient
func NewHourBankClient(bankURL string) *HourBankClient {
	return &HourBankClient{
		bankURL: bankURL,
	}
}

// Balance returns balance for current account
func (c *HourBankClient) Balance(account string) (uint64, error) {
	balanceURL := fmt.Sprintf("%s/api/account/%s/balance", c.bankURL, account)
	res, err := http.Get(balanceURL) /* #nosec */
	if err != nil {
		return uint64(0), err
	}

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return uint64(0), err
		}

		bodyString := string(bodyBytes)
		coinHoursInt, err := strconv.ParseInt(bodyString, 10, 64)
		if err != nil {
			return uint64(0), err
		}

		return uint64(coinHoursInt), nil
	}

	return uint64(0), fmt.Errorf("http response is not 200: %d", res.StatusCode)
}

// DepositHoursRequest models request for depositing coin hours
type DepositHoursRequest struct {
	UnspentOutputs []unspentOutput `json:"unspentOutput"`
	CoinHours      uint64          `json:"coinHours"`
}

type coinjoinInput struct {
	FromAddress string     // Address of transaction sender
	UxOut       string     // UxOut HEX of transaction sender
	Sign        cipher.Sig // Signature of input
	InputIndex  int        // Index of this input in the transaction inputs array
}

// SignableTransaction is the entity for signing and publishing a transaction
type SignableTransaction struct {
	Transaction coin.Transaction `json:"transaction"`
	Inputs      []coinjoinInput  `json:"inputs"`
}

// DepositHours puts SKY coin hours into the bank
func (c *HourBankClient) DepositHours(hours uint64, account string, wallet *wallet.Wallet) error {
	outputSet, err := apiClient.OutputsForAddresses([]string{string(account)})
	if err != nil {
		return err
	}

	ux, err := outputSet.HeadOutputs.ToUxArray()
	if err != nil {
		return err
	}

	outputs := make([]unspentOutput, len(outputSet.HeadOutputs))
	for i, ho := range outputSet.HeadOutputs {
		outputs[i] = unspentOutput{
			Hash:    ho.Hash,
			Address: ho.Address,
			Coins:   ux[i].Body.Coins,
			Hours:   ho.CalculatedHours,
		}
	}

	req := DepositHoursRequest{
		CoinHours:      hours,
		UnspentOutputs: outputs,
	}

	mb, err := modelToBytes(req)
	if err != nil {
		return err
	}

	depositTransactionURL := fmt.Sprintf("%s/api/transactions/deposit", c.bankURL)
	res, err := http.Post(depositTransactionURL, "application/json", mb) /* #nosec */
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var signableTransaction SignableTransaction
	if res.StatusCode == 200 {
		if err := json.NewDecoder(res.Body).Decode(&signableTransaction); err != nil {
			return err
		}
	} else {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}

	for _, input := range signableTransaction.Inputs {
		addr, err := cipher.DecodeBase58Address(input.FromAddress)
		if err != nil {
			return err
		}
		wltEntry, ok := wallet.GetEntry(addr)
		if ok {
			if err != nil {
				return err
			}
			err = signAddressInputs(addr, signableTransaction.Transaction.InnerHash, wltEntry.Secret.Hex(), signableTransaction.Inputs)
			if err != nil {
				return err
			}
		}
	}

	sigs := make([]cipher.Sig, len(signableTransaction.Transaction.In))
	for _, in := range signableTransaction.Inputs {
		sigs[in.InputIndex] = in.Sign
	}

	signableTransaction.Transaction.Sigs = sigs

	err = signableTransaction.Transaction.UpdateHeader()
	if err != nil {
		return err
	}

	depositReq := SignableTransaction{
		Transaction: signableTransaction.Transaction,
		Inputs:      signableTransaction.Inputs,
	}

	mb, err = modelToBytes(depositReq)
	if err != nil {
		return err
	}

	depositTransactionURL = c.bankURL + "/api/account/" + string(account) + "/deposit"
	_, err = http.Post(depositTransactionURL, "application/json", mb) /* #nosec */

	return err
}

// TransferHoursRequest represents a body of the TransferHours request
type TransferHoursRequest struct {
	AddressTo string `json:"addressTo"`
	Amount    int64  `json:"amount"`
}

// TransferHours performs hours transfer from one account to another
func (c *HourBankClient) TransferHours(accountFrom, accountTo string, amount uint64) error {
	req := TransferHoursRequest{
		AddressTo: string(accountTo),
		Amount:    int64(amount),
	}

	mb, err := modelToBytes(req)
	if err != nil {
		return err
	}

	transferCoinHoursURL := c.bankURL + "/api/account/" + string(accountFrom) + "/transfer"
	_, err = http.Post(transferCoinHoursURL, "application/json", mb) /* #nosec */

	return err
}

// WithdrawHoursRequest models request for withdrawing coin hours
type WithdrawHoursRequest struct {
	Address        string          `json:"address"`
	UnspentOutputs []unspentOutput `json:"unspentOutput"`
	CoinHours      uint64          `json:"coinHours"`
}

// WithdrawHours withdraws SKY coin hours from the bank
func (c *HourBankClient) WithdrawHours(hours uint64, account string, wallet *wallet.Wallet) error {
	outputSet, err := apiClient.OutputsForAddresses([]string{string(account)})
	if err != nil {
		return err
	}

	ux, err := outputSet.HeadOutputs.ToUxArray()
	if err != nil {
		return err
	}

	outputs := make([]unspentOutput, len(outputSet.HeadOutputs))
	for i, ho := range outputSet.HeadOutputs {
		outputs[i] = unspentOutput{
			Hash:    ho.Hash,
			Address: ho.Address,
			Coins:   ux[i].Body.Coins,
			Hours:   ho.CalculatedHours,
		}
	}

	req := WithdrawHoursRequest{
		Address:        string(account),
		CoinHours:      hours,
		UnspentOutputs: outputs,
	}

	mb, err := modelToBytes(req)
	if err != nil {
		return err
	}

	depositTransactionURL := fmt.Sprintf("%s/api/transactions/withdraw", c.bankURL)
	res, err := http.Post(depositTransactionURL, "application/json", mb) /* #nosec */
	if err != nil {
		return err
	}

	var signableTransaction SignableTransaction
	if res.StatusCode == 200 {
		if err := json.NewDecoder(res.Body).Decode(&signableTransaction); err != nil {
			return err
		}
	} else {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return errors.New(string(body))
	}

	for _, input := range signableTransaction.Inputs {
		addr, err := cipher.DecodeBase58Address(input.FromAddress)
		if err != nil {
			return err
		}
		wltEntry, ok := wallet.GetEntry(addr)
		if ok {
			if err != nil {
				return err
			}
			err = signAddressInputs(addr, signableTransaction.Transaction.InnerHash, wltEntry.Secret.Hex(), signableTransaction.Inputs)
			if err != nil {
				return err
			}
		}
	}

	sigs := make([]cipher.Sig, len(signableTransaction.Transaction.In))
	for _, in := range signableTransaction.Inputs {
		sigs[in.InputIndex] = in.Sign
	}
	signableTransaction.Transaction.Sigs = sigs

	err = signableTransaction.Transaction.UpdateHeader()
	if err != nil {
		return err
	}

	withdrawReq := SignableTransaction{
		Transaction: signableTransaction.Transaction,
		Inputs:      signableTransaction.Inputs,
	}

	mb, err = modelToBytes(withdrawReq)
	if err != nil {
		return err
	}

	depositTransactionURL = c.bankURL + "/api/account/" + string(account) + "/withdraw"
	_, err = http.Post(depositTransactionURL, "application/json", mb) /* #nosec */

	return err
}

func getCoinhourBankClient(c *cobra.Command) (*HourBankClient, error) {
	bankURL, err := c.Flags().GetString("bankURL")
	if err != nil {
		return nil, err
	}

	bankClient := NewHourBankClient(bankURL)
	return bankClient, nil
}

func getWallet(c *cobra.Command) (*wallet.Wallet, error) {
	walletFile, err := c.Flags().GetString("wallet-file")
	if err != nil {
		return nil, nil
	}

	wltPath, err := resolveWalletPath(cliConfig, walletFile)
	if err != nil {
		return nil, err
	}

	wlt, err := wallet.Load(wltPath)
	if err != nil {
		return nil, err
	}

	address, err := c.Flags().GetString("address")
	if err != nil {
		return nil, err
	}

	sourceAddr, err := cipher.DecodeBase58Address(address)
	if err != nil {
		return nil, err
	}

	if _, ok := wlt.GetEntry(sourceAddr); !ok {
		return nil, fmt.Errorf("sender address not in wallet")
	}

	password, err := c.Flags().GetString("password")
	if err != nil {
		return nil, err
	}
	pr := NewPasswordReader([]byte(password))

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

	if wlt.IsEncrypted() {
		p, err := pr.Password()
		if err != nil {
			return nil, err
		}

		return wlt.Unlock(p)
	}

	return wlt, nil
}

func modelToBytes(i interface{}) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	return b, json.NewEncoder(b).Encode(i)
}

// signAddressInputs signs all UxOuts from tx inputs (inputs) that corresponds to specified address by its secret key.
func signAddressInputs(address cipher.Address, txInnerHash cipher.SHA256, secKeyHex string, inputs []coinjoinInput) error {
	secKey, err := cipher.SecKeyFromHex(secKeyHex)
	if err != nil {
		return err
	}

	for i, in := range inputs {
		if in.FromAddress == address.String() {
			h := cipher.AddSHA256(txInnerHash, cipher.MustSHA256FromHex(in.UxOut)) // hash to sign
			inputs[i].Sign, err = cipher.SignHash(h, secKey)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
