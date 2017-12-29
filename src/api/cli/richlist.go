package cli

import (
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
		Flags: []gcli.Flag{
			gcli.IntFlag{
				Name:  "topn,n",
				Usage: "Returns richlist top number, by default returns all",
			},
			gcli.BoolFlag{
				Name:  "include-distribution,d",
				Usage: "Include distribution address or not, default false",
			},
		},
		Action: getRichlistCmd,
	}

}

func getRichlistCmd(c *gcli.Context) error {
	topn := c.Int("n")
	if topn == 0 {
		topn = -1
	}
	isDistribution := c.Bool("d")
	rpcClient := RpcClientFromContext(c)
	outputs, err := rpcClient.GetRichlist(topn, isDistribution)
	if err != nil {
		return err
	}

	return printJson(outputs)
}
