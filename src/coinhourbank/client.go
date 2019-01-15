package coinhourbank

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/api"
)

// HourBankClient provides simplified access to the coin hour bank
type HourBankClient struct {
	bankURL string
	client  *api.Client
	wallet  BankWallet
	address SourceAddress
}

// NewHourBankClient creates a new instance of HourBankClient
func NewHourBankClient(
	nodeURL,
	bankURL,
	walletPath string,
	walletPassword []byte,
	walletSourceAddress SourceAddress,
) (*HourBankClient, error) {

	skyClient := api.NewClient(nodeURL)
	wallet, err := NewBankWallet(skyClient, walletPath, &walletPassword)
	if err != nil {
		return nil, err
	}

	return newHourBankClient(skyClient, bankURL, wallet, walletSourceAddress)
}

func newHourBankClient(
	client *api.Client,
	bankURL string,
	wallet BankWallet,
	address SourceAddress,
) (*HourBankClient, error) {

	return &HourBankClient{
		bankURL: bankURL,
		client:  client,
		wallet:  wallet,
		address: address,
	}, nil
}

// Balance returns balance for current account
func (c *HourBankClient) Balance(account Account) (CoinHours, error) {
	balanceURL := fmt.Sprintf("%s/api/account/%s/balance", c.bankURL, account)
	res, err := http.Get(balanceURL) /* #nosec */
	if err != nil {
		return CoinHours(0), err
	}

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return CoinHours(0), err
		}

		bodyString := string(bodyBytes)
		coinHoursInt, err := strconv.ParseInt(bodyString, 10, 64)
		if err != nil {
			return CoinHours(0), err
		}

		return CoinHours(coinHoursInt), nil
	}

	return CoinHours(0), fmt.Errorf("http response is not 200: %d", res.StatusCode)
}

// DepositHours puts SKY coin hours into the bank
func (c *HourBankClient) DepositHours(hours CoinHours, account Account) error {
	outputs, err := c.wallet.Outputs([]string{string(c.address)})
	if err != nil {
		return err
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


	tx, inputs, err := c.wallet.SignTransaction(signableTransaction.Transaction, signableTransaction.Inputs)
	if err != nil {
		return err
	}

	depositReq := SignableTransaction{
		Transaction: tx,
		Inputs:      inputs,
	}

	mb, err = modelToBytes(depositReq)
	if err != nil {
		return err
	}

	depositTransactionURL = c.bankURL + "/api/account/" + string(account) + "/deposit"
	_, err = http.Post(depositTransactionURL, "application/json", mb) /* #nosec */

	return err
}

// TransferHours performs hours transfer from one account to another
func (c *HourBankClient) TransferHours(accountFrom Account, accountTo Account, amount CoinHours) error {
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

// WithdrawHours withdraws SKY coin hours from the bank
func (c *HourBankClient) WithdrawHours(hours CoinHours, account Account) error {
	outputs, err := c.wallet.Outputs([]string{string(c.address)})
	if err != nil {
		return err
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

	tx, inputs, err := c.wallet.SignTransaction(signableTransaction.Transaction, signableTransaction.Inputs)
	if err != nil {
		return err
	}

	withdrawReq := SignableTransaction{
		Transaction: tx,
		Inputs:      inputs,
	}

	mb, err = modelToBytes(withdrawReq)
	if err != nil {
		return err
	}

	depositTransactionURL = c.bankURL + "/api/account/" + string(account) + "/withdraw"
	_, err = http.Post(depositTransactionURL, "application/json", mb) /* #nosec */

	return err
}
