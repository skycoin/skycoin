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
// swagger:model blockchainMetadata_extend
type BlockchainMetadata struct {
	// swagger:allOf
	readable.BlockchainMetadata
	TimeSinceLastBlock wh.Duration `json:"time_since_last_block"`
}

// HealthResponse is returned by the /health endpoint
// swagger:response healthResponse
type HealthResponse struct {

	// require: true
	BlockchainMetadata BlockchainMetadata `json:"blockchain"`
	Version             readable.BuildInfo `json:"version"`
	CoinName            string             `json:"coin"`
	DaemonUserAgent     string             `json:"user_agent"`
	OpenConnections     int                `json:"open_connections"`
	OutgoingConnections int                `json:"outgoing_connections"`
	IncomingConnections int                `json:"incoming_connections"`
	// swagger:strfmt duration
	Uptime                wh.Duration `json:"uptime"`
	CSRFEnabled           bool        `json:"csrf_enabled"`
	CSPEnabled            bool        `json:"csp_enabled"`
	WalletAPIEnabled      bool        `json:"wallet_api_enabled"`
	GUIEnabled            bool        `json:"gui_enabled"`
	UnversionedAPIEnabled bool        `json:"unversioned_api_enabled"`
	JSON20RPCEnabled      bool        `json:"json_rpc_enabled"`
	// swagger:ignore
	UserVerifyTxn readable.VerifyTxn `json:"user_verify_transaction"`
	// swagger:ignore
	UnconfirmedVerifyTxn readable.VerifyTxn `json:"unconfirmed_verify_transaction"`
	StartedAt            int64              `json:"started_at"`
}

// swagger:response realHealthResponse
type RealHealthResponse struct {
	// in: body
	BlockchainMetadata struct {
		// swagger:allOf
		readable.BlockchainMetadata
		TimeSinceLastBlock wh.Duration `json:"time_since_last_block"`
	}`json:"blockchain"`
	CoinName            string             `json:"coin"`
	DaemonUserAgent     string             `json:"user_agent"`
	OpenConnections     int                `json:"open_connections"`
	OutgoingConnections int                `json:"outgoing_connections"`
	IncomingConnections int                `json:"incoming_connections"`
	// swagger:strfmt duration
	Uptime                wh.Duration `json:"uptime"`
	CSRFEnabled           bool        `json:"csrf_enabled"`
	CSPEnabled            bool        `json:"csp_enabled"`
	WalletAPIEnabled      bool        `json:"wallet_api_enabled"`
	GUIEnabled            bool        `json:"gui_enabled"`
	UnversionedAPIEnabled bool        `json:"unversioned_api_enabled"`
	JSON20RPCEnabled      bool        `json:"json_rpc_enabled"`
	// in: body
	UserVerifyTxn readable.VerifyTxn `json:"user_verify_transaction"`
	// in: body
	UnconfirmedVerifyTxn readable.VerifyTxn `json:"unconfirmed_verify_transaction"`
	StartedAt            int64              `json:"started_at"`
	// in: body
	Version             readable.BuildInfo `json:"version"`
}


// healthHandler returns node health data
// URI: /api/v1/health
// Method: GET
func healthHandler(c muxConfig, gateway Gatewayer) http.HandlerFunc {

	// swagger:route GET /api/v1/health health
	//
	// healthHandler returns node health data
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http, https
	//
	//     Security:
	//       api_key:
	//       oauth: read, write
	//
	//     Responses:
	//       default: genericError
	//       200: realHealthResponse
	//       422: validationError

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
			OpenConnections:       health.OutgoingConnections + health.IncomingConnections,
			OutgoingConnections:   health.OutgoingConnections,
			IncomingConnections:   health.IncomingConnections,
			Uptime:                wh.FromDuration(health.Uptime),
			CSRFEnabled:           !c.disableCSRF,
			CSPEnabled:            !c.disableCSP,
			UnversionedAPIEnabled: c.enableUnversionedAPI,
			GUIEnabled:            c.enableGUI,
			JSON20RPCEnabled:      c.enableJSON20RPC,
			WalletAPIEnabled:      walletAPIEnabled,
			UserVerifyTxn:         readable.NewVerifyTxn(params.UserVerifyTxn),
			UnconfirmedVerifyTxn:  readable.NewVerifyTxn(health.UnconfirmedVerifyTxn),
			StartedAt:             health.StartedAt.Unix(),
		})
	}
}
