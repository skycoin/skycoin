package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkManager(t *testing.T) {
	networkmanager := NewNetworkManager()

	// IsInstalled
	if !networkmanager.IsInstalled() {
		t.Skip("skipping test, program not installed")
	}
	assert.Nil(t, nil)
}
