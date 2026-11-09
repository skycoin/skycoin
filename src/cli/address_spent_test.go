package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/readable"
)

func TestSummarizeAddressSpent(t *testing.T) {
	r, err := summarizeAddressSpent("a", []readable.SpentOutput{
		{Coins: 1000000},
		{Coins: 2500000},
	})
	require.NoError(t, err)
	require.Equal(t, AddressSpentResult{Address: "a", Received: 2, Balance: "3.500000"}, r)

	r, err = summarizeAddressSpent("a", []readable.SpentOutput{
		{Coins: 1000000, SpentBlockSeq: 7},
		{Coins: 2000000},
	})
	require.NoError(t, err)
	require.Equal(t, AddressSpentResult{Address: "a", Spent: true, Received: 2, SpentOutputs: 1, Balance: "2.000000"}, r)

	r, err = summarizeAddressSpent("a", nil)
	require.NoError(t, err)
	require.False(t, r.Spent)
}
