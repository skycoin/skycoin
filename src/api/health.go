package api

import (
	"fmt"
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
	BlockchainMetadata BlockchainMetadata `json:"blockchain"`
	Version            visor.BuildInfo    `json:"version"`
	OpenConnections    int                `json:"open_connections"`
	Uptime             wh.Duration        `json:"uptime"`
}

// Returns node health data.
// URI: /api/v1/health
// Method: GET
func healthCheck(gateway Gatewayer) http.HandlerFunc {
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

		elapsedBlockTime := time.Now().UTC().Unix() - int64(health.BlockchainMetadata.Head.Time)
		timeSinceLastBlock := time.Second * time.Duration(elapsedBlockTime)

		wh.SendJSONOr500(logger, w, HealthResponse{
			BlockchainMetadata: BlockchainMetadata{
				BlockchainMetadata: health.BlockchainMetadata,
				TimeSinceLastBlock: wh.FromDuration(timeSinceLastBlock),
			},
			Version:         health.Version,
			OpenConnections: health.OpenConnections,
			Uptime:          wh.FromDuration(health.Uptime),
		})
	}
}
