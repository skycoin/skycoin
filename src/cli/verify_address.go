package cli

import (
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
)

func verifyAddressCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Verify a skycoin address",
		Use:                   "verifyAddress [skycoin address]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE: func(_ *cobra.Command, args []string) error {
			_, err := cipher.DecodeBase58Address(args[0])
			return err
		},
	}
}
