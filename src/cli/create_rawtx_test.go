package cli

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/util/fee"
)

func TestMain(m *testing.M) {
	params.UserVerifyTxn.BurnFactor = 10
}

func TestMakeChangeOut(t *testing.T) {
	// single destination test
	uxOuts := []transaction.UxBalance{
		{
			Hash:    cipher.MustSHA256FromHex("f569461182b0efe9a5c666e9a35c6602b351021c1803cc740aca548cf6db4cb2"),
			Address: cipher.MustDecodeBase58Address("k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND"),
			BkSeq:   10,
			Coins:   400e6,
			Hours:   1,
		},
		{
			Hash:    cipher.MustSHA256FromHex("bddf0aaf80f96c144f33ac8a27764a868d37e1c11e568063ebeb1367de859566"),
			Address: cipher.MustDecodeBase58Address("A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe"),
			BkSeq:   11,
			Coins:   300e6,
			Hours:   1,
		},
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

	chgOut := txOuts[1]
	t.Logf("chgOut:%+v\n", chgOut)
	require.Equal(t, chgAddr, chgOut.Address.String())
	require.Exactly(t, uint64(100e6), chgOut.Coins)

	spendOut := txOuts[0]
	require.Equal(t, spendAmt[0].Addr, spendOut.Address.String())
	require.Exactly(t, spendAmt[0].Coins, spendOut.Coins)

	require.Exactly(t, uint64(1), chgOut.Hours)
	require.Exactly(t, uint64(0), spendOut.Hours)

	// multiple destination test
	uxOuts = []transaction.UxBalance{
		{
			Hash:    cipher.MustSHA256FromHex("f569461182b0efe9a5c666e9a35c6602b351021c1803cc740aca548cf6db4cb2"),
			Address: cipher.MustDecodeBase58Address("k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND"),
			BkSeq:   10,
			Coins:   10e6,
			Hours:   8,
		},
		{
			Hash:    cipher.MustSHA256FromHex("bddf0aaf80f96c144f33ac8a27764a868d37e1c11e568063ebeb1367de859566"),
			Address: cipher.MustDecodeBase58Address("A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe"),
			BkSeq:   11,
			Coins:   5e6,
			Hours:   16,
		},
	}

	spendAmt = []SendAmount{
		{
			Addr:  "2PBmUva7J8WFsyWg979cREZkU3z2pkYjNkE",
			Coins: 1e6,
		},
		{
			Addr:  "2CgSQ4FbtfbP6fJqmy75WwkW2tNsKPL2zzp",
			Coins: 2e6,
		}}

	_, err = cipher.DecodeBase58Address(chgAddr)
	require.NoError(t, err)

	txOuts, err = makeChangeOut(uxOuts, chgAddr, spendAmt)
	require.NoError(t, err)
	require.NotEmpty(t, txOuts)

	// Should have a change output and two output to the destination in toAddrs
	require.Len(t, txOuts, 3)

	chgOut = txOuts[2]
	t.Logf("chgOut:%+v\n", chgOut)
	require.Equal(t, chgAddr, chgOut.Address.String())
	require.Exactly(t, uint64(12e6), chgOut.Coins)

	for i := range spendAmt {
		require.Equal(t, spendAmt[i].Addr, txOuts[i].Address.String())
		require.Exactly(t, spendAmt[i].Coins, txOuts[i].Coins)
	}

	require.Exactly(t, uint64(18), chgOut.Hours)
	require.Exactly(t, uint64(1), txOuts[0].Hours)
	require.Exactly(t, uint64(2), txOuts[1].Hours)
}

func TestMakeChangeOutMinOneCoinHourSend(t *testing.T) {
	uxOuts := []transaction.UxBalance{
		{
			Hash:    cipher.MustSHA256FromHex("f569461182b0efe9a5c666e9a35c6602b351021c1803cc740aca548cf6db4cb2"),
			Address: cipher.MustDecodeBase58Address("k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND"),
			BkSeq:   10,
			Coins:   400e6,
			Hours:   200,
		},
		{
			Hash:    cipher.MustSHA256FromHex("bddf0aaf80f96c144f33ac8a27764a868d37e1c11e568063ebeb1367de859566"),
			Address: cipher.MustDecodeBase58Address("A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe"),
			BkSeq:   11,
			Coins:   300e6,
			Hours:   100,
		},
	}

	spendAmt := []SendAmount{{
		Addr:  "2PBmUva7J8WFsyWg979cREZkU3z2pkYjNkE",
		Coins: 0.001e6,
	}}

	chgAddr := "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ"
	_, err := cipher.DecodeBase58Address(chgAddr)
	require.NoError(t, err)

	txOuts, err := makeChangeOut(uxOuts, chgAddr, spendAmt)
	require.NoError(t, err)
	require.NotEmpty(t, txOuts)

	// Should have a change output and an output to the destination in toAddrs
	require.Len(t, txOuts, 2)

	chgOut := txOuts[1]
	t.Logf("chgOut:%+v\n", chgOut)
	require.Equal(t, chgAddr, chgOut.Address.String())
	require.Exactly(t, uint64(699.999e6), chgOut.Coins)

	spendOut := txOuts[0]
	require.Equal(t, spendAmt[0].Addr, spendOut.Address.String())
	require.Exactly(t, spendAmt[0].Coins, spendOut.Coins)

	require.Exactly(t, uint64(269), chgOut.Hours)
	require.Exactly(t, uint64(1), spendOut.Hours)
}

func TestMakeChangeOutCoinHourCap(t *testing.T) {
	uxOuts := []transaction.UxBalance{
		{
			Hash:    cipher.MustSHA256FromHex("f569461182b0efe9a5c666e9a35c6602b351021c1803cc740aca548cf6db4cb2"),
			Address: cipher.MustDecodeBase58Address("k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND"),
			BkSeq:   10,
			Coins:   400e6,
			Hours:   2000,
		},
		{
			Hash:    cipher.MustSHA256FromHex("bddf0aaf80f96c144f33ac8a27764a868d37e1c11e568063ebeb1367de859566"),
			Address: cipher.MustDecodeBase58Address("A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe"),
			BkSeq:   11,
			Coins:   300e6,
			Hours:   1000,
		},
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

	chgOut := txOuts[1]
	t.Logf("chgOut:%+v\n", chgOut)
	require.Equal(t, chgAddr, chgOut.Address.String())
	require.Exactly(t, uint64(100e6), chgOut.Coins)

	spendOut := txOuts[0]
	require.Equal(t, spendAmt[0].Addr, spendOut.Address.String())
	require.Exactly(t, spendAmt[0].Coins, spendOut.Coins)

	require.Exactly(t, uint64(2100), chgOut.Hours)
	require.Exactly(t, uint64(600), spendOut.Hours)
}

func TestMakeChangeOutOneCoinHour(t *testing.T) {
	// As long as there is at least one coin hour left, creating a transaction
	// will still succeed
	uxOuts := []transaction.UxBalance{
		{
			Hash:    cipher.MustSHA256FromHex("f569461182b0efe9a5c666e9a35c6602b351021c1803cc740aca548cf6db4cb2"),
			BkSeq:   10,
			Address: cipher.MustDecodeBase58Address("k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND"),
			Coins:   400e6,
			Hours:   0,
		},
		{
			Hash:    cipher.MustSHA256FromHex("bddf0aaf80f96c144f33ac8a27764a868d37e1c11e568063ebeb1367de859566"),
			BkSeq:   11,
			Address: cipher.MustDecodeBase58Address("A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe"),
			Coins:   300e6,
			Hours:   1,
		},
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

	chgOut := txOuts[1]
	t.Logf("chgOut:%+v\n", chgOut)
	require.Equal(t, chgAddr, chgOut.Address.String())
	require.Exactly(t, uint64(100e6), chgOut.Coins)
	require.Exactly(t, uint64(0), chgOut.Hours)

	spendOut := txOuts[0]
	require.Equal(t, spendAmt[0].Addr, spendOut.Address.String())
	require.Exactly(t, spendAmt[0].Coins, spendOut.Coins)
	require.Exactly(t, uint64(0), spendOut.Hours)
}

func TestMakeChangeOutInsufficientCoinHours(t *testing.T) {
	// If there are no coin hours in the inputs, creating the txn will fail
	// because it will not be accepted by the network
	uxOuts := []transaction.UxBalance{
		{
			Hash:    cipher.MustSHA256FromHex("f569461182b0efe9a5c666e9a35c6602b351021c1803cc740aca548cf6db4cb2"),
			BkSeq:   10,
			Address: cipher.MustDecodeBase58Address("k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND"),
			Coins:   400e6,
			Hours:   0,
		},
		{
			Hash:    cipher.MustSHA256FromHex("bddf0aaf80f96c144f33ac8a27764a868d37e1c11e568063ebeb1367de859566"),
			BkSeq:   11,
			Address: cipher.MustDecodeBase58Address("A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe"),
			Coins:   300e6,
			Hours:   0,
		},
	}

	spendAmt := []SendAmount{{
		Addr:  "2PBmUva7J8WFsyWg979cREZkU3z2pkYjNkE",
		Coins: 600e6,
	}}

	chgAddr := "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ"
	_, err := cipher.DecodeBase58Address(chgAddr)
	require.NoError(t, err)

	_, err = makeChangeOut(uxOuts, chgAddr, spendAmt)
	testutil.RequireError(t, err, fee.ErrTxnNoFee.Error())
}

func TestChooseSpends(t *testing.T) {
	// Start with readable.UnspentOutputsSummary
	// Spends should be minimized

	// Insufficient HeadOutputs
	// Sufficient HeadOutputs, but insufficient after adjusting for OutgoingOutputs
	// Insufficient HeadOutputs, but sufficient after adjusting for IncomingOutputs
	// Sufficient HeadOutputs after adjusting for OutgoingOutputs

	var coins uint64 = 100e6

	hashA := testutil.RandSHA256(t).Hex()
	hashB := testutil.RandSHA256(t).Hex()
	hashC := testutil.RandSHA256(t).Hex()
	hashD := testutil.RandSHA256(t).Hex()

	addrA := testutil.MakeAddress().String()
	addrB := testutil.MakeAddress().String()
	addrC := testutil.MakeAddress().String()
	addrD := testutil.MakeAddress().String()

	cases := []struct {
		name     string
		err      error
		spendLen int
		ros      readable.UnspentOutputsSummary
	}{
		{
			"Insufficient HeadOutputs",
			transaction.ErrInsufficientBalance,
			0,
			readable.UnspentOutputsSummary{
				HeadOutputs: readable.UnspentOutputs{
					{
						Hash:              hashA,
						Address:           addrA,
						BkSeq:             22,
						Coins:             "75.000000",
						CalculatedHours:   100,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
					{
						Hash:              hashB,
						Address:           addrB,
						BkSeq:             19,
						Coins:             "13.000000",
						CalculatedHours:   0,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
				},
			},
		},

		{
			"Sufficient HeadOutputs, but insufficient after subtracting OutgoingOutputs",
			transaction.ErrInsufficientBalance,
			0,
			readable.UnspentOutputsSummary{
				HeadOutputs: readable.UnspentOutputs{
					{
						Hash:              hashA,
						Address:           addrA,
						BkSeq:             22,
						Coins:             "75.000000",
						CalculatedHours:   100,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
					{
						Hash:              hashB,
						Address:           addrB,
						BkSeq:             19,
						Coins:             "50.000000",
						CalculatedHours:   0,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
				},
				OutgoingOutputs: readable.UnspentOutputs{
					{
						Hash:              hashB,
						Address:           addrB,
						BkSeq:             19,
						Coins:             "50.000000",
						CalculatedHours:   0,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
				},
			},
		},

		{
			"Insufficient HeadOutputs, but sufficient after adding IncomingOutputs",
			ErrTemporaryInsufficientBalance,
			0,
			readable.UnspentOutputsSummary{
				HeadOutputs: readable.UnspentOutputs{
					{
						Hash:              hashA,
						Address:           addrA,
						BkSeq:             22,
						Coins:             "20.000000",
						CalculatedHours:   100,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
					{
						Hash:              hashB,
						Address:           addrB,
						BkSeq:             19,
						Coins:             "30.000000",
						CalculatedHours:   0,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
				},
				IncomingOutputs: readable.UnspentOutputs{
					{
						Hash:              hashC,
						Address:           addrC,
						BkSeq:             134,
						Coins:             "40.000000",
						CalculatedHours:   200,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
					{
						Hash:              hashD,
						Address:           addrD,
						BkSeq:             29,
						Coins:             "11.000000",
						CalculatedHours:   0,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
				},
			},
		},

		{
			"Sufficient HeadOutputs and still sufficient after subtracting OutgoingOutputs",
			nil,
			2,
			readable.UnspentOutputsSummary{
				HeadOutputs: readable.UnspentOutputs{
					{
						Hash:              hashA,
						Address:           addrA,
						BkSeq:             22,
						Coins:             "15.000000",
						CalculatedHours:   100,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
					{
						Hash:              hashB,
						Address:           addrB,
						BkSeq:             19,
						Coins:             "90.000000",
						CalculatedHours:   0,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
					{
						Hash:              hashC,
						Address:           addrC,
						BkSeq:             19,
						Coins:             "20.000000",
						CalculatedHours:   1,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
				},
				OutgoingOutputs: readable.UnspentOutputs{
					{
						Hash:              hashA,
						Address:           addrA,
						BkSeq:             22,
						Coins:             "15.000000",
						CalculatedHours:   100,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
				},
			},
		},

		{
			"Sufficient HeadOutputs and still sufficient after subtracting OutgoingOutputs but will have no coinhours",
			fee.ErrTxnNoFee,
			0,
			readable.UnspentOutputsSummary{
				HeadOutputs: readable.UnspentOutputs{
					{
						Hash:              hashA,
						Address:           addrA,
						BkSeq:             22,
						Coins:             "15.000000",
						CalculatedHours:   100,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
					{
						Hash:              hashB,
						Address:           addrB,
						BkSeq:             19,
						Coins:             "90.000000",
						CalculatedHours:   0,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
					{
						Hash:              hashC,
						Address:           addrC,
						BkSeq:             19,
						Coins:             "30.000000",
						CalculatedHours:   0,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
				},
				OutgoingOutputs: readable.UnspentOutputs{
					{
						Hash:              hashA,
						Address:           addrA,
						BkSeq:             22,
						Coins:             "15.000000",
						CalculatedHours:   100,
						SourceTransaction: testutil.RandSHA256(t).Hex(),
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			spends, err := chooseSpends(&tc.ros, coins)

			if tc.err != nil {
				testutil.RequireError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.spendLen, len(spends))

				var totalCoins uint64
				for _, ux := range spends {
					totalCoins += ux.Coins
				}

				require.True(t, coins <= totalCoins)
			}
		})
	}
}

func TestParseSendAmountsFromCSV(t *testing.T) {
	cases := []struct {
		name   string
		fields [][]string
		amts   []SendAmount
		err    error
	}{
		{
			name: "valid simple case",
			fields: [][]string{
				{"2Niqzo12tZ9ioZq5vwPHMVR4g7UVpp9TCmP", "123"},
				{"2UDzBKnxZf4d9pdrBJAqbtoeH641RFLYKxd", "123.456"},
				{"8LbGZ9Z9r7ELNKyrQmAbhLhLvrmLJjfotm", "123.456789"},
				{"7KU683yzoPE9rVuuFRQMZVhGwBBtwqTKT2", "0"},
			},
			amts: []SendAmount{
				{
					Addr:  "2Niqzo12tZ9ioZq5vwPHMVR4g7UVpp9TCmP",
					Coins: 123e6,
				},
				{
					Addr:  "2UDzBKnxZf4d9pdrBJAqbtoeH641RFLYKxd",
					Coins: 123456e3,
				},
				{
					Addr:  "8LbGZ9Z9r7ELNKyrQmAbhLhLvrmLJjfotm",
					Coins: 123456789,
				},
				{
					Addr:  "7KU683yzoPE9rVuuFRQMZVhGwBBtwqTKT2",
					Coins: 0,
				},
			},
		},

		{
			name: "invalid coins value",
			fields: [][]string{
				{"7KU683yzoPE9rVuuFRQMZVhGwBBtwqTKT2", "0.1234567"},
			},
			err: errors.New("[row 0] Invalid amount 0.1234567: Droplet string conversion failed: Too many decimal places"),
		},

		{
			name: "invalid address value",
			fields: [][]string{
				{"xxx", "0.123"},
			},
			err: errors.New("[row 0] Invalid address xxx: Invalid address length"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			amts, err := parseSendAmountsFromCSV(tc.fields)

			if tc.err != nil {
				require.Equal(t, tc.err, err)
				require.Nil(t, amts)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.amts, amts)
		})
	}
}
