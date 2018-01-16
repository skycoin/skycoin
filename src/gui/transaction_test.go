package gui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"math"

	"time"

	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor"
)

// GetAllUnconfirmedTxns returns all unconfirmed transactions
func (gw FakeGateway) GetAllUnconfirmedTxns() []visor.UnconfirmedTxn {
	args := gw.Called()
	return args.Get(0).([]visor.UnconfirmedTxn)
}

func createUnconfirmedTxn(t *testing.T) visor.UnconfirmedTxn {
	ut := visor.UnconfirmedTxn{}
	ut.Txn = coin.Transaction{}
	ut.Txn.InnerHash = testutil.RandSHA256(t)
	ut.Received = utc.Now().UnixNano()
	ut.Checked = ut.Received
	ut.Announced = time.Time{}.UnixNano()
	return ut
}

func TestGetPendingTxs(t *testing.T) {
	invalidTxn := createUnconfirmedTxn(t)
	invalidTxn.Txn.Out = append(invalidTxn.Txn.Out, coin.TransactionOutput{
		Coins: math.MaxInt64 + 1,
	})

	tt := []struct {
		name                          string
		method                        string
		url                           string
		status                        int
		err                           string
		getAllUnconfirmedTxnsResponse []visor.UnconfirmedTxn
		httpResponse                  []*visor.ReadableUnconfirmedTxn
	}{
		{
			"405",
			http.MethodPost,
			"/pendingTxs",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			[]visor.UnconfirmedTxn{},
			nil,
		},
		{
			"500 - bad unconfirmedTxn",
			http.MethodGet,
			"/pendingTxs",
			http.StatusInternalServerError,
			"500 Internal Server Error",
			[]visor.UnconfirmedTxn{
				invalidTxn,
			},
			nil,
		},
		{
			"200",
			http.MethodGet,
			"/pendingTxs",
			http.StatusOK,
			"",
			[]visor.UnconfirmedTxn{},
			[]*visor.ReadableUnconfirmedTxn{},
		},
	}

	for _, tc := range tt {
		gateway := &FakeGateway{
			t: t,
		}
		gateway.On("GetAllUnconfirmedTxns").Return(tc.getAllUnconfirmedTxnsResponse)

		req, err := http.NewRequest(tc.method, tc.url, nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(getPendingTxs(gateway))

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
				tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg []*visor.ReadableUnconfirmedTxn
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			require.Equal(t, tc.httpResponse, msg, tc.name)
		}
	}
}
