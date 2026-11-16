package liteclient

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet/sweep/plan"
)

type sweepFixture struct {
	addrs []Address
	uxs   []coin.UxOut
	req   SweepRequest
}

// newSweepFixture is a deterministic wallet of four addresses. The first two
// have sent and hold coins, the third holds coins and never sent.
func newSweepFixture(t *testing.T) sweepFixture {
	t.Helper()
	addrs, err := GenerateAddresses(testSeed(), 4)
	if err != nil {
		t.Fatal(err)
	}
	f := sweepFixture{addrs: addrs}
	for i, a := range addrs {
		f.req.Addresses = append(f.req.Addresses, SweepAddress{Address: a.Address})
		if i == 3 {
			continue
		}
		ux := coin.UxOut{
			Head: coin.UxHead{Time: 1000, BkSeq: 1},
			Body: coin.UxBody{SrcTransaction: cipher.SumSHA256([]byte{byte(i)}), Address: cipher.MustDecodeBase58Address(a.Address), Coins: uint64(i+1) * 1e6, Hours: 100},
		}
		f.uxs = append(f.uxs, ux)
		f.req.Outputs = append(f.req.Outputs, SweepOutput{Hash: ux.Hash().Hex(), Address: a.Address, Coins: ux.Body.Coins, Hours: 100})
	}
	f.req.Sent = []string{addrs[0].Address, addrs[1].Address}
	f.req.Seen = []string{addrs[0].Address, addrs[1].Address, addrs[2].Address}
	return f
}

func (f sweepFixture) plan(t *testing.T) SweepResult {
	t.Helper()
	body, err := json.Marshal(f.req)
	if err != nil {
		t.Fatal(err)
	}
	out, err := SweepPlan(string(body))
	if err != nil {
		t.Fatalf("SweepPlan: %v", err)
	}
	var res SweepResult
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatal(err)
	}
	return res
}

func TestSweepPlanSignsWithPrepareTransaction(t *testing.T) {
	f := newSweepFixture(t)
	f.req.Rule.Wallet = "w"
	f.req.Rule.Enabled = true
	f.req.Rule.Destination.Kind = "wallet"

	res := f.plan(t)
	if res.Exposed != 2 || res.Coins != 3e6 || res.Hours != 200 {
		t.Fatalf("exposed %d, coins %d, hours %d", res.Exposed, res.Coins, res.Hours)
	}
	if len(res.Need) != 1 || res.Need[0] != -1 || len(res.Sweeps) != 0 {
		t.Fatalf("want the level 1 destination asked for first, got %+v", res)
	}

	// The fresh address is the one after the stored ones.
	f.req.Destinations = map[string]string{"-1": f.addrs[3].Address}
	res = f.plan(t)
	if len(res.Sweeps) != 1 || len(res.Errors) != 0 {
		t.Fatalf("want one sweep, got %+v", res)
	}
	sw := res.Sweeps[0]
	if sw.Fee != 200 || len(sw.Outputs) != 1 || sw.Outputs[0].Hours != 0 || sw.Outputs[0].Coins != 3e6 || sw.Outputs[0].Address != f.addrs[3].Address {
		t.Fatalf("unexpected sweep %+v", sw)
	}

	// The caller adds the secrets and signs with the existing entry point.
	secrets := map[string]string{}
	for _, a := range f.addrs {
		secrets[a.Address] = a.Secret
	}
	inputs := make([]TransactionInput, len(sw.Inputs))
	byHash := map[string]coin.UxOut{}
	for _, u := range f.uxs {
		byHash[u.Hash().Hex()] = u
	}
	var spent coin.UxArray
	for i, in := range sw.Inputs {
		inputs[i] = TransactionInput{Hash: in.Hash, Address: in.Address, Secret: secrets[in.Address]}
		spent = append(spent, byHash[in.Hash])
	}
	inBody, _ := json.Marshal(inputs)      //nolint:errcheck
	outBody, _ := json.Marshal(sw.Outputs) //nolint:errcheck
	raw, err := PrepareTransaction(string(inBody), string(outBody))
	if err != nil {
		t.Fatalf("PrepareTransaction: %v", err)
	}
	txn, err := coin.DeserializeTransactionHex(raw)
	if err != nil {
		t.Fatal(err)
	}
	if err := txn.VerifyInputSignatures(spent); err != nil {
		t.Fatalf("signatures: %v", err)
	}
	if txn.Out[0].Hours != 0 || len(txn.In) != 2 {
		t.Fatalf("unexpected transaction %+v", txn)
	}
}

