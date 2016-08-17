package linux

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSysfs(t *testing.T) {
	sysfs := NewSysfs()

	outA := sysfs.Run("wlan0")
	t.Logf("%#v\n", outA)

	outB := sysfs.Run("lo")
	assert.Equal(t, outB.WirelessDirectoryExists, false)
	if outB.WirelessDirectoryExists {
		t.Fatal(errors.New("wifi: lo is not wireless interface"))
	}

	assert.Nil(t, nil)
}
