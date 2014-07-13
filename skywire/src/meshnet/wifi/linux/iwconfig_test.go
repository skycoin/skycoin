package linux

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIWConfigInfo(t *testing.T) {
	iwconfig := NewIWConfig()

	// IsInstalled
	if !iwconfig.IsInstalled() {
		t.Skip("skipping test, program not installed")
	}

	// InfoList
	infoList, err := iwconfig.InfoList()
	assert.Nil(t, err)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", infoList)

	// Scan
	// ScanRefresh
	// Parse

}
