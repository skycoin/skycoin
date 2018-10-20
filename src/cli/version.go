package cli

import (
	"fmt"
	"reflect"

	gcli "github.com/spf13/cobra"
)

func versionCmd() *gcli.Command {
	versionCmd := &gcli.Command{
		Use:   "version",
		Short: "List the current version of Skycoin components",
		RunE: func(c *gcli.Command, args []string) error {
			var ver = struct {
				Skycoin string `json:"skycoin"`
				Cli     string `json:"cli"`
				RPC     string `json:"rpc"`
				Wallet  string `json:"wallet"`
			}{
				Version,
				Version,
				Version,
				Version,
			}

			if jsonOutput {
				return printJSON(ver)
			}

			v := reflect.ValueOf(ver)
			t := reflect.TypeOf(ver)
			for i := 0; i < v.NumField(); i++ {
				fmt.Printf("%s:%v\n", t.Field(i).Tag.Get("json"), v.Field(i).Interface())
			}

			return nil
		},
	}

	versionCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Returns the results in JSON format")

	return versionCmd
}
