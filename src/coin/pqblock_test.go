package coin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPQBlockRuleApplies(t *testing.T) {
	require.False(t, PQBlockRule{ForkSeq: 1}.Applies(5), "no key, no rule")

	r := PQBlockRule{Pubkey: []byte{1}, ForkSeq: 10}
	require.False(t, r.Applies(0))
	require.False(t, r.Applies(9))
	require.True(t, r.Applies(10))
	require.True(t, r.Applies(11))

	r.ForkSeq = 0
	require.False(t, r.Applies(0), "genesis is exempt")
	require.True(t, r.Applies(1))
}
