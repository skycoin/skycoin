// Package ui2 src/ui2/hw_test.go c3-ui-wallet
package ui2

import (
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// skywalletWallet is the public part bw.js posts for a Skywallet's wallet.
func skywalletWallet(x sendFixture) BrowserWallet {
	return BrowserWallet{
		ID: x.spent, Label: "S", Type: "skywallet",
		Addresses: []BrowserAddress{{Address: x.spent, Index: 0}, {Address: x.unused, Index: 1}},
	}
}

func TestSkywalletPage(t *testing.T) {
	x := sendNode(t)
	h := Handler(x.n, Options{})
	body := get(t, h, "/hw").Body.String()
	for _, want := range []string{`data-hw-page`, `data-bw-page="hw"`, `hx-post="hw/list"`, `id="hw-dialog"`, `data-hw-op="generate-12"`, `data-hw-op="wipe"`, `data-msg-err-wrong-pin=`, Asset("hw.js")} {
		require.Contains(t, body, want)
	}
	// The PIN pad is the device's matrix layout, 7 8 9 on top.
	require.Regexp(t, regexp.MustCompile(`data-hw-key="7"[^<]*>[^<]*</button>\s*<button[^>]*data-hw-key="8"`), body)

	list := postForm(h, "/hw/list", url.Values{"wallets": {jsonOf(t, []BrowserWallet{skywalletWallet(x)})}}).Body.String()
	require.Contains(t, list, "skywallet")
	require.Contains(t, list, bwWalletURL(x.spent))
	require.Contains(t, postForm(h, "/hw/list", url.Values{"wallets": {"[]"}}).Body.String(), "No Skywallet wallet")
}

func TestSkywalletWalletView(t *testing.T) {
	x := sendNode(t)
	h := Handler(x.n, Options{})
	body := postForm(h, "/bw/wallet/view", url.Values{"wallet": {jsonOf(t, skywalletWallet(x))}}).Body.String()
	start := strings.Index(body, `id="wallet-options"`)
	require.Positive(t, start)
	opts := body[start : start+strings.Index(body[start:], "</section>")]
	for _, op := range []string{"rename", "hw-address", "hw-confirm", "delete"} {
		require.Contains(t, opts, `data-bw-op="`+op+`"`)
	}
	for _, op := range []string{`data-bw-op="seed"`, `data-bw-op="password"`, `data-bw-op="reset"`} {
		require.NotContains(t, opts, op, "a Skywallet's seed stays on the device")
	}
	require.NotRegexp(t, regexp.MustCompile(`\sname=`), opts, "no option field can be posted")
	require.Contains(t, opts, `<option value="1">`+x.unused+`</option>`)
	require.Contains(t, body, `id="hw-dialog"`)
	require.Contains(t, get(t, h, "/bw/send?id="+x.spent).Body.String(), `id="hw-dialog"`)
}
