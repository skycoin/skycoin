package gui

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"encoding/json"

	"math"

	"github.com/stretchr/testify/require"

	"errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
)

func makeBadBlock(t *testing.T) *coin.Block {
	genPublic, _ := cipher.GenerateKeyPair()
	genAddress := cipher.AddressFromPubKey(genPublic)
	var genCoins uint64 = 1000e6
	var genTime uint64 = 1000
	now := genTime + 100
	preBlock, err := coin.NewGenesisBlock(genAddress, genCoins, genTime)
	require.NoError(t, err)
	uxHash := testutil.RandSHA256(t)
	tx := coin.Transaction{}
	tx.PushOutput(genAddress, math.MaxInt64+1, 255)
	b, err := coin.NewBlock(*preBlock, now, uxHash, coin.Transactions{tx}, func(t *coin.Transaction) (uint64, error) {
		return 0, nil
	})
	require.NoError(t, err)
	return b
}

func TestGetBlock(t *testing.T) {

	badBlock := makeBadBlock(t)
	validHashString := testutil.RandSHA256(t).Hex()
	validSHA256, err := cipher.SHA256FromHex(validHashString)
	require.NoError(t, err)

	tt := []struct {
		name                        string
		method                      string
		status                      int
		err                         string
		hash                        string
		sha256                      cipher.SHA256
		seqStr                      string
		seq                         uint64
		gatewayGetBlockByHashResult coin.SignedBlock
		gatewayGetBlockByHashExists bool
		gatewayGetBlockBySeqResult  coin.SignedBlock
		gatewayGetBlockBySeqExists  bool
		response                    *visor.ReadableBlock
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - no seq and hash",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - should specify one filter, hash or seq",
		},
		{
			name:   "400 - seq and hash simultaneously",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - should only specify one filter, hash or seq",
			hash:   "hash",
			seqStr: "seq",
		},
		{
			name:   "400 - hash error: encoding/hex err invalid byte: U+0068 'h'",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - encoding/hex: invalid byte: U+0068 'h'",
			hash:   "hash",
		},
		{
			name:   "400 - hash error: encoding/hex: odd length hex string",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - encoding/hex: odd length hex string",
			hash:   "1234abc",
		},
		{
			name:   "400 - hash error: Invalid hex length",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid hex length",
			hash:   "1234abcd",
		},
		{
			name:   "404 - block by hash does not exist",
			method: http.MethodGet,
			status: http.StatusNotFound,
			err:    "404 Not Found",
			hash:   validHashString,
			sha256: validSHA256,
		},
		{
			name:   "200 - got block by hash",
			method: http.MethodGet,
			status: http.StatusOK,
			hash:   validHashString,
			sha256: validSHA256,
			gatewayGetBlockByHashExists: true,
			response: &visor.ReadableBlock{
				Head: visor.ReadableBlockHeader{
					BkSeq:             0x0,
					BlockHash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
					PreviousBlockHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Time:              0x0,
					Fee:               0x0,
					Version:           0x0,
					BodyHash:          "0000000000000000000000000000000000000000000000000000000000000000",
				},
				Body: visor.ReadableBlockBody{
					Transactions: []visor.ReadableTransaction{},
				},
			},
		},
		{
			name:   "400 - seq error: invalid syntax",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - strconv.ParseUint: parsing \"seq\": invalid syntax",
			seqStr: "seq",
		},
		{
			name:   "404 - block by seq does not exist",
			method: http.MethodGet,
			status: http.StatusNotFound,
			err:    "404 Not Found",
			seqStr: "1",
			seq:    1,
		},
		{
			name:   "500 - NewReadableBlock error",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error",
			seqStr: "1",
			seq:    1,
			gatewayGetBlockBySeqResult: coin.SignedBlock{
				Block: *badBlock,
			},
			gatewayGetBlockBySeqExists: true,
		},
		{
			name:   "200 - got block by seq",
			method: http.MethodGet,
			status: http.StatusOK,
			seqStr: "1",
			seq:    1,
			gatewayGetBlockBySeqExists: true,
			response: &visor.ReadableBlock{
				Head: visor.ReadableBlockHeader{
					BkSeq:             0x0,
					BlockHash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
					PreviousBlockHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Time:              0x0,
					Fee:               0x0,
					Version:           0x0,
					BodyHash:          "0000000000000000000000000000000000000000000000000000000000000000",
				},
				Body: visor.ReadableBlockBody{
					Transactions: []visor.ReadableTransaction{},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &GatewayerMock{}

			gateway.On("GetBlockByHash", tc.sha256).Return(tc.gatewayGetBlockByHashResult, tc.gatewayGetBlockByHashExists)
			gateway.On("GetBlockBySeq", tc.seq).Return(tc.gatewayGetBlockBySeqResult, tc.gatewayGetBlockBySeqExists)

			endpoint := "/block"

			v := url.Values{}
			if tc.hash != "" {
				v.Add("hash", tc.hash)
			}
			if tc.seqStr != "" {
				v.Add("seq", tc.seqStr)
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *visor.ReadableBlock
				err := json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.response, msg)
			}
		})
	}
}

