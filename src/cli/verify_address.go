package cli

import (
	gcli "github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
)

func verifyAddressCmd() *gcli.Command {
	return &gcli.Command{
		Short: "Verify a skycoin address",
		Use:   "verifyAddress [skycoin address]",
		Args:  gcli.ExactArgs(1),
        DisableFlagsInUseLine: true,
        SilenceUsage: true,
		RunE: func(c *gcli.Command, args []string) error {
			_, err := cipher.DecodeBase58Address(args[0])
			return err
		},
	}
}
