package commands

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// A node that was down at startup gets its name once it answers.
func TestRediscoverCoinsFillsPlaceholder(t *testing.T) {
	node := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"fiber":{"name":"skycoin","display_name":"Skycoin","ticker":"SKY","coin_hours_display_name":"Coin Hours"}}`)) //nolint:errcheck
	}))
	defer node.Close()

	coin := &discoveredCoin{ID: 0, CoinName: "Node 0", CoinSymbol: "N0", placeholder: true, remoteNodeURL: node.URL}
	rediscoverCoins(context.Background(), []*discoveredCoin{coin})

	require.False(t, coin.placeholder)
	require.Equal(t, "Skycoin", coin.CoinName)
	require.Equal(t, "SKY", coin.CoinSymbol)
}
