package utc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNow(t *testing.T) {
	now := Now()
	require.False(t, now.IsZero())
	time.Sleep(time.Millisecond * 10)
	now2 := Now()
	require.True(t, now2.After(now))
}

func TestUnixNow(t *testing.T) {
	now := Now()
	unow := UnixNow()
	require.True(t, now.Unix() == unow || now.Unix() == unow-1)
}
