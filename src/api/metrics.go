package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	wh "github.com/SkycoinProject/skycoin/src/util/http"
)

var (
	promUnspents = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "unspent_outputs",
			Help: "Number of unspent outputs",
		})
	promUnconfirmedTxns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "unconfirmed_txns",
			Help: "Number of unconfirmed transactions",
		})
	promTimeSinceLastBlock = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "time_since_last_block_seconds",
			Help: "Time since the last block created",
		})
	promOpenConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "open_connections",
			Help: "Number of open connections",
		})
	promOutgoingConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "outgoing_connections",
			Help: "Number of outgoing connections",
		})
	promIncomingConns = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "incoming_connections",
			Help: "Number of incoming connections",
		})
	promStartedAt = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "started_at",
			Help: "Node start time, in unixtime",
		})
	promLastBlockSeq = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "last_block_seq",
			Help: "Last block sequence number",
		})
)

func init() {
	prometheus.MustRegister(promUnspents)
	prometheus.MustRegister(promUnconfirmedTxns)
	prometheus.MustRegister(promTimeSinceLastBlock)
	prometheus.MustRegister(promOpenConns)
	prometheus.MustRegister(promOutgoingConns)
	prometheus.MustRegister(promIncomingConns)
	prometheus.MustRegister(promStartedAt)
	prometheus.MustRegister(promLastBlockSeq)
}

func metricsHandler(c muxConfig, gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health, err := getHealthData(c, gateway)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		promUnspents.Set(float64(health.BlockchainMetadata.Unspents))
		promUnconfirmedTxns.Set(float64(health.BlockchainMetadata.Unconfirmed))
		promTimeSinceLastBlock.Set(health.BlockchainMetadata.TimeSinceLastBlock.Seconds())
		promOpenConns.Set(float64(health.OpenConnections))
		promOutgoingConns.Set(float64(health.OutgoingConnections))
		promIncomingConns.Set(float64(health.IncomingConnections))
		promStartedAt.Set(float64(gateway.StartedAt().Unix()))
		promLastBlockSeq.Set(float64(health.BlockchainMetadata.Head.BkSeq))

		promhttp.Handler().ServeHTTP(w, r)
	}
}
