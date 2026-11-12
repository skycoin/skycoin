package ui2

import (
	"encoding/json"
	"html"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
)

func (f *fakeNode) CreateTransaction(r api.CreateTransactionRequest) (*api.CreateTransactionResponse, error) {
	f.created2 = append(f.created2, r)
	var t api.CreateTransactionResponse
	if len(r.UxOuts) > 0 && r.ChangeAddress == nil {
		// A joint transaction: every named output is spent but "dropped".
		for _, u := range r.UxOuts {
			if u != "dropped" {
				t.Transaction.In = append(t.Transaction.In, api.CreatedTransactionInput{UxID: u})
			}
		}
		t.EncodedTransaction = "JOINTRAW"
		return &t, nil
	}
	t.Transaction.Fee = "3"
	from := "ux-owner"
	if len(r.Addresses) > 0 {
		from = r.Addresses[0]
	}
	t.Transaction.In = []api.CreatedTransactionInput{{UxID: "ux1", Address: from}}
	for _, to := range r.To {
		t.Transaction.Out = append(t.Transaction.Out, api.CreatedTransactionOutput{Address: to.Address, Coins: to.Coins, Hours: "5"})
	}
	t.Transaction.Out = append(t.Transaction.Out, api.CreatedTransactionOutput{Address: *r.ChangeAddress, Coins: "0.5", Hours: "4"})
	t.EncodedTransaction = "UNSIGNEDRAW"
	return &t, nil
}

func TestPrepareBrowser(t *testing.T) {
	x := sendNode(t)
	w := browserWallet(x)
	fresh := newAddr()
	r, s := PrepareBrowser(x.n, BrowserSendForm{Wallet: w, Opts: SendForm{To: x.toSpent, Coins: "1.5"}, FreshChange: fresh, FreshReceive: newAddr()})
	require.Empty(t, r.Err)
	require.Equal(t, []string{warnToSpent}, r.Warnings)
	req := x.n.created2[0]
	require.Equal(t, fresh, *req.ChangeAddress, "the fresh change address by default")
	require.Equal(t, w.addrs(), req.Addresses)
	require.Equal(t, api.HoursSelection{Type: "auto", Mode: "share", ShareFactor: "0.5"}, req.HoursSelection)
	require.Equal(t, SignRequest{
		Inputs:  []SignInput{{Hash: "ux1", Address: x.spent}},
		Outputs: []SignOutput{{Address: x.toSpent, Coins: 1_500_000, Hours: 5}, {Address: fresh, Coins: 500_000, Hours: 4}},
		Dests:   []SignDest{{To: x.toSpent, Coins: 1_500_000}},
	}, s)
	require.Empty(t, r.Txid, "an unsigned transaction's id is not the one it will have")

	r, _ = PrepareBrowser(x.n, BrowserSendForm{Wallet: w, Opts: SendForm{To: x.toFresh, Coins: "1", Change: x.change}, FreshChange: fresh})
	require.Equal(t, []string{warnChangeReused}, r.Warnings)
	r, _ = PrepareBrowser(x.n, BrowserSendForm{Wallet: w, Opts: SendForm{To: x.toFresh, Coins: "1"}, FreshChange: "nope"})
	require.Equal(t, errBadChange.Error(), r.Err)

	receive := newAddr()
	r, s = PrepareBrowser(x.n, BrowserSendForm{Wallet: w, Opts: SendForm{From: x.spent}, FreshChange: fresh, FreshReceive: receive})
	require.Empty(t, r.Err)
	req = x.n.created2[len(x.n.created2)-1]
	require.Equal(t, []string{x.spent}, req.Addresses, "a move spends only the exposed address")
	require.Equal(t, receive, req.To[0].Address)
	require.Equal(t, "2.000000", req.To[0].Coins)
	require.Equal(t, receive, s.Dests[0].To)
	r, _ = PrepareBrowser(x.n, BrowserSendForm{Wallet: w, Opts: SendForm{From: x.toFresh}, FreshChange: fresh, FreshReceive: receive})
	require.Equal(t, errBadTo.Error(), r.Err)
}

