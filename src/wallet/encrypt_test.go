package wallet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptAndDecrypt(t *testing.T) {
	tt := []struct {
		name string
		data string
		key  string
		err  error
	}{
		{
			"one byte with empty key",
			"1",
			"",
			nil,
		},
		{
			"one byte, with key",
			"1",
			"key",
			nil,
		},
		{
			"encrypt data size < 32, with key",
			"hello,world",
			"pwd",
			nil,
		},
		{
			"encrypt data size = 32, with key",
			"22Be3vk9nSbXHYJ8JMkCPphe73NQvGhm",
			"pwd",
			nil,
		},
		{
			"encrypt data size > 32, with key",
			"22Be3vk9nSbXHYJ8JMkCPphe73NQvGhmabc",
			"pwd13",
			nil,
		},
		{
			"encrypt data size > 32, with key",
			"22Be3vk9nSbXHYJ8JMkCPphe73NQvGhmabc",
			"pwd13",
			nil,
		},
		{
			"encrypt data size > 2*32, with key",
			"9d4a8e0b9fe800e40790cab00710ef2a49dbdb6be2b0e0950c739fe0799d0",
			"8JMkCPphe73NQvGhmab",
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			enData, err := encrypt([]byte(tc.data), []byte(tc.key))
			require.NoError(t, err)
			str := string(enData)

			d, err := decrypt(str, []byte(tc.key))
			require.NoError(t, err)
			require.Equal(t, tc.data, string(d))
		})
	}

}
