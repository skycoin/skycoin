package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/skycoin/skycoin/src/params"
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
	OutgoingConnections   int                `json:"outgoing_connections"`
	IncomingConnections   int                `json:"incoming_connections"`
	Uptime                wh.Duration        `json:"uptime"`
	CSRFEnabled           bool               `json:"csrf_enabled"`
	CSPEnabled            bool               `json:"csp_enabled"`
	WalletAPIEnabled      bool               `json:"wallet_api_enabled"`
	GUIEnabled            bool               `json:"gui_enabled"`
	UnversionedAPIEnabled bool               `json:"unversioned_api_enabled"`
	JSON20RPCEnabled      bool               `json:"json_rpc_enabled"`
	UserVerifyTxn         readable.VerifyTxn `json:"user_verify_transaction"`
	UnconfirmedVerifyTxn  readable.VerifyTxn `json:"unconfirmed_verify_transaction"`
	APIStartedAt          int64              `json:api_started_at`
	CoinStartedAt         int64              `json:coin_started_at`
	DaemonStartedAt       int64              `json:daemon_started_at`
}

// healthHandler returns node health data
// URI: /api/v1/health
// Method: GET
func healthHandler(c muxConfig, csrfStore *CSRFStore, gateway Gatewayer) http.HandlerFunc {
	getAPIStartTime := c.health.GetAPIStartTime
	getCoinStartTime := c.health.GetCoinStartTime
	getDaemonStartTime := c.health.GetDaemonStartTime

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

		apiStartUnixTime := int64(0)
		if getAPIStartTime != nil {
			apiStartUnixTime = getAPIStartTime().Unix()
		}
		coinStartUnixTime := int64(0)
		if getCoinStartTime != nil {
			coinStartUnixTime = getCoinStartTime().Unix()
		}
		daemonStartUnixTime := int64(0)
		if getDaemonStartTime != nil {
			daemonStartUnixTime = getDaemonStartTime().Unix()
		}

		wh.SendJSONOr500(logger, w, HealthResponse{
			BlockchainMetadata: BlockchainMetadata{
				BlockchainMetadata: readable.NewBlockchainMetadata(health.BlockchainMetadata),
				TimeSinceLastBlock: wh.FromDuration(timeSinceLastBlock),
			},
			Version:               c.health.BuildInfo,
			CoinName:              c.health.CoinName,
			DaemonUserAgent:       userAgent,
			OpenConnections:       health.OutgoingConnections + health.IncomingConnections,
			OutgoingConnections:   health.OutgoingConnections,
			IncomingConnections:   health.IncomingConnections,
			Uptime:                wh.FromDuration(health.Uptime),
			CSRFEnabled:           csrfStore.Enabled,
			CSPEnabled:            !c.disableCSP,
			UnversionedAPIEnabled: c.enableUnversionedAPI,
			GUIEnabled:            c.enableGUI,
			JSON20RPCEnabled:      c.enableJSON20RPC,
			WalletAPIEnabled:      walletAPIEnabled,
			UserVerifyTxn:         readable.NewVerifyTxn(params.UserVerifyTxn),
			UnconfirmedVerifyTxn:  readable.NewVerifyTxn(health.UnconfirmedVerifyTxn),
			APIStartedAt:          apiStartUnixTime,
			CoinStartedAt:         coinStartUnixTime,
			DaemonStartedAt:       daemonStartUnixTime,
		})
	}
}
