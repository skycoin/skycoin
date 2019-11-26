package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/readable"
	"github.com/SkycoinProject/skycoin/src/testutil"
)

func TestGetBalanceOfAddresses(t *testing.T) {
	addrs := []string{
		testutil.MakeAddress().String(),
		testutil.MakeAddress().String(),
		testutil.MakeAddress().String(),
	}

	hashes := make([]string, 10)
	for i := 0; i < len(hashes); i++ {
		h := testutil.RandSHA256(t)
		hashes[i] = h.Hex()
	}

	cases := []struct {
		name   string
		outs   readable.UnspentOutputsSummary
		addrs  []string
		result *BalanceResult
		err    error
	}{
		{
			name: "confirmed == spendable == expected",
			outs: readable.UnspentOutputsSummary{
				HeadOutputs: readable.UnspentOutputs{
					{
						Hash:            hashes[0],
						Address:         addrs[0],
						Coins:           "99.900000",
						CalculatedHours: 3000,
					},
					{
						Hash:            hashes[1],
						Address:         addrs[0],
						Coins:           "0.100000",
						CalculatedHours: 120000,
					},
					{
						Hash:            hashes[2],
						Address:         addrs[2],
						Coins:           "23.111111",
						CalculatedHours: 123,
					},
				},
			},
			addrs: addrs,
			result: &BalanceResult{
				Confirmed: Balance{
					Coins: "123.111111",
					Hours: "123123",
				},
				Spendable: Balance{
					Coins: "123.111111",
					Hours: "123123",
				},
				Expected: Balance{
					Coins: "123.111111",
					Hours: "123123",
				},
				Addresses: []AddressBalances{
					{
						Confirmed: Balance{
							Coins: "100.000000",
							Hours: "123000",
						},
						Spendable: Balance{
							Coins: "100.000000",
							Hours: "123000",
						},
						Expected: Balance{
							Coins: "100.000000",
							Hours: "123000",
						},
						Address: addrs[0],
					},
					{
						Confirmed: Balance{
							Coins: "0.000000",
							Hours: "0",
						},
						Spendable: Balance{
							Coins: "0.000000",
							Hours: "0",
						},
						Expected: Balance{
							Coins: "0.000000",
							Hours: "0",
						},
						Address: addrs[1],
					},
					{
						Confirmed: Balance{
							Coins: "23.111111",
							Hours: "123",
						},
						Spendable: Balance{
							Coins: "23.111111",
							Hours: "123",
						},
						Expected: Balance{
							Coins: "23.111111",
							Hours: "123",
						},
						Address: addrs[2],
					},
				},
			},
		},

		{
			name: "confirmed != spendable != expected",
			outs: readable.UnspentOutputsSummary{
				HeadOutputs: readable.UnspentOutputs{
					{
						Hash:            hashes[0],
						Address:         addrs[0],
						Coins:           "89.900000",
						CalculatedHours: 3000,
					},
					{
						Hash:            hashes[1],
						Address:         addrs[0],
						Coins:           "0.100000",
						CalculatedHours: 97000,
					},
					{
						Hash:            hashes[5],
						Address:         addrs[0],
						Coins:           "10.000000",
						CalculatedHours: 23000,
					},
					{
						Hash:            hashes[2],
						Address:         addrs[2],
						Coins:           "1.000001",
						CalculatedHours: 23,
					},
					{
						Hash:            hashes[6],
						Address:         addrs[2],
						Coins:           "22.111110",
						CalculatedHours: 100,
					},
				},
				OutgoingOutputs: readable.UnspentOutputs{
					{
						Hash:            hashes[5],
						Address:         addrs[0],
						Coins:           "10.000000",
						CalculatedHours: 23000,
					},
					{
						Hash:            hashes[6],
						Address:         addrs[2],
						Coins:           "22.111110",
						CalculatedHours: 100,
					},
				},
				IncomingOutputs: readable.UnspentOutputs{
					{
						Hash:            hashes[3],
						Address:         addrs[1],
						Coins:           "1.000000",
						CalculatedHours: 333,
					},
					{
						Hash:            hashes[4],
						Address:         addrs[1],
						Coins:           "0.111111",
						CalculatedHours: 0,
					},
					{
						Hash:            hashes[7],
						Address:         addrs[2],
						Coins:           "44.999999",
						CalculatedHours: 433,
					},
				},
			},
			addrs: addrs,
			result: &BalanceResult{
				Confirmed: Balance{
					Coins: "123.111111",
					Hours: "123123",
				},
				Spendable: Balance{
					Coins: "91.000001",
					Hours: "100023",
				},
				Expected: Balance{
					Coins: "137.111111",
					Hours: "100789",
				},
				Addresses: []AddressBalances{
					{
						Confirmed: Balance{
							Coins: "100.000000",
							Hours: "123000",
						},
						Spendable: Balance{
							Coins: "90.000000",
							Hours: "100000",
						},
						Expected: Balance{
							Coins: "90.000000",
							Hours: "100000",
						},
						Address: addrs[0],
					},
					{
						Confirmed: Balance{
							Coins: "0.000000",
							Hours: "0",
						},
						Spendable: Balance{
							Coins: "0.000000",
							Hours: "0",
						},
						Expected: Balance{
							Coins: "1.111111",
							Hours: "333",
						},
						Address: addrs[1],
					},
					{
						Confirmed: Balance{
							Coins: "23.111111",
							Hours: "123",
						},
						Spendable: Balance{
							Coins: "1.000001",
							Hours: "23",
						},
						Expected: Balance{
							Coins: "46.000000",
							Hours: "456",
						},
						Address: addrs[2],
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := getBalanceOfAddresses(&tc.outs, tc.addrs)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.result, result)
		})
	}
}
