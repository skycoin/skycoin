package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/skycoin/hardware-wallet-go/src/skywallet/wire"

	"github.com/gogo/protobuf/proto"
	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
	messages "github.com/skycoin/hardware-wallet-protob/go"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/droplet"
)

// TransactionSignRequest is request data for /api/v1/transaction_sign
type TransactionSignRequest struct {
	TransactionInputs  []TransactionInput  `json:"transaction_inputs"`
	TransactionOutputs []TransactionOutput `json:"transaction_outputs"`
}

// TransactionInput is a skycoin transaction input
type TransactionInput struct {
	Index *uint32 `json:"index"` // pointer to differentiate between 0 and nil
	Hash  string  `json:"hash"`
}

// TransactionOutput is a skycoin transaction output
type TransactionOutput struct {
	AddressIndex *uint32 `json:"address_index"` // pointer to differentiate between 0 and nil
	Address      string  `json:"address"`
	Coins        string  `json:"coins"`
	Hours        string  `json:"hours"`
}

// TransactionSignResponse is data returned by POST /api/v1/transaction_sign
type TransactionSignResponse struct {
	Signatures *[]string `json:"signatures"`
}

// URI: /api/v1/transactionSign
// Method: POST
// Args: JSON Body
func transactionSign(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		if r.Header.Get("Content-Type") != ContentTypeJSON {
			resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req TransactionSignRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		if err := req.validate(); err != nil {
			logger.WithError(err).Error("invalid sign transaction request")
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		txnInputs, txnOutputs, err := req.TransactionParams()
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		// for integration tests
		if autoPressEmulatorButtons {
			err := gateway.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				logger.Error("transactionSign failed: %s", err.Error())
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		var msg wire.Message
		retCH := make(chan int)
		errCH := make(chan int)
		ctx := r.Context()

		go func() {
			msg, err = gateway.TransactionSign(txnInputs, txnOutputs)
			if err != nil {
				errCH <- 1
				return
			}
			retCH <- 1
		}()

		select {
		case <-retCH:
			HandleFirmwareResponseMessages(w, msg)
		case <-errCH:
			logger.Errorf("transactionSign failed: %s", err.Error())
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
		case <-ctx.Done():
			disConnErr := gateway.Disconnect()
			if disConnErr != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
			} else {
				resp := NewHTTPErrorResponse(499, "Client Closed Request")
				writeHTTPResponse(w, resp)
			}
		}
	}
}

func (r *TransactionSignRequest) validate() error {
	if len(r.TransactionInputs) == 0 {
		return errors.New("inputs are required")
	}

	for _, input := range r.TransactionInputs {
		if input.Hash == "" {
			return errors.New("input hash cannot be empty")
		}
	}

	for _, output := range r.TransactionOutputs {
		if output.Address == "" {
			return errors.New("address cannot be empty")
		}

		if output.Coins == "" {
			return errors.New("coins cannot be empty")
		}

		if output.Hours == "" {
			return errors.New("hours cannot be empty")
		}
	}

	return nil
}

// TransactionParams returns params for a transaction from the request data
func (r *TransactionSignRequest) TransactionParams() ([]*messages.SkycoinTransactionInput, []*messages.SkycoinTransactionOutput, error) {
	var transactionInputs []*messages.SkycoinTransactionInput
	var transactionOutputs []*messages.SkycoinTransactionOutput

	for _, input := range r.TransactionInputs {
		var transactionInput messages.SkycoinTransactionInput

		transactionInput.HashIn = proto.String(input.Hash)

		if input.Index != nil {
			transactionInput.Index = proto.Uint32(*input.Index)
		} else {
			transactionInput.Index = nil
		}
		transactionInputs = append(transactionInputs, &transactionInput)
	}

	for _, output := range r.TransactionOutputs {
		var transactionOutput messages.SkycoinTransactionOutput

		_, err := cipher.DecodeBase58Address(output.Address)
		if err != nil {
			return nil, nil, err
		}

		coins, err := droplet.FromString(output.Coins)
		if err != nil {
			return nil, nil, err
		}

		hours, err := strconv.ParseUint(output.Hours, 10, 64)
		if err != nil {
			return nil, nil, err
		}

		transactionOutput.Address = proto.String(output.Address)
		transactionOutput.Coin = proto.Uint64(coins)
		transactionOutput.Hour = proto.Uint64(hours)

		if output.AddressIndex != nil {
			transactionOutput.AddressIndex = proto.Uint32(*output.AddressIndex)
		}

		transactionOutputs = append(transactionOutputs, &transactionOutput)
	}

	return transactionInputs, transactionOutputs, nil
}
