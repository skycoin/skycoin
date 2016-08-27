package linux

import (
	//"github.com/stretchr/testify/assert"
	"testing"
)

func TestRFKill(t *testing.T) {
	/*
		rfkill := NewRFKill()
		var identifier string
		var err error

		// IsInstalled
		if !rfkill.IsInstalled() {
			t.Skip("skipping test, program not installed")
		}

		// ListAll
		rfks, err := rfkill.ListAll()
		assert.Nil(t, err)
		t.Logf("%#v\n", rfks)

		// SoftBlock
		identifier = "all"
		rfkill.SoftUnblock(identifier)
		assert.Equal(t, rfkill.IsBlocked(identifier), false)
		assert.Equal(t, rfkill.IsSoftBlocked(identifier), false)
		rfkill.SoftBlock(identifier)
		assert.Equal(t, rfkill.IsBlocked(identifier), true)
		assert.Equal(t, rfkill.IsSoftBlocked(identifier), true)

		identifier = "wlan"
		rfkill.SoftUnblock(identifier)
		assert.Equal(t, rfkill.IsBlocked(identifier), false)
		assert.Equal(t, rfkill.IsSoftBlocked(identifier), false)
		rfkill.SoftBlock(identifier)
		assert.Equal(t, rfkill.IsBlocked(identifier), true)
		assert.Equal(t, rfkill.IsSoftBlocked(identifier), true)

		// SoftUnblock
		identifier = "all"
		rfkill.SoftBlock(identifier)
		assert.Equal(t, rfkill.IsBlocked(identifier), true)
		assert.Equal(t, rfkill.IsSoftBlocked(identifier), true)
		rfkill.SoftUnblock(identifier)
		assert.Equal(t, rfkill.IsBlocked(identifier), false)
		assert.Equal(t, rfkill.IsSoftBlocked(identifier), false)

		identifier = "wlan"
		rfkill.SoftBlock(identifier)
		assert.Equal(t, rfkill.IsBlocked(identifier), true)
		assert.Equal(t, rfkill.IsSoftBlocked(identifier), true)
		rfkill.SoftUnblock(identifier)
		assert.Equal(t, rfkill.IsBlocked(identifier), false)
		assert.Equal(t, rfkill.IsSoftBlocked(identifier), false)

		// IsBlockedAfterUnblocking
		identifier = "all"
		rfkill.SoftBlock(identifier)
		assert.Equal(t, rfkill.IsBlockedAfterUnblocking(identifier), false)

		identifier = "wlan"
		rfkill.SoftBlock(identifier)
		assert.Equal(t, rfkill.IsBlockedAfterUnblocking(identifier), false)
	*/
}
