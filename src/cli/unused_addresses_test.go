package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnusedAddressesCmd(t *testing.T) {
	cmd := unusedAddressesCmd()
	require.NotNil(t, cmd)
	require.Equal(t, "unusedAddresses [wallet]", cmd.Use)

	// Test flags exist
	numFlag := cmd.Flags().Lookup("num")
	require.NotNil(t, numFlag)
	require.Equal(t, "n", numFlag.Shorthand)

	generateFlag := cmd.Flags().Lookup("generate")
	require.NotNil(t, generateFlag)
	require.Equal(t, "g", generateFlag.Shorthand)
}
