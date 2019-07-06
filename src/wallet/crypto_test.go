package wallet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecrets(t *testing.T) {
	s := make(Secrets)
	s.set("k1", "v1")

	v, ok := s.get("k1")
	require.True(t, ok)
	require.Equal(t, "v1", v)

	_, ok = s.get("k2")
	require.False(t, ok)

	s.set("k2", "v2")

	b, err := s.serialize()
	require.NoError(t, err)

	s1 := make(Secrets)
	err = s1.deserialize(b)
	require.NoError(t, err)
	require.Equal(t, s, s1)
}
