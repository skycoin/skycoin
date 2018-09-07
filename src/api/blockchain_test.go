package api

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
	tx := coin.Transaction{
		In: []cipher.SHA256{
			testutil.RandSHA256(t),
		},
	}
	tx.PushOutput(genAddress, math.MaxInt64+1, 255)
	b, err := coin.NewBlock(*preBlock, now, uxHash, coin.Transactions{tx}, func(t *coin.Transaction) (uint64, error) {
		return 0, nil
	})
	require.NoError(t, err)
	require.NotEqual(t, b.Head.BkSeq, uint64(0))
	require.NotEmpty(t, b.Body.Transactions)
	for i, txn := range b.Body.Transactions {
		require.NotEmpty(t, txn.In, "txn %d/%d", i+1, len(b.Body.Transactions))
	}
	return b
}

func TestGetBlock(t *testing.T) {
	badBlock := makeBadBlock(t)
	validHashString := testutil.RandSHA256(t).Hex()
	validSHA256, err := cipher.SHA256FromHex(validHashString)
	require.NoError(t, err)

	tt := []struct {
		name                               string
		method                             string
		status                             int
		err                                string
		hash                               string
		sha256                             cipher.SHA256
		seqStr                             string
		seq                                uint64
		verbose                            bool
		verboseStr                         string
		gatewayGetBlockByHashResult        *coin.SignedBlock
		gatewayGetBlockByHashErr           error
		gatewayGetBlockBySeqResult         *coin.SignedBlock
		gatewayGetBlockBySeqErr            error
		gatewayGetBlockByHashVerboseResult *visor.ReadableBlockVerbose
		gatewayGetBlockByHashVerboseErr    error
		gatewayGetBlockBySeqVerboseResult  *visor.ReadableBlockVerbose
		gatewayGetBlockBySeqVerboseErr     error
		response                           interface{}
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
			name:   "400 - seq error: invalid syntax",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid seq value \"badseq\"",
			seqStr: "badseq",
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
			err:    "500 Internal Server Error - Droplet string conversion failed: Value is too large",
			seqStr: "1",
			seq:    1,
			gatewayGetBlockBySeqResult: &coin.SignedBlock{
				Block: *badBlock,
			},
		},
		{
			name:                     "500 - get block by hash error",
			method:                   http.MethodGet,
			status:                   http.StatusInternalServerError,
			err:                      "500 Internal Server Error - GetSignedBlockByHash failed",
			hash:                     validHashString,
			sha256:                   validSHA256,
			gatewayGetBlockByHashErr: errors.New("GetSignedBlockByHash failed"),
		},
		{
			name:                    "500 - get block by seq error",
			method:                  http.MethodGet,
			status:                  http.StatusInternalServerError,
			err:                     "500 Internal Server Error - GetSignedBlockBySeq failed",
			seqStr:                  "1",
			seq:                     1,
			gatewayGetBlockBySeqErr: errors.New("GetSignedBlockBySeq failed"),
		},
		{
			name:                       "200 - get block by seq",
			method:                     http.MethodGet,
			status:                     http.StatusOK,
			seqStr:                     "1",
			seq:                        1,
			gatewayGetBlockBySeqResult: &coin.SignedBlock{},
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
			name:                        "200 - get block by hash",
			method:                      http.MethodGet,
			status:                      http.StatusOK,
			hash:                        validHashString,
			sha256:                      validSHA256,
			gatewayGetBlockByHashResult: &coin.SignedBlock{},
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
			name:       "200 - get block by hash verbose",
			method:     http.MethodGet,
			status:     http.StatusOK,
			hash:       validHashString,
			sha256:     validSHA256,
			verbose:    true,
			verboseStr: "1",
			gatewayGetBlockByHashVerboseResult: func() *visor.ReadableBlockVerbose {
				b, err := visor.NewReadableBlockVerbose(&coin.Block{}, nil)
				require.NoError(t, err)
				return b
			}(),
			response: &visor.ReadableBlockVerbose{
				Head: visor.ReadableBlockHeader{
					BkSeq:             0x0,
					BlockHash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
					PreviousBlockHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Time:              0x0,
					Fee:               0x0,
					Version:           0x0,
					BodyHash:          "0000000000000000000000000000000000000000000000000000000000000000",
				},
				Body: visor.ReadableBlockBodyVerbose{
					Transactions: []visor.ReadableBlockTransactionVerbose{},
				},
			},
		},

		{
			name:       "200 - get block by seq verbose",
			method:     http.MethodGet,
			status:     http.StatusOK,
			seq:        1,
			seqStr:     "1",
			verbose:    true,
			verboseStr: "1",
			gatewayGetBlockBySeqVerboseResult: func() *visor.ReadableBlockVerbose {
				b, err := visor.NewReadableBlockVerbose(&coin.Block{}, nil)
				require.NoError(t, err)
				return b
			}(),
			response: &visor.ReadableBlockVerbose{
				Head: visor.ReadableBlockHeader{
					BkSeq:             0x0,
					BlockHash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
					PreviousBlockHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Time:              0x0,
					Fee:               0x0,
					Version:           0x0,
					BodyHash:          "0000000000000000000000000000000000000000000000000000000000000000",
				},
				Body: visor.ReadableBlockBodyVerbose{
					Transactions: []visor.ReadableBlockTransactionVerbose{},
				},
			},
		},

		{
			name:                            "500 - get block by hash verbose error",
			method:                          http.MethodGet,
			status:                          http.StatusInternalServerError,
			hash:                            validHashString,
			sha256:                          validSHA256,
			verbose:                         true,
			verboseStr:                      "1",
			gatewayGetBlockByHashVerboseErr: errors.New("GetBlockByHashVerbose failed"),
			err:                             "500 Internal Server Error - GetBlockByHashVerbose failed",
		},

		{
			name:                           "500 - get block by seq verbose error",
			method:                         http.MethodGet,
			status:                         http.StatusInternalServerError,
			seq:                            1,
			seqStr:                         "1",
			verbose:                        true,
			verboseStr:                     "1",
			gatewayGetBlockBySeqVerboseErr: errors.New("GetBlockBySeqVerbose failed"),
			err:                            "500 Internal Server Error - GetBlockBySeqVerbose failed",
		},

		{
			name:       "404 - get block by hash verbose not found",
			method:     http.MethodGet,
			status:     http.StatusNotFound,
			hash:       validHashString,
			sha256:     validSHA256,
			verbose:    true,
			verboseStr: "1",
			err:        "404 Not Found",
		},

		{
			name:       "404 - get block by seq verbose not found",
			method:     http.MethodGet,
			status:     http.StatusNotFound,
			seq:        1,
			seqStr:     "1",
			verbose:    true,
			verboseStr: "1",
			err:        "404 Not Found",
		},

		{
			name:       "400 - invalid verbose flag",
			method:     http.MethodGet,
			status:     http.StatusBadRequest,
			seq:        1,
			seqStr:     "1",
			verbose:    true,
			verboseStr: "asdasdasd",
			err:        "400 Bad Request - Invalid value for verbose",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}

			gateway.On("GetSignedBlockByHash", tc.sha256).Return(tc.gatewayGetBlockByHashResult, tc.gatewayGetBlockByHashErr)
			gateway.On("GetSignedBlockBySeq", tc.seq).Return(tc.gatewayGetBlockBySeqResult, tc.gatewayGetBlockBySeqErr)
			gateway.On("GetBlockByHashVerbose", tc.sha256).Return(tc.gatewayGetBlockByHashVerboseResult, tc.gatewayGetBlockByHashVerboseErr)
			gateway.On("GetBlockBySeqVerbose", tc.seq).Return(tc.gatewayGetBlockBySeqVerboseResult, tc.gatewayGetBlockBySeqVerboseErr)

			endpoint := "/api/v1/block"

			v := url.Values{}
			if tc.hash != "" {
				v.Add("hash", tc.hash)
			}
			if tc.seqStr != "" {
				v.Add("seq", tc.seqStr)
			}
			if tc.verboseStr != "" {
				v.Add("verbose", tc.verboseStr)
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
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()))
			} else {
				if tc.verbose {
					var msg *visor.ReadableBlockVerbose
					err := json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.response.(*visor.ReadableBlockVerbose), msg)
				} else {
					var msg *visor.ReadableBlock
					err := json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.response.(*visor.ReadableBlock), msg)
				}
			}
		})
	}
}

