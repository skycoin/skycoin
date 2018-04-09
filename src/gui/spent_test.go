package gui

import (
	"testing"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
	"net/http"
	"bytes"
	"github.com/stretchr/testify/require"
	"encoding/json"
	"net/http/httptest"
	"strings"
)

func TestAdvancedSpend(t *testing.T) {
	tt := []struct {
		name                       string
		method                     string
		body                       *wallet.AdvancedSpendRequest
		status                     int
		err                        string
		gatewayAdvancedSpendResult *coin.Transaction
		gatewayAdvnacedSpendErr    error
		gatewayGetWalletResult     wallet.Wallets
		gatewayGetWalletErr        error
		advancedSpendResult        *AdvancedSpendResult
		csrfDisabled               bool
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - missing hours selection type",
			method: http.MethodPost,
			body:   &wallet.AdvancedSpendRequest{},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing hours selection type",
		},
		{
			name:   "400 - missing hours selection type",
			method: http.MethodPost,
			body: &wallet.AdvancedSpendRequest{
				HoursSelection: wallet.HoursSelection{
					Type: "",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing hours selection type",
		},
		{
			name:   "400 - missing hours selection mode",
			method: http.MethodPost,
			body: &wallet.AdvancedSpendRequest{
				HoursSelection: wallet.HoursSelection{
					Type: "auto",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing hours selection mode when type is auto",
		},
		{
			name:   "400 -  missing hours selection share factor",
			method: http.MethodPost,
			body: &wallet.AdvancedSpendRequest{
				HoursSelection: wallet.HoursSelection{
					Type: "auto",
					Mode: "split_even",
				},
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing hours selection share factor when mode is split_even",
		},
		{
			name:   "400 - no sender address provided",
			method: http.MethodPost,
			body: &wallet.AdvancedSpendRequest{
				HoursSelection: wallet.HoursSelection{
					Type: "manual",
				},
				Addresses:     []string{},
				ChangeAddress: "2mEgmYt6NZHA1erYqbAeXmGPD5gqLZ9toFv",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - no sender addresses found",
		},
		{
			name:   "400 - empty sender address",
			method: http.MethodPost,
			body: &wallet.AdvancedSpendRequest{
				HoursSelection: wallet.HoursSelection{
					Type: "manual",
				},
				Addresses:     []string{""},
				ChangeAddress: "2mEgmYt6NZHA1erYqbAeXmGPD5gqLZ9toFv",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - empty sender address",
		},
		{
			name:   "400 - invalid change address",
			method: http.MethodPost,
			body: &wallet.AdvancedSpendRequest{
				HoursSelection: wallet.HoursSelection{
					Type: "manual",
				},
				Addresses:     []string{"tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V"},
				ChangeAddress: "xxxx",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid change address: Invalid address length",
		},
		{
			name:   "400 - invalid change address",
			method: http.MethodPost,
			body: &wallet.AdvancedSpendRequest{
				HoursSelection: wallet.HoursSelection{
					Type: "manual",
				},
				Addresses:     []string{"tWPDM36ex9zLjJw1aPMfYTVPbYgkL2Xp9V"},
				ChangeAddress: "2mEgmYt6NZHA1erYqbAeXmGPD5gqLZ9toFx",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid change address: Invalid checksum",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.gatewayAdvancedSpendResult == nil {
				tc.gatewayAdvancedSpendResult = &coin.Transaction{}
			}

			gateway := &GatewayerMock{}
			gateway.On("AdvancedSpend", tc.body).Return(tc.gatewayAdvancedSpendResult, tc.gatewayAdvnacedSpendErr)
			gateway.On("GetWallets").Return(tc.gatewayGetWalletResult, tc.gatewayGetWalletErr)

			endpoint := "/spend/advanced"

			requestJSON, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBuffer(requestJSON))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/json")

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
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %d, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg AdvancedSpendResult
				err := json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, *tc.advancedSpendResult, msg)
			}
		})
	}
}
