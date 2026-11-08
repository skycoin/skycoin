package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
)

type fakeHashOutputs struct {
	outs readable.UnspentOutputs
}

func (f fakeHashOutputs) OutputsForHashes(_ []string) (*readable.UnspentOutputsSummary, error) {
	return &readable.UnspentOutputsSummary{HeadOutputs: f.outs}, nil
}

// jointFixture is a transaction with one input from each of two keys
type jointFixture struct {
	txn      coin.Transaction
	keys     []cipher.SecKey
	addrs    []cipher.Address
	lookup   fakeHashOutputs
	uxids    []string
	outAddrs []cipher.Address
}

func newJointFixture(t *testing.T) jointFixture {
	var f jointFixture
	hours := []uint64{1112, 5}
	for i := 0; i < 2; i++ {
		pk, sk := cipher.GenerateKeyPair()
		addr := cipher.AddressFromPubKey(pk)
		ux := testutil.RandSHA256(t)
		f.keys = append(f.keys, sk)
		f.addrs = append(f.addrs, addr)
		f.uxids = append(f.uxids, ux.Hex())
		f.lookup.outs = append(f.lookup.outs, readable.UnspentOutput{
			Hash: ux.Hex(), Address: addr.String(), Coins: "0.001000", CalculatedHours: hours[i],
		})
		require.NoError(t, f.txn.PushInput(ux))
	}
	f.outAddrs = []cipher.Address{f.addrs[0], f.addrs[1]}
	require.NoError(t, f.txn.PushOutput(f.addrs[0], 1000, 0))
	require.NoError(t, f.txn.PushOutput(f.addrs[1], 1000, 1000))
	f.txn.Sigs = make([]cipher.Sig, 2)
	require.NoError(t, f.txn.UpdateHeader())
	return f
}

func (f jointFixture) sign(i int) {
	h := cipher.AddSHA256(f.txn.InnerHash, f.txn.In[i])
	f.txn.Sigs[i] = cipher.MustSignHash(h, f.keys[i])
}

