package cli

import (
	"encoding/json"
	"fmt"

	gcli "github.com/urfave/cli"
)

func init() {
	cmd := gcli.Command{
		Name:      "send",
		Usage:     "Send skycoin from a wallet or an address to a recipient address",
		ArgsUsage: "[to address] [amount]",
		Description: `
        If you are sending from a wallet the coins will be taken recursively from all 
        addresses within the wallet starting with the first address until the amount of 
        the transaction is met. 
        
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
			// gcli.StringFlag{
			// 	Name:  "p",
			// 	Usage: "[password] Password for address or wallet.",
			// },
			gcli.BoolFlag{
				Name:  "json,j",
				Usage: "Returns the results in JSON format.",
			},
		},
		Action: func(c *gcli.Context) error {
			rawtx, err := createRawTransaction(c)
			if err != nil {
				return err
			}

			txid, err := broadcastTx(rawtx)
			if err != nil {
				return err
			}

			jsonFmt := c.Bool("json")
			if jsonFmt {
				var rlt = struct {
					Txid string `json:"txid"`
				}{
					txid,
				}
				d, err := json.MarshalIndent(rlt, "", "    ")
				if err != nil {
					return errJSONMarshal
				}
				fmt.Println(string(d))
			} else {
				fmt.Printf("txid:%s\n", txid)
			}

			return nil
		},
	}
	Commands = append(Commands, cmd)
}
