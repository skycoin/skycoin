package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher/bip32"
)

func verifyXpubCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Verify a bip32 xpub key",
		Use:                   "verifyXpub [xpub key]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE: func(_ *cobra.Command, args []string) error {
			pk, err := bip32.DeserializeEncodedPublicKey(args[0])
			if err != nil {
				return fmt.Errorf("invalid xpub key: %w", err)
			}

			fmt.Printf("Valid xpub key (depth=%d, child=%d)\n", pk.Depth, pk.ChildNumber())
			return nil
		},
	}
}
