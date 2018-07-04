package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/fee"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/wallet"
)

// CreateTransactionResponse is returned by /wallet/transaction
type CreateTransactionResponse struct {
	Transaction        CreatedTransaction `json:"transaction"`
	EncodedTransaction string             `json:"encoded_transaction"`
}

// NewCreateTransactionResponse creates a CreateTransactionResponse
func NewCreateTransactionResponse(txn *coin.Transaction, inputs []wallet.UxBalance) (*CreateTransactionResponse, error) {
	cTxn, err := NewCreatedTransaction(txn, inputs)
	if err != nil {
		return nil, err
	}

	return &CreateTransactionResponse{
		Transaction:        *cTxn,
		EncodedTransaction: hex.EncodeToString(txn.Serialize()),
	}, nil
}

// CreatedTransaction represents a transaction created by /wallet/transaction
type CreatedTransaction struct {
	Length    uint32 `json:"length"`
	Type      uint8  `json:"type"`
	TxID      string `json:"txid"`
	InnerHash string `json:"inner_hash"`
	Fee       string `json:"fee"`

	Sigs []string                   `json:"sigs"`
	In   []CreatedTransactionInput  `json:"inputs"`
	Out  []CreatedTransactionOutput `json:"outputs"`
}

// NewCreatedTransaction returns a CreatedTransaction
func NewCreatedTransaction(txn *coin.Transaction, inputs []wallet.UxBalance) (*CreatedTransaction, error) {
	if len(txn.In) != len(inputs) {
		return nil, errors.New("len(txn.In) != len(inputs)")
	}

	var outputHours uint64
	for _, o := range txn.Out {
		var err error
		outputHours, err = coin.AddUint64(outputHours, o.Hours)
		if err != nil {
			return nil, err
		}
	}

	var inputHours uint64
	for _, i := range inputs {
		var err error
		inputHours, err = coin.AddUint64(inputHours, i.Hours)
		if err != nil {
			return nil, err
		}
	}

	if inputHours < outputHours {
		return nil, errors.New("inputHours unexpectedly less than output hours")
	}

	fee := inputHours - outputHours

	sigs := make([]string, len(txn.Sigs))
	for i, s := range txn.Sigs {
		sigs[i] = s.Hex()
	}

	txid := txn.Hash()
	out := make([]CreatedTransactionOutput, len(txn.Out))
	for i, o := range txn.Out {
		co, err := NewCreatedTransactionOutput(o, txid)
		if err != nil {
			return nil, err
		}
		out[i] = *co
	}

	in := make([]CreatedTransactionInput, len(inputs))
	for i, o := range inputs {
		ci, err := NewCreatedTransactionInput(o)
		if err != nil {
			return nil, err
		}
		in[i] = *ci
	}

	return &CreatedTransaction{
		Length:    txn.Length,
		Type:      txn.Type,
		TxID:      txid.Hex(),
		InnerHash: txn.InnerHash.Hex(),
		Fee:       fmt.Sprint(fee),

		Sigs: sigs,
		In:   in,
		Out:  out,
	}, nil
}

