package deviceWallet

import(

	"testing"

	messages "github.com/skycoin/skycoin/protob"
	"github.com/stretchr/testify/require"
)

func TestGetAddressUsb(t *testing.T) {

	// need to connect the usb device
	DeviceSetMnemonic(DeviceTypeUsb, "cloud flower upset remain green metal below cup stem infant art thank")
	kind, data := DeviceAddressGen(DeviceTypeUsb, messages.SkycoinAddressType_AddressTypeSkycoin, 1)
	require.Equal(t, kind, uint16(2)) //Success message
	require.Equal(t, data[2:], []byte("2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw"))
	
	kind, data = DeviceAddressGen(DeviceTypeUsb, messages.SkycoinAddressType_AddressTypeSkycoin, 2)
	require.Equal(t, kind, uint16(2)) //Success message
	require.Equal(t, data[2:], []byte("zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs"))
}


func TestGetAddressEmulator(t *testing.T) {

	DeviceSetMnemonic(DeviceTypeEmulator, "cloud flower upset remain green metal below cup stem infant art thank")
	kind, data := DeviceAddressGen(DeviceTypeUsb, messages.SkycoinAddressType_AddressTypeSkycoin, 1)
	require.Equal(t, kind, uint16(2)) //Success message
	require.Equal(t, data[2:], []byte("2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw"))

	kind, data = DeviceAddressGen(DeviceTypeUsb, messages.SkycoinAddressType_AddressTypeSkycoin, 2)
	require.Equal(t, kind, uint16(2)) //Success message
	require.Equal(t, data[2:], []byte("zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs"))
}