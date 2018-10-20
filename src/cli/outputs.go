package cli

import (
	"fmt"

	gcli "github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/wallet"
)

func walletOutputsCmd() *gcli.Command {
	return &gcli.Command{
		Short: "Display outputs of specific wallet",
		Use:   "walletOutputs [wallet file]",
		Long: fmt.Sprintf(`Display outputs of specific wallet, the default wallet (%s) will be
    used if no wallet was specified, use ENV 'WALLET_NAME'
    to update default wallet file name, and 'WALLET_DIR' to update
    the default wallet directory`, cliConfig.FullWalletPath()),
		RunE: getWalletOutputsCmd,
		Args: gcli.MaximumNArgs(1),
	}
}

func addressOutputsCmd() *gcli.Command {
	return &gcli.Command{
		Short: "Display outputs of specific addresses",
		Use:   "addressOutputs [address list]",
		Long: `Display outputs of specific addresses, join multiple addresses with space,
       example: addressOutputs $addr1 $addr2 $addr3`,
		RunE: getAddressOutputsCmd,
		Args: gcli.MinimumNArgs(1),
	}
}

func getWalletOutputsCmd(c *gcli.Command, args []string) error {
	var err error
	w, err := resolveWalletPath(cliConfig, args[0])
	if err != nil {
		return err
	}

	outputs, err := GetWalletOutputsFromFile(apiClient, w)
	if err != nil {
		return err
	}

	return printJSON(webrpc.OutputsResult{
		Outputs: *outputs,
	})
}

func getAddressOutputsCmd(c *gcli.Command, args []string) error {
	addrs := make([]string, len(args))

	var err error
	for i := 0; i < len(args); i++ {
		addrs[i] = args[i]
		if _, err = cipher.DecodeBase58Address(addrs[i]); err != nil {
			return fmt.Errorf("invalid address: %v, err: %v", addrs[i], err)
		}
	}

	outputs, err := apiClient.OutputsForAddresses(addrs)
	if err != nil {
		return err
	}

	return printJSON(webrpc.OutputsResult{
		Outputs: *outputs,
	})
}

// PUBLIC

// GetWalletOutputsFromFile returns unspent outputs associated with all addresses in a wallet file
func GetWalletOutputsFromFile(c GetOutputser, walletFile string) (*readable.UnspentOutputsSummary, error) {
	wlt, err := wallet.Load(walletFile)
	if err != nil {
		return nil, err
	}

	return GetWalletOutputs(c, wlt)
}

// GetWalletOutputs returns unspent outputs associated with all addresses in a wallet.Wallet
func GetWalletOutputs(c GetOutputser, wlt *wallet.Wallet) (*readable.UnspentOutputsSummary, error) {
	cipherAddrs := wlt.GetAddresses()
	addrs := make([]string, len(cipherAddrs))
	for i := range cipherAddrs {
		addrs[i] = cipherAddrs[i].String()
	}

	return c.OutputsForAddresses(addrs)
}
