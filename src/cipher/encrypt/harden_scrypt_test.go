package encrypt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Regression tests for the audit DoS findings on ScryptChacha20poly1305.Decrypt:
// attacker-controlled metadata must not be able to panic or OOM the process.
func TestDecryptEmptyInputDoesNotPanic(t *testing.T) {
	// Previously sliced encData[:2] on an empty buffer and panicked.
	_, err := DefaultScryptChacha20poly1305.Decrypt([]byte{}, []byte("password"))
	require.Error(t, err)
}

func TestValidateScryptParams(t *testing.T) {
	// Default parameters must remain valid (do not break existing wallets).
	require.NoError(t, validateScryptParams(ScryptN, ScryptR, ScryptP, ScryptKeyLen))

	// Non-positive / degenerate parameters are rejected.
	require.Error(t, validateScryptParams(0, 8, 1, 32))
	require.Error(t, validateScryptParams(1<<20, 0, 1, 32))
	require.Error(t, validateScryptParams(1<<20, 8, 0, 32))
	require.Error(t, validateScryptParams(1<<20, 8, 1, 0))

	// Memory-exhausting parameters (e.g. N=1<<24,r=8 ⇒ ~16 GiB) are rejected
	// before scrypt.Key is ever called.
	require.Error(t, validateScryptParams(1<<24, 8, 1, 32))
	require.Error(t, validateScryptParams(1<<26, 8, 1, 32))
	// Oversized keyLen is rejected.
	require.Error(t, validateScryptParams(1<<20, 8, 1, 1<<25))
}
