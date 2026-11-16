package ui2

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/wallet/sweep"
)

func (f *fakeNode) sweepResponse(id string) *api.SweepResponse {
	resp := &api.SweepResponse{}
	resp.Status.Wallet = id
	resp.Status.WaitingCoins = "0.000000"
	if r, ok := f.sweeps[id]; ok {
		resp.Rule = &r
		resp.Status.Enabled = r.Enabled
	}
	if ttl, ok := f.unlocked[id]; ok {
		resp.Status.Unlocked = true
		if ttl > 0 {
			t := time.Date(2026, 11, 15, 12, 0, 0, 0, time.UTC)
			resp.Status.UnlockExpires = &t
		}
	}
	if !resp.Status.Enabled {
		resp.Status.Paused = "disabled"
	}
	if resp.Status.Enabled && !resp.Status.Unlocked {
		resp.Status.Paused = "paused: wallet locked, 7.000000 coins waiting on exposed addresses"
		resp.Status.WaitingCoins = "7.000000"
	}
	return resp
}

func (f *fakeNode) WalletSweep(id string) (*api.SweepResponse, error) {
	if _, err := f.Wallet(id); err != nil {
		return nil, err
	}
	return f.sweepResponse(id), nil
}

func (f *fakeNode) SetWalletSweep(r sweep.Rule) (*api.SweepResponse, error) {
	if f.sweepErr != nil {
		err := f.sweepErr
		f.sweepErr = nil
		return nil, err
	}
	if f.sweeps == nil {
		f.sweeps = map[string]sweep.Rule{}
	}
	f.sweeps[r.Wallet] = r
	return f.sweepResponse(r.Wallet), nil
}

func (f *fakeNode) WalletSweepUnlock(req api.SweepUnlockRequest) (*api.SweepResponse, error) {
	if req.Password != "pw" {
		return nil, api.NewClientError("401 Unauthorized", 401, "invalid password")
	}
	if f.unlocked == nil {
		f.unlocked = map[string]uint64{}
	}
	f.unlocked[req.ID] = req.TTLSeconds
	return f.sweepResponse(req.ID), nil
}

func (f *fakeNode) WalletSweepLock(id string) (*api.SweepResponse, error) {
	delete(f.unlocked, id)
	return f.sweepResponse(id), nil
}

func sweepNodeWithAccounts() *fakeNode {
	n := testNode()
	n.wallets[0].Accounts = []readable.WalletAccount{{Index: 0}, {Index: 1, Name: "safe"}}
	n.wallets = append(n.wallets,
		api.WalletResponse{Meta: readable.WalletMeta{Filename: "b.wlt", Label: "Cold", Type: "deterministic"}},
		api.WalletResponse{Meta: readable.WalletMeta{Filename: "t.wlt", Type: "bip44", Temp: true}},
		api.WalletResponse{Meta: readable.WalletMeta{Filename: "x.wlt", Type: "xpub"}},
	)
	for _, id := range []string{"b.wlt", "t.wlt", "x.wlt"} {
		n.balances[id] = &api.BalanceResponse{}
	}
	return n
}

func TestSweepSection(t *testing.T) {
	n := sweepNodeWithAccounts()
	h := Handler(n, Options{})

	body := get(t, h, "/wallet?id=a.wlt").Body.String()
	require.Contains(t, body, `id="auto-sweep"`)
	require.Contains(t, body, `hx-get="wallet/sweep?id=a.wlt"`)
	require.Contains(t, body, `hx-trigger="toggle once"`, "the panel loads when opened")
	require.Contains(t, body, "Auto-move coins off exposed addresses")

	for _, id := range []string{"t.wlt", "x.wlt"} {
		require.NotContains(t, get(t, h, "/wallet?id="+id).Body.String(), `id="auto-sweep"`, id)
	}
}

