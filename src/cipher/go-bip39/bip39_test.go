package bip39

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsMnemonicValid(t *testing.T) {
	m, err := NewDefaultMnemonic()
	require.NoError(t, err)
	require.True(t, IsMnemonicValid(m))

	// Truncated
	m = m[:len(m)-15]
	require.False(t, IsMnemonicValid(m))

	// Trailing whitespace
	m, err = NewDefaultMnemonic()
	require.NoError(t, err)
	m += " "
	require.False(t, IsMnemonicValid(m))

	m, err = NewDefaultMnemonic()
	require.NoError(t, err)
	m += "\n"
	require.False(t, IsMnemonicValid(m))

	// Preceding whitespace
	m, err = NewDefaultMnemonic()
	require.NoError(t, err)
	m = " " + m
	require.False(t, IsMnemonicValid(m))

	m, err = NewDefaultMnemonic()
	require.NoError(t, err)
	m = "\n" + m
	require.False(t, IsMnemonicValid(m))

	// Extra whitespace between words
	m, err = NewDefaultMnemonic()
	require.NoError(t, err)
	ms := strings.Split(m, " ")
	m = strings.Join(ms, "  ")
	require.False(t, IsMnemonicValid(m))

	// Contains invalid word
	m, err = NewDefaultMnemonic()
	require.NoError(t, err)
	ms = strings.Split(m, " ")
	ms[2] = "foo"
	m = strings.Join(ms, " ")
	require.False(t, IsMnemonicValid(m))

	// Invalid number of words
	m, err = NewDefaultMnemonic()
	require.NoError(t, err)
	ms = strings.Split(m, " ")
	m = strings.Join(ms[:len(ms)-1], " ")
	require.False(t, IsMnemonicValid(m))
}
