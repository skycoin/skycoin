package cli

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"

	gcli "github.com/spf13/cobra"
)

func transactionCmd() *gcli.Command {
	return &gcli.Command{
		Short: "Show detail info of specific transaction",
		Use:   "transaction [transaction id]",
        DisableFlagsInUseLine: true,
        SilenceUsage: true,
		Args:  gcli.MaximumNArgs(1),
		RunE: func(c *gcli.Command, args []string) error {
			txid := args[0]
			if txid == "" {
				return errors.New("txid is empty")
			}

			// validate the txid
			_, err := cipher.SHA256FromHex(txid)
			if err != nil {
				return errors.New("invalid txid")
			}

			txn, err := apiClient.Transaction(txid)
			if err != nil {
				return err
			}

			return printJSON(webrpc.TxnResult{
				Transaction: txn,
			})
		},
	}
}

func decodeRawTxCmd() *gcli.Command {
	return &gcli.Command{
		Short: "Decode raw transaction",
		Use:   "decodeRawTransaction [raw transaction]",
        DisableFlagsInUseLine: true,
        SilenceUsage: true,
		Args:  gcli.ExactArgs(1),
		RunE: func(c *gcli.Command, args []string) error {
			b, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("invalid raw transaction: %v", err)
			}

			txn, err := coin.TransactionDeserialize(b)
			if err != nil {
				return fmt.Errorf("Unable to deserialize transaction bytes: %v", err)
			}

			// Assume the transaction is not malformed and if it has no inputs
			// that it is the genesis block's transaction
			isGenesis := len(txn.In) == 0
			rTxn, err := readable.NewTransaction(txn, isGenesis)
			if err != nil {
				return err
			}

			return printJSON(rTxn)
		},
	}
}
