package gui

import (
	"net/http"

	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor"
)

type HealthResponse struct {
	BlockChainMetadata  *visor.BlockchainMetadata `json:"blockchain_metadata"`
	VersionData         visor.BuildInfo           `json:"version_data"`
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

		resp := &HealthResponse{
			BlockChainMetadata:  metadata,
			VersionData:         gateway.GetBuildInfo(),
			UnconfirmedTxCount:  txnCount,
			OpenConnectionCount: connectionCount,
		}

		if resp == nil {
			wh.Error404(w)
			return
		}

		wh.SendJSONOr500(logger, w, resp)
	}
}
