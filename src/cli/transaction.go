package cli

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"

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
				return errors.New("invalid txid")
			}

			rpcClient := RPCClientFromContext(c)

			tx, err := rpcClient.GetTransactionByID(txid)
			if err != nil {
				return err
			}

			return printJSON(tx)
		},
	}
}

func decodeRawTxCmd() gcli.Command {
	name := "decodeRawTransaction"
	return gcli.Command{
		Name:         name,
		Usage:        "Decode raw transaction",
		ArgsUsage:    "[raw transaction]",
		OnUsageError: onCommandUsageError(name),
		Action: func(c *gcli.Context) error {
			rawTxStr := c.Args().First()
			if rawTxStr == "" {
				errorWithHelp(c, errors.New("missing raw transaction value"))
				return nil
			}

			b, err := hex.DecodeString(rawTxStr)
			if err != nil {
				fmt.Printf("invalid raw transaction: %v\n", err)
				return err
			}

			tx, err := coin.TransactionDeserialize(b)
			if err != nil {
				fmt.Printf("Unable to deserialize transaction bytes: %v\n", err)
				return err
			}

			txStr, err := visor.TransactionToJSON(tx)
			if err != nil {
				fmt.Println(err)
				return err
			}

			fmt.Println(txStr)
			return nil
		},
	}
}
