package cli

import (
	"strconv"

	gcli "github.com/urfave/cli"
)

func richlistCmd() gcli.Command {
	name := "richlist"
	return gcli.Command{
		Name:      name,
		Usage:     "Display rich list as desc order",
		ArgsUsage: "[topn] [bool (include distribution address or not, default false)]",
		Description: `Display rich list, first argument is topn, second argument is bool(inlcude distribution address or not) 
        example: richlist 100 true`,
		OnUsageError: onCommandUsageError(name),
		Action:       getRichlistCmd,
	}

}

func getRichlistCmd(c *gcli.Context) error {
	var err error
	var isDistribution bool
	var topn int
	topnStr := c.Args().Get(0)
	//return all if no args
	if topnStr == "" {
		isDistribution = true
		topn = -1
	} else {
		topn, err = strconv.Atoi(topnStr)
		if err != nil {
			gcli.ShowSubcommandHelp(c)
			return err
		}
		isDistributionStr := c.Args().Get(1)
		if isDistributionStr == "" {
			isDistribution = false
		} else {
			isDistribution, err = strconv.ParseBool(isDistributionStr)
			if err != nil {
				gcli.ShowSubcommandHelp(c)
				return err
			}
		}
	}

	rpcClient := RpcClientFromContext(c)
	outputs, err := rpcClient.GetRichlist(topn, isDistribution)
	if err != nil {
		return err
	}

	return printJson(outputs)
}
