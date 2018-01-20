package wallet

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecretSetAndGet(t *testing.T) {
	type tst struct {
		i int
		b []byte
		s string
		f float32
	}

	tt := []struct {
		name  string
		key   string
		value interface{}
		err   error
	}{
		{
			"v=string",
			"key",
			"value",
			nil,
		},
		{
			"v=int",
			"key",
			9080,
			nil,
		},
		{
			"v=int8",
			"key",
			int8(10),
			nil,
		},
		{
			"v=int16",
			"key",
			int16(10),
			nil,
		},
		{
			"v=int32",
			"key",
			int32(10),
			nil,
		},
		{
			"v=int64",
			"key",
			int64(10),
			nil,
		},
		{
			"v=uint",
			"key",
			uint(10),
			nil,
		},
		{
			"v=uint8",
			"key",
			uint8(10),
			nil,
		},
		{
			"v=uint16",
			"key",
			uint16(10),
			nil,
		},
		{
			"v=uint32",
			"key",
			uint32(10),
			nil,
		},
		{
			"v=uint64",
			"key",
			uint64(10),
			nil,
		},
		{
			"v=float32",
			"key",
			float32(10.10),
			nil,
		},
		{
			"v=float64",
			"key",
			float64(10.10),
			nil,
		},
		{
			"v=struct",
			"key",
			struct {
				tst
				st tst
			}{
				tst: tst{
					10,
					[]byte{0x02},
					"v1",
					11.11,
				},
				st: tst{
					10,
					[]byte{0x03},
					"v2",
					12.12,
				},
			},
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := make(secrets)
			require.NoError(t, s.set(tc.key, tc.value))
			v := reflect.New(reflect.TypeOf(tc.value))
			require.Equal(t, tc.err, s.get(tc.key, v.Interface()))
		})
	}
}
