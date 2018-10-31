package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
)

// BlockchainMetadata extends visor.BlockchainMetadata to include the time since the last block
type BlockchainMetadata struct {
	readable.BlockchainMetadata
	TimeSinceLastBlock wh.Duration `json:"time_since_last_block"`
}

// HealthResponse is returned by the /health endpoint
type HealthResponse struct {
	BlockchainMetadata    BlockchainMetadata `json:"blockchain"`
	Version               readable.BuildInfo `json:"version"`
	CoinName              string             `json:"coin"`
	DaemonUserAgent       string             `json:"user_agent"`
	OpenConnections       int                `json:"open_connections"`
	Uptime                wh.Duration        `json:"uptime"`
	CSRFEnabled           bool               `json:"csrf_enabled"`
	CSPEnabled            bool               `json:"csp_enabled"`
	WalletAPIEnabled      bool               `json:"wallet_api_enabled"`
	GUIEnabled            bool               `json:"gui_enabled"`
	UnversionedAPIEnabled bool               `json:"unversioned_api_enabled"`
	JSON20RPCEnabled      bool               `json:"json_rpc_enabled"`
}

// healthHandler returns node health data
// URI: /api/v1/health
// Method: GET
func healthHandler(c muxConfig, csrfStore *CSRFStore, gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		health, err := gateway.GetHealth()
		if err != nil {
			err = fmt.Errorf("gateway.GetHealth failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		elapsedBlockTime := time.Now().UTC().Unix() - int64(health.BlockchainMetadata.HeadBlock.Head.Time)
		timeSinceLastBlock := time.Second * time.Duration(elapsedBlockTime)

		_, walletAPIEnabled := c.enabledAPISets[EndpointsWallet]

		userAgent, err := c.health.DaemonUserAgent.Build()
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, HealthResponse{
			BlockchainMetadata: BlockchainMetadata{
				BlockchainMetadata: readable.NewBlockchainMetadata(health.BlockchainMetadata),
				TimeSinceLastBlock: wh.FromDuration(timeSinceLastBlock),
			},
			Version:               c.health.BuildInfo,
			CoinName:              c.health.CoinName,
			DaemonUserAgent:       userAgent,
			OpenConnections:       health.OpenConnections,
			Uptime:                wh.FromDuration(health.Uptime),
			CSRFEnabled:           csrfStore.Enabled,
			CSPEnabled:            !c.disableCSP,
			UnversionedAPIEnabled: c.enableUnversionedAPI,
			GUIEnabled:            c.enableGUI,
			JSON20RPCEnabled:      c.enableJSON20RPC,
			WalletAPIEnabled:      walletAPIEnabled,
		})
	}
}
