package cli

import (
	"fmt"

	gcli "github.com/spf13/cobra"

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
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  getWalletOutputsCmd,
	}
}

func addressOutputsCmd() *gcli.Command {
	return &gcli.Command{
		Short: "Display outputs of specific addresses",
		Use:   "addressOutputs [address list]",
		Long: `Display outputs of specific addresses, join multiple addresses with space,
    example: addressOutputs $addr1 $addr2 $addr3`,
		Args:                  gcli.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  getAddressOutputsCmd,
	}
}

// OutputsResult the output json format
type OutputsResult struct {
	Outputs readable.UnspentOutputsSummary `json:"outputs"`
}

func getWalletOutputsCmd(_ *gcli.Command, args []string) error {
	var wltPath string
	if len(args) == 1 {
		wltPath = args[0]
	}
	var err error
	w, err := resolveWalletPath(cliConfig, wltPath)
	if err != nil {
		return err
	}

	outputs, err := GetWalletOutputsFromFile(apiClient, w)
	if err != nil {
		return err
	}

	return printJSON(OutputsResult{
		Outputs: *outputs,
	})
}

func getAddressOutputsCmd(_ *gcli.Command, args []string) error {
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

	return printJSON(OutputsResult{
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
func GetWalletOutputs(c GetOutputser, wlt wallet.Wallet) (*readable.UnspentOutputsSummary, error) {
	cipherAddrs := wlt.GetAddresses()
	addrs := make([]string, len(cipherAddrs))
	for i := range cipherAddrs {
		addrs[i] = cipherAddrs[i].String()
	}

	return c.OutputsForAddresses(addrs)
}