func TestSweepPanel(t *testing.T) {
	n := sweepNodeWithAccounts()
	h := Handler(n, Options{})

	body := get(t, h, "/wallet/sweep?id=a.wlt").Body.String()
	require.Contains(t, body, `name="enabled"`)
	require.NotContains(t, body, `name="enabled" value="1" checked`, "off by default")
	require.Contains(t, body, "burns ALL the coin hours")
	require.Contains(t, body, `name="burn_ok"`)
	require.Contains(t, body, `name="shared" value="1"`)
	require.Contains(t, body, "Account 1 · safe")
	require.Contains(t, body, `name="to" value="next" checked`, "a bip44 wallet sweeps to the next account by default")
	require.Contains(t, body, `<option value="b.wlt"`, "another wallet is offered")
	require.NotContains(t, body, `<option value="t.wlt"`, "a temporary wallet is not")
	require.Contains(t, body, `hx-post="wallet/sweep/unlock"`, "an encrypted wallet can be unlocked")
	require.Contains(t, body, `<option value="0">Until the node restarts</option>`)
	require.Contains(t, body, `hx-get="wallet/sweep/status?id=a.wlt&amp;enc=1" hx-trigger="every 10s"`)
	require.Contains(t, body, "data-poll-open")
	require.Contains(t, body, `badge off`)
	require.NotContains(t, body, ">disabled<", "the reason a rule that is off is not running is not shown")

	plain := get(t, h, "/wallet/sweep?id=b.wlt").Body.String()
	require.NotContains(t, plain, `name="shared"`)
	require.Contains(t, plain, `name="to" value="self" checked`)
	require.NotContains(t, plain, "wallet/sweep/unlock", "an unencrypted wallet needs no unlock")

	require.Contains(t, get(t, h, "/wallet/sweep?id=missing.wlt").Body.String(), "wallet not found")
}

func TestSweepRule(t *testing.T) {
	n := sweepNodeWithAccounts()
	h := Handler(n, Options{})
	post := func(v url.Values) string {
		v.Set("id", "a.wlt")
		return postForm(h, "/wallet/sweep/rule", v).Body.String()
	}

	body := post(url.Values{"enabled": {"1"}, "shared": {"0"}, "to": {"next"}})
	require.Contains(t, body, "Confirm that every move burns")
	require.Empty(t, n.sweeps, "nothing is set before the cost is confirmed")
	require.Contains(t, body, `name="shared" value="0" checked`, "the form keeps what was posted")
	require.Contains(t, body, `name="burn_ok"`, "and still asks for the confirmation")

	body = post(url.Values{"enabled": {"1"}, "burn_ok": {"1"}, "shared": {"0"}, "to": {"account"}, "to_account": {"1"}, "addresses": {" a1 \n a2 "}})
	require.Contains(t, body, "Saved.")
	r := n.sweeps["a.wlt"]
	require.True(t, r.Enabled)
	require.Equal(t, []uint32{0}, r.SharedAccounts)
	require.Equal(t, sweep.DestAccount, r.Destination.Kind)
	require.Equal(t, uint32(1), *r.Destination.Account)
	require.Equal(t, []string{"a1", "a2"}, r.Addresses)
	require.Contains(t, body, "badge warn-text", "a locked encrypted wallet is paused")
	require.Contains(t, body, "paused: wallet locked, 7.000000 coins waiting")
	require.NotContains(t, body, `name="burn_ok"`, "an enabled rule does not ask again")

	r.Gap = 50
	n.sweeps["a.wlt"] = r
	body = post(url.Values{"enabled": {"1"}, "to": {"wallet"}, "to_wallet": {"b.wlt"}})
	require.Contains(t, body, "Saved.")
	require.Equal(t, sweep.Destination{Kind: sweep.DestWallet, Wallet: "b.wlt"}, n.sweeps["a.wlt"].Destination)
	require.Equal(t, uint32(50), n.sweeps["a.wlt"].Gap, "fields the form does not show are kept")
	require.Empty(t, n.sweeps["a.wlt"].SharedAccounts)

	n.sweepErr = api.NewClientError("400 Bad Request", 400, "destination account 0 is shared")
	body = post(url.Values{"enabled": {"1"}, "shared": {"0"}, "to": {"account"}, "to_account": {"0"}})
	require.Contains(t, body, "destination account 0 is shared", "the node's refusal is shown")
	require.Equal(t, 1, strings.Count(body, "destination account 0 is shared"), "once")
	require.Contains(t, body, `name="to" value="account" checked`)

	require.Contains(t, post(url.Values{"to": {"address"}, "to_address": {" "}}), "Choose a destination")
	require.Contains(t, post(url.Values{"to": {"bogus"}}), "Choose a destination")
	require.Contains(t, post(url.Values{"to": {"account"}, "to_account": {"x"}}), "Choose a valid account")

	body = post(url.Values{"to": {"address"}, "to_address": {"2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv"}})
	require.Contains(t, body, "Saved.")
	require.False(t, n.sweeps["a.wlt"].Enabled, "unticking turns it off")
	require.Contains(t, body, "badge off")
}

