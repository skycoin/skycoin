package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func broadcastTxCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Broadcast a raw transaction to the network",
		Use:                   "broadcastTransaction [raw transaction]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE: func(_ *cobra.Command, args []string) error {
			rawtx := args[0]

			txid, err := apiClient.InjectEncodedTransaction(rawtx)
			if err != nil {
				return err
			}

			fmt.Println(txid)
			return nil
		},
	}

}
