package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWifiInterfacesList(t *testing.T) {
	ifaces, err := NewWifiInterfaces()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v\n", ifaces)
	assert.Nil(t, nil)
}

func TestWifiInterfaceModeManaged(t *testing.T) {

	// Test DHCP

	// Test Addresses
	/*
		ifaces, err := WifiInterfaces()
		iface := ifaces[0]
		iface.Interface.Addrs()[0]
		if addrs, err := iface.Addrs(); err == nil {
			for _, addr := range addrs {
				//inter.Name, addr
			}
		}
	*/
}

func TestWifiInterfaceModeAdhoc(t *testing.T) {

}
