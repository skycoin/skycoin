package encrypt

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScryptChacha20poly1305Encrypt(t *testing.T) {
	for i := uint(1); i < 20; i++ {
		name := fmt.Sprintf("N=1<<%v r=%v p=%v keyLen=%v", i, 8, 1, 32)
		t.Run(name, func(t *testing.T) {
			crypto := ScryptChacha20poly1305{N: 1 << i, R: 8, P: 1, KeyLen: 32}
			encData, err := crypto.Encrypt([]byte("plaintext"), []byte("password"))
			require.NoError(t, err)

			data, err := base64.StdEncoding.DecodeString(string(encData))
			require.NoError(t, err)
			// Checks the prefix
			ml := binary.LittleEndian.Uint16(data[:scryptChacha20MetaLengthSize])
			require.True(t, int(scryptChacha20MetaLengthSize+ml) <= len(data))
			var m meta
			require.NoError(t, json.Unmarshal(data[scryptChacha20MetaLengthSize:scryptChacha20MetaLengthSize+ml], &m))
			require.Equal(t, m.N, 1<<i)
			require.Equal(t, m.R, 8)
			require.Equal(t, m.P, 1)
			require.Equal(t, m.KeyLen, 32)
		})
	}
}

func TestScryptChacha20poly1305Decrypt(t *testing.T) {
	tt := []struct {
		name    string
		data    []byte
		encData []byte
		encPwd  []byte
		decPwd  []byte
		err     error
	}{
		{
			name:    "ok",
			data:    []byte("plaintext"),
			encData: []byte("dQB7Im4iOjUyNDI4OCwiciI6OCwicCI6MSwia2V5TGVuIjozMiwic2FsdCI6ImpiejUrSFNjTFFLWkI5T0tYblNNRmt2WDBPY3JxVGZ0ZFpDNm9KUFpaeHc9Iiwibm9uY2UiOiJLTlhOQmRQa1ZUWHZYNHdoIn3PQFmOot0ETxTuv//skTG7Q57UVamGCgG5"),
			encPwd:  []byte("pwd"),
			decPwd:  []byte("pwd"),
			err:     nil,
		},
		{
			name:    "invalid password",
			data:    []byte("plaintext"),
			encData: []byte("dQB7Im4iOjUyNDI4OCwiciI6OCwicCI6MSwia2V5TGVuIjozMiwic2FsdCI6ImpiejUrSFNjTFFLWkI5T0tYblNNRmt2WDBPY3JxVGZ0ZFpDNm9KUFpaeHc9Iiwibm9uY2UiOiJLTlhOQmRQa1ZUWHZYNHdoIn3PQFmOot0ETxTuv//skTG7Q57UVamGCgG5"),
			encPwd:  []byte("pwd"),
			decPwd:  []byte("wrong password"),
			err:     errors.New("chacha20poly1305: message authentication failed"),
		},
		{
			name:    "missing password",
			data:    []byte("plaintext"),
			encData: []byte("dQB7Im4iOjUyNDI4OCwiciI6OCwicCI6MSwia2V5TGVuIjozMiwic2FsdCI6ImpiejUrSFNjTFFLWkI5T0tYblNNRmt2WDBPY3JxVGZ0ZFpDNm9KUFpaeHc9Iiwibm9uY2UiOiJLTlhOQmRQa1ZUWHZYNHdoIn3PQFmOot0ETxTuv//skTG7Q57UVamGCgG5"),
			encPwd:  []byte("pwd"),
			decPwd:  nil,
			err:     errors.New("missing password"),
		},
	}

	for _, tc := range tt {
		name := fmt.Sprintf("N=1<<19 r=8 p=1 keyLen=32 %v", tc.name)
		t.Run(name, func(t *testing.T) {
			crypto := ScryptChacha20poly1305{}
			data, err := crypto.Decrypt(tc.encData, tc.decPwd)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.data, data)
		})
	}
}
