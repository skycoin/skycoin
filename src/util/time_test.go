package util

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNow(t *testing.T) {
    now := Now()
    assert.False(t, now.IsZero())
    now2 := Now()
    assert.True(t, now2.After(now))
}

func TestUnixNow(t *testing.T) {
    now := Now()
    unow := UnixNow()
    assert.True(t, now.Unix() == unow || now.Unix() == unow-1)
}

func TestZeroTime(t *testing.T) {
    z := ZeroTime()
    assert.True(t, z.IsZero())
}
