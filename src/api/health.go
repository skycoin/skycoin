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
<<<<<<< 8ce8115a116d9f110da8e874a7af84aca85a148a
	BlockchainMetadata   BlockchainMetadata `json:"blockchain"`
	Version              readable.BuildInfo `json:"version"`
	CoinName             string             `json:"coin"`
	DaemonUserAgent      string             `json:"user_agent"`
	OpenConnections      int                `json:"open_connections"`
	OutgoingConnections  int                `json:"outgoing_connections"`
	IncomingConnections  int                `json:"incoming_connections"`
	Uptime               wh.Duration        `json:"uptime"`
	CSRFEnabled          bool               `json:"csrf_enabled"`
	CSPEnabled           bool               `json:"csp_enabled"`
	WalletAPIEnabled     bool               `json:"wallet_api_enabled"`
	GUIEnabled           bool               `json:"gui_enabled"`
	UserVerifyTxn        readable.VerifyTxn `json:"user_verify_transaction"`
	UnconfirmedVerifyTxn readable.VerifyTxn `json:"unconfirmed_verify_transaction"`
	StartedAt            int64              `json:"started_at"`
=======
	BlockchainMetadata    BlockchainMetadata `json:"blockchain"`
	Version               readable.BuildInfo `json:"version"`
	CoinName              string             `json:"coin"`
	DaemonUserAgent       string             `json:"user_agent"`
	OpenConnections       int                `json:"open_connections"`
	OutgoingConnections   int                `json:"outgoing_connections"`
	IncomingConnections   int                `json:"incoming_connections"`
	Uptime                wh.Duration        `json:"uptime"`
	CSRFEnabled           bool               `json:"csrf_enabled"`
	HeadercheckEnabled    bool               `json:"headercheck_enabled"`
	CSPEnabled            bool               `json:"csp_enabled"`
	WalletAPIEnabled      bool               `json:"wallet_api_enabled"`
	GUIEnabled            bool               `json:"gui_enabled"`
	UnversionedAPIEnabled bool               `json:"unversioned_api_enabled"`
	JSON20RPCEnabled      bool               `json:"json_rpc_enabled"`
	UserVerifyTxn         readable.VerifyTxn `json:"user_verify_transaction"`
	UnconfirmedVerifyTxn  readable.VerifyTxn `json:"unconfirmed_verify_transaction"`
	StartedAt             int64              `json:"started_at"`
>>>>>>> implement the correct check
}

// healthHandler returns node health data
// URI: /api/v1/health
// Method: GET
func healthHandler(c muxConfig, gateway Gatewayer) http.HandlerFunc {
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
<<<<<<< 8ce8115a116d9f110da8e874a7af84aca85a148a
			Version:              c.health.BuildInfo,
			CoinName:             c.health.CoinName,
			DaemonUserAgent:      userAgent,
			OpenConnections:      health.OutgoingConnections + health.IncomingConnections,
			OutgoingConnections:  health.OutgoingConnections,
			IncomingConnections:  health.IncomingConnections,
			Uptime:               wh.FromDuration(health.Uptime),
			CSRFEnabled:          !c.disableCSRF,
			CSPEnabled:           !c.disableCSP,
			GUIEnabled:           c.enableGUI,
			WalletAPIEnabled:     walletAPIEnabled,
			UserVerifyTxn:        readable.NewVerifyTxn(params.UserVerifyTxn),
			UnconfirmedVerifyTxn: readable.NewVerifyTxn(health.UnconfirmedVerifyTxn),
			StartedAt:            health.StartedAt.Unix(),
=======
			Version:               c.health.BuildInfo,
			CoinName:              c.health.CoinName,
			DaemonUserAgent:       userAgent,
			OpenConnections:       health.OutgoingConnections + health.IncomingConnections,
			OutgoingConnections:   health.OutgoingConnections,
			IncomingConnections:   health.IncomingConnections,
			Uptime:                wh.FromDuration(health.Uptime),
			CSRFEnabled:           !c.disableCSRF,
			HeadercheckEnabled:    !c.disableHeadercheck,
			CSPEnabled:            !c.disableCSP,
			UnversionedAPIEnabled: c.enableUnversionedAPI,
			GUIEnabled:            c.enableGUI,
			JSON20RPCEnabled:      c.enableJSON20RPC,
			WalletAPIEnabled:      walletAPIEnabled,
			UserVerifyTxn:         readable.NewVerifyTxn(params.UserVerifyTxn),
			UnconfirmedVerifyTxn:  readable.NewVerifyTxn(health.UnconfirmedVerifyTxn),
			StartedAt:             health.StartedAt.Unix(),
>>>>>>> implement the correct check
		})
	}
}
