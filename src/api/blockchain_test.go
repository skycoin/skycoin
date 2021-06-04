package api

import (
	"io"
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
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
)

func TestGetBlockchainMetadata(t *testing.T) {
	cases := []struct {
		name                        string
		method                      string
		status                      int
		err                         string
		getBlockchainMetadataResult *visor.BlockchainMetadata
		getBlockchainMetadataErr    error
		result                      readable.BlockchainMetadata
	}{
		{
			name:   "405 method not allowed",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                     "500 - GetBlockchainMetadata error",
			method:                   http.MethodGet,
			status:                   http.StatusInternalServerError,
			err:                      "500 Internal Server Error - gateway.GetBlockchainMetadata failed: GetBlockchainMetadata error",
			getBlockchainMetadataErr: errors.New("GetBlockchainMetadata error"),
		},
		{
			name:   "500 - nil visor.BlockchainMetadata",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gateway.GetBlockchainMetadata metadata is nil",
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			getBlockchainMetadataResult: &visor.BlockchainMetadata{
				HeadBlock:   coin.SignedBlock{},
				Unspents:    12,
				Unconfirmed: 13,
			},
			result: readable.BlockchainMetadata{
				Head: readable.BlockHeader{
					Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
					PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
					BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
					UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
				},
				Unspents:    12,
				Unconfirmed: 13,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("GetBlockchainMetadata").Return(tc.getBlockchainMetadataResult, tc.getBlockchainMetadataErr)

			endpoint := "/api/v1/blockchain/metadata"
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			setCSRFParameters(t, tokenValid, req)

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v` %s", status, tc.status, strings.TrimSpace(rr.Body.String()))

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg readable.BlockchainMetadata
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestGetBlockchainProgress(t *testing.T) {
	addr1 := testutil.MakeAddress()
	addr2 := testutil.MakeAddress()

	cases := []struct {
		name                        string
		method                      string
		status                      int
		err                         string
		headBkSeq                   uint64
		headBkSeqErr                error
		getBlockchainProgressResult *daemon.BlockchainProgress
		result                      readable.BlockchainProgress
	}{
		{
			name:   "405 method not allowed",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},

		{
			name:         "500 - HeadBkSeq error",
			method:       http.MethodGet,
			status:       http.StatusInternalServerError,
			err:          "500 Internal Server Error - gateway.HeadBkSeq failed: HeadBkSeq error",
			headBkSeqErr: errors.New("HeadBkSeq error"),
		},

		{
			name:   "500 - nil daemon.BlockchainProgress",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gateway.GetBlockchainProgress progress is nil",
		},

		{
			name:      "200",
			method:    http.MethodGet,
			status:    http.StatusOK,
			headBkSeq: 99,
			getBlockchainProgressResult: &daemon.BlockchainProgress{
				Peers: []daemon.PeerBlockchainHeight{
					{
						Address: addr1.String(),
						Height:  101,
					},
					{
						Address: addr2.String(),
						Height:  102,
					},
				},
				Current: 99,
				Highest: 102,
			},
			result: readable.BlockchainProgress{
				Peers: []readable.PeerBlockchainHeight{
					{
						Address: addr1.String(),
						Height:  101,
					},
					{
						Address: addr2.String(),
						Height:  102,
					},
				},
				Current: 99,
				Highest: 102,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("HeadBkSeq").Return(tc.headBkSeq, true, tc.headBkSeqErr)
			gateway.On("GetBlockchainProgress", tc.headBkSeq).Return(tc.getBlockchainProgressResult)

			endpoint := "/api/v1/blockchain/progress"
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			setCSRFParameters(t, tokenValid, req)

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v` %s", status, tc.status, strings.TrimSpace(rr.Body.String()))

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg readable.BlockchainProgress
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func makeBadBlock(t *testing.T) *coin.Block {
	genPublic, _ := cipher.GenerateKeyPair()
	genAddress := cipher.AddressFromPubKey(genPublic)
	var genCoins uint64 = 1000e6
	var genTime uint64 = 1000
	now := genTime + 100
	preBlock, err := coin.NewGenesisBlock(genAddress, genCoins, genTime)
	require.NoError(t, err)
	uxHash := testutil.RandSHA256(t)
	txn := coin.Transaction{
		In: []cipher.SHA256{
			testutil.RandSHA256(t),
		},
	}
	err = txn.PushOutput(genAddress, math.MaxInt64+1, 255)
	require.NoError(t, err)
	b, err := coin.NewBlock(*preBlock, now, uxHash, coin.Transactions{txn}, func(t *coin.Transaction) (uint64, error) {
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

	type verboseResult struct {
		Block  *coin.SignedBlock
		Inputs [][]visor.TransactionInput
	}

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
		gatewayGetBlockByHashVerboseResult verboseResult
		gatewayGetBlockByHashVerboseErr    error
		gatewayGetBlockBySeqVerboseResult  verboseResult
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
			name:   "500 - readable.NewBlock error",
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
			response: &readable.Block{
				Head: readable.BlockHeader{
					BkSeq:        0x0,
					Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
					PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Time:         0x0,
					Fee:          0x0,
					Version:      0x0,
					BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
					UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
				},
				Body: readable.BlockBody{
					Transactions: []readable.Transaction{},
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
			response: &readable.Block{
				Head: readable.BlockHeader{
					BkSeq:        0x0,
					Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
					PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Time:         0x0,
					Fee:          0x0,
					Version:      0x0,
					BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
					UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
				},
				Body: readable.BlockBody{
					Transactions: []readable.Transaction{},
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
			gatewayGetBlockByHashVerboseResult: verboseResult{
				Block:  &coin.SignedBlock{},
				Inputs: nil,
			},
			response: &readable.BlockVerbose{
				Head: readable.BlockHeader{
					BkSeq:        0x0,
					Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
					PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Time:         0x0,
					Fee:          0x0,
					Version:      0x0,
					BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
					UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
				},
				Body: readable.BlockBodyVerbose{
					Transactions: []readable.BlockTransactionVerbose{},
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
			gatewayGetBlockBySeqVerboseResult: verboseResult{
				Block:  &coin.SignedBlock{},
				Inputs: nil,
			},
			response: &readable.BlockVerbose{
				Head: readable.BlockHeader{
					BkSeq:        0x0,
					Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
					PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Time:         0x0,
					Fee:          0x0,
					Version:      0x0,
					BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
					UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
				},
				Body: readable.BlockBodyVerbose{
					Transactions: []readable.BlockTransactionVerbose{},
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
			gatewayGetBlockByHashVerboseErr: errors.New("GetSignedBlockByHashVerbose failed"),
			err:                             "500 Internal Server Error - GetSignedBlockByHashVerbose failed",
		},

		{
			name:                           "500 - get block by seq verbose error",
			method:                         http.MethodGet,
			status:                         http.StatusInternalServerError,
			seq:                            1,
			seqStr:                         "1",
			verbose:                        true,
			verboseStr:                     "1",
			gatewayGetBlockBySeqVerboseErr: errors.New("GetSignedBlockBySeqVerbose failed"),
			err:                            "500 Internal Server Error - GetSignedBlockBySeqVerbose failed",
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
			gateway.On("GetSignedBlockByHashVerbose", tc.sha256).Return(tc.gatewayGetBlockByHashVerboseResult.Block,
				tc.gatewayGetBlockByHashVerboseResult.Inputs, tc.gatewayGetBlockByHashVerboseErr)
			gateway.On("GetSignedBlockBySeqVerbose", tc.seq).Return(tc.gatewayGetBlockBySeqVerboseResult.Block,
				tc.gatewayGetBlockBySeqVerboseResult.Inputs, tc.gatewayGetBlockBySeqVerboseErr)

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
			req.Header.Add("Content-Type", ContentTypeForm)

			setCSRFParameters(t, tokenValid, req)

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()))
			} else {
				if tc.verbose {
					var msg *readable.BlockVerbose
					err := json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.response.(*readable.BlockVerbose), msg)
				} else {
					var msg *readable.Block
					err := json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.response.(*readable.Block), msg)
				}
			}
		})
	}
}

func TestGetBlocks(t *testing.T) {
	type httpBody struct {
		Start   string
		End     string
		Seqs    string
		Verbose string
	}

	type verboseResult struct {
		Blocks []coin.SignedBlock
		Inputs [][][]visor.TransactionInput
	}

	tt := []struct {
		name                                 string
		method                               string
		status                               int
		err                                  string
		body                                 *httpBody
		start                                uint64
		end                                  uint64
		seqs                                 []uint64
		verbose                              bool
		gatewayGetBlocksInRangeResult        []coin.SignedBlock
		gatewayGetBlocksInRangeError         error
		gatewayGetBlocksInRangeVerboseResult verboseResult
		gatewayGetBlocksInRangeVerboseError  error
		gatewayGetBlocksResult               []coin.SignedBlock
		gatewayGetBlocksError                error
		gatewayGetBlocksVerboseResult        verboseResult
		gatewayGetBlocksVerboseError         error
		response                             interface{}
	}{
		{
			name:   "405",
			method: http.MethodDelete,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - empty start, end and seqs",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - At least one of seqs or start or end are required",
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
			name:   "400 - seqs combined with start",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - seqs cannot be used with start or end",
			body: &httpBody{
				Seqs:  "1,2,3",
				Start: "1",
			},
		},
		{
			name:   "400 - seqs combined with end",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - seqs cannot be used with start or end",
			body: &httpBody{
				Seqs: "1,2,3",
				End:  "1",
			},
		},
		{
			name:   "400 - bad seqs",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid sequence \"a\" at seqs[2]",
			body: &httpBody{
				Seqs: "1,2,a",
			},
		},
		{
			name:   "400 - bad seqs",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid sequence \"\" at seqs[1]",
			body: &httpBody{
				Seqs: "1,,2",
			},
		},
		{
			name:   "400 - bad seqs",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid sequence \"foo\" at seqs[0]",
			body: &httpBody{
				Seqs: "foo",
			},
		},
		{
			name:   "400 - duplicate seqs",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Duplicate sequence 2 at seqs[3]",
			body: &httpBody{
				Seqs: "1,2,3,2",
			},
		},

		{
			name:   "404 - block seq not found",
			method: http.MethodGet,
			status: http.StatusNotFound,
			err:    "404 Not Found - block does not exist seq=4",
			body: &httpBody{
				Seqs: "1,2,4",
			},
			seqs:                  []uint64{1, 2, 4},
			gatewayGetBlocksError: visor.NewErrBlockNotExist(4),
		},

		{
			name:   "404 - block seq not found verbose",
			method: http.MethodGet,
			status: http.StatusNotFound,
			err:    "404 Not Found - block does not exist seq=4",
			body: &httpBody{
				Seqs:    "1,2,4",
				Verbose: "1",
			},
			seqs:                         []uint64{1, 2, 4},
			verbose:                      true,
			gatewayGetBlocksVerboseError: visor.NewErrBlockNotExist(4),
		},

		{
			name:   "500 - gatewayGetBlocksInRangeError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gatewayGetBlocksInRangeError",
			body: &httpBody{
				Start: "1",
				End:   "3",
			},
			start:                        1,
			end:                          3,
			gatewayGetBlocksInRangeError: errors.New("gatewayGetBlocksInRangeError"),
		},
		{
			name:   "500 - gatewayGetBlocksInRangeVerboseError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gatewayGetBlocksInRangeVerboseError",
			body: &httpBody{
				Start:   "1",
				End:     "3",
				Verbose: "1",
			},
			start:                               1,
			end:                                 3,
			verbose:                             true,
			gatewayGetBlocksInRangeVerboseError: errors.New("gatewayGetBlocksInRangeVerboseError"),
		},

		{
			name:   "500 - gatewayGetBlocksError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gatewayGetBlocksError",
			body: &httpBody{
				Seqs: "1,2,3",
			},
			seqs:                  []uint64{1, 2, 3},
			gatewayGetBlocksError: errors.New("gatewayGetBlocksError"),
		},
		{
			name:   "500 - gatewayGetBlocksVerboseError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gatewayGetBlocksVerboseError",
			body: &httpBody{
				Seqs:    "1,2,3",
				Verbose: "1",
			},
			seqs:                         []uint64{1, 2, 3},
			verbose:                      true,
			gatewayGetBlocksVerboseError: errors.New("gatewayGetBlocksVerboseError"),
		},

		{
			name:   "200 range",
			method: http.MethodGet,
			status: http.StatusOK,
			body: &httpBody{
				Start: "1",
				End:   "3",
			},
			start:                         1,
			end:                           3,
			gatewayGetBlocksInRangeResult: []coin.SignedBlock{{}},
			response: &readable.Blocks{
				Blocks: []readable.Block{
					readable.Block{
						Head: readable.BlockHeader{
							Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
							PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
							BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
							UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
						},
						Body: readable.BlockBody{
							Transactions: []readable.Transaction{},
						},
					},
				},
			},
		},
		{
			name:   "200 range verbose",
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
			gatewayGetBlocksInRangeVerboseResult: verboseResult{
				Blocks: []coin.SignedBlock{{}},
				Inputs: [][][]visor.TransactionInput{{}},
			},
			response: &readable.BlocksVerbose{
				Blocks: []readable.BlockVerbose{
					readable.BlockVerbose{
						Head: readable.BlockHeader{
							Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
							PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
							BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
							UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
						},
						Body: readable.BlockBodyVerbose{
							Transactions: []readable.BlockTransactionVerbose{},
						},
					},
				},
			},
		},

		{
			name:   "200 seqs",
			method: http.MethodGet,
			status: http.StatusOK,
			body: &httpBody{
				Seqs: "1,2,3",
			},
			seqs:                   []uint64{1, 2, 3},
			gatewayGetBlocksResult: []coin.SignedBlock{{}},
			response: &readable.Blocks{
				Blocks: []readable.Block{
					readable.Block{
						Head: readable.BlockHeader{
							Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
							PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
							BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
							UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
						},
						Body: readable.BlockBody{
							Transactions: []readable.Transaction{},
						},
					},
				},
			},
		},
		{
			name:   "200 seqs verbose",
			method: http.MethodGet,
			status: http.StatusOK,
			body: &httpBody{
				Seqs:    "1,2,3",
				Verbose: "1",
			},
			seqs:    []uint64{1, 2, 3},
			verbose: true,
			gatewayGetBlocksVerboseResult: verboseResult{
				Blocks: []coin.SignedBlock{{}},
				Inputs: [][][]visor.TransactionInput{{}},
			},
			response: &readable.BlocksVerbose{
				Blocks: []readable.BlockVerbose{
					readable.BlockVerbose{
						Head: readable.BlockHeader{
							Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
							PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
							BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
							UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
						},
						Body: readable.BlockBodyVerbose{
							Transactions: []readable.BlockTransactionVerbose{},
						},
					},
				},
			},
		},

		{
			name:   "200 seqs POST",
			method: http.MethodPost,
			status: http.StatusOK,
			body: &httpBody{
				Seqs: "1,2,3",
			},
			seqs:                   []uint64{1, 2, 3},
			gatewayGetBlocksResult: []coin.SignedBlock{{}},
			response: &readable.Blocks{
				Blocks: []readable.Block{
					readable.Block{
						Head: readable.BlockHeader{
							Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
							PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
							BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
							UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
						},
						Body: readable.BlockBody{
							Transactions: []readable.Transaction{},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("GetBlocksInRange", tc.start, tc.end).Return(tc.gatewayGetBlocksInRangeResult, tc.gatewayGetBlocksInRangeError)
			gateway.On("GetBlocksInRangeVerbose", tc.start, tc.end).Return(tc.gatewayGetBlocksInRangeVerboseResult.Blocks,
				tc.gatewayGetBlocksInRangeVerboseResult.Inputs, tc.gatewayGetBlocksInRangeVerboseError)
			gateway.On("GetBlocks", tc.seqs).Return(tc.gatewayGetBlocksResult, tc.gatewayGetBlocksError)
			gateway.On("GetBlocksVerbose", tc.seqs).Return(tc.gatewayGetBlocksVerboseResult.Blocks,
				tc.gatewayGetBlocksVerboseResult.Inputs, tc.gatewayGetBlocksVerboseError)

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
				if tc.body.Seqs != "" {
					v.Add("seqs", tc.body.Seqs)
				}
			}

			var reqBody io.Reader
			if len(v) > 0 {
				if tc.method == http.MethodPost {
					reqBody = strings.NewReader(v.Encode())
				} else {
					endpoint += "?" + v.Encode()
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, reqBody)
			require.NoError(t, err)

			if tc.method == http.MethodPost {
				req.Header.Set("Content-Type", ContentTypeForm)
			}

			setCSRFParameters(t, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				if tc.verbose {
					var msg *readable.BlocksVerbose
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.response, msg)
				} else {
					var msg *readable.Blocks
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

	type verboseResult struct {
		Blocks []coin.SignedBlock
		Inputs [][][]visor.TransactionInput
	}

	tt := []struct {
		name                              string
		method                            string
		status                            int
		err                               string
		body                              httpBody
		num                               uint64
		verbose                           bool
		gatewayGetLastBlocksResult        []coin.SignedBlock
		gatewayGetLastBlocksError         error
		gatewayGetLastBlocksVerboseResult verboseResult
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
			num:                        1,
			gatewayGetLastBlocksResult: []coin.SignedBlock{{}},
			response: &readable.Blocks{
				Blocks: []readable.Block{
					readable.Block{
						Head: readable.BlockHeader{
							Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
							PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
							BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
							UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
						},
						Body: readable.BlockBody{
							Transactions: []readable.Transaction{},
						},
					},
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
			gatewayGetLastBlocksVerboseResult: verboseResult{
				Blocks: []coin.SignedBlock{{}},
				Inputs: [][][]visor.TransactionInput{{}},
			},
			response: &readable.BlocksVerbose{
				Blocks: []readable.BlockVerbose{
					readable.BlockVerbose{
						Head: readable.BlockHeader{
							Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
							PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
							BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
							UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
						},
						Body: readable.BlockBodyVerbose{
							Transactions: []readable.BlockTransactionVerbose{},
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/last_blocks"
			gateway := &MockGatewayer{}

			gateway.On("DaemonConfig").Return(daemon.DaemonConfig{MaxLastBlocksCount: 256})
			gateway.On("GetLastBlocks", tc.num).Return(tc.gatewayGetLastBlocksResult, tc.gatewayGetLastBlocksError)
			gateway.On("GetLastBlocksVerbose", tc.num).Return(tc.gatewayGetLastBlocksVerboseResult.Blocks,
				tc.gatewayGetLastBlocksVerboseResult.Inputs, tc.gatewayGetLastBlocksVerboseError)

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

			setCSRFParameters(t, tokenValid, req)

			rr := httptest.NewRecorder()

			handler := newServerMux(defaultMuxConfig(), gateway)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				if tc.verbose {
					var msg *readable.BlocksVerbose
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.response, msg)
				} else {
					var msg *readable.Blocks
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.response, msg)
				}
			}
		})
	}
}
