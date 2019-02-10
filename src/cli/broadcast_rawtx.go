package cli

import (
	"fmt"

	gcli "github.com/spf13/cobra"
)

func broadcastTxCmd() *gcli.Command {
	return &gcli.Command{
		Short:                 "Broadcast a raw transaction to the network",
		Use:                   "broadcastTransaction [raw transaction]",
		Args:                  gcli.ExactArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE: func(_ *gcli.Command, args []string) error {
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
