package wallet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecrets(t *testing.T) {
	s := make(Secrets)
	s.Set("k1", "v1")

	v, ok := s.Get("k1")
	require.True(t, ok)
	require.Equal(t, "v1", v)

	_, ok = s.Get("k2")
	require.False(t, ok)

	s.Set("k2", "v2")

	b, err := s.Serialize()
	require.NoError(t, err)

	s1 := make(Secrets)
	err = s1.Deserialize(b)
	require.NoError(t, err)
	require.Equal(t, s, s1)
}
