package crypto

import (
	"testing"

	"github.com/SkycoinProject/skycoin/src/wallet/secrets"
	"github.com/stretchr/testify/require"
)

// TODO: avoid the dependency of secrets package
func TestSecrets(t *testing.T) {
	s := make(secrets.Secrets)
	s.Set("k1", "v1")

	v, ok := s.Get("k1")
	require.True(t, ok)
	require.Equal(t, "v1", v)

	_, ok = s.Get("k2")
	require.False(t, ok)

	s.Set("k2", "v2")

	b, err := s.Serialize()
	require.NoError(t, err)

	s1 := make(secrets.Secrets)
	err = s1.Deserialize(b)
	require.NoError(t, err)
	require.Equal(t, s, s1)
}
