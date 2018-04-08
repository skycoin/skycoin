package gui

import (
	"net/http"
	"time"

	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor"
)

// BlockchainMetadata extends visor.BlockchainMetadata to include the time since the last block
type BlockchainMetadata struct {
	*visor.BlockchainMetadata
	TimeSinceLastBlock wh.Duration `json:"time_since_last_block"`
}

// HealthResponse is returned by the /health endpoint
type HealthResponse struct {
	Blockchain      BlockchainMetadata `json:"blockchain"`
	Version         visor.BuildInfo    `json:"version"`
	OpenConnections int                `json:"open_connections"`
	Uptime          wh.Duration        `json:"uptime"`
}

// Returns node health data.
// URI: /health
// Method: GET
func healthCheck(gateway Gatewayer, startedAt time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		metadata, err := gateway.GetBlockchainMetadata()
		if err != nil {
			logger.WithError(err).Error("gateway.GetBlockchainMetadata failed")
			wh.Error500Msg(w, err.Error())
			return
		}

		elapsedBlockTime := time.Now().UTC().Unix() - int64(metadata.Head.Time)
		timeSinceLastBlock := time.Second * time.Duration(elapsedBlockTime)

		resp := &HealthResponse{
			Blockchain: BlockchainMetadata{
				BlockchainMetadata: metadata,
				TimeSinceLastBlock: wh.FromDuration(timeSinceLastBlock),
			},
			Version:         gateway.GetBuildInfo(),
			OpenConnections: len(gateway.GetConnections().Connections),
			Uptime:          wh.FromDuration(time.Since(startedAt)),
		}

		wh.SendJSONOr500(logger, w, resp)
	}
}
