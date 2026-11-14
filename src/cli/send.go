// Package cli src/cli/send.go c3-coin-cli
package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func sendCmd() *cobra.Command {
	sendCmd := &cobra.Command{
		Args:  cobra.MinimumNArgs(1),
		Short: "Send skycoin from a wallet or an address to a recipient address",
		Use:   "send [wallet] [to address] [amount]",
		Long: `Send skycoin from a wallet or an address to a recipient address.

    Note: the [amount] argument is the coins you will spend, 1 coins = 1e6 droplets.

    The [to address] and [amount] arguments can be replaced with the --many/-m option.

    If you are sending from a wallet without specifying an address,
    the transaction will use one or more of the addresses within the wallet.

    A send that needs more outputs than fit in one transaction is split into several,
    which are all broadcast and print one txid each.

    Use caution when using the “-p” command. If you have command history enabled
    your wallet encryption password can be recovered from the history log.
    If you do not include the “-p” option you will be prompted to enter your password
    after you enter your command.`,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			txns, err := createRawTxnCmdHandler(c, args)
			if err != nil {
				printHelp(c)
				return err
			}

			jsonOutput, err := c.Flags().GetBool("json")
			if err != nil {
				return err
			}

			// The parts spend different outputs, so none waits for another.
			txids := make([]string, 0, len(txns))
			for _, txn := range txns {
				txid, err := apiClient.InjectTransaction(txn)
				if err != nil {
					if len(txids) > 0 {
						return fmt.Errorf("sent %d of %d transactions (%s), then: %w", len(txids), len(txns), strings.Join(txids, ", "), err)
					}
					return err
				}
				txids = append(txids, txid)
			}

			if jsonOutput {
				out := struct {
					Txid  string   `json:"txid"`
					Txids []string `json:"txids,omitempty"`
				}{Txid: txids[0]}
				if len(txids) > 1 {
					out.Txids = txids
				}
				return printJSON(out)
			}

			for _, txid := range txids {
				fmt.Printf("txid:%s\n", txid)
			}
			return nil
		},
	}

	sendCmd.Flags().StringP("from-address", "a", "", "From address in wallet")
	sendCmd.Flags().StringP("change-address", "c", "", `Specify the change address.
Defaults to a new address added to the wallet (an input address for collection wallets).`)
	sendCmd.Flags().StringP("many", "m", "", `use JSON string to set multiple receive addresses and coins,
example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'`)
	sendCmd.Flags().StringP("password", "p", "", "Wallet password")
	sendCmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format.")
	sendCmd.Flags().String("csv", "", "CSV file containing addresses and amounts to send")

	return sendCmd
}