// ToTransaction converts a CreatedTransaction back to a coin.Transaction
func (r *CreatedTransaction) ToTransaction() (*coin.Transaction, error) {
	t := coin.Transaction{}

	t.Length = r.Length
	t.Type = r.Type

	var err error
	t.InnerHash, err = cipher.SHA256FromHex(r.InnerHash)
	if err != nil {
		return nil, err
	}

	sigs := make([]cipher.Sig, len(r.Sigs))
	for i, s := range r.Sigs {
		sigs[i], err = cipher.SigFromHex(s)
		if err != nil {
			return nil, err
		}
	}

	t.Sigs = sigs

	in := make([]cipher.SHA256, len(r.In))
	for i, n := range r.In {
		in[i], err = cipher.SHA256FromHex(n.UxID)
		if err != nil {
			return nil, err
		}
	}

	t.In = in

	out := make([]coin.TransactionOutput, len(r.Out))
	for i, o := range r.Out {
		addr, err := cipher.DecodeBase58Address(o.Address)
		if err != nil {
			return nil, err
		}

		coins, err := droplet.FromString(o.Coins)
		if err != nil {
			return nil, err
		}

		hours, err := strconv.ParseUint(o.Hours, 10, 64)
		if err != nil {
			return nil, err
		}

		out[i] = coin.TransactionOutput{
			Address: addr,
			Coins:   coins,
			Hours:   hours,
		}
	}

	t.Out = out

	hash, err := cipher.SHA256FromHex(r.TxID)
	if err != nil {
		return nil, err
	}
	if t.Hash() != hash {
		return nil, errors.New("ReadableTransaction.Hash does not match parsed transaction hash")
	}

	return &t, nil
}

// CreatedTransactionOutput is a transaction output
type CreatedTransactionOutput struct {
	UxID    string `json:"uxid"`
	Address string `json:"address"`
	Coins   string `json:"coins"`
	Hours   string `json:"hours"`
}

// NewCreatedTransactionOutput creates CreatedTransactionOutput
func NewCreatedTransactionOutput(out coin.TransactionOutput, txid cipher.SHA256) (*CreatedTransactionOutput, error) {
	coins, err := droplet.ToString(out.Coins)
	if err != nil {
		return nil, err
	}

	return &CreatedTransactionOutput{
		UxID:    out.UxID(txid).Hex(),
		Address: out.Address.String(),
		Coins:   coins,
		Hours:   fmt.Sprint(out.Hours),
	}, nil
}

// CreatedTransactionInput is a verbose transaction input
type CreatedTransactionInput struct {
	UxID            string `json:"uxid"`
	Address         string `json:"address,omitempty"`
	Coins           string `json:"coins,omitempty"`
	Hours           string `json:"hours,omitempty"`
	CalculatedHours string `json:"calculated_hours,omitempty"`
	Time            uint64 `json:"timestamp,omitempty"`
	Block           uint64 `json:"block,omitempty"`
	TxID            string `json:"txid,omitempty"`
}

// NewCreatedTransactionInput creates CreatedTransactionInput
func NewCreatedTransactionInput(out wallet.UxBalance) (*CreatedTransactionInput, error) {
	coins, err := droplet.ToString(out.Coins)
	if err != nil {
		return nil, err
	}

	if out.SrcTransaction.Null() {
		return nil, errors.New("NewCreatedTransactionInput UxOut.SrcTransaction is not initialized")
	}

	addr := out.Address.String()
	hours := fmt.Sprint(out.InitialHours)
	calculatedHours := fmt.Sprint(out.Hours)
	txID := out.SrcTransaction.Hex()

	return &CreatedTransactionInput{
		UxID:            out.Hash.Hex(),
		Address:         addr,
		Coins:           coins,
		Hours:           hours,
		CalculatedHours: calculatedHours,
		Time:            out.Time,
		Block:           out.BkSeq,
		TxID:            txID,
	}, nil
}

// createTransactionRequest is sent to /wallet/transaction
type createTransactionRequest struct {
	IgnoreUnconfirmed bool                           `json:"ignore_unconfirmed"`
	HoursSelection    hoursSelection                 `json:"hours_selection"`
	Wallet            createTransactionRequestWallet `json:"wallet"`
	ChangeAddress     *wh.Address                    `json:"change_address,omitempty"`
	To                []receiver                     `json:"to"`
}

// createTransactionRequestWallet defines a wallet to spend from and optionally which addresses in the wallet
type createTransactionRequestWallet struct {
	ID        string       `json:"id"`
	UxOuts    []wh.SHA256  `json:"unspents,omitempty"`
	Addresses []wh.Address `json:"addresses,omitempty"`
	Password  string       `json:"password"`
}

