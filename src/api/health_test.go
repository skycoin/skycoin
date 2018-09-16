package api

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"encoding/json"
	"net/http/httptest"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/visor"
)

func TestHealthHandler(t *testing.T) {
	cases := []struct {
		name             string
		method           string
		code             int
		err              string
		getHealthErr     error
		cfg              muxConfig
		csrfEnabled      bool
		walletAPIEnabled bool
	}{
		{
			name:   "405 method not allowed",
			method: http.MethodPost,
			code:   http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
			cfg:    defaultMuxConfig(),
		},

		{
			name:         "gateway.GetHealth error",
			method:       http.MethodGet,
			code:         http.StatusInternalServerError,
			err:          "500 Internal Server Error - gateway.GetHealth failed: GetHealth failed",
			getHealthErr: errors.New("GetHealth failed"),
			cfg:          defaultMuxConfig(),
		},

		{
			name:             "valid response",
			method:           http.MethodGet,
			code:             http.StatusOK,
			cfg:              defaultMuxConfig(),
			walletAPIEnabled: true,
		},

		{
			name:   "valid response, opposite config",
			method: http.MethodGet,
			code:   http.StatusOK,
			cfg: muxConfig{
				host:                 configuredHost,
				appLoc:               ".",
				disableCSP:           false,
				enableGUI:            true,
				enableUnversionedAPI: true,
				enableJSON20RPC:      true,
				enabledAPISets: map[string]struct{}{
					EndpointsStatus: struct{}{},
					EndpointsRead:   struct{}{},
				},
			},
			csrfEnabled:      true,
			walletAPIEnabled: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			unspents := uint64(10)
			unconfirmed := uint64(20)

			uxHash, err := cipher.SHA256FromHex("8a3e0aac619551ae009cfb28c2b36bb1300925f74da770d1512072314f6a4c80")
			require.NoError(t, err)

			prevHash, err := cipher.SHA256FromHex("001eb7911b6a6ab7c75feb88726dd2bc8b87133aebc82201c4404537eb74f7ac")
			require.NoError(t, err)

			bodyHash, err := cipher.SHA256FromHex("36be8d70d1e9f70b340ea7ecf0b247c27086bad10568044c1196fe150f6cea1b")
			require.NoError(t, err)

			metadata := visor.BlockchainMetadata{
				HeadBlock: coin.SignedBlock{
					Block: coin.Block{
						Head: coin.BlockHeader{
							BkSeq:    21175,
							UxHash:   uxHash,
							PrevHash: prevHash,
							Time:     1523168686,
							Fee:      2,
							Version:  0,
							BodyHash: bodyHash,
						},
					},
				},
				Unspents:    unspents,
				Unconfirmed: unconfirmed,
			}

			buildInfo := readable.BuildInfo{
				Version: "1.0.0",
				Commit:  "abcdef",
				Branch:  "develop",
			}
			tc.cfg.buildInfo = buildInfo

			health := &daemon.Health{
				BlockchainMetadata: metadata,
				OpenConnections:    3,
				Uptime:             time.Second * 4,
			}

			gateway := &MockGatewayer{}

			if tc.getHealthErr != nil {
				gateway.On("GetHealth").Return(nil, tc.getHealthErr)
			} else {
				gateway.On("GetHealth").Return(health, nil)
			}

			endpoint := "/api/v1/health"
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: tc.csrfEnabled,
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(tc.cfg, gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)
			if tc.code != http.StatusOK {
				require.Equal(t, tc.code, rr.Code)
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()))
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

			require.Equal(t, metadata.HeadBlock.Block.Head.BkSeq, r.BlockchainMetadata.Head.BkSeq)
			require.Equal(t, metadata.HeadBlock.Block.Head.Time, r.BlockchainMetadata.Head.Time)
			require.Equal(t, metadata.HeadBlock.Block.Head.Fee, r.BlockchainMetadata.Head.Fee)
			require.Equal(t, metadata.HeadBlock.Block.Head.Version, r.BlockchainMetadata.Head.Version)
			require.Equal(t, metadata.HeadBlock.Block.Head.PrevHash.Hex(), r.BlockchainMetadata.Head.PreviousBlockHash)
			require.Equal(t, metadata.HeadBlock.Block.Head.Hash().Hex(), r.BlockchainMetadata.Head.BlockHash)
			require.Equal(t, metadata.HeadBlock.Block.Head.BodyHash.Hex(), r.BlockchainMetadata.Head.BodyHash)

			require.Equal(t, tc.csrfEnabled, r.CSRFEnabled)
			require.Equal(t, !tc.cfg.disableCSP, r.CSPEnabled)
			require.Equal(t, tc.cfg.enableUnversionedAPI, r.UnversionedAPIEnabled)
			require.Equal(t, tc.cfg.enableGUI, r.GUIEnabled)
			require.Equal(t, tc.cfg.enableJSON20RPC, r.JSON20RPCEnabled)
			require.Equal(t, tc.walletAPIEnabled, r.WalletAPIEnabled)
		})
	}
}
