package cli

import (
	"github.com/skycoin/skycoin/src/api"
	gcli "github.com/urfave/cli"
)

func richListCmd() gcli.Command {
	name := "richList"
	return gcli.Command{
		Name:  name,
		Usage: "Returns top 20 address balances (based on unspent outputs). Distribution wallets are not currently included.",
		//ArgsUsage:    "[top N wallets (default 20)] [include distribution wallets (true / false) (defalut false)]",
		OnUsageError: onCommandUsageError(name),
		Action:       getRichList,
	}
}

func getRichList(c *gcli.Context) error {
	client := APIClientFromContext(c)

	//num := c.Args().Get(0)
	//if num == "" {
	//	num = "10"
	//}

	//n, err := strconv.ParseInt(num, 10, 32)
	//if err != nil {
	//	return fmt.Errorf("invalid number or top wallets, %s", err)
	//}

	//incDist := c.Args().Get(1)
	//incDistFlag, err2 := strconv.ParseBool(incDist)
	//if err2 != nil {
	//	return fmt.Errorf("invalid boolean value, %s", err2)
	//}

	params := &api.RichlistParams{
		N:                   20,
		IncludeDistribution: false,
	}

	richList, err3 := client.Richlist(params)
	if err3 != nil {
		return err3
	}

	return printJSON(richList)
}
