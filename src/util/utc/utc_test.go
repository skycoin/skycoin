package utc

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
