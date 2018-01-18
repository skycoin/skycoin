package wallet

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/chacha20poly1305"
)

func TestScryptChacha20poly1305Encrypt(t *testing.T) {
	tt := []struct {
		name     string
		opts     Options
		addrsNum uint64
		pwd      []byte
		err      error
	}{
		{
			"ok address=0",
			Options{
				Seed: "seed",
			},
			0,
			[]byte("pwd"),
			nil,
		},
		{
			"ok address=1",
			Options{
				Seed: "seed",
			},
			1,
			[]byte("pwd"),
			nil,
		},
		{
			"ok address=2",
			Options{
				Seed: "seed",
			},
			2,
			[]byte("pwd"),
			nil,
		},
		{
			"missing password",
			Options{
				Seed: "seed",
			},
			0,
			nil,
			ErrMissingPassword,
		},
		{
			"wallet already encrypted",
			Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			0,
			[]byte("pwd"),
			ErrWalletEncrypted,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			crypto := newScryptChacha20poly1305Crypto(scryptN, scryptR, scryptP, scryptKeyLen, scryptSaltLen)
			w, err := NewWallet("t.wlt", tc.opts)
			require.NoError(t, err)
			_, err = w.GenerateAddresses(tc.addrsNum)
			require.NoError(t, err)

			err = crypto.Encrypt(w, tc.pwd)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			require.NoError(t, w.validate())

			// Checks if wallet is encrypted
			require.True(t, w.IsEncrypted())
			require.Equal(t, CryptoTypeScryptChacha20poly1305, w.cryptoType())

			// Checks authenticated fields
			auth, err := w.authenticated()
			require.NoError(t, err)
			require.NotNil(t, auth)
			require.Equal(t, scryptN, auth.N)
			require.Equal(t, scryptR, auth.R)
			require.Equal(t, scryptP, auth.P)
			require.Equal(t, scryptKeyLen, auth.KeyLen)
			require.Len(t, auth.Nonce, chacha20poly1305.NonceSize)
			require.Len(t, auth.Salt, scryptSaltLen)

			// Checks if the seeds already encrypted, must not be the same as original
			require.NotEqual(t, tc.opts.Seed, w.seed())

			// Checks if the entries are encrypted
			// generates unencrypted entries from options.Seed
			lsed, keys := cipher.GenerateDeterministicKeyPairsSeed([]byte(tc.opts.Seed), int(tc.addrsNum))
			require.NotEqual(t, hex.EncodeToString(lsed), w.lastSeed())

			for i, k := range keys {
				addr := cipher.AddressFromSecKey(k)
				require.Equal(t, addr, w.Entries[i].Address)
				require.Equal(t, cipher.SecKey{}, w.Entries[i].Secret)
			}
		})
	}
}

func TestScryptChacha20poly1305Decrypt(t *testing.T) {
	tt := []struct {
		name     string
		opts     Options
		addrsNum uint64
		pwd      []byte
		err      error
	}{
		{
			"ok address=1",
			Options{
				Seed:       "seed",
				Encrypt:    true,
				Password:   []byte("pwd"),
				CryptoType: CryptoTypeScryptChacha20poly1305,
			},
			1,
			[]byte("pwd"),
			nil,
		},
		{
			"ok address=2",
			Options{
				Seed:       "seed",
				Encrypt:    true,
				Password:   []byte("pwd"),
				CryptoType: CryptoTypeScryptChacha20poly1305,
			},
			2,
			[]byte("pwd"),
			nil,
		},
		{
			"ok address=3",
			Options{
				Seed:       "seed",
				Encrypt:    true,
				Password:   []byte("pwd"),
				CryptoType: CryptoTypeScryptChacha20poly1305,
			},
			2,
			[]byte("pwd"),
			nil,
		},
		{
			"wallet not encrypted",
			Options{
				Seed: "seed",
			},
			0,
			[]byte("pwd"),
			ErrWalletNotEncrypted,
		},
		{
			"authentication failed",
			Options{
				Seed:       "seed",
				Encrypt:    true,
				Password:   []byte("pwd"),
				CryptoType: CryptoTypeScryptChacha20poly1305,
			},
			2,
			[]byte("wrong password"),
			ErrAuthenticationFailed{errors.New("chacha20poly1305: message authentication failed")},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			crypto := newScryptChacha20poly1305Crypto(scryptN, scryptR, scryptP, scryptKeyLen, scryptSaltLen)
			// Creates unencrypted wallet first
			encrypt := tc.opts.Encrypt
			pwd := tc.opts.Password

			tc.opts.Encrypt = false
			tc.opts.Password = nil

			w, err := NewWallet("t.wlt", tc.opts)
			require.NoError(t, err)

			// Generates addresses
			_, err = w.GenerateAddresses(tc.addrsNum)
			require.NoError(t, err)

			// Encrypts wallet if the original options.Encrypt is true
			if encrypt {
				require.NoError(t, crypto.Encrypt(w, pwd))
			}

			dw, err := crypto.Decrypt(w, tc.pwd)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			require.False(t, dw.IsEncrypted())

			lsd, keys := cipher.GenerateDeterministicKeyPairsSeed([]byte(tc.opts.Seed), int(tc.addrsNum))

			// Checks if the seeds are matched
			require.Equal(t, tc.opts.Seed, dw.seed())
			require.Equal(t, hex.EncodeToString(lsd), dw.lastSeed())

			for i, k := range keys {
				addr := cipher.AddressFromSecKey(k)
				require.Equal(t, addr, dw.Entries[i].Address)
				require.Equal(t, k, dw.Entries[i].Secret)
			}

			// Checks the original wallet
			require.True(t, w.IsEncrypted())
			require.NotEqual(t, w.seed(), dw.seed())
			require.NotEqual(t, w.lastSeed(), dw.lastSeed())
			for i := range w.Entries {
				require.NotEqual(t, w.Entries[i].Secret, dw.Entries[i].Secret)
			}
		})
	}
}
