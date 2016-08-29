package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUDevAdm(t *testing.T) {
	udevadm := NewUDevAdm()

	// IsInstalled
	if !udevadm.IsInstalled() {
		t.Skip("skipping test, program not installed")
	}

	// Run
	out, err := udevadm.Run("wlan0")
	assert.Nil(t, err)

	t.Logf("%#v\n", out)
}
