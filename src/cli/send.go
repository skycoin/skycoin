package cli

import (
	"fmt"

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

    Use caution when using the “-p” command. If you have command history enabled
    your wallet encryption password can be recovered from the history log.
    If you do not include the “-p” option you will be prompted to enter your password
    after you enter your command.`,
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			rawTxn, err := createRawTxnCmdHandler(c, args)
			if err != nil {
				printHelp(c)
				return err
			}

			txid, err := apiClient.InjectTransaction(rawTxn)
			if err != nil {
				return err
			}

			jsonOutput, err := c.Flags().GetBool("json")
			if err != nil {
				return err
			}
			if jsonOutput {
				return printJSON(struct {
					Txid string `json:"txid"`
				}{
					Txid: txid,
				})
			}

			fmt.Printf("txid:%s\n", txid)
			return nil
		},
	}

	sendCmd.Flags().StringP("from-address", "a", "", "From address in wallet")
	sendCmd.Flags().StringP("change-address", "c", "", `Specify the change address.
Defaults to one of the spending addresses (deterministic wallets) or to a new change address (bip44 wallets).`)
	sendCmd.Flags().StringP("many", "m", "", `use JSON string to set multiple receive addresses and coins,
example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'`)
	sendCmd.Flags().StringP("password", "p", "", "Wallet password")
	sendCmd.Flags().BoolP("json", "j", false, "Returns the results in JSON format.")
	sendCmd.Flags().String("csv", "", "CSV file containing addresses and amounts to send")

	return sendCmd
}
