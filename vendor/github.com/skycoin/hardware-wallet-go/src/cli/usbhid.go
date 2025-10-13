package cli

import (

	"github.com/spf13/cobra"
	skyWallet "github.com/skycoin/hardware-wallet-go/src/skywallet"
)

func init() {
	getUsbDetails.Flags().StringVar(&deviceType, "deviceType", "USB", "Device type to send instructions to, hardware wallet (USB) or emulator.")
}

var getUsbDetails = &cobra.Command{
		Use:   "getUsbDetails",
		Short: "Ask host usb about details for the hardware wallet",
		RunE: func(_ *cobra.Command, _ []string) error {
			device := skyWallet.NewDevice(skyWallet.DeviceTypeFromString(deviceType))
			if device == nil {
				return nil
			}
			defer device.Close()

			infos, err := device.GetUsbInfo()
			if err != nil {
				log.Errorln(err)
			}
			for infoIdx := range infos {
				log.Infoln("-----------------------------------------")
				if infos[infoIdx].VendorID == skyWallet.SkycoinVendorID {
					log.Printf("%-13d%-5s%s", infos[infoIdx].VendorID, "==>", "Skycoin Foundation")
				}
				if infos[infoIdx].ProductID == skyWallet.SkycoinHwProductID {
					log.Printf("%-13d%-5s%s", infos[infoIdx].ProductID, "==>", "Hardware Wallet")
				}
				log.Printf("%-13s%-5s%s", "Device path", "==>", infos[infoIdx].Path)
			}
			return nil
		},
	}
