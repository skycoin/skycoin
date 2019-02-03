package mathutil

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddUint64(t *testing.T) {
	n, err := AddUint64(10, 11)
	require.NoError(t, err)
	require.Equal(t, uint64(21), n)

	_, err = AddUint64(math.MaxUint64, 1)
	require.Error(t, err)
}

func TestAddUint32(t *testing.T) {
	n, err := AddUint32(10, 11)
	require.NoError(t, err)
	require.Equal(t, uint32(21), n)

	_, err = AddUint32(math.MaxUint32, 1)
	require.Error(t, err)
}

func TestMultUint64(t *testing.T) {
	n, err := MultUint64(10, 11)
	require.NoError(t, err)
	require.Equal(t, uint64(110), n)

	_, err = MultUint64(math.MaxUint64/2, 3)
	require.Error(t, err)
}

func TestUint64ToInt64(t *testing.T) {
	cases := []struct {
		a   uint64
		b   int64
		err error
	}{
		{
			a: 0,
			b: 0,
		},
		{
			a: 1,
			b: 1,
		},
		{
			a: math.MaxInt64,
			b: math.MaxInt64,
		},
		{
			a:   math.MaxUint64,
			err: ErrUint64OverflowsInt64,
		},
		{
			a:   math.MaxInt64 + 1,
			err: ErrUint64OverflowsInt64,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprint(tc.a), func(t *testing.T) {
			x, err := Uint64ToInt64(tc.a)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.Equal(t, tc.b, x)
			}
		})
	}
}

func TestInt64ToUint64(t *testing.T) {
	cases := []struct {
		a   int64
		b   uint64
		err error
	}{
		{
			a: 0,
			b: 0,
		},
		{
			a: 1,
			b: 1,
		},
		{
			a: math.MaxInt64,
			b: math.MaxInt64,
		},
		{
			a:   -math.MaxInt64,
			err: ErrInt64UnderflowsUint64,
		},
		{
			a:   -1,
			err: ErrInt64UnderflowsUint64,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprint(tc.a), func(t *testing.T) {
			x, err := Int64ToUint64(tc.a)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.Equal(t, tc.b, x)
			}
		})
	}
}

func TestIntToUint32(t *testing.T) {
	cases := []struct {
		a   int
		b   uint32
		err error
	}{
		{
			a: 0,
			b: 0,
		},
		{
			a:   -1,
			err: ErrIntUnderflowsUint32,
		},
		{
			a: math.MaxInt32,
			b: math.MaxInt32,
		},
		{
			a: 999,
			b: 999,
		},
		// 64bit test defined in Test64BitIntToUint32
	}

	for _, tc := range cases {
		t.Run(fmt.Sprint(tc.a), func(t *testing.T) {
			x, err := IntToUint32(tc.a)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.Equal(t, tc.b, x)
			}
		})
	}
}
