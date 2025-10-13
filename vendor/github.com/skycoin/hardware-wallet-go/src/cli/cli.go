package cli

import (
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/spf13/cobra"
)

const (
	Version = "1.7.0"
)

var log = logging.MustGetLogger("skycoin-hw-cli")

//RootCmd is the root command
var RootCmd = &cobra.Command{
		Use:     "skycoin-hw-cli",
		Short:   "the skycoin hardware wallet command line interface",
		Version: Version,
	}

func init() {
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
	)
}
