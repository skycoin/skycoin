package skycoin

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
)

func TestInitAPISetsGUIEnablesWallet(t *testing.T) {
	c := &Config{Node: NodeConfig{
		EnabledAPISets: api.EndpointsRead + "," + api.EndpointsTransaction,
		EnableGUI:      true,
	}}
	require.NoError(t, c.initAPISets())
	require.True(t, c.Node.EnableGUI)
	require.Contains(t, c.Node.enabledAPISets, api.EndpointsWallet)

	c.Node.DisabledAPISets = api.EndpointsWallet
	require.Error(t, c.initAPISets())

	c.Node.EnableGUI = false
	c.Node.LaunchBrowser = true
	require.NoError(t, c.initAPISets())
	require.NotContains(t, c.Node.enabledAPISets, api.EndpointsWallet)
	require.False(t, c.Node.LaunchBrowser)
}
