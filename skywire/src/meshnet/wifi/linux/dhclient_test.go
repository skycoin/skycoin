package linux

import (
	//"github.com/stretchr/testify/assert"
	"testing"
)

func TestDHClient(t *testing.T) {
	dhclient := NewDHClient()

	// IsInstalled
	if !dhclient.IsInstalled() {
		t.Skip("skipping test, program not installed")
	}
}
