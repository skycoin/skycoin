package params

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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
			divisor, err := DropletPrecisionToDivisor(tc.precision)
			require.NoError(t, err)
			require.Equal(t, tc.divisor, divisor, "%d != %d", tc.divisor, divisor)
		})
	}

	_, err := DropletPrecisionToDivisor(7)
	require.Equal(t, ErrInvalidMaxDropletPrecision, err)
}
