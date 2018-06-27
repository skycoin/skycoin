package api

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
		name         string
		method       string
		code         int
		getHealthErr error
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
			name:         "gateway.GetHealth error",
			method:       http.MethodGet,
			code:         http.StatusInternalServerError,
			getHealthErr: errors.New("GetHealth failed"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			unspents := uint64(10)
			unconfirmed := uint64(20)

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

			health := &daemon.Health{
				BlockchainMetadata: metadata,
				OpenConnections:    3,
				Version:            buildInfo,
				Uptime:             time.Second * 4,
			}

			gateway := NewGatewayerMock()
			if tc.getHealthErr != nil {
				gateway.On("GetHealth").Return(nil, tc.getHealthErr)
			} else {
				gateway.On("GetHealth").Return(health, nil)
			}

			endpoint := "/api/v1/health"
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			cfg := muxConfig{
				host:   configuredHost,
				appLoc: ".",
			}
			handler := newServerMux(cfg, gateway, &CSRFStore{}, nil)
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
			require.Equal(t, health.Uptime, r.Uptime.Duration)

			require.Equal(t, health.OpenConnections, r.OpenConnections)

			require.Equal(t, unconfirmed, r.BlockchainMetadata.Unconfirmed)
			require.Equal(t, unspents, r.BlockchainMetadata.Unspents)
			require.True(t, r.BlockchainMetadata.TimeSinceLastBlock.Duration > time.Duration(0))
			require.Equal(t, metadata.Head, r.BlockchainMetadata.Head)
		})
	}
}
