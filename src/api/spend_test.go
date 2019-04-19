package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/wallet"
)

type rawHoursSelection struct {
	Type        string  `json:"type"`
	Mode        string  `json:"mode"`
	ShareFactor *string `json:"share_factor,omitempty"`
}

type rawReceiver struct {
	Address string `json:"address"`
	Coins   string `json:"coins"`
	Hours   string `json:"hours,omitempty"`
}

type rawCreateTxnRequest struct {
	UxOuts         []string          `json:"unspents,omitempty"`
	Addresses      []string          `json:"addresses,omitempty"`
	HoursSelection rawHoursSelection `json:"hours_selection"`
	ChangeAddress  string            `json:"change_address,omitempty"`
	To             []rawReceiver     `json:"to"`
	Password       string            `json:"password"`
}

func TestCreateTransaction(t *testing.T) {
	changeAddress := testutil.MakeAddress()
	destinationAddress := testutil.MakeAddress()
	emptyAddress := cipher.Address{}

	txn := &coin.Transaction{
		Length:    100,
		Type:      0,
		InnerHash: testutil.RandSHA256(t),
		In:        []cipher.SHA256{testutil.RandSHA256(t)},
		Out: []coin.TransactionOutput{
			{
				Address: destinationAddress,
				Coins:   1e6,
				Hours:   100,
			},
		},
	}

	inputs := []visor.TransactionInput{
		{
			UxOut: coin.UxOut{
				Head: coin.UxHead{
					Time:  uint64(time.Now().UTC().Unix()),
					BkSeq: 9999,
				},
				Body: coin.UxBody{
					SrcTransaction: testutil.RandSHA256(t),
					Address:        testutil.MakeAddress(),
					Coins:          1e6,
					Hours:          100,
				},
			},
			CalculatedHours: 200,
		},
	}

	createdTxn, err := NewCreatedTransaction(txn, inputs)
	require.NoError(t, err)

	createTxnResponse := CreateTransactionResponse{
		Transaction:        *createdTxn,
		EncodedTransaction: txn.MustSerializeHex(),
	}

	validBody := &rawCreateTxnRequest{
		HoursSelection: rawHoursSelection{
			Type: transaction.HoursSelectionTypeManual,
		},
		To: []rawReceiver{
			{
				Address: destinationAddress.String(),
				Coins:   "100",
				Hours:   "10",
			},
		},
		ChangeAddress: changeAddress.String(),
		UxOuts:        []string{testutil.RandSHA256(t).Hex(), testutil.RandSHA256(t).Hex()},
	}

	walletInput := testutil.RandSHA256(t)

	tt := []struct {
		name    string
		method  string
		status  int
		body    *rawCreateTxnRequest
		rawBody string

		gatewayCreateTransactionResult *coin.Transaction
		gatewayCreateTransactionInputs []visor.TransactionInput
		gatewayCreateTransactionErr    error

		csrfDisabled bool
		contentType  string

		httpResponse HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodGet,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:         "415",
			method:       http.MethodPost,
			status:       http.StatusUnsupportedMediaType,
			contentType:  ContentTypeForm,
			httpResponse: NewHTTPErrorResponse(http.StatusUnsupportedMediaType, ""),
		},

		{
			name:         "400 - missing hours selection type",
			method:       http.MethodPost,
			body:         &rawCreateTxnRequest{},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "missing hours_selection.type"),
		},

		{
			name:   "400 - invalid hours selection type",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: "foo",
				},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid hours_selection.type"),
		},

		{
			name:   "400 - missing hours selection mode",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeAuto,
				},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "missing hours_selection.mode"),
		},

		{
			name:   "400 - invalid hours selection mode",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeAuto,
					Mode: "foo",
				},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid hours_selection.mode"),
		},

		{
			name:   "400 - missing hours selection share factor",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeAuto,
					Mode: transaction.HoursSelectionModeShare,
				},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "missing hours_selection.share_factor when hours_selection.mode is share"),
		},

		{
			name:   "400 - share factor set but mode is not share",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeManual,
					ShareFactor: newStrPtr("0.5"),
				},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "hours_selection.share_factor can only be used when hours_selection.mode is share"),
		},

		{
			name:   "400 - negative share factor",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: newStrPtr("-1"),
				},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "hours_selection.share_factor cannot be negative"),
		},

		{
			name:   "400 - share factor greater than 1",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: newStrPtr("1.1"),
				},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "hours_selection.share_factor cannot be more than 1"),
		},

		{
			name:   "400 - empty sender address",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: changeAddress.String(),
				Addresses:     []string{""},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid address: Invalid base58 string"),
		},

		{
			name:   "400 - invalid sender address",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: changeAddress.String(),
				Addresses:     []string{"xxx"},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid address: Invalid address length"),
		},

		{
			name:   "400 - invalid change address",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: "xxx",
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid address: Invalid address length"),
		},

		{
			name:   "400 - empty change address",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: emptyAddress.String(),
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "change_address must not be the null address"),
		},

		{
			name:   "400 - auto type destination has hours",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Hours:   "100",
						Coins:   "1.01",
					},
				},
				ChangeAddress: changeAddress.String(),
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "to[0].hours must not be specified for auto hours_selection.mode"),
		},

		{
			name:   "400 - manual type destination missing hours",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.01",
					},
				},
				ChangeAddress: changeAddress.String(),
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "to[0].hours must be specified for manual hours_selection.mode"),
		},

		{
			name:   "400 - manual type has mode set",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
					Mode: transaction.HoursSelectionModeShare,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.01",
						Hours:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "hours_selection.mode cannot be used for manual hours_selection.type"),
		},

		{
			name:   "400 - address is empty",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.01",
						Hours:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
				Addresses:     []string{emptyAddress.String()},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "addresses[0] is empty"),
		},

		{
			name:   "400 - to address is empty",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: emptyAddress.String(),
						Coins:   "1.01",
						Hours:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "to[0].address is empty"),
		},

		{
			name:   "400 - to coins is zero",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "0",
						Hours:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "to[0].coins must not be zero"),
		},

		{
			name:   "400 - invalid to coins",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "0.1a",
						Hours:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "can't convert 0.1a to decimal"),
		},

		{
			name:   "400 - invalid to hours",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "0.1",
						Hours:   "100.1",
					},
				},
				ChangeAddress: changeAddress.String(),
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid hours value: strconv.ParseUint: parsing \"100.1\": invalid syntax"),
		},

		{
			name:   "400 - empty string to coins",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "",
						Hours:   "",
					},
				},
				ChangeAddress: changeAddress.String(),
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "can't convert  to decimal"),
		},

		{
			name:   "400 - coins has too many decimals",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.1234",
					},
				},
				ChangeAddress: changeAddress.String(),
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "to[0].coins has too many decimal places"),
		},

		{
			name:   "400 - empty to",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: changeAddress.String(),
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "to is empty"),
		},

		{
			name:   "400 - manual duplicate outputs",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: changeAddress.String(),
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.2",
						Hours:   "100",
					},
					{
						Address: destinationAddress.String(),
						Coins:   "1.2",
						Hours:   "100",
					},
				},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "to contains duplicate values"),
		},

		{
			name:   "400 - auto duplicate outputs",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				ChangeAddress: changeAddress.String(),
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.2",
					},
					{
						Address: destinationAddress.String(),
						Coins:   "1.2",
					},
				},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "to contains duplicate values"),
		},

		{
			name:   "400 - both uxouts and addresses specified",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				ChangeAddress: changeAddress.String(),
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.2",
					},
				},
				Addresses: []string{destinationAddress.String()},
				UxOuts:    []string{walletInput.Hex()},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "unspents and addresses cannot be combined"),
		},

		{
			name:   "400 - missing uxouts and addresses",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				ChangeAddress: changeAddress.String(),
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.2",
					},
				},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "one of addresses or unspents must not be empty"),
		},

		{
			name:   "400 - duplicate uxouts",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				ChangeAddress: changeAddress.String(),
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.2",
					},
				},
				UxOuts: []string{walletInput.Hex(), walletInput.Hex()},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "unspents contains duplicate values"),
		},

		{
			name:   "400 - duplicate addresses",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				ChangeAddress: changeAddress.String(),
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.2",
					},
				},
				Addresses: []string{destinationAddress.String(), destinationAddress.String()},
			},
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "addresses contains duplicate values"),
		},

		{
			name:   "200 - auto type split even",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
				Addresses:     []string{changeAddress.String()},
			},
			status:                         http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			httpResponse: HTTPResponse{
				Data: createTxnResponse,
			},
		},

		{
			name:   "200 - manual type zero hours",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "100",
						Hours:   "0",
					},
				},
				ChangeAddress: changeAddress.String(),
				Addresses:     []string{changeAddress.String()},
			},
			status:                         http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			httpResponse: HTTPResponse{
				Data: createTxnResponse,
			},
		},

		{
			name:   "200 - manual type nonzero hours",
			method: http.MethodPost,
			body: &rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "100",
						Hours:   "10",
					},
				},
				ChangeAddress: changeAddress.String(),
				Addresses:     []string{changeAddress.String()},
			},
			status:                         http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			httpResponse: HTTPResponse{
				Data: createTxnResponse,
			},
		},

		{
			name:                           "200 - manual type nonzero hours - csrf disabled",
			method:                         http.MethodPost,
			body:                           validBody,
			status:                         http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			httpResponse: HTTPResponse{
				Data: createTxnResponse,
			},
			csrfDisabled: true,
		},

		{
			name:                        "500 - misc error",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusInternalServerError,
			gatewayCreateTransactionErr: errors.New("unhandled error"),
			httpResponse:                NewHTTPErrorResponse(http.StatusInternalServerError, "unhandled error"),
		},

		{
			name:                        "400 - no fee",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusBadRequest,
			gatewayCreateTransactionErr: fee.ErrTxnNoFee,
			httpResponse:                NewHTTPErrorResponse(http.StatusBadRequest, "Transaction has zero coinhour fee"),
		},

		{
			name:                        "400 - insufficient coin hours",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusBadRequest,
			gatewayCreateTransactionErr: fee.ErrTxnInsufficientCoinHours,
			httpResponse:                NewHTTPErrorResponse(http.StatusBadRequest, "Insufficient coinhours for transaction outputs"),
		},

		{
			name:                        "400 - uxout doesn't exist",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusBadRequest,
			gatewayCreateTransactionErr: blockdb.NewErrUnspentNotExist("foo"),
			httpResponse:                NewHTTPErrorResponse(http.StatusBadRequest, "unspent output of foo does not exist"),
		},

		{
			name:                        "400 - visor error",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusBadRequest,
			gatewayCreateTransactionErr: visor.ErrSpendingUnconfirmed,
			httpResponse:                NewHTTPErrorResponse(http.StatusBadRequest, "Please spend after your pending transaction is confirmed"),
		},

		{
			name:                        "400 - txn create error",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusBadRequest,
			gatewayCreateTransactionErr: transaction.ErrInsufficientBalance,
			httpResponse:                NewHTTPErrorResponse(http.StatusBadRequest, "balance is not sufficient"),
		},

		{
			name:         "400 - invalid json",
			method:       http.MethodPost,
			rawBody:      "{ca",
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid character 'c' looking for beginning of object key string"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}

			// If the rawRequestBody can be deserialized to CreateTransactionRequest, use it to mock gateway.WalletCreateTransaction
			serializedBody, err := json.Marshal(tc.body)
			require.NoError(t, err)
			var body walletCreateTransactionRequest
			err = json.Unmarshal(serializedBody, &body)
			if err == nil {
				x := gateway.On("CreateTransaction", body.TransactionParams(), body.VisorParams())
				x.Return(tc.gatewayCreateTransactionResult, tc.gatewayCreateTransactionInputs, tc.gatewayCreateTransactionErr)
			}

			endpoint := "/api/v2/transaction"

			bodyText := []byte(tc.rawBody)
			if len(bodyText) == 0 {
				bodyText, err = json.Marshal(tc.body)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBuffer(bodyText))
			require.NoError(t, err)

			contentType := tc.contentType
			if contentType == "" {
				contentType = ContentTypeJSON
			}

			req.Header.Add("Content-Type", contentType)

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = tc.csrfDisabled

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v` (%v)", status, tc.status, rr.Body)

			var rsp ReceivedHTTPResponse
			err = json.Unmarshal(rr.Body.Bytes(), &rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var msg CreateTransactionResponse
				err := json.Unmarshal(rsp.Data, &msg)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data.(CreateTransactionResponse), msg)
			}
		})
	}
}

func TestWalletCreateTransaction(t *testing.T) {
	type rawWalletCreateTxnRequest struct {
		rawCreateTxnRequest
		WalletID string `json:"wallet_id"`
		Password string `json:"password"`
		Unsigned bool   `json:"unsigned"`
	}

	changeAddress := testutil.MakeAddress()
	destinationAddress := testutil.MakeAddress()
	emptyAddress := cipher.Address{}

	txn := &coin.Transaction{
		Length:    100,
		Type:      0,
		InnerHash: testutil.RandSHA256(t),
		In:        []cipher.SHA256{testutil.RandSHA256(t)},
		Out: []coin.TransactionOutput{
			{
				Address: destinationAddress,
				Coins:   1e6,
				Hours:   100,
			},
		},
	}

	inputs := []visor.TransactionInput{
		{
			UxOut: coin.UxOut{
				Head: coin.UxHead{
					Time:  uint64(time.Now().UTC().Unix()),
					BkSeq: 9999,
				},
				Body: coin.UxBody{
					SrcTransaction: testutil.RandSHA256(t),
					Address:        testutil.MakeAddress(),
					Coins:          1e6,
					Hours:          100,
				},
			},
			CalculatedHours: 200,
		},
	}

	createdTxn, err := NewCreatedTransaction(txn, inputs)
	require.NoError(t, err)

	createTxnResponse := &CreateTransactionResponse{
		Transaction:        *createdTxn,
		EncodedTransaction: txn.MustSerializeHex(),
	}

	validBody := rawWalletCreateTxnRequest{
		rawCreateTxnRequest: rawCreateTxnRequest{
			HoursSelection: rawHoursSelection{
				Type: transaction.HoursSelectionTypeManual,
			},
			To: []rawReceiver{
				{
					Address: destinationAddress.String(),
					Coins:   "100",
					Hours:   "10",
				},
			},
			ChangeAddress: changeAddress.String(),
		},
		WalletID: "foo.wlt",
	}

	walletInput := testutil.RandSHA256(t)

	type testCase struct {
		name                           string
		method                         string
		body                           rawWalletCreateTxnRequest
		rawBody                        string
		status                         int
		err                            string
		gatewayCreateTransactionResult *coin.Transaction
		gatewayCreateTransactionInputs []visor.TransactionInput
		gatewayCreateTransactionErr    error
		createTransactionResponse      *CreateTransactionResponse
		csrfDisabled                   bool
		contentType                    string
	}

	baseCases := []testCase{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			body:   rawWalletCreateTxnRequest{},
			err:    "405 Method Not Allowed",
		},

		{
			name:        "415",
			method:      http.MethodPost,
			status:      http.StatusUnsupportedMediaType,
			body:        rawWalletCreateTxnRequest{},
			contentType: ContentTypeForm,
			err:         "415 Unsupported Media Type",
		},

		{
			name:   "400 - missing hours selection type",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing hours_selection.type",
		},

		{
			name:   "400 - invalid hours selection type",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: "foo",
					},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid hours_selection.type",
		},

		{
			name:   "400 - missing hours selection mode",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeAuto,
					},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing hours_selection.mode",
		},

		{
			name:   "400 - invalid hours selection mode",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeAuto,
						Mode: "foo",
					},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid hours_selection.mode",
		},

		{
			name:   "400 - missing hours selection share factor",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeAuto,
						Mode: transaction.HoursSelectionModeShare,
					},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing hours_selection.share_factor when hours_selection.mode is share",
		},

		{
			name:   "400 - share factor set but mode is not share",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type:        transaction.HoursSelectionTypeManual,
						ShareFactor: newStrPtr("0.5"),
					},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - hours_selection.share_factor can only be used when hours_selection.mode is share",
		},

		{
			name:   "400 - negative share factor",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type:        transaction.HoursSelectionTypeAuto,
						Mode:        transaction.HoursSelectionModeShare,
						ShareFactor: newStrPtr("-1"),
					},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - hours_selection.share_factor cannot be negative",
		},

		{
			name:   "400 - share factor greater than 1",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type:        transaction.HoursSelectionTypeAuto,
						Mode:        transaction.HoursSelectionModeShare,
						ShareFactor: newStrPtr("1.1"),
					},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - hours_selection.share_factor cannot be more than 1",
		},

		{
			name:   "400 - empty sender address",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					Addresses:     []string{""},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid address: Invalid base58 string",
		},

		{
			name:   "400 - invalid sender address",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: changeAddress.String(),
					Addresses:     []string{"xxx"},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid address: Invalid address length",
		},

		{
			name:   "400 - invalid change address",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: "xxx",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid address: Invalid address length",
		},

		{
			name:   "400 - empty change address",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: emptyAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - change_address must not be the null address",
		},

		{
			name:   "400 - auto type destination has hours",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type:        transaction.HoursSelectionTypeAuto,
						Mode:        transaction.HoursSelectionModeShare,
						ShareFactor: newStrPtr("0.5"),
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Hours:   "100",
							Coins:   "1.01",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to[0].hours must not be specified for auto hours_selection.mode",
		},

		{
			name:   "400 - manual type destination missing hours",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "1.01",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to[0].hours must be specified for manual hours_selection.mode",
		},

		{
			name:   "400 - manual type has mode set",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
						Mode: transaction.HoursSelectionModeShare,
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "1.01",
							Hours:   "100",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - hours_selection.mode cannot be used for manual hours_selection.type",
		},

		{
			name:   "400 - missing wallet ID",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "1.01",
							Hours:   "100",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing wallet_id",
		},

		{
			name:   "400 - address is empty",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "1.01",
							Hours:   "100",
						},
					},
					ChangeAddress: changeAddress.String(),
					Addresses:     []string{emptyAddress.String()},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - addresses[0] is empty",
		},

		{
			name:   "400 - to address is empty",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					To: []rawReceiver{
						{
							Address: emptyAddress.String(),
							Coins:   "1.01",
							Hours:   "100",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to[0].address is empty",
		},

		{
			name:   "400 - to coins is zero",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "0",
							Hours:   "100",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to[0].coins must not be zero",
		},

		{
			name:   "400 - invalid to coins",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "0.1a",
							Hours:   "100",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - can't convert 0.1a to decimal",
		},

		{
			name:   "400 - invalid to hours",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "0.1",
							Hours:   "100.1",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid hours value: strconv.ParseUint: parsing \"100.1\": invalid syntax",
		},

		{
			name:   "400 - empty string to coins",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type:        transaction.HoursSelectionTypeAuto,
						Mode:        transaction.HoursSelectionModeShare,
						ShareFactor: newStrPtr("0.5"),
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "",
							Hours:   "",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - can't convert  to decimal",
		},

		{
			name:   "400 - coins has too many decimals",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type:        transaction.HoursSelectionTypeAuto,
						Mode:        transaction.HoursSelectionModeShare,
						ShareFactor: newStrPtr("0.5"),
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "1.1234",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to[0].coins has too many decimal places",
		},

		{
			name:   "400 - empty to",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to is empty",
		},

		{
			name:   "400 - manual duplicate outputs",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: changeAddress.String(),
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "1.2",
							Hours:   "100",
						},
						{
							Address: destinationAddress.String(),
							Coins:   "1.2",
							Hours:   "100",
						},
					},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to contains duplicate values",
		},

		{
			name:   "400 - auto duplicate outputs",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type:        transaction.HoursSelectionTypeAuto,
						Mode:        transaction.HoursSelectionModeShare,
						ShareFactor: newStrPtr("0.5"),
					},
					ChangeAddress: changeAddress.String(),
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "1.2",
						},
						{
							Address: destinationAddress.String(),
							Coins:   "1.2",
						},
					},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to contains duplicate values",
		},

		{
			name:   "400 - both uxouts and addresses specified",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type:        transaction.HoursSelectionTypeAuto,
						Mode:        transaction.HoursSelectionModeShare,
						ShareFactor: newStrPtr("0.5"),
					},
					ChangeAddress: changeAddress.String(),
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "1.2",
						},
					},
					Addresses: []string{destinationAddress.String()},
					UxOuts:    []string{walletInput.Hex()},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - unspents and addresses cannot be combined",
		},

		{
			name:   "400 - duplicate uxouts",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type:        transaction.HoursSelectionTypeAuto,
						Mode:        transaction.HoursSelectionModeShare,
						ShareFactor: newStrPtr("0.5"),
					},
					ChangeAddress: changeAddress.String(),
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "1.2",
						},
					},
					UxOuts: []string{walletInput.Hex(), walletInput.Hex()},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - unspents contains duplicate values",
		},

		{
			name:   "400 - duplicate addresses",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type:        transaction.HoursSelectionTypeAuto,
						Mode:        transaction.HoursSelectionModeShare,
						ShareFactor: newStrPtr("0.5"),
					},
					ChangeAddress: changeAddress.String(),
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "1.2",
						},
					},
					Addresses: []string{destinationAddress.String(), destinationAddress.String()},
				},
				WalletID: "foo.wlt",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - addresses contains duplicate values",
		},

		{
			name:   "200 - auto type split even",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type:        transaction.HoursSelectionTypeAuto,
						Mode:        transaction.HoursSelectionModeShare,
						ShareFactor: newStrPtr("0.5"),
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "100",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status:                         http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			createTransactionResponse:      createTxnResponse,
		},

		{
			name:   "200 - manual type zero hours",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "100",
							Hours:   "0",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status:                         http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			createTransactionResponse:      createTxnResponse,
		},

		{
			name:   "200 - manual type nonzero hours",
			method: http.MethodPost,
			body: rawWalletCreateTxnRequest{
				rawCreateTxnRequest: rawCreateTxnRequest{
					HoursSelection: rawHoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					To: []rawReceiver{
						{
							Address: destinationAddress.String(),
							Coins:   "100",
							Hours:   "10",
						},
					},
					ChangeAddress: changeAddress.String(),
				},
				WalletID: "foo.wlt",
			},
			status:                         http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			createTransactionResponse:      createTxnResponse,
		},

		{
			name:                           "200 - manual type nonzero hours - csrf disabled",
			method:                         http.MethodPost,
			body:                           validBody,
			status:                         http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			createTransactionResponse:      createTxnResponse,
			csrfDisabled:                   true,
		},

		{
			name:                        "500 - misc error",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusInternalServerError,
			gatewayCreateTransactionErr: errors.New("unhandled error"),
			err:                         "500 Internal Server Error - unhandled error",
		},

		{
			name:                        "400 - no fee",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusBadRequest,
			gatewayCreateTransactionErr: fee.ErrTxnNoFee,
			err:                         "400 Bad Request - Transaction has zero coinhour fee",
		},

		{
			name:                        "400 - insufficient coin hours",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusBadRequest,
			gatewayCreateTransactionErr: fee.ErrTxnInsufficientCoinHours,
			err:                         "400 Bad Request - Insufficient coinhours for transaction outputs",
		},

		{
			name:                        "400 - uxout doesn't exist",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusBadRequest,
			gatewayCreateTransactionErr: blockdb.NewErrUnspentNotExist("foo"),
			err:                         "400 Bad Request - unspent output of foo does not exist",
		},

		{
			name:    "400 - invalid json",
			method:  http.MethodPost,
			body:    rawWalletCreateTxnRequest{},
			rawBody: "{ca",
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - invalid character 'c' looking for beginning of object key string",
		},

		{
			name:                        "400 - other wallet error",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusBadRequest,
			gatewayCreateTransactionErr: wallet.ErrWalletEncrypted,
			err:                         "400 Bad Request - wallet is encrypted",
		},

		{
			name:                        "404 - wallet not found",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusNotFound,
			gatewayCreateTransactionErr: wallet.ErrWalletNotExist,
			err:                         "404 Not Found - wallet doesn't exist",
		},

		{
			name:                        "403 - wallet API disabled",
			method:                      http.MethodPost,
			body:                        validBody,
			status:                      http.StatusForbidden,
			gatewayCreateTransactionErr: wallet.ErrWalletAPIDisabled,
			err:                         "403 Forbidden",
		},
	}

	cases := make([]testCase, len(baseCases)*2)
	copy(cases, baseCases)
	copy(cases[len(baseCases):], baseCases)
	for i := range baseCases {
		cases[i].body.Unsigned = true
	}

	cases = append(cases, testCase{
		name:   "400 - password provided for unsigned request",
		method: http.MethodPost,
		body: rawWalletCreateTxnRequest{
			rawCreateTxnRequest: rawCreateTxnRequest{
				HoursSelection: rawHoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.01",
						Hours:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
			},
			WalletID: "foo.wlt",
			Password: "foo",
			Unsigned: true,
		},
		status: http.StatusBadRequest,
		err:    "400 Bad Request - password must not be used for unsigned transactions",
	})

	for _, tc := range cases {
		name := fmt.Sprintf("unsigned=%v %s", tc.body.Unsigned, tc.name)
		t.Run(name, func(t *testing.T) {
			gateway := &MockGatewayer{}

			// If the rawRequestBody can be deserialized to CreateTransactionRequest, use it to mock gateway.WalletCreateTransaction
			serializedBody, err := json.Marshal(tc.body)
			require.NoError(t, err)
			var body walletCreateTransactionRequest
			err = json.Unmarshal(serializedBody, &body)
			if err == nil {
				if tc.body.Unsigned {
					x := gateway.On("WalletCreateTransaction", body.WalletID, body.TransactionParams(), body.VisorParams())
					x.Return(tc.gatewayCreateTransactionResult, tc.gatewayCreateTransactionInputs, tc.gatewayCreateTransactionErr)
				} else {
					x := gateway.On("WalletCreateTransactionSigned", body.WalletID, []byte(body.Password), body.TransactionParams(), body.VisorParams())
					x.Return(tc.gatewayCreateTransactionResult, tc.gatewayCreateTransactionInputs, tc.gatewayCreateTransactionErr)

				}
			}

			endpoint := "/api/v1/wallet/transaction"

			bodyText := []byte(tc.rawBody)
			if len(bodyText) == 0 {
				bodyText, err = json.Marshal(tc.body)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBuffer(bodyText))
			require.NoError(t, err)

			contentType := tc.contentType
			if contentType == "" {
				contentType = ContentTypeJSON
			}

			req.Header.Add("Content-Type", contentType)

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = tc.csrfDisabled

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()))
			} else {
				var msg CreateTransactionResponse
				err := json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.NotNil(t, tc.createTransactionResponse)
				require.Equal(t, *tc.createTransactionResponse, msg)
			}
		})
	}
}

func newStrPtr(s string) *string {
	return &s
}

func TestWalletSignTransaction(t *testing.T) {
	destinationAddress := testutil.MakeAddress()

	signedTxn := coin.Transaction{
		Length:    100,
		Type:      0,
		InnerHash: testutil.RandSHA256(t),
		Sigs:      []cipher.Sig{testutil.RandSig(t), testutil.RandSig(t)},
		In:        []cipher.SHA256{testutil.RandSHA256(t), testutil.RandSHA256(t)},
		Out: []coin.TransactionOutput{
			{
				Address: destinationAddress,
				Coins:   1e6,
				Hours:   100,
			},
		},
	}

	txn := signedTxn
	txn.Sigs = make([]cipher.Sig, len(txn.In))

	inputs := []visor.TransactionInput{
		{
			UxOut: coin.UxOut{
				Head: coin.UxHead{
					Time:  uint64(time.Now().UTC().Unix()),
					BkSeq: 9999,
				},
				Body: coin.UxBody{
					SrcTransaction: testutil.RandSHA256(t),
					Address:        testutil.MakeAddress(),
					Coins:          1e6,
					Hours:          100,
				},
			},
			CalculatedHours: 200,
		},
		{
			UxOut: coin.UxOut{
				Head: coin.UxHead{
					Time:  uint64(time.Now().UTC().Unix()),
					BkSeq: 9999,
				},
				Body: coin.UxBody{
					SrcTransaction: testutil.RandSHA256(t),
					Address:        testutil.MakeAddress(),
					Coins:          1e6,
					Hours:          100,
				},
			},
			CalculatedHours: 200,
		},
	}

	signedTxnResp, err := NewCreateTransactionResponse(&signedTxn, inputs)
	require.NoError(t, err)

	validBody := &WalletSignTransactionRequest{
		WalletID:           "foo.wlt",
		EncodedTransaction: txn.MustSerializeHex(),
	}

	tt := []struct {
		name                         string
		method                       string
		body                         *WalletSignTransactionRequest
		rawBody                      string
		status                       int
		gatewaySignTransactionResult *coin.Transaction
		gatewaySignTransactionInputs []visor.TransactionInput
		gatewaySignTransactionErr    error
		csrfDisabled                 bool
		contentType                  string
		httpResponse                 HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodGet,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:         "415",
			method:       http.MethodPost,
			status:       http.StatusUnsupportedMediaType,
			contentType:  ContentTypeForm,
			httpResponse: NewHTTPErrorResponse(http.StatusUnsupportedMediaType, ""),
		},

		{
			name:   "400 wallet ID required",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			body: &WalletSignTransactionRequest{
				EncodedTransaction: validBody.EncodedTransaction,
			},
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "wallet_id is required"),
		},

		{
			name:   "400 encoded_transaction is required",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			body: &WalletSignTransactionRequest{
				WalletID: "foo.wlt",
			},
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "encoded_transaction is required"),
		},

		{
			name:   "400 decode transaction failed",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			body: &WalletSignTransactionRequest{
				WalletID:           "foo.wlt",
				EncodedTransaction: "abc",
			},
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "Decode transaction failed: encoding/hex: odd length hex string"),
		},

		{
			name:   "400 too many sign indexes",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			body: &WalletSignTransactionRequest{
				WalletID:           "foo.wlt",
				EncodedTransaction: validBody.EncodedTransaction,
				SignIndexes:        []int{0, 1, 2},
			},
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "Too many values in sign_indexes"),
		},

		{
			name:   "400 sign indexes out of range",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			body: &WalletSignTransactionRequest{
				WalletID:           "foo.wlt",
				EncodedTransaction: validBody.EncodedTransaction,
				SignIndexes:        []int{5},
			},
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "Value in sign_indexes exceeds range of transaction inputs array"),
		},

		{
			name:   "400 duplicate sign indexes",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			body: &WalletSignTransactionRequest{
				WalletID:           "foo.wlt",
				EncodedTransaction: validBody.EncodedTransaction,
				SignIndexes:        []int{1, 1},
			},
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "Duplicate value in sign_indexes"),
		},

		{
			name:                      "500 - misc error",
			method:                    http.MethodPost,
			body:                      validBody,
			status:                    http.StatusInternalServerError,
			gatewaySignTransactionErr: errors.New("unhandled error"),
			httpResponse:              NewHTTPErrorResponse(http.StatusInternalServerError, "unhandled error"),
		},

		{
			name:                      "400 - wallet not encrypted",
			method:                    http.MethodPost,
			body:                      validBody,
			status:                    http.StatusBadRequest,
			gatewaySignTransactionErr: wallet.ErrWalletNotEncrypted,
			httpResponse:              NewHTTPErrorResponse(http.StatusBadRequest, "wallet is not encrypted"),
		},

		{
			name:                      "400 - wallet encrypted",
			method:                    http.MethodPost,
			body:                      validBody,
			status:                    http.StatusBadRequest,
			gatewaySignTransactionErr: wallet.ErrWalletEncrypted,
			httpResponse:              NewHTTPErrorResponse(http.StatusBadRequest, "wallet is encrypted"),
		},

		{
			name:                      "400 - violates hard constraint",
			method:                    http.MethodPost,
			body:                      validBody,
			status:                    http.StatusBadRequest,
			gatewaySignTransactionErr: visor.NewErrTxnViolatesHardConstraint(errors.New("bad txn")),
			httpResponse:              NewHTTPErrorResponse(http.StatusBadRequest, "Transaction violates hard constraint: bad txn"),
		},

		{
			name:                      "400 - unspents do not exist",
			method:                    http.MethodPost,
			body:                      validBody,
			status:                    http.StatusBadRequest,
			gatewaySignTransactionErr: blockdb.NewErrUnspentNotExist("foo"),
			httpResponse:              NewHTTPErrorResponse(http.StatusBadRequest, "unspent output of foo does not exist"),
		},

		{
			name:         "400 - invalid json",
			method:       http.MethodPost,
			rawBody:      "{ca",
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "invalid character 'c' looking for beginning of object key string"),
		},

		{
			name:                      "404 - wallet not found",
			method:                    http.MethodPost,
			body:                      validBody,
			status:                    http.StatusNotFound,
			gatewaySignTransactionErr: wallet.ErrWalletNotExist,
			httpResponse:              NewHTTPErrorResponse(http.StatusNotFound, "wallet doesn't exist"),
		},

		{
			name:                      "403 - wallet API disabled",
			method:                    http.MethodPost,
			body:                      validBody,
			status:                    http.StatusForbidden,
			gatewaySignTransactionErr: wallet.ErrWalletAPIDisabled,
			httpResponse:              NewHTTPErrorResponse(http.StatusForbidden, "wallet api is disabled"),
		},

		{
			name:                         "200 - no password",
			method:                       http.MethodPost,
			body:                         validBody,
			status:                       http.StatusOK,
			gatewaySignTransactionResult: &signedTxn,
			gatewaySignTransactionInputs: inputs,
			httpResponse: HTTPResponse{
				Data: *signedTxnResp,
			},
		},

		{
			name:                         "200 - no password csrf disabled",
			method:                       http.MethodPost,
			body:                         validBody,
			status:                       http.StatusOK,
			gatewaySignTransactionResult: &signedTxn,
			gatewaySignTransactionInputs: inputs,
			httpResponse: HTTPResponse{
				Data: *signedTxnResp,
			},
			csrfDisabled: true,
		},

		{
			name:   "200 - password",
			method: http.MethodPost,
			body: &WalletSignTransactionRequest{
				WalletID:           "foo.wlt",
				Password:           "foo",
				EncodedTransaction: validBody.EncodedTransaction,
			},
			status:                       http.StatusOK,
			gatewaySignTransactionResult: &signedTxn,
			gatewaySignTransactionInputs: inputs,
			httpResponse: HTTPResponse{
				Data: *signedTxnResp,
			},
		},

		{
			name:   "200 - sign indexes",
			method: http.MethodPost,
			body: &WalletSignTransactionRequest{
				WalletID:           "foo.wlt",
				SignIndexes:        []int{1},
				EncodedTransaction: validBody.EncodedTransaction,
			},
			status:                       http.StatusOK,
			gatewaySignTransactionResult: &signedTxn,
			gatewaySignTransactionInputs: inputs,
			httpResponse: HTTPResponse{
				Data: *signedTxnResp,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}

			var txn *coin.Transaction
			if tc.body != nil {
				// Decode the transaction used in the request body, but ignore an error in case the
				// transaction is intentionally malformed
				txnx, err := coin.DeserializeTransactionHex(tc.body.EncodedTransaction)
				if err == nil {
					txn = &txnx
				}
			}

			if tc.body != nil {
				gateway.On("WalletSignTransaction", tc.body.WalletID, []byte(tc.body.Password), txn, tc.body.SignIndexes).Return(tc.gatewaySignTransactionResult, tc.gatewaySignTransactionInputs, tc.gatewaySignTransactionErr)
			}

			endpoint := "/api/v2/wallet/transaction/sign"

			bodyText := []byte(tc.rawBody)
			if len(bodyText) == 0 {
				var err error
				bodyText, err = json.Marshal(tc.body)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBuffer(bodyText))
			require.NoError(t, err)

			contentType := tc.contentType
			if contentType == "" {
				contentType = ContentTypeJSON
			}

			req.Header.Add("Content-Type", contentType)

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = tc.csrfDisabled

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			var rsp ReceivedHTTPResponse
			err = json.Unmarshal(rr.Body.Bytes(), &rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var cRsp CreateTransactionResponse
				err := json.Unmarshal(rsp.Data, &cRsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data.(CreateTransactionResponse), cRsp)
			}
		})
	}
}
