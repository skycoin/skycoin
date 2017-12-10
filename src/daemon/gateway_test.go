package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"fmt"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/daemon/strand"
	"github.com/skycoin/skycoin/src/wallet"
)

func TestFbyAddresses(t *testing.T) {
	uxs := make(coin.UxArray, 5)
	addrs := make([]cipher.Address, 5)
	for i := 0; i < 5; i++ {
		addrs[i] = testutil.MakeAddress()
		uxs[i] = coin.UxOut{
			Body: coin.UxBody{
				Address: addrs[i],
			},
		}
	}

	tests := []struct {
		name    string
		addrs   []string
		outputs []coin.UxOut
		want    []coin.UxOut
	}{
		// TODO: Add test cases.
		{
			"filter with one address",
			[]string{addrs[0].String()},
			uxs[:2],
			uxs[:1],
		},
		{
			"filter with multiple addresses",
			[]string{addrs[0].String(), addrs[1].String()},
			uxs[:3],
			uxs[:2],
		},
	}
	for _, tt := range tests {
		// fmt.Printf("want:%+v\n", tt.want)
		outs := FbyAddresses(tt.addrs)(tt.outputs)
		require.Equal(t, outs, coin.UxArray(tt.want))
	}
}

func TestFbyHashes(t *testing.T) {
	uxs := make(coin.UxArray, 5)
	addrs := make([]cipher.Address, 5)
	for i := 0; i < 5; i++ {
		addrs[i] = testutil.MakeAddress()
		uxs[i] = coin.UxOut{
			Body: coin.UxBody{
				Address: addrs[i],
			},
		}
	}

	type args struct {
		hashes []string
	}
	tests := []struct {
		name    string
		hashes  []string
		outputs coin.UxArray
		want    coin.UxArray
	}{
		// TODO: Add test cases.
		{
			"filter with one hash",
			[]string{uxs[0].Hash().Hex()},
			uxs[:2],
			uxs[:1],
		},
		{
			"filter with multiple hash",
			[]string{uxs[0].Hash().Hex(), uxs[1].Hash().Hex()},
			uxs[:3],
			uxs[:2],
		},
	}
	for _, tt := range tests {
		outs := FbyHashes(tt.hashes)(tt.outputs)
		require.Equal(t, outs, coin.UxArray(tt.want))
	}
}

// Gateway RPC interface wrapper for daemon state
type FakeGateway struct {
	Config GatewayConfig
	drpc   RPC
	vrpc   visor.RPC

	// Backref to Daemon
	d *Daemon
	// Backref to Visor
	v *visor.Visor
	// Requests are queued on this channel
	requests chan strand.Request
}

// NewGateway create and init an Gateway instance.
func NewFakeGateway(c GatewayConfig, D *Daemon) *FakeGateway {
	return &FakeGateway{
		Config:   c,
		drpc:     RPC{},
		vrpc:     visor.MakeRPC(D.Visor.v),
		d:        D,
		v:        D.Visor.v,
		requests: make(chan strand.Request, c.BufferSize),
	}
}

func (gw *FakeGateway) Spend(wltID string, coins uint64, dest cipher.Address) (*coin.Transaction, error) {
	var tx *coin.Transaction
	var err error
	gw.strand("Spend", func() {
		// create spend validator
		unspent := gw.v.Blockchain.Unspent()
		sv := newSpendValidator(gw.v.Unconfirmed, unspent)
		// create and sign transaction
		tx, err = gw.vrpc.CreateAndSignTransaction(wltID, sv, unspent, gw.v.Blockchain.Time(), coins, dest)
		if err != nil {
			err = fmt.Errorf("Create transaction failed: %v", err)
			return
		}

		// inject transaction
		if err = gw.d.Visor.InjectTransaction(*tx, gw.d.Pool); err != nil {
			err = fmt.Errorf("Inject transaction failed: %v", err)
		}
	})

	return tx, err
}

// GetWalletBalance returns balance pair of specific wallet
func (gw *FakeGateway) GetWalletBalance(wltID string) (wallet.BalancePair, error) {
	var balance wallet.BalancePair
	var err error
	gw.strand("GetWalletBalance", func() {
		var addrs []cipher.Address
		addrs, err = gw.vrpc.GetWalletAddresses(wltID)
		if err != nil {
			return
		}
		auxs := gw.vrpc.GetUnspent(gw.v).GetUnspentsOfAddrs(addrs)

		var spendUxs coin.AddressUxOuts
		spendUxs, err = gw.vrpc.GetUnconfirmedSpends(gw.v, addrs)
		if err != nil {
			err = fmt.Errorf("get unconfimed spending failed when checking wallet balance: %v", err)
			return
		}

		var recvUxs coin.AddressUxOuts
		recvUxs, err = gw.vrpc.GetUnconfirmedReceiving(gw.v, addrs)
		if err != nil {
			err = fmt.Errorf("get unconfirmed receiving failed when when checking wallet balance: %v", err)
			return
		}

		coins1, hours1 := gw.v.AddressBalance(auxs)
		coins2, hours2 := gw.v.AddressBalance(auxs.Sub(spendUxs).Add(recvUxs))
		balance = wallet.BalancePair{
			Confirmed: wallet.Balance{Coins: coins1, Hours: hours1},
			Predicted: wallet.Balance{Coins: coins2, Hours: hours2},
		}
	})

	return balance, err
}
