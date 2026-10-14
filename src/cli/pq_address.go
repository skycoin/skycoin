package cli

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip39"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/util/droplet"
)

func pqAddressGenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Generate post-quantum (ML-DSA-44) addresses from a seed",
		Use:   "pqAddressGen",
		Long: `Generate post-quantum addresses, which are version 1 addresses that hash an
    ML-DSA-44 public key. Each address has a secret key derived from the seed.

    These addresses only work on a chain past its post-quantum address fork
    (pq_address_fork_seq). Keep the seed and the secret keys secret.`,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, _ []string) error {
			seed, err := c.Flags().GetString("seed")
			if err != nil {
				return err
			}
			num, err := c.Flags().GetUint32("num")
			if err != nil {
				return err
			}
			if seed == "" {
				if seed, err = bip39.NewDefaultMnemonic(); err != nil {
					return err
				}
			}

			type entry struct {
				Address   string `json:"address"`
				SecretKey string `json:"secret_key"`
			}
			out := struct {
				Seed    string  `json:"seed"`
				Entries []entry `json:"entries"`
			}{Seed: seed}
			for i := uint32(0); i < num; i++ {
				sk := pqsig.DeriveSecKey([]byte(seed), i)
				pk, err := sk.PubKey()
				if err != nil {
					return err
				}
				out.Entries = append(out.Entries, entry{
					Address:   cipher.AddressFromPQPubKey(pk).String(),
					SecretKey: sk.Hex(),
				})
			}
			return printJSON(out)
		},
	}
	cmd.Flags().StringP("seed", "s", "", "seed to derive the keys from, a new mnemonic if empty")
	cmd.Flags().Uint32P("num", "n", 1, "number of addresses")
	return cmd
}

func pqSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Send coins from a post-quantum address",
		Use:   "pqSend [secret key] [to address] [amount]",
		Long: `Send coins from the post-quantum address of a secret key (see pqAddressGen).

    Change goes back to the same address. Reusing a post-quantum address is safe,
    because its ML-DSA-44 public key does not give away the secret key.`,
		Args:         cobra.ExactArgs(3),
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			sk, err := pqsig.SecKeyFromHex(args[0])
			if err != nil {
				return err
			}
			to, err := cipher.DecodeBase58Address(args[1])
			if err != nil {
				return err
			}
			coins, err := droplet.FromString(args[2])
			if err != nil {
				return err
			}

			txn, err := createPQTxn(apiClient, sk, to, coins)
			if err != nil {
				return err
			}
			txid, err := apiClient.InjectTransaction(txn)
			if err != nil {
				return err
			}
			fmt.Printf("txid:%s\n", txid)
			return nil
		},
	}
	return cmd
}

// createPQTxn makes and signs a transaction spending the post-quantum address of sk.
func createPQTxn(c GetOutputser, sk pqsig.SecKey, to cipher.Address, coins uint64) (*coin.Transaction, error) {
	pk, err := sk.PubKey()
	if err != nil {
		return nil, err
	}
	from := cipher.AddressFromPQPubKey(pk)

	outputs, err := c.OutputsForAddresses([]string{from.String()})
	if err != nil {
		return nil, err
	}
	uxs, err := outputs.SpendableOutputs().ToUxArray()
	if err != nil {
		return nil, err
	}
	if len(uxs) == 0 {
		return nil, errors.New("the address has no spendable coins")
	}
	head, err := outputs.Head.ToCoinBlockHeader()
	if err != nil {
		return nil, err
	}

	share := decimal.New(5, -1)
	txn, uxb, err := transaction.Create(transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type:        transaction.HoursSelectionTypeAuto,
			Mode:        transaction.HoursSelectionModeShare,
			ShareFactor: &share,
		},
		To:            []coin.TransactionOutput{{Address: to, Coins: coins}},
		ChangeAddress: &from,
	}, coin.NewAddressUxOuts(uxs), head.Time)
	if err != nil {
		return nil, err
	}

	txn.InnerHash = txn.HashInner()
	for i := range uxb {
		if err := txn.SignPQInput(sk, i); err != nil {
			return nil, err
		}
	}
	if err := txn.UpdateHeader(); err != nil {
		return nil, err
	}
	if err := txn.Verify(); err != nil {
		return nil, err
	}
	return txn, nil
}