func TestBrowserSendHandlers(t *testing.T) {
	x := sendNode(t)
	x.n.browser = true
	h := Handler(x.n, Options{})
	w := browserWallet(x)

	body := get(t, h, "/bw/send?id="+w.ID).Body.String()
	require.Regexp(t, regexp.MustCompile(`<input type="password" id="bw-send-password" autocomplete="current-password">`), body, "the password field has no name")
	require.Contains(t, body, `hx-trigger="bw-review"`)

	fresh := newAddr()
	rr := postForm(h, "/bw/send/review", url.Values{
		"wallet": {jsonOf(t, w)}, "to": {x.toFresh}, "coins": {"1"}, "fresh_change": {fresh}, "fresh_receive": {newAddr()},
	})
	body = rr.Body.String()
	m := regexp.MustCompile(`data-bw-sign="([^"]*)"`).FindStringSubmatch(body)
	require.Len(t, m, 2, body)
	var s SignRequest
	require.NoError(t, json.Unmarshal([]byte(html.UnescapeString(m[1])), &s))
	require.Equal(t, []SignDest{{To: x.toFresh, Coins: 1_000_000}}, s.Dests)
	require.Contains(t, body, `id="bw-sign-send"`)

	rr = postForm(h, "/bw/send/inject", url.Values{"id": {w.ID}, "rawtx": {"SIGNED"}})
	require.Contains(t, rr.Body.String(), `href="bw/wallet?id=`+w.ID+`"`)
	require.Equal(t, []string{"SIGNED"}, x.n.injected)

	body = postForm(h, "/bw/wallet/view", url.Values{"wallet": {jsonOf(t, w)}}).Body.String()
	require.Contains(t, body, `href="bw/send?id=`+w.ID+`&amp;from=`+x.spent+`"`, "the move link")
	require.True(t, strings.Contains(body, `href="bw/send?id=`+w.ID+`"`), "the send button")
}

func TestBrowserAdvancedSend(t *testing.T) {
	x := sendNode(t)
	x.n.browser = true
	h := Handler(x.n, Options{})
	w := browserWallet(x)
	fresh := newAddr()
	last := func() api.CreateTransactionRequest { return x.n.created2[len(x.n.created2)-1] }
	base := func() url.Values {
		return url.Values{"wallet": {jsonOf(t, w)}, "fresh_change": {fresh}, "fresh_receive": {newAddr()}}
	}

	v := base()
	v["to"] = []string{x.toFresh, x.toReused}
	v["coins"] = []string{"1", "0.25"}
	v["hours"] = []string{"2", "3"}
	v.Set("hours_mode", "manual")
	v["from_out"] = []string{"uxA"}
	body := postForm(h, "/bw/send/review", v).Body.String()
	req := last()
	require.Equal(t, "manual", req.HoursSelection.Type)
	require.Equal(t, []string{"uxA"}, req.UxOuts)
	require.Empty(t, req.Addresses)
	require.Equal(t, "3", req.To[1].Hours)
	m := regexp.MustCompile(`data-bw-sign="([^"]*)"`).FindStringSubmatch(body)
	require.Len(t, m, 2, body)
	var s SignRequest
	require.NoError(t, json.Unmarshal([]byte(html.UnescapeString(m[1])), &s))
	require.Equal(t, []SignDest{{To: x.toFresh, Coins: 1_000_000}, {To: x.toReused, Coins: 250_000}}, s.Dests)

	v = base()
	v.Set("bulk", x.toFresh+",2\n"+x.toSpent+",1")
	v["from_addr"] = []string{x.spent}
	postForm(h, "/bw/send/review", v)
	require.Len(t, last().To, 2)
	require.Equal(t, []string{x.spent}, last().Addresses)
	v["from_addr"] = []string{x.toFresh}
	require.Contains(t, postForm(h, "/bw/send/review", v).Body.String(), "not in this wallet")

	v = base()
	v.Set("to", x.toFresh)
	v.Set("coins", "1")
	v.Set("unsigned", "1")
	body = postForm(h, "/bw/send/review", v).Body.String()
	require.Contains(t, body, "data-bw-unsigned-tx=")
	require.Contains(t, body, "UNSIGNEDRAW")
	require.NotContains(t, body, `id="bw-sign-send"`)

	body = postForm(h, "/bw/send/sources", url.Values{"wallet": {jsonOf(t, w)}}).Body.String()
	require.Contains(t, body, `data-fill-all="5.000000"`)
	require.Contains(t, body, `name="from_addr" value="`+x.spent+`"`)

	page := get(t, h, "/bw/send?id="+w.ID).Body.String()
	require.Contains(t, page, `hx-post="bw/send/sources"`)
	require.Contains(t, page, `data-bw-unsigned`)
	require.Contains(t, page, `data-fill-link`)
}
