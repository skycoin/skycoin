package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/fee"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/mathutil"
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
func NewCreateTransactionResponse(txn *coin.Transaction, inputs []visor.TransactionInput) (*CreateTransactionResponse, error) {
	cTxn, err := NewCreatedTransaction(txn, inputs)
	if err != nil {
		return nil, err
	}

	txnHex, err := txn.SerializeHex()
	if err != nil {
		return nil, err
	}

	return &CreateTransactionResponse{
		Transaction:        *cTxn,
		EncodedTransaction: txnHex,
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
func NewCreatedTransaction(txn *coin.Transaction, inputs []visor.TransactionInput) (*CreatedTransaction, error) {
	if len(txn.In) != len(inputs) {
		return nil, errors.New("len(txn.In) != len(inputs)")
	}

	var outputHours uint64
	for _, o := range txn.Out {
		var err error
		outputHours, err = mathutil.AddUint64(outputHours, o.Hours)
		if err != nil {
			return nil, err
		}
	}

	var inputHours uint64
	for _, i := range inputs {
		var err error
		inputHours, err = mathutil.AddUint64(inputHours, i.CalculatedHours)
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

	txID := txn.Hash()
	out := make([]CreatedTransactionOutput, len(txn.Out))
	for i, o := range txn.Out {
		co, err := NewCreatedTransactionOutput(o, txID)
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
		TxID:      txID.Hex(),
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
		return nil, fmt.Errorf("readable.Transaction.Hash %s does not match parsed transaction hash %s", t.Hash().Hex(), hash.Hex())
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
func NewCreatedTransactionInput(out visor.TransactionInput) (*CreatedTransactionInput, error) {
	coins, err := droplet.ToString(out.UxOut.Body.Coins)
	if err != nil {
		return nil, err
	}

	if out.UxOut.Body.SrcTransaction.Null() {
		return nil, errors.New("NewCreatedTransactionInput UxOut.SrcTransaction is not initialized")
	}

	addr := out.UxOut.Body.Address.String()
	hours := fmt.Sprint(out.UxOut.Body.Hours)
	calculatedHours := fmt.Sprint(out.CalculatedHours)
	txID := out.UxOut.Body.SrcTransaction.Hex()

	return &CreatedTransactionInput{
		UxID:            out.UxOut.Hash().Hex(),
		Address:         addr,
		Coins:           coins,
		Hours:           hours,
		CalculatedHours: calculatedHours,
		Time:            out.UxOut.Head.Time,
		Block:           out.UxOut.Head.BkSeq,
		TxID:            txID,
	}, nil
}

// createTransactionRequest is sent to POST /api/v2/transaction
type createTransactionRequest struct {
	IgnoreUnconfirmed bool           `json:"ignore_unconfirmed"`
	HoursSelection    hoursSelection `json:"hours_selection"`
	ChangeAddress     *wh.Address    `json:"change_address,omitempty"`
	To                []receiver     `json:"to"`
	UxOuts            []wh.SHA256    `json:"unspents,omitempty"`
	Addresses         []wh.Address   `json:"addresses,omitempty"`
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
	if r.ChangeAddress != nil && r.ChangeAddress.Null() {
		return errors.New("change_address must not be the null address")
	}

	switch r.HoursSelection.Type {
	case transaction.HoursSelectionTypeAuto:
		for i, to := range r.To {
			if to.Hours != nil {
				return fmt.Errorf("to[%d].hours must not be specified for auto hours_selection.mode", i)
			}
		}

		switch r.HoursSelection.Mode {
		case transaction.HoursSelectionModeShare:
		case "":
			return errors.New("missing hours_selection.mode")
		default:
			return errors.New("invalid hours_selection.mode")
		}

	case transaction.HoursSelectionTypeManual:
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
		if r.HoursSelection.Mode == transaction.HoursSelectionModeShare {
			return errors.New("missing hours_selection.share_factor when hours_selection.mode is share")
		}
	} else {
		if r.HoursSelection.Mode != transaction.HoursSelectionModeShare {
			return errors.New("hours_selection.share_factor can only be used when hours_selection.mode is share")
		}

		switch {
		case r.HoursSelection.ShareFactor.LessThan(decimal.New(0, 0)):
			return errors.New("hours_selection.share_factor cannot be negative")
		case r.HoursSelection.ShareFactor.GreaterThan(decimal.New(1, 0)):
			return errors.New("hours_selection.share_factor cannot be more than 1")
		}
	}

	if len(r.UxOuts) != 0 && len(r.Addresses) != 0 {
		return errors.New("unspents and addresses cannot be combined")
	}

	addressMap := make(map[cipher.Address]struct{}, len(r.Addresses))
	for i, a := range r.Addresses {
		if a.Null() {
			return fmt.Errorf("addresses[%d] is empty", i)
		}

		if _, ok := addressMap[a.Address]; ok {
			return errors.New("addresses contains duplicate values")
		}

		addressMap[a.Address] = struct{}{}
	}

	// Check for duplicate spending uxouts
	uxouts := make(map[cipher.SHA256]struct{}, len(r.UxOuts))
	for _, o := range r.UxOuts {
		if _, ok := uxouts[o.SHA256]; ok {
			return errors.New("unspents contains duplicate values")
		}

		uxouts[o.SHA256] = struct{}{}
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

		if to.Coins.Value()%params.UserVerifyTxn.MaxDropletDivisor() != 0 {
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

		txo := coin.TransactionOutput{
			Address: to.Address.Address,
			Coins:   to.Coins.Value(),
			Hours:   hours,
		}

		if _, ok := outputs[txo]; ok {
			return errors.New("to contains duplicate values")
		}

		outputs[txo] = struct{}{}
	}

	return nil
}

// TransactionParams converts createTransactionRequest to transaction.Params
func (r createTransactionRequest) TransactionParams() transaction.Params {
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

	return transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type:        r.HoursSelection.Type,
			Mode:        r.HoursSelection.Mode,
			ShareFactor: r.HoursSelection.ShareFactor,
		},
		ChangeAddress: changeAddress,
		To:            to,
	}
}

func (r createTransactionRequest) VisorParams() visor.CreateTransactionParams {
	return visor.CreateTransactionParams{
		IgnoreUnconfirmed: r.IgnoreUnconfirmed,
		Addresses:         r.addresses(),
		UxOuts:            r.uxOuts(),
	}
}

func (r createTransactionRequest) addresses() []cipher.Address {
	if len(r.Addresses) == 0 {
		return nil
	}
	addresses := make([]cipher.Address, len(r.Addresses))
	for i, a := range r.Addresses {
		addresses[i] = a.Address
	}
	return addresses
}

func (r createTransactionRequest) uxOuts() []cipher.SHA256 {
	if len(r.UxOuts) == 0 {
		return nil
	}
	uxouts := make([]cipher.SHA256, len(r.UxOuts))
	for i, o := range r.UxOuts {
		uxouts[i] = o.SHA256
	}
	return uxouts
}

// transactionHandlerV2 creates a transaction from provided outputs and parameters
// Method: POST
// URI: /api/v2/transaction
// Args: JSON body
func transactionHandlerV2(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req createTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		if err := req.Validate(); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		// Check that addresses or unspents are not empty
		// This is not checked in Validate() because POST /api/v1/wallet/transaction
		// allows both to be empty
		if len(req.Addresses) == 0 && len(req.UxOuts) == 0 {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "one of addresses or unspents must not be empty")
			writeHTTPResponse(w, resp)
			return
		}

		txn, inputs, err := gateway.CreateTransaction(req.TransactionParams(), req.VisorParams())
		if err != nil {
			var resp HTTPResponse
			switch err.(type) {
			case blockdb.ErrUnspentNotExist, transaction.Error, visor.UserError, wallet.Error:
				resp = NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			default:
				switch err {
				case fee.ErrTxnNoFee, fee.ErrTxnInsufficientCoinHours:
					resp = NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
				default:
					resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				}
			}
			writeHTTPResponse(w, resp)
			return
		}

		txnResp, err := NewCreateTransactionResponse(txn, inputs)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("NewCreateTransactionResponse failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: txnResp,
		})
	}
}

// walletCreateTransactionRequest is sent to POST /api/v1/wallet/transaction
type walletCreateTransactionRequest struct {
	Unsigned bool   `json:"unsigned"`
	WalletID string `json:"wallet_id"`
	Password string `json:"password"`
	createTransactionRequest
}

// Validate validates walletCreateTransactionRequest data
func (r walletCreateTransactionRequest) Validate() error {
	if r.WalletID == "" {
		return errors.New("missing wallet_id")
	}

	if r.Unsigned && len(r.Password) != 0 {
		return errors.New("password must not be used for unsigned transactions")
	}

	return r.createTransactionRequest.Validate()
}

// walletCreateTransactionHandler creates a transaction
// Method: POST
// URI: /api/v1/wallet/transaction
// Args: JSON body
func walletCreateTransactionHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		if !isContentTypeJSON(r.Header.Get("Content-Type")) {
			wh.Error415(w)
			return
		}

		var req walletCreateTransactionRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.WithError(err).Error("Invalid create transaction request")
			wh.Error400(w, err.Error())
			return
		}

		if err := req.Validate(); err != nil {
			logger.WithError(err).Error("Invalid create transaction request")
			wh.Error400(w, err.Error())
			return
		}

		var txn *coin.Transaction
		var inputs []visor.TransactionInput
		if req.Unsigned {
			txn, inputs, err = gateway.WalletCreateTransaction(req.WalletID, req.TransactionParams(), req.VisorParams())
		} else {
			txn, inputs, err = gateway.WalletCreateTransactionSigned(req.WalletID, []byte(req.Password), req.TransactionParams(), req.VisorParams())
		}
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
			case blockdb.ErrUnspentNotExist,
				transaction.Error,
				visor.UserError:
				wh.Error400(w, err.Error())
			default:
				switch err {
				case fee.ErrTxnNoFee,
					fee.ErrTxnInsufficientCoinHours:
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

// WalletSignTransactionRequest is the request body object for /api/v2/wallet/transaction/sign
type WalletSignTransactionRequest struct {
	WalletID           string `json:"wallet_id"`
	Password           string `json:"password"`
	EncodedTransaction string `json:"encoded_transaction"`
	SignIndexes        []int  `json:"sign_indexes"`
}

// walletSignTransactionHandler signs an unsigned transaction
// Method: POST
// URI: /api/v2/wallet/transaction/sign
// Args: JSON body
func walletSignTransactionHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req WalletSignTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		if req.WalletID == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "wallet_id is required")
			writeHTTPResponse(w, resp)
			return
		}

		if req.EncodedTransaction == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "encoded_transaction is required")
			writeHTTPResponse(w, resp)
			return
		}

		txn, err := decodeTxn(req.EncodedTransaction)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("Decode transaction failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		// Check that number of sign_indexes does not exceed number of inputs
		if len(req.SignIndexes) > len(txn.In) {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "Too many values in sign_indexes")
			writeHTTPResponse(w, resp)
			return
		}

		// Check that values in sign_indexes are in the range of txn inputs
		for _, i := range req.SignIndexes {
			if i < 0 || i >= len(txn.In) {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, "Value in sign_indexes exceeds range of transaction inputs array")
				writeHTTPResponse(w, resp)
				return
			}
		}

		// Check for duplicate values in sign_indexes
		signIndexesMap := make(map[int]struct{}, len(req.SignIndexes))
		for _, i := range req.SignIndexes {
			if _, ok := signIndexesMap[i]; ok {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, "Duplicate value in sign_indexes")
				writeHTTPResponse(w, resp)
				return
			}
			signIndexesMap[i] = struct{}{}
		}

		signedTxn, inputs, err := gateway.WalletSignTransaction(req.WalletID, []byte(req.Password), txn, req.SignIndexes)
		if err != nil {
			var resp HTTPResponse
			switch err.(type) {
			case wallet.Error:
				switch err {
				case wallet.ErrWalletNotExist:
					resp = NewHTTPErrorResponse(http.StatusNotFound, err.Error())
				case wallet.ErrWalletAPIDisabled:
					resp = NewHTTPErrorResponse(http.StatusForbidden, err.Error())
				default:
					resp = NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
				}
			case visor.ErrTxnViolatesSoftConstraint,
				visor.ErrTxnViolatesHardConstraint,
				visor.ErrTxnViolatesUserConstraint,
				blockdb.ErrUnspentNotExist:
				resp = NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			default:
				resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			}
			writeHTTPResponse(w, resp)
			return
		}

		txnResp, err := NewCreateTransactionResponse(signedTxn, inputs)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: txnResp,
		})
	}
}
