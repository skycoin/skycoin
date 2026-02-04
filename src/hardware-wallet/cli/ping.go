package cli

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"

	messages "github.com/skycoin/hardware-wallet-protob/go"

	skyWallet "github.com/skycoin/skycoin/src/hardware-wallet/skywallet"
)

func init() {
	pingCmd.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
}

var pingCmd = &cobra.Command{
	Use:   "ping [message]",
	Short: "Send a ping message to the device and display the response.",
	Long: `Send a ping message to the device. The device will echo it back in a Success response.

Special test commands (prefix with TEST:):
  TEST:SHA256  - Test SHA256 hash
  TEST:RIPEMD  - Test RIPEMD160 hash
  TEST:B58     - Test Base58 encoding
  TEST:SQR     - Test field squaring
  TEST:MUL     - Test field multiplication
  TEST:PUBKEY1 - Test pubkey from seckey=1
  TEST:PUBKEY2 - Test pubkey from seckey=2
  TEST:ECMULT  - Test EC multiplication
  TEST:GPOINT  - Get generator point
  TEST:ADDR    - Test address generation`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		device := skyWallet.NewDevice(skyWallet.DeviceTypeFromString(deviceType))
		if device == nil {
			return fmt.Errorf("failed to create device")
		}
		defer device.Close()

		// Check that device is in firmware mode (not bootloader)
		if err := requireFirmwareMode(device); err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("Hint: The device may be in bootloader mode. Flash firmware first, or use --skip to bypass this check")
			return err
		}

		if os.Getenv("AUTO_PRESS_BUTTONS") == "1" && device.Driver.DeviceType() == skyWallet.DeviceTypeEmulator && runtime.GOOS == "linux" {
			err := device.SetAutoPressButton(true, skyWallet.ButtonRight)
			if err != nil {
				return err
			}
		}

		msg, err := device.Ping(args[0])
		if err != nil {
			return err
		}

		switch msg.Kind {
		case uint16(messages.MessageType_MessageType_Success):
			success := &messages.Success{}
			err = proto.Unmarshal(msg.Data, success)
			if err != nil {
				return err
			}
			fmt.Printf("Response: %s\n", success.GetMessage())
		case uint16(messages.MessageType_MessageType_Failure):
			msgData, err := skyWallet.DecodeFailMsg(msg)
			if err != nil {
				return err
			}
			fmt.Println("Failure:", msgData)
		default:
			return fmt.Errorf("received unexpected message type: %s", messages.MessageType(msg.Kind))
		}
		return nil
	},
}
