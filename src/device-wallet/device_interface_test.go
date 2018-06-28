package deviceWallet

import(
	"testing"

	"github.com/skycoin/skycoin/src/util/logging"
	messages "github.com/skycoin/skycoin/protob"
	"github.com/stretchr/testify/require"
)

var logger = logging.MustGetLogger("deviceWallet")

func TestGetAddressUsb(t *testing.T) {
	if (DeviceConnected(DeviceTypeUsb) == false) {
		logger.Fatal("TestGetAddressUsb do not work if Usb device is not connected")
		return
	}

	require.True(t, DeviceConnected(DeviceTypeUsb)) 
	WipeDevice(DeviceTypeUsb)
	// need to connect the usb device
	DeviceSetMnemonic(DeviceTypeUsb, "cloud flower upset remain green metal below cup stem infant art thank")
	kind, address := DeviceAddressGen(DeviceTypeUsb, messages.SkycoinAddressType_AddressTypeSkycoin, 1)
	logger.Info(address)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	
	kind, address = DeviceAddressGen(DeviceTypeUsb, messages.SkycoinAddressType_AddressTypeSkycoin, 2)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address, "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")
}


func TestGetAddressEmulator(t *testing.T) {	
	if (DeviceConnected(DeviceTypeEmulator) == false) {
		logger.Fatal("TestGetAddressEmulator do not work if Emulator device is not running")
		return
	}

	require.True(t, DeviceConnected(DeviceTypeEmulator)) 
	WipeDevice(DeviceTypeEmulator)
	DeviceSetMnemonic(DeviceTypeEmulator, "cloud flower upset remain green metal below cup stem infant art thank")
	kind, address := DeviceAddressGen(DeviceTypeEmulator, messages.SkycoinAddressType_AddressTypeSkycoin, 1)
	logger.Info(address)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address, "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")

	kind, address = DeviceAddressGen(DeviceTypeEmulator, messages.SkycoinAddressType_AddressTypeSkycoin, 2)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address, "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")
}