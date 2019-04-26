package daemon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/util/useragent"
	"github.com/skycoin/skycoin/src/visor"
)

func TestDivideHashes(t *testing.T) {
	hashes := make([]cipher.SHA256, 10)
	for i := 0; i < 10; i++ {
		hashes[i] = cipher.SumSHA256(cipher.RandByte(512))
	}

	testCases := []struct {
		name  string
		init  []cipher.SHA256
		n     int
		array [][]cipher.SHA256
	}{
		{
			"has one odd",
			hashes[:],
			3,
			[][]cipher.SHA256{
				[]cipher.SHA256{
					hashes[0],
					hashes[1],
					hashes[2],
				},
				[]cipher.SHA256{
					hashes[3],
					hashes[4],
					hashes[5],
				},
				[]cipher.SHA256{
					hashes[6],
					hashes[7],
					hashes[8],
				},
				[]cipher.SHA256{
					hashes[9],
				},
			},
		},
		{
			"only one value",
			hashes[:1],
			1,
			[][]cipher.SHA256{
				[]cipher.SHA256{
					hashes[0],
				},
			},
		},
		{
			"empty value",
			hashes[:0],
			0,
			[][]cipher.SHA256{},
		},
		{
			"with 3 value",
			hashes[:3],
			3,
			[][]cipher.SHA256{
				[]cipher.SHA256{
					hashes[0],
					hashes[1],
					hashes[2],
				},
			},
		},
		{
			"with 8 value",
			hashes[:8],
			3,
			[][]cipher.SHA256{
				[]cipher.SHA256{
					hashes[0],
					hashes[1],
					hashes[2],
				},
				[]cipher.SHA256{
					hashes[3],
					hashes[4],
					hashes[5],
				},
				[]cipher.SHA256{
					hashes[6],
					hashes[7],
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rlt := divideHashes(tc.init, tc.n)
			require.Equal(t, tc.array, rlt)
		})
	}
}

func TestVerifyUserTxnAgainstPeer(t *testing.T) {
	now := uint64(time.Now().UTC().Unix())

	cases := []struct {
		name         string
		err          error
		txn          coin.Transaction
		head         *coin.SignedBlock
		inputs       coin.UxArray
		verifyParams params.VerifyTxn
	}{
		{
			name: "invalid droplet precision",
			err:  params.ErrInvalidDecimals,
			txn: coin.Transaction{
				Out: []coin.TransactionOutput{
					{
						Coins: 111100,
					},
				},
			},
			verifyParams: params.VerifyTxn{
				MaxDropletPrecision: 3,
			},
		},

		{
			name: "invalid txn size",
			err:  visor.ErrTxnExceedsMaxBlockSize,
			txn: coin.Transaction{
				Out: []coin.TransactionOutput{
					{
						Coins: 1e6,
					},
				},
			},
			verifyParams: params.VerifyTxn{
				MaxDropletPrecision: 3,
				MaxTransactionSize:  1,
			},
		},

		{
			name: "invalid burn fee",
			err:  fee.ErrTxnInsufficientFee,
			txn: coin.Transaction{
				Out: []coin.TransactionOutput{
					{
						Coins: 1e6,
						Hours: 100,
					},
				},
			},
			head: &coin.SignedBlock{
				Block: coin.Block{
					Head: coin.BlockHeader{
						Time: now,
					},
				},
			},
			inputs: coin.UxArray{
				{
					Head: coin.UxHead{
						Time: now,
					},
					Body: coin.UxBody{
						Coins: 1e6,
						Hours: 150,
					},
				},
			},
			verifyParams: params.VerifyTxn{
				MaxDropletPrecision: 3,
				MaxTransactionSize:  99999999,
				BurnFactor:          2,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := verifyUserTxnAgainstPeer(tc.txn, tc.head, tc.inputs, tc.verifyParams)
			require.Equal(t, tc.err, err)
		})
	}
}

func TestCheckBroadcastTxnRecipients(t *testing.T) {
	// contains a connection not introduced and a connection without user agent
	connections := NewConnections()

	// one connection connected but not introduced
	_, err := connections.connected("1.1.1.1:9999", 3)
	require.NoError(t, err)

	// one connection introduced without user agent
	_, err = connections.connected("2.2.2.2:9999", 1)
	require.NoError(t, err)
	_, err = connections.introduced("2.2.2.2:9999", 1, &IntroductionMessage{
		Mirror:          6666,
		ListenPort:      9999,
		ProtocolVersion: 2,
	})
	require.NoError(t, err)

	// contains a connection not introduced, a connection without user agent
	// and a connection with user agent
	connections2 := NewConnections()

	// one connection connected but not introduced
	_, err = connections2.connected("1.1.1.1:9999", 3)
	require.NoError(t, err)

	// one connection introduced without user agent
	_, err = connections2.connected("2.2.2.2:9999", 1)
	require.NoError(t, err)
	_, err = connections2.introduced("2.2.2.2:9999", 1, &IntroductionMessage{
		Mirror:          6666,
		ListenPort:      9999,
		ProtocolVersion: 2,
	})
	require.NoError(t, err)

	// one connection introduced with user agent
	_, err = connections2.connected("3.3.3.3:9999", 2)
	require.NoError(t, err)
	_, err = connections2.introduced("3.3.3.3:9999", 2, &IntroductionMessage{
		Mirror:          7777,
		ListenPort:      9999,
		ProtocolVersion: 2,
		UserAgent:       useragent.MustParse("skycoin:0.25.1"),
		UnconfirmedVerifyTxn: params.VerifyTxn{
			MaxDropletPrecision: 4, // the default precision for unspecified peers is 3
			MaxTransactionSize:  32768,
			BurnFactor:          2,
		},
	})
	require.NoError(t, err)

	cases := []struct {
		name        string
		err         error
		accepts     int
		connections *Connections
		ids         []uint64
		txn         coin.Transaction
		head        *coin.SignedBlock
		inputs      coin.UxArray
	}{
		{
			name:        "accepted by connection with no user agent",
			err:         nil,
			accepts:     1,
			connections: connections,
			ids:         []uint64{1, 2, 999}, // includes unknown gnet id to make sure it doesn't crash on bad input
			txn: coin.Transaction{
				Out: []coin.TransactionOutput{
					{
						Coins: 1e6,
					},
				},
			},
			head: &coin.SignedBlock{},
			inputs: coin.UxArray{
				{
					Body: coin.UxBody{
						Coins: 1e6,
						Hours: 100,
					},
				},
			},
		},

		{
			name:        "not accepted by connection with no user agent",
			err:         ErrNoPeerAcceptsTxn,
			accepts:     0,
			connections: connections,
			ids:         []uint64{1, 2},
			txn: coin.Transaction{
				Out: []coin.TransactionOutput{
					{
						Coins: 1e2,
					},
				},
			},
			head: &coin.SignedBlock{},
			inputs: coin.UxArray{
				{
					Body: coin.UxBody{
						Coins: 1e2,
						Hours: 100,
					},
				},
			},
		},

		{
			name:        "connections contains a connection that specifies user agent, which will propagate a txn even if soft-invalid",
			err:         nil,
			accepts:     1,
			connections: connections2,
			ids:         []uint64{1, 2, 3},
			txn: coin.Transaction{
				Out: []coin.TransactionOutput{
					{
						Coins: 1e1,
					},
				},
			},
			head: &coin.SignedBlock{},
			inputs: coin.UxArray{
				{
					Body: coin.UxBody{
						Coins: 1e1,
						Hours: 100,
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			accepts, err := checkBroadcastTxnRecipients(tc.connections, tc.ids, tc.txn, tc.head, tc.inputs)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.accepts, accepts)
		})
	}
}