func TestSweepPlanMappingsAndFallback(t *testing.T) {
	f := newSweepFixture(t)
	cold := "2Hx321N7dmkWuWiTN7pVzCCARtGsKAdK6Pi"
	f.req.Rule = ruleJSON(t, `{"wallet":"w","destination":{"kind":"wallet"},"mappings":[{"from":{"address":"`+f.addrs[1].Address+`"},"to":{"kind":"address","address":"`+cold+`"}}]}`)
	f.req.Destinations = map[string]string{"0": cold, "-1": f.addrs[3].Address}

	res := f.plan(t)
	if res.Exposed != 2 || res.Uncovered != 1 || len(res.Sweeps) != 1 || res.Sweeps[0].Destination != cold || res.Sweeps[0].Outputs[0].Coins != 2e6 {
		t.Fatalf("level 1 off: want only the mapped address swept, got %+v", res)
	}

	f.req.Rule.Enabled = true
	res = f.plan(t)
	if res.Uncovered != 0 || len(res.Sweeps) != 2 {
		t.Fatalf("level 1 on: want the rest swept too, got %+v", res)
	}

	off := false
	f.req.Rule.IncludeSpent = &off
	res = f.plan(t)
	if res.Exposed != 1 || len(res.Sweeps) != 1 || res.Sweeps[0].Destination != cold {
		t.Fatalf("include spent off: want only the listed address, got %+v", res)
	}

	f.req.Covered = []string{f.req.Outputs[1].Hash}
	res = f.plan(t)
	if len(res.Sweeps) != 0 || res.Coins != 0 {
		t.Fatalf("an output a pending sweep spends is not swept again, got %+v", res)
	}
}

func TestSweepPlanRefusesExposedDestination(t *testing.T) {
	f := newSweepFixture(t)
	f.req.Rule.Wallet = "w"
	f.req.Rule.Enabled = true
	f.req.Rule.Destination.Kind = "wallet"
	f.req.Destinations = map[string]string{"-1": f.addrs[1].Address}
	res := f.plan(t)
	if len(res.Sweeps) != 0 || len(res.Errors) != 1 || !strings.Contains(res.Errors[0], "the destination is exposed") {
		t.Fatalf("want the exposed destination refused, got %+v", res)
	}

	// A bip44 destination in a shared account is listed with its place and refused.
	f = newSweepFixture(t)
	for i := range f.req.Addresses {
		f.req.Addresses[i].Bip44, f.req.Addresses[i].Child = true, uint32(i) //nolint:gosec // test index
	}
	f.req.Rule = ruleJSON(t, `{"wallet":"w","enabled":true,"destination":{"kind":"wallet"},"shared_accounts":[0]}`)
	f.req.Destinations = map[string]string{"-1": f.addrs[3].Address}
	res = f.plan(t)
	if len(res.Errors) != 1 || res.Gap["0"] == 0 {
		t.Fatalf("want a refusal and a gap to fill, got %+v", res)
	}
}

func TestSweepPlanNoHours(t *testing.T) {
	f := newSweepFixture(t)
	for i := range f.req.Outputs {
		f.req.Outputs[i].Hours = 0
	}
	f.req.Rule = ruleJSON(t, `{"wallet":"w","enabled":true,"destination":{"kind":"wallet"}}`)
	f.req.Destinations = map[string]string{"-1": f.addrs[3].Address}
	res := f.plan(t)
	if !res.NoHours || len(res.Sweeps) != 0 {
		t.Fatalf("want a wait for coin hours, got %+v", res)
	}
}

func TestSweepCheck(t *testing.T) {
	f := newSweepFixture(t)
	check := func(rule string) error {
		f.req.Rule = ruleJSON(t, rule)
		body, err := json.Marshal(f.req)
		if err != nil {
			t.Fatal(err)
		}
		return SweepCheck(string(body))
	}
	if err := check(`{"wallet":"w","enabled":true,"destination":{"kind":"wallet"}}`); err != nil {
		t.Fatal(err)
	}
	cases := map[string]string{
		`{"wallet":"w","destination":{"kind":"wallet"},"mappings":[{"from":{"address":"` + f.addrs[2].Address + `"},"to":{"kind":"wallet"}}]}`: "never sent",
		`{"wallet":"w","destination":{"kind":"address","address":"` + f.addrs[0].Address + `"}}`:                                               "the destination is exposed",
		`{"wallet":"w","destination":{"kind":"wallet"},"mappings":[{"from":{"account":0},"to":{"kind":"wallet"}}]}`:                            "has no account 0",
		`{"wallet":"w","destination":{"kind":"nope"}}`:                                                                                         "unknown destination kind",
	}
	for rule, want := range cases {
		if err := check(rule); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: got %v, want %q", rule, err, want)
		}
	}
}

func ruleJSON(t *testing.T, s string) plan.Rule {
	t.Helper()
	var r plan.Rule
	if err := json.Unmarshal([]byte(s), &r); err != nil {
		t.Fatal(err)
	}
	return r
}
