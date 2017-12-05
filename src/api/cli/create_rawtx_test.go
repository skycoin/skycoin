package cli

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor/historydb"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testWebrpcAddr = "127.0.0.1:8081"

type fakeGateway struct{}

func (fg fakeGateway) GetUnspentOutputs(filters ...daemon.OutputsFilter) (visor.ReadableOutputSet, error) {

	outs := []coin.UxOut{
		coin.UxOut{
			Head: coin.UxHead{
				Time:  0,
				BkSeq: 0,
			},
			Body: coin.UxBody{
				SrcTransaction: cipher.SHA256{},
				Address:        cipher.Address{},
				Coins:          500e6,
				Hours:          100,
			},
		},
	}
	rbOuts, err := visor.NewReadableOutputs(outs)
	if err != nil {
		return visor.ReadableOutputSet{}, err
	}

	return visor.ReadableOutputSet{
		HeadOutputs: rbOuts,
	}, nil
}

func (fg fakeGateway) GetAddrUxOuts(addr cipher.Address) ([]*historydb.UxOutJSON, error) {
	return nil, nil
}

func (fg fakeGateway) GetBlocks(start, end uint64) (*visor.ReadableBlocks, error) {
	return nil, nil
}

func (fg fakeGateway) GetBlocksInDepth(vs []uint64) (*visor.ReadableBlocks, error) {
	return nil, nil
}

func (fg fakeGateway) GetLastBlocks(num uint64) (*visor.ReadableBlocks, error) {
	return nil, nil
}

func (fg fakeGateway) GetTransaction(txid cipher.SHA256) (*visor.Transaction, error) {
	return nil, nil
}

func (fg fakeGateway) InjectTransaction(tx coin.Transaction) error {
	return nil
}

func (fg fakeGateway) GetTimeNow() uint64 {
	return 0
}

func setupWebRPC(t *testing.T) *webrpc.WebRPC {
	rpc, err := webrpc.New(testWebrpcAddr, &fakeGateway{})
	require.NoError(t, err)
	rpc.WorkerNum = 1
	rpc.ChanBuffSize = 2
	return rpc
}

func TestCreateRawTx(t *testing.T) {
	s := setupWebRPC(t)
	c := webrpc.Client{
		Addr: s.Addr,
	}

	go func() {
		err := s.Run()
		require.NoError(t, err)
	}()

	defer func() {
		err := s.Shutdown()
		require.NoError(t, err)
	}()

	tests := []struct {
		name    string
		inAddrs []string
		chgAddr string
		toAddrs []SendAmount
		wlt     wallet.Wallet

		err      error
		expected string
	}{
		{
			name:    "invalid address",
			inAddrs: []string{"foo-bar-buzz"},
			err:     errors.New("invalid address: foo-bar-buzz [code: -32602]"),
		},
		{
			name:    "insufficient balance",
			inAddrs: []string{"2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"},
			chgAddr: "k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND",
			toAddrs: []SendAmount{
				SendAmount{
					Addr:  "A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe",
					Coins: 100e6,
				},
				SendAmount{
					Addr:  "2PBmUva7J8WFsyWg979cREZkU3z2pkYjNkE",
					Coins: 500e6,
				},
			},
			err: errors.New("balance in wallet is not sufficient"),
		},
		{
			name:    "address not in wallet",
			inAddrs: []string{"2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"},
			chgAddr: "k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND",
			toAddrs: []SendAmount{
				SendAmount{
					Addr:  "A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe",
					Coins: 100e6,
				},
				SendAmount{
					Addr:  "2PBmUva7J8WFsyWg979cREZkU3z2pkYjNkE",
					Coins: 100e6,
				},
			},
			wlt: wallet.Wallet{},
			err: errors.New("is not in wallet"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := CreateRawTx(&c, &tc.wlt, tc.inAddrs, tc.chgAddr, tc.toAddrs)
			if tc.name == "address not in wallet" {
				require.Contains(t, err.Error(), tc.err.Error())
			} else if tc.err != nil {
				require.Equal(t, err.Error(), tc.err.Error())
			}
		})
	}
}

