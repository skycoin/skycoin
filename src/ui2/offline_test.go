package ui2

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
)

func (f *fakeNode) OutputsForHashes(hashes []string) (*readable.UnspentOutputsSummary, error) {
	out := &readable.UnspentOutputsSummary{}
	for _, h := range hashes {
		if f.knownOutputs[h] {
			out.HeadOutputs = append(out.HeadOutputs, readable.UnspentOutput{Hash: h})
		}
	}
	return out, nil
}

// VerifyTransaction decodes the transaction the way the node would, giving
// input i the address f.inputOwners[i].
func (f *fakeNode) VerifyTransaction(r api.VerifyTransactionRequest) (*api.VerifyTransactionResponse, error) {
	txn, err := coin.DeserializeTransactionHexWithPQ(r.EncodedTransaction)
	if err != nil {
		return nil, err
	}
	out := &api.VerifyTransactionResponse{Unsigned: r.Unsigned}
	out.Transaction.TxID = txn.Hash().Hex()
	out.Transaction.Fee = "5"
	for i, h := range txn.In {
		out.Transaction.In = append(out.Transaction.In, api.CreatedTransactionInput{UxID: h.Hex(), Address: f.inputOwners[i], Coins: "1.000000", CalculatedHours: "10"})
	}
	for _, o := range txn.Out {
		out.Transaction.Out = append(out.Transaction.Out, api.CreatedTransactionOutput{Address: o.Address.String(), Coins: "1.000000", Hours: "2"})
	}
	return out, nil
}

// WalletSignTransaction marks the asked inputs signed.
func (f *fakeNode) WalletSignTransaction(r api.WalletSignTransactionRequest) (*api.CreateTransactionResponse, error) {
	f.signReqs = append(f.signReqs, r)
	txn, err := coin.DeserializeTransactionHexWithPQ(r.EncodedTransaction)
	if err != nil {
		return nil, err
	}
	for _, i := range r.SignIndexes {
		txn.Sigs[i] = cipher.MustSigFromHex(strings.Repeat("11", 65))
	}
	raw, err := txn.SerializeHex()
	if err != nil {
		return nil, err
	}
	return &api.CreateTransactionResponse{EncodedTransaction: raw}, nil
}

// rawTx is an unsigned transaction spending n inputs, with input i signed
// when signed[i].
func rawTx(t *testing.T, signed ...bool) string {
	t.Helper()
	var txn coin.Transaction
	for i := range signed {
		require.NoError(t, txn.PushInput(cipher.SumSHA256([]byte{byte(i)})))
	}
	require.NoError(t, txn.PushOutput(cipher.MustDecodeBase58Address(newAddr()), 1e6, 2))
	txn.Sigs = make([]cipher.Sig, len(signed))
	for i, s := range signed {
		if s {
			txn.Sigs[i] = cipher.MustSigFromHex(strings.Repeat("22", 65))
		}
	}
	require.NoError(t, txn.UpdateHeader())
	raw, err := txn.SerializeHex()
	require.NoError(t, err)
	return raw
}

