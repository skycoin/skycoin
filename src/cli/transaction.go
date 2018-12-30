package cli

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"

	"github.com/spf13/cobra"
)

func transactionCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Show detail info of specific transaction",
		Use:                   "transaction [transaction id]",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
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

func decodeRawTxCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Decode raw transaction",
		Use:                   "decodeRawTransaction [raw transaction]",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
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

func addressTransactionsCmd() cobra.Command {
	return cobra.Command{
		Short: "Show detail for transaction associated with one or more specified addresses",
		Use:   "addressTransactions [address list]",
		Long: `Display transactions for specific addresses, seperate multiple addresses with a space,
        example: addressTransactions addr1 addr2 addr3`,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  getAddressTransactionsCmd,
	}
}

func getAddressTransactionsCmd(c *cobra.Command, args []string) error {
	// Build the list of addresses from the command line arguments
	addrs := make([]string, len(args))
	var err error
	for i := 0; i < len(args); i++ {
		addrs[i] = args[i]
		if _, err = cipher.DecodeBase58Address(addrs[i]); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", addrs[i], err)
		}
	}

	// If one or more addresses have beeb provided, request their transactions - otherwise report an error
	if len(addrs) > 0 {
		outputs, err := apiClient.TransactionsVerbose(addrs)
		if err != nil {
			return err
		}

		return printJSON(outputs)
	}

	return fmt.Errorf("at least one address must be specified. Example: %s addr1 addr2 addr3", c.Name())
}
