package cli

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/spf13/cobra"
)

func generateKeyPairCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Generates a new key pair (Public and Private) as Hex strings.",
		Long:                  "Generates a new key pair (Public and Private) based on the Skycoin crypto library and returns them as Hex strings. If these are intended for actual use, ensure they are kept secure.",
		Use:                   "generateKeyPair",
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  generateKeyPair,
	}
}

func generateKeyPair(_ *cobra.Command, _ []string) error {
	p, s := cipher.GenerateKeyPair()
	a, err := cipher.AddressFromSecKey(s)
	if err != nil {
		return err
	}

	if a.Verify(p) != nil {
		return err
	}

	fmt.Println("New key pair generated. If you intend to use these, ensure they are kept secure.")
	fmt.Printf("Public:  %s\n", p.Hex())
	fmt.Printf("Private: %s\n", s.Hex())
	return nil
}