func TestReviewAndSign(t *testing.T) {
	n := testNode()
	n.inputOwners = []string{"addrRecv0", "stranger"}
	raw := rawTx(t, false, false)

	v := ReviewTx(n, "a.wlt", raw)
	require.Empty(t, v.Err)
	require.Equal(t, 0, v.Signed)
	require.Equal(t, 1, v.ToSign)
	require.True(t, v.Encrypted)
	require.True(t, v.Inputs[0].Mine)
	require.False(t, v.Inputs[1].Mine)
	require.Equal(t, "10", v.MinFee, "20 input hours at the default burn factor of 2")

	h := Handler(n, Options{})
	body := postForm(h, "/tx/review", url.Values{"id": {"a.wlt"}, "rawtx": {raw}}).Body.String()
	require.Contains(t, body, `hx-post="tx/sign"`)
	require.Contains(t, body, "0 of 2 inputs signed")
	require.NotContains(t, body, `hx-post="tx/broadcast"`)

	body = postForm(h, "/tx/sign", url.Values{"id": {"a.wlt"}, "rawtx": {raw}, "password": {"pw"}}).Body.String()
	require.Equal(t, []int{0}, n.signReqs[0].SignIndexes, "only the wallet's own unsigned input")
	require.Equal(t, "pw", n.signReqs[0].Password)
	require.Contains(t, body, "1 of 2 inputs signed")
	require.Contains(t, body, "already signed its inputs")

	full := rawTx(t, true, true)
	body = postForm(h, "/tx/review", url.Values{"rawtx": {full}}).Body.String()
	require.Contains(t, body, `hx-post="tx/broadcast"`)
	require.NotContains(t, body, `hx-post="tx/sign"`)
	postForm(h, "/tx/broadcast", url.Values{"rawtx": {" " + full + " "}})
	require.Equal(t, full, n.injected[len(n.injected)-1])

	require.Contains(t, postForm(h, "/tx/review", url.Values{"rawtx": {"zz"}}).Body.String(), "not a valid transaction")
	require.Contains(t, postForm(h, "/tx/sign", url.Values{"id": {"a.wlt"}, "rawtx": {rawTx(t, true, false)}}).Body.String(), "owns none of the unsigned inputs")

	body = get(t, h, "/tx?id=a.wlt&tx="+raw).Body.String()
	require.Contains(t, body, raw)
	require.Contains(t, body, `<option value="a.wlt" selected>`)
}

func TestJoint(t *testing.T) {
	n := testNode()
	n.knownOutputs = map[string]bool{"theirs1": true}
	h := Handler(n, Options{})
	require.Contains(t, get(t, h, "/joint?id=a.wlt").Body.String(), `hx-get="joint/outputs"`)
	require.Contains(t, get(t, h, "/joint/outputs?id=a.wlt").Body.String(), `value="ux1"`)

	to := newAddr()
	v := url.Values{"id": {"a.wlt"}, "mine": {"ux1"}, "theirs": {"THEIRS1, "}, "to": {to, ""}, "coins": {"1.5", ""}, "hours": {"3", ""}}
	body := postForm(h, "/joint/create", v).Body.String()
	require.Contains(t, body, "JOINTRAW")
	r := n.created2[len(n.created2)-1]
	require.Equal(t, []string{"ux1", "theirs1"}, r.UxOuts)
	require.Equal(t, "manual", r.HoursSelection.Type)
	require.Equal(t, []api.Receiver{{Address: to, Coins: "1.5", Hours: "3"}}, r.To)

	v.Set("theirs", "unknown")
	require.Contains(t, postForm(h, "/joint/create", v).Body.String(), "not unspent outputs on this node")
	v.Set("theirs", "dropped")
	n.knownOutputs["dropped"] = true
	require.Contains(t, postForm(h, "/joint/create", v).Body.String(), "left out some of the outputs")
	v.Set("hours", "x")
	require.Contains(t, postForm(h, "/joint/create", v).Body.String(), "whole number of hours")
	require.Contains(t, postForm(h, "/joint/create", url.Values{"id": {"a.wlt"}}).Body.String(), "at least one output")
	require.Equal(t, http.StatusOK, get(t, h, "/joint/row").Code)
}

func TestUnsignedSend(t *testing.T) {
	x := sendNode(t)
	h := Handler(x.n, Options{})
	body := postForm(h, "/send/unsigned", url.Values{"id": {"s.wlt"}, "to": {x.toFresh}, "coins": {"1"}}).Body.String()
	require.True(t, x.n.txReqs[len(x.n.txReqs)-1].Unsigned)
	require.Contains(t, body, `href="tx?id=s.wlt&amp;tx=RAWTX"`)
	require.Contains(t, get(t, h, "/send?id=s.wlt").Body.String(), `hx-post="send/unsigned"`)
}

func TestBrowserTxPage(t *testing.T) {
	x := sendNode(t)
	x.n.browser = true
	x.n.err = api.NewClientError("403 Forbidden", http.StatusForbidden, "Endpoint is disabled")
	body := get(t, Handler(x.n, Options{}), "/tx").Body.String()
	require.Contains(t, body, `hx-post="tx/review"`, "a browser wallet node still reviews and broadcasts")
	require.NotContains(t, body, `name="id"`)
	require.NotContains(t, body, `href="joint"`)
}
