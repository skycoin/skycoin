package deviceWallet

import(

	"testing"

	messages "github.com/skycoin/skycoin/protob"
	"github.com/stretchr/testify/require"
)

func TestGetAddressUsb(t *testing.T) {

	// need to connect the usb device
	DeviceSetMnemonic(DeviceTypeUsb, "cloud flower upset remain green metal below cup stem infant art thank")
	kind, address := DeviceAddressGen(DeviceTypeUsb, messages.SkycoinAddressType_AddressTypeSkycoin, 1)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	
	kind, address = DeviceAddressGen(DeviceTypeUsb, messages.SkycoinAddressType_AddressTypeSkycoin, 2)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address, "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")
}


func TestGetAddressEmulator(t *testing.T) {

	DeviceSetMnemonic(DeviceTypeEmulator, "cloud flower upset remain green metal below cup stem infant art thank")
	kind, address := DeviceAddressGen(DeviceTypeEmulator, messages.SkycoinAddressType_AddressTypeSkycoin, 1)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")

	kind, address = DeviceAddressGen(DeviceTypeEmulator, messages.SkycoinAddressType_AddressTypeSkycoin, 2)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address, "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")
}