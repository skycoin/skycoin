package api

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"encoding/json"
	"net/http/httptest"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/useragent"
	"github.com/skycoin/skycoin/src/visor"
)

func TestMain(m *testing.M) {
	params.UserVerifyTxn.BurnFactor = 10
}

func TestHealthHandler(t *testing.T) {
	cases := []struct {
		name                     string
		method                   string
		code                     int
		err                      string
		getBlockchainMetadataErr error
		getConnectionsErr        error
		cfg                      muxConfig
		walletAPIEnabled         bool
	}{
		{
			name:   "405 method not allowed",
			method: http.MethodPost,
			code:   http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
			cfg:    defaultMuxConfig(),
		},

		{
			name:                     "gateway.GetBlockchainMetadata error",
			method:                   http.MethodGet,
			code:                     http.StatusInternalServerError,
			err:                      "500 Internal Server Error - gateway.GetBlockchainMetadata failed: GetBlockchainMetadata failed",
			getBlockchainMetadataErr: errors.New("GetBlockchainMetadata failed"),
			cfg:                      defaultMuxConfig(),
		},

		{
			name:              "gateway.GetConnections error",
			method:            http.MethodGet,
			code:              http.StatusInternalServerError,
			err:               "500 Internal Server Error - gateway.GetConnections failed: GetConnections failed",
			getConnectionsErr: errors.New("GetConnections failed"),
			cfg:               defaultMuxConfig(),
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
				health: HealthConfig{
					BlockPublisher: true,
				},
				host:        configuredHost,
				appLoc:      ".",
				disableCSRF: false,
				disableCSP:  false,
				enableGUI:   true,
				enabledAPISets: map[string]struct{}{
					EndpointsStatus: struct{}{},
					EndpointsRead:   struct{}{},
				},
			},
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
			tc.cfg.health.BuildInfo = buildInfo

			tc.cfg.health.Fiber.Name = "skycoin"
			tc.cfg.health.DaemonUserAgent = useragent.Data{
				Coin:    "skycoin",
				Version: "0.25.0",
				Remark:  "test",
			}

			conns := []daemon.Connection{
				{
					ConnectionDetails: daemon.ConnectionDetails{
						Outgoing: true,
						State:    daemon.ConnectionStateConnected,
					},
				},
				{
					ConnectionDetails: daemon.ConnectionDetails{
						Outgoing: false,
						State:    daemon.ConnectionStateIntroduced,
					},
				},
				{
					ConnectionDetails: daemon.ConnectionDetails{
						Outgoing: true,
						State:    daemon.ConnectionStateIntroduced,
					},
				},
				{
					ConnectionDetails: daemon.ConnectionDetails{
						Outgoing: false,
						State:    daemon.ConnectionStateConnected,
					},
				},
				{
					ConnectionDetails: daemon.ConnectionDetails{
						Outgoing: true,
						State:    daemon.ConnectionStateConnected,
					},
				},
			}

			gateway := &MockGatewayer{}

			if tc.getBlockchainMetadataErr != nil {
				gateway.On("GetBlockchainMetadata").Return(nil, tc.getBlockchainMetadataErr)
			} else {
				gateway.On("GetBlockchainMetadata").Return(&metadata, nil)
			}

			if tc.getConnectionsErr != nil {
				gateway.On("GetConnections", mock.Anything).Return(nil, tc.getConnectionsErr)
			} else {
				gateway.On("GetConnections", mock.Anything).Return(conns, nil)
			}

			startedAt := time.Now().Add(time.Second * -4)

			gateway.On("StartedAt").Return(startedAt)

			dc := daemon.DaemonConfig{
				UnconfirmedVerifyTxn: params.VerifyTxn{
					BurnFactor:          params.UserVerifyTxn.BurnFactor * 2,
					MaxTransactionSize:  params.UserVerifyTxn.MaxTransactionSize * 2,
					MaxDropletPrecision: params.UserVerifyTxn.MaxDropletPrecision - 1,
				},
			}

			gateway.On("DaemonConfig").Return(dc)

			endpoint := "/api/v1/health"
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(tc.cfg, gateway)
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
			// Put uptime in a range in case the test scheduler is slow and there is a pause.
			// Should be at least 4 seconds since startedAt was computed to be 4 seconds ago.
			require.True(t, time.Second*4 <= r.Uptime.Duration && time.Second*8 >= r.Uptime.Duration, "uptime should be a little over 4 seconds")

			require.Equal(t, 3, r.OutgoingConnections)
			require.Equal(t, 2, r.IncomingConnections)
			require.Equal(t, len(conns), r.OpenConnections)
			require.Equal(t, "skycoin", r.CoinName)
			require.Equal(t, "skycoin:0.25.0(test)", r.DaemonUserAgent)
			require.Equal(t, tc.cfg.health.BlockPublisher, r.BlockPublisher)

			require.Equal(t, unconfirmed, r.BlockchainMetadata.Unconfirmed)
			require.Equal(t, unspents, r.BlockchainMetadata.Unspents)
			require.True(t, r.BlockchainMetadata.TimeSinceLastBlock.Duration > time.Duration(0))

			require.Equal(t, metadata.HeadBlock.Block.Head.BkSeq, r.BlockchainMetadata.Head.BkSeq)
			require.Equal(t, metadata.HeadBlock.Block.Head.Time, r.BlockchainMetadata.Head.Time)
			require.Equal(t, metadata.HeadBlock.Block.Head.Fee, r.BlockchainMetadata.Head.Fee)
			require.Equal(t, metadata.HeadBlock.Block.Head.Version, r.BlockchainMetadata.Head.Version)
			require.Equal(t, metadata.HeadBlock.Block.Head.PrevHash.Hex(), r.BlockchainMetadata.Head.PreviousHash)
			require.Equal(t, metadata.HeadBlock.Block.Head.Hash().Hex(), r.BlockchainMetadata.Head.Hash)
			require.Equal(t, metadata.HeadBlock.Block.Head.BodyHash.Hex(), r.BlockchainMetadata.Head.BodyHash)

			require.Equal(t, !tc.cfg.disableCSRF, r.CSRFEnabled)
			require.Equal(t, !tc.cfg.disableCSP, r.CSPEnabled)
			require.Equal(t, tc.cfg.enableGUI, r.GUIEnabled)
			require.Equal(t, tc.walletAPIEnabled, r.WalletAPIEnabled)

			require.Equal(t, uint32(10), r.UserVerifyTxn.BurnFactor)
			require.Equal(t, uint32(32*1024), r.UserVerifyTxn.MaxTransactionSize)
			require.Equal(t, uint8(3), r.UserVerifyTxn.MaxDropletPrecision)

			require.Equal(t, dc.UnconfirmedVerifyTxn.BurnFactor, r.UnconfirmedVerifyTxn.BurnFactor)
			require.Equal(t, dc.UnconfirmedVerifyTxn.MaxTransactionSize, r.UnconfirmedVerifyTxn.MaxTransactionSize)
			require.Equal(t, dc.UnconfirmedVerifyTxn.MaxDropletPrecision, r.UnconfirmedVerifyTxn.MaxDropletPrecision)
			require.True(t, time.Now().Unix() > r.StartedAt)

		})
	}
}
