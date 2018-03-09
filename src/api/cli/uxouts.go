package cli

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	gcli "github.com/urfave/cli"
)

func uxoutsCmd() gcli.Command {
	name := "uxouts"
	return gcli.Command{
		Name:         name,
		Usage:        "Display uxouts of uxids",
		ArgsUsage:    "[uxid array]",
		Description:  "Display uxouts of uxids",
		OnUsageError: onCommandUsageError(name),
		Action:       getUxoutsCmd,
	}
}

func getUxoutsCmd(c *gcli.Context) error {
	rpcClient := RpcClientFromContext(c)

	var uxids []string
	var err error
	for i := 0; i < c.NArg(); i++ {
		uxid := c.Args().Get(i)
		_, err := cipher.SHA256FromHex(uxid)
		if err != nil {
			return fmt.Errorf("invalid uxid %v: %v", uxid, err)
		}

		uxids = append(uxids, uxid)
	}

	outputs, err := rpcClient.GetUxouts(uxids)
	if err != nil {
		return err
	}

	return printJson(outputs)
}
