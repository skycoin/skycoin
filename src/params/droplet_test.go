package params

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	_require "github.com/skycoin/skycoin/src/testutil/require"
)

func TestDropletPrecisionToDivisor(t *testing.T) {
	cases := []struct {
		precision uint8
		divisor   uint64
	}{
		{0, 1e6},
		{1, 1e5},
		{2, 1e4},
		{3, 1e3},
		{4, 1e2},
		{5, 1e1},
		{6, 1},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("DropletPrecisionToDivisor(%d)=%d", tc.precision, tc.divisor)
		t.Run(name, func(t *testing.T) {
			divisor := DropletPrecisionToDivisor(tc.precision)
			require.Equal(t, tc.divisor, divisor, "%d != %d", tc.divisor, divisor)
		})
	}

	_require.PanicsWithLogMessage(t, "precision must be <= droplet.Exponent", func() {
		DropletPrecisionToDivisor(7)
	})
}