func TestNewTransaction(t *testing.T) {
	_, sk := cipher.GenerateKeyPair()
	addr := cipher.AddressFromSecKey(sk)

	utxos := []UnspentOut{
		UnspentOut{
			visor.ReadableOutput{
				Hash:              "c8b8eac053a5640bae40144cbc3dda02746071e3c7d00a4b5dfd06d28f928ec4",
				SourceTransaction: "b3c6f0f87c5282ff7ff5e6d637c2581e6a56826a76ec3dd221d02786881e3d14",
				Address:           addr.String(),
				Coins:             "2500",
				Hours:             800291,
			},
		},
	}

	outs := []coin.TransactionOutput{
		coin.TransactionOutput{
			Address: addr,
			Coins:   2500e6,
			Hours:   400145,
		},
	}

	tx, err := NewTransaction(utxos, []cipher.SecKey{sk}, outs)
	require.NoError(t, err)
	require.NoError(t, tx.Verify())
}

func makeReadableOutput(addr, coins string, hours uint64) visor.ReadableOutput {
	return visor.ReadableOutput{
		Address: addr,
		Coins:   coins,
		Hours:   hours,
	}
}

func TestGetSufficientUnspents(t *testing.T) {
	tests := []struct {
		Name       string
		Unspents   webrpc.OutputsResult
		BeingSpent uint64

		Expected []UnspentOut
		Err      error
	}{
		{
			Name: "sufficient HeadOutputs",
			Unspents: webrpc.OutputsResult{
				Outputs: visor.ReadableOutputSet{
					HeadOutputs: []visor.ReadableOutput{
						makeReadableOutput("k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND", "200", 0),
						makeReadableOutput("A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe", "400", 0),
					},
				},
			},
			BeingSpent: 300e6,
			Expected: []UnspentOut{
				UnspentOut{
					makeReadableOutput("k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND", "200", 0),
				},
				UnspentOut{
					makeReadableOutput("A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe", "400", 0),
				},
			},
		},
		{
			Name: "insufficient HeadOutputs",
			Unspents: webrpc.OutputsResult{
				Outputs: visor.ReadableOutputSet{
					HeadOutputs: []visor.ReadableOutput{
						makeReadableOutput("k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND", "200", 0),
						makeReadableOutput("A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe", "400", 0),
					},
				},
			},
			BeingSpent: 700e6,
			Expected:   nil,
			Err:        ErrInsufficientBalance,
		},
		{
			Name: "temporary insufficient HeadOutputs + UncfmIncomingOutputs",
			Unspents: webrpc.OutputsResult{
				Outputs: visor.ReadableOutputSet{
					HeadOutputs: []visor.ReadableOutput{
						makeReadableOutput("k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND", "200", 0),
						makeReadableOutput("A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe", "400", 0),
					},
					IncomingOutputs: []visor.ReadableOutput{
						makeReadableOutput("PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m", "200", 0),
					},
				},
			},
			BeingSpent: 700e6,
			Expected:   nil,
			Err:        ErrTemporaryInsufficientBalance,
		},
		{
			Name: "insufficient HeadOutputs + UncfmIncomingOutputs",
			Unspents: webrpc.OutputsResult{
				Outputs: visor.ReadableOutputSet{
					HeadOutputs: []visor.ReadableOutput{
						makeReadableOutput("k3rmz3PGbTxd7KL8AL5CeHrWy35C1UcWND", "200", 0),
						makeReadableOutput("A2h4iWC1SDGmS6UPezatFzEUwirLJtjFUe", "400", 0),
					},
					IncomingOutputs: []visor.ReadableOutput{
						makeReadableOutput("PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m", "200", 0),
					},
				},
			},
			BeingSpent: 900e6,
			Expected:   nil,
			Err:        ErrInsufficientBalance,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			unspents, err := getSufficientUnspents(&tc.Unspents, tc.BeingSpent)
			if err != nil {
				require.Equal(t, err, tc.Err)
			} else {
				assert.True(t, reflect.DeepEqual(unspents, tc.Expected))
			}
		})
	}

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
