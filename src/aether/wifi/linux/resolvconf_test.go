package linux

import (
	//"github.com/stretchr/testify/assert"
	//"io/ioutil"
	//"net"
	//"strings"
	"testing"
)

func TestResolvConf(t *testing.T) {
	/*
		resolvconf := NewResolvConf()

		// IsInstalled
		if !resolvconf.IsInstalled() {
			t.Skip("skipping test, program not installed")
		}

		// Set with a List
		nameserversList := []net.IP{
			net.ParseIP("1.1.1.1"),
			net.ParseIP("2.2.2.2"),
		}
		resolvconf.Set("wlan0", "darknet", nameserversList)

		// Set Verified
		out, err := ioutil.ReadFile("/run/resolvconf/interface/wlan0.darknet")
		assert.Nil(t, err)
		outs := string(out)
		assert.True(t, strings.Contains(outs, "1.1.1.1"))
		assert.True(t, strings.Contains(outs, "2.2.2.2"))

		resolvconf.Update()

		// Delete
		// File should exist
		errd := resolvconf.Delete("wlan0", "darknet")
		assert.Nil(t, errd)

		// File should be removed
	*/
}
