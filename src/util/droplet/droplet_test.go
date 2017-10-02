package droplet

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		s string
		n uint64
		e error
	}{
		{
			s: "0",
			n: 0,
		},
		{
			s: "0.",
			n: 0,
		},
		{
			s: "0.0",
			n: 0,
		},
		{
			s: "0.000000",
			n: 0,
		},
		{
			s: "0.0000000",
			n: 0,
		},
		{
			s: "0.0000001",
			e: ErrTooManyDecimals,
		},
		{
			s: "0.000001",
			n: 1,
		},
		{
			s: "0.0000010",
			n: 1,
		},
		{
			s: "1",
			n: 1e6,
		},
		{
			s: "1.000001",
			n: 1e6 + 1,
		},
		{
			s: "-1",
			e: ErrNegativeValue,
		},
		{
			s: "10000",
			n: 1e4 * 1e6,
		},
		{
			s: "123456789.123456",
			n: 123456789123456,
		},
		{
			s: "123.000456",
			n: 123000456,
		},
		{
			s: "100SKY",
			e: errors.New("can't convert 100SKY to decimal"),
		},
		{
			s: "",
			e: errors.New("can't convert  to decimal"),
		},
		{
			s: "999999999999999999999999999999999999999999",
			e: ErrTooLarge,
		},
		{
			s: "9223372036854.775807",
			n: 9223372036854775807,
		},
		{
			s: "-9223372036854.775807",
			e: ErrNegativeValue,
		},
		{
			s: "9223372036854775808",
			e: ErrTooLarge,
		},
		{
			s: "9223372036854775807.000001",
			e: ErrTooLarge,
		},
		{
			s: "9223372036854775807",
			e: ErrTooLarge,
		},
		{
			s: "9223372036854775806.000001",
			e: ErrTooLarge,
		},
	}

	for _, tcc := range cases {
		tc := tcc
		t.Run(tc.s, func(t *testing.T) {
			t.Parallel()

			n, err := FromString(tc.s)

			if tc.e == nil {
				require.NoError(t, err)
				require.Equal(t, tc.n, n, "result: %d", n)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.e, err)
				require.Equal(t, uint64(0), n, "result: %d", n)
			}
		})
	}
}

func TestToString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		n uint64
		s string
		e error
	}{
		{
			n: 0,
			s: "0.000000",
		},
		{
			n: 1,
			s: "0.000001",
		},
		{
			n: 1e6,
			s: "1.000000",
		},
		{
			n: 100100,
			s: "0.100100",
		},
		{
			n: 1001000,
			s: "1.001000",
		},
		{
			n: 999,
			s: "0.000999",
		},
		{
			n: 999000000,
			s: "999.000000",
		},
		{
			n: 123000456,
			s: "123.000456",
		},
		{
			n: 9223372036854775808,
			e: ErrTooLarge,
		},
	}

	for _, tcc := range cases {
		tc := tcc
		t.Run(tc.s, func(t *testing.T) {
			t.Parallel()

			s, err := ToString(tc.n)

			if tc.e == nil {
				require.NoError(t, err)
				require.Equal(t, tc.s, s)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.e, err)
				require.Equal(t, "", s)
			}
		})
	}
}
