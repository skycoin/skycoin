package cli

import (
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher/pqsig"
)

func pqKeyGenCmd() *cobra.Command {
	return &cobra.Command{
		Short: "Generate a post-quantum (ML-DSA-44) block publisher key pair",
		Use:   "pqKeyGen",
		Long: `Generate a post-quantum (ML-DSA-44) block publisher key pair.

    Nodes take the public key as blockchain_pq_pubkey_str in fiber.toml, and from
    pq_fork_seq they require every block to carry a signature by it. The block
    publisher takes the secret key as -blockchain-pq-secret-key. Keep it secret.`,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			sk, err := pqsig.NewSecKey()
			if err != nil {
				return err
			}
			pk, err := sk.PubKey()
			if err != nil {
				return err
			}
			return printJSON(struct {
				SecretKey string `json:"secret_key"`
				PublicKey string `json:"public_key"`
			}{sk.Hex(), pk.Hex()})
		},
	}
}
