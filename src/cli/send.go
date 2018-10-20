package cli

import (
	"fmt"

	gcli "github.com/spf13/cobra"
)

func sendCmd() *gcli.Command {
	sendCmd := &gcli.Command{
		Short: "Send skycoin from a wallet or an address to a recipient address",
		Use:   "send [to address] [amount]",
		Long: `Note: the [amount] argument is the coins you will spend, 1 coins = 1e6 droplets.

    If you are sending from a wallet without specifying an address,
    the transaction will use one or more of the addresses within the wallet.

    Use caution when using the “-p” command. If you have command history enabled
    your wallet encryption password can be recovered from the history log.
    If you do not include the “-p” option you will be prompted to enter your password
    after you enter your command.`,
		RunE: func(c *gcli.Command, args []string) error {
			rawTxn, err := createRawTxnCmdHandler(args)
			if err != nil {
				printHelp(c)
				return err
			}

			txid, err := apiClient.InjectTransaction(rawTxn)
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

	sendCmd.Flags().StringVarP(&walletFile, "wallet-file", "f", "", "wallet file or path. If no path is specified your default wallet path will be used.")
	sendCmd.Flags().StringVarP(&address, "address", "a", "", "From address")
	sendCmd.Flags().StringVarP(&changeAddress, "change-address", "c", "", `Specify different change address.
By default the from address or a wallets coinbase address will be used.`)
	sendCmd.Flags().StringVarP(&many, "many", "m", "", `use JSON string to set multiple receive addresses and coins,
example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'`)
	sendCmd.Flags().StringVarP(&password, "password", "p", "", "Wallet password")
	sendCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Returns the results in JSON format.")
	sendCmd.Flags().StringVar(&csvFile, "csv-file", "", "CSV file containing addresses and amounts to send")

	return sendCmd
}