func TestGetBlocks(t *testing.T) {
	type httpBody struct {
		Start   string
		End     string
		Verbose string
	}

	tt := []struct {
		name                          string
		method                        string
		status                        int
		err                           string
		body                          *httpBody
		start                         uint64
		end                           uint64
		verbose                       bool
		gatewayGetBlocksResult        *visor.ReadableBlocks
		gatewayGetBlocksError         error
		gatewayGetBlocksVerboseResult *visor.ReadableBlocksVerbose
		gatewayGetBlocksVerboseError  error
		response                      interface{}
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
			name:   "400 - bad verbose",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid value for verbose",
			body: &httpBody{
				Start:   "1",
				End:     "2",
				Verbose: "foo",
			},
		},
		{
			name:   "500 - gatewayGetBlocksError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gatewayGetBlocksError",
			body: &httpBody{
				Start: "1",
				End:   "3",
			},
			start:                 1,
			end:                   3,
			gatewayGetBlocksError: errors.New("gatewayGetBlocksError"),
		},
		{
			name:   "500 - gatewayGetBlocksVerboseError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gatewayGetBlocksVerboseError",
			body: &httpBody{
				Start:   "1",
				End:     "3",
				Verbose: "1",
			},
			start:                        1,
			end:                          3,
			verbose:                      true,
			gatewayGetBlocksVerboseError: errors.New("gatewayGetBlocksVerboseError"),
		},

		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			body: &httpBody{
				Start: "1",
				End:   "3",
			},
			start:                  1,
			end:                    3,
			gatewayGetBlocksResult: &visor.ReadableBlocks{Blocks: []visor.ReadableBlock{visor.ReadableBlock{}}},
			response:               &visor.ReadableBlocks{Blocks: []visor.ReadableBlock{visor.ReadableBlock{}}},
		},
		{
			name:   "200 verbose",
			method: http.MethodGet,
			status: http.StatusOK,
			body: &httpBody{
				Start:   "1",
				End:     "3",
				Verbose: "1",
			},
			start:   1,
			end:     3,
			verbose: true,
			gatewayGetBlocksVerboseResult: &visor.ReadableBlocksVerbose{
				Blocks: []visor.ReadableBlockVerbose{
					visor.ReadableBlockVerbose{},
				},
			},
			response: &visor.ReadableBlocksVerbose{
				Blocks: []visor.ReadableBlockVerbose{
					visor.ReadableBlockVerbose{},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("GetBlocks", tc.start, tc.end).Return(tc.gatewayGetBlocksResult, tc.gatewayGetBlocksError)
			gateway.On("GetBlocksVerbose", tc.start, tc.end).Return(tc.gatewayGetBlocksVerboseResult, tc.gatewayGetBlocksVerboseError)

			endpoint := "/api/v1/blocks"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.Start != "" {
					v.Add("start", tc.body.Start)
				}
				if tc.body.End != "" {
					v.Add("end", tc.body.End)
				}
				if tc.body.Verbose != "" {
					v.Add("verbose", tc.body.Verbose)
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
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				if tc.verbose {
					var msg *visor.ReadableBlocksVerbose
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.response, msg)
				} else {
					var msg *visor.ReadableBlocks
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.response, msg)
				}
			}
		})
	}
}

