package skydexui

import (
	"strconv"
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
