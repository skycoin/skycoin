package pqsig

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignVerifyBlock(t *testing.T) {
	sk, err := NewSecKey()
	require.NoError(t, err)
	pk, err := sk.PubKey()
	require.NoError(t, err)
	require.Len(t, pk, PubKeySize)

	h := sha256.Sum256([]byte("header"))
	sig, err := SignBlock(sk, h[:])
	require.NoError(t, err)
	require.Len(t, sig, SigSize)
	require.NoError(t, VerifyBlock(pk, h[:], sig))

	again, err := SignBlock(sk, h[:])
	require.NoError(t, err)
	require.Equal(t, sig, again, "signatures are deterministic")

	other := sha256.Sum256([]byte("other header"))
	require.Equal(t, ErrInvalidSig, VerifyBlock(pk, other[:], sig))
	require.Equal(t, ErrMissingSig, VerifyBlock(pk, h[:], nil))

	sk2, err := NewSecKey()
	require.NoError(t, err)
	pk2, err := sk2.PubKey()
	require.NoError(t, err)
	require.Equal(t, ErrInvalidSig, VerifyBlock(pk2, h[:], sig))
}

func TestHexRoundTrip(t *testing.T) {
	sk, err := NewSecKey()
	require.NoError(t, err)
	pk, err := sk.PubKey()
	require.NoError(t, err)

	sk2, err := SecKeyFromHex(sk.Hex())
	require.NoError(t, err)
	require.Equal(t, sk, sk2)
	pk2, err := PubKeyFromHex(pk.Hex())
	require.NoError(t, err)
	require.Equal(t, pk, pk2)

	_, err = PubKeyFromHex("00")
	require.Error(t, err)
	_, err = SecKeyFromHex("00")
	require.Error(t, err)
}