func TestSweepUnlockLock(t *testing.T) {
	n := sweepNodeWithAccounts()
	n.sweeps = map[string]sweep.Rule{"a.wlt": {Wallet: "a.wlt", Enabled: true, Destination: sweep.Destination{Kind: sweep.DestAccount}}}
	h := Handler(n, Options{})
	do := func(op string, v url.Values) string {
		v.Set("id", "a.wlt")
		return postForm(h, "/wallet/sweep/"+op, v).Body.String()
	}

	require.Contains(t, do("unlock", url.Values{}), "Enter the password.")
	require.Contains(t, do("unlock", url.Values{"password": {"bad"}, "ttl": {"900"}}), "invalid password")

	body := do("unlock", url.Values{"password": {"pw"}, "ttl": {"3600"}})
	require.Contains(t, body, "The node can now sign moves")
	require.Equal(t, uint64(3600), n.unlocked["a.wlt"])
	require.Contains(t, body, "Until 2026-11-15 12:00:00 UTC")
	require.Contains(t, body, "badge on")
	require.Contains(t, body, `hx-post="wallet/sweep/lock"`)
	require.NotContains(t, body, `type="password"`, "no password field while unlocked")

	body = do("lock", url.Values{})
	require.Contains(t, body, "The node forgot the password.")
	require.NotContains(t, n.unlocked, "a.wlt")

	do("unlock", url.Values{"password": {"pw"}, "ttl": {"0"}})
	require.Contains(t, get(t, h, "/wallet/sweep/status?id=a.wlt&enc=1").Body.String(), "Until the node restarts")
}

func TestSweepStatus(t *testing.T) {
	n := sweepNodeWithAccounts()
	h := Handler(n, Options{})
	last := time.Date(2026, 11, 14, 9, 30, 0, 0, time.UTC)
	n.sweeps = map[string]sweep.Rule{"b.wlt": {Wallet: "b.wlt", Enabled: true}}
	n.unlocked = map[string]uint64{"b.wlt": 0}
	sn := &statusNode{fakeNode: n, status: func(s *api.SweepStatus) {
		s.Exposed, s.Contested, s.Confirmed = 21, 1, 4
		s.Destinations = []string{"destAddr"}
		s.LastSweep = &last
		s.LastTxids = []string{"txLast"}
		s.Pending = []string{"txPending"}
		s.Displaced = []string{"txDisplaced"}
		s.LastError = "inject failed"
	}}
	body := get(t, Handler(sn, Options{}), "/wallet/sweep/status?id=b.wlt").Body.String()
	for _, s := range []string{"badge on", "21", "destAddr", "2026-11-14 09:30:00 UTC", "txLast", "txPending", "txDisplaced", "Contested outputs", "inject failed", `data-copy="txLast"`} {
		require.Contains(t, body, s)
	}
	require.True(t, strings.HasPrefix(strings.TrimSpace(body), `<div class="sweep-status"`), "the poll swaps only the status")

	require.Contains(t, get(t, h, "/wallet/sweep/status?id=missing.wlt").Body.String(), "wallet not found")
	require.Equal(t, 405, get(t, h, "/wallet/sweep/rule?id=a.wlt").Code)
}

// statusNode lets a test fill in the status the node reports.
type statusNode struct {
	*fakeNode
	status func(*api.SweepStatus)
}

func (s *statusNode) WalletSweep(id string) (*api.SweepResponse, error) {
	resp, err := s.fakeNode.WalletSweep(id)
	if err == nil {
		s.status(&resp.Status)
	}
	return resp, err
}
