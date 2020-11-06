package cli

import (
	"github.com/spf13/cobra"
)

// WalletEntry represents an entry in a wallet file
type WalletEntry struct {
	Name       string `json:"name"`
	Label      string `json:"label"`
	AddressNum int    `json:"address_num"`
}

func listWalletsCmd() *cobra.Command {
	return &cobra.Command{
		Short: "Lists all wallets stored in the wallet directory",
		Use:   "listWallets",
		Long: `Lists all wallets stored in the wallet directory.

    The [wallet dir] argument is optional. If not provided, defaults to $DATA_DIR/wallets`,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.MaximumNArgs(0),
		RunE:                  listWallets,
	}
}

func listWallets(_ *cobra.Command, _ []string) error {
	fdn, err := apiClient.WalletFolderName()
	if err != nil {
		return err
	}

	var wlts = struct {
		Directory string        `json:"directory"`
		Wallets   []WalletEntry `json:"wallets"`
	}{
		Directory: fdn.Address,
		Wallets:   []WalletEntry{},
	}

	wltsRsp, err := apiClient.Wallets()
	if err != nil {
		return err
	}

	for _, w := range wltsRsp {
		wlts.Wallets = append(wlts.Wallets, WalletEntry{
			Name:       w.Meta.Filename,
			Label:      w.Meta.Label,
			AddressNum: len(w.Entries),
		})
	}

	return printJSON(wlts)
}
