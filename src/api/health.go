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
	BlockchainMetadata   BlockchainMetadata `json:"blockchain"`
	Version              readable.BuildInfo `json:"version"`
	CoinName             string             `json:"coin"`
	DaemonUserAgent      string             `json:"user_agent"`
	OpenConnections      int                `json:"open_connections"`
	OutgoingConnections  int                `json:"outgoing_connections"`
	IncomingConnections  int                `json:"incoming_connections"`
	Uptime               wh.Duration        `json:"uptime"`
	CSRFEnabled          bool               `json:"csrf_enabled"`
	HeaderCheckEnabled   bool               `json:"header_check_enabled"`
	CSPEnabled           bool               `json:"csp_enabled"`
	WalletAPIEnabled     bool               `json:"wallet_api_enabled"`
	GUIEnabled           bool               `json:"gui_enabled"`
	UserVerifyTxn        readable.VerifyTxn `json:"user_verify_transaction"`
	UnconfirmedVerifyTxn readable.VerifyTxn `json:"unconfirmed_verify_transaction"`
	StartedAt            int64              `json:"started_at"`
}

// healthHandler returns node health data
// URI: /api/v1/health
// Method: GET
func healthHandler(c muxConfig, gateway Gatewayer) http.HandlerFunc {

	// swagger:operation GET /api/v1/health health
	//
	// Returns node health data.
	//
	// ---
	//
	// produces:
	// - application/json
	// responses:
	//   200:
	//     description: This endpoint returns node health data.
	//     schema:
	//       properties:
	//         blockchain:
	//           type: object
	//           properties:
	//             head:
	//               type: object
	//               properties:
	//                 seq:
	//                   type: integer
	//                   format: int64
	//                 block_hash:
	//                   type: string
	//                 previous_block_hash:
	//                   type: string
	//                 timestamp:
	//                   type: integer
	//                   format: int64
	//                 fee:
	//                   type: integer
	//                   format: int64
	//                 version:
	//                   type: integer
	//                   format: int64
	//                 tx_body_hash:
	//                   type: string
	//                 ux_hash:
	//                   type: string
	//             unspents:
	//               type: integer
	//               format: int64
	//             unconfirmed:
	//               type: integer
	//               format: int64
	//             time_since_last_block:
	//               type: string
	//         version:
	//           type: object
	//           properties:
	//             version:
	//               type: string
	//             commit:
	//               type: string
	//             branch:
	//               type: string
	//         coin:
	//           type: string
	//         user_agent:
	//           type: string
	//         open_connections:
	//           type: integer
	//           format: int64
	//         outgoing_connections:
	//           type: integer
	//           format: int64
	//         incoming_connections:
	//           type: integer
	//           format: int64
	//         uptime:
	//           type: string
	//         csrf_enabled:
	//             type: boolean
	//         csp_enabled:
	//             type: boolean
	//         wallet_api_enabled:
	//             type: boolean
	//         gui_enabled:
	//             type: boolean
	//         unversioned_api_enabled:
	//             type: boolean
	//         json_rpc_enabled:
	//             type: boolean
	//         user_verify_transaction:
	//           type: object
	//           properties:
	//             burn_factor:
	//               type: integer
	//               format: int64
	//             max_transaction_size:
	//               type: integer
	//               format: int64
	//             max_decimal:
	//               type: integer
	//               format: int64
	//         unconfirmed_verify_transaction:
	//           type: object
	//           properties:
	//             burn_factor:
	//               type: integer
	//               format: int64
	//             max_transaction_size:
	//               type: integer
	//               format: int64
	//             max_decimal:
	//               type: integer
	//               format: int64
	//         started_at:
	//             type: integer
	//             format: int64
	//
	//   default:
	//     $ref: '#/responses/genericError'

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
			Version:              c.health.BuildInfo,
			CoinName:             c.health.CoinName,
			DaemonUserAgent:      userAgent,
			OpenConnections:      health.OutgoingConnections + health.IncomingConnections,
			OutgoingConnections:  health.OutgoingConnections,
			IncomingConnections:  health.IncomingConnections,
			Uptime:               wh.FromDuration(health.Uptime),
			CSRFEnabled:          !c.disableCSRF,
			HeaderCheckEnabled:   !c.disableHeaderCheck,
			CSPEnabled:           !c.disableCSP,
			GUIEnabled:           c.enableGUI,
			WalletAPIEnabled:     walletAPIEnabled,
			UserVerifyTxn:        readable.NewVerifyTxn(params.UserVerifyTxn),
			UnconfirmedVerifyTxn: readable.NewVerifyTxn(health.UnconfirmedVerifyTxn),
			StartedAt:            health.StartedAt.Unix(),
		})
	}
}
