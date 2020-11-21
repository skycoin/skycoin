package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/wallet"
)

func walletOutputsCmd() *cobra.Command {
	return &cobra.Command{
		Short:                 "Display outputs of specific wallet",
		Use:                   "walletOutputs [wallet]",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.ExactArgs(1),
		RunE:                  getWalletOutputsCmd,
	}
}

func addressOutputsCmd() *cobra.Command {
	return &cobra.Command{
		Short: "Display outputs of specific addresses",
		Use:   "addressOutputs [address list]",
		Long: `Display outputs of specific addresses, join multiple addresses with space,
    example: addressOutputs $addr1 $addr2 $addr3`,
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  getAddressOutputsCmd,
	}
}

// OutputsResult the output json format
type OutputsResult struct {
	Outputs readable.UnspentOutputsSummary `json:"outputs"`
}

func getWalletOutputsCmd(_ *cobra.Command, args []string) error {
	outputs, err := GetWalletOutputsFromFile(apiClient, args[0])
	if err != nil {
		return err
	}

	return printJSON(OutputsResult{
		Outputs: *outputs,
	})
}

func getAddressOutputsCmd(_ *cobra.Command, args []string) error {
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
