package ui2

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnits(t *testing.T) {
	for _, c := range []struct {
		v      string
		places int
		want   string
	}{
		{"100000", 8, "0.001"},
		{"150000000", 8, "1.5"},
		{"0", 8, "0"},
		{"1000000000000000000", 18, "1"},
		{"5", 18, "0.000000000000000005"},
		{"-2500000", 6, "-2.5"},
	} {
		v, _ := new(big.Int).SetString(c.v, 10)
		require.Equal(t, c.want, formatUnits(v, c.places), c.v)
	}
	var u struct{ A, B Units }
	require.NoError(t, json.Unmarshal([]byte(`{"A":123,"B":"456"}`), &u))
	require.Equal(t, "123", u.A.String())
	require.Equal(t, "456", u.B.String())

	v, err := parseUnits("0.0005", 8)
	require.NoError(t, err)
	require.Equal(t, "50000", v.String())
	v, err = parseUnits("2", 18)
	require.NoError(t, err)
	require.Equal(t, "2000000000000000000", v.String())
	for _, bad := range []string{"", "0", "0.000000001", "-1", "1e3", "x"} {
		_, err := parseUnits(bad, 8)
		require.Error(t, err, bad)
	}
}
