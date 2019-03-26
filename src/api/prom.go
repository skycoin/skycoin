package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	wh "github.com/skycoin/skycoin/src/util/http"
	"net/http"
)

var (
	promMetricHealthSeq = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "health_head_seq",
		Help: "Health -> head sequence in the block chain",
	})
	promMetricHealthFee = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_head_fee",
		Help: "Health -> head fee in the block chain",
	})
	promMetricHealthUnspents = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_unspents",
		Help: "Health -> unspent transactions in the block chain",
	})
	promMetricHealthUnconfirmed = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_unconfirmed",
		Help: "Health -> unconfirmed transactions in the block chain",
	})
	promMetricHealthOpenConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_open_connections",
		Help: "Health -> open connections in the node",
	})
	promMetricHealthOutgoingConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_outgoing_connections",
		Help: "Health -> outgoing connections in the node",
	})
	promMetricHealthIncomingConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_incoming_connections",
		Help: "Health -> incoming connections in the node",
	})
	promMetricHealthUserVerifyBurnFactor = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_user_verify_burn_factor",
		Help: "Health -> user verify burn factor in the block chain",
	})
	promMetricHealthUserVerifyMaxTransactionSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_user_verify_max_transaction_size",
		Help: "Health -> user verify max transaction size in the block chain",
	})
	promMetricHealthUserVerifyMaxDecimals = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_user_verify_max_decimals",
		Help: "Health -> user verify max decimals in the block chain",
	})
	promMetricHealthUnconfirmedVerifyBurnFactor = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_unconfirmed_verify_burn_factor",
		Help: "Health -> unconfirmed verify burn factor in the block chain",
	})
	promMetricHealthUnconfirmedVerifyMaxTransactionSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_unconfirmed_verify_max_transaction_size",
		Help: "Health -> health unconfirmed verify max transaction size in the block chain",
	})
	promMetricHealthUnconfirmedVerifyMaxDecimals = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "health_unconfirmed_verify_max_decimals",
		Help: "Health -> health unconfirmed verify max decimals in the block chain",
	})
)

func metricsMiddleware(c muxConfig, gateway Gatewayer, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		health, err := getHealthData(c, gateway)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}
		promMetricHealthSeq.Set(float64(health.BlockchainMetadata.Head.BkSeq))
		promMetricHealthFee.Set(float64(health.BlockchainMetadata.Head.Fee))
		promMetricHealthUnspents.Set(float64(health.BlockchainMetadata.Unspents))
		promMetricHealthUnconfirmed.Set(float64(health.BlockchainMetadata.Unconfirmed))
		promMetricHealthOpenConnections.Set(float64(health.OpenConnections))
		promMetricHealthOutgoingConnections.Set(float64(health.OutgoingConnections))
		promMetricHealthIncomingConnections.Set(float64(health.IncomingConnections))
		promMetricHealthUserVerifyBurnFactor.Set(float64(health.UserVerifyTxn.BurnFactor))
		promMetricHealthUserVerifyMaxTransactionSize.Set(float64(health.UserVerifyTxn.MaxTransactionSize))
		promMetricHealthUserVerifyMaxDecimals.Set(float64(health.UserVerifyTxn.MaxDropletPrecision))
		promMetricHealthUnconfirmedVerifyBurnFactor.Set(float64(health.UnconfirmedVerifyTxn.BurnFactor))
		promMetricHealthUnconfirmedVerifyMaxTransactionSize.Set(float64(health.UnconfirmedVerifyTxn.MaxTransactionSize))
		promMetricHealthUnconfirmedVerifyMaxDecimals.Set(float64(health.UnconfirmedVerifyTxn.MaxDropletPrecision))
		next.ServeHTTP(w, r)
	})
}

// metricsHandler returns node health data
// URI: /api/v2/metrics
// Method: GET
func metricsHandler(c muxConfig, gateway Gatewayer) http.HandlerFunc {
	return metricsMiddleware(c, gateway, promhttp.Handler())
}

func init() {
	// NOTE(denisacostaq@gmail.com): Register prometheus metrics
	prometheus.MustRegister(promMetricHealthSeq)
	prometheus.MustRegister(promMetricHealthFee)
	prometheus.MustRegister(promMetricHealthUnspents)
	prometheus.MustRegister(promMetricHealthUnconfirmed)
	prometheus.MustRegister(promMetricHealthOpenConnections)
	prometheus.MustRegister(promMetricHealthOutgoingConnections)
	prometheus.MustRegister(promMetricHealthIncomingConnections)
	prometheus.MustRegister(promMetricHealthUserVerifyBurnFactor)
	prometheus.MustRegister(promMetricHealthUserVerifyMaxTransactionSize)
	prometheus.MustRegister(promMetricHealthUserVerifyMaxDecimals)
	prometheus.MustRegister(promMetricHealthUnconfirmedVerifyBurnFactor)
	prometheus.MustRegister(promMetricHealthUnconfirmedVerifyMaxTransactionSize)
	prometheus.MustRegister(promMetricHealthUnconfirmedVerifyMaxDecimals)
}
