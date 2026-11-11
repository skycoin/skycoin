package explorerui

import (
	"strconv"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/util/droplet"
)

// Coin is how the node names its coin and hours, from /api/v1/health.
type Coin struct {
	FullName      string
	Ticker        string
	Hours         string
	HoursSingular string
	MaxDecimals   int
	QRPrefix      string
}

// defaultCoin follows CoinIdentifiers in explorer/src/app/app.config.ts.
var defaultCoin = Coin{
	FullName:      "Skycoin",
	Ticker:        "SKY",
	Hours:         "Coin Hours",
	HoursSingular: "Coin Hour",
	MaxDecimals:   6,
	QRPrefix:      "skycoin:",
}

// namedAddresses follows explorer/src/app/app.config.ts.
var namedAddresses = map[string]string{
	"2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8": "Binance hot wallet",
}

// AddrName is the name given to an address in parentheses, or empty.
func AddrName(addr string) string {
	if n, ok := namedAddresses[addr]; ok {
		return " (" + n + ")"
	}
	return ""
}

// group puts a comma between each three digits, as the Angular number pipe
// does in its default en-US locale.
func group(digits string) string {
	if len(digits) <= 3 {
		return digits
	}
	var b strings.Builder
	lead := len(digits) % 3
	if lead > 0 {
		b.WriteString(digits[:lead])
	}
	for i := lead; i < len(digits); i += 3 {
		if b.Len() > 0 {
			b.WriteByte(',')
		}
		b.WriteString(digits[i : i+3])
	}
	return b.String()
}

// Int is n with its thousands grouped.
func Int(n uint64) string { return group(strconv.FormatUint(n, 10)) }

// IntS is a decimal string of a whole number with its thousands grouped, the
// string itself when it is not one.
func IntS(s string) string {
	if _, err := strconv.ParseUint(s, 10, 64); err != nil {
		return s
	}
	return group(s)
}

// Coins is an amount of droplets with at most dec decimals, rounded half up
// as the Angular number pipe does. fixed keeps the trailing zeros.
func Coins(droplets int64, dec int, fixed bool) string {
	dec = min(max(dec, 0), droplet.Exponent)
	neg := droplets < 0
	u := uint64(droplets) //nolint:gosec // the sign is kept apart
	if neg {
		u = uint64(-droplets) //nolint:gosec
	}
	unit := uint64(1)
	for range droplet.Exponent - dec {
		unit *= 10
	}
	u = (u + unit/2) / unit
	scale := uint64(1)
	for range dec {
		scale *= 10
	}
	s := group(strconv.FormatUint(u/scale, 10))
	if dec > 0 {
		frac := strconv.FormatUint(u%scale, 10)
		frac = strings.Repeat("0", dec-len(frac)) + frac
		if !fixed {
			frac = strings.TrimRight(frac, "0")
		}
		if frac != "" {
			s += "." + frac
		}
	}
	if neg && s != "0" {
		s = "-" + s
	}
	return s
}

// Droplets reads a coin amount such as "12.5" as droplets, 0 when it is not
// one.
func Droplets(s string) int64 {
	neg := strings.HasPrefix(s, "-")
	d, err := droplet.FromString(strings.TrimPrefix(s, "-"))
	if err != nil {
		return 0
	}
	if neg {
		return -int64(d) //nolint:gosec
	}
	return int64(d) //nolint:gosec
}

// When is a moment shown as UTC on the server and in the reader's own zone
// once explorer.js runs.
type When struct {
	t time.Time
}

func unixWhen(sec uint64) When { return When{time.Unix(int64(sec), 0).UTC()} } //nolint:gosec

// IsZero reports a moment that was never set.
func (w When) IsZero() bool { return w.t.Unix() <= 0 }

// ISO is the moment for a <time datetime>.
func (w When) ISO() string { return w.t.Format(time.RFC3339) }

// Date is yyyy-MM-dd.
func (w When) Date() string { return w.t.Format("2006-01-02") }

// Clock is HH:mm:ss.
func (w When) Clock() string { return w.t.Format("15:04:05") }

// Year is yyyy.
func (w When) Year() string { return w.t.Format("2006") }

// Day is MMM dd in capitals.
func (w When) Day() string { return strings.ToUpper(w.t.Format("Jan 02")) }

// Minute is HH:mm.
func (w When) Minute() string { return w.t.Format("15:04") }

func itoa(n int) string { return strconv.Itoa(n) }

func utoa(n uint64) string { return strconv.FormatUint(n, 10) }
