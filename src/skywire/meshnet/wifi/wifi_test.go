package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWifiInterfacesList(t *testing.T) {
	ifaces, err := WifiInterfaces()
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
