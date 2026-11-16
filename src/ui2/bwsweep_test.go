package ui2

import (
	"net/url"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBrowserSweepSection(t *testing.T) {
	x := sendNode(t)
	x.n.browser = true
	h := Handler(x.n, Options{})

	w := browserWallet(x)
	body := postForm(h, "/bw/wallet/view", url.Values{"wallet": {jsonOf(t, w)}}).Body.String()
	require.Contains(t, body, `data-bw-sweep`)
	require.Contains(t, body, "only while it is open and the wallet is unlocked here", "the page says when it sweeps")
	require.Contains(t, body, `data-bw-sweep-unlock`, "an encrypted wallet unlocks in the page")
	require.Contains(t, body, "Until I leave this page")
	require.Contains(t, body, `data-f="shared"`, "a bip44 wallet can mark its xPub shared")
	require.Contains(t, body, "burns ALL the coin hours")
	section := regexp.MustCompile(`(?s)<section class="card" id="auto-sweep".*?</section>`).FindString(body)
	require.NotEmpty(t, section)
	require.NotRegexp(t, regexp.MustCompile(`<(input|textarea|select)[^>]*\sname=`), section, "nothing in the panel can be posted")

	w.Encrypted, w.Type = false, "deterministic"
	body = postForm(h, "/bw/wallet/view", url.Values{"wallet": {jsonOf(t, w)}}).Body.String()
	require.NotContains(t, body, `data-bw-sweep-unlock`, "a wallet without a password needs no unlock")
	require.NotContains(t, body, `data-f="shared"`)

	w.Temporary = true
	require.NotContains(t, postForm(h, "/bw/wallet/view", url.Values{"wallet": {jsonOf(t, w)}}).Body.String(), `data-bw-sweep`)

	w.Temporary, w.Type = false, "skywallet"
	body = postForm(h, "/bw/wallet/view", url.Values{"wallet": {jsonOf(t, w)}}).Body.String()
	require.NotContains(t, body, `data-bw-sweep`)
	require.Contains(t, body, "the device must confirm every signature")
}

func TestBrowserSweepInject(t *testing.T) {
	x := sendNode(t)
	x.n.browser = true
	h := Handler(x.n, Options{})
	body := postForm(h, "/bw/sweep/inject", url.Values{"rawtx": {"RAWTX7"}}).Body.String()
	require.Contains(t, body, `data-txid="tx7"`)
	require.Equal(t, []string{"RAWTX7"}, x.n.injected)
	x.n.refuse = "RAWTX8"
	require.Contains(t, postForm(h, "/bw/sweep/inject", url.Values{"rawtx": {"RAWTX8"}}).Body.String(), `data-error="400 Bad Request - refused"`)
}
