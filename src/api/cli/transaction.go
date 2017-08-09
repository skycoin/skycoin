package cli

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"

	gcli "github.com/urfave/cli"
)

func transactionCmd() gcli.Command {
	name := "transaction"
	return gcli.Command{
		Name:         name,
		Usage:        "Show detail info of specific transaction",
		ArgsUsage:    "[transaction id]",
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			txid := c.Args().First()
			if txid == "" {
				return errors.New("txid is empty")
			}

			// validate the txid
			_, err := cipher.SHA256FromHex(txid)
			if err != nil {
				return errors.New("error txid")
			}

			rpcClient := c.App.Metadata["rpc"].(*RpcClient)

			tx, err := rpcClient.GetTransactionByID(txid)
			if err != nil {
				return err
			}

			return printJson(tx)
		},
	}
	// Commands = append(Commands, cmd)
}
