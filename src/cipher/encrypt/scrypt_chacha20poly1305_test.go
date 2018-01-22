package encrypt

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScryptChacha20poly1305Encrypt(t *testing.T) {
	for i := uint(20); i < 21; i++ {
		name := fmt.Sprintf("N=1<<%v r=%v p=%v keyLen=%v", i, 8, 1, 32)
		t.Run(name, func(t *testing.T) {
			crypto := NewScryptChacha20poly1305(1<<i, 8, 1, 32)
			encData, err := crypto.Encrypt([]byte("plaintext"), []byte("password"))
			require.NoError(t, err)

			data, err := base64.StdEncoding.DecodeString(string(encData))
			require.NoError(t, err)
			// Checks the prefix
			ml := binary.LittleEndian.Uint16(data[:metaLengthSize])
			require.True(t, ml <= math.MaxUint16)
			require.True(t, int(metaLengthSize+ml) <= len(data))
			var m meta
			require.NoError(t, json.Unmarshal(data[metaLengthSize:metaLengthSize+ml], &m))
			require.Equal(t, m.N, 1<<i)
			require.Equal(t, m.R, 8)
			require.Equal(t, m.P, 1)
			require.Equal(t, m.KeyLen, 32)
		})
	}
}

func TestScryptChacha20poly1305Decrypt(t *testing.T) {
	tt := []struct {
		name   string
		data   []byte
		encPwd []byte
		decPwd []byte
		err    error
	}{
		{
			"ok",
			[]byte("plaintext"),
			[]byte("pwd"),
			[]byte("pwd"),
			nil,
		},
		{
			"invalid password",
			[]byte("plaintext"),
			[]byte("pwd"),
			[]byte("wrong password"),
			errors.New("chacha20poly1305: message authentication failed"),
		},
		{
			"missing password",
			[]byte("plaintext"),
			[]byte("pwd"),
			nil,
			errors.New("missing password"),
		},
	}

	for _, tc := range tt {
		for i := uint(20); i < 21; i++ {
			name := fmt.Sprintf("N=1<<%v r=8 p=1 keyLen=32 %v", i, tc.name)
			t.Run(name, func(t *testing.T) {
				crypto := NewScryptChacha20poly1305(1<<i, 8, 1, 32)
				encData, err := crypto.Encrypt(tc.data, tc.encPwd)
				require.NoError(t, err)

				data, err := crypto.Decrypt(encData, tc.decPwd)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.Equal(t, tc.data, data)
			})
		}
	}
}
