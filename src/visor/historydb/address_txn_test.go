package historydb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
)

func TestAddAddressTxns(t *testing.T) {
	var preAddrs []cipher.Address
	var preTxHashes []cipher.SHA256

	type pair struct {
		addr   cipher.Address
		txHash cipher.SHA256
	}

	type expectPair struct {
		addr cipher.Address
		txs  []cipher.SHA256
	}

	for i := 0; i < 3; i++ {
		preAddrs = append(preAddrs, makeAddress())
		preTxHashes = append(preTxHashes, cipher.SumSHA256([]byte(fmt.Sprintf("tx%d", i))))
	}

	var testCases = []struct {
		name     string
		addPairs []pair
		expect   []expectPair
	}{
		{
			"address with single tx",
			[]pair{
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[0],
				},
				{
					addr:   preAddrs[1],
					txHash: preTxHashes[1],
				},
			},
			[]expectPair{
				{
					preAddrs[0],
					preTxHashes[:1],
				},
				{
					preAddrs[1],
					preTxHashes[1:2],
				},
			},
		},
		{
			"address with multiple transactions",
			[]pair{
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[0],
				},
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[1],
				},
			},
			[]expectPair{
				{
					preAddrs[0],
					preTxHashes[:2],
				},
			},
		},
		{
			"add address with multiple same transactions",
			[]pair{
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[0],
				},
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[0],
				},
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[0],
				},
			},
			[]expectPair{
				{
					preAddrs[0],
					preTxHashes[:1],
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, td := prepareDB(t)
			defer td()

			addrTxns := &addressTxns{}

			err := db.Update("", func(tx *dbutil.Tx) error {
				for _, pr := range tc.addPairs {
					err := addrTxns.add(tx, pr.addr, pr.txHash)
					require.NoError(t, err)
				}
				return nil
			})
			require.NoError(t, err)

			for _, e := range tc.expect {
				err := db.View("", func(tx *dbutil.Tx) error {
					hashes, err := addrTxns.get(tx, e.addr)
					require.NoError(t, err)
					require.Equal(t, e.txs, hashes)
					return nil
				})
				require.NoError(t, err)
			}

		})
	}
}

func TestGetAddressTxns(t *testing.T) {
	var preAddrs []cipher.Address
	var preTxHashes []cipher.SHA256

	type pair struct {
		addr   cipher.Address
		txHash cipher.SHA256
	}

	type expectPair struct {
		addr cipher.Address
		txs  []cipher.SHA256
	}

	for i := 0; i < 3; i++ {
		preAddrs = append(preAddrs, makeAddress())
		preTxHashes = append(preTxHashes, cipher.SumSHA256([]byte(fmt.Sprintf("tx%d", i))))
	}

	var testCases = []struct {
		name     string
		addPairs []pair
		expect   []expectPair
	}{
		{
			"address with single tx",
			[]pair{
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[0],
				},
				{
					addr:   preAddrs[1],
					txHash: preTxHashes[1],
				},
			},
			[]expectPair{
				{
					preAddrs[0],
					preTxHashes[:1],
				},
				{
					preAddrs[1],
					preTxHashes[1:2],
				},
			},
		},
		{
			"address with multiple transactions",
			[]pair{
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[0],
				},
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[1],
				},
			},
			[]expectPair{
				{
					preAddrs[0],
					preTxHashes[:2],
				},
			},
		},
		{
			"add address with multiple same transactions",
			[]pair{
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[0],
				},
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[0],
				},
				{
					addr:   preAddrs[0],
					txHash: preTxHashes[0],
				},
			},
			[]expectPair{
				{
					preAddrs[0],
					preTxHashes[:1],
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, td := prepareDB(t)
			defer td()

			addrTxns := &addressTxns{}

			err := db.Update("", func(tx *dbutil.Tx) error {
				for _, pr := range tc.addPairs {
					err := addrTxns.add(tx, pr.addr, pr.txHash)
					require.NoError(t, err)
				}

				return nil
			})
			require.NoError(t, err)

			err = db.View("", func(tx *dbutil.Tx) error {
				for _, e := range tc.expect {
					hashes, err := addrTxns.get(tx, e.addr)
					require.NoError(t, err)
					require.Equal(t, e.txs, hashes)
				}
				return nil
			})
			require.NoError(t, err)

		})
	}
}
