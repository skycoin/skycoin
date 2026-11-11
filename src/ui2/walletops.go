package ui2

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"github.com/skycoin/skycoin/src/wallet"
)

// Notices a wallet page shows after a change, as catalog keys. The page
// shows only these, so a link cannot put other text on it.
var notices = map[string]string{
	"renamed":   "wallet.done-renamed",
	"added":     "wallet.done-added",
	"scanned":   "wallet.done-scanned",
	"encrypted": "wallet.done-encrypted",
	"decrypted": "wallet.done-decrypted",
}

// Notice is the message a wallet page opens with, empty for none.
type Notice struct {
	Key string
	N   string
}

// ParseNotice reads the notice a wallet page was sent to with.
func ParseNotice(get func(string) string) Notice {
	k, ok := notices[get("done")]
	if !ok {
		return Notice{}
	}
	n := ""
	if v, err := strconv.Atoi(get("n")); err == nil && v >= 0 {
		n = strconv.Itoa(v)
	}
	return Notice{Key: k, N: n}
}

func walletDoneURL(id, done string, n int) string {
	u := walletURL(id) + "&done=" + url.QueryEscape(done)
	if n >= 0 {
		u += "&n=" + strconv.Itoa(n)
	}
	return u
}

// Errors the wallet forms report, as catalog keys.
var (
	errNoName      = errors.New("wallet.err-name")
	errNum         = errors.New("wallet.err-num")
	errNoPassword  = errors.New("wallet.err-password")
	errPasswordTwo = errors.New("create.err-password-confirm")
)

// MaxNewAddresses bounds one request for new addresses or a scan, as the
// Angular wallet's number of addresses dialog does.
const MaxNewAddresses = 100

func parseNum(s string, def int) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return def, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > MaxNewAddresses {
		return 0, errNum
	}
	return n, nil
}

func addressCount(n Node, id string) int {
	w, err := n.Wallet(id)
	if err != nil {
		return 0
	}
	return len(w.Entries)
}

// walletOp is a form that changes wallet id and then reloads its page with
// a notice. It returns the notice and its count, -1 for none.
type walletOp func(n Node, id string, post func(string) string) (string, int, error)

var walletOps = map[string]walletOp{
	"rename": func(n Node, id string, post func(string) string) (string, int, error) {
		label := strings.TrimSpace(post("label"))
		if label == "" {
			return "", 0, errNoName
		}
		return "renamed", -1, n.UpdateWallet(id, label)
	},
	"address": func(n Node, id string, post func(string) string) (string, int, error) {
		num, err := parseNum(post("num"), 1)
		if err != nil {
			return "", 0, err
		}
		w, err := n.Wallet(id)
		if err != nil {
			return "", 0, err
		}
		opts := []wallet.Option{wallet.OptionGenerateN(uint64(num))} //nolint:gosec
		if w.Meta.Type == "bip44" {
			opts = append(opts, wallet.OptionExternal())
		}
		addrs, err := n.NewWalletAddress(id, post("password"), opts...)
		return "added", len(addrs), err
	},
	"scan": func(n Node, id string, post func(string) string) (string, int, error) {
		num, err := parseNum(post("num"), 20)
		if err != nil {
			return "", 0, err
		}
		before := addressCount(n, id)
		if _, err := n.ScanWalletAddresses(id, num, post("password")); err != nil {
			return "", 0, err
		}
		return "scanned", max(addressCount(n, id)-before, 0), nil
	},
	"encrypt": func(n Node, id string, post func(string) string) (string, int, error) {
		pw := post("password")
		switch {
		case pw == "":
			return "", 0, errNoPassword
		case pw != post("password2"):
			return "", 0, errPasswordTwo
		}
		_, err := n.EncryptWallet(id, pw)
		return "encrypted", -1, err
	},
	"decrypt": func(n Node, id string, post func(string) string) (string, int, error) {
		if post("password") == "" {
			return "", 0, errNoPassword
		}
		_, err := n.DecryptWallet(id, post("password"))
		return "decrypted", -1, err
	},
}

// walletOpHandler runs the op the path names and sends the page back to
// the wallet, or renders the error in the form.
func walletOpHandler(n Node) func(w http.ResponseWriter, r *http.Request) templ.Component {
	return func(w http.ResponseWriter, r *http.Request) templ.Component {
		op, ok := walletOps[r.PathValue("op")]
		if !ok {
			http.NotFound(w, r)
			return nil
		}
		id := r.PostFormValue("id")
		done, num, err := op(n, id, r.PostFormValue)
		if err != nil {
			return formError(err)
		}
		w.Header().Set("HX-Redirect", walletDoneURL(id, done, num))
		return nil
	}
}
