package skydexui

import (
	"math"
	"strconv"
	"strings"
	"time"
)

// shortKey keeps head and tail characters of a long key, as the React client
// shows keys and ids. The full key goes in the title.
func shortKey(k string, head, tail int) string {
	if len(k) <= head+tail+2 {
		return k
	}
	return k[:head] + "…" + k[len(k)-tail:]
}

// timeAgo follows ConnectGate.jsx: ts is in milliseconds.
func timeAgo(ts int64, now time.Time) string {
	if ts <= 0 {
		return ""
	}
	s := max(int64(0), (now.UnixMilli()-ts)/1000)
	switch {
	case s < 60:
		return "just now"
	case s < 3600:
		return strconv.FormatInt(s/60, 10) + "m ago"
	case s < 86400:
		return strconv.FormatInt(s/3600, 10) + "h ago"
	}
	return strconv.FormatInt(s/86400, 10) + "d ago"
}

// num prints an amount in plain decimal notation.
func num(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

// uriSchemes are the BIP 21 style schemes, from QrButton.jsx, whose amount
// parameter is a decimal coin value. Any other coin gets an address only QR.
var uriSchemes = map[string]string{
	"BTC":  "bitcoin",
	"LTC":  "litecoin",
	"DOGE": "dogecoin",
	"DASH": "dash",
	"BCH":  "bitcoincash",
	"SKY":  "skycoin",
}

// buildURI is the payment URI for coin, empty when there is no scheme for it
// or no amount to put in it.
func buildURI(address string, amount float64, coin string) string {
	scheme := uriSchemes[strings.ToUpper(coin)]
	if scheme == "" || address == "" || !(amount > 0) || math.IsInf(amount, 0) {
		return ""
	}
	a := strings.TrimRight(strconv.FormatFloat(amount, 'f', 8, 64), "0")
	return scheme + ":" + address + "?amount=" + strings.TrimSuffix(a, ".")
}

// amountText is "<amount> <coin>", empty without either.
func amountText(amount float64, coin string) string {
	if amount == 0 || coin == "" {
		return ""
	}
	return num(amount) + " " + coin
}

// shortID keeps the head of a long id, as the React tables do.
func shortID(id string) string {
	if len(id) <= 12 {
		return id
	}
	return id[:8] + "…"
}
