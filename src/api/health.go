package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/skycoin/skycoin/src/daemon"
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

func getHealthData(c muxConfig, gateway Gatewayer) (HealthResponse, error) {
	metadata, err := gateway.GetBlockchainMetadata()
	if err != nil {
		err = fmt.Errorf("gateway.GetBlockchainMetadata failed: %v", err)
		return HealthResponse{}, err
	}

	conns, err := gateway.GetConnections(func(c daemon.Connection) bool {
		return c.State != daemon.ConnectionStatePending
	})
	if err != nil {
		err = fmt.Errorf("gateway.GetConnections failed: %v", err)
		return HealthResponse{}, err
	}

	outgoingConns := 0
	incomingConns := 0
	for _, c := range conns {
		if c.Outgoing {
			outgoingConns++
		} else {
			incomingConns++
		}
	}

	elapsedBlockTime := time.Now().UTC().Unix() - int64(metadata.HeadBlock.Head.Time)
	timeSinceLastBlock := time.Second * time.Duration(elapsedBlockTime)

	_, walletAPIEnabled := c.enabledAPISets[EndpointsWallet]

	userAgent, err := c.health.DaemonUserAgent.Build()
	if err != nil {
		return HealthResponse{}, err
	}
	return HealthResponse{
		BlockchainMetadata: BlockchainMetadata{
			BlockchainMetadata: readable.NewBlockchainMetadata(*metadata),
			TimeSinceLastBlock: wh.FromDuration(timeSinceLastBlock),
		},
		Version:              c.health.BuildInfo,
		CoinName:             c.health.CoinName,
		DaemonUserAgent:      userAgent,
		OpenConnections:      len(conns),
		OutgoingConnections:  outgoingConns,
		IncomingConnections:  incomingConns,
		CSRFEnabled:          !c.disableCSRF,
		HeaderCheckEnabled:   !c.disableHeaderCheck,
		CSPEnabled:           !c.disableCSP,
		GUIEnabled:           c.enableGUI,
		WalletAPIEnabled:     walletAPIEnabled,
		UserVerifyTxn:        readable.NewVerifyTxn(params.UserVerifyTxn),
		UnconfirmedVerifyTxn: readable.NewVerifyTxn(gateway.DaemonConfig().UnconfirmedVerifyTxn),
		Uptime:               wh.FromDuration(time.Since(gateway.StartedAt())),
		StartedAt:            gateway.StartedAt().Unix(),
	}, nil
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
		if health, err := getHealthData(c, gateway); err != nil {
			wh.Error500(w, err.Error())
		} else {
			wh.SendJSONOr500(logger, w, health)
		}
	}
}
