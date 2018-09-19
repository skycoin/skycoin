package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/wallet"
)

func TestCreateTransaction(t *testing.T) {
	type rawRequestWallet struct {
		ID        string   `json:"id"`
		UxOuts    []string `json:"unspents,omitempty"`
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
		Wallet         rawRequestWallet  `json:"wallet"`
		ChangeAddress  string            `json:"change_address,omitempty"`
		To             []rawReceiver     `json:"to"`
		Password       string            `json:"password"`
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

	inputs := []wallet.UxBalance{
		{
			Hash:           testutil.RandSHA256(t),
			Time:           uint64(time.Now().UTC().Unix()),
			BkSeq:          9999,
			SrcTransaction: testutil.RandSHA256(t),
			Address:        testutil.MakeAddress(),
			Coins:          1e6,
			Hours:          200,
			InitialHours:   100,
		},
	}

	createdTxn, err := NewCreatedTransaction(txn, inputs)
	require.NoError(t, err)

	createTxnResponse := &CreateTransactionResponse{
		Transaction:        *createdTxn,
		EncodedTransaction: hex.EncodeToString(txn.Serialize()),
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
		Wallet: rawRequestWallet{
			ID: "foo.wlt",
		},
	}

	walletInput := testutil.RandSHA256(t)

	tt := []struct {
		name                           string
		method                         string
		body                           *rawRequest
		status                         int
		err                            string
		gatewayCreateTransactionResult *coin.Transaction
		gatewayCreateTransactionInputs []wallet.UxBalance
		gatewayCreateTransactionErr    error
		createTransactionResponse      *CreateTransactionResponse
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
					Mode: wallet.HoursSelectionModeShare,
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing hours_selection.share_factor when hours_selection.mode is share",
		},

		{
			name:   "400 - share factor set but mode is not share",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeManual,
					ShareFactor: newStrPtr("0.5"),
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - hours_selection.share_factor can only be used when hours_selection.mode is share",
		},

		{
			name:   "400 - negative share factor",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeShare,
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
					Mode:        wallet.HoursSelectionModeShare,
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
				Wallet: rawRequestWallet{
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
				Wallet: rawRequestWallet{
					Addresses: []string{"xxx"},
				},
				ChangeAddress: changeAddress.String(),
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid address: Invalid address length",
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
			err:    "400 Bad Request - change_address must not be the null address",
		},

		{
			name:   "400 - auto type destination has hours",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeShare,
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
					Mode: wallet.HoursSelectionModeShare,
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
				Wallet:        rawRequestWallet{},
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
				Wallet: rawRequestWallet{
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
				Wallet: rawRequestWallet{
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
				Wallet: rawRequestWallet{
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
				Wallet: rawRequestWallet{
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
				Wallet: rawRequestWallet{
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
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeShare,
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
				Wallet: rawRequestWallet{
					ID: "foo.wlt",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - can't convert  to decimal",
		},

		{
			name:   "400 - coins has too many decimals",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.1234",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawRequestWallet{
					ID: "foo.wlt",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to[0].coins has too many decimal places",
		},

		{
			name:   "400 - empty to",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawRequestWallet{
					ID: "foo.wlt",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to is empty",
		},

		{
			name:   "400 - manual duplicate outputs",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type: wallet.HoursSelectionTypeManual,
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawRequestWallet{
					ID: "foo.wlt",
				},
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
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to contains duplicate values",
		},

		{
			name:   "400 - auto duplicate outputs",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawRequestWallet{
					ID: "foo.wlt",
				},
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
			status: http.StatusBadRequest,
			err:    "400 Bad Request - to contains duplicate values",
		},

		{
			name:   "400 - both wallet uxouts and wallet addresses specified",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawRequestWallet{
					ID:        "foo.wlt",
					Addresses: []string{destinationAddress.String()},
					UxOuts:    []string{walletInput.Hex()},
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.2",
					},
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - wallet.unspents and wallet.addresses cannot be combined",
		},

		{
			name:   "400 - duplicate wallet uxouts",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawRequestWallet{
					ID:     "foo.wlt",
					UxOuts: []string{walletInput.Hex(), walletInput.Hex()},
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.2",
					},
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - wallet.unspents contains duplicate values",
		},

		{
			name:   "400 - duplicate wallet addresses",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawRequestWallet{
					ID:        "foo.wlt",
					Addresses: []string{destinationAddress.String(), destinationAddress.String()},
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "1.2",
					},
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - wallet.addresses contains duplicate values",
		},

		{
			name:   "200 - auto type split even",
			method: http.MethodPost,
			body: &rawRequest{
				HoursSelection: rawHoursSelection{
					Type:        wallet.HoursSelectionTypeAuto,
					Mode:        wallet.HoursSelectionModeShare,
					ShareFactor: newStrPtr("0.5"),
				},
				To: []rawReceiver{
					{
						Address: destinationAddress.String(),
						Coins:   "100",
					},
				},
				ChangeAddress: changeAddress.String(),
				Wallet: rawRequestWallet{
					ID: "foo.wlt",
				},
			},
			status: http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			createTransactionResponse:      createTxnResponse,
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
				Wallet: rawRequestWallet{
					ID: "foo.wlt",
				},
			},
			status: http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			createTransactionResponse:      createTxnResponse,
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
				Wallet: rawRequestWallet{
					ID: "foo.wlt",
				},
			},
			status: http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			createTransactionResponse:      createTxnResponse,
		},

		{
			name:   "200 - manual type nonzero hours - csrf disabled",
			method: http.MethodPost,
			body:   validBody,
			status: http.StatusOK,
			gatewayCreateTransactionResult: txn,
			gatewayCreateTransactionInputs: inputs,
			createTransactionResponse:      createTxnResponse,
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
			name:   "400 - uxout doesn't exist",
			method: http.MethodPost,
			body:   validBody,
			status: http.StatusBadRequest,
			gatewayCreateTransactionErr: blockdb.NewErrUnspentNotExist("foo"),
			err: "400 Bad Request - unspent output of foo does not exist",
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
			gateway := &MockGatewayer{}

			// If the rawRequestBody can be deserialized to CreateTransactionRequest, use it to mock gateway.CreateTransaction
			serializedBody, err := json.Marshal(tc.body)
			require.NoError(t, err)
			var body createTransactionRequest
			err = json.Unmarshal(serializedBody, &body)
			if err == nil {
				gateway.On("CreateTransaction", body.ToWalletParams()).Return(tc.gatewayCreateTransactionResult, tc.gatewayCreateTransactionInputs, tc.gatewayCreateTransactionErr)
			}

			endpoint := "/api/v1/wallet/transaction"

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
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

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
