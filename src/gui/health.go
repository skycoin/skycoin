package gui

import (
	"net/http"

	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor"
)

const Version = "0.21.1"

type HealthResponse struct {
	BlockChainMetadata  *visor.BlockchainMetadata `json:"blockchain_metadata"`
	VersionData         string                    `json:"version_data"`
	UnconfirmedTxCount  int                       `json:"unconfirmed_tx_count"`
	OpenConnectionCount int                       `json:"open_connection_count"`
}

// Health status of application
func healthCheck(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		txnCount := len(gateway.GetAllUnconfirmedTxns())
		connectionCount := len(gateway.GetConnections().Connections)
		metadata := gateway.GetBlockchainMetadata()

		resp := HealthResponse{
			BlockChainMetadata:  metadata,
			VersionData:         Version,
			UnconfirmedTxCount:  txnCount,
			OpenConnectionCount: connectionCount,
		}

		wh.SendOr404(w, &resp)
	}
}
