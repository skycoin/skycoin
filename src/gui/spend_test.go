package gui

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil" //http,json helpers
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

func TestCreateTransaction(t *testing.T) {
	type rawWalletRequest struct {
		ID        string   `json:"id"`
		Addresses []string `json:"addresses,omitempty"`
		Password  string   `json:"password"`
	}

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

	type rawRequest struct {
		HoursSelection rawHoursSelection `json:"hours_selection"`
		Wallet         rawWalletRequest  `json:"wallet"`
		ChangeAddress  string            `json:"change_address,omitempty"`
		To             []rawReceiver     `json:"to"`
		Password       string            `json:"password"`
	}

	changeAddress := testutil.MakeAddress()
	// walletAddress := testutil.MakeAddress()
	destinationAddress := testutil.MakeAddress()
	emptyAddress := cipher.Address{}

	txn := &coin.Transaction{}
	visorTxn := &visor.Transaction{
		Txn: *txn,
	}
	readableTxn, err := visor.NewReadableTransaction(visorTxn)
	require.NoError(t, err)
	createTxnResult := &CreateTransactionResult{
		Transaction: *readableTxn,
	}

	validBody := &rawRequest{
		HoursSelection: rawHoursSelection{
			Type: wallet.HoursSelectionTypeManual,
		},
		To: []rawReceiver{
			{
				Address: destinationAddress.String(),
				Coins:   "100",
				Hours:   "10",
			},
		},
		ChangeAddress: changeAddress.String(),
		Wallet: rawWalletRequest{
			ID: "foo.wlt",
		},
	}

	tt := []struct {
		name                           string
		method                         string
		body                           *rawRequest
		status                         int
		err                            string
		gatewayCreateTransactionResult *coin.Transaction
		gatewayCreateTransactionErr    error
		createTransactionResult        *CreateTransactionResult
		csrfDisabled                   bool
		contentType                    string
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},

		{
			name:        "415",
			method:      http.MethodPost,
			status:      http.StatusUnsupportedMediaType,
			contentType: "application/x-www-form-urlencoded",
			err:         "415 Unsupported Media Type",
		},

		{
			name:   "400 - missing hours selection type",
			method: http.MethodPost,
			body:   &rawRequest{},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing hours_selection.type",
		},

		{
			name:   "400 - invalid hours selection type",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: "foo",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid hours_selection.type",
		},

		{
			name:   "400 - missing hours selection mode",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeAuto,
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing hours_selection.mode",
		},

		{
			name:   "400 - invalid hours selection mode",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeAuto,
					Mode: "foo",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid hours_selection.mode",
		},

		{
			name:   "400 - missing hours selection share factor",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeAuto,
					Mode: wallet.HoursSelectionModeSplitEven,
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing hours_selection.share_factor when hours_selection.mode is split_even",
		},

		{
			name:   "400 - share factor set but mode is not split_even",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeMatchCoins,
					ShareFactor: newStrPtr("0.5"),
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - hours_selection.share_factor can only be used when hours_selection.mode is split_even",
		},

		{
			name:   "400 - negative share factor",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeSplitEven,
					ShareFactor: newStrPtr("-1"),
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - hours_selection.share_factor cannot be negative",
		},

		{
			name:   "400 - share factor greater than 1",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeSplitEven,
					ShareFactor: newStrPtr("1.1"),
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - hours_selection.share_factor cannot be more than 1",
		},

		{
			name:   "400 - empty sender address",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: rawWalletRequest{
					Addresses: []string{""},
				},
				ChangeAddress: changeAddress.String(),
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid address: Invalid base58 string",
		},

		{
			name:   "400 - invalid sender address",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				Wallet: rawWalletRequest{
					Addresses: []string{"xxx"},
				},
				ChangeAddress: changeAddress.String(),
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid address: Invalid address length",
		},

		{
			name:   "400 - missing change address",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing change_address",
		},

		{
			name:   "400 - invalid change address",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				ChangeAddress: "xxx",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid address: Invalid address length",
		},

		{
			name:   "400 - empty change address",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				ChangeAddress: emptyAddress.String(),
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - change_address is an empty address",
		},

		{
			name:   "400 - auto type destination has hours",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeAuto,
					Mode: wallet.HoursSelectionModeMatchCoins,
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
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to[0].hours must not be specified for auto hours_selection.mode",
		},

		{
			name:   "400 - auto type destination has hours",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeAuto,
					Mode: wallet.HoursSelectionModeMatchCoins,
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
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to[0].hours must not be specified for auto hours_selection.mode",
		},

		{
			name:   "400 - manual type destination missing hours",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.01",
					},
				},
				ChangeAddress: changeAddress.String(),
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to[0].hours must be specified for manual hours_selection.mode",
		},

		{
			name:   "400 - manual type has mode set",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
					Mode: wallet.HoursSelectionModeMatchCoins,
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
			status: http.StatusBadRequest,
			err:    "400 Bad Request - hours_selection.mode cannot be used for manual hours_selection.type",
		},

		{
			name:   "400 - missing wallet ID",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.01",
						Hours:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet:        rawWalletRequest{},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing wallet.id",
		},

		{
			name:   "400 - wallet address is empty",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.01",
						Hours:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawWalletRequest{
					ID:        "foo.wlt",
					Addresses: []string{emptyAddress.String()},
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - wallet.addresses[0] is empty",
		},

		{
			name:   "400 - to address is empty",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: emptyAddress.String(),
						Coins:   "1.01",
						Hours:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawWalletRequest{
					ID: "foo.wlt",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to[0].address is empty",
		},

		{
			name:   "400 - to coins is zero",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "0",
						Hours:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawWalletRequest{
					ID: "foo.wlt",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to[0].coins must not be zero",
		},

		{
			name:   "400 - invalid to coins",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "0.1a",
						Hours:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawWalletRequest{
					ID: "foo.wlt",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - can't convert 0.1a to decimal",
		},

		{
			name:   "400 - invalid to hours",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "0.1",
						Hours:   "100.1",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawWalletRequest{
					ID: "foo.wlt",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid hours value: strconv.ParseUint: parsing \"100.1\": invalid syntax",
		},

		{
			name:   "400 - empty string to coins",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeAuto,
					Mode: wallet.HoursSelectionModeMatchCoins,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "",
						Hours:   "",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawWalletRequest{
					ID: "foo.wlt",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - can't convert  to decimal",
		},

		{
			name:   "400 - empty to",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawWalletRequest{
					ID: "foo.wlt",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to is empty",
		},

		{
			name:   "200 - auto type match coins",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeAuto,
					Mode: wallet.HoursSelectionModeMatchCoins,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawWalletRequest{
					ID: "foo.wlt",
				},
			},
			status: http.StatusOK,
			gatewayCreateTransactionResult: txn,
			createTransactionResult:        createTxnResult,
		},

		{
			name:   "200 - auto type split even",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeSplitEven,
					ShareFactor: newStrPtr("0.5"),
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawWalletRequest{
					ID: "foo.wlt",
				},
			},
			status: http.StatusOK,
			gatewayCreateTransactionResult: txn,
			createTransactionResult:        createTxnResult,
		},

		{
			name:   "200 - manual type zero hours",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "100",
						Hours:   "0",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawWalletRequest{
					ID: "foo.wlt",
				},
			},
			status: http.StatusOK,
			gatewayCreateTransactionResult: txn,
			createTransactionResult:        createTxnResult,
		},

		{
			name:   "200 - manual type nonzero hours",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "100",
						Hours:   "10",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawWalletRequest{
					ID: "foo.wlt",
				},
			},
			status: http.StatusOK,
			gatewayCreateTransactionResult: txn,
			createTransactionResult:        createTxnResult,
		},

		{
			name:   "200 - manual type nonzero hours - csrf disabled",
			method: http.MethodPost,
			body:   validBody,
			status: http.StatusOK,
			gatewayCreateTransactionResult: txn,
			createTransactionResult:        createTxnResult,
			csrfDisabled:                   true,
		},

		{
			name:   "500 - misc error",
			method: http.MethodPost,
			body:   validBody,
			status: http.StatusInternalServerError,
			gatewayCreateTransactionErr: errors.New("unhandled error"),
			err: "500 Internal Server Error - unhandled error",
		},

		{
			name:   "400 - no fee",
			method: http.MethodPost,
			body:   validBody,
			status: http.StatusBadRequest,
			gatewayCreateTransactionErr: fee.ErrTxnNoFee,
			err: "400 Bad Request - Transaction has zero coinhour fee",
		},

		{
			name:   "400 - insufficient coin hours",
			method: http.MethodPost,
			body:   validBody,
			status: http.StatusBadRequest,
			gatewayCreateTransactionErr: fee.ErrTxnInsufficientCoinHours,
			err: "400 Bad Request - Insufficient coinhours for transaction outputs",
		},

		{
			name:   "400 - other wallet error",
			method: http.MethodPost,
			body:   validBody,
			status: http.StatusBadRequest,
			gatewayCreateTransactionErr: wallet.ErrWalletEncrypted,
			err: "400 Bad Request - wallet is encrypted",
		},

		{
			name:   "404 - wallet not found",
			method: http.MethodPost,
			body:   validBody,
			status: http.StatusNotFound,
			gatewayCreateTransactionErr: wallet.ErrWalletNotExist,
			err: "404 Not Found - wallet doesn't exist",
		},

		{
			name:   "403 - wallet API disabled",
			method: http.MethodPost,
			body:   validBody,
			status: http.StatusForbidden,
			gatewayCreateTransactionErr: wallet.ErrWalletAPIDisabled,
			err: "403 Forbidden",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.gatewayCreateTransactionResult == nil {
				tc.gatewayCreateTransactionResult = &coin.Transaction{}
			}

			gateway := &GatewayerMock{}

			// If the rawRequestBody is can be deserialized to CreateTransactionRequest, use it to mock gateway.CreateTransaction
			serializedBody, err := json.Marshal(tc.body)
			require.NoError(t, err)
			var body CreateTransactionRequest
			err = json.Unmarshal(serializedBody, &body)
			if err == nil {
				gateway.On("CreateTransaction", body.ToWalletParams()).Return(tc.gatewayCreateTransactionResult, tc.gatewayCreateTransactionErr)
			}

			endpoint := "/wallet/transaction"

			requestJSON, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBuffer(requestJSON))
			require.NoError(t, err)

			contentType := tc.contentType
			if contentType == "" {
				contentType = "application/json"
			}

			req.Header.Add("Content-Type", contentType)

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}

			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()))
			} else {
				var msg CreateTransactionResult
				err := json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.NotNil(t, tc.createTransactionResult)
				require.Equal(t, *tc.createTransactionResult, msg)
			}
		})
	}
}

func newStrPtr(s string) *string {
	return &s
}
