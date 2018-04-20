package gui

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/fee"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

// CreateTransactionResponse is returned by /wallet/transaction
type CreateTransactionResponse struct {
	Transaction        visor.ReadableTransaction `json:"transaction"`
	EncodedTransaction string                    `json:"encoded_transaction"`
}

// createTransactionRequest is sent to /wallet/transaction
type createTransactionRequest struct {
	HoursSelection hoursSelection                 `json:"hours_selection"`
	Wallet         createTransactionRequestWallet `json:"wallet"`
	ChangeAddress  *wh.Address                    `json:"change_address"`
	To             []receiver                     `json:"to"`
}

// createTransactionRequestWallet defines a wallet to spend from and optionally which addresses in the wallet
type createTransactionRequestWallet struct {
	ID        string       `json:"id"`
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

	if r.ChangeAddress == nil {
		return errors.New("missing change_address")
	} else if r.ChangeAddress.Empty() {
		return errors.New("change_address is an empty address")
	}

	if r.Wallet.ID == "" {
		return errors.New("missing wallet.id")
	}

	for i, a := range r.Wallet.Addresses {
		if a.Empty() {
			return fmt.Errorf("wallet.addresses[%d] is empty", i)
		}
	}

	if len(r.To) == 0 {
		return errors.New("to is empty")
	}

	for i, to := range r.To {
		if to.Address.Empty() {
			return fmt.Errorf("to[%d].address is empty", i)
		}

		if to.Coins == 0 {
			return fmt.Errorf("to[%d].coins must not be zero", i)
		}

		if to.Coins.Value()%visor.MaxDropletDivisor() != 0 {
			return fmt.Errorf("to[%d].coins has too many decimal places", i)
		}
	}

	// Check for duplicate outputs, a transaction can't have outputs with
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

	return nil
}

// ToWalletParams converts createTransactionRequest to wallet.CreateTransactionParams
func (r createTransactionRequest) ToWalletParams() wallet.CreateTransactionParams {
	addresses := make([]cipher.Address, len(r.Wallet.Addresses))
	for i, a := range r.Wallet.Addresses {
		addresses[i] = a.Address
	}

	walletParams := wallet.CreateTransactionWalletParams{
		ID:        r.Wallet.ID,
		Addresses: addresses,
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

	var changeAddress cipher.Address
	if r.ChangeAddress != nil {
		changeAddress = r.ChangeAddress.Address
	}

	return wallet.CreateTransactionParams{
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

		txn, err := gateway.CreateTransaction(params.ToWalletParams())
		if err != nil {
			switch err.(type) {
			case wallet.Error:
				switch err {
				case wallet.ErrWalletAPIDisabled:
					wh.Error403(w)
				case wallet.ErrWalletNotExist:
					wh.Error404Msg(w, err.Error())
				default:
					wh.Error400(w, err.Error())
				}
			default:
				switch err {
				case fee.ErrTxnNoFee, fee.ErrTxnInsufficientCoinHours:
					wh.Error400(w, err.Error())
				default:
					wh.Error500Msg(w, err.Error())
				}
			}
			return
		}

		readableTxn, err := visor.NewReadableTransaction(&visor.Transaction{
			Txn: *txn,
		})
		if err != nil {
			err = fmt.Errorf("visor.NewReadableTransaction failed: %v", err)
			logger.WithError(err).Error()
			wh.Error500Msg(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, CreateTransactionResponse{
			Transaction:        *readableTxn,
			EncodedTransaction: hex.EncodeToString(txn.Serialize()),
		})
	}
}