func TestGetBlocks(t *testing.T) {
	type httpBody struct {
		Start string
		End   string
	}

	tt := []struct {
		name                   string
		method                 string
		status                 int
		err                    string
		body                   *httpBody
		start                  uint64
		end                    uint64
		gatewayGetBlocksResult *visor.ReadableBlocks
		gatewayGetBlocksError  error
		response               *visor.ReadableBlocks
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - empty start/end",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid start value \"\"",
		},
		{
			name:   "400 - bad start",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid start value \"badStart\"",
			body: &httpBody{
				Start: "badStart",
			},
		},
		{
			name:   "400 - bad end",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid end value \"badEnd\"",
			body: &httpBody{
				Start: "1",
				End:   "badEnd",
			},
			start: 1,
		},
		{
			name:   "400 - gatewayGetBlocksError",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Get blocks failed: gatewayGetBlocksError",
			body: &httpBody{
				Start: "1",
				End:   "3",
			},
			start: 1,
			end:   3,
			gatewayGetBlocksError: errors.New("gatewayGetBlocksError"),
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			body: &httpBody{
				Start: "1",
				End:   "3",
			},
			start: 1,
			end:   3,
			gatewayGetBlocksResult: &visor.ReadableBlocks{Blocks: []visor.ReadableBlock{visor.ReadableBlock{}}},
			response:               &visor.ReadableBlocks{Blocks: []visor.ReadableBlock{visor.ReadableBlock{}}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &GatewayerMock{}
			gateway.On("GetBlocks", tc.start, tc.end).Return(tc.gatewayGetBlocksResult, tc.gatewayGetBlocksError)

			endpoint := "/blocks"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.Start != "" {
					v.Add("start", tc.body.Start)
				}
				if tc.body.End != "" {
					v.Add("end", tc.body.End)
				}
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *visor.ReadableBlocks
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.response, msg)
			}
		})
	}
}

func TestGetLastBlocks(t *testing.T) {
	type httpBody struct {
		Num string
	}
	tt := []struct {
		name                       string
		method                     string
		url                        string
		status                     int
		err                        string
		body                       httpBody
		num                        uint64
		gatewayGetLastBlocksResult *visor.ReadableBlocks
		gatewayGetLastBlocksError  error
		response                   *visor.ReadableBlocks
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
			body: httpBody{
				Num: "1",
			},
			num: 1,
		},
		{
			name:   "400 - empty num value",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Param: num is empty",
			num:    1,
		},
		{
			name:   "400 - bad num value",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - strconv.ParseUint: parsing \"badNumValue\": invalid syntax",
			body: httpBody{
				Num: "badNumValue",
			},
			num: 1,
		},
		{
			name:   "400 - gatewayGetLastBlocksError",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Get last 1 blocks failed: gatewayGetLastBlocksError",
			body: httpBody{
				Num: "1",
			},
			num: 1,
			gatewayGetLastBlocksError: errors.New("gatewayGetLastBlocksError"),
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			body: httpBody{
				Num: "1",
			},
			num: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/last_blocks"
			gateway := NewGatewayerMock()

			gateway.On("GetLastBlocks", tc.num).Return(tc.gatewayGetLastBlocksResult, tc.gatewayGetLastBlocksError)

			v := url.Values{}
			if tc.body.Num != "" {
				v.Add("num", tc.body.Num)
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()

			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *visor.ReadableBlocks
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.response, msg)
			}
		})
	}
}
