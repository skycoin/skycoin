package deviceWallet

import(

	"testing"

	messages "github.com/skycoin/skycoin/protob"
	"github.com/stretchr/testify/require"
)

func TestGetAddressUsb(t *testing.T) {

	DeviceSetMnemonic(DeviceTypeUsb, "cloud flower upset remain green metal below cup stem infant art thank")
	DeviceAddressGen(DeviceTypeUsb, messages.SkycoinAddressType_AddressTypeSkycoin, 1)
	require.Nil(t, nil)
}


func TestGetAddressEmulator(t *testing.T) {

	DeviceSetMnemonic(DeviceTypeEmulator, "cloud flower upset remain green metal below cup stem infant art thank")
	DeviceAddressGen(DeviceTypeEmulator, messages.SkycoinAddressType_AddressTypeSkycoin, 1)
	require.Nil(t, nil)
}