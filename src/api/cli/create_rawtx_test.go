package cli

import (
	"strconv"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/stretchr/testify/require"
)

func TestCreateRawTx(t *testing.T) {
	// Need fake gateway

}

func TestMakeChangeOut(t *testing.T) {
	uxOuts := []UnspentOut{
		{visor.ReadableOutput{
			Hash:              "",
			SourceTransaction: "",
			Address:           "k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND",
			Coins:             strconv.Itoa(400),
			Hours:             200,
		}},
		{visor.ReadableOutput{
			Hash:              "",
			SourceTransaction: "",
			Address:           "A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe",
			Coins:             strconv.Itoa(300),
			Hours:             100,
		}},
	}

	spendAmt := []SendAmount{{
		Addr:  "2PBmUva7J8WFsyWg979cREZkU3z2pkYjNkE",
		Coins: 600e6,
	}}

	chgAddr := "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ"
	_, err := cipher.DecodeBase58Address(chgAddr)
	require.NoError(t, err)

	txOuts, err := makeChangeOut(uxOuts, chgAddr, spendAmt)
	require.NoError(t, err)
	require.NotEmpty(t, txOuts)

	// Should have a change output and an output to the destination in toAddrs
	require.Len(t, txOuts, 2)

	chgOut := txOuts[0]
	t.Logf("chgOut:%+v\n", chgOut)
	require.Equal(t, chgAddr, chgOut.Address.String())
	require.Exactly(t, uint64(100e6), chgOut.Coins)
	require.Exactly(t, uint64(300/8), chgOut.Hours)

	spendOut := txOuts[1]
	require.Equal(t, spendAmt[0].Addr, spendOut.Address.String())
	require.Exactly(t, spendAmt[0].Coins, spendOut.Coins)
	require.Exactly(t, uint64(300/4), spendOut.Hours)
}
