package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"

	"github.com/skycoin/skycoin/src/api/webrpc"
)

func sendCmd() gcli.Command {
	name := "send"
	return gcli.Command{
		Name:      name,
		Usage:     "Send skycoin from a wallet or an address to a recipient address",
		ArgsUsage: "[to address] [amount]",
		Description: `
		Note: the [amount] argument is the coins you will spend, 1 coins = 1e6 droplets.

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
			gcli.StringFlag{
				Name: "m",
				Usage: `[send to many] use JSON string to set multiple recive addresses and coins,
				example: -m '[{"addr":"$addr1", "coins": "10.2"}, {"addr":"$addr2", "coins": "20"}]'`,
			},
			gcli.BoolFlag{
				Name:  "json,j",
				Usage: "Returns the results in JSON format.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			rpcClient := RpcClientFromContext(c)

			rawtx, err := createRawTxCmdHandler(c)
			if err != nil {
				errorWithHelp(c, err)
				return nil
			}

			txid, err := rpcClient.InjectTransaction(rawtx)
			if err != nil {
				return err
			}

			jsonFmt := c.Bool("json")
			if jsonFmt {
				return printJson(struct {
					Txid string `json:"txid"`
				}{
					Txid: txid,
				})
			}

			fmt.Printf("txid:%s\n", txid)
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}

func advancedSendCmd() gcli.Command {
	name := "advancedSend"
	return gcli.Command{
		Name:      name,
		Usage:     "Create a custom send transaction from a wallet or address to a recipient address",
		ArgsUsage: "[to address] [amount] [hours]",
		Description: `
		Note: [amount] argument is the coins you will spend, 1 coins = 1e6 droplets.
		      [hours] argument is number of hours to send to the destination address`,
		Flags: []gcli.Flag{
			gcli.StringFlag{
				Name:  "f",
				Usage: "[wallet file or path] From wallet. If no path is specified your default wallet path will be used.",
			},
			gcli.StringFlag{
				Name:  "a",
				Usage: "[addresses] From address(es)",
			},
			gcli.StringFlag{
				Name:  "c",
				Usage: `[changeAddress] Specify change address, by default the from address with most coins is used`,
			},
			gcli.StringFlag{
				Name: "u",
				Usage: `[unspent address hashes] hash of unspent address(es) to use in the transaction, for multiple separate using a comma
				example: -u "7a1555d60ec1d2a8861376d4b028b321dd9c7d6b438fb628d0ec2e675d91afcb, 4e905cffbbfe0f0f4b7c62bb6c1419ad74fc923eb2d60c239d1f4cf92dce3c5e"`,
			},
			gcli.StringFlag{
				Name: "m",
				Usage: `[send to many] use JSON string to set multiple recive addresses and coins,
				example: -m '[{"addr":"$addr1", "coins": "10.2", "hours": "1"}, {"addr":"$addr2", "coins": "20"}]'`,
			},
			gcli.BoolFlag{
				Name:  "json,j",
				Usage: "Returns the results in JSON format.",
			},
		},
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			rpcClient := RpcClientFromContext(c)

			rawtx, err := createAdvancedRawTxCmdHandler(c)
			if err != nil {
				errorWithHelp(c, err)
				return nil
			}

			txid, err := rpcClient.InjectTransaction(rawtx)
			if err != nil {
				return err
			}

			jsonFmt := c.Bool("json")
			if jsonFmt {
				return printJson(struct {
					Txid string `json:"txid"`
				}{
					Txid: txid,
				})
			}

			fmt.Printf("txid:%s\n", txid)
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}

// @TODO: following functions are not used anywhere, remove them if not required
// SendFromWallet sends from any address or combination of addresses from a wallet. Returns txid.
func SendFromWallet(c *webrpc.Client, walletFile, chgAddr string, toAddrs []SendAmount) (string, error) {
	rawTx, err := CreateRawTxFromWallet(c, walletFile, chgAddr, toAddrs)
	if err != nil {
		return "", err
	}

	return c.InjectTransaction(rawTx)
}

// SendFromAddress sends from a specific address in a wallet. Returns txid.
func SendFromAddress(c *webrpc.Client, addr, walletFile, chgAddr string, toAddrs []SendAmount) (string, error) {
	rawTx, err := CreateRawTxFromAddress(c, addr, walletFile, chgAddr, toAddrs)
	if err != nil {
		return "", err
	}

	return c.InjectTransaction(rawTx)
}
