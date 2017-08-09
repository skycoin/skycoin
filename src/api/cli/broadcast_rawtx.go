package cli

import (
	"fmt"

	"github.com/skycoin/skycoin/src/api/webrpc"
	gcli "github.com/urfave/cli"
)

func broadcastTxCMD() gcli.Command {
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

			txid, err := BroadcastTx(rawtx)
			if err != nil {
				return err
			}

			fmt.Println(txid)
			return nil
		},
	}
	// Commands = append(Commands, cmd)
}

// PUBLIC

// Returns TxId
func BroadcastTx(rawtx string) (string, error) {
	params := []string{rawtx}
	rlt := webrpc.TxIDJson{}

	if err := DoRpcRequest(&rlt, "inject_transaction", params, "1"); err != nil {
		return "", err
	}

	return rlt.Txid, nil
}
