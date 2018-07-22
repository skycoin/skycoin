package deviceWallet

import (
	"testing"

	"github.com/stretchr/testify/require"

	messages "github.com/skycoin/skycoin/protob"
)

func TestGetAddressUsb(t *testing.T) {
	if DeviceConnected(DeviceTypeUsb) == false {
		t.Skip("TestGetAddressUsb do not work if Usb device is not connected")
		return
	}

	require.True(t, DeviceConnected(DeviceTypeUsb))
	WipeDevice(DeviceTypeUsb)
	// need to connect the usb device
	DeviceSetMnemonic(DeviceTypeUsb, "cloud flower upset remain green metal below cup stem infant art thank")
	kind, address := DeviceAddressGen(DeviceTypeUsb, 2, 0)
	logger.Info(address)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address[0], "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.Equal(t, address[1], "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")
}

func TestGetAddressEmulator(t *testing.T) {
	if DeviceConnected(DeviceTypeEmulator) == false {
		t.Skip("TestGetAddressEmulator do not work if Emulator device is not running")
		return
	}

	require.True(t, DeviceConnected(DeviceTypeEmulator))
	WipeDevice(DeviceTypeEmulator)
	DeviceSetMnemonic(DeviceTypeEmulator, "cloud flower upset remain green metal below cup stem infant art thank")
	kind, address := DeviceAddressGen(DeviceTypeEmulator, 2, 0)
	logger.Info(address)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address[0], "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.Equal(t, address[1], "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")
}
