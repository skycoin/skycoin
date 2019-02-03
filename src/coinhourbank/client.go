package coinhourbank

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
)

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

// DepositHours puts SKY coin hours into the bank
func (c *HourBankClient) DepositHours(hours uint64, account string, spendableOutputSet coin.UxArray, wallet *wallet.Wallet) error {
	outputs := make([]unspentOutput, len(spendableOutputSet))
	for i, ho := range spendableOutputSet {
		outputs[i] = unspentOutput{
			Hash:    ho.Hash().Hex(),
			Address: ho.Body.Address.String(),
			Coins:   ho.Body.Coins,
			Hours:   ho.Body.Hours,
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