// hoursSelection defines options for hours distribution
type hoursSelection struct {
	Type        string           `json:"type"`
	Mode        string           `json:"mode"`
	ShareFactor *decimal.Decimal `json:"share_factor,omitempty"`
}

// receiver specifies a spend destination
type receiver struct {
	Address wh.Address `json:"address"`
	Coins   wh.Coins   `json:"coins"`
	Hours   *wh.Hours  `json:"hours,omitempty"`
}

// Validate validates createTransactionRequest data
func (r createTransactionRequest) Validate() error {
	switch r.HoursSelection.Type {
	case wallet.HoursSelectionTypeAuto:
		for i, to := range r.To {
			if to.Hours != nil {
				return fmt.Errorf("to[%d].hours must not be specified for auto hours_selection.mode", i)
			}
		}

		switch r.HoursSelection.Mode {
		case wallet.HoursSelectionModeShare:
		case "":
			return errors.New("missing hours_selection.mode")
		default:
			return errors.New("invalid hours_selection.mode")
		}

	case wallet.HoursSelectionTypeManual:
		for i, to := range r.To {
			if to.Hours == nil {
				return fmt.Errorf("to[%d].hours must be specified for manual hours_selection.mode", i)
			}
		}

		if r.HoursSelection.Mode != "" {
			return errors.New("hours_selection.mode cannot be used for manual hours_selection.type")
		}

	case "":
		return errors.New("missing hours_selection.type")
	default:
		return errors.New("invalid hours_selection.type")
	}

	if r.HoursSelection.ShareFactor == nil {
		if r.HoursSelection.Mode == wallet.HoursSelectionModeShare {
			return errors.New("missing hours_selection.share_factor when hours_selection.mode is share")
		}
	} else {
		if r.HoursSelection.Mode != wallet.HoursSelectionModeShare {
			return errors.New("hours_selection.share_factor can only be used when hours_selection.mode is share")
		}

		switch {
		case r.HoursSelection.ShareFactor.LessThan(decimal.New(0, 0)):
			return errors.New("hours_selection.share_factor cannot be negative")
		case r.HoursSelection.ShareFactor.GreaterThan(decimal.New(1, 0)):
			return errors.New("hours_selection.share_factor cannot be more than 1")
		}
	}

	if r.ChangeAddress != nil && r.ChangeAddress.Null() {
		return errors.New("change_address must not be the null address")
	}

	if r.Wallet.ID == "" {
		return errors.New("missing wallet.id")
	}

	addressMap := make(map[cipher.Address]struct{}, len(r.Wallet.Addresses))
	for i, a := range r.Wallet.Addresses {
		if a.Null() {
			return fmt.Errorf("wallet.addresses[%d] is empty", i)
		}

		addressMap[a.Address] = struct{}{}
	}

	if len(addressMap) != len(r.Wallet.Addresses) {
		return errors.New("wallet.addresses contains duplicate values")
	}

	if len(r.To) == 0 {
		return errors.New("to is empty")
	}

	for i, to := range r.To {
		if to.Address.Null() {
			return fmt.Errorf("to[%d].address is empty", i)
		}

		if to.Coins == 0 {
			return fmt.Errorf("to[%d].coins must not be zero", i)
		}

		if to.Coins.Value()%visor.MaxDropletDivisor() != 0 {
			return fmt.Errorf("to[%d].coins has too many decimal places", i)
		}
	}

	// Check for duplicate created outputs, a transaction can't have outputs with
	// the same (address, coins, hours)
	// Auto mode would distribute hours to the outputs and could hypothetically
	// avoid assigning duplicate hours in many cases, but the complexity for doing
	// so is very high, so also reject duplicate (address, coins) for auto mode.
	outputs := make(map[coin.TransactionOutput]struct{}, len(r.To))
	for _, to := range r.To {
		var hours uint64
		if to.Hours != nil {
			hours = to.Hours.Value()
		}

		outputs[coin.TransactionOutput{
			Address: to.Address.Address,
			Coins:   to.Coins.Value(),
			Hours:   hours,
		}] = struct{}{}
	}

	if len(outputs) != len(r.To) {
		return errors.New("to contains duplicate values")
	}

	if len(r.Wallet.UxOuts) != 0 && len(r.Wallet.Addresses) != 0 {
		return errors.New("wallet.unspents and wallet.addresses cannot be combined")
	}

	// Check for duplicate spending uxouts
	uxouts := make(map[cipher.SHA256]struct{}, len(r.Wallet.UxOuts))
	for _, o := range r.Wallet.UxOuts {
		uxouts[o.SHA256] = struct{}{}
	}

	if len(uxouts) != len(r.Wallet.UxOuts) {
		return errors.New("wallet.unspents contains duplicate values")
	}

	return nil
}

