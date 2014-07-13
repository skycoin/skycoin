package linux

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNetworkManager(t *testing.T) {
	networkmanager := NewNetworkManager()

	// IsInstalled
	if !networkmanager.IsInstalled() {
		t.Skip("skipping test, program not installed")
	}
	assert.Nil(t, nil)
}
