package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/readable"
)

type fakeUxOuts map[string][]readable.SpentOutput

func (f fakeUxOuts) AddressUxOuts(addr string) ([]readable.SpentOutput, error) {
	if addr == "down" {
		return nil, errors.New("connection refused")
	}
	return f[addr], nil
}

func TestWarnSpentAddresses(t *testing.T) {
	f := fakeUxOuts{
		"spent":    {{SpentBlockSeq: 0}, {SpentBlockSeq: 42}},
		"received": {{SpentBlockSeq: 0}},
	}
	var buf bytes.Buffer
	warnSpentAddresses(f, &buf, []string{"spent", "received", "new", "spent", "down"})

	out := buf.String()
	require.Equal(t, 1, strings.Count(out, "spent has spent before"))
	require.NotContains(t, out, "received has")
	require.NotContains(t, out, "new has")
	require.Contains(t, out, "could not check whether down has spent before")
}
