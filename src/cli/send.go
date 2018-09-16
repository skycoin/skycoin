package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"
)

func sendCmd() gcli.Command {
	name := "send"
	return gcli.Command{
		Name:      name,
		Usage:     "Send skycoin from a wallet or an address to a recipient address",
		ArgsUsage: "[to address] [amount]",
		Description: `
		Note: the [amount] argument is the coins you will spend, 1 coins = 1e6 droplets.

        If you are sending from a wallet without specifying an address,
        the transaction will use one or more of the addresses within the wallet.

        Use caution when using the “-p” command. If you have command history enabled
        your wallet encryption password can be recovered from the history log.
        If you do not include the “-p” option you will be prompted to enter your password
        after you enter your command.`,
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "f",
				Usage: "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.",
			},
			gcli.StringFlag{
				Name:  "a",
				Usage: "[address] From address",
			},
			gcli.StringFlag{
				Name: "c",
				Usage: `[changeAddress] Specify change address, by default the from address or
				the wallet's coinbase address will be used`,
			},
			gcli.StringFlag{
				Name:  "p",
				Usage: "[password] Wallet password",
			},
			gcli.StringFlag{
				Name: "m",
				Usage: `[send to many] use JSON string to set multiple recive addresses and coins,
				example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'`,
			},
			gcli.BoolFlag{
				Name:  "json,j",
				Usage: "Returns the results in JSON format.",
			},
			gcli.StringFlag{
				Name:  "csv",
				Usage: "[filepath] CSV file containing addresses and amounts to send",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			apiClient := APIClientFromContext(c)

			rawTxn, err := createRawTxnCmdHandler(c)
			if err != nil {
				printHelp(c)
				return err
			}

			txid, err := apiClient.InjectTransaction(rawTxn)
			if err != nil {
				return err
			}

			jsonFmt := c.Bool("json")
			if jsonFmt {
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
}
