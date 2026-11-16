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

func TestApplySweepMaps(t *testing.T) {
	set := func(names ...string) func(string) bool {
		m := make(map[string]bool)
		for _, n := range names {
			m[n] = true
		}
		return func(n string) bool { return m[n] }
	}
	two, five := uint32(2), uint32(5)
	base := sweep.Rule{Wallet: "w", Mappings: []sweep.Mapping{{From: sweep.Source{Account: &two}, To: sweep.Destination{Kind: sweep.DestWallet}}}}

	r, err := applySweepFlags(base, sweepFlags{maps: []string{
		"2=wallet:cold.wlt",
		"account:3=next-account",
		"addrA=account:5",
		"address:addrB=address:addrC",
		"addrD=addrE",
		"addrF=wallet",
	}}, set("map"))
	require.NoError(t, err)
	three := uint32(3)
	require.Equal(t, []sweep.Mapping{
		{From: sweep.Source{Account: &two}, To: sweep.Destination{Kind: sweep.DestWallet, Wallet: "cold.wlt"}},
		{From: sweep.Source{Account: &three}, To: sweep.Destination{Kind: sweep.DestAccount}},
		{From: sweep.Source{Address: "addrA"}, To: sweep.Destination{Kind: sweep.DestAccount, Account: &five}},
		{From: sweep.Source{Address: "addrB"}, To: sweep.Destination{Kind: sweep.DestAddress, Address: "addrC"}},
		{From: sweep.Source{Address: "addrD"}, To: sweep.Destination{Kind: sweep.DestAddress, Address: "addrE"}},
		{From: sweep.Source{Address: "addrF"}, To: sweep.Destination{Kind: sweep.DestWallet}},
	}, r.Mappings)
	require.Nil(t, r.IncludeSpent)
	require.False(t, r.Enabled, "mappings alone do not turn level 1 on")

	r, err = applySweepFlags(base, sweepFlags{}, set("clear-maps"))
	require.NoError(t, err)
	require.Empty(t, r.Mappings)

	r, err = applySweepFlags(base, sweepFlags{includeSpent: false}, set("include-spent"))
	require.NoError(t, err)
	require.NotNil(t, r.IncludeSpent)
	require.False(t, *r.IncludeSpent)
	r, err = applySweepFlags(r, sweepFlags{includeSpent: true}, set("include-spent"))
	require.NoError(t, err)
	require.Nil(t, r.IncludeSpent, "true is the default and is left out")

	for _, bad := range []string{"nodest", "=wallet", "2=", "2=account:x"} {
		_, err = applySweepFlags(base, sweepFlags{maps: []string{bad}}, set("map"))
		require.Error(t, err, bad)
	}
	_, err = applySweepFlags(base, sweepFlags{maps: []string{"2=wallet"}}, set("map", "clear-maps"))
	require.Error(t, err)
}