// ToWalletParams converts createTransactionRequest to wallet.CreateTransactionParams
func (r createTransactionRequest) ToWalletParams() wallet.CreateTransactionParams {
	addresses := make([]cipher.Address, len(r.Wallet.Addresses))
	for i, a := range r.Wallet.Addresses {
		addresses[i] = a.Address
	}

	uxouts := make([]cipher.SHA256, len(r.Wallet.UxOuts))
	for i, o := range r.Wallet.UxOuts {
		uxouts[i] = o.SHA256
	}

	walletParams := wallet.CreateTransactionWalletParams{
		ID:        r.Wallet.ID,
		Addresses: addresses,
		UxOuts:    uxouts,
		Password:  []byte(r.Wallet.Password),
	}

	to := make([]coin.TransactionOutput, len(r.To))
	for i, t := range r.To {
		var hours uint64
		if t.Hours != nil {
			hours = t.Hours.Value()
		}

		to[i] = coin.TransactionOutput{
			Address: t.Address.Address,
			Coins:   t.Coins.Value(),
			Hours:   hours,
		}
	}

	var changeAddress *cipher.Address
	if r.ChangeAddress != nil {
		changeAddress = &r.ChangeAddress.Address
	}

	return wallet.CreateTransactionParams{
		IgnoreUnconfirmed: r.IgnoreUnconfirmed,
		HoursSelection: wallet.HoursSelection{
			Type:        r.HoursSelection.Type,
			Mode:        r.HoursSelection.Mode,
			ShareFactor: r.HoursSelection.ShareFactor,
		},
		Wallet:        walletParams,
		ChangeAddress: changeAddress,
		To:            to,
	}
}

func createTransactionHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			wh.Error415(w)
			return
		}

		var params createTransactionRequest
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			logger.WithError(err).Error("Invalid create transaction request")
			wh.Error400(w, err.Error())
			return
		}

		if err := params.Validate(); err != nil {
			logger.WithError(err).Error("Invalid create transaction request")
			wh.Error400(w, err.Error())
			return
		}

		txn, inputs, err := gateway.CreateTransaction(params.ToWalletParams())
		if err != nil {
			switch err.(type) {
			case wallet.Error:
				switch err {
				case wallet.ErrWalletAPIDisabled:
					wh.Error403(w, "")
				case wallet.ErrWalletNotExist:
					wh.Error404(w, err.Error())
				default:
					wh.Error400(w, err.Error())
				}
			case blockdb.ErrUnspentNotExist:
				wh.Error400(w, err.Error())
			default:
				switch err {
				case fee.ErrTxnNoFee,
					fee.ErrTxnInsufficientCoinHours,
					wallet.ErrSpendingUnconfirmed:
					wh.Error400(w, err.Error())
				default:
					wh.Error500(w, err.Error())
				}
			}
			return
		}

		txnResp, err := NewCreateTransactionResponse(txn, inputs)
		if err != nil {
			err = fmt.Errorf("NewCreateTransactionResponse failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, txnResp)
	}
}
