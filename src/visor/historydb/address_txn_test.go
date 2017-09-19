package historydb

import (
	"fmt"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/stretchr/testify/require"
)

func TestNewAddressTxns(t *testing.T) {
	db, td := testutil.PrepareDB(t)
	defer td()

	_, err := newAddressTxnsBkt(db)
	require.Nil(t, err)

	// the address_txns bucket must be exist
	db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("address_txns"))
		require.NotNil(t, bkt)
		return nil
	})
}

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
			db, td := testutil.PrepareDB(t)
			defer td()

			_, err := newAddressTxnsBkt(db)
			require.Nil(t, err)

			require.Nil(t, db.Update(func(tx *bolt.Tx) error {
				bkt := tx.Bucket(addressTxnsBktName)
				for _, pr := range tc.addPairs {
					require.Nil(t, setAddressTxns(bkt, pr.addr, pr.txHash))
				}
				return nil
			}))

			for _, e := range tc.expect {
				db.View(func(tx *bolt.Tx) error {
					bkt := tx.Bucket(addressTxnsBktName)
					v := bkt.Get(e.addr.Bytes())
					require.NotNil(t, v)
					var hashes []cipher.SHA256
					require.Nil(t, encoder.DeserializeRaw(v, &hashes))
					require.Equal(t, e.txs, hashes)
					return nil
				})
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
			db, td := testutil.PrepareDB(t)
			defer td()

			addrTxnsBkt, err := newAddressTxnsBkt(db)
			require.Nil(t, err)

			require.Nil(t, db.Update(func(tx *bolt.Tx) error {
				bkt := tx.Bucket(addressTxnsBktName)

				for _, pr := range tc.addPairs {
					if err := setAddressTxns(bkt, pr.addr, pr.txHash); err != nil {
						return err
					}
				}

				return nil
			}))

			for _, e := range tc.expect {
				hashes, err := addrTxnsBkt.Get(e.addr)
				require.Nil(t, err)
				require.Equal(t, e.txs, hashes)
			}

		})
	}
}
