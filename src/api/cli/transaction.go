package cli

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"

	"github.com/skycoin/skycoin/src/api/webrpc"
	gcli "github.com/urfave/cli"
)

func transactionCMD() gcli.Command {
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

			tx, err := GetTransactionByID(txid)
			if err != nil {
				return err
			}

			return printJson(tx)
		},
	}
	// Commands = append(Commands, cmd)
}

func GetTransactionByID(txid string) (*webrpc.TxnResult, error) {
	txn := webrpc.TxnResult{}
	if err := DoRpcRequest(&txn, "get_transaction", []string{txid}, "1"); err != nil {
		return nil, err
	}

	return &txn, nil
}
