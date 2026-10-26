package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/wallet"
)

func sweepCmd() *cobra.Command {
	sweepCmd := &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Short: "Move every coin in a wallet to one new address",
		Use:   "sweep [wallet]",
		Long: `Move every spendable coin in a wallet to one new address in the same wallet.

    An address that has spent has its public key on the blockchain, and the
    coins it still holds are exposed to anyone who can derive the private key
    from that key. Sweeping moves them to an address whose key has never been
    revealed. Spend from the new address afterwards, not the old ones.

    --to sends the coins to another address instead, such as one in a new
    wallet. Collection wallets cannot generate an address and need --to.

    Use caution when using the "-p" option. If you have command history enabled
    your wallet encryption password can be recovered from the history log. If you
    do not include the "-p" option you will be prompted to enter your password
    after you enter your command.`,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			to, err := c.Flags().GetString("to")
			if err != nil {
				return err
			}
			password, err := c.Flags().GetString("password")
			if err != nil {
				return err
			}
			rawOnly, err := c.Flags().GetBool("raw")
			if err != nil {
				return err
			}
			jsonOutput, err := c.Flags().GetBool("json")
			if err != nil {
				return err
			}

			txn, dest, err := sweepWallet(apiClient, args[0], to, NewPasswordReader([]byte(password)), params.MainNetDistribution)
			if err != nil {
				return err
			}

			rawTxn, err := txn.SerializeHex()
			if err != nil {
				return err
			}
			if rawOnly {
				if jsonOutput {
					return printJSON(struct {
						RawTx   string `json:"rawtx"`
						Address string `json:"address"`
					}{rawTxn, dest})
				}
				fmt.Println(rawTxn)
				return nil
			}

			txid, err := apiClient.InjectTransaction(txn)
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(struct {
					Txid    string `json:"txid"`
					Address string `json:"address"`
				}{txid, dest})
			}
			fmt.Printf("txid:%s\naddress:%s\n", txid, dest)
			return nil
		},
	}

	sweepCmd.Flags().String("to", "", "Send the coins to this address instead of a new one in the wallet")
	sweepCmd.Flags().StringP("password", "p", "", "Wallet password")
	sweepCmd.Flags().Bool("raw", false, "Print the raw transaction instead of broadcasting it")
	sweepCmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format.")

	return sweepCmd
}

// sweepWallet makes a transaction moving every spendable coin in walletFile to one
// address: to, or else a new address it adds to the wallet file.
func sweepWallet(c GetOutputser, walletFile, to string, pr PasswordReader, distParams params.Distribution) (*coin.Transaction, string, error) {
	if to != "" {
		if _, err := cipher.DecodeBase58Address(to); err != nil {
			return nil, "", ErrAddress
		}
	}

	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, "", WalletLoadError{err}
	}
	if to == "" && wlt.Type() == wallet.WalletTypeCollection {
		return nil, "", errors.New("a collection wallet cannot generate an address, use --to")
	}

	var password []byte
	if wlt.IsEncrypted() {
		if password, err = pr.Password(); err != nil {
			return nil, "", err
		}
	}

	addrs, err := wlt.GetAddresses()
	if err != nil {
		return nil, "", err
	}
	inAddrs := make([]string, len(addrs))
	for i, a := range addrs {
		inAddrs[i] = a.String()
	}

	outputs, err := c.OutputsForAddresses(inAddrs)
	if err != nil {
		return nil, "", err
	}
	uxs, err := outputs.SpendableOutputs().ToUxArray()
	if err != nil {
		return nil, "", err
	}
	total, err := uxs.Coins()
	if err != nil {
		return nil, "", err
	}
	if total == 0 {
		return nil, "", errors.New("the wallet has no spendable coins")
	}

	added := false
	if to == "" {
		if to, added, err = newChangeAddress(wlt, password, inAddrs); err != nil {
			return nil, "", err
		}
	}

	txn, err := CreateRawTxn(c, wlt, inAddrs, to, []SendAmount{{Addr: to, Coins: total}}, password, distParams)
	if err != nil {
		return nil, "", err
	}
	if added {
		if err := saveWalletFile(wlt, walletFile); err != nil {
			return nil, "", err
		}
	}
	return txn, to, nil
}