func TestParseJointOutputs(t *testing.T) {
	addr := testutil.MakeAddress().String()
	cases := []struct {
		name  string
		specs []string
		want  []api.Receiver
		err   string
	}{
		{name: "two outputs", specs: []string{addr + ",0.001,0", " " + addr + " , 1 , 1000 "},
			want: []api.Receiver{{Address: addr, Coins: "0.001", Hours: "0"}, {Address: addr, Coins: "1", Hours: "1000"}}},
		{name: "missing hours", specs: []string{addr + ",1"}, err: "want address,coins,hours"},
		{name: "bad address", specs: []string{"nope,1,1"}, err: "invalid address"},
		{name: "bad coins", specs: []string{addr + ",x,1"}, err: "invalid coins"},
		{name: "zero coins", specs: []string{addr + ",0,1"}, err: "must not be zero"},
		{name: "too precise", specs: []string{addr + ",0.0000001,1"}, err: "invalid coins"},
		{name: "below divisor", specs: []string{addr + ",0.0001,1"}, err: "too many decimal places"},
		{name: "bad hours", specs: []string{addr + ",1,-3"}, err: "invalid hours"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseJointOutputs(tc.specs)
			if tc.err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestMakeJointTxnRequest(t *testing.T) {
	f := newJointFixture(t)
	change := testutil.MakeAddress().String()
	out := func(coins string) []api.Receiver {
		return []api.Receiver{{Address: f.addrs[1].String(), Coins: coins, Hours: "1000"}}
	}

	req, err := makeJointTxnRequest(f.lookup, f.uxids, out("0.002"), "")
	require.NoError(t, err)
	require.Equal(t, "manual", req.HoursSelection.Type)
	require.Equal(t, f.uxids, req.UxOuts)
	require.Nil(t, req.ChangeAddress)

	_, err = makeJointTxnRequest(f.lookup, f.uxids, out("0.001"), "")
	require.EqualError(t, err, "outputs leave 0.001000 coins unspent, set --change-address or add them to an output")

	req, err = makeJointTxnRequest(f.lookup, f.uxids, out("0.001"), change)
	require.NoError(t, err)
	require.Equal(t, change, *req.ChangeAddress)

	_, err = makeJointTxnRequest(f.lookup, f.uxids, out("0.003"), "")
	require.EqualError(t, err, "outputs spend 0.003000 coins but the inputs hold 0.002000")

	missing := testutil.RandSHA256(t).Hex()
	_, err = makeJointTxnRequest(f.lookup, []string{missing}, out("0.001"), "")
	require.EqualError(t, err, "unspent output "+missing+" does not exist")

	_, err = makeJointTxnRequest(f.lookup, []string{"zz"}, out("0.001"), "")
	require.Error(t, err)
	_, err = makeJointTxnRequest(f.lookup, nil, out("0.001"), "")
	require.Error(t, err)
}

func TestCheckAllInputsSpent(t *testing.T) {
	got := []api.CreatedTransactionInput{{UxID: "a"}, {UxID: "b"}}
	require.NoError(t, checkAllInputsSpent([]string{"b", "a"}, got))
	err := checkAllInputsSpent([]string{"a", "c"}, got)
	require.Error(t, err)
	require.Contains(t, err.Error(), "left them out: c")
}

func TestWalletSignIndexes(t *testing.T) {
	f := newJointFixture(t)
	wallet := func(addrs ...cipher.Address) *api.WalletResponse {
		w := &api.WalletResponse{}
		for _, a := range addrs {
			w.Accounts = append(w.Accounts, readable.WalletAccount{ChangeEntries: []readable.WalletEntry{{Address: a.String()}}})
		}
		w.Entries = []readable.WalletEntry{{Address: testutil.MakeAddress().String()}}
		return w
	}

	idx, err := walletSignIndexes(f.lookup, &f.txn, wallet(f.addrs[1]))
	require.NoError(t, err)
	require.Equal(t, []int{1}, idx)

	idx, err = walletSignIndexes(f.lookup, &f.txn, wallet(f.addrs[0]))
	require.NoError(t, err)
	require.Equal(t, []int{0}, idx)

	// Owning every unsigned input keeps the old request, with no sign_indexes
	idx, err = walletSignIndexes(f.lookup, &f.txn, wallet(f.addrs[0], f.addrs[1]))
	require.NoError(t, err)
	require.Nil(t, idx)

	f.sign(0)
	idx, err = walletSignIndexes(f.lookup, &f.txn, wallet(f.addrs[1]))
	require.NoError(t, err)
	require.Nil(t, idx)

	_, err = walletSignIndexes(f.lookup, &f.txn, wallet(f.addrs[0]))
	require.EqualError(t, err, "the wallet owns none of the unsigned inputs")
}

func TestTxnReview(t *testing.T) {
	f := newJointFixture(t)
	inputs, err := lookupInputs(f.lookup, f.uxids)
	require.NoError(t, err)
	mine := map[string]struct{}{f.addrs[0].String(): {}}

	s := txnReview(&f.txn, inputs, mine)
	require.Contains(t, s, "[0] "+f.addrs[0].String()+"  0.001000 coins  1112 hours  unsigned  YOURS")
	require.Contains(t, s, "[1] "+f.addrs[1].String()+"  0.001000 coins  5 hours  unsigned  other")
	require.Contains(t, s, "[1] "+f.addrs[1].String()+"  0.001000 coins  1000 hours  other")
	require.Contains(t, s, "Fee: 117 hours burned (minimum 112)")
	require.Contains(t, s, "Signed: 0 of 2 inputs\n")

	f.sign(0)
	f.sign(1)
	s = txnReview(&f.txn, inputs, nil)
	require.Contains(t, s, "1112 hours  signed\n")
	require.Contains(t, s, "Signed: 2 of 2 inputs, ready to broadcast")

	inputs[1].Known = false
	require.Contains(t, txnReview(&f.txn, inputs, nil), "Fee: unknown")
}

func TestFeeNote(t *testing.T) {
	txn := api.CreatedTransaction{Fee: "117", In: []api.CreatedTransactionInput{{CalculatedHours: "1112"}, {CalculatedHours: "5"}}}
	require.Equal(t, " (minimum 112, give the other 5 to an output or they are burned)", feeNote(txn))
	txn.Fee = "112"
	require.Equal(t, "", feeNote(txn))
}
