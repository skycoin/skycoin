package cli

import (
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/util/logging"
)

const (
	// Version is the CLI version string.
	Version = "1.7.0"
)

var log = logging.MustGetLogger("skycoin-hw-cli")

// RootCmd is the root command
var RootCmd = &cobra.Command{
	Use:     "skycoin-hw-cli",
	Short:   "the skycoin hardware wallet command line interface",
	Version: Version,
}

func init() {
	// Add global --skip flag to bypass device mode checks
	RootCmd.PersistentFlags().BoolVar(&skipModeCheck, "skip", false, "Skip device mode verification (firmware/bootloader check)")

	RootCmd.AddCommand(
		applySettingsCmd,
		setMnemonicCmd,
		featuresCmd,
		generateMnemonicCmd,
		addressGenCmd,
		firmwareUpdate,
		signMessageCmd,
		checkMessageSignatureCmd,
		setPinCode,
		removePinCode,
		wipeCmd,
		backupCmd,
		recoveryCmd,
		cancelCmd,
		transactionSignCmd,
		getRawEntropyCmd,
		getMixedEntropyCmd,
		getUsbDetails,
		pingCmd,
	)
}
