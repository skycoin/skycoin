// +build !386,!amd64p32,!arm,!armbe,!mips,!mipsle,!mips64p32,!mips64p32le,!ppc,!s390,!sparc

package mathutil

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test64BitIntToUint32(t *testing.T) {
	// Remaining tests defined in TestIntToUint32
	cases := []struct {
		a   int
		b   uint32
		err error
	}{
		{
			a:   math.MaxUint32 + 1,
			err: ErrIntOverflowsUint32,
		},
		{
			a: math.MaxUint32,
			b: math.MaxUint32,
		},
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
