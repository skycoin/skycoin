package cli

import (
	"fmt"
	"strconv"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
	gcli "github.com/urfave/cli"
)

func walletOutputsCmd(cfg Config) gcli.Command {
	name := "walletOutputs"
	return gcli.Command{
		Name:      name,
		Usage:     "Display outputs of specific wallet",
		ArgsUsage: "[wallet file]",
		Description: fmt.Sprintf(`Display outputs of specific wallet, the default
		wallet (%s) will be
		used if no wallet was specified, use ENV 'WALLET_NAME'
		to update default wallet file name, and 'WALLET_DIR' to update
		the default wallet directory`, cfg.FullWalletPath()),
		OnUsageError: onCommandUsageError(name),
		Action:       getWalletOutputsCmd,
	}
}

func addressOutputsCmd() gcli.Command {
	name := "addressOutputs"
	return gcli.Command{
		Name:      name,
		Usage:     "Display outputs of specific addresses",
		ArgsUsage: "[address list]",
		Description: `Display outputs of specific addresses, join multiple addresses with space,
        example: addressOutputs $addr1 $addr2 $addr3`,
		OnUsageError: onCommandUsageError(name),
		Action:       getAddressOutputsCmd,
	}

}

func richlistCmd() gcli.Command {
	name := "richlist"
	return gcli.Command{
		Name:      name,
		Usage:     "Display rich list as desc order",
		ArgsUsage: "[topn] [bool (include distribution address or not, default false)]",
		Description: `Display rich list, first argument is topn, second argument is bool(inlcude distribution address or not) 
        example: richlist 100 true`,
		OnUsageError: onCommandUsageError(name),
		Action:       getTopnOutputsCmd,
	}

}

func getWalletOutputsCmd(c *gcli.Context) error {
	cfg := ConfigFromContext(c)
	rpcClient := RpcClientFromContext(c)

	w := ""
	if c.NArg() > 0 {
		w = c.Args().First()
	}

	var err error
	w, err = resolveWalletPath(cfg, w)
	if err != nil {
		return err
	}

	outputs, err := GetWalletOutputsFromFile(rpcClient, w)
	if err != nil {
		return err
	}

	return printJson(outputs)
}

func getAddressOutputsCmd(c *gcli.Context) error {
	rpcClient := RpcClientFromContext(c)

	addrs := make([]string, c.NArg())
	var err error
	for i := 0; i < c.NArg(); i++ {
		addrs[i] = c.Args().Get(i)
		if _, err = cipher.DecodeBase58Address(addrs[i]); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", addrs[i], err)
		}
	}

	outputs, err := rpcClient.GetUnspentOutputs(addrs)
	if err != nil {
		return err
	}

	return printJson(outputs)
}

// PUBLIC

func GetWalletOutputsFromFile(c *webrpc.Client, walletFile string) (*webrpc.OutputsResult, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, err
	}

	return GetWalletOutputs(c, wlt)
}

func GetWalletOutputs(c *webrpc.Client, wlt *wallet.Wallet) (*webrpc.OutputsResult, error) {
	cipherAddrs := wlt.GetAddresses()
	addrs := make([]string, len(cipherAddrs))
	for i := range cipherAddrs {
		addrs[i] = cipherAddrs[i].String()
	}

	return c.GetUnspentOutputs(addrs)
}

func getTopnOutputsCmd(c *gcli.Context) error {
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
	outputs, err := rpcClient.GetTopnOutputs(topn, isDistribution)
	if err != nil {
		return err
	}

	return printJson(outputs)
}
