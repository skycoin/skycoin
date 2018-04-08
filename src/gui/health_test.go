package gui

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"encoding/json"
	"net/http/httptest"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
)

func TestHealthCheckHandler(t *testing.T) {

	cases := []struct {
		name                     string
		method                   string
		code                     int
		getBlockchainMetadataErr error
	}{
		{
			name:   "valid response",
			method: http.MethodGet,
			code:   http.StatusOK,
		},
		{
			name:   "403 method not allowed",
			method: http.MethodPost,
			code:   http.StatusMethodNotAllowed,
		},
		{
			name:   "gateway.GetBlockchainMetadata error",
			method: http.MethodGet,
			code:   http.StatusInternalServerError,
			getBlockchainMetadataErr: errors.New("GetBlockchainMetadata failed"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			unspents := uint64(10)
			unconfirmed := uint64(20)

			connections := &daemon.Connections{
				Connections: []*daemon.Connection{
					{
						ID:   1,
						Addr: "127.0.0.1:4343",
					},
					{
						ID:   2,
						Addr: "127.0.0.1:5454",
					},
					{
						ID:   3,
						Addr: "127.0.0.1:6565",
					},
				},
			}

			metadata := &visor.BlockchainMetadata{
				Head: visor.ReadableBlockHeader{
					BkSeq:             21175,
					BlockHash:         "8a3e0aac619551ae009cfb28c2b36bb1300925f74da770d1512072314f6a4c80",
					PreviousBlockHash: "001eb7911b6a6ab7c75feb88726dd2bc8b87133aebc82201c4404537eb74f7ac",
					Time:              1523168686,
					Fee:               2,
					Version:           0,
					BodyHash:          "36be8d70d1e9f70b340ea7ecf0b247c27086bad10568044c1196fe150f6cea1b",
				},
				Unspents:    unspents,
				Unconfirmed: unconfirmed,
			}

			buildInfo := visor.BuildInfo{
				Version: "1.0.0",
				Commit:  "abcdef",
				Branch:  "develop",
			}

			gateway := NewGatewayerMock()
			gateway.On("GetConnections").Return(connections)
			gateway.On("GetBuildInfo").Return(buildInfo)
			if tc.getBlockchainMetadataErr != nil {
				gateway.On("GetBlockchainMetadata").Return(nil, tc.getBlockchainMetadataErr)
			} else {
				gateway.On("GetBlockchainMetadata").Return(metadata, nil)
			}

			endpoint := "/health"
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			cfg := muxConfig{
				host:   configuredHost,
				appLoc: ".",
			}
			handler := newServerMux(cfg, gateway, &CSRFStore{})
			handler.ServeHTTP(rr, req)

			if tc.code != http.StatusOK {
				require.Equal(t, tc.code, rr.Code)
				return
			}

			require.Equal(t, http.StatusOK, rr.Code)

			r := &HealthResponse{}
			err = json.Unmarshal(rr.Body.Bytes(), r)
			require.NoError(t, err)

			require.Equal(t, buildInfo.Version, r.Version.Version)
			require.Equal(t, buildInfo.Commit, r.Version.Commit)
			require.Equal(t, buildInfo.Branch, r.Version.Branch)

			require.Equal(t, len(connections.Connections), r.OpenConnections)

			require.Equal(t, unconfirmed, r.Blockchain.Unconfirmed)
			require.Equal(t, unspents, r.Blockchain.Unspents)
			require.True(t, r.Blockchain.TimeSinceLastBlock.Duration > time.Duration(0))
			require.Equal(t, metadata.Head, r.Blockchain.Head)
		})
	}
}