func TestGetLastBlocks(t *testing.T) {
	type httpBody struct {
		Num     string
		Verbose string
	}
	tt := []struct {
		name                              string
		method                            string
		status                            int
		err                               string
		body                              httpBody
		num                               uint64
		verbose                           bool
		gatewayGetLastBlocksResult        *visor.ReadableBlocks
		gatewayGetLastBlocksError         error
		gatewayGetLastBlocksVerboseResult *visor.ReadableBlocksVerbose
		gatewayGetLastBlocksVerboseError  error
		response                          interface{}
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
			body: httpBody{
				Num: "1",
			},
		},
		{
			name:   "400 - empty num value",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid num value \"\"",
		},
		{
			name:   "400 - bad num value",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid num value \"badNumValue\"",
			body: httpBody{
				Num: "badNumValue",
			},
		},
		{
			name:   "400 - bad verbose",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid value for verbose",
			body: httpBody{
				Num:     "1",
				Verbose: "foo",
			},
		},
		{
			name:   "500 - gatewayGetLastBlocksError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gatewayGetLastBlocksError",
			body: httpBody{
				Num: "1",
			},
			num:                       1,
			gatewayGetLastBlocksError: errors.New("gatewayGetLastBlocksError"),
		},
		{
			name:   "500 - gatewayGetLastBlocksVerboseError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gatewayGetLastBlocksVerboseError",
			body: httpBody{
				Num:     "1",
				Verbose: "1",
			},
			num:                              1,
			verbose:                          true,
			gatewayGetLastBlocksVerboseError: errors.New("gatewayGetLastBlocksVerboseError"),
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			body: httpBody{
				Num: "1",
			},
			num: 1,
			gatewayGetLastBlocksResult: &visor.ReadableBlocks{
				Blocks: []visor.ReadableBlock{
					visor.ReadableBlock{},
				},
			},
			response: &visor.ReadableBlocks{
				Blocks: []visor.ReadableBlock{
					visor.ReadableBlock{},
				},
			},
		},
		{
			name:   "200 verbose",
			method: http.MethodGet,
			status: http.StatusOK,
			body: httpBody{
				Num:     "1",
				Verbose: "1",
			},
			num:     1,
			verbose: true,
			gatewayGetLastBlocksVerboseResult: &visor.ReadableBlocksVerbose{
				Blocks: []visor.ReadableBlockVerbose{
					visor.ReadableBlockVerbose{},
				},
			},
			response: &visor.ReadableBlocksVerbose{
				Blocks: []visor.ReadableBlockVerbose{
					visor.ReadableBlockVerbose{},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/last_blocks"
			gateway := &MockGatewayer{}

			gateway.On("GetLastBlocks", tc.num).Return(tc.gatewayGetLastBlocksResult, tc.gatewayGetLastBlocksError)
			gateway.On("GetLastBlocksVerbose", tc.num).Return(tc.gatewayGetLastBlocksVerboseResult, tc.gatewayGetLastBlocksVerboseError)

			v := url.Values{}
			if tc.body.Num != "" {
				v.Add("num", tc.body.Num)
			}
			if tc.body.Verbose != "" {
				v.Add("verbose", tc.body.Verbose)
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

			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				if tc.verbose {
					var msg *visor.ReadableBlocksVerbose
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.response, msg)
				} else {
					var msg *visor.ReadableBlocks
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.response, msg)
				}
			}
		})
	}
}
