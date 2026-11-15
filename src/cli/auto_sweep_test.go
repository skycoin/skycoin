package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/wallet/sweep"
)

func TestApplySweepFlags(t *testing.T) {
	set := func(names ...string) func(string) bool {
		m := make(map[string]bool)
		for _, n := range names {
			m[n] = true
		}
		return func(n string) bool { return m[n] }
	}
	two := uint32(2)
	base := sweep.Rule{Wallet: "w", SharedAccounts: []uint32{0}, Destination: sweep.Destination{Kind: sweep.DestAccount, Account: &two}, Enabled: true}

	r, err := applySweepFlags(base, sweepFlags{}, set())
	require.NoError(t, err)
	require.Equal(t, base, r, "no flags keep the rule")

	r, err = applySweepFlags(base, sweepFlags{shared: []uint{0, 1}, toAddress: "x", disable: true}, set("shared-account", "to-address"))
	require.NoError(t, err)
	require.Equal(t, []uint32{0, 1}, r.SharedAccounts)
	require.Equal(t, sweep.Destination{Kind: sweep.DestAddress, Address: "x"}, r.Destination)
	require.False(t, r.Enabled)

	r, err = applySweepFlags(base, sweepFlags{}, set("to-next-account"))
	require.NoError(t, err)
	require.Equal(t, sweep.Destination{Kind: sweep.DestAccount}, r.Destination)

	_, err = applySweepFlags(base, sweepFlags{}, set("to-wallet", "to-address"))
	require.Error(t, err)
	_, err = applySweepFlags(base, sweepFlags{enable: true, disable: true}, set())
	require.Error(t, err)
}
