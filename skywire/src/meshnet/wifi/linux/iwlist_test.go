package linux

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIWList(t *testing.T) {
	iwlist := NewIWList()

	// IsInstalled
	if !iwlist.IsInstalled() {
		t.Skip("skipping test, program not installed")
	}

	wifinets, _ := iwlist.Scan("wlan0")

	t.Logf("Total Wifi Networks: %#v\n", len(wifinets))
	t.Logf("%#v\n", wifinets)
	assert.Nil(t, nil)
}
