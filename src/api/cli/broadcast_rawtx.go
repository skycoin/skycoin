package cli

import (
	"fmt"

	gcli "github.com/urfave/cli"
)

func broadcastTxCmd() gcli.Command {
	name := "broadcastTransaction"
	return gcli.Command{
		Name:         name,
		Usage:        "Broadcast a raw transaction to the network",
		ArgsUsage:    "[raw transaction]",
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			rawtx := c.Args().First()
			if rawtx == "" {
				gcli.ShowSubcommandHelp(c)
				return nil
			}

			rpcClient := RpcClientFromContext(c)
			txid, err := rpcClient.InjectTransaction(rawtx)
			if err != nil {
				return err
			}

			fmt.Println(txid)
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}
